package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type (
	Mass   int  // The mass of something.
	Module Mass // A module to be launched.
	Fuel   Mass // Required to launch modules.
)

// Return the amount of Fuel required to launch this Module.
func (mod Module) FuelRequired() Fuel {
	return Fuel(mod/3 - 2)
}

// Compute and display the sum of the fuel requirements for all of the modules
// given on stdin.
func main() {
	modules, err := parse(os.Stdin)
	if err != nil {
		panic(err)
	}
	total := Fuel(0)
	for _, mod := range modules {
		total += mod.FuelRequired()
	}
	fmt.Printf("The sum of the fuel requirements is %d.\n", total)
}

// Parse one Module per line of input.
func parse(r io.Reader) ([]Module, error) {
	modules := make([]Module, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		i, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		modules = append(modules, Module(i))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return modules, nil
}
