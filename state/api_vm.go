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
	stack := l.stack
	subProto := stack.closure.proto.Protos[idx]
	newClosure := newLuaClosure(subProto)
	stack.push(newClosure)

	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack != 1 {
			// the upvalue is not from the outer function(maybe the outer's outer)
			// function g()
			//   local a, b
			//   function f()
			//       local c
			//       function h() a = b c = b end  -- a and b's Instack == 0, c's Instack = 1
			//   end
			// end
			newClosure.upvals[i] = stack.closure.upvals[uvIdx]
			continue
		}
		if stack.openuvs == nil {
			stack.openuvs = make(map[int]*upvalue, len(subProto.Upvalues))
		}
		if openuv, found := stack.openuvs[uvIdx]; found {
			newClosure.upvals[i] = openuv
			continue
		}
		// new closure
		newClosure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
		stack.openuvs[uvIdx] = newClosure.upvals[i]
	}
}
