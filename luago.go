package main

import (
	"fmt"
	"io/ioutil"
	"luago/binchunk"
	"luago/state"
	"luago/vm"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		proto := binchunk.Undump(data)
		list(proto)
		luaMain(proto)
	}
}

func luaMain(proto *binchunk.Prototype) {
	nRegs := int(proto.MaxStackSize)
	ls := state.New(nRegs+8, proto)
	ls.SetTop(nRegs)
	for {
		pc := ls.PC()
		inst := vm.Instruction(ls.Fetch())
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
		inst.Execute(ls)
		fmt.Printf("[%02d] %-8s ", pc+1, inst.OpName())
		printStack(ls)
	}
}
