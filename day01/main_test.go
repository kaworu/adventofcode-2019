package main

import "testing"

func TestFuelRequired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name string
		Mass
		Expected Mass
	}{
		{"first example", 12, 2},
		{"second example", 14, 2},
		{"third example", 1969, 654},
		{"fourth example", 100756, 33583},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			if r := FuelRequired(tc.Mass); r != tc.Expected {
				t.Errorf("FuelRequired(%d) = %d; expected %d", tc.Mass, r, tc.Expected)
			}
		})
	}
}

func TestTotalFuelRequired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name string
		Mass
		Expected Mass
	}{
		{"first example", 14, 2},
		{"second example", 1969, 966},
		{"third example", 100756, 50346},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			mf, ff := TotalFuelRequired(tc.Mass)
			if r := mf + ff; r != tc.Expected {
				t.Errorf("TotalFuelRequired(%v) = %d; expected %d", tc.Mass, r, tc.Expected)
			}
		})
	}
}
