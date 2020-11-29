package main

import (
	"strings"
	"testing"
)

var encoded = "0222112222120000"
var layers = []Layer{
	Layer{2, 2, []Pixel{0, 2, 2, 2}},
	Layer{2, 2, []Pixel{1, 1, 2, 2}},
	Layer{2, 2, []Pixel{2, 2, 1, 2}},
	Layer{2, 2, []Pixel{0, 0, 0, 0}},
}

func TestParse(t *testing.T) {
	decoded, err := Parse(2, 2, strings.NewReader(encoded))
	if err != nil {
		t.Errorf("Parse(%v) error: %s", encoded, err)
	} else if len(decoded) != len(layers) {
		t.Errorf("got %d layers; want %d", len(decoded), len(layers))
	} else {
		for i := range layers {
			if !LayerEquals(decoded[i], layers[i]) {
				t.Errorf("layers[%d] = %v; want %v", i, decoded[i], layers[i])
			}
		}
	}
}

func TestFlatten(t *testing.T) {
	want := Layer{2, 2, []Pixel{0, 1, 1, 0}}
	flat, err := Flatten(layers)
	if err != nil {
		t.Errorf("Flatten(%v) error: %s", encoded, err)
	} else if !LayerEquals(flat, want) {
		t.Errorf("Flatten(%v) = %v; want = %v", encoded, flat, want)
	}
}

// LayerEquals returns true if the two given layers are the same, false
// otherwise.
func LayerEquals(a, b Layer) bool {
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
