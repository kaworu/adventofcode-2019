package main

import "testing"

func TestNewPasswordVersion1(t *testing.T) {
	tests := []struct {
		Name  string
		Input uint64
		Valid bool
	}{
		{
			Name:  "first example",
			Input: 111111,
			Valid: true,
		},
		{
			Name:  "second example",
			Input: 223450,
			Valid: false,
		},
		{
			Name:  "third example",
			Input: 123789,
			Valid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			d := &Digits6{}
			d.Set(tc.Input)
			p := NewPassword(d)
			valid := (p.CompareTo(d) == 0 && p.HasRepeatingDigits())
			switch {
			case !valid && tc.Valid:
				t.Errorf("expected %v to be a valid password", tc.Input)
			case valid && !tc.Valid:
				t.Errorf("expected %v to be an invalid password", tc.Input)
			}
		})
	}
}

func TestNewPasswordVersion2(t *testing.T) {
	tests := []struct {
		Name  string
		Input uint64
		Valid bool
	}{
		{
			Name:  "first example",
			Input: 112233,
			Valid: true,
		},
		{
			Name:  "second example",
			Input: 123444,
			Valid: false,
		},
		{
			Name:  "third example",
			Input: 111122,
			Valid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			d := &Digits6{}
			d.Set(tc.Input)
			p := NewPassword(d)
			valid := (p.CompareTo(d) == 0 && p.HasDouble())
			switch {
			case !valid && tc.Valid:
				t.Errorf("expected %v to be a valid password", tc.Input)
			case valid && !tc.Valid:
				t.Errorf("expected %v to be an invalid password", tc.Input)
			}
		})
	}
}
