// Package llm is the aiworkflows context's LLMClient adapter against a
// running Ollama HTTP service. Dev points it at `http://127.0.0.1:11434`
// (the host-bound port from infra/compose/docker-compose.dev.yml).
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
)

const (
	defaultURL     = "http://127.0.0.1:11434"
	defaultModel   = "gemma4:e4b"
	defaultTimeout = 60 * time.Second
)

// Client is the Ollama HTTP client. The struct is intentionally
// constructor-only so callers cannot bypass the URL/model defaults.
type Client struct {
	url   string
	model string
	http  *http.Client
}

var _ aiapp.LLMClient = (*Client)(nil)

// Config holds optional overrides. Zero values fall back to the
// package-level defaults so callers can pass an empty Config.
type Config struct {
	URL     string
	Model   string
	Timeout time.Duration
}

// NewClient constructs an Ollama client. Pass an empty Config to use
// defaults; pass partial config to override individual fields.
func NewClient(cfg Config) *Client {
	url := cfg.URL
	if url == "" {
		url = defaultURL
	}
	model := cfg.Model
	if model == "" {
		model = defaultModel
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &Client{
		url:   url,
		model: model,
		http:  &http.Client{Timeout: timeout},
	}
}

// Model returns the configured default model name. Useful for logging.
func (c *Client) Model() string { return c.model }

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Generate sends a single non-streaming completion request. Streaming
// is intentionally disabled — the workflow step only cares about the
// final response and per-token streaming would complicate retry logic.
func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(generateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama post: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("ollama status %d: %s", resp.StatusCode, string(raw))
	}

	var out generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	return out.Response, nil
}
