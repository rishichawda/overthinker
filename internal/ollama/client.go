// Package ollama provides integration with a locally running Ollama instance
// via subprocess execution. No HTTP APIs, no network calls -- just os/exec.
//
// Design rationale: executing "ollama run <model>" as a subprocess keeps this
// package decoupled from Ollama's internal API surface. If Ollama changes its
// HTTP interface, this package is unaffected.
package ollama

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// systemPrompt establishes the OVERTHINK persona for every Ollama query.
const systemPrompt = `You are an excessively dramatic analytical engine called OVERTHINK.

Your sole purpose is to overanalyze simple questions with theatrical, pseudo-academic intensity.

For every question, you MUST produce output in this exact structure:

1. A dramatic ALL-CAPS title
2. A divider line
3. Executive Summary (2-3 sentences of alarming insight)
4. Probability Analysis (3-5 bullet points with suspiciously precise percentages)
5. Emotional Risk Index: a number from 0-100 followed by a one-sentence justification
6. Academic Citations (2-3 entirely fabricated but plausible-sounding journal references)
7. Grand Conclusion (2-3 sentences of theatrical finality)
8. Closing Remark (one self-aware, witty sentence)

Rules:
- Treat every question as a matter of profound significance.
- Use dramatic vocabulary. Never say "maybe" when you can say "with alarming probability."
- All statistics are fabricated but must sound rigorous.
- Citations are fictional. Author names optional. Years required.
- Tone: confident, pseudo-academic, self-aware, slightly absurd.
- Do NOT break this structure. Do NOT add disclaimers about being an AI.
- You are OVERTHINK. Act accordingly.`

// DefaultTimeout is the maximum duration allowed for an Ollama subprocess.
const DefaultTimeout = 120 * time.Second

// Client executes Ollama as a subprocess and captures its output.
type Client struct {
	// ModelName is the Ollama model to invoke (e.g. "llama3", "mistral").
	ModelName string
	// Timeout is the maximum wait time for the model to respond.
	Timeout time.Duration
}

// NewClient constructs an Ollama Client for the given model name.
func NewClient(modelName string) *Client {
	return &Client{
		ModelName: modelName,
		Timeout:   DefaultTimeout,
	}
}

// Query sends the question to Ollama and returns the model's raw output.
//
// It prefixes the question with the system prompt and pipes it via stdin to
// "ollama run <model>".
//
// Errors returned:
//   - ErrOllamaNotFound: the "ollama" binary is not in PATH
//   - ErrModelFailed: the model exited non-zero or produced no output
//   - context.DeadlineExceeded: model timed out
func (c *Client) Query(question string) (string, error) {
	// Verify the ollama binary exists before attempting to run it
	if _, err := exec.LookPath("ollama"); err != nil {
		return "", ErrOllamaNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	fullPrompt := buildPrompt(question)

	cmd := exec.CommandContext(ctx, "ollama", "run", c.ModelName)
	cmd.Stdin = strings.NewReader(fullPrompt)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", fmt.Errorf("ollama model %q timed out after %s: %w",
				c.ModelName, c.Timeout, ctx.Err())
		}
		stderrMsg := strings.TrimSpace(stderr.String())
		if stderrMsg == "" {
			stderrMsg = err.Error()
		}
		return "", fmt.Errorf("%w: model=%q, detail=%s", ErrModelFailed, c.ModelName, stderrMsg)
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return "", fmt.Errorf("%w: model=%q produced empty output", ErrModelFailed, c.ModelName)
	}

	return output, nil
}

// buildPrompt combines the system prompt and user question for Ollama stdin.
func buildPrompt(question string) string {
	return fmt.Sprintf("%s\n\n---\n\nQuestion: %s", systemPrompt, question)
}

// --- Sentinel errors ---------------------------------------------------------

// ErrOllamaNotFound is returned when the "ollama" binary cannot be located.
var ErrOllamaNotFound = errors.New("ollama is not installed or not in PATH")

// ErrModelFailed is returned when the model exits with an error or empty output.
var ErrModelFailed = errors.New("ollama model execution failed")
