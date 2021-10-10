package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const (
	// Add is the addition opcode.
	Add = 1
	// Mult is the multiplication opcode.
	Mult = 2
	// Halt terminate the program.
	Halt = 99
)

const (
	// Output is the index of the program result in memory.
	Output = 0
	// Noun is the value placed in address 1.
	Noun = 1
	// Verb is the value placed in address 2.
	Verb = 2
)

// Intcode is a value in the computer's memory.
type Intcode int64

// Memory represent the state of an Intcode computer.
type Memory []Intcode

// Copy return an copy of the program.
func (mem Memory) Copy() Memory {
	cpy := make(Memory, len(mem))
	copy(cpy, mem)
	return cpy
}

// Setup prepare the program for execution by setting its Noun and Verb
// indices.
func (mem Memory) Setup(noun, verb Intcode) {
	mem[Noun] = noun
	mem[Verb] = verb
}

// Execute run the Intcode gravity assist program.
// It returns an error if an unexpected opcode is encountered.
func (mem Memory) Execute() error {
	var pc int64 = 0 // instruction pointer
	for {
		opcode := mem[pc]
		switch opcode {
		case Add:
			lpos := mem[pc+1]
			rpos := mem[pc+2]
			dest := mem[pc+3]
			mem[dest] = mem[lpos] + mem[rpos]
		case Mult:
			lpos := mem[pc+1]
			rpos := mem[pc+2]
			dest := mem[pc+3]
			mem[dest] = mem[lpos] * mem[rpos]
		case Halt:
			return nil
		default:
			return fmt.Errorf("unsupported opcode: %d", opcode)
		}
		pc += 4
	}
}

// main execute the Intcode program given on stdin and output the value left at
// position 0 after the program halts.
func main() {
	// parse the puzzle input, i.e. the initial state of the Intcode program.
	initial, err := Parse(os.Stdin)
	if err != nil {
		log.Fatalf("input error: %s\n", err)
	}

	// execute each possible noun and verb combination in its own goroutine
	// using a couple of channel so that programs finding an interesting output
	// can be collected.
	alarm := make(chan Memory, 1)
	landing := make(chan Memory, 1)
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			mem := initial.Copy()
			go func(noun, verb Intcode) { // capture noun and verb
				mem.Setup(noun, verb)
				err := mem.Execute()
				switch {
				case err != nil:
					// we crashed the program. Maybe that's because of this
					// noun and verb combination, so we ignore it.
				case noun == 12 && verb == 2:
					alarm <- mem // 1202 program alarm
				case mem[Output] == 19690720:
					landing <- mem // Moon landing by Appollo 11
				}
			}(Intcode(noun), Intcode(verb))
		}
	}
	fst := <-alarm
	snd := <-landing

	fmt.Printf("The value left at position 0 after the program halts is %d,\n", fst[Output])
	fmt.Printf("and when the output is 19690720: 100 * noun + verb = %d.\n", 100*snd[Noun]+snd[Verb])
}

// Parse an Intcode program.
// It returns the parsed Intcode program's initial memory and any read or
// conversion error encountered.
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
