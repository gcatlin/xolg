package main

import "fmt"

func printValue(value Value) {
	fmt.Printf("%g", value)
}

func constantInstruction(name string, chunk *Chunk, offset int) {
	byte0 := int(chunk.code[offset+1])
	constant := byte0
	fmt.Printf("%-16s %4d '", name, constant)
	printValue(chunk.constants[constant])
	fmt.Printf("'")
}

func constantXInstruction(name string, chunk *Chunk, offset int) {
	byte0 := int(chunk.code[offset+1])
	byte1 := int(chunk.code[offset+2])
	byte2 := int(chunk.code[offset+3])
	constant := byte0<<0 | byte1<<8 | byte2<<16
	fmt.Printf("%-16s %4d '", name, constant)
	printValue(chunk.constants[constant])
	fmt.Printf("'")
}

func simpleInstruction(name string, offset int) {
	fmt.Printf("%-16s      ", name)
}

func unknownInstruction(instr Op, offset int) {
	fmt.Printf("Unknown opcode: %d", instr)
}

func (c *Chunk) disassemble(name string) {
	fmt.Printf("=== %s ===\n", name)
	fmt.Printf("OFFSET B0 B1 B2 B3 LINE   OPCODE\n")
	fmt.Printf("------ -- -- -- -- -----  ----------------\n")
	for i, max := 0, len(c.code); i < max; {
		i = c.disassembleInstruction(i)
		fmt.Println()
	}
	fmt.Println()
}

func (c *Chunk) disassembleInstruction(offset int) int {
	line := c.getLine(offset)
	instr := Op(c.code[offset])
	size := InstrSize[instr]
	if size == 0 {
		size = 1
	}

	// Instruction bytes
	fmt.Printf("%06X ", offset)
	for i := 0; i < size; i++ {
		fmt.Printf("%02X ", c.code[offset+i])
	}
	for i := size; i < 4; i++ {
		fmt.Printf("   ")
	}

	// Line numbers
	if offset > 0 && line == c.getLine(offset-1) {
		fmt.Printf("    |  ")
	} else {
		fmt.Printf("%5d  ", line)
	}

	switch instr {
	case OpConstant:
		constantInstruction("OpConstant", c, offset)
	case OpConstantX:
		constantXInstruction("OpConstantX", c, offset)
	case OpAdd:
		simpleInstruction("OpAdd", offset)
	case OpSub:
		simpleInstruction("OpSub", offset)
	case OpMul:
		simpleInstruction("OpMul", offset)
	case OpDiv:
		simpleInstruction("OpDiv", offset)
	case OpNegate:
		simpleInstruction("OpNegate", offset)
	case OpReturn:
		simpleInstruction("OpReturn", offset)
	default:
		unknownInstruction(instr, offset)
	}

	return offset + size
}
