package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	// Up is the north direction in the grid.
	Up = 0
	// Right is the east direction in the grid.
	Right = iota
	// Down is the south direction in the grid.
	Down = iota
	// Left is the west direction in the grid.
	Left = iota
)

// Point represent a position in the grid.
// y grows Up, x grows Right.
type Point struct {
	x, y int64
}

// Grid is where the wires are connected.
// The values are the count of wire at the given point.
type Grid map[Point]uint64

// Step represent a path component.
type Step struct {
	Direction byte // Direction is either Up, Right, Down or Left.
	Count     int  // Count is the step's number of port.
}

// Path represent a circuit wire connection to a central port on the grid.
type Path struct {
	ID    uint64 // ID is the wire identifier.
	Steps []Step /// Steps is the path description of the wire in the grid.
}

// Add returns p + other.
func (p Point) Add(other Point) Point {
	return Point{x: p.x + other.x, y: p.y + other.y}
}

// Distance compute and returns the Manhattan distance between two points.
func (p Point) Distance(other Point) int64 {
	return abs(p.x-other.x) + abs(p.y-other.y)
}

// Closest find the point in the given slice having the least distance from p.
// It return the point found and its distance relative to p. If the slice is
// empty, the returned point is undefined and distance is -1.
func (p Point) Closest(others []Point) (Point, int64) {
	var c Point                  // closest point
	var cd int64 = math.MaxInt64 // closest distance
	if len(others) == 0 {
		return c, -1
	}
	for _, other := range others {
		d := p.Distance(other)
		if d < cd {
			cd = d
			c = other
		}
	}
	return c, cd
}

// CentralPort returns the point in the grid where from where all wire
// connections begins.
func (grid Grid) CentralPort() Point {
	// The central port is the Point zero value for any Grid (by convention).
	return Point{x: 0, y: 0}
}

// Connect place a given wire path into the grid.
// It returns the points where the wire has crossed another already connected
// wire, in the order encountered.
func (grid Grid) Connect(path Path) []Point {
	var intersections []Point
	p := grid.CentralPort()
	for _, step := range path.Steps {
		var delta Point
		switch step.Direction {
		case Up:
			delta.y = 1
		case Right:
			delta.x = 1
		case Down:
			delta.y = -1
		case Left:
			delta.x = -1
		}
		for i := 0; i < step.Count; i++ {
			p = p.Add(delta)
			grid[p] |= path.ID
			if grid[p] != path.ID {
				intersections = append(intersections, p)
			}
		}
	}
	return intersections
}

// main compute and display the Manhattan distance from the central port to the
// closest intersection of the wires description given on stdin.
func main() {
	grid := make(Grid)
	wires, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}
	grid.Connect(wires[0])
	intersections := grid.Connect(wires[1])
	_, distance := grid.CentralPort().Closest(intersections)
	fmt.Printf("The Manhattan distance fron the central port to the closest intersection is %d.\n", distance)
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
			path.Steps = append(path.Steps, step)
		}
		paths = append(paths, path)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	for i := range paths {
		paths[i].ID = 1 << i
	}
	return paths, nil
}

// parseStep is a parsing helper for Parse.
// It parse and returns one step any parsing error encountered.
func parseStep(s string) (Step, error) {
	var step Step

	// the smallest step would be something like U1
	if len(s) < 2 {
		return step, fmt.Errorf("step too short")
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
		return step, fmt.Errorf("unrecognized direction %c", s[0])
	}

	// parse the step count
	i, err := strconv.Atoi(s[1:])
	if err != nil {
		return step, err
	}
	if i < 0 {
		return step, fmt.Errorf("negative step count")
	}
	step.Count = i

	return step, nil
}

// abs compute and returns the absolute value of n.
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
