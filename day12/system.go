package main

// System represent the Jupiter orbit of its four largest moons.
type System struct {
	Moons []*Moon
}

// Simulate the motion of the moons a given number of steps.
func (s *System) Simulate(steps int) {
	for i := 0; i < steps; i++ {
		// gravity
		for i, m := range s.Moons {
			for j := i + 1; j < len(s.Moons); j++ {
				Accelerate(m, s.Moons[j])
			}
		}
		// velocity
		for _, m := range s.Moons {
			m.Move()
		}
	}
}

// Total energy returns the total energy in the system.
func (s *System) TotalEnergy() int {
	energy := 0
	for _, m := range s.Moons {
		energy += m.Energy()
	}
	return energy
}
