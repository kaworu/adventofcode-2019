package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Mass is the weight of something.
type Mass int64

// FuelRequired compute the amount of fuel required to carry the given mass. //
// It returns the Mass of the fuel required.
func FuelRequired(m Mass) Mass {
	return m/3 - 2
}

// TotalFuelRequired compute the amount of fuel required to launch all the
// given modules.
// It returns the total Mass of fuel required for the modules and the Mass of
// fuel required for the fuel.
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

// main compute and display the sum of the fuel requirements for all of the
// modules given on stdin.
func main() {
	modules, err := Parse(os.Stdin)
	if err != nil {
		log.Fatalf("input error: %s\n", err)
	}
	mf, ff := TotalFuelRequired(modules...)
	fmt.Printf("The sum of the fuel requirements is %d,\n", mf)
	fmt.Printf("and when also taking into account the mass of the added fuel it is %d.\n", mf+ff)
}

// Parse read the given input to produce a slice of Mass.
// It returns the parsed collection of Mass and any read of conversion error
// encountered.
func Parse(r io.Reader) ([]Mass, error) {
	var modules []Mass
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		i, err := strconv.ParseUint(line, 10, 63) // 63 bit size fit in int64
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
