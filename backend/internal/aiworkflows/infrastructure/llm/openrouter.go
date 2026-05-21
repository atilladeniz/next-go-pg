package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
)

// OpenRouterClient is the LLMClient implementation against the
// OpenRouter gateway (https://openrouter.ai). The wire protocol is
// OpenAI-compatible, so any model OpenRouter routes to (Llama, Claude,
// Gemini, GPT, etc.) is reachable with a single API key. We default
// to a free-tier auto-router; the model is configurable via
// `OPENROUTER_MODEL`.
type OpenRouterClient struct {
	url     string
	apiKey  string
	model   string
	referer string
	title   string
	http    *http.Client
}

var _ aiapp.LLMClient = (*OpenRouterClient)(nil)

// OpenRouterConfig holds optional overrides. Most fields fall back to
// env-driven defaults the constructor consumes; pass an empty Config to
// pick up defaults entirely.
type OpenRouterConfig struct {
	URL     string        // default https://openrouter.ai
	APIKey  string        // REQUIRED — no zero-value default
	Model   string        // default openai/gpt-oss-120b
	Referer string        // optional, used in OpenRouter analytics
	Title   string        // optional, used in OpenRouter analytics
	Timeout time.Duration // default 60s (cloud is fast vs local)
}

const (
	defaultOpenRouterURL     = "https://openrouter.ai"
	defaultOpenRouterModel   = "openai/gpt-oss-120b"
	defaultOpenRouterTimeout = 60 * time.Second
)

// NewOpenRouterClient constructs the client. Returns an error if the API
// key is empty — the cloud LLM is useless without it and we'd rather
// fail fast at composition time than at first request.
func NewOpenRouterClient(cfg OpenRouterConfig) (*OpenRouterClient, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("openrouter: api key is required")
	}
	url := cfg.URL
	if url == "" {
		url = defaultOpenRouterURL
	}
	model := cfg.Model
	if model == "" {
		model = defaultOpenRouterModel
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultOpenRouterTimeout
	}
	return &OpenRouterClient{
		url:     url,
		apiKey:  cfg.APIKey,
		model:   model,
		referer: cfg.Referer,
		title:   cfg.Title,
		http:    &http.Client{Timeout: timeout},
	}, nil
}

// Model returns the configured model identifier. Useful for logging.
func (c *OpenRouterClient) Model() string { return c.model }

// Ping verifies the API key is accepted by OpenRouter. Used at startup
// to fail-fast instead of letting the first workflow run hit a 401.
// Calls `/api/v1/auth/key` — a cheap auth-info endpoint that doesn't
// burn any inference budget.
func (c *OpenRouterClient) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"/api/v1/auth/key", nil)
	if err != nil {
		return fmt.Errorf("build ping request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("openrouter unreachable: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("openrouter: api key rejected (401) — check OPENROUTER_API_KEY")
	}
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("openrouter ping status %d: %s", resp.StatusCode, string(raw))
	}
	return nil
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    any    `json:"code,omitempty"`
	} `json:"error,omitempty"`
}

// Generate sends a non-streaming chat completion. The prompt becomes a
// single user message; OpenRouter then routes to whichever provider
// backs `c.model`.
func (c *OpenRouterClient) Generate(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+"/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if c.referer != "" {
		req.Header.Set("HTTP-Referer", c.referer)
	}
	if c.title != "" {
		req.Header.Set("X-Title", c.title)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("openrouter post: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("openrouter status %d: %s", resp.StatusCode, string(raw))
	}

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("openrouter error: %s", out.Error.Message)
	}
	if len(out.Choices) == 0 {
		return "", errors.New("openrouter: empty choices")
	}
	return out.Choices[0].Message.Content, nil
}
