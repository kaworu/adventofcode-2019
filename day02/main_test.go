package main

import "testing"

func TestExecute(t *testing.T) {
	tests := []struct {
		Name     string
		Program  Memory
		Expected Memory
	}{
		{
			Name: "detailed example",
			Program: Memory{
				1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50,
			},
			Expected: Memory{
				3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50,
			},
		},
		{
			Name: "first example",
			Program: Memory{
				1, 0, 0, 0, 99,
			},
			Expected: Memory{
				2, 0, 0, 0, 99,
			},
		},
		{
			Name: "second example",
			Program: Memory{
				2, 3, 0, 3, 99,
			},
			Expected: Memory{
				2, 3, 0, 6, 99,
			},
		},
		{
			Name: "third example",
			Program: Memory{
				2, 4, 4, 5, 99, 0,
			},
			Expected: Memory{
				2, 4, 4, 5, 99, 9801,
			},
		},
		{
			Name: "fourth example",
			Program: Memory{
				1, 1, 1, 4, 99, 5, 6, 0, 99,
			},
			Expected: Memory{
				30, 1, 1, 4, 2, 5, 6, 0, 99,
			},
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
func Equal(a, b Memory) bool {
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
