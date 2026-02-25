# Development Guide

This document covers the internal architecture and design decisions for `overthink`. If you're interested in contributing or understanding how the tool works under the hood, start here.

---

## Project Structure

```
overthink/
│
├── cmd/overthink/        # Entry point & CLI orchestration
├── internal/
│   ├── engine/           # Core analysis pipeline
│   │   ├── analyzer.go      (orchestrates all analysis stages)
│   │   ├── probability.go   (generates pseudo-statistics)
│   │   ├── risk.go          (keyword-weighted risk index)
│   │   ├── charts.go        (colored ASCII bar rendering)
│   │   ├── citations.go     (fabricated academic pools)
│   │   ├── color.go         (ANSI escape codes & formatting)
│   │   └── formatter.go     (io.Writer terminal output)
│   │
│   ├── ollama/           # Subprocess client
│   │   └── client.go        (os/exec + graceful fallback)
│   │
│   └── utils/            # Shared utilities
│       └── random.go        (time-seeded RNG)
│
├── go.mod               # No external dependencies
└── README.md            # User-facing documentation
```

---

## Design Philosophy

| Choice | Why |
|---|---|
| **No external dependencies** | Single static binary. No version drama. No `go.sum` hostage situations. |
| **Standard library only** | Go's `flag`, `fmt`, `os/exec`, `math/rand` are more than enough. |
| **Time-seeded randomness** | Every run is different. Overthinking fatigue is real; we combat it with variety. |
| **ANSI colors, always on** | It's 2026. Every terminal supports color. No `--color` flag nonsense. |
| **Subprocess over HTTP** | Immune to Ollama API changes. Works with any version. Forever. |
| **`io.Writer` based** | Decouples output from stdout. Test-friendly. Redirect anywhere. |
| **All content in code** | No config files, no assets, no external data. Pure offline-first Go. |

---

## Architecture Deep Dive

### Entry Point: `cmd/overthink/main.go`

Handles CLI flag parsing using the standard `flag` package (Cobra is overkill for one optional flag). Dispatches to either the local engine or Ollama mode.

### The Analysis Pipeline: `internal/engine/`

#### `analyzer.go`

The orchestrator. Calls all other modules in sequence:

1. **Title generation** — Uses dramatic prefixes/nouns + meaningful keywords from the question
2. **Summary generation** — Pick a template, fill it with fake confidence
3. **Probabilities** — Generate 3–5 outcomes that sum to exactly 100%
4. **Risk calculation** — Keyword-weighted index (0–100)
5. **Citations** — Pick 2–4 fabricated journals and years
6. **Conclusion** — Pick from a pool of theatrical finales
7. **Closing line** — Self-aware quip

All randomness is seeded from `time.Now().UnixNano()`, so every run is unique.

#### `probability.go`

Generates suspiciously precise percentages. Uses weighted random numbers normalized to 100.0%. The final entry gets the remainder to ensure exact closure.

#### `risk.go`

Keyword-weighted Risk Index. Maps ~40 anxiety-triggering words (ex, quit, dead, regret, etc.) to point values. Accumulates points from the question + a random baseline. Capped at 100.

#### `risk.go` → color mapping

Green (< 40), Yellow (40–69), Red (70+). Used by the formatter.

#### `charts.go`

Pure Unicode rendering. `filledBlock = "█"`, `emptyBlock = "░"`, `dividerChar = "─"`. Takes an ANSI color code and applies it to filled bars.

#### `citations.go`

Pool of ~20 fabricated journal names + author suffixes. Generates 2–4 entries with fake years (2008–2024).

#### `color.go`

ANSI escape code constants + helper functions. All colors applied unconditionally—modern terminals support them.

#### `formatter.go`

Renders the complete `AnalysisResult` to an `io.Writer`. Applies colors, formats sections, renders bars. Testable because it doesn't depend on `os.Stdout`.

### Ollama Integration: `internal/ollama/`

#### `client.go`

Constructs a full system prompt, pipes it via stdin to `ollama run <model>`, captures stdout. If anything fails—not installed, model missing, timeout—returns an error. The main CLI gracefully falls back.

### Utilities: `internal/utils/`

#### `random.go`

`NewRand()` returns a `*rand.Rand` seeded from the current nanosecond. Pool helpers: `PickString()`, `ShuffleStrings()`.

---

## Testing

Currently, there are no automated tests. To add test coverage, follow these patterns:

```bash
# Create a test file (e.g., probability_test.go)
# Use Go's built-in testing package and standard test conventions
go test ./...
```

### Test Ideas

If you want to contribute tests, these areas are good candidates:

- **Probability generation** — Verify percentages sum to 100.0% and stay in expected ranges
- **Risk calculation** — Confirm keyword detection and scoring logic
- **Title generation** — Ensure generated titles use meaningful question keywords
- **Random seeding** — Verify time-based generation produces different output on separate runs

---

## Contributing

The codebase is intentionally minimal and self-contained:

1. **Add features in logical places.** New outcome types? Add to pools in `analyzer.go`. New probability outcomes? Update `probability.go`.
2. **Keep dependencies at zero.** Use only the Go standard library. No external packages.
3. **Maintain the tone.** The tool is self-aware and theatrical. Keep dramatic language throughout.
4. **Test manually.** Run the tool with different questions to verify behavior before making changes.

---

## Performance Notes

The tool is designed for single invocations, not high throughput:

- **Startup time:** ~1ms (pure Go binary, no interpreted overhead)
- **Analysis time:** <1ms (pure computation, no I/O except final print)
- **Ollama calls:** 1–5 seconds per request (depends on model size and hardware)

No optimization is needed for typical use. If you're calling it 1000 times per second, you have different problems.

---

## Future Possibilities

Some ideas that maintain the design principles:

- **`--format json`** — Output structured data (without colors)
- **History file** — `~/.overthink/history.json` to persist past analyses
- **Custom pools** — Environment variables or config files (breaks offline-first principle; discouraged)
- **Animated spinner** — Fake "computing" for theatrical effect
- **Theme color schemes** — Different ANSI palettes for different moods

Any of these can be added without external dependencies.

---

## License

MIT
