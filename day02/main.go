package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	Add  = 1  // Add is the addition opcode.
	Mult = 2  // Mult is the multiplication opcode.
	Halt = 99 // Halt terminate the program.
)

// Intcode is a gravity assist program composed by integers.
type Intcode []int

// Alarm restore the program to the "1202 program alarm" state it had just
// before the last computer caught fire.
func (program Intcode) Alarm() {
	program[1] = 12
	program[2] = 2
}

// Execute run the Intcode gravity assist program.
// It returns an error if an unexpected opcode is encountered.
func (program Intcode) Execute() error {
	position := 0
Loop:
	for {
		opcode := program[position]
		switch opcode {
		case Add:
			lpos := program[position+1]
			rpos := program[position+2]
			dest := program[position+3]
			program[dest] = program[lpos] + program[rpos]
		case Mult:
			lpos := program[position+1]
			rpos := program[position+2]
			dest := program[position+3]
			program[dest] = program[lpos] * program[rpos]
		case Halt:
			break Loop
		default:
			return errors.New(fmt.Sprintf("unsupported opcode %d", opcode))
		}
		position += 4
	}
	return nil
}

// main execute the Intcode program given on stdin and output The value left at
// position 0 after the program halts.
func main() {
	program, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}
	program.Alarm()
	if err = program.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "intcode error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("The value left at position 0 after the program halts is %d.\n", program[0])
}

// Parse an Intcode program.
// It returns the parsed Intcode program and any read of convertion error
// encountered.
func Parse(r io.Reader) (Intcode, error) {
	var program Intcode
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanIntcodes)
	for scanner.Scan() {
		ic, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		program = append(program, ic)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return program, nil
}

// ScanIntcodes is a split function for Scanner.
// It returns each Intcode of text.
func ScanIntcodes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Heavily inspired by ScanLines, the default Scanner split function.
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, ",\n"); i >= 0 {
		// We have a full Intcode
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated Intcode. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
