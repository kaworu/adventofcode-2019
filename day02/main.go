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

// Intcode opcodes
const (
	Add  = 1  // Add is the addition opcode.
	Mult = 2  // Mult is the multiplication opcode.
	Halt = 99 // Halt terminate the program.
)

// Intcode memory indices
const (
	Output = 0    // Output is the index of the program result in memory.
	Noun   = iota // Noun is the value placed in adress 1
	Verb   = iota // Verb is the value placed in adress 2
)

// Intcode represent the memory state of an Intcode computer.
type Intcode []int

// Copy return an copy of the program.
func (prog Intcode) Copy() Intcode {
	clone := make(Intcode, len(prog))
	copy(clone, prog)
	return clone
}

// Setup prepare the program for execution by setting its Noun and Verb
// indices.
func (prog Intcode) Setup(noun, verb int) {
	prog[Noun] = noun
	prog[Verb] = verb
}

// Execute run the Intcode gravity assist program.
// It returns an error if an unexpected opcode is encountered.
func (prog Intcode) Execute() error {
	ip := 0 // instruction pointer
Loop:
	for {
		opcode := prog[ip]
		switch opcode {
		case Add:
			lpos, rpos, dest := prog[ip+1], prog[ip+2], prog[ip+3]
			prog[dest] = prog[lpos] + prog[rpos]
		case Mult:
			lpos, rpos, dest := prog[ip+1], prog[ip+2], prog[ip+3]
			prog[dest] = prog[lpos] * prog[rpos]
		case Halt:
			break Loop
		default:
			return fmt.Errorf("unsupported opcode %d", opcode)
		}
		ip += 4
	}
	return nil
}

// main execute the Intcode program given on stdin and output the value left at
// position 0 after the program halts.
func main() {
	// parse the puzzle input, i.e. the initial state of the Intcode program.
	initial, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}

	// execute each possible noun and verb combination in its own goroutine
	// using a couple of channel so that programs finding an interesting output
	// can be collected.
	alarm, landing := make(chan Intcode, 1), make(chan Intcode, 1)
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			go func(noun, verb int) { // capture noun and verb
				prog := initial.Copy()
				prog.Setup(noun, verb)
				err := prog.Execute()
				switch {
				case err != nil:
					// we crashed the program. Maybe that's because of this
					// noun and verb combination, so we ignore it.
				case noun == 12 && verb == 2:
					alarm <- prog // 1202 program alarm
				case prog[Output] == 19690720:
					landing <- prog // Moon landing by Appollo 11
				}
			}(noun, verb)
		}
	}
	fst, snd := <-alarm, <-landing

	fmt.Printf("The value left at position 0 after the program halts is %d,\n", fst[Output])
	fmt.Printf("and when the output is 19690720: 100 * noun + verb = %d.\n", 100*snd[Noun]+snd[Verb])
}

// Parse an Intcode program.
// It returns the parsed Intcode program's initial memory and any read of
// convertion error encountered.
func Parse(r io.Reader) (Intcode, error) {
	var prog Intcode
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanIntcodes)
	for scanner.Scan() {
		ic, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		prog = append(prog, ic)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return prog, nil
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
