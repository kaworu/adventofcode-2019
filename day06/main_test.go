package main

import (
	"strings"
	"testing"
)

func TestOrbitCount(t *testing.T) {
	input := `COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
`
	expected := struct {
		direct, indirect int
	}{11, 31}

	uom, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parsing error")
	}
	direct, indirect := uom.OrbitCount()
	if direct != expected.direct {
		t.Errorf("expected %v direct orbits, got %v", expected.direct, direct)
	}
	if indirect != expected.indirect {
		t.Errorf("expected %v indirect orbits, got %v", expected.indirect, indirect)
	}
}
