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

// fakePublisher records domain events for assertion.
type fakePublisher struct {
	events []domain.DomainEvent
}

func (p *fakePublisher) Publish(_ context.Context, events ...domain.DomainEvent) error {
	p.events = append(p.events, events...)
	return nil
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

func TestIncrementStatField_bumpsCounter_savesAndPublishesEvent(t *testing.T) {
	repo := newFakeStatsRepo()
	pub := &fakePublisher{}
	uc := application.IncrementStatField{Repo: repo, Events: pub}

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
	if len(pub.events) != 1 {
		t.Fatalf("published events = %d, want 1", len(pub.events))
	}
	ev, ok := pub.events[0].(domain.StatIncremented)
	if !ok {
		t.Fatalf("event type = %T, want domain.StatIncremented", pub.events[0])
	}
	if ev.Field != domain.StatFieldProjects {
		t.Errorf("event.Field = %v, want StatFieldProjects", ev.Field)
	}
	if ev.NewValue != 8 {
		t.Errorf("event.NewValue = %d, want 8", ev.NewValue)
	}
	if ev.Delta != 5 {
		t.Errorf("event.Delta = %d, want 5", ev.Delta)
	}
}

func TestIncrementStatField_negativeClampedAtZero_eventReflectsClamp(t *testing.T) {
	repo := newFakeStatsRepo()
	pub := &fakePublisher{}
	uc := application.IncrementStatField{Repo: repo, Events: pub}

	uid, _ := domain.NewUserID("user-3")
	got, err := uc.Execute(context.Background(), uid, domain.StatFieldNotifications, -100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Notifications != 0 {
		t.Errorf("Notifications = %d, want 0 (clamped)", got.Notifications)
	}
	// The event records the actual post-clamp delta (-2, not -100).
	ev := pub.events[0].(domain.StatIncremented)
	if ev.Delta != -2 {
		t.Errorf("event.Delta = %d, want -2 (post-clamp from seed 2)", ev.Delta)
	}
	if ev.NewValue != 0 {
		t.Errorf("event.NewValue = %d, want 0", ev.NewValue)
	}
}

func TestIncrementStatField_repoSaveFailurePropagated_noPublish(t *testing.T) {
	repo := newFakeStatsRepo()
	repo.failSave = true
	pub := &fakePublisher{}
	uc := application.IncrementStatField{Repo: repo, Events: pub}

	uid, _ := domain.NewUserID("user-4")
	if _, err := uc.Execute(context.Background(), uid, domain.StatFieldActivity, 1); err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(pub.events) != 0 {
		t.Errorf("publisher was called despite save failure: %+v", pub.events)
	}
}
