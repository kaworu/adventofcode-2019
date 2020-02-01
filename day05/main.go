package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	// Add is the addition opcode.
	Add = 1
	// Mult is the multiplication opcode.
	Mult = 2
	// Read take a single integer as input.
	Read = 3
	// Write ouputs the value of its only parameter.
	Write = 4
	// JumpIfTrue update the instruction pointer of the first parameter is
	// non-zero.
	JumpIfTrue = 5
	// JumpIfFalse update the instruction pointer of the first parameter is
	// zero.
	JumpIfFalse = 6
	// LessThan is the "lesser than" comparison opcode.
	LessThan = 7
	// Equals is the equality comparison opcode.
	Equals = 8
	// Halt terminate the program.
	Halt = 99
)

const (
	// Position mode where an Intcode is an address.
	Position = 0
	// Immediate mode where an Intcode is an immediate value.
	Immediate = 1
)

// Intcode is a value in the computer's memory.
type Intcode int64

type (
	// Memory represent the state of an Intcode computer.
	Memory []Intcode
	// Input is the slice from where the Read instruction will get its data.
	Input []Intcode
	// Output is the slice where the Write instruction will append its data.
	Output []Intcode
)

// Copy return an copy of the program.
func (mem Memory) Copy() Memory {
	cpy := make(Memory, len(mem))
	copy(cpy, mem)
	return cpy
}

// Execute run the Intcode program with the provided Input. It returns the
// Output generated by the evaluation and an error on failure.
func (mem Memory) Execute(in Input) (Output, error) {
	var out Output
	var pc int64 = 0 // instruction pointer
	for {
		opcode := mem[pc] % 100
		lhsm := (mem[pc] / 100) % 10  // first operand mode
		rhsm := (mem[pc] / 1000) % 10 // second operand mode
		switch opcode {
		case Add, Mult, LessThan, Equals, JumpIfTrue, JumpIfFalse:
			lhs := mem[pc+1]
			if lhsm == Position {
				lhs = mem[lhs]
			}
			rhs := mem[pc+2]
			if rhsm == Position {
				rhs = mem[rhs]
			}
			dst := mem[pc+3]
			switch opcode {
			case Add:
				mem[dst] = lhs + rhs
				pc += 4
			case Mult:
				mem[dst] = lhs * rhs
				pc += 4
			case LessThan:
				if lhs < rhs {
					mem[dst] = 1
				} else {
					mem[dst] = 0
				}
				pc += 4
			case Equals:
				if lhs == rhs {
					mem[dst] = 1
				} else {
					mem[dst] = 0
				}
				pc += 4
			case JumpIfTrue:
				if lhs != 0 {
					pc = int64(rhs)
				} else {
					pc += 3
				}
			case JumpIfFalse:
				if lhs == 0 {
					pc = int64(rhs)
				} else {
					pc += 3
				}
			}
		case Read:
			if len(in) == 0 {
				return nil, fmt.Errorf("Read instruction: empty input")
			}
			// NOTE: Read is always in position mode
			dst := mem[pc+1]
			mem[dst] = in[0]
			in = in[1:] // "pop" the input we've read.
			pc += 2
		case Write:
			src := mem[pc+1]
			// NOTE: Write use the first mode for its only operand.
			if lhsm == Position {
				src = mem[src]
			}
			out = append(out, src)
			pc += 2
		case Halt:
			return out, nil
		default:
			return nil, fmt.Errorf("unsupported opcode %d", opcode)
		}
	}
}

// main execute the Thermal Environment Supervision Terminal Intcode diagnostic
// program given on stdin and display its output.
func main() {
	// parse the puzzle input, i.e. the initial state of the Intcode program.
	mem, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}

	// air conditioner unit
	output, err := mem.Copy().Execute(Input{1})
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("The TEST diagnostic program output for the air conditioner unit is: %v\n", output)

	// thermal radiator controller
	output, err = mem.Copy().Execute(Input{5})
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("The TEST diagnostic program output for the thermal radiator controller is: %v\n", output)
}

// Parse an Intcode program.
// It returns the parsed Intcode program's initial memory and any read or
// convertion error encountered.
func Parse(r io.Reader) (Memory, error) {
	var mem Memory
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanIntcodes)
	for scanner.Scan() {
		ic, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return nil, err
		}
		mem = append(mem, Intcode(ic))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return mem, nil
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
