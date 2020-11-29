package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

// FindCycleAndSimulate compute the given system's total energy after the
// requested number of steps along with the step count needed to reach a state
// that exactly matches a previous state.
func FindCycleAndSimulate(system *System, steps int) (energy, count int) {
	// map an axis state to its step number
	seen := make(map[string]int)
	// cycle step count for x, y, and z axis.
	var cycle [3]int
	// getters for each axis of a Vec3d, allowing us to loop through axis.
	axis := []func(Vec3d) int{
		func(v Vec3d) int { return v.x },
		func(v Vec3d) int { return v.y },
		func(v Vec3d) int { return v.z },
	}
	var buf bytes.Buffer
	for i := 0; ; i++ {
		// Compute the total system energy after the request number of steps.
		if i == steps {
			energy = system.TotalEnergy()
		}
		// compute each axis current state (as they are independent from each
		// other) in the search for a previously seen one so that we can detect
		// the cycle step count for this axis.
		for a, get := range axis {
			if cycle[a] > 0 {
				// we already know the cycle step count for this axis, so we
				// may skip the computation.
				continue
			}
			// start each state with the axis index. We do so in order to avoid
			// a collision with the same state from another axis.
			buf.Reset()
			buf.WriteString(strconv.Itoa(a))
			// write each moon's current axis position and velocity into the
			// state. After that we have a full picture of the current axis
			// state.
			for _, m := range system.Moons {
				buf.WriteString(strconv.Itoa(get(m.pos)))
				buf.WriteString(strconv.Itoa(get(m.vel)))
			}
			state := buf.String()
			if n, ok := seen[state]; ok {
				// we found the same state from a previous computation, now we
				// can get the cycle step count.
				cycle[a] = i - n
			} else {
				// store this state for later searches.
				seen[state] = i
			}
		}
		// we're done when we've found the cycle step count for each axis and
		// also having reached at least the request number of steps to compute
		// the total system energy.
		if cycle[0] > 0 && cycle[1] > 0 && cycle[2] > 0 && i >= steps {
			count = lcm3(cycle[0], cycle[1], cycle[2])
			break
		}
		// advance the simulation.
		system.Simulate(1)
	}
	return
}

// main compute and display the total energy in the system given on stdin after
// 1000 steps, along with the number of steps it take to reach the first state
// that exactly matches a previous state.
func main() {
	steps := 1000 // see README.md
	system, err := Parse(os.Stdin)
	if err != nil {
		log.Fatalf("input error: %s\n", err)
	}
	energy, count := FindCycleAndSimulate(system, steps)
	fmt.Printf("The total energy in the system after %d steps is %d,\n", steps, energy)
	fmt.Printf("and it take %d steps to reach a cycle.\n", count)
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

// gcd returns the greatest common divisor between a and b.
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// lcm3 returns the  least common multiple between a, b, and c.
func lcm3(a, b, c int) int {
	t := a * b / gcd(a, b)
	r := t * c / gcd(t, c)
	return r
}
