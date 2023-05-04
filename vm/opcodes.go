package vm

import "luago/api"

// encoding modes(format)
type OpEncodingMode byte

/*
	31    25|24   16|15    8|7           0

iABC  +--B:9--+--C:9--+--A:8--+--opcode:6--+
iABx  +---Bx:18-------+--A:8--+--opcode:6--+
iAsBx +--sBx:18-------+--A:8--+--opcode:6--+
iAx   +--------Ax:26----------+--opcode:6--+
*/
const (
	IABC = iota
	IABx
	IAsBx
	IAx
)

// op codes
type OpCode byte

const (
	OP_MOVE = iota
	OP_LOADK
	OP_LOADKX
	OP_LOADBOOL
	OP_LOADNIL
	OP_GETUPVAL
	OP_GETTABUP
	OP_GETTABLE
	OP_SETTABUP
	OP_SETUPVAL
	OP_SETTABLE
	OP_NEWTABLE

	OP_SELF
	OP_ADD
	OP_SUB
	OP_MUL
	OP_MOD
	OP_POW
	OP_DIV
	OP_IDIV
	OP_BAND
	OP_BOR
	OP_BXOR
	OP_SHL
	OP_SHR
	OP_UNM
	OP_BNOT
	OP_NOT
	OP_LEN
	OP_CONCAT
	OP_JMP
	OP_EQ
	OP_LT
	OP_LE
	OP_TEST
	OP_TESTSET
	OP_CALL
	OP_TAILCALL
	OP_RETURN
	OP_FORLOOP
	OP_FORPREP
	OP_TFORCALL
	OP_TFORLOOP
	OP_SETLIST
	OP_CLOSURE
	OP_VARARG
	OP_EXTRAARG
)

// argument types
type OpArgMode byte

const (
	OpArgN OpArgMode = iota // argument is not used
	OpArgU                  // argument is used
	OpArgR                  // argument is a register(iABC) or a jump offset(iAsBx)
	OpArgK                  // argument is a index of constant or register
)

type opcode struct {
	testFlag byte           // operator is a test (next instruction must be a jump)
	setAFlag byte           // instruction set register A
	argBMode OpArgMode      // B arg mode
	argCMode OpArgMode      // C arg mode
	opMode   OpEncodingMode // op model
	name     string

	action func(Instruction, api.LuaVM)
}

var opcodes = []opcode{
	{0, 1, OpArgR, OpArgN, IABC, "MOVE", move},
	{0, 1, OpArgK, OpArgN, IABx, "LOADK", loadk},
	{0, 1, OpArgN, OpArgN, IABx, "LOADKX", loadKx},
	{0, 1, OpArgU, OpArgU, IABC, "LOADBOOL", loadBool},
	{0, 1, OpArgU, OpArgN, IABC, "LOADNIL", loadNil},

	{0, 1, OpArgU, OpArgN, IABC, "GETUPVAL", nil},
	{0, 1, OpArgU, OpArgK, IABC, "GETTABUP", nil},
	{0, 1, OpArgR, OpArgK, IABC, "GETTABLE", getTable},

	{0, 0, OpArgK, OpArgK, IABC, "SETTABUP", nil},
	{0, 0, OpArgU, OpArgN, IABC, "SETUPVAL", nil},
	{0, 0, OpArgK, OpArgK, IABC, "SETTABLE", setTable},
	{0, 1, OpArgU, OpArgU, IABC, "NEWTABLE", newTable},

	{0, 1, OpArgR, OpArgK, IABC, "SELF", nil},
	{0, 1, OpArgK, OpArgK, IABC, "ADD", add},
	{0, 1, OpArgK, OpArgK, IABC, "SUB", sub},
	{0, 1, OpArgK, OpArgK, IABC, "MUL", mul},
	{0, 1, OpArgK, OpArgK, IABC, "MOD", mod},
	{0, 1, OpArgK, OpArgK, IABC, "POW", pow},
	{0, 1, OpArgK, OpArgK, IABC, "DIV", div},
	{0, 1, OpArgK, OpArgK, IABC, "IDIV", idiv},
	{0, 1, OpArgK, OpArgK, IABC, "BAND", band},
	{0, 1, OpArgK, OpArgK, IABC, "BOR", bor},
	{0, 1, OpArgK, OpArgK, IABC, "BXOR", bxor},
	{0, 1, OpArgK, OpArgK, IABC, "SHL", shl},
	{0, 1, OpArgK, OpArgK, IABC, "SHR", shr},

	{0, 1, OpArgR, OpArgN, IABC, "UNM", unm},
	{0, 1, OpArgR, OpArgN, IABC, "BNOT", bnot},
	{0, 1, OpArgR, OpArgN, IABC, "NOT", not},
	{0, 1, OpArgR, OpArgN, IABC, "LEN", _len},
	{0, 1, OpArgR, OpArgR, IABC, "CONCAT", concat},

	{0, 0, OpArgR, OpArgN, IAsBx, "JMP", jmp},

	{1, 0, OpArgK, OpArgK, IABC, "EQ", eq},
	{1, 0, OpArgK, OpArgK, IABC, "LT", lt},
	{1, 0, OpArgK, OpArgK, IABC, "LE", le},
	{1, 0, OpArgN, OpArgU, IABC, "TEST", test},
	{1, 1, OpArgR, OpArgU, IABC, "TESTSET", testSet},

	{0, 1, OpArgU, OpArgU, IABC, "CALL", nil},
	{0, 1, OpArgU, OpArgU, IABC, "TESTCALL", nil},
	{0, 0, OpArgU, OpArgN, IABC, "RETURN", nil},

	{0, 1, OpArgR, OpArgN, IAsBx, "FORLOOP", forLoop},
	{0, 1, OpArgR, OpArgN, IAsBx, "FORPREP", forPrep},

	{0, 0, OpArgN, OpArgU, IABC, "TFORCALL", nil},
	{0, 1, OpArgR, OpArgN, IAsBx, "TFORLOOP", nil},

	{0, 0, OpArgU, OpArgU, IABC, "SETLIST", setList},
	{0, 1, OpArgU, OpArgN, IABx, "CLOSURE", nil},
	{0, 1, OpArgU, OpArgN, IABC, "VARARG", nil},
	{0, 0, OpArgU, OpArgU, IAx, "EXTRAARG", nil},
}
