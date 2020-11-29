package main

import "testing"

func TestInputOutput(t *testing.T) {
	tests := []struct {
		name   string
		prog   Memory
		input  Input
		output Output
	}{
		{
			name:   "identity program",
			prog:   Memory{3, 0, 4, 0, 99},
			input:  Input{42},
			output: Output{42},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mem := tc.prog.Copy()
			output, err := mem.Execute(tc.input)
			switch {
			case err != nil:
				t.Errorf("%v.Execute(%v) error", tc.prog, tc.input)
			case !IntcodeEqual(output, tc.output):
				t.Errorf("%v.Execute(%v) = %v; want %v", tc.prog, tc.input, output, tc.output)
			}
		})
	}
}

func TestModes(t *testing.T) {
	tests := []struct {
		name string
		prog Memory
		want Memory
	}{
		{
			name: "multiply then halt",
			prog: Memory{1002, 4, 3, 4, 33},
			want: Memory{1002, 4, 3, 4, 99},
		},
		{
			name: "add immediate with negative value",
			prog: Memory{1101, 100, -1, 4, 0},
			want: Memory{1101, 100, -1, 4, 99},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mem := tc.prog.Copy()
			_, err := mem.Execute( /* Input */ nil)
			switch {
			case err != nil:
				t.Errorf("%v.Execute() error", tc.prog)
			case !IntcodeEqual(mem, tc.want):
				t.Errorf("%v.Execute() final state = %v; want %v", tc.prog, mem, tc.want)
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
		name string
		prog Memory
		exec func(in Input) Output
	}{
		{
			name: "position mode input==8",
			prog: Memory{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			exec: eq8,
		},
		{
			name: "position mode input<8",
			prog: Memory{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8},
			exec: lt8,
		},
		{
			name: "immediate mode input==8",
			prog: Memory{3, 3, 1108, -1, 8, 3, 4, 3, 99},
			exec: eq8,
		},
		{
			name: "immediate mode input<8",
			prog: Memory{3, 3, 1107, -1, 8, 3, 4, 3, 99},
			exec: lt8,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for i := -10; i <= 10; i++ {
				mem := tc.prog.Copy()
				in := Input{Intcode(i)}
				want := tc.exec(in)
				output, err := mem.Execute(in)
				switch {
				case err != nil:
					t.Errorf("%v.Execute(%v) error: %s", tc.prog, in, err)
					return
				case !IntcodeEqual(output, want):
					t.Errorf("%v.Execute(%v) = %v; want %v", tc.prog, in, output, want)
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
