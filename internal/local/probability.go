package local

import (
	"math"
	"math/rand"

	"github.com/rishichawda/overthinker/internal/engine"
	"github.com/rishichawda/overthinker/internal/utils"
)

// outcomeLabels is the canonical pool of dramatic outcome descriptions.
var outcomeLabels = []string{
	"chance of immediate regret",
	"chance of mild existential dread",
	"chance of ambiguous, unresolvable outcome",
	"chance of catastrophic nostalgia",
	"chance of unexpected, inconvenient clarity",
	"chance of productive downward spiral",
	"chance of overanalyzing the analysis itself",
	"chance of dramatic internal monologue",
	"chance of second-guessing this decision tomorrow",
	"chance of googling the same question in 3 days",
	"chance of late-night retroactive justification",
	"chance of unsolicited opinion from a friend",
	"chance of creating a pros/cons list that solves nothing",
	"chance of consulting a horoscope",
	"chance of blaming Mercury retrograde",
	"chance of writing a journal entry about this",
	"chance of inexplicable calm followed by panic",
	"chance of doing it anyway regardless of this report",
}

// generateProbabilities produces 3-5 pseudo-statistical probability entries
// that sum to exactly 100.0%.
func generateProbabilities(rng *rand.Rand) []engine.Probability {
	count := 3 + rng.Intn(3)

	shuffled := utils.ShuffleStrings(rng, outcomeLabels)
	chosen := shuffled[:count]

	weights := make([]float64, count)
	total := 0.0
	for i := 0; i < count; i++ {
		w := 10.0 + rng.Float64()*60.0
		weights[i] = w
		total += w
	}

	probs := make([]engine.Probability, count)
	running := 0.0
	for i := 0; i < count-1; i++ {
		pct := math.Round((weights[i]/total)*1000) / 10
		running += pct
		probs[i] = engine.Probability{Label: chosen[i], Percentage: pct}
	}
	probs[count-1] = engine.Probability{
		Label:      chosen[count-1],
		Percentage: math.Round((100.0-running)*10) / 10,
	}

	return probs
}
