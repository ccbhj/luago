package api

type LuaVM interface {
	LuaState

	PC() int             // only for testing
	AddPC(n int)         // modify the PC(use for jumping)
	Fetch() uint32       // fetch current instruction and move PC to the next instruction
	GetConst(idx int)    // push const at idx to stack top
	GetRK(rk int)        // push const or stack value to stack top
	LoadProto(idx int)   // load and push proto to stack top
	RegisterCount() int  // return how many registers
	LoadVararg(n int)    // load n arguments
	CloseUpvalues(a int) // close upvalues
}
