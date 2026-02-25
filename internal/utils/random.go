package utils

import (
	"math/rand"
	"time"
)

// NewRand creates a new random source seeded from the current time in nanoseconds.
// Every invocation produces a different sequence, so each analysis is unique.
func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// PickString selects a pseudo-random element from a string slice.
func PickString(rng *rand.Rand, pool []string) string {
	if len(pool) == 0 {
		panic("overthink: cannot select from empty pool")
	}
	return pool[rng.Intn(len(pool))]
}

// ShuffleStrings returns a new slice with elements in pseudo-random order.
// The original slice is not modified.
func ShuffleStrings(rng *rand.Rand, s []string) []string {
	result := make([]string, len(s))
	copy(result, s)
	rng.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}
