package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

// main compute and display the total energy in the system given on stdin after
// 1000 steps.
func main() {
	steps := 1000
	system, err := Parse(os.Stdin)
	if err != nil {
		log.Fatalf("input error: %s\n", err)
	}
	system.Simulate(steps)
	energy := system.TotalEnergy()
	fmt.Printf("The total energy in the system after %d steps is %d.\n", steps, energy)
}

// Parse read the the position of the four largest moons of Jupiter (Io,
// Europa, Ganymede, and Callisto). It returns a system containing moons and
// any read or parsing error encountered.
func Parse(r io.Reader) (*System, error) {
	// strict regexp but that's good enough.
	posRegexp, err := regexp.Compile(`^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>$`)
	if err != nil {
		return nil, err
	}

	var s System
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		matches := posRegexp.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("invalid position: %s", line)
		}
		// ignore conversion errors as the regexp has already ensured we get a
		// number.
		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])
		z, _ := strconv.Atoi(matches[3])
		s.Moons = append(s.Moons, NewMoon(x, y, z))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &s, nil
}
