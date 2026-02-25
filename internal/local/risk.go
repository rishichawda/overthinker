package local

import (
	"math/rand"
	"strings"
)

// riskKeywords maps thematic keywords to a risk score increment.
var riskKeywords = map[string]int{
	// Romantic peril
	"ex": 25, "text": 10, "love": 15, "date": 12,
	"relationship": 18, "breakup": 28, "feelings": 14,
	"heart": 16, "miss": 20, "crush": 13,

	// Professional anxiety
	"quit": 22, "job": 15, "career": 12, "boss": 10,
	"fire": 20, "fired": 25, "resign": 22, "startup": 18, "salary": 10,

	// Existential dread
	"life": 8, "meaning": 20, "purpose": 18, "late": 15, "old": 10,
	"future": 12, "dead": 30, "die": 28, "worth": 16, "point": 14,
	"regret": 22, "mistake": 18, "wrong": 12, "mess": 10,
	"failing": 20, "failed": 22, "failure": 25,

	// Financial anxiety
	"money": 12, "debt": 20, "broke": 18, "invest": 8, "savings": 10,

	// Social pressure
	"family": 15, "friend": 8, "alone": 20, "lonely": 22,
	"trust": 14, "lie": 16, "truth": 10, "tell": 8,

	// Decision paralysis
	"should": 5, "could": 4, "would": 4, "maybe": 8,
	"start": 6, "stop": 8, "leave": 14, "stay": 10,
	"change": 10, "try": 5, "move": 12, "wait": 6,
	"never": 12, "always": 8, "finally": 10,
}

// calculateRiskIndex computes the Emotional Risk Index (0-100) for a given question.
func calculateRiskIndex(question string, rng *rand.Rand) int {
	lower := strings.ToLower(question)
	words := strings.Fields(lower)

	base := 20 + rng.Intn(20)

	accumulated := 0
	seen := make(map[string]bool)
	for _, word := range words {
		clean := strings.Trim(word, ".,?!;:'\"")
		if score, ok := riskKeywords[clean]; ok && !seen[clean] {
			accumulated += score
			seen[clean] = true
		}
	}

	total := base + accumulated
	if total > 100 {
		total = 100
	}
	return total
}
