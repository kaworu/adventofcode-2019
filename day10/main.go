package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// Asteroid is a celestial body in the asteroid belt.
type Asteroid struct {
	x, y int
}

// LineOfSight represent the direction or angle from one asteroid to another.
// The zero value represent the direction from an asteroid to itself. A
// LineOfSight hold the invariant gcd(x, y) == 1.
type LineOfSight struct {
	x, y int
}

// LineOfSight compute and returns the direction from a to o.
func (a Asteroid) LineOfSight(o Asteroid) LineOfSight {
	dx, dy := o.x-a.x, o.y-a.y
	div := gcd(dx, dy)
	switch {
	case div < 0:
		div = -div
	case div == 0:
		return LineOfSight{x: 0, y: 0}
	}
	return LineOfSight{x: dx / div, y: dy / div}
}

// Detect returns the count of others asteroid in direct line of sight from a.
// When two asteroids are in the same line of sight they will be counted as one
// since the farthest  one is "hidden behind" the closest one.
func (a Asteroid) Detect(others []Asteroid) int {
	detected := make(map[LineOfSight]struct{}) // "set-like" map
	for _, o := range others {
		if a != o {
			los := a.LineOfSight(o)
			detected[los] = struct{}{}
		}
	}
	return len(detected)
}

// BestLocation find the asteroid which would be the best place to build a new
// monitoring station. It returns the asteroid found to be the best, the count
// of other asteroids in line of sight from it, and an error when asteroids is
// empty.
func BestLocation(asteroids []Asteroid) (Asteroid, int, error) {
	max := 0
	loc := Asteroid{}
	if len(asteroids) == 0 {
		return loc, max, fmt.Errorf("empty slice argument")
	}
	for i, a := range asteroids {
		if n := a.Detect(asteroids); i == 0 || n > max {
			loc, max = a, n
		}
	}
	return loc, max, nil
}

// main find and display the best asteroid to build a new monitoring station
// and the count of other asteroids in line of sight from it.
func main() {
	asteroids, err := Parse(os.Stdin)
	if err != nil {
		log.Fatalf("input error: %s\n", err)
	}
	a, n, err := BestLocation(asteroids)
	if err != nil {
		log.Fatalf("BestLocation(): %s\n", err)
	}
	fmt.Printf("%d other asteroids can be detected from %v.\n", n, a)
}

// Parse the map of asteroid in the region. It returns the complete asteroid
// list and any read or parsing error encountered.
func Parse(r io.Reader) ([]Asteroid, error) {
	var asteroids []Asteroid
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanBytes)
	x, y := 0, 0
	for scanner.Scan() {
		switch c := scanner.Text(); c {
		case "#":
			asteroids = append(asteroids, Asteroid{x, y})
			fallthrough
		case ".":
			x++
		case "\n":
			y++
			x = 0
		default:
			return nil, fmt.Errorf("unexpected position at (%d,%d): %s", x, y, c)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return asteroids, nil
}

// gcd compute and returns the greatest common divisor between a and b using
// the Euclidean algorithm.
// See https://en.wikipedia.org/wiki/Euclidean_algorithm#Implementations
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
