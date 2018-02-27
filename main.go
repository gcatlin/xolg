package main

const HEX = "%02X"

const (
	OpConstant byte = iota
	OpConstantX
	OpReturn
)

var InstrSize = []int{
	OpConstant:  2,
	OpConstantX: 4,
	OpReturn:    1,
}

type Value float64

type Chunk struct {
	code      []byte
	constants []Value
	lines     []int
	offsets   []int
}

func NewChunk() *Chunk {
	return &Chunk{
		code:      make([]byte, 0, 16),
		constants: make([]Value, 0, 4),
		lines:     make([]int, 0, 4),
		offsets:   make([]int, 0, 4),
	}
}

func (c *Chunk) addConstant(v Value) int {
	c.constants = append(c.constants, v)
	return len(c.constants) - 1
}

func (c *Chunk) getLine(offset int) int {
	found := false
	low, mid, high := 0, 0, len(c.lines)-1
	for low <= high {
		mid = (low + high) / 2
		if offset < c.offsets[mid] {
			high = mid - 1
		} else if offset > c.offsets[mid] {
			low = mid + 1
		} else {
			found = true
			break
		}
	}
	if !found && offset <= c.offsets[mid] {
		mid--
	}
	return c.lines[mid]
}

func (c *Chunk) write(b byte, line int) {
	n := len(c.lines)
	if n == 0 || c.lines[n-1] != line {
		c.lines = append(c.lines, line)
		c.offsets = append(c.offsets, len(c.code))
	}
	c.code = append(c.code, b)
}

func (c *Chunk) writeConstant(v Value, line int) {
	constant := c.addConstant(v)
	if len(c.constants) <= 0xFF {
		c.write(OpConstant, line)
		c.write(byte(constant), line)
		return
	}
	c.write(OpConstantX, line)
	c.write(byte(constant>>0), line)
	c.write(byte(constant>>8), line)
	c.write(byte(constant>>16), line)
}

func main() {
	c := NewChunk()
	c.writeConstant(1.2, 123)
	c.writeConstant(3.4, 124)
	c.writeConstant(5.6, 124)
	c.write(OpReturn, 125)
	c.disassemble("test chunk")
}
