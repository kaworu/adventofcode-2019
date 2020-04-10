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

// LineOfSight represent the direction or angle from one Asteroid to another.
// The zero value represent the direction from an Asteroid to itself. A
// LineOfSight hold the invariant gcd(x, y) == 1.
type LineOfSight struct {
	x, y int
}

// BeamShot is a giant laser setting to vaporize the victim Asteroid.
type BeamShot struct {
	From, To Asteroid
	LineOfSight
	Distance int
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

// Vaporize returns a channel of Asteroid to be vaporized (in order) by a giant
// laser installed at the given station.
func Vaporize(station Asteroid, asteroids []Asteroid) <-chan Asteroid {
	victims := make(chan Asteroid)
	go func() {
		defer close(victims)
		// create a BeamShot for every asteroid that is going to be vaporized.
		// The ring represent the victims indexed by their shooting angle.
		ring := make(map[LineOfSight][]*BeamShot)
		// count how many asteroid we're shooting at, because asteroids may contain
		// the station which we obviously don't want to vaporize.
		n := 0
		for _, o := range asteroids {
			if o != station {
				shot := station.BeamShot(o)
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

		// build and sort the angles at which the laster is going to BeamShot.
		lines := make([]LineOfSight, 0, len(ring))
		for k := range ring {
			lines = append(lines, k)
		}
		sort.Slice(lines, func(i, j int) bool {
			return lines[i].Angle() < lines[j].Angle()
		})

		// vaporize the victims in order. We rotate through every BeamShot
		// angles and pick the closest Asteroid to be seen until there are none
		// left.
		for n > 0 {
			for _, l := range lines {
				if shots := ring[l]; len(shots) > 0 {
					victims <- shots[0].To
					ring[l] = shots[1:]
					n--
				}
			}
		}
	}()
	return victims
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
	fmt.Printf("%d other asteroids can be detected from %v,\n", n, a)

	victims := Vaporize(a, asteroids)
	for i := 0; i < 199; i++ {
		if _, ok := <-victims; !ok {
			log.Fatalf("only %d asteroid(s) vaporized; expected at least 200\n", i)
		}
	}
	vaporized := <-victims // the 200th
	fmt.Printf("and the 200th asteroid to be vaporized is at %v.\n", vaporized)
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
