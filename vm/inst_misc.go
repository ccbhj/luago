package vm

import . "luago/api"

func move(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(b, a)
}

func jmp(i Instruction, vm LuaVM) {
	a, jmpTo := i.AsBx()
	vm.AddPC(jmpTo)
	if a != 0 {
		vm.CloseUpvalues(a)
	}
}
