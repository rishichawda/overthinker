package ollama

import "github.com/rishichawda/overthinker/internal/engine"

// responseSchema is the JSON schema passed to Ollama's structured output API.
// It constrains the model to emit a valid, machine-readable JSON object that
// maps directly onto OllamaResponse — no text parsing required.
const responseSchema = `{
	"type": "object",
	"properties": {
		"title": {
			"type": "string",
			"description": "An ALL-CAPS dramatic title summarising the situation"
		},
		"summary": {
			"type": "string",
			"description": "2-3 sentences of alarming pseudo-academic insight"
		},
		"probabilities": {
			"type": "array",
			"description": "3-5 entries that must sum to exactly 100",
			"items": {
				"type": "object",
				"properties": {
					"label":      { "type": "string" },
					"percentage": { "type": "number" }
				},
				"required": ["label", "percentage"]
			}
		},
		"risk_index": {
			"type": "integer",
			"description": "Emotional Risk Index, 0-100",
			"minimum": 0,
			"maximum": 100
		},
		"risk_justification": {
			"type": "string",
			"description": "One sentence justifying the risk index score"
		},
		"citations": {
			"type": "array",
			"description": "2-3 entirely fabricated but plausible academic citations",
			"items": {
				"type": "object",
				"properties": {
					"source": { "type": "string" }
				},
				"required": ["source"]
			}
		},
		"conclusion": {
			"type": "string",
			"description": "2-3 sentences of theatrical finality"
		},
		"closing_remark": {
			"type": "string",
			"description": "One self-aware, witty closing sentence"
		}
	},
	"required": [
		"title", "summary", "probabilities",
		"risk_index", "risk_justification",
		"citations", "conclusion", "closing_remark"
	]
}`

// OllamaResponse is the structured JSON object the LLM must return.
type OllamaResponse struct {
	Title             string              `json:"title"`
	Summary           string              `json:"summary"`
	Probabilities     []probabilityEntry  `json:"probabilities"`
	RiskIndex         int                 `json:"risk_index"`
	RiskJustification string              `json:"risk_justification"`
	Citations         []citationEntry     `json:"citations"`
	Conclusion        string              `json:"conclusion"`
	ClosingRemark     string              `json:"closing_remark"`
}

type probabilityEntry struct {
	Label      string  `json:"label"`
	Percentage float64 `json:"percentage"`
}

type citationEntry struct {
	Source string `json:"source"`
}

// toAnalysisResult converts the structured LLM response into the shared
// engine.AnalysisResult used by the formatter.
func (r *OllamaResponse) toAnalysisResult() *engine.AnalysisResult {
	probs := make([]engine.Probability, len(r.Probabilities))
	for i, p := range r.Probabilities {
		probs[i] = engine.Probability{Label: p.Label, Percentage: p.Percentage}
	}

	citations := make([]engine.Citation, len(r.Citations))
	for i, c := range r.Citations {
		citations[i] = engine.Citation{Index: i + 1, Source: c.Source}
	}

	riskIndex := r.RiskIndex
	if riskIndex < 0 {
		riskIndex = 0
	}
	if riskIndex > 100 {
		riskIndex = 100
	}

	return &engine.AnalysisResult{
		Title:         r.Title,
		Summary:       r.Summary,
		Probabilities: probs,
		RiskIndex:     riskIndex,
		Citations:     citations,
		Conclusion:    r.Conclusion,
		ClosingLine:   r.ClosingRemark,
	}
}
