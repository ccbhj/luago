package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"luago/api"
	"luago/state"
)

func _iPairAux(ls api.LuaState) int {
	i := ls.ToInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) == api.LUA_TNIL {
		return 1
	}
	return 2
}

func ipairs(ls api.LuaState) int {
	ls.PushGoFunction(_iPairAux, 0)
	ls.PushValue(1)
	ls.PushNil()
	return 3
}

func pairs(ls api.LuaState) int {
	ls.PushGoFunction(next, 0)
	ls.PushValue(1)
	ls.PushNil()
	return 3
}

func next(ls api.LuaState) int {
	ls.SetTop(2)
	if ls.Next(1) {
		return 2
	} else {
		ls.PushNil()
		return 1
	}
}

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

func getMetatable(ls api.LuaState) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func setMetatable(ls api.LuaState) int {
	ls.SetMetatable(1)
	return 1
}

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		ls := state.New()
		ls.Register("print", print)
		ls.Register("getmetatable", getMetatable)
		ls.Register("setmetatable", setMetatable)
		ls.Register("next", next)
		ls.Register("pairs", pairs)
		ls.Register("ipairs", ipairs)
		ls.Load(data, os.Args[1], "b")
		ls.Call(0, 0)
	}
}
