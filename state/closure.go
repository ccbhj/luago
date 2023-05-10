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
	return &closure{proto: proto}
}

func newGoClosure(fn GoFunction, nUpvalue int) *closure {
	c := &closure{goFunc: fn}
	if nUpvalue > 0 {
		c.upvals = make([]*upvalue, nUpvalue)
	}
	return c
}
