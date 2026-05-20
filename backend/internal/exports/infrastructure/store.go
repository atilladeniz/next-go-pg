// Package infrastructure holds the exports context's adapters: an
// in-memory result store and (via the jobs/ subpackage) the River
// worker + enqueuer.
package infrastructure

import (
	"time"

	exportsapp "github.com/atilladeniz/next-go-pg/backend/internal/exports/application"
)

// MemoryStore is an in-process result store. Single-replica only; for
// horizontal scaling, swap in an S3/RustFS-backed adapter behind the
// same Store port.
type MemoryStore struct {
	exports map[string]*exportsapp.Result
}

var _ exportsapp.Store = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{exports: make(map[string]*exportsapp.Result)}
}

func (s *MemoryStore) Save(id string, result *exportsapp.Result) {
	s.exports[id] = result
}

func (s *MemoryStore) Get(id string) (*exportsapp.Result, bool) {
	result, ok := s.exports[id]
	if !ok {
		return nil, false
	}
	if time.Now().After(result.ExpiresAt) {
		delete(s.exports, id)
		return nil, false
	}
	return result, true
}

func (s *MemoryStore) Delete(id string) {
	delete(s.exports, id)
}
