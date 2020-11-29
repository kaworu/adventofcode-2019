package main

import "testing"

func TestVec3dAdd(t *testing.T) {
	// NOTE: test cases taken from the first example.
	tests := []struct {
		name          string
		lhs, rhs, sum Vec3d
	}{
		{
			name: "zeroes",
			lhs:  Vec3d{x: 0, y: 0, z: 0},
			rhs:  Vec3d{x: 0, y: 0, z: 0},
			sum:  Vec3d{x: 0, y: 0, z: 0},
		},
		{
			name: "first planet first step",
			lhs:  Vec3d{x: -1, y: 0, z: 2},  // pos at step 0
			rhs:  Vec3d{x: 3, y: -1, z: -1}, // vel at step 1
			sum:  Vec3d{x: 2, y: -1, z: 1},  // pos at step 1
		},
		{
			name: "second planet second step",
			lhs:  Vec3d{x: 3, y: -7, z: -4}, // pos at step 1
			rhs:  Vec3d{x: -2, y: 5, z: 6},  // vel at step 2
			sum:  Vec3d{x: 1, y: -2, z: 2},  // pos at step 2
		},
		{
			name: "third planet third step",
			lhs:  Vec3d{x: 1, y: -4, z: -1}, // pos at step 2
			rhs:  Vec3d{x: 1, y: 5, z: -4},  // vel at step 3
			sum:  Vec3d{x: 2, y: 1, z: -5},  // pos at step 3
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.lhs.Add(tc.rhs)
			if !equals(got, tc.sum) {
				t.Errorf("%+v.Add(%+v) = %+v; want %+v", tc.lhs, tc.rhs, got, tc.sum)
			}
			// Test commutativity, i.e. lhs + rhs = rhs + lhs
			got = tc.rhs.Add(tc.lhs)
			if !equals(got, tc.sum) {
				t.Errorf("%+v.Add(%+v) = %+v; want %+v", tc.rhs, tc.lhs, got, tc.sum)
			}
		})
	}
}

func TestVec3dSub(t *testing.T) {
	// NOTE: test cases taken from the first example. Same test vectors as
	// TestVec3dAdd, hence lhs + rhs = sum hold.
	tests := []struct {
		name          string
		lhs, rhs, sum Vec3d
	}{
		{
			name: "zeroes",
			lhs:  Vec3d{x: 0, y: 0, z: 0},
			rhs:  Vec3d{x: 0, y: 0, z: 0},
			sum:  Vec3d{x: 0, y: 0, z: 0},
		},
		{
			name: "first planet first step",
			lhs:  Vec3d{x: -1, y: 0, z: 2},  // pos at step 0
			rhs:  Vec3d{x: 3, y: -1, z: -1}, // vel at step 1
			sum:  Vec3d{x: 2, y: -1, z: 1},  // pos at step 1
		},
		{
			name: "second planet second step",
			lhs:  Vec3d{x: 3, y: -7, z: -4}, // pos at step 1
			rhs:  Vec3d{x: -2, y: 5, z: 6},  // vel at step 2
			sum:  Vec3d{x: 1, y: -2, z: 2},  // pos at step 2
		},
		{
			name: "third planet third step",
			lhs:  Vec3d{x: 1, y: -4, z: -1}, // pos at step 2
			rhs:  Vec3d{x: 1, y: 5, z: -4},  // vel at step 3
			sum:  Vec3d{x: 2, y: 1, z: -5},  // pos at step 3
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.sum.Sub(tc.lhs)
			if !equals(got, tc.rhs) {
				t.Errorf("%+v.Sub(%+v) = %+v; want %+v", tc.sum, tc.lhs, got, tc.rhs)
			}
			got = tc.sum.Sub(tc.rhs)
			if !equals(got, tc.lhs) {
				t.Errorf("%+v.Sub(%+v) = %+v; want %+v", tc.sum, tc.rhs, got, tc.lhs)
			}
		})
	}
}

func TestVec3dDiff(t *testing.T) {
	// NOTE: test cases inspired by the first example.
	tests := []struct {
		name           string
		lhs, rhs, diff Vec3d
	}{
		{
			name: "zeroes",
			lhs:  Vec3d{x: 0, y: 0, z: 0},
			rhs:  Vec3d{x: 0, y: 0, z: 0},
			diff: Vec3d{x: 0, y: 0, z: 0},
		},
		{
			name: "first planet first step",
			lhs:  Vec3d{x: -1, y: 0, z: 2},  // pos at step 0
			rhs:  Vec3d{x: 3, y: -1, z: -1}, // vel at step 1
			diff: Vec3d{x: -1, y: 1, z: 1},
		},
		{
			name: "second planet second step",
			lhs:  Vec3d{x: 3, y: -7, z: -4}, // pos at step 1
			rhs:  Vec3d{x: -2, y: 5, z: 6},  // vel at step 2
			diff: Vec3d{x: 1, y: -1, z: -1},
		},
		{
			name: "third planet third step",
			lhs:  Vec3d{x: 1, y: -4, z: -1}, // pos at step 2
			rhs:  Vec3d{x: 1, y: 5, z: -4},  // vel at step 3
			diff: Vec3d{x: 0, y: -1, z: 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.lhs.Diff(tc.rhs)
			if !equals(got, tc.diff) {
				t.Errorf("%+v.Diff(%+v) = %+v; want %+v", tc.lhs, tc.rhs, got, tc.diff)
			}
		})
	}
}

func TestVec3dSumAbs(t *testing.T) {
	// NOTE: test cases taken from the first example.
	tests := []struct {
		name string
		vec  Vec3d
		want int
	}{
		{
			name: "zeroes",
			vec:  Vec3d{x: 0, y: 0, z: 0},
			want: 0,
		},
		{
			name: "first planet pos tenth step",
			vec:  Vec3d{x: 2, y: 1, z: -3}, // pos at step 10
			want: 6,
		},
		{
			name: "second planet pos tenth step",
			vec:  Vec3d{x: 1, y: -8, z: -0}, // pos at step 10
			want: 9,
		},
		{
			name: "third planet pos tenth step",
			vec:  Vec3d{x: 3, y: -6, z: 1}, // pos at step 10
			want: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.vec.SumAbs()
			if got != tc.want {
				t.Errorf("%+v.SumAbs() = %v; want %d", tc.vec, got, tc.want)
			}
		})
	}
}

// equals returns true when a and b are the same Vec3d, false otherwise.
func equals(a, b Vec3d) bool {
	return a.x == b.x && a.y == b.y && a.z == b.z
}
