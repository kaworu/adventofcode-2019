package main

// Moon represent one of the Galilean satellite orbiting around Jupiter.
type Moon struct {
	pos, vel Vec3d
}

// NewMoon returns a Moon at the given position.
func NewMoon(x, y, z int) *Moon {
	return &Moon{
		pos: Vec3d{x: x, y: y, z: z},
		vel: Vec3d{x: 0, y: 0, z: 0},
	}
}

// Accelerate apply gravity to both a and b.
func Accelerate(a, b *Moon) {
	gravity := a.pos.Diff(b.pos)
	b.vel = b.vel.Add(gravity)
	a.vel = a.vel.Sub(gravity)
}

// Move apply m's velocity to its position.
func (m *Moon) Move() {
	m.pos = m.pos.Add(m.vel)
}

// Energy returns m's potential energy multiplied by its kinetic energy.
func (m *Moon) Energy() int {
	return m.pos.SumAbs() * m.vel.SumAbs()
}
