package main

import "testing"

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		program []Intcode
		input   Input
		want    Output
	}{
		{
			name:    "Quine",
			program: []Intcode{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
			input:   Input{},
			want:    []Intcode{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
		},
		{ // MULT(34915192, 34915192) -> 7; WRITE [7]; HALT
			name:    "big",
			program: []Intcode{1102, 34915192, 34915192, 7, 4, 7, 99, 0},
			input:   Input{},
			want:    []Intcode{34915192 * 34915192},
		},
		{ // WRITE 1125899906842624; HALT
			name:    "1125899906842624",
			program: []Intcode{104, 1125899906842624, 99},
			input:   Input{},
			want:    []Intcode{1125899906842624},
		},
	}

	var c Computer
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := c.Execute(tc.program, tc.input)
			switch {
			case err != nil:
				t.Errorf("Execute(%v, %v) error: %s", tc.program, tc.input, err)
			case !IntcodeEqual(output, tc.want):
				t.Errorf("Execute(%v, %v) = %v; want %v", tc.program, tc.input, output, tc.want)
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
