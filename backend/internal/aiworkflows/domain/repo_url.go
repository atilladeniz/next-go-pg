// Package domain is the aiworkflows bounded context's aggregate model.
// Pure Go, no persistence, no HTTP, no workflow-engine SDK, no other
// context's domain. The Hatchet SDK lives in infrastructure/.
package domain

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// RepoURL is a validated public Git repository URL accepted by the
// summarization workflow. Constructed via NewRepoURL so the shell-safe
// invariant is enforced at the boundary.
type RepoURL string

// shellMetachars are rejected in raw input to keep us safe if any code
// path ever shells out to `git`. The infrastructure-layer cloner uses
// go-git instead of exec.Command, but we still defend in depth: the
// domain refuses to admit a value that could ever become a shell hazard.
const shellMetachars = ";|&`$<>(){}!*?\\\"'\n\r\t "

// NewRepoURL parses and validates a raw repo URL. http/https schemes
// only; host required; the optional `.git` suffix is accepted.
func NewRepoURL(raw string) (RepoURL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("repo url must not be empty")
	}
	if strings.ContainsAny(raw, shellMetachars) {
		return "", errors.New("repo url contains forbidden characters")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("repo url is not a valid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("repo url must use http or https (got %q)", u.Scheme)
	}
	if u.Host == "" {
		return "", errors.New("repo url is missing host")
	}
	return RepoURL(raw), nil
}

func (r RepoURL) String() string { return string(r) }
