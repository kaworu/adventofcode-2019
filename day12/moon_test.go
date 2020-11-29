package main

import "testing"

func TestNewMoon(t *testing.T) {
	tests := []struct {
		name string
		v    Vec3d
	}{
		{
			name: "first moon",
			v:    Vec3d{x: -1, y: 0, z: 2},
		},
		{
			name: "second moon",
			v:    Vec3d{x: 2, y: -10, z: -7},
		},
		{
			name: "third moon",
			v:    Vec3d{x: 4, y: -8, z: 8},
		},
		{
			name: "fourth moon",
			v:    Vec3d{x: 3, y: 5, z: -1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			moon := NewMoon(tc.v.x, tc.v.y, tc.v.z)
			if !equals(moon.pos, tc.v) {
				t.Errorf("moon.pos=%+v; want %+v", moon.pos, tc.v)
			}
		})
	}
}

func TestMoonAccelerate(t *testing.T) {
	m1 := NewMoon(5, 3, 0)
	m2 := NewMoon(3, 5, 0)
	Accelerate(m1, m2)

	got, want := m1.vel, Vec3d{x: -1, y: 1, z: 0}
	if !equals(got, want) {
		t.Errorf("got velocity %+v; want %+v", got, want)
	}

	got, want = m2.vel, Vec3d{x: 1, y: -1, z: 0}
	if !equals(got, want) {
		t.Errorf("got velocity %+v; want %+v", got, want)
	}
}

func TestMoonMove(t *testing.T) {
	tests := []struct {
		name           string
		pos, vel, want Vec3d
	}{
		{
			name: "first moon",
			pos:  Vec3d{x: -1, y: 0, z: 2},
			vel:  Vec3d{x: 3, y: -1, z: -1},
			want: Vec3d{x: 2, y: -1, z: 1},
		},
		{
			name: "second moon",
			pos:  Vec3d{x: 2, y: -10, z: -7},
			vel:  Vec3d{x: 1, y: 3, z: 3},
			want: Vec3d{x: 3, y: -7, z: -4},
		},
		{
			name: "third moon",
			pos:  Vec3d{x: 4, y: -8, z: 8},
			vel:  Vec3d{x: -3, y: 1, z: -3},
			want: Vec3d{x: 1, y: -7, z: 5},
		},
		{
			name: "fourth moon",
			pos:  Vec3d{x: 3, y: 5, z: -1},
			vel:  Vec3d{x: -1, y: -3, z: 1},
			want: Vec3d{x: 2, y: 2, z: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			moon := &Moon{pos: tc.pos, vel: tc.vel}
			moon.Move()
			if !equals(moon.pos, tc.want) {
				t.Errorf("got %+v; want %+v", moon.pos, tc.pos)
			}
		})
	}
}

func TestMoonEnergy(t *testing.T) {
	tests := []struct {
		name     string
		pos, vel Vec3d
		want     int
	}{
		{
			name: "first moon",
			pos:  Vec3d{x: 2, y: 1, z: -3},
			vel:  Vec3d{x: -3, y: -2, z: 1},
			want: 36,
		},
		{
			name: "second moon",
			pos:  Vec3d{x: 1, y: -8, z: 0},
			vel:  Vec3d{x: -1, y: 1, z: 3},
			want: 45,
		},
		{
			name: "third moon",
			pos:  Vec3d{x: 3, y: -6, z: 1},
			vel:  Vec3d{x: 3, y: 2, z: -3},
			want: 80,
		},
		{
			name: "fourth moon",
			pos:  Vec3d{x: 2, y: 0, z: 4},
			vel:  Vec3d{x: 1, y: -1, z: -1},
			want: 18,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			moon := &Moon{pos: tc.pos, vel: tc.vel}
			got := moon.Energy()
			if got != tc.want {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}
