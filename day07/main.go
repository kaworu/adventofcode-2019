package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
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

type (
	// Intcode is a value in the computer's memory.
	Intcode int64
	// Memory represent the state of an Intcode computer.
	Memory []Intcode
	// Amplifier is a Intcode computer.
	Amplifier struct {
		Memory
		Input  <-chan Intcode
		Output chan<- Intcode
		pc     int64 // instruction pointer
	}
)

// mod returns i modulo n.
func mod(i, n int) int {
	return ((i % n) + n) % n
}

// Permutations returns all permutations of the given set.
func Permutations(set []Intcode) [][]Intcode {
	// https://en.wikipedia.org/wiki/Heap%27s_algorithm
	var all [][]Intcode
	var generate func(k int, xs []Intcode)
	generate = func(k int, xs []Intcode) {
		if k == 1 {
			p := make([]Intcode, len(xs))
			copy(p, xs)
			all = append(all, p)
		} else {
			generate(k-1, xs)
			for i := 0; i < k-1; i++ {
				if k%2 == 0 {
					tmp := xs[i]
					xs[i] = xs[k-1]
					xs[k-1] = tmp
				} else {
					tmp := xs[0]
					xs[0] = xs[k-1]
					xs[k-1] = tmp
				}
				generate(k-1, xs)
			}
		}
	}
	// work on a local copy so that set is unmodified after this function
	// returns.
	cpy := make([]Intcode, len(set))
	copy(cpy, set)
	generate(len(cpy), cpy)
	return all
}

// Copy return an copy of the program.
func (mem Memory) Copy() Memory {
	cpy := make(Memory, len(mem))
	copy(cpy, mem)
	return cpy
}

// Execute run the Intcode program in the Amplifier's memory.
// It returns an error on failure.
func (amp *Amplifier) Execute() error {
	mem := amp.Memory
	for {
		opcode := mem[amp.pc] % 100
		lhsm := (mem[amp.pc] / 100) % 10  // first operand mode
		rhsm := (mem[amp.pc] / 1000) % 10 // second operand mode
		switch opcode {
		case Add, Mult, LessThan, Equals, JumpIfTrue, JumpIfFalse:
			lhs := mem[amp.pc+1]
			if lhsm == Position {
				lhs = mem[lhs]
			}
			rhs := mem[amp.pc+2]
			if rhsm == Position {
				rhs = mem[rhs]
			}
			dst := mem[amp.pc+3]
			switch opcode {
			case Add:
				mem[dst] = lhs + rhs
				amp.pc += 4
			case Mult:
				mem[dst] = lhs * rhs
				amp.pc += 4
			case LessThan:
				if lhs < rhs {
					mem[dst] = 1
				} else {
					mem[dst] = 0
				}
				amp.pc += 4
			case Equals:
				if lhs == rhs {
					mem[dst] = 1
				} else {
					mem[dst] = 0
				}
				amp.pc += 4
			case JumpIfTrue:
				if lhs != 0 {
					amp.pc = int64(rhs)
				} else {
					amp.pc += 3
				}
			case JumpIfFalse:
				if lhs == 0 {
					amp.pc = int64(rhs)
				} else {
					amp.pc += 3
				}
			}
		case Read:
			// NOTE: Read is always in position mode
			dst := mem[amp.pc+1]
			mem[dst] = <-amp.Input
			amp.pc += 2
		case Write:
			src := mem[amp.pc+1]
			// NOTE: Write use the first mode for its only operand.
			if lhsm == Position {
				src = mem[src]
			}
			amp.Output <- src
			amp.pc += 2
		case Halt:
			return nil
		default:
			return fmt.Errorf("unsupported opcode %d", opcode)
		}
	}
}

// FeedbackLoop run the provided Amplifier Controller Software on a feedback
// loop of Amplifiers configured according to the given phase setting sequence.
// It returns the last Amplifier's output.
func FeedbackLoop(apc Memory, seq []Intcode) Intcode {
	n := len(seq)
	amps := make([]Amplifier, n)
	// Setup the feedback loop. i.e. each Amplifier to have its input being the
	// previous Amplifier output. The channels are all setup such that the
	// phase setting is provided first. The first Amplifier input is a special
	// case where we have to additionally setup the initial input value zero.
	for i := range amps {
		c := make(chan Intcode, 2)
		c <- seq[i] // phase setting
		if i == 0 {
			c <- 0
		}
		prev := mod(i-1, n)
		amps[i].Memory = apc.Copy()
		amps[i].Input = c
		amps[prev].Output = c
	}
	// Run each Amplifier in a goroutine and wait for all of them to halt.
	// Note that part one require to setup the Amplifiers in series (not in a
	// feedback loop). We blindly trust the Amplifier Controller Software to
	// halt after receiving exactly one input (i.e. without looping) for part
	// one (i.e. when the phase settings are between zero and four inclusive).
	var wg sync.WaitGroup
	wg.Add(n)
	for i := range amps {
		i := i // capture i
		go func() {
			defer wg.Done()
			amps[i].Execute() // FIXME: error are ignored
		}()
	}
	wg.Wait()
	// The first Amplifier input channel is the last Amplifier output channel.
	return <-amps[0].Input
}

// HighestSignal run each possible phase setting sequences permutations in
// a feedback loop of Amplifiers to find signals that can be sent to the
// thruster.
// It return the highest signal.
func HighestSignal(apc Memory, phases []Intcode) Intcode {
	sequences := Permutations(phases)
	signals := make(chan Intcode, len(sequences))
	for _, seq := range sequences {
		seq := seq // capture seq
		go func() {
			signals <- FeedbackLoop(apc, seq)
		}()
	}
	// Read every output values in order to find the greatest one to be
	// returned.
	var max Intcode
	for i := 0; i < len(sequences); i++ {
		x := <-signals
		if i == 0 || x > max {
			max = x
		}
	}
	return max
}

// Main parse the Amplifier Controller Software Intcode program, and then run
// each possible phase setting sequences permutations in a feedback loop of
// Amplifiers to find the highest signal that can be sent to the thruster.
func main() {
	// parse the puzzle input, i.e. the Amplifier Controller Software Intcode
	// program.
	mem, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}
	// part one - in series
	ps := []Intcode{0, 1, 2, 3, 4} // phase settings
	max := HighestSignal(mem, ps)
	fmt.Printf("The highest signal that can be sent to the thrusters using the phase settings %v is %v.\n", ps, max)
	// part two - feedback loop
	ps = []Intcode{5, 6, 7, 8, 9}
	max = HighestSignal(mem, ps)
	fmt.Printf("The highest signal that can be sent to the thrusters using the phase settings %v is %v.\n", ps, max)
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
