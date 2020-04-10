package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {

	input := strings.Trim(`
.#..#
.....
#####
....#
...##
`, "\n")

	expected := []Asteroid{
		Asteroid{1, 0}, Asteroid{4, 0},
		Asteroid{0, 2}, Asteroid{1, 2}, Asteroid{2, 2}, Asteroid{3, 2}, Asteroid{4, 2},
		Asteroid{4, 3},
		Asteroid{3, 4}, Asteroid{4, 4},
	}

	decoded, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error: %s", err)
	}
	if len(decoded) != len(expected) {
		t.Fatalf("got %d asteroids; expected %d", len(decoded), len(expected))
	}
	for i := range expected {
		if decoded[i] != expected[i] {
			t.Errorf("asteroids[%d] = %v; expected %v", i, decoded[i], expected[i])
		}
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		Asteroid
		Detected int
	}{
		{Asteroid{1, 0}, 7},
		{Asteroid{4, 0}, 7},
		{Asteroid{0, 2}, 6},
		{Asteroid{1, 2}, 7},
		{Asteroid{2, 2}, 7},
		{Asteroid{3, 2}, 7},
		{Asteroid{4, 2}, 5},
		{Asteroid{4, 3}, 7},
		{Asteroid{3, 4}, 8},
		{Asteroid{4, 4}, 7},
	}

	all := make([]Asteroid, len(tests))
	for i, tc := range tests {
		all[i] = tc.Asteroid
	}

	for _, tc := range tests {
		if n := tc.Asteroid.Detect(all); n != tc.Detected {
			t.Errorf("%v.Detect() = %d; expected %d", tc.Asteroid, n, tc.Detected)
		}
	}
}

func TestBestLocation(t *testing.T) {
	tests := []struct {
		Map      string
		Best     Asteroid
		Detected int
	}{
		{
			Map: strings.Trim(`
......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####
`, "\n"),
			Best:     Asteroid{5, 8},
			Detected: 33,
		},
		{
			Map: strings.Trim(`
#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.
`, "\n"),
			Best:     Asteroid{1, 2},
			Detected: 35,
		},
		{
			Map: strings.Trim(`
.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..
`, "\n"),
			Best:     Asteroid{6, 3},
			Detected: 41,
		},
		{
			Map: strings.Trim(`
.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##
`, "\n"),
			Best:     Asteroid{11, 13},
			Detected: 210,
		},
	}

	for i, tc := range tests {
		asteroids, err := Parse(strings.NewReader(tc.Map))
		if err != nil {
			t.Errorf("tests[%d]: Parse() error: %s", i, err)
			continue
		}
		best, detected, err := BestLocation(asteroids)
		switch {
		case err != nil:
			t.Errorf("tests[%d]: BestLocation() error: %s", i, err)
		case best != tc.Best:
			t.Errorf("tests[%d]: best = %v; expected %v", i, best, tc.Best)
		case detected != tc.Detected:
			t.Errorf("tests[%d]: detected = %d; expected %v", i, detected, tc.Detected)
		}
	}
}

func TestVaporize(t *testing.T) {
	belt := strings.Trim(`
.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##
`, "\n")
	station := Asteroid{11, 13}
	nth := map[int]Asteroid{
		1:   Asteroid{11, 12},
		2:   Asteroid{12, 1},
		3:   Asteroid{12, 2},
		10:  Asteroid{12, 8},
		20:  Asteroid{16, 0},
		50:  Asteroid{16, 9},
		100: Asteroid{10, 16},
		199: Asteroid{9, 6},
		200: Asteroid{8, 2},
		201: Asteroid{10, 9},
		299: Asteroid{11, 1},
	}

	asteroids, err := Parse(strings.NewReader(belt))
	if err != nil {
		t.Fatalf("Parse() error: %s", err)
	}
	victims := Vaporize(station, asteroids)
	i := 0
	for victim := range victims {
		i++
		if expected, ok := nth[i]; ok && victim != expected {
			t.Errorf("laser.Vaporize() = %v; expected %v", victim, expected)
		}
	}
	if i != len(asteroids)-1 {
		t.Errorf("vaporized %v asteroids; expected %v", i, len(asteroids)-1)
	}
}
