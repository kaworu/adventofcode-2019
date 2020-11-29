package main

import "testing"

func TestExecute(t *testing.T) {
	tests := []struct {
		name string
		prog Memory
		want Memory
	}{
		{
			name: "detailed example",
			prog: Memory{
				1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50,
			},
			want: Memory{
				3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50,
			},
		},
		{
			name: "first example",
			prog: Memory{
				1, 0, 0, 0, 99,
			},
			want: Memory{
				2, 0, 0, 0, 99,
			},
		},
		{
			name: "second example",
			prog: Memory{
				2, 3, 0, 3, 99,
			},
			want: Memory{
				2, 3, 0, 6, 99,
			},
		},
		{
			name: "third example",
			prog: Memory{
				2, 4, 4, 5, 99, 0,
			},
			want: Memory{
				2, 4, 4, 5, 99, 9801,
			},
		},
		{
			name: "fourth example",
			prog: Memory{
				1, 1, 1, 4, 99, 5, 6, 0, 99,
			},
			want: Memory{
				30, 1, 1, 4, 2, 5, 6, 0, 99,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mem := tc.prog.Copy()
			if err := mem.Execute(); err != nil {
				t.Fatalf("%v.Execute() error: %s", tc.prog, err)
			}
			if !Equal(mem, tc.want) {
				t.Fatalf("%v.Execute() = %v; want %v", tc.prog, mem, tc.want)
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
