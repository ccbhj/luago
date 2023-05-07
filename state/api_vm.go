package state

func (l *luaState) PC() int {
	return l.stack.pc
}

func (l *luaState) AddPC(n int) {
	l.stack.pc += n
}

func (l *luaState) Fetch() uint32 {
	i := l.stack.closure.proto.Code[l.stack.pc]
	l.stack.pc++
	return i
}

func (l *luaState) GetConst(idx int) {
	cnst := l.stack.closure.proto.Constants[idx]
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

func (l *luaState) RegisterCount() int {
	return int(l.stack.closure.proto.MaxStackSize)
}

func (l *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(l.stack.varargs)
	}
	l.stack.check(n)
	l.stack.pushN(l.stack.varargs, n)
}

func (l *luaState) LoadProto(idx int) {
	proto := l.stack.closure.proto.Protos[idx]
	closure := newLuaClosure(proto)
	l.stack.push(closure)
}
