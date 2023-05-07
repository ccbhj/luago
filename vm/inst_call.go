package vm

import . "luago/api"

func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a++

	vm.LoadProto(bx)
	vm.Replace(a)
}

func call(i Instruction, vm LuaVM) {
	fn, nArgs, nRet := i.ABC()
	fn++

	vm.Call(_pushFuncAndArgs(fn, nArgs, vm), nRet-1)
	_popResults(fn, nRet, vm)
}

func _return(i Instruction, vm LuaVM) {
	offset, n, _ := i.ABC()
	offset++

	if n == 1 {
		// no return value
		return
	}
	if n > 1 {
		// more than one return values
		vm.CheckStack(n - 1)
		for i := offset; i <= offset+n-2; i++ {
			vm.PushValue(i)
		}
		return
	}
	// collect return values from the last function call
	_fixStack(offset, vm)
}

func vararg(i Instruction, vm LuaVM) {
	off, n, _ := i.ABC()
	off++

	if n != 1 {
		vm.LoadVararg(n - 1)
		_popResults(off, n, vm)
	}
}

func tailcall(i Instruction, vm LuaVM) {
	fn, nArgs, _ := i.ABC()
	fn++
	nReturns := 0

	nArgs = _pushFuncAndArgs(fn, nArgs, vm)
	vm.Call(nArgs, nReturns-1)
	_popResults(fn, nReturns, vm)
}

func self(i Instruction, vm LuaVM) {
	// R(A+1) := R(B)
	// R(A) := R(B)[RK(C)]
	dst, objIdx, fnIdx := i.ABC()
	dst++
	objIdx++

	vm.Copy(objIdx, dst+1)
	vm.GetRK(fnIdx)
	vm.GetTable(objIdx)
	vm.Replace(dst)
}

/*
+--------+
| args_n | <- (`fn` + `nArgs` - 1)
+--------+
| ...... |
+--------+
| args_2 |            ========>   fn(args_1, args_2, ..., args_n)
+--------+
| args_1 |
+--------+
|function| <- `fn`
+--------+
*/
func _pushFuncAndArgs(fn, nArgs int, vm LuaVM) int {
	if nArgs >= 1 {
		// push all arguments to the top
		vm.CheckStack(nArgs)
		for i := fn; i < fn+nArgs; i++ {
			vm.PushValue(i)
		}
		return nArgs - 1
	}
	// collect all the return values from another function call
	// see _popResults:87
	_fixStack(fn, vm)
	return vm.GetTop() - vm.RegisterCount() - 1
}

func _fixStack(fn int, vm LuaVM) {
	x := int(vm.ToInteger(-1)) // see _popResults:87
	vm.Pop(1)

	vm.CheckStack(x - fn)
	for i := fn; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-fn)
}

func _popResults(fn, nRet int, vm LuaVM) {
	if nRet == 1 {
		return
	}
	if nRet > 1 {
		for i := fn + nRet - 2; i >= fn; i-- {
			vm.Replace(i)
		}
		return
	}
	// we need all of the values returned by fn
	// so we left them in the top of the stack to use them later
	// push a integer to mark them
	vm.CheckStack(1)
	vm.PushInteger(int64(fn))
}
