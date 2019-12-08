package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		Name     string
		Input    string
		Expected []Wire
	}{
		{
			Name:  "detailed example",
			Input: "R8,U5,L5,D3\nU7,R6,D4,L4",
			Expected: []Wire{
				Wire{
					ID: 0x1,
					Steps: []Step{
						Step{Right, 8},
						Step{Up, 5},
						Step{Left, 5},
						Step{Down, 3},
					},
				},
				Wire{
					ID: 0x2,
					Steps: []Step{
						Step{Up, 7},
						Step{Right, 6},
						Step{Down, 4},
						Step{Left, 4},
					},
				},
			},
		},
		{
			Name:  "first example",
			Input: "R75,D30,R83,U83,L12,D49,R71,U7,L72\nU62,R66,U55,R34,D71,R55,D58,R83",
			Expected: []Wire{
				Wire{
					ID: 0x1,
					Steps: []Step{
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
				},
				Wire{
					ID: 0x2,
					Steps: []Step{
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
		},
		{
			Name:  "second example",
			Input: "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51\nU98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
			Expected: []Wire{
				Wire{
					ID: 0x1,
					Steps: []Step{
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
				},
				Wire{
					ID: 0x2,
					Steps: []Step{
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
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			wires, err := Parse(strings.NewReader(tc.Input))
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
			if !WireEqual(wires, tc.Expected) {
				t.Errorf("got %v, expected %v", wires, tc.Expected)
			}
		})
	}
}

func TestConnect(t *testing.T) {
	tests := []struct {
		Name     string
		Wires    []Wire
		Expected []Point
	}{
		{
			Name: "detailed example",
			Wires: []Wire{
				Wire{
					ID: 0x1,
					Steps: []Step{
						Step{Right, 8},
						Step{Up, 5},
						Step{Left, 5},
						Step{Down, 3},
					},
				},
				Wire{
					ID: 0x2,
					Steps: []Step{
						Step{Up, 7},
						Step{Right, 6},
						Step{Down, 4},
						Step{Left, 4},
					},
				},
			},
			Expected: []Point{
				Point{x: 6, y: 5}, Point{x: 3, y: 3},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			grid := make(Grid)
			intersections := grid.Connect(tc.Wires[0])
			if len(intersections) != 0 {
				t.Errorf("Unexpected intersections at first wire connection")
			}
			intersections = grid.Connect(tc.Wires[1])
			if !PointsEqual(intersections, tc.Expected) {
				t.Errorf("got %v, expected %v", intersections, tc.Expected)
			}
		})
	}
}

func TestClosest(t *testing.T) {
	tests := []struct {
		Name     string
		Wires    []Wire
		Expected int64 // distance
	}{
		{
			Name: "detailed example",
			Wires: []Wire{
				Wire{
					ID: 0x1,
					Steps: []Step{
						Step{Right, 8},
						Step{Up, 5},
						Step{Left, 5},
						Step{Down, 3},
					},
				},
				Wire{
					ID: 0x2,
					Steps: []Step{
						Step{Up, 7},
						Step{Right, 6},
						Step{Down, 4},
						Step{Left, 4},
					},
				},
			},
			Expected: 6,
		},
		{
			Name: "first example",
			Wires: []Wire{
				Wire{
					ID: 0x1,
					Steps: []Step{
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
				},
				Wire{
					ID: 0x2,
					Steps: []Step{
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
			Expected: 159,
		},
		{
			Name: "second example",
			Wires: []Wire{
				Wire{
					ID: 0x1,
					Steps: []Step{
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
				},
				Wire{
					ID: 0x2,
					Steps: []Step{
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
			Expected: 135,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			grid := make(Grid)
			grid.Connect(tc.Wires[0])
			intersections := grid.Connect(tc.Wires[1])
			_, dist := grid.CentralPort().Closest(intersections)
			if dist != tc.Expected {
				t.Errorf("got %v, expected %v", dist, tc.Expected)
			}
		})
	}
}

// WireEqual compare two wires and returns true if they are the same, false
// otherwise.
func WireEqual(xs, ys []Wire) bool {
	if len(xs) != len(ys) {
		return false
	}
	for i, x := range xs {
		y := ys[i]
		// ID field comparison
		if x.ID != y.ID {
			return false
		}
		// Steps field comparison
		if len(x.Steps) != len(y.Steps) {
			return false
		}
		for j, s := range x.Steps {
			if s != y.Steps[j] {
				return false
			}
		}
	}
	return true
}

// PointsEqual compare two Point slices and returns true if they are equal,
// false otherwise.
func PointsEqual(ps, qs []Point) bool {
	if len(ps) != len(qs) {
		return false
	}
	for i, p := range ps {
		q := qs[i]
		if p.x != q.x || p.y != q.y {
			return false
		}
	}
	return true
}
