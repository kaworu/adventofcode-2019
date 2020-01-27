package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Body is a "small Solar System body", i.e. an object in space. A Body orbits
// another Body and they are represented as a tree. The root of the tree is the
// COM (Center of Mass).
// See https://english.stackexchange.com/a/281983
type Body struct {
	desc     string
	orbits   *Body
	orbiting []*Body
}

// UniversalOrbitMap maps bodies names to their entry in the tree.
type UniversalOrbitMap map[string]*Body

// IsCenterOfMass returns true when b is the Center of Mass, false otherwise.
func (b *Body) IsCenterOfMass() bool {
	// When b doesn't orbit another body, then it is the Center of Mass.
	return b.orbits == nil
}

// Depth returns the distance between b and the Center of Mass.
func (b *Body) Depth() int {
	d := 0
	for c := b; !c.IsCenterOfMass(); c = c.orbits {
		d++
	}
	return d
}

// OrbitCount compute and return the b's total of direct and indirect orbits.
// The depth argument is the distance between b and the Center of Mass.
func (b *Body) OrbitCount(depth int) (direct, indirect int) {
	// we have a direct orbit iff we're not the COM.
	if !b.IsCenterOfMass() {
		direct = 1
		indirect = depth - 1
	}

	// Compute recursively the direct and indirect orbits.
	for _, o := range b.orbiting {
		d, i := o.OrbitCount(depth + 1)
		direct += d
		indirect += i
	}

	return
}

// OrbitalTransfers compute and returns the minimum of orbital transfers
// required to move from the object b is orbiting to the object o is orbiting.
// It returns -1 when either b or o is the Center of Mass, or when they are not
// part of the same UniversalOrbitMap.
func (b *Body) OrbitalTransfers(o *Body) int {
	if b.IsCenterOfMass() || o.IsCenterOfMass() {
		return -1
	}

	// let n0, n1 be the objects b (respectively o) is orbiting.
	n0, n1 := b.orbits, o.orbits
	// let n0 be the deepest node.
	d0, d1 := n0.Depth(), n1.Depth()
	diff := d0 - d1
	if diff < 0 {
		n0, n1 = n1, n0
		diff = -diff
	}

	// Bring n0 "up" to the same level as n1.
	for i := 0; i < diff; i++ {
		n0 = n0.orbits
	}

	// Travel "up" from both sides until we get a match.
	for i := 0; n0 != nil && n1 != nil; i++ {
		if n0 == n1 {
			return i*2 + diff
		}
		n0, n1 = n0.orbits, n1.orbits
	}

	return -1
}

// OrbitCount compute and return the Center of Mass's total of direct and
// indirect orbits. When the map has no Center of Mass, (-1, -1) is returned.
func (uom UniversalOrbitMap) OrbitCount() (direct, indirect int) {
	if com, ok := uom["COM"]; ok {
		direct, indirect = com.OrbitCount(0)
	} else {
		direct, indirect = -1, -1
	}
	return
}

// main parse the universal orbit map, then compute and display the total of
// direct and indirect orbits.
func main() {
	uom, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}
	you, ok := uom["YOU"]
	if !ok {
		fmt.Fprint(os.Stderr, "YOU not found\n")
		os.Exit(1)
	}
	san, ok := uom["SAN"]
	if !ok {
		fmt.Fprint(os.Stderr, "SAN not found\n")
		os.Exit(1)
	}
	direct, indirect := uom.OrbitCount()
	distance := you.OrbitalTransfers(san)
	fmt.Printf("the total number of direct and indirect orbits is %v,\n", direct+indirect)
	fmt.Printf("and the minimum of orbital transfers required is %v.\n", distance)
}

// Parse a map of the local orbits.
// It returns the UniversalOrbitMap and any read or parsing error encountered.
func Parse(r io.Reader) (UniversalOrbitMap, error) {
	uom := make(UniversalOrbitMap)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		desc := strings.Split(line, ")")
		if len(desc) != 2 {
			return nil, fmt.Errorf("%s: invalid orbit line", line)
		}
		// c)o means o directly orbits c.
		cdesc := desc[0]
		odesc := desc[1]
		c, ok := uom[cdesc]
		if !ok {
			c = &Body{desc: cdesc}
		}
		o, ok := uom[odesc]
		if !ok {
			o = &Body{desc: odesc, orbits: c}
		}
		if o.orbits == nil {
			o.orbits = c
		}
		if o.orbits != c {
			return nil, fmt.Errorf("expected %v to orbits %v, got %v instead", o.desc, c.desc, o.orbits.desc)
		}
		o.orbits = c
		c.orbiting = append(c.orbiting, o)
		uom[c.desc] = c
		uom[o.desc] = o
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return uom, nil
}
