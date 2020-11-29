package main

import (
	"testing"
)

func TestTurning(t *testing.T) {
	tests := []struct {
		from Heading
		dir  Direction
		want Heading
	}{
		{from: North, dir: Right, want: East},
		{from: North, dir: Left, want: West},
		{from: East, dir: Right, want: South},
		{from: East, dir: Left, want: North},
		{from: South, dir: Right, want: West},
		{from: South, dir: Left, want: East},
		{from: West, dir: Right, want: North},
		{from: West, dir: Left, want: South},
	}
	for _, tc := range tests {
		if h := tc.from.Turning(tc.dir); h != tc.want {
			t.Errorf("%v.Turning(%v) = %v; want %v", tc.from, tc.dir, h, tc.want)
		}
	}
}

func TestMoving(t *testing.T) {
	origin := Point{x: 0, y: 0}
	tests := []struct {
		from    Point
		heading Heading
		want    Point
	}{
		{from: origin, heading: North, want: Point{x: 0, y: -1}},
		{from: origin, heading: East, want: Point{x: 1, y: 0}},
		{from: origin, heading: South, want: Point{x: 0, y: 1}},
		{from: origin, heading: West, want: Point{x: -1, y: 0}},
		{from: origin, heading: West + 1, want: origin}, // invalid heading.
	}
	for _, tc := range tests {
		if p := tc.from.Moving(tc.heading); p != tc.want {
			t.Errorf("%v.Moving(%v) = %v; want %v", tc.from, tc.heading, p, tc.want)
		}
	}
}

func TestPaint(t *testing.T) {
	ship := NewSpacecraft()
	robot := NewRobot()
	// This program simulate the test output sequence and verify that each
	// input matches the expectation. Tests are done directly in Intcode.
	program := []Intcode{
		// start of the program: goto main
		1105, 1, 32,
		// .data
		0, // [3]: Read register
		0, // [4]: Equals register
		// steps
		0, 1, 0, // expected 0, and output 1 (paint White) then 0 (turn Left)
		0, 0, 0, // expected 0, and output 0 (paint Black) then 0 (turn Left)
		0, 1, 0, // expected 0, and output 1 (paint White) then 0 (turn Left)
		0, 1, 0, // expected 0, and output 1 (paint White) then 0 (turn Left)
		1, 0, 1, // expected 1, and output 0 (paint Black) then 1 (turn Right)
		0, 1, 0, // expected 0, and output 1 (paint White) then 0 (turn Left)
		0, 1, 0, // expected 0, and output 1 (paint White) then 0 (turn Left)
		// .text
		// [26] halt:
		Halt, // signal the end of the test vector.
		// [27] fail: setup rbo[0] to -1 so that we can debug which step went
		// wrong, then trigger an illegal instruction.
		21101, -1, 0, 0, // rbo[0] = -1
		0, // invalid instruction.
		// [32] main: rbo will always point to the start of a step line
		109, 5, // setup rbo to the first step
		// [34] loop: until rbo hit the halt instruction.
		1208, 0, 99, 4, // [4] = (rbo[0] == 99 ? 1 : 0)
		11005, 4, 26, // goto halt if [4] != 0
		3, 3, // mem[3] = Read
		// check that our read match the expected value (i.e. rbo[0])
		208, 0, 3, 4, // [4] = (rbo[0] == [3] ? 1 : 0)
		1006, 4, 27, // goto fail if [4] == 0
		204, 1, // Write rbo[1]
		204, 2, // Write rbo[2]
		109, 3, // rbo += 3
		1105, 1, 34, // goto loop
	}
	err := robot.Paint(ship, program)
	if err != nil {
		// Here we could easily check what went wrong by inspecting the robot's
		// computer memory (find the step starting with -1).
		t.Errorf("Robot.Paint() error: %s", err)
	}
}
