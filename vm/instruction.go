package vm

import (
	"fmt"
	"luago/api"
)

const MAXARG_Bx = 1<<18 - 1       // 2^18 - 1 = 262143
const MAXARG_sBx = MAXARG_Bx >> 1 //  262143 / 2 = 131071

type Instruction uint32

func (i Instruction) Opcode() int {
	return int(i & 0x3F)
}

// iABC  +--B:9--+--C:9--+--A:8--+--opcode:6--+
func (i Instruction) ABC() (a, b, c int) {
	a = int(i >> 6 & 0xFF)
	c = int(i >> 14 & 0x1FF)
	b = int(i >> 23 & 0x1FF)
	return
}

// iABx  +---Bx:18-------+--A:8--+--opcode:6--+
func (i Instruction) ABx() (a, bx int) {
	a = int(i >> 6 & 0xFF)
	bx = int(i >> 14)
	return
}

// iAsBx +--sBx:18-------+--A:8--+--opcode:6--+
func (i Instruction) AsBx() (a, sbx int) {
	a, bx := i.ABx()
	return a, bx - MAXARG_sBx
}

func (i Instruction) Ax() int {
	return int(i >> 6)
}

func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

func (i Instruction) OpMode() byte {
	return byte(opcodes[i.Opcode()].opMode)
}

func (i Instruction) BMode() byte {
	return byte(opcodes[i.Opcode()].argBMode)
}

func (i Instruction) CMode() byte {
	return byte(opcodes[i.Opcode()].argCMode)
}

func (i Instruction) Execute(vm api.LuaVM) {
	act := opcodes[i.Opcode()].action
	if act != nil {
		act(i, vm)
	} else {
		panic(fmt.Sprintf("op %d action not found", i))
	}
}
