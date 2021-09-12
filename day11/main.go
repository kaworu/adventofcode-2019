package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Headings
const (
	North = iota
	East
	South
	West
)

// Directions
const (
	Left  = iota // Left 90 degrees
	Right        // Right 90 degrees
)

// Colors
const (
	// Black is the staring color of all panels.
	Black = iota
	// White is the only color beside black that a Robot can paint a panel
	// with.
	White
)

// Heading is either North, East, South, or West.
type Heading uint8

// Direction is either Left or Right.
type Direction int8

// Color represent the different paint of our ship's panels.
type Color uint8

// Point represent a position in the grid.
// y grows Down, x grows Right, the zero value is the central port.
type Point struct {
	x, y int64
}

// Spacecraft represent our ship with its panels on the side.
type Spacecraft struct {
	panels map[Point]Color
}

// Robot represent the emergency hull painting robot.
type Robot struct {
	brain Computer
	Heading
	Point
}

// String implements Stringer for Heading.
func (h Heading) String() string {
	switch h {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	default:
		return "<invalid heading>"
	}
}

// String implements Stringer for Direction.
func (d Direction) String() string {
	switch d {
	case Left:
		return "Left"
	case Right:
		return "Right"
	default:
		return "<invalid direction>"
	}
}

// String implements Stringer for Color.
func (c Color) String() string {
	switch c {
	case White:
		return "White"
	case Black:
		return "Black"
	default:
		return "<invalid color>"
	}
}

// Turning returns the new heading when turning in the provided Direction.
func (h Heading) Turning(d Direction) Heading {
	switch d {
	case Left:
		return (h - 1 + 4) % 4
	case Right:
		return (h + 1) % 4
	default:
		return h
	}
}

// PointOfOrigin return the painting robot's starting point, the one at the
// center of the ship.
func PointOfOrigin() Point {
	return Point{x: 0, y: 0}
}

// Moving returns the landing point when moving in the provided Heading from p.
func (p Point) Moving(h Heading) Point {
	switch h {
	case North:
		return Point{p.x, p.y - 1}
	case East:
		return Point{p.x + 1, p.y}
	case South:
		return Point{p.x, p.y + 1}
	case West:
		return Point{p.x - 1, p.y}
	}
	// invalid heading means we're not moving.
	return p
}

// NewSpacecraft create a brand new ship with its panels all painted in Black.
func NewSpacecraft() *Spacecraft {
	return &Spacecraft{
		panels: make(map[Point]Color),
	}
}

// PaintedPanelCount returns the number of panels that were painted at least once.
func (ship *Spacecraft) PaintedPanelCount() int {
	return len(ship.panels)
}

// PaintItBlack paint every already painted panels in Black.
func (ship *Spacecraft) PaintItBlack() {
	for p := range ship.panels {
		ship.panels[p] = Black
	}
}

// String implements Stringer for Spacecraft displaying its panels.
func (ship Spacecraft) String() string {
	min, max := PointOfOrigin(), PointOfOrigin()
	for p := range ship.panels {
		if p.x < min.x {
			min.x = p.x
		} else if p.x > max.x {
			max.x = p.x
		}
		if p.y < min.y {
			min.y = p.y
		} else if p.y > max.y {
			max.y = p.y
		}
	}

	var buf bytes.Buffer
	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			c := ship.panels[Point{x: x, y: y}]
			if c == White {
				buf.WriteString("#")
			} else {
				buf.WriteString(".")
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

// NewRobot create a Robot running the given paint program.
func NewRobot() *Robot {
	return &Robot{
		brain: Computer{
			input:  make(chan Intcode),
			output: make(chan Intcode),
		},
	}
}

// Paint deploy the robot on the given Spacecraft in order to paint its panels
// following the provided program.
func (r *Robot) Paint(ship *Spacecraft, program []Intcode) error {
	// position and heading setup.
	r.Point = PointOfOrigin()
	r.Heading = North // The robot starts facing up

	// camera setup
	shoot := func() Color {
		return ship.panels[r.Point]
	}

	// brush setup
	paint := func(c Color) {
		ship.panels[r.Point] = c
	}

	// motor setup
	turn := func(d Direction) {
		r.Heading = r.Turning(d)
		r.Point = r.Moving(r.Heading)
	}

	// brain setup
	halt := make(chan error)
	go func() {
		halt <- r.brain.Execute(program)
		close(halt)
	}()

	nread := 0
	for {
		picture := shoot()
		select {
		case err := <-halt:
			return err
		case r.brain.input <- Intcode(picture):
			// the computer has read our picture shoot.
		case o := <-r.brain.output:
			nread++
			if nread%2 == 1 {
				// First, it will output a value indicating the color to paint
				// the panel the robot is over.
				paint(Color(o))
			} else {
				// Second, it will output a value indicating the direction the
				// robot should turn.
				turn(Direction(o))
			}
		}
	}
}

// main execute the latest Basic Operation Of System Test Intcode program given
// on stdin and display its keycode output.
func main() {
	program, err := Parse(os.Stdin)
	if err != nil {
		log.Fatalf("input error: %s\n", err)
	}

	ship := NewSpacecraft()
	robot := NewRobot()
	err = robot.Paint(ship, program)
	if err != nil {
		log.Fatalf("painting error: %s\n", err)
	}
	fmt.Printf("%d panels were paint at least once,\n", ship.PaintedPanelCount())

	// redo the painting with the robot starting on a white panel.
	ship.PaintItBlack()
	ship.panels[PointOfOrigin()] = White
	err = robot.Paint(ship, program)
	if err != nil {
		log.Fatalf("painting error: %s\n", err)
	}
	fmt.Printf("and here is your ship after the robot started on a white panel:\n%v", ship)
}

// Parse an Intcode program.
// It returns the parsed Intcode program and any read or conversion error
// encountered.
func Parse(r io.Reader) ([]Intcode, error) {
	var prog []Intcode
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanIntcodes)
	for scanner.Scan() {
		ic, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return nil, err
		}
		prog = append(prog, Intcode(ic))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return prog, nil
}

// ScanIntcodes is a split function for Scanner.
// It returns each Intcode of text.
func ScanIntcodes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Heavily inspired by ScanLines, the default Scanner split function.
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, ",\n"); i >= 0 {
		// We have a full Intcode
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated Intcode. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
