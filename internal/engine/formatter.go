package engine

import (
	"fmt"
	"io"
	"strings"
)

// Formatter handles all terminal output for the overthink engine.
// It writes to an io.Writer, making it testable and redirectable.
type Formatter struct {
	w io.Writer
}

// NewFormatter constructs a Formatter that writes to the given writer.
func NewFormatter(w io.Writer) *Formatter {
	return &Formatter{w: w}
}

// Print renders a complete AnalysisResult in strict output order:
//  1. DRAMATIC TITLE
//  2. Divider line
//  3. Executive Summary
//  4. Probability Analysis
//  5. Emotional Risk Index + ASCII bar
//  6. ASCII Visualization (probability bars)
//  7. Academic Citations
//  8. Grand Conclusion
//  9. Closing Line
func (f *Formatter) Print(result *AnalysisResult) {
	f.line("")
	f.linef("  %s", boldCyan(result.Title))
	f.line(dim(RenderDivider(len(result.Title) + 2)))
	f.line("")
	f.section("Executive Summary", result.Summary)
	f.line("")
	f.printProbabilities(result.Probabilities)
	f.line("")
	fillColor := riskFillColor(result.RiskIndex)
	f.line(RenderRiskBar(result.RiskIndex, fillColor))
	f.line("")
	f.printCitations(result.Citations)
	f.line("")
	f.section("Grand Conclusion", result.Conclusion)
	f.line("")
	f.linef("  %s", italic(bold("--> "+result.ClosingLine)))
	f.line("")
}

// PrintOllamaOutput renders attribution header then the raw Ollama response.
func (f *Formatter) PrintOllamaOutput(model, output string) {
	f.line("")
	f.linef("  %s", boldCyan("[ Thinker: "+model+" ]"))
	f.line(dim(RenderDivider(60)))
	f.line("")
	f.line(strings.TrimSpace(output))
	f.line("")
}

// PrintWarning prints a formatted warning message.
// Used when Ollama is unavailable and the engine falls back to local mode.
func (f *Formatter) PrintWarning(msg string) {
	f.linef("%s  Warning: %s", colorBrightYellow+colorBold, msg+colorReset)
	f.linef("   Falling back to the built-in overthinking engine.")
	f.line("")
}

// --- Private rendering helpers -----------------------------------------------

func (f *Formatter) section(heading, body string) {
	f.linef("%s:", boldYellow(heading))
	f.linef("  %s", body)
}

func (f *Formatter) printProbabilities(probs []Probability) {
	f.linef("%s:", boldYellow("Probability Analysis"))
	f.line("")
	for _, p := range probs {
		f.linef("  %s%5.1f%s%%  %s", colorBrightCyan, p.Percentage, colorReset, dim(p.Label))
	}
	f.line("")
	f.linef("  %s", dim("Visual Breakdown:"))
	f.line("")
	f.line(RenderProbabilityBars(probs, colorBrightCyan))
}

func (f *Formatter) printCitations(citations []Citation) {
	f.linef("%s:", boldYellow("Academic Citations"))
	for _, c := range citations {
		f.linef("  %s  %s", dimCyan(fmt.Sprintf("[%d]", c.Index)), c.Source)
	}
}

func (f *Formatter) line(s string) {
	fmt.Fprintln(f.w, s)
}

func (f *Formatter) linef(format string, args ...any) {
	fmt.Fprintf(f.w, format+"\n", args...)
}
