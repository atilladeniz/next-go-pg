// Package composition is the composition root. It builds the application
// dependency graph (DB, repositories, use cases, handlers, SSE broker,
// background workers) and is consumed by cmd/server. Nothing under
// internal/ depends on this package; it depends on everything inward.
package composition
