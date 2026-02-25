package engine

// AnalysisResult is the complete output produced by any Thinker implementation.
// All fields are populated before being handed to the Formatter.
type AnalysisResult struct {
	Title         string
	Summary       string
	Probabilities []Probability
	RiskIndex     int
	Citations     []Citation
	Conclusion    string
	ClosingLine   string
}

// Probability represents a single entry in the pseudo-statistical breakdown.
// The Label describes the outcome; Percentage is a suspiciously precise number.
type Probability struct {
	Label      string
	Percentage float64
}

// Citation represents a single fabricated academic reference.
// All citations are entirely fictional. Any resemblance to real journals
// is a symptom of academic overexposure.
type Citation struct {
	Index  int
	Source string
}
