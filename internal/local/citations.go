package local

import (
	"fmt"
	"math/rand"

	"github.com/rishichawda/overthinker/internal/engine"
	"github.com/rishichawda/overthinker/internal/utils"
)

// journalNames is the authoritative pool of imaginary academic publications.
var journalNames = []string{
	"Journal of Existential Hesitation",
	"International Review of Questionable Decisions",
	"Proceedings of the Annual Regret Symposium",
	"Journal of Romantic Miscalculation",
	"Quarterly Bulletin of Applied Catastrophizing",
	"Annals of Unnecessary Second-Guessing",
	"Transactions on Cognitive Overload",
	"Institute for Advanced Overanalysis",
	"Review of Premature Conclusions",
	"Journal of Speculative Self-Sabotage",
	"Archives of Temporal Panic",
	"Reports on Unresolved Ambiguity",
	"Compendium of Midnight Decisions",
	"Survey of Avoidant Coping Strategies",
	"Journal of Theoretical What-Ifs",
	"Bulletin of the Society for Spiraling Thoughts",
	"Proceedings on Human Indecision (Special Issue)",
	"Cambridge Handbook of Feelings You Cannot Name",
	"Oxford Review of Things You Almost Said",
	"Wiley Encyclopedia of Overthought Outcomes",
}

// authorSuffixes provides fictitious author credentials for maximum credibility.
var authorSuffixes = []string{
	"et al.",
	"& Associates",
	"(Independent Research Division)",
	"(Posthumous Edition)",
	"(Retracted, then re-instated)",
	"(Peer-reviewed by one very tired colleague)",
}

// generateCitations produces 2-4 fabricated academic citations.
func generateCitations(rng *rand.Rand) []engine.Citation {
	count := 2 + rng.Intn(3)

	shuffled := utils.ShuffleStrings(rng, journalNames)
	selected := shuffled[:count]

	citations := make([]engine.Citation, count)
	for i, journal := range selected {
		year := 2008 + rng.Intn(17)
		suffix := utils.PickString(rng, authorSuffixes)
		citations[i] = engine.Citation{
			Index:  i + 1,
			Source: fmt.Sprintf("%s %s (%d)", journal, suffix, year),
		}
	}
	return citations
}
