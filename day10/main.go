package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
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

// BeamShot is a GiantLaser setting to vaporize the victim asteroid.
type BeamShot struct {
	From, To Asteroid
	LineOfSight
	Distance int
}

// GiantLaser is the main tool for the only solution, i.e. complete
// vaporization.
type GiantLaser struct {
	base    Asteroid
	victims []Asteroid
}

// String returns "x,y" to match the README format for Asteroid positions.
func (a Asteroid) String() string {
	return fmt.Sprintf("%d,%d", a.x, a.y)
}

// Angle of this line of sight, up being 0 and increasing clockwise.
func (l LineOfSight) Angle() float64 {
	return 180 - (180/math.Pi)*math.Atan2(float64(l.x), float64(l.y))
}

// BeamShot compute and returns the shot needed to vaporize o from a.
func (a Asteroid) BeamShot(o Asteroid) BeamShot {
	dx, dy := o.x-a.x, o.y-a.y
	div := gcd(dx, dy)
	switch {
	case div < 0:
		div = -div
	case div == 0:
		return BeamShot{From: a, To: o}
	}
	los := LineOfSight{x: dx / div, y: dy / div}
	return BeamShot{From: a, To: o, LineOfSight: los, Distance: div}
}

// Detect returns the count of others asteroid in direct line of sight from a.
// When two asteroids are in the same line of sight they will be counted as one
// since the farthest one is "hidden behind" the closest one.
func (a Asteroid) Detect(asteroids []Asteroid) int {
	detected := make(map[LineOfSight]struct{}) // "set-like" map
	for _, o := range asteroids {
		if a != o {
			detected[a.BeamShot(o).LineOfSight] = struct{}{}
		}
	}
	return len(detected)
}

// NewGiantLaser build a GiantLaser at the given base in order to vaporize the
// given asteroids.
func NewGiantLaser(base Asteroid, asteroids []Asteroid) *GiantLaser {
	// create a BeamShot for every victim this GiantLaser is going to hit. The
	// ring represent the victims indexed by their shooting angle.
	ring := make(map[LineOfSight][]*BeamShot)
	// count how many asteroid we're shooting at, because asteroids may contain
	// the base which we obviously don't want to vaporize.
	n := 0
	for _, o := range asteroids {
		if o != base {
			shot := base.BeamShot(o)
			ring[shot.LineOfSight] = append(ring[shot.LineOfSight], &shot)
			n++
		}
	}

	// because the laser only has enough power to vaporize one asteroid at a
	// time before continuing its rotation, the victims are sorted by distance
	// (closest first) for every shoot angle.
	for _, shots := range ring {
		sort.Slice(shots, func(i, j int) bool {
			return shots[i].Distance < shots[j].Distance
		})
	}

	// build and sort the angles where we're going to have a BeamShot.
	lines := make([]LineOfSight, 0, len(ring))
	for k := range ring {
		lines = append(lines, k)
	}
	sort.Slice(lines, func(i, j int) bool {
		return lines[i].Angle() < lines[j].Angle()
	})

	// build the victims in the shooting order. We rotate through every
	// BeamShot angles and pick the closest Asteroid to be seen until there are
	// none left.
	victims := make([]Asteroid, 0, n)
	for len(victims) < n {
		for _, l := range lines {
			if shots := ring[l]; len(shots) > 0 {
				victims = append(victims, shots[0].To)
				ring[l] = shots[1:]
			}
		}
	}

	return &GiantLaser{base: base, victims: victims}
}

// Vaporize make the GiantLaser shot at an asteroid. It returns the vaporized
// Asteroid and true if there was one to vaporize, the Asteroid zero value and
// false otherwise.
func (l *GiantLaser) Vaporize() (Asteroid, bool) {
	if len(l.victims) == 0 {
		return Asteroid{}, false
	}
	vaporized := l.victims[0]
	l.victims = l.victims[1:]
	return vaporized, true
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
	laser := NewGiantLaser(a, asteroids)
	bet := laser.victims[199]
	fmt.Printf("%d other asteroids can be detected from %v,\n", n, a)
	fmt.Printf("and the 200th asteroid to be vaporized is at %v.\n", bet)
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
