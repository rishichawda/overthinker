package engine

import (
	"math"
	"math/rand"

	"github.com/rishichawda/overthinker/internal/utils"
)

// Probability represents a single entry in the pseudo-statistical breakdown.
// The Label describes the outcome; Percentage is a suspiciously precise number.
type Probability struct {
	Label      string
	Percentage float64
}

// outcomeLabels is the canonical pool of dramatic outcome descriptions.
// Designed to sound plausible while remaining universally applicable to any
// question a human might ask at 2am.
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

// GenerateProbabilities produces 3-5 pseudo-statistical probability entries
// that sum to exactly 100.0%. Values are seeded deterministically so identical
// questions yield identical, reproducible anxiety.
func GenerateProbabilities(rng *rand.Rand) []Probability {
	// Pick between 3 and 5 outcomes for variety
	count := 3 + rng.Intn(3)

	// Shuffle the label pool and take the first count entries
	shuffled := utils.ShuffleStrings(rng, outcomeLabels)
	chosen := shuffled[:count]

	// Generate raw random weights, biased toward mid-range (10-70) for plausibility
	weights := make([]float64, count)
	total := 0.0
	for i := 0; i < count; i++ {
		w := 10.0 + rng.Float64()*60.0
		weights[i] = w
		total += w
	}

	// Normalize to percentages with 1 decimal place for that authentic academic feel
	probs := make([]Probability, count)
	running := 0.0
	for i := 0; i < count-1; i++ {
		pct := math.Round((weights[i]/total)*1000) / 10
		running += pct
		probs[i] = Probability{Label: chosen[i], Percentage: pct}
	}
	// Assign remainder to last entry so the total is provably 100.0%
	probs[count-1] = Probability{
		Label:      chosen[count-1],
		Percentage: math.Round((100.0-running)*10) / 10,
	}

	return probs
}
