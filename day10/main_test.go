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
