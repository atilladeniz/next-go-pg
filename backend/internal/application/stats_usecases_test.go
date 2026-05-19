package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
)

// fakeStatsRepo is an in-memory StatsRepository for testing use cases
// without booting GORM or Postgres.
type fakeStatsRepo struct {
	store     map[domain.UserID]*domain.UserStats
	getCalls  int
	saveCalls int
	failGet   bool
	failSave  bool
}

func newFakeStatsRepo() *fakeStatsRepo {
	return &fakeStatsRepo{store: map[domain.UserID]*domain.UserStats{}}
}

func (r *fakeStatsRepo) GetOrCreate(_ context.Context, userID domain.UserID) (*domain.UserStats, error) {
	r.getCalls++
	if r.failGet {
		return nil, errors.New("get failed")
	}
	if s, ok := r.store[userID]; ok {
		return s, nil
	}
	s := &domain.UserStats{UserID: userID, ProjectCount: 3, ActivityToday: 10, Notifications: 2}
	r.store[userID] = s
	return s, nil
}

func (r *fakeStatsRepo) Save(_ context.Context, stats *domain.UserStats) error {
	r.saveCalls++
	if r.failSave {
		return errors.New("save failed")
	}
	r.store[stats.UserID] = stats
	return nil
}

// fakeBroadcaster records broadcasted events for assertion.
type fakeBroadcaster struct {
	events []event
}

type event struct {
	name    string
	payload string
}

func (b *fakeBroadcaster) Broadcast(eventName, payload string) {
	b.events = append(b.events, event{name: eventName, payload: payload})
}

func TestGetUserStats_returnsStoreEntry(t *testing.T) {
	repo := newFakeStatsRepo()
	uc := application.GetUserStats{Repo: repo}

	uid, _ := domain.NewUserID("user-1")
	stats, err := uc.Execute(context.Background(), uid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.UserID != uid {
		t.Errorf("UserID = %q, want %q", stats.UserID, uid)
	}
	if stats.ProjectCount != 3 {
		t.Errorf("ProjectCount = %d, want 3 (seeded default)", stats.ProjectCount)
	}
	if repo.getCalls != 1 {
		t.Errorf("repo.GetOrCreate called %d times, want 1", repo.getCalls)
	}
}

func TestIncrementStatField_bumpsCounter_savesAndBroadcasts(t *testing.T) {
	repo := newFakeStatsRepo()
	broker := &fakeBroadcaster{}
	uc := application.IncrementStatField{Repo: repo, Events: broker}

	uid, _ := domain.NewUserID("user-2")
	got, err := uc.Execute(context.Background(), uid, domain.StatFieldProjects, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ProjectCount != 8 { // 3 (seed) + 5
		t.Errorf("ProjectCount = %d, want 8", got.ProjectCount)
	}
	if repo.saveCalls != 1 {
		t.Errorf("repo.Save called %d times, want 1", repo.saveCalls)
	}
	if len(broker.events) != 1 {
		t.Fatalf("broadcast count = %d, want 1", len(broker.events))
	}
	if broker.events[0].name != "stats-updated" {
		t.Errorf("event name = %q, want stats-updated", broker.events[0].name)
	}
	if broker.events[0].payload != `{"field":"projects"}` {
		t.Errorf("payload = %q, unexpected", broker.events[0].payload)
	}
}

func TestIncrementStatField_negativeClampedAtZero(t *testing.T) {
	repo := newFakeStatsRepo()
	uc := application.IncrementStatField{Repo: repo, Events: &fakeBroadcaster{}}

	uid, _ := domain.NewUserID("user-3")
	got, err := uc.Execute(context.Background(), uid, domain.StatFieldNotifications, -100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Notifications != 0 {
		t.Errorf("Notifications = %d, want 0 (clamped)", got.Notifications)
	}
}

func TestIncrementStatField_repoSaveFailurePropagated_noBroadcast(t *testing.T) {
	repo := newFakeStatsRepo()
	repo.failSave = true
	broker := &fakeBroadcaster{}
	uc := application.IncrementStatField{Repo: repo, Events: broker}

	uid, _ := domain.NewUserID("user-4")
	if _, err := uc.Execute(context.Background(), uid, domain.StatFieldActivity, 1); err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(broker.events) != 0 {
		t.Errorf("broadcaster was called despite save failure: %+v", broker.events)
	}
}
