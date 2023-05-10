package state

import (
	. "luago/api"
	"luago/binchunk"
)

type upvalue struct {
	val *luaValue
}

type closure struct {
	proto  *binchunk.Prototype
	goFunc GoFunction
	upvals []*upvalue
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	closure := &closure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		closure.upvals = make([]*upvalue, nUpvals)
	}
	return closure
}

func newGoClosure(fn GoFunction, nUpvalue int) *closure {
	c := &closure{goFunc: fn}
	if nUpvalue > 0 {
		c.upvals = make([]*upvalue, nUpvalue)
	}
	return c
}
