package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Path
	}{
		{
			name:  "detailed example",
			input: "R8,U5,L5,D3\nU7,R6,D4,L4",
			want: []Path{
				{
					Step{Right, 8},
					Step{Up, 5},
					Step{Left, 5},
					Step{Down, 3},
				},
				{
					Step{Up, 7},
					Step{Right, 6},
					Step{Down, 4},
					Step{Left, 4},
				},
			},
		},
		{
			name:  "first example",
			input: "R75,D30,R83,U83,L12,D49,R71,U7,L72\nU62,R66,U55,R34,D71,R55,D58,R83",
			want: []Path{
				{
					Step{Right, 75},
					Step{Down, 30},
					Step{Right, 83},
					Step{Up, 83},
					Step{Left, 12},
					Step{Down, 49},
					Step{Right, 71},
					Step{Up, 7},
					Step{Left, 72},
				},
				{
					Step{Up, 62},
					Step{Right, 66},
					Step{Up, 55},
					Step{Right, 34},
					Step{Down, 71},
					Step{Right, 55},
					Step{Down, 58},
					Step{Right, 83},
				},
			},
		},
		{
			name:  "second example",
			input: "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51\nU98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
			want: []Path{
				{
					Step{Right, 98},
					Step{Up, 47},
					Step{Right, 26},
					Step{Down, 63},
					Step{Right, 33},
					Step{Up, 87},
					Step{Left, 62},
					Step{Down, 20},
					Step{Right, 33},
					Step{Up, 53},
					Step{Right, 51},
				},
				{
					Step{Up, 98},
					Step{Right, 91},
					Step{Down, 20},
					Step{Right, 16},
					Step{Down, 67},
					Step{Right, 40},
					Step{Up, 7},
					Step{Right, 15},
					Step{Up, 6},
					Step{Right, 7},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			paths, err := Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("Parse() error: %s", err)
			}
			if len(paths) != len(tc.want) {
				t.Fatalf("len(paths) = %d; want %d", len(paths), len(tc.want))
			}
			for i := range paths {
				if !PathEqual(paths[i], tc.want[i]) {
					t.Errorf("paths[%d] = %v; want %v", i, paths[i], tc.want[i])
				}
			}
		})
	}
}

func TestClosest(t *testing.T) {
	tests := []struct {
		name  string
		input string
		dist  int64
		steps int64
	}{
		{
			name:  "detailed example",
			input: "R8,U5,L5,D3\nU7,R6,D4,L4",
			dist:  6,
			steps: 30,
		},
		{
			name:  "first example",
			input: "R75,D30,R83,U83,L12,D49,R71,U7,L72\nU62,R66,U55,R34,D71,R55,D58,R83",
			dist:  159,
			steps: 610,
		},
		{
			name:  "second example",
			input: "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51\nU98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
			dist:  135,
			steps: 410,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			paths, err := Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Errorf("Parse() error: %s", err)
			}
			if len(paths) != 2 {
				t.Errorf("got %d wires; want 2", len(paths))
			}
			fst := NewWire(paths[0])
			snd := NewWire(paths[1])
			md, ms := Connect(fst, snd)
			if md != tc.dist {
				t.Errorf("distance = %d; want %d", md, tc.dist)
			}
			if ms != tc.steps {
				t.Errorf("steps = %d; want %d", ms, tc.steps)
			}
		})
	}
}

// PathEqual compare two wire paths and returns true if they are the same,
// false otherwise.
func PathEqual(ps, qs Path) bool {
	if len(ps) != len(qs) {
		return false
	}
	for i := range ps {
		if ps[i] != qs[i] {
			return false
		}
	}
	return true
}
