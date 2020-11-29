package main

import "testing"

func TestFuelRequired(t *testing.T) {
	tests := []struct {
		name string
		Mass
		want Mass
	}{
		{
			name: "first example",
			Mass: 12,
			want: 2,
		},
		{
			name: "second example",
			Mass: 13,
			want: 2,
		},
		{
			name: "third example",
			Mass: 1969,
			want: 654,
		},
		{
			name: "fourth example",
			Mass: 100756,
			want: 33583,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if r := FuelRequired(tc.Mass); r != tc.want {
				t.Errorf("FuelRequired(%d) = %d; want %d", tc.Mass, r, tc.want)
			}
		})
	}
}

func TestTotalFuelRequired(t *testing.T) {
	tests := []struct {
		name string
		Mass
		want Mass
	}{
		{
			name: "first example",
			Mass: 14,
			want: 2,
		},
		{name: "second example",
			Mass: 1969,
			want: 966,
		},
		{
			name: "third example",
			Mass: 100756,
			want: 50346,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf, ff := TotalFuelRequired(tc.Mass)
			if r := mf + ff; r != tc.want {
				t.Errorf("TotalFuelRequired(%v) = %d; want %d", tc.Mass, r, tc.want)
			}
		})
	}
}
