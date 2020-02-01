package main

import (
	"math/rand"
	"testing"
	"time"
)

var tests = []struct {
	Name     string
	Program  Memory
	Sequence []Intcode
	Expected Intcode
}{
	{
		Name:     "first example",
		Program:  Memory{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0},
		Sequence: []Intcode{4, 3, 2, 1, 0},
		Expected: 43210,
	},
	{
		Name:     "second example",
		Program:  Memory{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0},
		Sequence: []Intcode{0, 1, 2, 3, 4},
		Expected: 54321,
	},
	{
		Name:     "third example",
		Program:  Memory{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33, 1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0},
		Sequence: []Intcode{1, 0, 4, 3, 2},
		Expected: 65210,
	},
}

func TestSeries(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			output := Series(tc.Program, tc.Sequence)
			if output != tc.Expected {
				t.Errorf("expected %v as output, got %v", tc.Expected, output)
			}
		})
	}
}

func TestHighestSignal(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			output := HighestSignal(tc.Program, Shuffled(tc.Sequence))
			if output != tc.Expected {
				t.Errorf("expected %v as output, got %v", tc.Expected, output)
			}
		})
	}
}

// Shuffled return a copy of the given slice with its element in a random
// order. Note that rand.Seed() must has been called before this function.
func Shuffled(xs []Intcode) []Intcode {
	ys := make([]Intcode, len(xs))
	copy(ys, xs)
	rand.Shuffle(len(xs), func(i, j int) { ys[i], ys[j] = ys[j], ys[i] })
	return ys
}
