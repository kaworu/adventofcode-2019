package main

// A 3-element vector represented by int x, y, z coordinates.
type Vec3d struct {
	x, y, z int
}

// Add return a Vec3d that is the vector addition result of v + o.
func (v Vec3d) Add(o Vec3d) Vec3d {
	return Vec3d{
		x: v.x + o.x,
		y: v.y + o.y,
		z: v.z + o.z,
	}
}

// Sub return a Vec3d that is the vector subtraction result of v - o.
func (v Vec3d) Sub(o Vec3d) Vec3d {
	return Vec3d{
		x: v.x - o.x,
		y: v.y - o.y,
		z: v.z - o.z,
	}
}

// Diff return a Vec3d that is the comparison result of v and o. Each result
// coordinate is set to -1, 0, or 1 when v's coordinate is found to be
// respectively lesser than, equal to, or greater than o's coordinate.
func (v Vec3d) Diff(o Vec3d) Vec3d {
	return Vec3d{
		x: idiff(v.x, o.x),
		y: idiff(v.y, o.y),
		z: idiff(v.z, o.z),
	}
}

// SumAbs returns the sum of the absolute values of v's coordinates.
func (v Vec3d) SumAbs() int {
	return iabs(v.x) + iabs(v.y) + iabs(v.z)
}

// iabs returns the absolute value of x.
// NOTE: broken for the minimum integer value.
func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// idiff returns -1, 0, or 1 when i is found to be respectively lesser than,
// equals to, or greater than j.
func idiff(i, j int) int {
	switch {
	case i > j:
		return +1
	case i < j:
		return -1
	default:
		return 0
	}
}
