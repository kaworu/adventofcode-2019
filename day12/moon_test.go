package main

import "testing"

func TestNewMoon(t *testing.T) {
	tests := []struct {
		Name string
		pos  Vec3d
	}{
		{
			Name: "first moon",
			pos:  Vec3d{x: -1, y: 0, z: 2},
		},
		{
			Name: "second moon",
			pos:  Vec3d{x: 2, y: -10, z: -7},
		},
		{
			Name: "third moon",
			pos:  Vec3d{x: 4, y: -8, z: 8},
		},
		{
			Name: "fourth moon",
			pos:  Vec3d{x: 3, y: 5, z: -1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			moon := NewMoon(tc.pos.x, tc.pos.y, tc.pos.z)
			if !equals(moon.pos, tc.pos) {
				t.Errorf("moon.pos=%+v; expected %+v", moon.pos, tc.pos)
			}
		})
	}
}

func TestMoonAccelerate(t *testing.T) {
	m1 := NewMoon(5, 3, 0)
	m2 := NewMoon(3, 5, 0)
	Accelerate(m1, m2)

	vel1 := Vec3d{x: -1, y: 1, z: 0}
	if !equals(m1.vel, vel1) {
		t.Errorf("got velocity %+v; expected %+v", m1.vel, vel1)
	}

	vel2 := Vec3d{x: 1, y: -1, z: 0}
	if !equals(m2.vel, vel2) {
		t.Errorf("got velocity %+v; expected %+v", m2.vel, vel2)
	}
}

func TestMoonMove(t *testing.T) {
	tests := []struct {
		Name          string
		pos, vel, res Vec3d
	}{
		{
			Name: "first moon",
			pos:  Vec3d{x: -1, y: 0, z: 2},
			vel:  Vec3d{x: 3, y: -1, z: -1},
			res:  Vec3d{x: 2, y: -1, z: 1},
		},
		{
			Name: "second moon",
			pos:  Vec3d{x: 2, y: -10, z: -7},
			vel:  Vec3d{x: 1, y: 3, z: 3},
			res:  Vec3d{x: 3, y: -7, z: -4},
		},
		{
			Name: "third moon",
			pos:  Vec3d{x: 4, y: -8, z: 8},
			vel:  Vec3d{x: -3, y: 1, z: -3},
			res:  Vec3d{x: 1, y: -7, z: 5},
		},
		{
			Name: "fourth moon",
			pos:  Vec3d{x: 3, y: 5, z: -1},
			vel:  Vec3d{x: -1, y: -3, z: 1},
			res:  Vec3d{x: 2, y: 2, z: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			moon := &Moon{pos: tc.pos, vel: tc.vel}
			moon.Move()
			if !equals(moon.pos, tc.res) {
				t.Errorf("got position %+v; expected %+v", moon.pos, tc.pos)
			}
		})
	}
}

func TestMoonEnergy(t *testing.T) {
	tests := []struct {
		Name     string
		pos, vel Vec3d
		expected int
	}{
		{
			Name:     "first moon",
			pos:      Vec3d{x: 2, y: 1, z: -3},
			vel:      Vec3d{x: -3, y: -2, z: 1},
			expected: 36,
		},
		{
			Name:     "second moon",
			pos:      Vec3d{x: 1, y: -8, z: 0},
			vel:      Vec3d{x: -1, y: 1, z: 3},
			expected: 45,
		},
		{
			Name:     "third moon",
			pos:      Vec3d{x: 3, y: -6, z: 1},
			vel:      Vec3d{x: 3, y: 2, z: -3},
			expected: 80,
		},
		{
			Name:     "fourth moon",
			pos:      Vec3d{x: 2, y: 0, z: 4},
			vel:      Vec3d{x: 1, y: -1, z: -1},
			expected: 18,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			moon := &Moon{pos: tc.pos, vel: tc.vel}
			got := moon.Energy()
			if got != tc.expected {
				t.Errorf("got %v; expected %v", got, tc.expected)
			}
		})
	}
}
