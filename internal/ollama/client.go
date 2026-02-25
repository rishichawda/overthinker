// Package ollama provides integration with a locally running Ollama server
// via the official Ollama Go HTTP API client (default: localhost:11434).
//
// It uses Ollama's structured-output feature (Format: JSON schema) to obtain
// a machine-readable response that requires no text parsing.
package ollama

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	ollamaapi "github.com/ollama/ollama/api"
	"github.com/rishichawda/overthinker/internal/engine"
)

// systemPrompt establishes the OVERTHINK persona.
// Section structure is no longer described here — the JSON schema enforces it.
const systemPrompt = `You are an excessively dramatic analytical engine called OVERTHINK.

Your sole purpose is to overanalyze simple questions with theatrical, pseudo-academic intensity.

Rules:
- Treat every question as a matter of profound significance.
- Use dramatic vocabulary. Never say "maybe" when you can say "with alarming probability."
- All statistics are fabricated but must sound rigorous.
- Probability percentages must sum to exactly 100.
- Citations are entirely fictional. Author names and years required.
- Tone: confident, pseudo-academic, self-aware, slightly absurd.
- Do NOT add disclaimers about being an AI.
- You are OVERTHINK. Act accordingly.`

// DefaultTimeout is the maximum duration allowed for an Ollama request.
const DefaultTimeout = 120 * time.Second

// OllamaHost is the default Ollama server address.
const OllamaHost = "http://localhost:11434"

// Client queries the Ollama HTTP API and implements engine.Thinker.
type Client struct {
	// ModelName is the Ollama model to invoke (e.g. "llama3", "mistral").
	ModelName string
	// Timeout is the maximum wait time for the model to respond.
	Timeout time.Duration
	// Host is the Ollama server base URL.
	Host string
}

// NewClient constructs an Ollama Client for the given model name.
func NewClient(modelName string) *Client {
	return &Client{
		ModelName: modelName,
		Timeout:   DefaultTimeout,
		Host:      OllamaHost,
	}
}

// Analyze implements engine.Thinker. It queries the Ollama server using
// structured JSON output constrained by responseSchema, then deserialises the
// response directly into an AnalysisResult — no text parsing required.
//
// Errors returned:
//   - ErrOllamaNotFound: the Ollama server is not reachable
//   - ErrModelNotFound: the requested model is not available on the server
//   - ErrModelFailed: the model returned an error, empty, or unparseable output
//   - context.DeadlineExceeded: request timed out
func (c *Client) Analyze(question string) (*engine.AnalysisResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	serverURL, err := url.Parse(c.Host)
	if err != nil {
		return nil, fmt.Errorf("invalid Ollama host %q: %w", c.Host, err)
	}

	client := ollamaapi.NewClient(serverURL, http.DefaultClient)

	if err := c.checkServer(ctx, client); err != nil {
		return nil, err
	}
	if err := c.checkModel(ctx, client); err != nil {
		return nil, err
	}

	var sb strings.Builder

	req := &ollamaapi.GenerateRequest{
		Model:  c.ModelName,
		System: systemPrompt,
		Prompt: fmt.Sprintf("Question: %s", question),
		Format: json.RawMessage(responseSchema),
		Stream: boolPtr(true),
	}

	err = client.Generate(ctx, req, func(resp ollamaapi.GenerateResponse) error {
		sb.WriteString(resp.Response)
		return nil
	})
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("ollama model %q timed out after %s: %w",
				c.ModelName, c.Timeout, ctx.Err())
		}
		return nil, fmt.Errorf("%w: model=%q, detail=%s", ErrModelFailed, c.ModelName, err.Error())
	}

	raw := strings.TrimSpace(sb.String())
	if raw == "" {
		return nil, fmt.Errorf("%w: model=%q produced empty output", ErrModelFailed, c.ModelName)
	}

	var response OllamaResponse
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return nil, fmt.Errorf("%w: model=%q returned invalid JSON: %s",
			ErrModelFailed, c.ModelName, err.Error())
	}

	return response.toAnalysisResult(), nil
}

// checkServer pings the Ollama server to verify it is reachable.
func (c *Client) checkServer(ctx context.Context, client *ollamaapi.Client) error {
	if err := client.Heartbeat(ctx); err != nil {
		return ErrOllamaNotFound
	}
	return nil
}

// checkModel verifies the requested model is available on the Ollama server.
func (c *Client) checkModel(ctx context.Context, client *ollamaapi.Client) error {
	resp, err := client.List(ctx)
	if err != nil {
		return nil // non-fatal; let generate surface the error
	}
	for _, m := range resp.Models {
		if m.Name == c.ModelName || strings.HasPrefix(m.Name, c.ModelName+":") {
			return nil
		}
	}
	return fmt.Errorf("%w: %q", ErrModelNotFound, c.ModelName)
}

// boolPtr returns a pointer to a bool — required by GenerateRequest.Stream.
func boolPtr(b bool) *bool { return &b }

// --- Sentinel errors ---------------------------------------------------------

// ErrOllamaNotFound is returned when the Ollama server is not reachable.
var ErrOllamaNotFound = errors.New("ollama server is not running (start it with: ollama serve)")

// ErrModelNotFound is returned when the requested model is not installed.
var ErrModelNotFound = errors.New("ollama model not found (install it with: ollama pull)")

// ErrModelFailed is returned when the model returns an error, empty, or unparseable output.
var ErrModelFailed = errors.New("ollama model execution failed")
