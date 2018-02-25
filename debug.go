package main

import "fmt"

func printValue(value Value) {
	fmt.Printf("%g", value)
}

func constantInstruction(name string, chunk *Chunk, offset int) int {
	byte0 := int(chunk.code[offset+1])
	constant := byte0
	fmt.Printf("%-16s %4d '", name, constant)
	printValue(chunk.constants[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func constantXInstruction(name string, chunk *Chunk, offset int) int {
	byte0 := int(chunk.code[offset+1])
	byte1 := int(chunk.code[offset+2])
	byte2 := int(chunk.code[offset+3])
	constant := byte0<<0 | byte1<<8 | byte2<<16
	fmt.Printf("%-16s %4d '", name, constant)
	printValue(chunk.constants[constant])
	fmt.Printf("'\n")
	return offset + 4
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%-16s\n", name)
	return offset + 1
}

func unknownInstruction(instr byte, offset int) int {
	fmt.Printf("Unknown opcode: %d\n", instr)
	return offset + 1
}

func (c *Chunk) disassemble(name string) {
	fmt.Printf("=== %s ===\n", name)
	fmt.Printf("OFFSET B0 B1 B2 B3 LINE   OPCODE\n")
	fmt.Printf("------ -- -- -- -- -----  ----------------\n")
	for i, max := 0, len(c.code); i < max; {
		i = c.disassembleInstruction(i)
	}
}

func (c *Chunk) disassembleInstruction(offset int) int {
	line := c.getLine(offset)
	instr := c.code[offset]
	size := InstrSize[instr]

	// Instruction bytes
	fmt.Printf("%06d ", offset)
	for i := 0; i < size; i++ {
		fmt.Printf(HEX+" ", c.code[offset+i])
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
		constantInstruction("OP_CONSTANT", c, offset)
	case OpConstantX:
		constantInstruction("OP_CONSTANT_X", c, offset)
	case OpReturn:
		simpleInstruction("OP_RETURN", offset)
	default:
		unknownInstruction(instr, offset)
	}
	return offset + size
}
