package main

import (
	"strings"
	"testing"
)

func TestPart1Examples(t *testing.T) {
	tests := []struct {
		name  string
		input string
		steps int
		want  int
	}{
		{
			name:  "first example",
			input: "<x=-1, y=0, z=2>\n<x=2, y=-10, z=-7>\n<x=4, y=-8, z=8>\n<x=3, y=5, z=-1>",
			steps: 10,
			want:  179,
		},
		{
			name:  "second example",
			input: "<x=-8, y=-10, z=0>\n<x=5, y=5, z=10>\n<x=2, y=-7, z=3>\n<x=9, y=-8, z=-3>",
			steps: 100,
			want:  1940,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Errorf("Parse() failed: %v", err)
				return
			}
			s.Simulate(tc.steps)
			got := s.TotalEnergy()
			if got != tc.want {
				t.Errorf("got %v total energy; want %v", got, tc.want)
			}
		})
	}
}

func TestPart2Examples(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		steps  int
		energy int
		count  int
	}{
		{
			name:   "first example",
			input:  "<x=-1, y=0, z=2>\n<x=2, y=-10, z=-7>\n<x=4, y=-8, z=8>\n<x=3, y=5, z=-1>",
			steps:  10,
			energy: 179,
			count:  2772,
		},
		{
			name:   "second example",
			input:  "<x=-8, y=-10, z=0>\n<x=5, y=5, z=10>\n<x=2, y=-7, z=3>\n<x=9, y=-8, z=-3>",
			steps:  100,
			energy: 1940,
			count:  4686774924,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Errorf("Parse() failed: %v", err)
				return
			}
			energy, count := FindCycleAndSimulate(s, tc.steps)
			if energy != tc.energy {
				t.Errorf("got %v total energy; want %v", energy, tc.energy)
			}
			if count != tc.count {
				t.Errorf("got %v cycle step count; want %v", count, tc.count)
			}
		})
	}
}
