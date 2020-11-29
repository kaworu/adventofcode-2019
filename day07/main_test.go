package main

import (
	"math/rand"
	"testing"
	"time"
)

var tests = []struct {
	name     string
	program  Memory
	sequence []Intcode
	want     Intcode
}{
	{
		name:     "part1 first example",
		program:  Memory{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0},
		sequence: []Intcode{4, 3, 2, 1, 0},
		want:     43210,
	},
	{
		name:     "part1 second example",
		program:  Memory{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0},
		sequence: []Intcode{0, 1, 2, 3, 4},
		want:     54321,
	},
	{
		name:     "part1 third example",
		program:  Memory{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33, 1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0},
		sequence: []Intcode{1, 0, 4, 3, 2},
		want:     65210,
	},
	{
		name:     "part2 first example",
		program:  Memory{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26, 27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5},
		sequence: []Intcode{9, 8, 7, 6, 5},
		want:     139629729,
	},
	{
		name:     "part2 second example",
		program:  Memory{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54, -5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4, 53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10},
		sequence: []Intcode{9, 7, 8, 5, 6},
		want:     18216,
	},
}

func TestFeedbackLoop(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := FeedbackLoop(tc.program, tc.sequence)
			if output != tc.want {
				t.Errorf("output = %v; want %v", output, tc.want)
			}
		})
	}
}

func TestHighestSignal(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := HighestSignal(tc.program, shuffled(tc.sequence))
			if output != tc.want {
				t.Errorf("output = %v; want %v", output, tc.want)
			}
		})
	}
}

// shuffled return a copy of the given slice with its element in a random
// order. Note that rand.Seed() must has been called before this function.
func shuffled(xs []Intcode) []Intcode {
	ys := make([]Intcode, len(xs))
	copy(ys, xs)
	rand.Shuffle(len(xs), func(i, j int) { ys[i], ys[j] = ys[j], ys[i] })
	return ys
}
