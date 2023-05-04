package state

import "luago/binchunk"

type luaState struct {
	stack *luaStack
	proto *binchunk.Prototype
	pc    int
}

func New(stackSize int, pb *binchunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),
		proto: pb,
		pc:    0,
	}
}
