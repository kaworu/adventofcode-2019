package main

import "testing"

func TestDraw(t *testing.T) {
	// This program simply output the sequence 1,2,3,6,5,4
	program := []Intcode{
		104, 1,
		104, 2,
		104, 3,
		104, 6,
		104, 5,
		104, 4,
		Halt,
	}
	screen, err := draw(program)
	if err != nil {
		t.Fatal(err)
	}

	pos := Point{1, 2}
	want := Tile(3)
	tile, ok := screen[pos]
	switch {
	case !ok:
		t.Errorf("expected a tile at %+v", pos)
	case tile != want:
		t.Errorf("got %v; want %v", tile, want)
	}

	pos = Point{6, 5}
	want = Tile(4)
	tile, ok = screen[pos]
	switch {
	case !ok:
		t.Errorf("expected a tile at %+v", pos)
	case tile != want:
		t.Errorf("got %v; want %v", tile, want)
	}
}
