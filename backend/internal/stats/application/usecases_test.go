package application_test

import (
	"context"
	"errors"
	"testing"

	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
	statsapp "github.com/atilladeniz/next-go-pg/backend/internal/stats/application"
	stats "github.com/atilladeniz/next-go-pg/backend/internal/stats/domain"
)

// fakeRepo is an in-memory stats Repository for testing.
type fakeRepo struct {
	store     map[shared.UserID]*stats.UserStats
	getCalls  int
	saveCalls int
	failGet   bool
	failSave  bool
}

func newFakeRepo() *fakeRepo { return &fakeRepo{store: map[shared.UserID]*stats.UserStats{}} }

func (r *fakeRepo) GetOrCreate(_ context.Context, userID shared.UserID) (*stats.UserStats, error) {
	r.getCalls++
	if r.failGet {
		return nil, errors.New("get failed")
	}
	if s, ok := r.store[userID]; ok {
		return s, nil
	}
	s := stats.NewUserStats(userID)
	r.store[userID] = s
	return s, nil
}

func (r *fakeRepo) Save(_ context.Context, agg *stats.UserStats) error {
	r.saveCalls++
	if r.failSave {
		return errors.New("save failed")
	}
	r.store[agg.UserID] = agg
	return nil
}

// fakePublisher records dispatched domain events.
type fakePublisher struct {
	events []shared.DomainEvent
}

func (p *fakePublisher) Publish(_ context.Context, events ...shared.DomainEvent) error {
	p.events = append(p.events, events...)
	return nil
}

func TestGetUserStats_returnsStoreEntry(t *testing.T) {
	repo := newFakeRepo()
	uc := statsapp.GetUserStats{Repo: repo}

	uid, _ := shared.NewUserID("user-1")
	got, err := uc.Execute(context.Background(), uid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserID != uid {
		t.Errorf("UserID = %q, want %q", got.UserID, uid)
	}
	if got.ProjectCount != 3 {
		t.Errorf("ProjectCount = %d, want 3 (seeded default)", got.ProjectCount)
	}
}

func TestIncrementStatField_bumpsCounter_savesAndPublishes(t *testing.T) {
	repo := newFakeRepo()
	pub := &fakePublisher{}
	uc := statsapp.IncrementStatField{Repo: repo, Events: pub}

	uid, _ := shared.NewUserID("user-2")
	got, err := uc.Execute(context.Background(), uid, stats.StatFieldProjects, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ProjectCount != 8 {
		t.Errorf("ProjectCount = %d, want 8", got.ProjectCount)
	}
	if len(pub.events) != 1 {
		t.Fatalf("published events = %d, want 1", len(pub.events))
	}
	ev, ok := pub.events[0].(stats.StatIncremented)
	if !ok {
		t.Fatalf("event type = %T, want stats.StatIncremented", pub.events[0])
	}
	if ev.Field != stats.StatFieldProjects {
		t.Errorf("event.Field = %v, want StatFieldProjects", ev.Field)
	}
	if ev.NewValue != 8 || ev.Delta != 5 {
		t.Errorf("event.NewValue=%d Delta=%d, want 8/5", ev.NewValue, ev.Delta)
	}
}

func TestIncrementStatField_negativeClampedAtZero_eventReflectsClamp(t *testing.T) {
	repo := newFakeRepo()
	pub := &fakePublisher{}
	uc := statsapp.IncrementStatField{Repo: repo, Events: pub}

	uid, _ := shared.NewUserID("user-3")
	got, err := uc.Execute(context.Background(), uid, stats.StatFieldNotifications, -100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Notifications != 0 {
		t.Errorf("Notifications = %d, want 0 (clamped)", got.Notifications)
	}
	ev := pub.events[0].(stats.StatIncremented)
	if ev.Delta != -2 {
		t.Errorf("event.Delta = %d, want -2 (post-clamp from seed 2)", ev.Delta)
	}
	if ev.NewValue != 0 {
		t.Errorf("event.NewValue = %d, want 0", ev.NewValue)
	}
}

func TestIncrementStatField_saveFailure_noPublish(t *testing.T) {
	repo := newFakeRepo()
	repo.failSave = true
	pub := &fakePublisher{}
	uc := statsapp.IncrementStatField{Repo: repo, Events: pub}

	uid, _ := shared.NewUserID("user-4")
	if _, err := uc.Execute(context.Background(), uid, stats.StatFieldActivity, 1); err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(pub.events) != 0 {
		t.Errorf("publisher called despite save failure: %+v", pub.events)
	}
}
