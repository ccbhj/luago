package vm

import . "luago/api"

func getUpval(i Instruction, vm LuaVM) {
	// R(A) := UpValue[B]
	dst, idx, _ := i.ABC()
	dst++
	idx++
	vm.Copy(LuaUpvalueIndex(idx), dst)
}

func setUpval(i Instruction, vm LuaVM) {
	// UpValue[B] := R(A)
	from, idx, _ := i.ABC()
	from++
	idx++

	vm.Copy(from, LuaUpvalueIndex(idx))
}

func getTabUp(i Instruction, vm LuaVM) {
	// R[A] := UpValue[B][RK[C]]
	dst, upvalIdx, keyIdx := i.ABC()
	dst++
	upvalIdx++

	vm.GetRK(keyIdx)
	vm.GetTable(LuaUpvalueIndex(upvalIdx))
	vm.Replace(dst)
}

func setTabUp(i Instruction, vm LuaVM) {
	// UpValue[A][RK[B]] = RK[C]
	tabIdx, keyIdx, valIdx := i.ABC()
	tabIdx++
	valIdx++

	vm.GetRK(keyIdx)
	vm.GetRK(valIdx)
	vm.SetTable(LuaUpvalueIndex(tabIdx))
}
