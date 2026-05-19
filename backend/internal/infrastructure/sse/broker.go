package sse

import (
	"fmt"
	"net/http"
	"sync"
)

// Event represents an SSE event
type Event struct {
	Type string
	Data string
}

// Broker manages SSE client connections
type Broker struct {
	clients    map[chan Event]bool
	register   chan chan Event
	unregister chan chan Event
	broadcast  chan Event
	mu         sync.RWMutex
}

// NewBroker creates a new SSE broker
func NewBroker() *Broker {
	b := &Broker{
		clients:    make(map[chan Event]bool),
		register:   make(chan chan Event),
		unregister: make(chan chan Event),
		broadcast:  make(chan Event),
	}
	go b.run()
	return b
}

func (b *Broker) run() {
	for {
		select {
		case client := <-b.register:
			b.mu.Lock()
			b.clients[client] = true
			b.mu.Unlock()

		case client := <-b.unregister:
			b.mu.Lock()
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client)
			}
			b.mu.Unlock()

		case event := <-b.broadcast:
			b.mu.RLock()
			for client := range b.clients {
				select {
				case client <- event:
				default:
					// Client buffer full, skip
				}
			}
			b.mu.RUnlock()
		}
	}
}

// Broadcast sends an event to all connected clients
func (b *Broker) Broadcast(eventType string, data string) {
	b.broadcast <- Event{Type: eventType, Data: data}
}

// ServeHTTP handles SSE connections
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if flusher is supported
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create client channel
	client := make(chan Event, 10)
	b.register <- client

	// Cleanup on disconnect
	defer func() {
		b.unregister <- client
	}()

	// Send initial connection event
	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"ok\"}\n\n")
	flusher.Flush()

	// Listen for events or client disconnect
	for {
		select {
		case event := <-client:
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, event.Data)
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

// ClientCount returns the number of connected clients
func (b *Broker) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}
