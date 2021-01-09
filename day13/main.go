package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Point represent a position in the arcade cabinet's screen.
// y grows Down, x grows Right, the zero value is top-left tile.
type Point struct {
	x, y int64
}

// Tile is a square tiles on the screen grid.
type Tile int

//go:generate stringer -type=Tile
const (
	Empty Tile = iota
	Wall
	Block
	Paddle
	Ball
)

// Screen is the arcade cabinet's screen.
type Screen map[Point]Tile

// NewTile returns the Tile represented by the given Intcode along with true
// when the conversion succeeded, false otherwise.
func NewTile(i Intcode) (Tile, bool) {
	if i >= Intcode(Empty) && i <= Intcode(Ball) {
		return Tile(i), true
	}
	return Empty, false
}

// Count return the count of tiles in the screen matching t.
func (s Screen) Count(t Tile) int {
	n := 0
	for _, tile := range s {
		if tile == t {
			n++
		}
	}
	return n
}

// draw execute the given program and return the Screen state once it is
// executed along with any error encountered.
func draw(program []Intcode) (Screen, error) {
	c := Computer{
		input:  make(chan Intcode),
		output: make(chan Intcode),
	}
	close(c.input)

	halt := make(chan error)
	go func() {
		halt <- c.Execute(program)
		close(halt)
	}()

	screen := make(Screen)
	var i, x, y int64
	for {
		select {
		case err := <-halt:
			switch {
			case err != nil:
				return nil, err
			case i%3 == 1:
				return nil, errors.New("halted before y position output")
			case i%3 == 2:
				return nil, errors.New("halted before tile id output")
			}
			return screen, nil
		case o := <-c.output:
			switch {
			case i%3 == 0:
				if o < 0 {
					return nil, fmt.Errorf("%d: invalid x position", o)
				}
				x = int64(o)
			case i%3 == 1:
				if o < 0 {
					return nil, fmt.Errorf("%d: invalid y position", o)
				}
				y = int64(o)
			case i%3 == 2:
				tid, ok := NewTile(o)
				if !ok {
					return nil, fmt.Errorf("%d: invalid tile id", o)
				}
				screen[Point{x, y}] = tid
			}
		}
		i++
	}
}

// main compute and display the count of block tiles on the screen when the
// game exits.
func main() {
	program, err := Parse(os.Stdin)
	if err != nil {
		log.Fatalf("input error: %s\n", err)
	}

	screen, err := draw(program)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("There are %d block tiles on the screen when the game exits.\n", screen.Count(Block))
}

// Parse an Intcode program.
// It returns the parsed Intcode program and any read or convertion error
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
