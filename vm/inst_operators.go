package vm

import . "luago/api"

func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, c := i.ABC()
	a++

	vm.GetRK(b)
	vm.GetRK(c)

	vm.Arith(op)
	vm.Replace(a)
}

func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := i.ABC()

	a++
	b++

	vm.PushValue(b)
	vm.Arith(op)
	vm.Replace(a)
}

func newUnOpInstFn(op ArithOp) func(Instruction, LuaVM) {
	return func(i Instruction, vm LuaVM) {
		_unaryArith(i, vm, op)
	}
}

func newBiOpInstFn(op ArithOp) func(Instruction, LuaVM) {
	return func(i Instruction, vm LuaVM) {
		_binaryArith(i, vm, op)
	}
}

var (
	add  = newBiOpInstFn(LUA_OPADD)
	sub  = newBiOpInstFn(LUA_OPSUB)
	mul  = newBiOpInstFn(LUA_OPMUL)
	mod  = newBiOpInstFn(LUA_OPMOD)
	pow  = newBiOpInstFn(LUA_OPPOW)
	div  = newBiOpInstFn(LUA_OPDIV)
	idiv = newBiOpInstFn(LUA_OPIDIV)
	band = newBiOpInstFn(LUA_OPBAND)
	bor  = newBiOpInstFn(LUA_OPBOR)
	bxor = newBiOpInstFn(LUA_OPBXOR)
	shl  = newBiOpInstFn(LUA_OPSHL)
	shr  = newBiOpInstFn(LUA_OPSHR)
	unm  = newUnOpInstFn(LUA_OPUNM)
	bnot = newUnOpInstFn(LUA_OPBNOT)
)

func _len(i Instruction, vm LuaVM) {
	dst, src, _ := i.ABC()
	dst++
	src++

	vm.Len(src)
	vm.Replace(dst)
}

func concat(i Instruction, vm LuaVM) {
	dst, start, end := i.ABC()
	dst++
	start++
	end++

	// concat n values between [`start`,`end`] at stack
	n := end - start + 1
	vm.CheckStack(n)
	for i := start; i <= end; i++ {
		vm.PushValue(i)
	}
	vm.Concat(n)
	vm.Replace(dst)
}

func _compare(i Instruction, vm LuaVM, op CompareOp) {
	// if (RK(a) op RK(b) ~= expect) then pc++
	expect, a, b := i.ABC()

	vm.GetRK(a)
	vm.GetRK(b)
	if vm.Compare(-2, -1, op) != (expect != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}

func newCmpInstFn(op CompareOp) func(Instruction, LuaVM) {
	return func(i Instruction, vm LuaVM) {
		_compare(i, vm, op)
	}
}

var (
	eq = newCmpInstFn(LUA_OPEQ)
	lt = newCmpInstFn(LUA_OPLT)
	le = newCmpInstFn(LUA_OPLE)
)

func not(i Instruction, vm LuaVM) {
	dst, src, _ := i.ABC()
	dst++
	src++

	vm.PushBoolean(!vm.ToBoolean(src))
	vm.Replace(dst)
}

func testSet(i Instruction, vm LuaVM) {
	// if (R(B) == bool(C)) then R(A) := R(B) else pc++
	dst, src, val := i.ABC()
	dst++
	src++

	if vm.ToBoolean(src) == (val != 0) {
		vm.Copy(src, dst)
	} else {
		vm.AddPC(1)
	}
}

func test(i Instruction, vm LuaVM) {
	// if not (R(A) == bool(C)) then pc++
	a, _, c := i.ABC()
	a++
	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}

func forPrep(i Instruction, vm LuaVM) {
	localVar, nextPCOff := i.AsBx()
	localVar++

	// R(A) -= R(A+2)
	vm.PushValue(localVar)
	vm.PushValue(localVar + 2)
	vm.Arith(LUA_OPSUB)
	vm.Replace(localVar)

	// pc += sbx
	vm.AddPC(nextPCOff)
}

func forLoop(i Instruction, vm LuaVM) {
	localVar, nextPCOff := i.AsBx()
	localVar++

	// R(A) += R(A+2)
	vm.PushValue(localVar + 2)
	vm.PushValue(localVar)
	vm.Arith(LUA_OPADD)
	vm.Replace(localVar)

	// R(A) <?= R(A+1)
	isPostiveStep := vm.ToNumber(localVar+2) >= 0
	if (isPostiveStep && vm.Compare(localVar, localVar+1, LUA_OPLE)) ||
		(!isPostiveStep && vm.Compare(localVar+1, localVar, LUA_OPLE)) {
		vm.AddPC(nextPCOff)           // pc += sBx
		vm.Copy(localVar, localVar+3) // R(A+3) = R(A)
	}
}
