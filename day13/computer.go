package main

import (
	"errors"
	"fmt"
)

// Opcodes
const (
	// Add is the addition opcode.
	Add = 1
	// Mult is the multiplication opcode.
	Mult = 2
	// Read take a single integer as input.
	Read = 3
	// Write ouputs the value of its only parameter.
	Write = 4
	// JumpIfTrue update the instruction pointer if the first parameter is
	// non-zero.
	JumpIfTrue = 5
	// JumpIfFalse update the instruction pointer if the first parameter is
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

// Mode
const (
	// Position mode where an Intcode is an address.
	Position = 0
	// Immediate mode where an Intcode is an immediate value.
	Immediate = 1
	// Relative mode where an Intcode is an offset with respect to the relative
	// base offset.
	Relative = 2
)

// Opcode represent the Intcode Computer operation codes.
type Opcode uint8

// Mode represent the Intcode Computer parameter modes.
type Mode uint8

// Intcode is a value in the Computer's memory.
type Intcode int64

// Computer implements a complete Intcode computer.
type Computer struct {
	mem    []Intcode // memory
	input  chan Intcode
	output chan Intcode
	pc     int // instruction pointer
	rbo    int // relative base offset
}

// ErrAbordedExecution is returned by Execute when the computer's input channel
// was close when a Read instruction was reached.
var ErrAbordedExecution = errors.New("execution aborded")

// Execute run the Intcode program on the Computer. It returns an error on
// failure.
func (c *Computer) Execute(program []Intcode) error {
	// Setup the initial memory and registers
	c.mem = make([]Intcode, len(program))
	copy(c.mem, program)
	c.pc = 0
	c.rbo = 0
	for {
		opcode, m1, m2, m3, err := c.instruction()
		if err != nil {
			return err
		}
		switch opcode {
		case Add, Mult, LessThan, Equals: // binary operators
			lhs, err := c.load(m1, c.pc+1)
			if err != nil {
				return err
			}
			rhs, err := c.load(m2, c.pc+2)
			if err != nil {
				return err
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
				return err
			}
			c.pc += 4
		case JumpIfTrue, JumpIfFalse: // jumps opcodes
			cond, err := c.load(m1, c.pc+1)
			if err != nil {
				return err
			}
			addr, err := c.load(m2, c.pc+2)
			if err != nil {
				return err
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
				return err
			}
			c.rbo += int(off)
			c.pc += 2
		case Read:
			r, ok := <-c.input
			if !ok {
				return ErrAbordedExecution
			}
			_, err = c.store(m1, c.pc+1, r)
			if err != nil {
				return err
			}
			c.pc += 2
		case Write:
			code, err := c.load(m1, c.pc+1)
			if err != nil {
				return err
			}
			c.output <- code
			c.pc += 2
		case Halt:
			return nil
		default:
			return fmt.Errorf("unsupported opcode: %d", opcode)
		}
	}
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
		return 0, fmt.Errorf("invalid mode: %v", mode)
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
