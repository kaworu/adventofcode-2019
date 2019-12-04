package main

import "testing"

func TestFuelRequired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name string
		Mod  Module
		Exp  Fuel
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
			if r := tc.Mod.FuelRequired(); r != tc.Exp {
				t.Errorf("%v.FuelRequired() = %d, expected %d", tc.Mod, r, tc.Exp)
			}
		})
	}
}
