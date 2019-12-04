package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

// The weight of something.
type Mass int

// Returns the amount of fuel required to carry the given mass.
func FuelRequired(m Mass) Mass {
	return m/3 - 2
}

// Return the amount of fuel required to launch all the given modules, along
// with the amount of fuel required for the fuel.
func TotalFuelRequired(modules ...Mass) (mf, ff Mass) {
	for _, module := range modules {
		fuel := FuelRequired(module)
		mf += fuel
		fuel = FuelRequired(fuel)
		for fuel > 0 {
			ff += fuel
			fuel = FuelRequired(fuel)
		}
	}
	return
}

// Compute and display the sum of the fuel requirements for all of the modules
// given on stdin.
func main() {
	modules, err := parse(os.Stdin)
	if err != nil {
		panic(err)
	}
	mf, ff := TotalFuelRequired(modules...)
	fmt.Printf("The sum of the fuel requirements is %d,\n", mf)
	fmt.Printf("and when also taking into account the mass of the added fuel it is %d.\n", mf+ff)
}

// Parse one module per line of input.
func parse(r io.Reader) ([]Mass, error) {
	modules := make([]Mass, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		i, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		modules = append(modules, Mass(i))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return modules, nil
}
