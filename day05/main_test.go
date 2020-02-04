package main

import "testing"

func TestInputOutput(t *testing.T) {
	tests := []struct {
		Name     string
		Program  Memory
		Input    Input
		Expected Output
	}{
		{
			Name:     "identity program",
			Program:  Memory{3, 0, 4, 0, 99},
			Input:    Input{42},
			Expected: Output{42},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			mem := tc.Program.Copy()
			output, err := mem.Execute(tc.Input)
			switch {
			case err != nil:
				t.Errorf("%v.Execute(%v) error", tc.Program, tc.Input)
			case !IntcodeEqual(output, tc.Expected):
				t.Errorf("%v.Execute(%v) = %v; expected %v", tc.Program, tc.Input, output, tc.Expected)
			}
		})
	}
}

func TestModes(t *testing.T) {
	tests := []struct {
		Name       string
		Program    Memory
		FinalState Memory
	}{
		{
			Name:       "multiply then halt",
			Program:    Memory{1002, 4, 3, 4, 33},
			FinalState: Memory{1002, 4, 3, 4, 99},
		},
		{
			Name:       "add immediate with negative value",
			Program:    Memory{1101, 100, -1, 4, 0},
			FinalState: Memory{1101, 100, -1, 4, 99},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			mem := tc.Program.Copy()
			_, err := mem.Execute( /* Input */ nil)
			switch {
			case err != nil:
				t.Errorf("%v.Execute() error", tc.Program)
			case !IntcodeEqual(mem, tc.FinalState):
				t.Errorf("%v.Execute() final state = %v; expected %v", tc.Program, mem, tc.FinalState)
			}
		})
	}
}

func TestComparisons(t *testing.T) {
	eq8 := func(in Input) Output {
		if in[0] == 8 {
			return Output{1}
		}
		return Output{0}
	}
	lt8 := func(in Input) Output {
		if in[0] < 8 {
			return Output{1}
		}
		return Output{0}
	}
	tests := []struct {
		Name    string
		Program Memory
		Execute func(in Input) Output
	}{
		{
			Name:    "position mode input==8",
			Program: Memory{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			Execute: eq8,
		},
		{
			Name:    "position mode input<8",
			Program: Memory{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8},
			Execute: lt8,
		},
		{
			Name:    "immediate mode input==8",
			Program: Memory{3, 3, 1108, -1, 8, 3, 4, 3, 99},
			Execute: eq8,
		},
		{
			Name:    "immediate mode input<8",
			Program: Memory{3, 3, 1107, -1, 8, 3, 4, 3, 99},
			Execute: lt8,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			for i := -10; i <= 10; i++ {
				mem := tc.Program.Copy()
				in := Input{Intcode(i)}
				expected := tc.Execute(in)
				output, err := mem.Execute(in)
				switch {
				case err != nil:
					t.Errorf("%v.Execute(%v) error: %s", tc.Program, in, err)
					return
				case !IntcodeEqual(output, expected):
					t.Errorf("%v.Execute(%v) = %v; expected %v", tc.Program, in, output, expected)
					return
				}
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
