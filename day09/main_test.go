package main

import "testing"

func TestExecute(t *testing.T) {
	tests := []struct {
		Name     string
		Program  []Intcode
		Input    Input
		Expected Output
	}{
		{
			Name:     "Quine",
			Program:  []Intcode{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
			Input:    Input{},
			Expected: []Intcode{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
		},
		{ // MULT(34915192, 34915192) -> 7; WRITE [7]; HALT
			Name:     "big",
			Program:  []Intcode{1102, 34915192, 34915192, 7, 4, 7, 99, 0},
			Input:    Input{},
			Expected: []Intcode{34915192 * 34915192},
		},
		{ // WRITE 1125899906842624; HALT
			Name:     "1125899906842624",
			Program:  []Intcode{104, 1125899906842624, 99},
			Input:    Input{},
			Expected: []Intcode{1125899906842624},
		},
	}

	var c Computer
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			output, err := c.Execute(tc.Program, tc.Input)
			switch {
			case err != nil:
				t.Errorf("Execute(%v, %v) error: %s", tc.Program, tc.Input, err)
			case !IntcodeEqual(output, tc.Expected):
				t.Errorf("Execute(%v, %v) = %v; expected %v", tc.Program, tc.Input, output, tc.Expected)
			}
		})
	}
}

// IntcodeEqual compare two Intcode slices and returns true if they are the
// same, false otherwise.
func IntcodeEqual(xs, ys []Intcode) bool {
	if len(xs) != len(ys) {
		return false
	}
	for i := range xs {
		if xs[i] != ys[i] {
			return false
		}
	}
	return true
}
