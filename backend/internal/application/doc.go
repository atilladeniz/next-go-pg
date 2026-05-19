// Package application is the use-case layer.
//
// It declares the ports (interfaces) consumed by use cases and exposes each
// use case as a struct with an Execute(ctx, ...) method. Implementations of
// the ports live under internal/infrastructure/. Handlers, jobs, and the SSE
// broker depend on this package — never the reverse.
package application
