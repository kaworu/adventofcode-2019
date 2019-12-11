package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Digits6 represent number composed of six digits.
// NOTE: Digits6 are stored least significant digit first.
type Digits6 [6]uint8

// Password are six digits number obeying all the criteria to be a valid
// password according to the Elves key facts.
type Password struct {
	Digits6
	ri int // adjacent digits index such as Digits6[ri] == Digits6[ri - 1]
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
// Since Password obey all the criteria and the given lower bound min may not,
// the returned Password is the minimum Digits6 value that satisfy
// password.CompareTo(min) >= 0. In other words, the returned Password Digits6
// is the closest valid Password greater than or equals to min.
func NewPassword(min *Digits6) *Password {
	p := &Password{}
	// copy min over p.Digits6 in a way that satisfy the fourth rule, i.e.
	// going from the highest significant digit to the lowest the digit never
	// decrease.
	p.Digits6[len(p.Digits6)-1] = min[len(min)-1]
	for i := len(p.Digits6) - 1; i > 0; i-- {
		if min[i-1] < min[i] {
			p.Digits6[i-1] = min[i]
		} else {
			p.Digits6[i-1] = min[i-1]
		}
		if i > p.ri && p.Digits6[i] == p.Digits6[i-1] {
			// We found two equals adjacent digits, keep track of them in ri.
			p.ri = i
		}
	}
	if p.ri == 0 {
		// Here p.Digits6 doesn't satisfy the two equals adjacent digits. We
		// have p.Digits6[0] > p.Digits6[1] because they never decrease (by
		// construction) and they can't be equals (otherwise we would have
		// ri = 1). Thus, setting p.Digits6[0] = p.Digits6[1] will set p to the
		// closest valid Password below the one we are looking for. From there
		// we Inc()rease the password to find the target.
		p.Digits6[0] = p.Digits6[1]
		p.ri = 1
		p.Inc()
	}
	return p
}

// Inc set the Password to the next valid Password, i.e. the smallest Password
// satisfying all the criteria greater than the current Password.  If p is
// exactly the maximum Digits6 value (i.e. 999999), then Inc will wrap p around
// and set it to 000000.
func (p *Password) Inc() {
	// Find the digit to increase. Note that an incrementing loop is used here
	// instead of range so that 999999 will set i = len(p.Digits6).
	var hi uint8 // the digit value we increased to
	var i int
	for i = 0; i < len(p.Digits6); i++ {
		if p.Digits6[i] < 9 {
			hi = p.Digits6[i] + 1
			p.Digits6[i] = hi
			break
		}
	}

	// Now i is the index of the digit increased, and hi hold its value. We can
	// set all the least significant digits than i to hi to satisfy the fourth
	// rule. Note that when p == 999999 then i = len(p.Digits6) and hi = 0,
	// effectively wrapping around and setting p.Digits6 to 000000.
	for j := 0; j < i; j++ {
		p.Digits6[j] = hi
	}

	// From here p.Digits6 has been increased according to the fourth rule. We
	// still need to enforce the third rule, i.e. two equals adjacent digits.

	// Set i to be the index of the highest digit changed.
	if i == len(p.Digits6) {
		i--
	}

	// If we changed at least the first digit from the previously two equals
	// adjacent digits we need to check for the third rule.
	if i >= p.ri-1 {
		switch {
		// check if we ith digit was increased to exactly the same value as
		// (i+1)th digit.
		case i+1 < len(p.Digits6) && p.Digits6[i] == p.Digits6[i+1]:
			p.ri = i + 1
		// If more than one digit was changed, we know we have
		// p.Digits6[i] == p.Digits6[i-1] == hi. Thus, i is the highest digit
		// index of the two equals adjacent digits.
		case i > 0:
			p.ri = i
		default: // i == 0
			// Here we only increased the least significant digit and it
			// doesn't match p.Digits6[1] (otherwise we would have fallen into
			// the first case. We're not a valid Password and so we need to
			// Inc()rease more.
			p.Inc()
			return
		}
	}
}

// Compute and display the count of Password within the range given on stdin.
func main() {
	bounds, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}
	n := 0
	for p := NewPassword(&bounds.Min); p.CompareTo(&bounds.Max) <= 0; p.Inc() {
		n++
	}
	fmt.Printf("There %d different passwords.\n", n)
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
