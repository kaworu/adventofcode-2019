package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	encoded := "0222112222120000"
	expected := []*Layer{
		&Layer{2, 2, []Pixel{0, 2, 2, 2}},
		&Layer{2, 2, []Pixel{1, 1, 2, 2}},
		&Layer{2, 2, []Pixel{2, 2, 1, 2}},
		&Layer{2, 2, []Pixel{0, 0, 0, 0}},
	}
	decoded, err := Parse(2, 2, strings.NewReader(encoded))
	if err != nil {
		t.Errorf("parsing failed: %s", err)
	} else if len(expected) != len(decoded) {
		t.Errorf("expected %d layers, got %d instead", len(expected), len(decoded))
	} else {
		for i := range expected {
			if !LayerEquals(expected[i], decoded[i]) {
				t.Errorf("expected %v for layer %d but got %v", expected[i], i, decoded[i])
			}
		}
	}
}

func TestFlatten(t *testing.T) {
	layers := []*Layer{
		&Layer{2, 2, []Pixel{0, 2, 2, 2}},
		&Layer{2, 2, []Pixel{1, 1, 2, 2}},
		&Layer{2, 2, []Pixel{2, 2, 1, 2}},
		&Layer{2, 2, []Pixel{0, 0, 0, 0}},
	}
	expected := &Layer{2, 2, []Pixel{0, 1, 1, 0}}
	flat := Flatten(layers)
	if !LayerEquals(expected, flat) {
		t.Errorf("expected %v but got %v", expected, flat)
	}
}

// LayerEquals returns true if the two given layers are the same, false
// otherwise.
func LayerEquals(a, b *Layer) bool {
	if a.width != b.width || a.height != b.height {
		return false
	}
	for i := range a.pixels {
		if a.pixels[i] != b.pixels[i] {
			return false
		}
	}
	return true
}
