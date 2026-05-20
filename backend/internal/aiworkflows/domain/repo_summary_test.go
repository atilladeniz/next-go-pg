package domain_test

import (
	"strings"
	"testing"
	"time"

	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

func mustRepoURL(t *testing.T, raw string) ai.RepoURL {
	t.Helper()
	u, err := ai.NewRepoURL(raw)
	if err != nil {
		t.Fatalf("NewRepoURL(%q) unexpected error: %v", raw, err)
	}
	return u
}

func mustUserID(t *testing.T) shared.UserID {
	t.Helper()
	uid, err := shared.NewUserID("user-1")
	if err != nil {
		t.Fatalf("NewUserID: %v", err)
	}
	return uid
}

func mustFileSummary(t *testing.T, filename, body string) ai.FileSummary {
	t.Helper()
	fs, err := ai.NewFileSummary(filename, body)
	if err != nil {
		t.Fatalf("NewFileSummary: %v", err)
	}
	return fs
}

func TestNewRepoURL(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{"https github", "https://github.com/owner/repo", false},
		{"https with .git", "https://github.com/owner/repo.git", false},
		{"http allowed", "http://gitea.example.com/owner/repo", false},
		{"trim whitespace", "  https://github.com/owner/repo  ", false},
		{"empty", "", true},
		{"only whitespace", "   ", true},
		{"ssh scheme", "ssh://git@github.com/owner/repo", true},
		{"git scheme", "git://github.com/owner/repo", true},
		{"no scheme", "github.com/owner/repo", true},
		{"shell semicolon", "https://github.com/owner/repo;rm -rf /", true},
		{"shell backtick", "https://github.com/owner/repo`whoami`", true},
		{"shell dollar", "https://github.com/owner/repo$IFS", true},
		{"newline", "https://github.com/owner/repo\nhi", true},
		{"missing host", "https://", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ai.NewRepoURL(tc.raw)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewRepoURL(%q) err=%v, wantErr=%v", tc.raw, err, tc.wantErr)
			}
		})
	}
}

func TestNewStatus(t *testing.T) {
	t.Parallel()
	for _, ok := range []string{"pending", "running", "completed", "failed", "cancelled"} {
		if _, err := ai.NewStatus(ok); err != nil {
			t.Errorf("NewStatus(%q) unexpected error: %v", ok, err)
		}
	}
	if _, err := ai.NewStatus("foobar"); err == nil {
		t.Errorf("NewStatus(\"foobar\") expected error, got nil")
	}
}

func TestStatusIsTerminal(t *testing.T) {
	t.Parallel()
	cases := map[ai.Status]bool{
		ai.StatusPending:   false,
		ai.StatusRunning:   false,
		ai.StatusCompleted: true,
		ai.StatusFailed:    true,
		ai.StatusCancelled: true,
	}
	for s, want := range cases {
		if got := s.IsTerminal(); got != want {
			t.Errorf("%s.IsTerminal()=%v, want %v", s, got, want)
		}
	}
}

func TestNewFileSummary(t *testing.T) {
	t.Parallel()
	if _, err := ai.NewFileSummary("", "body"); err == nil {
		t.Errorf("NewFileSummary with empty filename expected error")
	}
	if _, err := ai.NewFileSummary("   ", "body"); err == nil {
		t.Errorf("NewFileSummary with whitespace filename expected error")
	}
	fs, err := ai.NewFileSummary("README.md", "")
	if err != nil {
		t.Fatalf("NewFileSummary with empty body unexpectedly errored: %v", err)
	}
	if fs.Filename() != "README.md" {
		t.Errorf("Filename mismatch: %q", fs.Filename())
	}
	if fs.Summary() != "" {
		t.Errorf("Summary should be empty, got %q", fs.Summary())
	}
}

func TestRepoSummary_HappyPath(t *testing.T) {
	t.Parallel()
	uid := mustUserID(t)
	url := mustRepoURL(t, "https://github.com/owner/repo")

	r := ai.NewRepoSummary(uid, url)
	if r.Status != ai.StatusPending {
		t.Fatalf("expected initial status pending, got %s", r.Status)
	}
	if len(r.PullEvents()) != 0 {
		t.Fatalf("constructor must not emit events")
	}

	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	if err := r.MarkStarted(now); err != nil {
		t.Fatalf("MarkStarted: %v", err)
	}
	if r.Status != ai.StatusRunning {
		t.Errorf("status after MarkStarted = %s, want running", r.Status)
	}
	events := r.PullEvents()
	if len(events) != 1 {
		t.Fatalf("MarkStarted expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(ai.SummaryStarted); !ok {
		t.Errorf("MarkStarted event type = %T, want SummaryStarted", events[0])
	}

	if err := r.AppendFileSummary(mustFileSummary(t, "main.go", "entry point"), 2); err != nil {
		t.Fatalf("AppendFileSummary: %v", err)
	}
	if err := r.AppendFileSummary(mustFileSummary(t, "lib.go", "lib funcs"), 2); err != nil {
		t.Fatalf("AppendFileSummary: %v", err)
	}
	if len(r.Files) != 2 {
		t.Errorf("Files len = %d, want 2", len(r.Files))
	}
	events = r.PullEvents()
	if len(events) != 2 {
		t.Fatalf("expected 2 FileSummarized events, got %d", len(events))
	}
	first, ok := events[0].(ai.FileSummarized)
	if !ok {
		t.Fatalf("event[0] type = %T", events[0])
	}
	if first.FileIndex != 1 || first.FileCount != 2 || first.Filename != "main.go" {
		t.Errorf("event[0] = %+v", first)
	}
	second := events[1].(ai.FileSummarized)
	if second.FileIndex != 2 || second.Filename != "lib.go" {
		t.Errorf("event[1] = %+v", second)
	}

	if err := r.MarkCompleted("overall summary", now.Add(time.Minute)); err != nil {
		t.Fatalf("MarkCompleted: %v", err)
	}
	if r.Status != ai.StatusCompleted {
		t.Errorf("status = %s, want completed", r.Status)
	}
	if r.Summary != "overall summary" {
		t.Errorf("Summary = %q", r.Summary)
	}
	completed := r.PullEvents()
	if len(completed) != 1 {
		t.Fatalf("MarkCompleted expected 1 event, got %d", len(completed))
	}
	if _, ok := completed[0].(ai.SummaryCompleted); !ok {
		t.Errorf("event type = %T, want SummaryCompleted", completed[0])
	}
}

func TestRepoSummary_IllegalTransitions(t *testing.T) {
	t.Parallel()
	uid := mustUserID(t)
	url := mustRepoURL(t, "https://github.com/owner/repo")
	now := time.Now()

	// MarkStarted from non-pending fails.
	r := ai.NewRepoSummary(uid, url)
	if err := r.MarkStarted(now); err != nil {
		t.Fatalf("first MarkStarted: %v", err)
	}
	if err := r.MarkStarted(now); err == nil {
		t.Errorf("second MarkStarted expected error, got nil")
	}

	// AppendFileSummary on pending fails.
	r2 := ai.NewRepoSummary(uid, url)
	if err := r2.AppendFileSummary(mustFileSummary(t, "x.go", ""), 1); err == nil {
		t.Errorf("AppendFileSummary on pending expected error")
	}

	// MarkCompleted on pending fails.
	r3 := ai.NewRepoSummary(uid, url)
	if err := r3.MarkCompleted("s", now); err == nil {
		t.Errorf("MarkCompleted on pending expected error")
	}

	// MarkFailed terminal -> error.
	r4 := ai.NewRepoSummary(uid, url)
	_ = r4.MarkStarted(now)
	_ = r4.MarkCompleted("s", now)
	if err := r4.MarkFailed("oops", now); err == nil {
		t.Errorf("MarkFailed after completed expected error")
	}

	// MarkCancelled terminal -> error.
	r5 := ai.NewRepoSummary(uid, url)
	_ = r5.MarkStarted(now)
	_ = r5.MarkFailed("nope", now)
	if err := r5.MarkCancelled(now); err == nil {
		t.Errorf("MarkCancelled after failed expected error")
	}
}

func TestRepoSummary_Cancel(t *testing.T) {
	t.Parallel()
	uid := mustUserID(t)
	url := mustRepoURL(t, "https://github.com/owner/repo")
	now := time.Now()

	// Cancel from pending.
	r := ai.NewRepoSummary(uid, url)
	if err := r.MarkCancelled(now); err != nil {
		t.Fatalf("MarkCancelled from pending: %v", err)
	}
	if r.Status != ai.StatusCancelled {
		t.Errorf("status = %s, want cancelled", r.Status)
	}
	events := r.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 SummaryCancelled event")
	}
	if _, ok := events[0].(ai.SummaryCancelled); !ok {
		t.Errorf("event type = %T", events[0])
	}

	// Cancel from running.
	r2 := ai.NewRepoSummary(uid, url)
	_ = r2.MarkStarted(now)
	_ = r2.PullEvents() // drain
	if err := r2.MarkCancelled(now); err != nil {
		t.Fatalf("MarkCancelled from running: %v", err)
	}
}

func TestRepoSummary_EventNames(t *testing.T) {
	t.Parallel()
	type named interface{ EventName() string }
	cases := []struct {
		event named
		want  string
	}{
		{ai.SummaryStarted{}, "aiworkflows.summary_started"},
		{ai.FileSummarized{}, "aiworkflows.file_summarized"},
		{ai.SummaryCompleted{}, "aiworkflows.summary_completed"},
		{ai.SummaryFailed{}, "aiworkflows.summary_failed"},
		{ai.SummaryCancelled{}, "aiworkflows.summary_cancelled"},
	}
	for _, tc := range cases {
		got := tc.event.EventName()
		if got != tc.want {
			t.Errorf("EventName = %q, want %q", got, tc.want)
		}
		if !strings.HasPrefix(got, "aiworkflows.") {
			t.Errorf("event %T name %q missing context prefix", tc.event, got)
		}
	}
}
