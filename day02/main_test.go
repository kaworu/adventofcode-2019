package main

import "testing"

func TestExecute(t *testing.T) {
	tests := []struct {
		Name     string
		Program  Intcode
		Expected Intcode
	}{
		{
			"detailed example",
			Intcode{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			Intcode{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50},
		},
		{
			"first example",
			Intcode{1, 0, 0, 0, 99},
			Intcode{2, 0, 0, 0, 99},
		},
		{
			"second example",
			Intcode{2, 3, 0, 3, 99},
			Intcode{2, 3, 0, 6, 99},
		},
		{
			"third example",
			Intcode{2, 4, 4, 5, 99, 0},
			Intcode{2, 4, 4, 5, 99, 9801},
		},
		{
			"fourth example",
			Intcode{1, 1, 1, 4, 99, 5, 6, 0, 99},
			Intcode{30, 1, 1, 4, 2, 5, 6, 0, 99},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			if err := tc.Program.Execute(); err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
			if !Equal(tc.Program, tc.Expected) {
				t.Errorf("got %v, expected %v", tc.Program, tc.Expected)
			}
		})
	}
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
