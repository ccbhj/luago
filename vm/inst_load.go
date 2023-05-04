package vm

import . "luago/api"

func loadNil(i Instruction, vm LuaVM) {
	offset, n, _ := i.ABC()
	offset++

	vm.PushNil() // prepare a nil value into stack
	for i := offset; i <= offset+n; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

func loadBool(i Instruction, vm LuaVM) {
	idx, val, shouldIncrPC := i.ABC()
	idx++

	vm.PushBoolean(val != 0)
	vm.Replace(idx)
	if shouldIncrPC != 0 {
		vm.AddPC(1)
	}
}

func loadk(i Instruction, vm LuaVM) {
	rIdx, kIdx := i.ABx()
	rIdx++

	vm.GetConst(kIdx)
	vm.Replace(rIdx)
}

func loadKx(i Instruction, vm LuaVM) {
	rIdx, _ := i.ABx()
	rIdx++
	ax := Instruction(vm.Fetch()).Ax()

	vm.GetConst(ax)
	vm.Replace(rIdx)
}
