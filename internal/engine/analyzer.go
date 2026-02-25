package engine

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/rishichawda/overthinker/internal/utils"
)

// AnalysisResult is the complete, peer-reviewed output of the overthink engine.
// Every field is populated deterministically from the input question.
type AnalysisResult struct {
	Title         string
	Summary       string
	Probabilities []Probability
	RiskIndex     int
	Citations     []Citation
	Conclusion    string
	ClosingLine   string
}

// --- Title Generation --------------------------------------------------------

var dramaticPrefixes = []string{
	"THE INEVITABLE",
	"THE CATASTROPHIC",
	"THE UNRESOLVED",
	"THE IRREVERSIBLE",
	"THE DEEPLY ALARMING",
	"THE STATISTICALLY SIGNIFICANT",
	"THE EXISTENTIALLY CHARGED",
	"THE CHRONICALLY UNRESOLVED",
	"THE QUIETLY DEVASTATING",
	"THE ACADEMICALLY CONCERNING",
	"THE PERENNIALLY UNFINISHED",
	"THE SUSPICIOUSLY FAMILIAR",
	"THE UNCOMFORTABLY RELATABLE",
	"THE STRUCTURALLY INEVITABLE",
}

var dramaticNouns = []string{
	"EMOTIONAL CASCADE",
	"COGNITIVE SPIRAL",
	"EXISTENTIAL TRAJECTORY",
	"PSYCHOLOGICAL UNDERTOW",
	"DECISION VORTEX",
	"ANALYTICAL PARADOX",
	"TEMPORAL RECKONING",
	"NEUROLOGICAL EVENT",
	"PHILOSOPHICAL QUANDARY",
	"INTERNAL MONOLOGUE",
	"CONSEQUENCE MATRIX",
	"ANXIETY FEEDBACK LOOP",
	"UNCERTAINTY GRADIENT",
	"NARRATIVE ARC",
	"RISK TOPOLOGY",
}

// stopWords are common words excluded from title keyword extraction.
var stopWords = map[string]bool{
	"the": true, "and": true, "for": true, "are": true, "but": true,
	"not": true, "you": true, "all": true, "can": true, "had": true,
	"her": true, "was": true, "one": true, "our": true, "out": true,
	"get": true, "has": true, "him": true, "his": true, "how": true,
	"its": true, "may": true, "now": true, "see": true, "two": true,
	"who": true, "did": true, "does": true, "any": true, "too": true,
	"that": true, "with": true, "this": true, "from": true, "they": true,
	"will": true, "have": true, "been": true, "into": true, "your": true,
	"when": true, "what": true, "more": true, "also": true, "than": true,
	"then": true, "some": true, "even": true, "just": true, "like": true,
	"over": true, "such": true, "here": true, "very": true, "much": true,
}

// generateTitle constructs an ALL-CAPS title from the input question.
// Extracts meaningful content words to create context-specific drama.
func generateTitle(question string, rng *rand.Rand) string {
	prefix := utils.PickString(rng, dramaticPrefixes)
	noun := utils.PickString(rng, dramaticNouns)

	words := strings.Fields(strings.ToLower(question))
	var meaningful []string
	for _, w := range words {
		clean := strings.Trim(w, ".,?!;:'\"")
		if len(clean) > 3 && !stopWords[clean] {
			meaningful = append(meaningful, strings.ToUpper(clean))
		}
	}

	if len(meaningful) == 0 {
		return fmt.Sprintf("%s %s OF THIS SITUATION", prefix, noun)
	}

	max := 4
	if len(meaningful) < max {
		max = len(meaningful)
	}
	subject := strings.Join(meaningful[:max], " ")
	return fmt.Sprintf("%s %s OF %s", prefix, noun, subject)
}

// --- Summary Generation ------------------------------------------------------

var summaryTemplates = []string{
	"After exhaustive cognitive simulation spanning 847 theoretical scenarios, the system has identified measurable turbulence in your current trajectory.",
	"A thorough multi-pass analysis reveals structural instability in the decision space surrounding this inquiry. The data is not encouraging.",
	"Cross-referencing your question against seventeen known behavioral archetypes, the system has flagged a statistically non-trivial probability of regret.",
	"Initial triage of this question triggered three separate alarm protocols. The situation has been escalated to the Dramatic Analysis Unit.",
	"Preliminary modeling indicates this question belongs to a well-documented category of decisions that humans make, reconsider, and then make again.",
	"The system has processed your inquiry using an advanced cascade of speculative heuristics. The results are both definitive and deeply ambiguous.",
	"Upon reflection -- 0.003 seconds of it -- the analytical engine has concluded that this question deserves far more attention than you've given it.",
	"Your question was run against the full corpus of human second-guessing. Several concerning patterns emerged immediately.",
	"After consulting internal uncertainty tables and applying a proprietary regret coefficient, a risk profile has been assembled. You won't love it.",
	"The cognitive simulation completed successfully. The news is mixed. The emotional implications are not.",
}

func generateSummary(rng *rand.Rand) string {
	return utils.PickString(rng, summaryTemplates)
}

// --- Conclusion Generation ---------------------------------------------------

var conclusionTemplates = []string{
	"Historical precedent strongly suggests you will proceed regardless of these findings. The system respects your autonomy and documents its objections.",
	"All available evidence points toward a path you've already emotionally chosen. This report exists to provide intellectual cover for that choice.",
	"The analysis is complete. The conclusion is inevitable. The action you take will be the one you were always going to take.",
	"Based on prior behavioral patterns across comparable datasets, the outcome of this decision was determined approximately six minutes before you ran this command.",
	"While the risk index is elevated, humans have historically proceeded under far worse conditions. This is both reassuring and alarming.",
	"The system recommends caution, restraint, and careful deliberation. The system acknowledges these recommendations will be ignored within 48 hours.",
	"After extensive analysis, the most scientifically defensible conclusion is: it depends. On things you haven't told us. And possibly on Mercury.",
	"This report has been generated. The implications have been flagged. The consequences remain, as always, entirely your responsibility.",
	"The data suggests two equally valid paths forward. You already know which one you'll take. So does the system.",
	"In the fullness of time, this decision will seem either obviously correct or obviously catastrophic. The system looks forward to being cited either way.",
}

func generateConclusion(rng *rand.Rand) string {
	return utils.PickString(rng, conclusionTemplates)
}

// --- Closing Line Generation -------------------------------------------------

var closingLines = []string{
	"You opened the chat window before running this command, didn't you?",
	"The system notes this is your third overthought decision this week. Statistically speaking, that's fine.",
	"This report will self-justify in approximately 72 hours.",
	"For what it's worth: the fact that you asked means you already know the answer.",
	"The system wishes you clarity, but expects you'll settle for validation.",
	"Proceed with caution. Or don't. The system will generate a report either way.",
	"If this were easy, you wouldn't need a dramatic analysis engine. You're welcome.",
	"Consider this report peer-reviewed by everyone who has ever been in your situation.",
	"The system has done its part. The rest is, unfortunately, up to you.",
	"A follow-up report is available whenever you spiral again. The system will be here.",
	"You already know what you're going to do. This report told you it was okay.",
	"Whatever you decide, the system supports you -- and will absolutely say 'I told you so.'",
	"Take a breath. Then do the thing you were going to do anyway. That's all any of us can do.",
	"Overthinking: complete. Action: TBD by the most chaotic part of your brain.",
	"The system detected 3 instances of the word 'should' in your future internal monologue. You're going to be fine.",
}

func generateClosingLine(rng *rand.Rand) string {
	return utils.PickString(rng, closingLines)
}

// --- Public API --------------------------------------------------------------

// Analyze is the primary entry point for the overthink engine.
// It accepts the raw user question and returns a fully populated AnalysisResult.
// Each invocation uses a fresh random seed, so output varies on every run.
func Analyze(question string) *AnalysisResult {
	rng := utils.NewRand()

	return &AnalysisResult{
		Title:         generateTitle(question, rng),
		Summary:       generateSummary(rng),
		Probabilities: GenerateProbabilities(rng),
		RiskIndex:     CalculateRiskIndex(question, rng),
		Citations:     GenerateCitations(rng),
		Conclusion:    generateConclusion(rng),
		ClosingLine:   generateClosingLine(rng),
	}
}
