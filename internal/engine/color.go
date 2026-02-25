package engine

// ANSI escape codes for terminal coloring.
// Applied unconditionally. Modern terminals universally support them.
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorItalic = "\033[3m"

	colorBrightRed    = "\033[91m"
	colorBrightGreen  = "\033[92m"
	colorBrightYellow = "\033[93m"
	colorBrightCyan   = "\033[96m"
)

func bold(s string) string       { return colorBold + s + colorReset }
func dim(s string) string        { return colorDim + s + colorReset }
func italic(s string) string     { return colorItalic + s + colorReset }
func boldCyan(s string) string   { return colorBrightCyan + colorBold + s + colorReset }
func boldYellow(s string) string { return colorBrightYellow + colorBold + s + colorReset }
func dimCyan(s string) string    { return colorBrightCyan + colorDim + s + colorReset }

// riskFillColor returns the ANSI color for the filled portion of the risk bar.
// Green for calm, yellow for concerning, red for alarming.
func riskFillColor(score int) string {
	switch {
	case score >= 70:
		return colorBrightRed
	case score >= 40:
		return colorBrightYellow
	default:
		return colorBrightGreen
	}
}
