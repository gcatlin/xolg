package main

import "fmt"

type Op byte

const (
	OpConstant Op = iota
	OpConstantX
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpNegate
	OpReturn
	__op_count__
)

type InterpretResult int

const (
	InterpretOk InterpretResult = iota
	InterpretCompileError
	InterpretRuntimeError
)

var InstrSize = [__op_count__]int{
	OpConstant:  2,
	OpConstantX: 4,
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
		c.writeOp(OpConstant, line)
		c.write(byte(constant), line)
		return
	}
	c.writeOp(OpConstantX, line)
	c.write(byte(constant>>0), line)
	c.write(byte(constant>>8), line)
	c.write(byte(constant>>16), line)
}

func (c *Chunk) writeOp(op Op, line int) {
	c.write(byte(op), line)
}

type VM struct {
	chunk *Chunk
	ip    int // offset into Chunk
	stack []Value
}

func NewVM() *VM {
	return &VM{
		stack: make([]Value, 0, 256),
	}
}

func (vm *VM) resetStack() {
	vm.stack = vm.stack[:0]
}

func (vm *VM) interpret(c *Chunk) InterpretResult {
	vm.chunk = c
	vm.ip = 0
	return vm.run()
}

func (vm *VM) next() Op {
	instr := Op(vm.chunk.code[vm.ip])
	vm.ip++
	return instr
}

func (vm *VM) push(v Value) {
	vm.stack = append(vm.stack, v)
}

func (vm *VM) pop() Value {
	v := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return v
}

func (vm *VM) readConstant() Value {
	return vm.chunk.constants[vm.next()]
}

func (vm *VM) readConstantX() Value {
	byte0 := int(vm.next())
	byte1 := int(vm.next())
	byte2 := int(vm.next())
	constant := byte0<<0 | byte1<<8 | byte2<<16
	return vm.chunk.constants[constant]
}

func (vm *VM) run() InterpretResult {
	for {
		fmt.Printf("          ")
		for _, v := range vm.stack {
			fmt.Printf("[ ")
			printValue(v)
			fmt.Printf(" ]")
		}
		fmt.Printf("\n")
		vm.chunk.disassembleInstruction(vm.ip)
		instr := vm.next()
		switch instr {
		case OpConstant:
			constant := vm.readConstant()
			vm.push(constant)
			printValue(constant)
			fmt.Printf("\n")
		case OpConstantX:
			constant := vm.readConstantX()
			printValue(constant)
			fmt.Printf("\n")
		case OpAdd:
			vm.push(vm.pop() + vm.pop())
		case OpSub:
			vm.push(-vm.pop() + vm.pop())
		case OpMul:
			vm.push(vm.pop() + vm.pop())
		case OpDiv:
			rhs := vm.pop()
			lhs := vm.pop()
			vm.push(lhs / rhs)
		case OpNegate:
			vm.push(-vm.pop())
		case OpReturn:
			printValue(vm.pop())
			fmt.Printf("\n")
			return InterpretOk
		}
	}
}

func main() {
	vm := NewVM()
	c := NewChunk()
	c.writeConstant(1.2, 123)
	c.writeConstant(3.4, 123)
	c.writeOp(OpAdd, 123)
	c.writeConstant(5.6, 123)
	c.writeOp(OpDiv, 123)
	c.writeOp(OpNegate, 123)
	c.writeOp(OpReturn, 123)
	c.disassemble("test chunk")
	vm.interpret(c)
}
