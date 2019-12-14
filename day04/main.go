package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Digits6 represent number composed of six digits. By definition, a Digits6
// always satisfy the first criterion.
// NOTE: Digits6 are stored least significant digit first.
type Digits6 [6]uint8

// Password are six digits number satisfying the fourth criterion, i.e. going
// from the highest digit to the lowest the digits never decrease.
type Password struct {
	Digits6
	// ri is the highest digits index satisfying Digits6[ri] == Digits6[ri-1].
	// When there are no equal adjacent digits then ri == 0.
	ri int
}

// Bounds represent the range given in the puzzle input.
type Bounds struct {
	Min, Max Digits6
}

// Set will change digits to represent the value n.
// Note that since Digits6 are limited to six digits, only the six least
// significant digits of n will be used (in base 10).
func (ds *Digits6) Set(n uint64) {
	for i := range ds {
		ds[i] = uint8(n % 10)
		n /= 10
	}
}

// CompareTo returns an integer comparing two Digits6.
// It returns 0 if ds==es, -1 if ds < es, and +1 if ds > es.
func (ds *Digits6) CompareTo(es *Digits6) int {
	// NOTE: going from the end of the array to the start since Digits6 are
	// least significant digit first.
	for i := len(ds) - 1; i >= 0; i-- {
		switch {
		case ds[i] < es[i]:
			return -1
		case ds[i] > es[i]:
			return +1
		}
	}
	return 0
}

// NewPassword create a Password from a lower bound Digits6 number.
// Since Password must satisfy the fourth criterion and the given lower bound
// min may not, the returned Password is the minimum valid value satisfying
// climb.CompareTo(min) >= 0. In other words, the returned Password is the
// closest valid Password greater than or equals to min.
func NewPassword(min *Digits6) *Password {
	p := &Password{}
	// copy min over p.Digits6 in a way that satisfy the fourth rule, i.e.
	// going from the highest significant digit to the lowest the digit never
	// decrease.
	p.Digits6[len(p.Digits6)-1] = min[len(min)-1]
	for i := len(p.Digits6) - 2; i >= 0; i-- {
		if min[i] < min[i+1] {
			p.Digits6[i] = min[i+1]
		} else {
			p.Digits6[i] = min[i]
		}
		// When we find two equal adjacent digits, keep track of them in ri.
		if p.ri == 0 && p.Digits6[i] == p.Digits6[i+1] {
			p.ri = i + 1
		}
	}
	return p
}

// Inc set the Password to the next valid Password, i.e. the smallest Password
// satisfying the fourth criterion greater than the current Password.
// NOTE: When p is the maximum Digits6 value (i.e. 999999), then Inc will wrap
// p around set it to 000000.
func (p *Password) Inc() {
	// Find the digit to increase. Note that an incrementing loop is used here
	// instead of range so that 999999 will set i = len(p.Digits6).
	var hi uint8 // the digit value increased to
	var i int
	for i = 0; i < len(p.Digits6); i++ {
		if p.Digits6[i] < 9 {
			hi = p.Digits6[i] + 1
			p.Digits6[i] = hi
			break
		}
	}

	// Now i is the index of the digit increased and hi hold its value. We can
	// set all the least significant digits than i to hi to satisfy the fourth
	// criterion. Note that when p == 999999 then i = len(p.Digits6) and
	// hi = 0, effectively wrapping around and setting p.Digits6 to 000000.
	for j := 0; j < i; j++ {
		p.Digits6[j] = hi
	}

	// From here p.Digits6 has been increased according to the fourth
	// criterion. We do some equal adjacent digits bookkeeping to ease password
	// detection.

	// Set i to be the index of the highest digit changed.
	if i == len(p.Digits6) {
		i--
	}

	// When we changed one value from the previously two equal adjacent digits
	// we need to update ri. Note that ri == 0, i.e. there was not double
	// previously, will follow the same codepath.
	if i >= p.ri-1 {
		switch {
		// check if we ith digit was increased to exactly the same value as
		// (i+1)th digit.
		case i+1 < len(p.Digits6) && p.Digits6[i] == p.Digits6[i+1]:
			p.ri = i + 1
		// Here we have either i == 0 or p.Digits6[i] == p.Digits6[i-1] == hi.
		// Thus, i is the highest digit index of the two equals adjacent digits
		// (or zero if there is no double).
		default:
			p.ri = i
		}
	}
}

// HasRepeatingDigits returns true when p has at least two equal adjacent
// digits, false otherwise.
func (p *Password) HasRepeatingDigits() bool {
	return p.ri > 0
}

// HasDouble returns true when p has two equal adjacent digits that are not
// part of a larger group, false otherwise.
func (p *Password) HasDouble() bool {
	// value of the highest repeating digit.
	v := p.Digits6[p.ri]
	// when p.ri > 0 then we already have two equal adjacent digits.
	n := 2
	for i := p.ri - 2; i >= 0; i-- {
		switch {
		case p.Digits6[i] == v: // same number as the previous.
			n++
		case n == 2: // different from the previous and we just seen a double.
			return true
		default: // "reset" the double counter
			v = p.Digits6[i]
			n = 1
		}
	}
	// From here, our last chance to satisfy the rule is to have the two least
	// significant digits forming a double.
	return n == 2 && p.Digits6[0] == p.Digits6[1]
}

// Compute and display the count of Password within the range given on stdin.
func main() {
	bounds, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}
	v1 := 0 // count of version 1 passwords, satisfying HasRepeatingDigits()
	v2 := 0 // count of version 2 passwords, satisfying HasDouble()
	for p := NewPassword(&bounds.Min); p.CompareTo(&bounds.Max) <= 0; p.Inc() {
		if p.HasRepeatingDigits() {
			v1++
		}
		if p.HasDouble() {
			v2++
		}
	}
	fmt.Printf("There are %d different passwords with repeating digits,\n", v1)
	fmt.Printf("And %d different passwords with double not part of a larger group.\n", v2)
}

// Parse a bound description formatted as "min-max".
// It returns the Bounds and any read or parsing error encountered.
func Parse(r io.Reader) (*Bounds, error) {
	bounds := &Bounds{}
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return bounds, fmt.Errorf("empty input")
	}
	line := scanner.Text()
	if err := scanner.Err(); err != nil {
		return bounds, err
	}
	numbers := strings.Split(line, "-")
	if len(numbers) != 2 {
		return bounds, fmt.Errorf("%s: not a range", line)
	}
	min, err := strconv.ParseUint(numbers[0], 10, 64)
	if err != nil {
		return bounds, err
	}
	max, err := strconv.ParseUint(numbers[1], 10, 64)
	if err != nil {
		return bounds, err
	}
	if max > 999999 || max < min {
		return bounds, fmt.Errorf("%s: invalid range", line)
	}
	bounds.Min.Set(min)
	bounds.Max.Set(max)
	return bounds, nil
}
