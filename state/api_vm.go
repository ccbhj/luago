package state

import "luago/llog"

func (l *luaState) PC() int {
	return l.pc
}

func (l *luaState) AddPC(n int) {
	l.pc += n
}

func (l *luaState) Fetch() uint32 {
	if l.pc < 0 || l.pc >= len(l.proto.Code) {
		llog.Fatal("pc(%d) out of range code table[0, %d)", l.pc, len(l.proto.Code))
	}
	ins := l.proto.Code[l.pc]
	l.pc++
	return ins
}

func (l *luaState) GetConst(idx int) {
	if idx < 0 || idx >= len(l.proto.Constants) {
		llog.Fatal("idx(%d) out of range constants table[0, %d)", idx, len(l.proto.Constants))
	}
	cnst := l.proto.Constants[idx]
	l.stack.push(cnst)
}

// rk must be the OpArgK in an iABC instruction
func (l *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		l.GetConst(rk & 0xFF)
		return
	}
	// register
	l.PushValue(rk + 1)
}
