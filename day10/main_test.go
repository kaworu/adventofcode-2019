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

	want := []Asteroid{
		{1, 0}, {4, 0},
		{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2},
		{4, 3},
		{3, 4}, {4, 4},
	}

	decoded, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error: %s", err)
	}
	if len(decoded) != len(want) {
		t.Fatalf("got %d asteroids; want %d", len(decoded), len(want))
	}
	for i := range want {
		if decoded[i] != want[i] {
			t.Errorf("asteroids[%d] = %v; want %v", i, decoded[i], want[i])
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
			t.Errorf("%v.Detect() = %d; want %d", tc.Asteroid, n, tc.Detected)
		}
	}
}

func TestBestLocation(t *testing.T) {
	tests := []struct {
		plan     string
		best     Asteroid
		detected int
	}{
		{
			plan: strings.Trim(`
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
			best:     Asteroid{5, 8},
			detected: 33,
		},
		{
			plan: strings.Trim(`
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
			best:     Asteroid{1, 2},
			detected: 35,
		},
		{
			plan: strings.Trim(`
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
			best:     Asteroid{6, 3},
			detected: 41,
		},
		{
			plan: strings.Trim(`
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
			best:     Asteroid{11, 13},
			detected: 210,
		},
	}

	for i, tc := range tests {
		asteroids, err := Parse(strings.NewReader(tc.plan))
		if err != nil {
			t.Errorf("tests[%d]: Parse() error: %s", i, err)
			continue
		}
		best, detected, err := BestLocation(asteroids)
		switch {
		case err != nil:
			t.Errorf("tests[%d]: BestLocation() error: %s", i, err)
		case best != tc.best:
			t.Errorf("tests[%d]: best = %v; want %v", i, best, tc.best)
		case detected != tc.detected:
			t.Errorf("tests[%d]: detected = %d; want %v", i, detected, tc.detected)
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
		1:   {11, 12},
		2:   {12, 1},
		3:   {12, 2},
		10:  {12, 8},
		20:  {16, 0},
		50:  {16, 9},
		100: {10, 16},
		199: {9, 6},
		200: {8, 2},
		201: {10, 9},
		299: {11, 1},
	}

	asteroids, err := Parse(strings.NewReader(belt))
	if err != nil {
		t.Fatalf("Parse() error: %s", err)
	}
	victims := Vaporize(station, asteroids)
	i := 0
	for victim := range victims {
		i++
		if want, ok := nth[i]; ok && victim != want {
			t.Errorf("laser.Vaporize() = %v; want %v", victim, want)
		}
	}
	if i != len(asteroids)-1 {
		t.Errorf("vaporized %v asteroids; want %v", i, len(asteroids)-1)
	}
}
