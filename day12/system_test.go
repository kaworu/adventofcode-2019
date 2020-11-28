package main

import "testing"

func TestSystemSimulate(t *testing.T) {
	// the system from the first example at step 0.
	s := &System{
		Moons: []*Moon{
			NewMoon(-1, 0, 2),
			NewMoon(2, -10, -7),
			NewMoon(4, -8, 8),
			NewMoon(3, 5, -1),
		},
	}

	s.Simulate(1)

	step1 := []struct {
		pos, vel Vec3d
	}{
		{
			pos: Vec3d{x: 2, y: -1, z: 1},
			vel: Vec3d{x: 3, y: -1, z: -1},
		},
		{
			pos: Vec3d{x: 3, y: -7, z: -4},
			vel: Vec3d{x: 1, y: 3, z: 3},
		},
		{
			pos: Vec3d{x: 1, y: -7, z: 5},
			vel: Vec3d{x: -3, y: 1, z: -3},
		}, {
			pos: Vec3d{x: 2, y: 2, z: 0},
			vel: Vec3d{x: -1, y: -3, z: 1},
		},
	}

	for i, m := range s.Moons {
		if !equals(m.pos, step1[i].pos) {
			t.Errorf("moon %d pos after step 1 = %+v; expected %+v", i, m.pos, step1[i].pos)
		}
		if !equals(m.vel, step1[i].vel) {
			t.Errorf("moon %d vel after step 1 = %+v; expected %+v", i, m.vel, step1[i].vel)
		}
	}
}

func TestSystemTotalEnergy(t *testing.T) {
	// the system from the first example at step 10.
	s := &System{
		Moons: []*Moon{
			&Moon{
				pos: Vec3d{x: 2, y: 1, z: -3},
				vel: Vec3d{x: -3, y: -2, z: 1},
			},
			&Moon{
				pos: Vec3d{x: 1, y: -8, z: 0},
				vel: Vec3d{x: -1, y: 1, z: 3},
			},
			&Moon{
				pos: Vec3d{x: 3, y: -6, z: 1},
				vel: Vec3d{x: 3, y: 2, z: -3},
			},
			&Moon{
				pos: Vec3d{x: 2, y: 0, z: 4},
				vel: Vec3d{x: 1, y: -1, z: -1},
			},
		},
	}
	got := s.TotalEnergy()
	expected := 179
	if got != expected {
		t.Errorf("got %v; expected %v", got, expected)
	}
}
