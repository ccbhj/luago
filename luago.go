package main

import (
	"fmt"
	"io/ioutil"
	"luago/api"
	"luago/state"
	"os"
)

func print(ls api.LuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print(ls.ToString(i))
		} else {
			fmt.Printf(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Printf("\t")
		}
	}
	fmt.Println()
	return 0
}

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		ls := state.New()
		ls.Register("print", print)
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
