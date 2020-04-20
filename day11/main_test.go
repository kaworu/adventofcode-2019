package main

import (
	"testing"
)

func TestTurning(t *testing.T) {
	tests := []struct {
		From, To Heading
		Direction
	}{
		{From: North, Direction: Right, To: East},
		{From: North, Direction: Left, To: West},
		{From: East, Direction: Right, To: South},
		{From: East, Direction: Left, To: North},
		{From: South, Direction: Right, To: West},
		{From: South, Direction: Left, To: East},
		{From: West, Direction: Right, To: North},
		{From: West, Direction: Left, To: South},
	}
	for _, tc := range tests {
		if h := tc.From.Turning(tc.Direction); h != tc.To {
			t.Errorf("%v.Turning(%v) = %v; expected %v", tc.From, tc.Direction, h, tc.To)
		}
	}
}

func TestMoving(t *testing.T) {
	origin := Point{x: 0, y: 0}
	tests := []struct {
		From, To Point
		Heading
	}{
		{From: origin, Heading: North, To: Point{x: 0, y: -1}},
		{From: origin, Heading: East, To: Point{x: 1, y: 0}},
		{From: origin, Heading: South, To: Point{x: 0, y: 1}},
		{From: origin, Heading: West, To: Point{x: -1, y: 0}},
		{From: origin, Heading: West + 1, To: origin}, // invalid heading.
	}
	for _, tc := range tests {
		if p := tc.From.Moving(tc.Heading); p != tc.To {
			t.Errorf("%v.Moving(%v) = %v; expected %v", tc.From, tc.Heading, p, tc.To)
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
