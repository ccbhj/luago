package state

import (
	. "luago/api"
	"luago/binchunk"
)

type closure struct {
	proto  *binchunk.Prototype
	goFunc GoFunction
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	return &closure{proto: proto}
}

func newGoClosure(fn GoFunction) *closure {
	return &closure{goFunc: fn}
}
