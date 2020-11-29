package main

import "testing"

func TestNewPasswordVersion1(t *testing.T) {
	tests := []struct {
		name  string
		input uint64
		valid bool
	}{
		{
			name:  "first example",
			input: 111111,
			valid: true,
		},
		{
			name:  "second example",
			input: 223450,
			valid: false,
		},
		{
			name:  "third example",
			input: 123789,
			valid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &Digits6{}
			d.Set(tc.input)
			p := NewPassword(d)
			valid := (p.CompareTo(d) == 0 && p.HasRepeatingDigits())
			switch {
			case !valid && tc.valid:
				t.Errorf("expected %v to be a valid password", tc.input)
			case valid && !tc.valid:
				t.Errorf("expected %v to be an invalid password", tc.input)
			}
		})
	}
}

func TestNewPasswordVersion2(t *testing.T) {
	tests := []struct {
		name  string
		input uint64
		valid bool
	}{
		{
			name:  "first example",
			input: 112233,
			valid: true,
		},
		{
			name:  "second example",
			input: 123444,
			valid: false,
		},
		{
			name:  "third example",
			input: 111122,
			valid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &Digits6{}
			d.Set(tc.input)
			p := NewPassword(d)
			valid := (p.CompareTo(d) == 0 && p.HasDouble())
			switch {
			case !valid && tc.valid:
				t.Errorf("expected %v to be a valid password", tc.input)
			case valid && !tc.valid:
				t.Errorf("expected %v to be an invalid password", tc.input)
			}
		})
	}
}
