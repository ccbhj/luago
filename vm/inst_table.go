package vm

import (
	. "luago/api"
)

const CFIELDS_PER_FLUSH = 50

func newTable(i Instruction, vm LuaVM) {
	dst, nArr, nRec := i.ABC()
	dst++

	vm.CreateTable(Fb2int(nArr), Fb2int(nRec))
	vm.Replace(dst)
}

func getTable(i Instruction, vm LuaVM) {
	dst, tbl, key := i.ABC()
	dst++
	tbl++

	vm.GetRK(key)
	vm.GetTable(tbl)
	vm.Replace(dst)
}

func setTable(i Instruction, vm LuaVM) {
	tbl, key, val := i.ABC()
	tbl++

	vm.GetRK(key)
	vm.GetRK(val)
	vm.SetTable(tbl)
}

func setList(i Instruction, vm LuaVM) {
	list, batchOff, nbatch := i.ABC()
	list++

	isBatchOffZero := batchOff == 0
	if isBatchOffZero {
		// returned values from function call
		batchOff = int(vm.ToInteger(-1)) - list - 1
		vm.Pop(1)
	}

	if nbatch > 0 {
		nbatch--
	} else {
		nbatch = Instruction(vm.Fetch()).Ax()
	}

	idx := int64(nbatch * CFIELDS_PER_FLUSH)
	for j := 1; j <= batchOff; j++ {
		idx++
		vm.PushValue(list + j)
		vm.SetI(list, idx)
	}

	if isBatchOffZero {
		for j := vm.RegisterCount() + 1; j <= vm.GetTop(); j++ {
			idx++
			vm.PushValue(j)
			vm.SetI(list, idx)
		}
		// clear the stack
		vm.SetTop(vm.RegisterCount())
	}
}
