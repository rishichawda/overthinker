package engine

import (
	"fmt"
	"strings"
)

const (
	// chartWidth is the total character width of all rendered bar charts.
	chartWidth = 40

	// filledBlock is the Unicode block character used for the filled bar portion.
	filledBlock = "\u2588"
	// emptyBlock is used for the unfilled portion of a bar.
	emptyBlock = "\u2591"
	// dividerChar is the box-drawing horizontal line used for section dividers.
	dividerChar = "\u2500"
)

// RenderRiskBar renders a colored horizontal bar for the Emotional Risk Index.
// fillColor is an ANSI color code applied to the filled portion of the bar.
// The label line shows the numeric score; the bar line shows the visual.
func RenderRiskBar(score int, fillColor string) string {
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	filled := (score * chartWidth) / 100
	empty := chartWidth - filled
	bar := fillColor + strings.Repeat(filledBlock, filled) + colorReset + dim(strings.Repeat(emptyBlock, empty))
	label := fmt.Sprintf("%sEmotional Risk Index: %s%d%s/100%s",
		colorBold,
		fillColor, score, colorReset+colorBold,
		colorReset)
	return label + "\n" + bar
}

// RenderProbabilityBars renders a compact bar chart for each probability entry.
// barColor is an ANSI color code applied to the filled portion of each bar.
func RenderProbabilityBars(probs []Probability, barColor string) string {
	var sb strings.Builder
	for _, p := range probs {
		filled := int((p.Percentage / 100.0) * float64(chartWidth))
		empty := chartWidth - filled
		bar := barColor + strings.Repeat(filledBlock, filled) + colorReset + dim(strings.Repeat(emptyBlock, empty))
		sb.WriteString(fmt.Sprintf("  %s%5.1f%%%s  %s  %s\n",
			barColor, p.Percentage, colorReset,
			bar, dim(p.Label)))
	}
	return sb.String()
}

// RenderDivider returns a horizontal divider line of the given character width.
func RenderDivider(width int) string {
	return strings.Repeat(dividerChar, width)
}
