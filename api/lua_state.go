package api

type GoFunction func(LuaState) int

type LuaState interface {
	// basic stack manipulation
	GetTop() int
	AbsIndex(idx int) int
	// Check grows the stack's size to be at least `n` free slots
	CheckStack(n int) bool
	Pop(n int)
	Copy(fromIdx, toIdx int)
	// PushValue get the value at idx and push them to the top
	PushValue(idx int)
	// Replace pop the value on stack top and replace the value at `idx` with it
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx, n int)
	SetTop(idx int)

	// access function
	TypeName(t LuaType) string
	Type(idx int) LuaType
	IsNone(idx int) bool
	IsNil(idx int) bool
	IsNoneOrNil(idx int) bool
	IsBoolean(idx int) bool
	IsInteger(idx int) bool
	IsNumber(idx int) bool
	IsString(idx int) bool

	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64, bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)

	// push function (Go -> stack)
	PushNil()
	PushBoolean(bool)
	PushInteger(int64)
	PushNumber(float64)
	PushString(string)

	// Arith evaluate the op using values in stack
	// and push the result into the stack top
	// for math and bitwise op
	Arith(op ArithOp)
	Compare(idx1, idx2 int, op CompareOp) bool
	Len(idx int)
	// Concat perform `n` times concatention at stack top
	Concat(n int)

	// Table APIs
	NewTable()
	CreateTable(nArr, nRec int)
	GetTable(idx int) LuaType
	GetField(idx int, k string) LuaType
	GetI(idx int, i int64) LuaType
	// Table set function
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, n int64)

	Load(chunk []byte, chunkName, mode string) int
	Call(nArgs, nResults int)

	PushGoFunction(f GoFunction)
	IsGoFunction(idx int) bool
	ToGoFunction(idx int) GoFunction

	PushGlobalTable()
	GetGlobal(name string) LuaType
	SetGlobal(name string)
	Register(name string, f GoFunction)
}
