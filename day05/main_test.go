package main

import "testing"

func TestExecute(t *testing.T) {
	tests := []struct {
		Name       string
		Program    Memory
		Input      Input
		Expected   Output
		FinalState Memory
	}{
		{
			Name:       "identity program",
			Program:    Memory{3, 0, 4, 0, 99},
			Input:      Input{42},
			Expected:   Output{42},
			FinalState: Memory{42, 0, 4, 0, 99},
		},
		{
			Name:       "multiply then halt",
			Program:    Memory{1002, 4, 3, 4, 33},
			Input:      nil,
			Expected:   nil,
			FinalState: Memory{1002, 4, 3, 4, 99},
		},
		{
			Name:       "add immediate with negative value",
			Program:    Memory{1101, 100, -1, 4, 0},
			Input:      nil,
			Expected:   nil,
			FinalState: Memory{1101, 100, -1, 4, 99},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			mem := tc.Program.Copy()
			output, err := mem.Execute(tc.Input)
			if err != nil {
				t.Errorf("Unexpected execute error: %s", err)
			}
			if !IntcodeEqual(output, tc.Expected) {
				t.Errorf("expected %v as output, got %v", tc.Expected, output)
			}
			if !IntcodeEqual(mem, tc.FinalState) {
				t.Errorf("expected %v as final state, got %v", tc.FinalState, mem)
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
