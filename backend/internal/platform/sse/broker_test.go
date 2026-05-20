package sse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// waitFor polls f every 1ms until it returns true or the timeout
// elapses. Used instead of arbitrary sleeps for state synchronization.
func waitFor(t *testing.T, timeout time.Duration, f func() bool) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return false
}

func TestBroadcastDeliversToAllClients(t *testing.T) {
	b := NewBroker()
	defer func() { _ = b.Shutdown(context.Background()) }()

	c1 := make(chan Event, 1)
	c2 := make(chan Event, 1)
	b.register <- c1
	b.register <- c2

	if !waitFor(t, time.Second, func() bool { return b.ClientCount() == 2 }) {
		t.Fatalf("expected 2 clients, got %d", b.ClientCount())
	}

	b.Broadcast("stats-updated", `{"field":"projects"}`)

	for i, c := range []chan Event{c1, c2} {
		select {
		case ev := <-c:
			if ev.Type != "stats-updated" || ev.Data != `{"field":"projects"}` {
				t.Errorf("client %d: unexpected event %+v", i, ev)
			}
		case <-time.After(time.Second):
			t.Fatalf("client %d did not receive event", i)
		}
	}
}

func TestUnregisterStopsDelivery(t *testing.T) {
	b := NewBroker()
	defer func() { _ = b.Shutdown(context.Background()) }()

	c1 := make(chan Event, 1)
	c2 := make(chan Event, 1)
	b.register <- c1
	b.register <- c2

	if !waitFor(t, time.Second, func() bool { return b.ClientCount() == 2 }) {
		t.Fatalf("setup: expected 2 clients")
	}

	b.unregister <- c1
	if !waitFor(t, time.Second, func() bool { return b.ClientCount() == 1 }) {
		t.Fatalf("expected client count to drop to 1, got %d", b.ClientCount())
	}

	b.Broadcast("ping", "1")

	select {
	case _, ok := <-c1:
		if ok {
			t.Errorf("unregistered client should not receive events")
		}
		// closed channel read returning ok=false is the expected
		// outcome from unregister having closed c1.
	case <-time.After(50 * time.Millisecond):
		// Also acceptable: nothing arrived.
	}

	select {
	case ev := <-c2:
		if ev.Type != "ping" {
			t.Errorf("unexpected event on c2: %+v", ev)
		}
	case <-time.After(time.Second):
		t.Fatalf("c2 should still receive events")
	}
}

func TestSlowClientDoesNotBlockOthers(t *testing.T) {
	b := NewBroker()
	defer func() { _ = b.Shutdown(context.Background()) }()

	// Slow client: capacity 1, immediately full.
	slow := make(chan Event, 1)
	slow <- Event{Type: "filler"}

	fast := make(chan Event, 10)

	b.register <- slow
	b.register <- fast
	if !waitFor(t, time.Second, func() bool { return b.ClientCount() == 2 }) {
		t.Fatalf("setup: expected 2 clients")
	}

	b.Broadcast("update", "x")

	select {
	case ev := <-fast:
		if ev.Type != "update" {
			t.Errorf("fast client got wrong event: %+v", ev)
		}
	case <-time.After(time.Second):
		t.Fatalf("fast client did not receive broadcast — slow client blocked dispatch")
	}
}

func TestShutdownIsIdempotentAndUnblocks(t *testing.T) {
	b := NewBroker()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := b.Shutdown(ctx); err != nil {
		t.Fatalf("first shutdown: %v", err)
	}
	if err := b.Shutdown(ctx); err != nil {
		t.Fatalf("second shutdown should be a no-op error, got: %v", err)
	}

	// Broadcast after shutdown must not panic or block.
	done := make(chan struct{})
	go func() {
		b.Broadcast("ignored", "x")
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Broadcast blocked after Shutdown")
	}
}

func TestShutdownDrainsConnectedClients(t *testing.T) {
	b := NewBroker()

	c := make(chan Event, 1)
	b.register <- c
	if !waitFor(t, time.Second, func() bool { return b.ClientCount() == 1 }) {
		t.Fatalf("setup: expected 1 client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := b.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	if b.ClientCount() != 0 {
		t.Errorf("expected drained client count 0, got %d", b.ClientCount())
	}

	// The client channel must be closed by the broker (signal to the
	// HTTP handler that it should return).
	select {
	case _, ok := <-c:
		if ok {
			t.Errorf("expected client channel to be closed after shutdown")
		}
	default:
		t.Errorf("expected client channel to be closed and readable, but read would block")
	}
}

// flushRecorder wraps httptest.ResponseRecorder so it satisfies
// http.Flusher (the broker's ServeHTTP refuses to run otherwise).
// Reads of the body are guarded by a mutex because the broker writes
// concurrently with the test reading.
type flushRecorder struct {
	rec *httptest.ResponseRecorder
	mu  sync.Mutex
}

func (f *flushRecorder) Header() http.Header { return f.rec.Header() }
func (f *flushRecorder) WriteHeader(s int)   { f.rec.WriteHeader(s) }
func (f *flushRecorder) Write(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.rec.Body.Write(p)
}
func (f *flushRecorder) Flush() {}
func (f *flushRecorder) body() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.rec.Body.String()
}

func TestServeHTTPSendsHeartbeats(t *testing.T) {
	b := NewBroker(WithHeartbeatInterval(10 * time.Millisecond))
	defer func() { _ = b.Shutdown(context.Background()) }()

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/events", nil).WithContext(ctx)
	rec := &flushRecorder{rec: httptest.NewRecorder()}

	done := make(chan struct{})
	go func() {
		b.ServeHTTP(rec, req)
		close(done)
	}()

	// Wait long enough for at least one heartbeat tick (~3 intervals).
	if !waitFor(t, time.Second, func() bool {
		return strings.Contains(rec.body(), ": heartbeat")
	}) {
		cancel()
		<-done
		t.Fatalf("expected at least one heartbeat in body, got: %q", rec.body())
	}

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("ServeHTTP did not return after request cancellation")
	}

	body := rec.body()
	if !strings.Contains(body, "event: connected") {
		t.Errorf("expected initial 'connected' event in body, got: %q", body)
	}
}

func TestServeHTTPViaTestServerStreamsBroadcasts(t *testing.T) {
	b := NewBroker(WithHeartbeatInterval(time.Hour))
	defer func() { _ = b.Shutdown(context.Background()) }()

	srv := httptest.NewServer(http.HandlerFunc(b.ServeHTTP))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	resp, err := http.DefaultClient.Do(req) //nolint:bodyclose // closed via defer below
	if err != nil {
		t.Fatalf("client do: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Wait until the broker registers our client.
	if !waitFor(t, time.Second, func() bool { return b.ClientCount() >= 1 }) {
		t.Fatalf("client did not register")
	}

	b.Broadcast("stats-updated", `{"field":"projects"}`)

	// Read at least until we see the broadcast event.
	got := make(chan string, 1)
	go func() {
		buf := make([]byte, 4096)
		var collected strings.Builder
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				collected.Write(buf[:n])
				if strings.Contains(collected.String(), "event: stats-updated") {
					got <- collected.String()
					return
				}
			}
			if err != nil {
				got <- collected.String()
				return
			}
		}
	}()

	select {
	case s := <-got:
		if !strings.Contains(s, "event: stats-updated") {
			t.Errorf("expected broadcast in stream, got: %q", s)
		}
	case <-ctx.Done():
		t.Fatalf("timeout waiting for broadcast")
	}
}
