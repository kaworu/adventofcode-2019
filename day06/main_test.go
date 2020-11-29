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
	want := struct {
		direct, indirect int
	}{11, 31}

	uom, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parsing error")
	}
	direct, indirect := uom.OrbitCount()
	if direct != want.direct {
		t.Errorf("got %v direct orbits; want %v", direct, want.direct)
	}
	if indirect != want.indirect {
		t.Errorf("got %v inddirect orbits; want %v", indirect, want.indirect)
	}
}

func TestOrbitalTransfers(t *testing.T) {
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
K)YOU
I)SAN
`
	want := 4

	uom, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parsing error: %v", err)
	}
	you, ok := uom["YOU"]
	if !ok {
		t.Fatalf("Parsing error: YOU can not be found")
	}
	san, ok := uom["SAN"]
	if !ok {
		t.Fatalf("Parsing error: SAN can not be found")
	}
	if n := you.OrbitalTransfers(san); n != want {
		t.Errorf("got %v orbital transfers; want %v", n, want)
	}
}
