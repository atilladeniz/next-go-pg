// Package sse owns the platform-wide Server-Sent Events broker. It
// is a fan-out: every connected HTTP client receives every broadcast
// event, with drop-on-buffer-full back-pressure so a slow client
// cannot stall the rest.
package sse

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/pkg/metrics"
)

// defaultHeartbeatInterval is the keepalive cadence in production.
// Proxies / load balancers (nginx, Cloud Run, Cloudflare) typically
// idle-kill connections around 30–60 s; 25 s keeps us comfortably
// under that.
const defaultHeartbeatInterval = 25 * time.Second

// Event is a single SSE message.
type Event struct {
	Type string
	Data string
}

// Broker manages SSE client connections.
type Broker struct {
	clients    map[chan Event]struct{}
	register   chan chan Event
	unregister chan chan Event
	broadcast  chan Event
	done       chan struct{} // closed by Shutdown to signal run() to exit
	closed     chan struct{} // closed by run() once it has exited

	mu       sync.RWMutex
	shutdown bool

	heartbeatInterval time.Duration
}

// Option configures a Broker. Production callers should not need any.
type Option func(*Broker)

// WithHeartbeatInterval overrides the SSE keepalive cadence. Mainly
// used by tests; the default is 25 s.
func WithHeartbeatInterval(d time.Duration) Option {
	return func(b *Broker) { b.heartbeatInterval = d }
}

// NewBroker creates a new SSE broker and starts its dispatch loop.
// Call Shutdown to stop the loop and release the goroutine.
func NewBroker(opts ...Option) *Broker {
	b := &Broker{
		clients:           make(map[chan Event]struct{}),
		register:          make(chan chan Event),
		unregister:        make(chan chan Event),
		broadcast:         make(chan Event),
		done:              make(chan struct{}),
		closed:            make(chan struct{}),
		heartbeatInterval: defaultHeartbeatInterval,
	}
	for _, o := range opts {
		o(b)
	}
	go b.run()
	return b
}

func (b *Broker) run() {
	defer close(b.closed)
	for {
		select {
		case <-b.done:
			b.drainClients()
			return

		case client := <-b.register:
			b.mu.Lock()
			b.clients[client] = struct{}{}
			b.mu.Unlock()
			metrics.SSEConnectionsActive.Inc()

		case client := <-b.unregister:
			b.mu.Lock()
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client)
				metrics.SSEConnectionsActive.Dec()
			}
			b.mu.Unlock()

		case event := <-b.broadcast:
			b.mu.RLock()
			for client := range b.clients {
				select {
				case client <- event:
					metrics.SSEMessagesSent.Inc()
				default:
					// Client buffer full — drop rather than block
					// the dispatch loop. The slow client will catch
					// up on the next event or, more likely, be
					// disconnected by its HTTP client.
					metrics.SSEEventsDropped.Inc()
				}
			}
			b.mu.RUnlock()
		}
	}
}

func (b *Broker) drainClients() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for client := range b.clients {
		delete(b.clients, client)
		close(client)
		metrics.SSEConnectionsActive.Dec()
	}
}

// Broadcast sends an event to every connected client. No-op after Shutdown.
func (b *Broker) Broadcast(eventType string, data string) {
	if b.isShutdown() {
		return
	}
	select {
	case b.broadcast <- Event{Type: eventType, Data: data}:
	case <-b.done:
		// Shutdown raced with this Broadcast — drop the event.
	}
}

// Shutdown signals the dispatch loop to exit and waits for it (or
// for ctx to expire). Safe to call multiple times.
func (b *Broker) Shutdown(ctx context.Context) error {
	b.mu.Lock()
	if b.shutdown {
		b.mu.Unlock()
		// Already shutting down — still wait for run() to finish.
		select {
		case <-b.closed:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	b.shutdown = true
	b.mu.Unlock()

	close(b.done)
	select {
	case <-b.closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *Broker) isShutdown() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.shutdown
}

// ServeHTTP streams SSE events to a single client. A periodic
// heartbeat comment keeps idle connections alive behind proxies.
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := make(chan Event, 10)

	// Register. If the broker is mid-shutdown, the register channel
	// has no reader — fall through to a 503 instead of blocking.
	select {
	case b.register <- client:
	case <-b.done:
		http.Error(w, "SSE broker shutting down", http.StatusServiceUnavailable)
		return
	case <-r.Context().Done():
		return
	}

	defer func() {
		// Best-effort unregister. If shutdown raced ahead, the run
		// loop already closed our channel in drainClients.
		select {
		case b.unregister <- client:
		case <-b.done:
		}
	}()

	fmt.Fprint(w, "event: connected\ndata: {\"status\":\"ok\"}\n\n")
	flusher.Flush()

	ticker := time.NewTicker(b.heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-client:
			if !ok {
				// Broker drained us during shutdown.
				return
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, event.Data)
			flusher.Flush()

		case <-ticker.C:
			// SSE comment line — ignored by clients, kept alive by proxies.
			fmt.Fprint(w, ": heartbeat\n\n")
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

// ClientCount returns the number of currently connected clients.
func (b *Broker) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}
