package engine

// Thinker is the common interface for all analysis backends.
// Implementations include the local deterministic engine (internal/local)
// and the Ollama LLM client (internal/ollama).
type Thinker interface {
	Analyze(question string) (*AnalysisResult, error)
}
