package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	// Up is the north direction in the grid.
	Up = iota
	// Right is the east direction in the grid.
	Right
	// Down is the south direction in the grid.
	Down
	// Left is the west direction in the grid.
	Left
)

// Step represent a path component.
type Step struct {
	Direction uint8 // Direction is either Up, Right, Down or Left.
	Count     int64 // Count is the step's number of port.
}

// Path represent a circuit wire connection description.
type Path []Step

// Point represent a position in the grid.
// y grows Up, x grows Right, the zero value is the central port.
type Point struct {
	x, y int64
}

// Segment are straight connections between two points.
type Segment struct {
	From, To   Point
	xmin, xmax int64
	ymin, ymax int64
	steps      int64
}

// Wire represent a wire connected into the grid from the central port.
type Wire []*Segment

// CentralPort returns the point in the grid where from where all wire
// connections begins.
func CentralPort() Point {
	// By convention, the central port is the Point zero value.
	return Point{x: 0, y: 0}
}

// Add returns p + other.
func (p Point) Add(other Point) Point {
	return Point{x: p.x + other.x, y: p.y + other.y}
}

// Distance compute and returns the Manhattan distance between two points.
func (p Point) Distance(other Point) int64 {
	return abs(p.x-other.x) + abs(p.y-other.y)
}

// NewSegment create a new Segment given the starting and ending points.
func NewSegment(from, to Point) *Segment {
	seg := &Segment{From: from, To: to}
	if from.x < to.x {
		seg.xmin, seg.xmax = from.x, to.x
	} else {
		seg.xmin, seg.xmax = to.x, from.x
	}
	if from.y < to.y {
		seg.ymin, seg.ymax = from.y, to.y
	} else {
		seg.ymin, seg.ymax = to.y, from.y
	}
	seg.steps = from.Distance(to)
	return seg
}

// IntersectWith check if the other Segment and seg share a point. If they do,
// the intersection point and true is returned. Otherwise the Point zero value
// and false is returned.
func (seg *Segment) IntersectWith(other *Segment) (Point, bool) {
	switch {
	case seg.xmin <= other.xmin && seg.xmax >= other.xmin &&
		seg.ymin <= other.ymax && seg.ymin >= other.ymin:
		return Point{x: other.xmin, y: seg.ymin}, true
	case other.xmin <= seg.xmin && other.xmax >= seg.xmin &&
		other.ymin <= seg.ymax && other.ymin >= seg.ymin:
		return Point{x: seg.xmin, y: other.ymin}, true
	}
	return Point{}, false
}

// NewWire place a given wire path into the grid and return the resulting Wire.
func NewWire(path Path) Wire {
	var wire Wire
	start := CentralPort() // the current position, starting at the central port.
	for _, step := range path {
		stop := start
		switch step.Direction {
		case Up:
			stop.y += step.Count
		case Right:
			stop.x += step.Count
		case Down:
			stop.y -= step.Count
		case Left:
			stop.x -= step.Count
		}
		wire = append(wire, NewSegment(start, stop))
		start = stop
	}
	return wire
}

// Connect link a couple of wire on the grid. It returns the the Manhattan
// distance from the central port to the closest intersection (md) and the
// fewest combined steps the wires must take to reach an intersection (ms).
func Connect(a, b Wire) (md, ms int64) {
	md = -1
	ms = -1
	cp := CentralPort()
	var astep int64 = 0
	for _, aseg := range a {
		var bstep int64 = 0
		for _, bseg := range b {
			// we omit the intersection at the central port, hence p != cp.
			if p, ok := aseg.IntersectWith(bseg); ok && p != cp {
				// min distance
				d := cp.Distance(p)
				if md == -1 || d < md {
					md = d
				}
				// min combined step
				s := astep + aseg.From.Distance(p) +
					bstep + bseg.From.Distance(p)
				if ms == -1 || s < ms {
					ms = s
				}
			}
			bstep += bseg.steps
		}
		astep += aseg.steps
	}
	return
}

// main compute and display the Manhattan distance from the central port to the
// closest intersection of the wires description given on stdin.
func main() {
	paths, err := Parse(os.Stdin)
	if err != nil {
		log.Fatalf("input error: %s\n", err)
	}
	fst := NewWire(paths[0])
	snd := NewWire(paths[1])
	md, ms := Connect(fst, snd)
	fmt.Printf("The Manhattan distance fron the central port to the closest intersection is %v,\n", md)
	fmt.Printf("and the fewest combined steps the wires must take to reach an intersection is %v.\n", ms)
}

// Parse a couple of wire paths.
// It returns the parsed paths and any read or parsing error encountered.
func Parse(r io.Reader) ([]Path, error) {
	var paths []Path
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var path Path
		line := scanner.Text()
		for _, part := range strings.Split(line, ",") {
			step, err := parseStep(part)
			if err != nil {
				return nil, err
			}
			path = append(path, step)
		}
		paths = append(paths, path)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return paths, nil
}

// parseStep is a parsing helper for Parse.
// It parse and returns one step any parsing error encountered.
func parseStep(s string) (Step, error) {
	var step Step

	// the smallest step would be something like U1
	if len(s) < 2 {
		return step, fmt.Errorf("step too short: %s", s)
	}

	// parse the direction
	switch s[0] {
	case 'U':
		step.Direction = Up
	case 'R':
		step.Direction = Right
	case 'D':
		step.Direction = Down
	case 'L':
		step.Direction = Left
	default:
		return step, fmt.Errorf("unrecognized direction: %c", s[0])
	}

	// parse the step count
	i, err := strconv.ParseUint(s[1:], 10, 63) // 63 bit size fit in int64
	if err != nil {
		return step, err
	}
	step.Count = int64(i)

	return step, nil
}

// abs compute and returns the absolute value of n.
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
