// Package git is the aiworkflows context's RepoCloner adapter. It uses
// the pure-Go `go-git` library — no shell-out, no `git` binary
// dependency, no shell-injection surface (defense in depth with the
// RepoURL value object that already rejects shell metachars).
package git

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	gogit "github.com/go-git/go-git/v5"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
)

// Cloner is the RepoCloner adapter. MaxBytes caps the total unpacked
// repository size to defend against malicious or pathologically large
// repos. SingleBranch keeps the clone shallow (HEAD of default branch
// only) so the per-run disk footprint stays small.
type Cloner struct {
	BaseDir  string // parent directory for the working copies, e.g. os.TempDir()
	MaxBytes int64  // total unpacked size cap; 0 = no cap
}

var _ aiapp.RepoCloner = (*Cloner)(nil)

// NewCloner constructs a Cloner with sensible defaults.
func NewCloner(baseDir string, maxBytes int64) *Cloner {
	if baseDir == "" {
		baseDir = os.TempDir()
	}
	return &Cloner{BaseDir: baseDir, MaxBytes: maxBytes}
}

// Clone performs a shallow clone of url into a freshly-created temp dir
// under BaseDir. The returned ClonedRepo.Cleanup removes the directory.
// Caller MUST invoke Cleanup even on error — we honour the contract by
// only returning Cleanup-bearing values on success.
func (c *Cloner) Clone(ctx context.Context, url ai.RepoURL) (aiapp.ClonedRepo, error) {
	dir, err := os.MkdirTemp(c.BaseDir, "repo-summary-*")
	if err != nil {
		return aiapp.ClonedRepo{}, fmt.Errorf("mkdir temp: %w", err)
	}

	cleanup := func() error { return os.RemoveAll(dir) }

	_, err = gogit.PlainCloneContext(ctx, dir, false, &gogit.CloneOptions{
		URL:               url.String(),
		Depth:             1,
		SingleBranch:      true,
		ShallowSubmodules: true,
		Progress:          io.Discard,
	})
	if err != nil {
		_ = cleanup()
		return aiapp.ClonedRepo{}, fmt.Errorf("clone %s: %w", url.String(), err)
	}

	if c.MaxBytes > 0 {
		size, err := dirSize(dir)
		if err != nil {
			_ = cleanup()
			return aiapp.ClonedRepo{}, fmt.Errorf("measure clone: %w", err)
		}
		if size > c.MaxBytes {
			_ = cleanup()
			return aiapp.ClonedRepo{}, fmt.Errorf("clone size %d bytes exceeds limit %d", size, c.MaxBytes)
		}
	}

	return aiapp.ClonedRepo{Path: dir, Cleanup: cleanup}, nil
}

func dirSize(root string) (int64, error) {
	var total int64
	err := filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total, err
}
