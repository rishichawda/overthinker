// Command overthink accepts a question and produces an absurdly over-analyzed
// response -- complete with fake statistics, fabricated citations, and theatrical
// philosophical conclusions.
//
// Usage:
//
//	overthink "Should I text my ex?"
//	overthink --thinker llama3 "Should I quit my job?"
//
// If --thinker is provided, the question is sent to a locally running Ollama
// server via the HTTP API. If Ollama is unavailable or fails, the built-in
// local engine takes over.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rishichawda/overthinker/internal/engine"
	"github.com/rishichawda/overthinker/internal/local"
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
	thinkerFlag := flag.String("thinker", "", "Ollama model name to use for analysis")

	flag.Usage = func() { fmt.Fprint(os.Stderr, usageText) }
	flag.Parse()

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

	if *thinkerFlag != "" {
		runWithOllama(question, *thinkerFlag, formatter)
		return
	}

	result, _ := local.New().Analyze(question)
	formatter.Print(result)
}

// runWithOllama queries the Ollama model and renders the result with full
// formatting. On any error it warns and falls back to the local engine.
func runWithOllama(question, model string, formatter *engine.Formatter) {
	result, err := ollama.NewClient(model).Analyze(question)
	if err != nil {
		formatter.PrintWarning(fmt.Sprintf("%v", err))
		result, _ = local.New().Analyze(question)
		formatter.Print(result)
		return
	}

	formatter.PrintModelHeader(model)
	formatter.Print(result)
}
