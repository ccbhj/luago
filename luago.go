package main

import (
	"io/ioutil"
	"luago/state"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		ls := state.New()
		ls.Load(data, os.Args[1], "b")
		ls.Call(0, 0)
	}
}

// func luaMain(proto *binchunk.Prototype) {
// 	nRegs := int(proto.MaxStackSize)
// 	ls := state.New(nRegs+8, proto)
// 	ls.SetTop(nRegs)
// 	for {
// 		pc := ls.PC()
// 		inst := vm.Instruction(ls.Fetch())
// 		if inst.Opcode() == vm.OP_RETURN {
// 			break
// 		}
// 		inst.Execute(ls)
// 		fmt.Printf("[%02d] %-8s ", pc+1, inst.OpName())
// 		printStack(ls)
// 	}
// }
