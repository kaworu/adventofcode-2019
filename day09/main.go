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
	// RelativeBaseOffset adjust the relative base.
	RelativeBaseOffset = 9
	// Halt terminate the program.
	Halt = 99
)

const (
	// Position mode where an Intcode is an address.
	Position = 0
	// Immediate mode where an Intcode is an immediate value.
	Immediate = 1
	// Relative mode where an Intcode is an offset with respect to the relative
	// base offset.
	Relative = 2
)

type (
	// Opcode represent the Intcode Computer operation codes.
	Opcode uint8
	// Mode represent the Intcode Computer parameter modes.
	Mode uint8
	// Intcode is a value in the Computer's memory.
	Intcode int64
)

type (
	// Input is the slice from where the Read instruction will get its data.
	Input []Intcode
	// Output is the slice where the Write instruction will append its data.
	Output []Intcode
)

// Computer implement a complete Intcode computer.
type Computer struct {
	mem []Intcode // memory
	pc  int       // instruction pointer
	rbo int       // relative base offset
}

// expand the Computer's memory with zero values up to i. Once it returns, it
// is guaranteed that c.mem[i] will not be out of bounds.
func (c *Computer) expand(i int) {
	if i >= len(c.mem) {
		s := nextPow2(i + 1) // the new memory size
		buf := make([]Intcode, s)
		copy(buf, c.mem)
		c.mem = buf
	}
}

// fetch read the value at the address i in the Computer's memory. It returns
// the value read and an error when the address is invalid.
func (c *Computer) fetch(i int) (Intcode, error) {
	if i < 0 {
		return 0, fmt.Errorf("invalid memory read at %d", i)
	}
	c.expand(i)
	return c.mem[i], nil
}

// put write the given value at the address i in the Computer's memory. It
// returns the value written and an error when the address is invalid.
func (c *Computer) put(i int, val Intcode) (Intcode, error) {
	if i < 0 {
		return 0, fmt.Errorf("invalid memory write at %d", i)
	}
	c.expand(i)
	c.mem[i] = val
	return val, nil
}

// load an Intcode parameter honoring the given mode. It returns the value read
// along with any address or mode error encountered.
func (c *Computer) load(mode Mode, i int) (Intcode, error) {
	param, err := c.fetch(i)
	if err != nil {
		return 0, err
	}
	switch mode {
	case Position:
		return c.fetch(int(param))
	case Immediate:
		return param, nil
	case Relative:
		return c.fetch(c.rbo + int(param))
	default:
		return 0, fmt.Errorf("invalid mode %v", mode)
	}
}

// store a value at an Intcode parameter honoring the given mode. It returns
// the value written along with any address or mode error encountered.
func (c *Computer) store(mode Mode, i int, val Intcode) (Intcode, error) {
	param, err := c.fetch(i)
	if err != nil {
		return 0, err
	}
	switch mode {
	case Position:
		return c.put(int(param), val)
	case Relative:
		return c.put(c.rbo+int(param), val)
	default:
		// NOTE from day05: Parameters that an instruction writes to will
		// never be in immediate mode.
		return 0, fmt.Errorf("invalid mode %v", mode)
	}
}

// instruction returns the current operation code, mode of the first parameter,
// mode of the second parameter, mode of the third parameter, along with any
// memory read error encountered.
func (c *Computer) instruction() (Opcode, Mode, Mode, Mode, error) {
	i, err := c.fetch(c.pc)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	opcode := Opcode(i % 100)
	m1 := Mode(i / 100 % 10)
	m2 := Mode(i / 1000 % 10)
	m3 := Mode(i / 10000 % 10)
	return opcode, m1, m2, m3, nil
}

// Execute run the Intcode program on the Computer with the provided Input. It
// returns the Output generated by the evaluation and an error on failure.
func (c *Computer) Execute(program []Intcode, in Input) (Output, error) {
	// Setup the initial memory and registers
	c.mem = make([]Intcode, len(program))
	copy(c.mem, program)
	c.pc = 0
	c.rbo = 0
	var output []Intcode
	for {
		opcode, m1, m2, m3, err := c.instruction()
		if err != nil {
			return nil, err
		}
		switch opcode {
		case Add, Mult, LessThan, Equals: // binary operators
			lhs, err := c.load(m1, c.pc+1)
			if err != nil {
				return nil, err
			}
			rhs, err := c.load(m2, c.pc+2)
			if err != nil {
				return nil, err
			}
			var result Intcode
			switch opcode {
			case Add:
				result = lhs + rhs
			case Mult:
				result = lhs * rhs
			case LessThan:
				if lhs < rhs {
					result = 1
				}
			case Equals:
				if lhs == rhs {
					result = 1
				}
			}
			_, err = c.store(m3, c.pc+3, result)
			if err != nil {
				return nil, err
			}
			c.pc += 4
		case JumpIfTrue, JumpIfFalse: // jumps opcodes
			cond, err := c.load(m1, c.pc+1)
			if err != nil {
				return nil, err
			}
			addr, err := c.load(m2, c.pc+2)
			if err != nil {
				return nil, err
			}
			var jump bool
			switch opcode {
			case JumpIfTrue:
				jump = cond != 0
			case JumpIfFalse:
				jump = cond == 0
			}
			if jump {
				c.pc = int(addr)
			} else {
				c.pc += 3
			}
		case RelativeBaseOffset:
			off, err := c.load(m1, c.pc+1)
			if err != nil {
				return nil, err
			}
			c.rbo += int(off)
			c.pc += 2
		case Read:
			if len(in) == 0 {
				return nil, fmt.Errorf("Read instruction: empty input")
			}
			_, err = c.store(m1, c.pc+1, in[0])
			if err != nil {
				return nil, err
			}
			in = in[1:] // "pop" the input we've read.
			c.pc += 2
		case Write:
			src, err := c.load(m1, c.pc+1)
			if err != nil {
				return nil, err
			}
			output = append(output, src)
			c.pc += 2
		case Halt:
			return output, nil
		default:
			return nil, fmt.Errorf("unsupported opcode %d", opcode)
		}
	}
}

// main execute the latest Basic Operation Of System Test Intcode program given
// on stdin and display its keycode output.
func main() {
	boost, err := Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}

	var c Computer
	output, err := c.Execute(boost, Input{1})
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution error: %s\n", err)
		os.Exit(1)
	}
	if len(output) != 1 {
		fmt.Fprintf(os.Stderr, "unexpected output from the BOOST program: %v.\n", output)
		os.Exit(1)
	}
	fmt.Printf("The BOOST keycode is %v,\n", output[0])

	output, err = c.Execute(boost, Input{2})
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution error: %s\n", err)
		os.Exit(1)
	}
	if len(output) != 1 {
		fmt.Fprintf(os.Stderr, "unexpected output from the BOOST program: %v.\n", output)
		os.Exit(1)
	}
	fmt.Printf("and the coordinates of the distress signal are %v.\n", output[0])
}

// Parse an Intcode program.
// It returns the parsed Intcode program and any read or convertion error
// encountered.
func Parse(r io.Reader) ([]Intcode, error) {
	var prog []Intcode
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanIntcodes)
	for scanner.Scan() {
		ic, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return nil, err
		}
		prog = append(prog, Intcode(ic))
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

// nextPow2 returns the smallest power of two greater or equal to n.
func nextPow2(n int) int {
	if n == 0 || n&(n-1) == 0 {
		return n
	}
	count := 0
	for n != 0 {
		n >>= 1
		count++
	}
	return 1 << count
}
