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

// OrbitCount compute and return the Center Of Mass's total of direct and
// indirect orbits.
func (uom UniversalOrbitMap) OrbitCount() (direct, indirect int) {
	if com, ok := uom["COM"]; ok {
		direct, indirect = com.OrbitCount(0)
	} else {
		direct, indirect = -1, -1
	}
	return
}

// OrbitCount compute and return the b's total of direct and indirect orbits.
// The depth argument is the distance between b and the Center of Mass.
func (b *Body) OrbitCount(depth int) (direct, indirect int) {
	// we have a direct orbit iff we're orbiting around another Body, i.e. if
	// we're not the COM.
	if b.orbits != nil {
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

// main parse the universal orbit map, then compute and display the total of
// direct and indirect orbits.
func main() {
	uom, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}
	direct, indirect := uom.OrbitCount()
	fmt.Printf("the total number of direct and indirect orbits is %v.\n", direct+indirect)
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
