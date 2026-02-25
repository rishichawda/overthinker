// Command overthink accepts a question and produces an absurdly over-analyzed
// response -- complete with fake statistics, fabricated citations, and theatrical
// philosophical conclusions.
//
// Usage:
//
//	overthink "Should I text my ex?"
//	overthink --thinker llama3 "Should I quit my job?"
//
// If --thinker is provided, the question is passed to a locally running Ollama
// instance. If Ollama is unavailable or fails, the built-in engine takes over.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rishichawda/overthinker/internal/engine"
	"github.com/rishichawda/overthinker/internal/ollama"
)

const usageText = `overthink -- a dramatic overanalysis engine

Usage:
  overthink [flags] "<your question>"

Flags:
  --thinker <model>   Use a local Ollama model (e.g. llama3, mistral)
                      Falls back to built-in engine if Ollama is unavailable.

Examples:
  overthink "Should I text my ex?"
  overthink "Is it too late to start coding?"
  overthink --thinker llama3 "Should I quit my job?"

If no question is provided, this message is printed and the program exits.
`

func main() {
	// Using the standard flag package: one optional flag does not justify Cobra.
	thinkerFlag := flag.String("thinker", "", "Ollama model name to use for analysis")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usageText)
	}

	flag.Parse()

	// Join all remaining arguments to support unquoted multi-word questions
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	question := strings.TrimSpace(strings.Join(args, " "))
	if question == "" {
		flag.Usage()
		os.Exit(1)
	}

	formatter := engine.NewFormatter(os.Stdout)

	// Ollama path
	if *thinkerFlag != "" {
		handleOllamaMode(question, *thinkerFlag, formatter)
		return
	}

	// Local engine path
	result := engine.Analyze(question)
	formatter.Print(result)
}

// handleOllamaMode attempts to query the specified Ollama model.
// On any failure it emits a graceful warning and falls back to the local engine.
func handleOllamaMode(question, modelName string, formatter *engine.Formatter) {
	client := ollama.NewClient(modelName)
	output, err := client.Query(question)
	if err != nil {
		formatter.PrintWarning(fmt.Sprintf("%v", err))
		result := engine.Analyze(question)
		formatter.Print(result)
		return
	}

	formatter.PrintOllamaOutput(modelName, output)
}
