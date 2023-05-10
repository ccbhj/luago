package state

import (
	"luago/api"
	"luago/binchunk"
	"luago/llog"
	"luago/utils"
	"luago/vm"
)

func (l *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	l.stack.push(c)
	if len(proto.Upvalues) > 0 {
		env := l.registry.get(api.LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	return 0
}

func (l *luaState) Call(nArgs, nResults int) {
	val := l.stack.get(-(nArgs + 1))
	c, ok := val.(*closure)
	if !ok {
		panic("not a function")
	}
	if c.goFunc != nil {
		llog.Debug("call go function\n")
		l.callGoClosure(nArgs, nResults, c)
		return
	}
	llog.Debug("call %s<%d,%d>\n", c.proto.Source,
		c.proto.LineDefined, c.proto.LastLineDefined)
	l.callLuaClosure(nArgs, nResults, c)
}

func (l *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	newStack := newLuaStack(nRegs + 20)
	newStack.closure = c

	funcAndArgs := l.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs

	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	l.pushLuaStack(newStack)
	l.runLuaClosure()
	l.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		l.stack.check(len(results))
		l.stack.pushN(results, nResults)
	}
}

func (l *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(l.Fetch())
		llog.Debug("[%02d] %-8s ", l.PC()+1, inst.OpName())
		inst.Execute(l)
		utils.PrintStack(l)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

func (l *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	newStack := newLuaStack(nArgs + 20)
	newStack.closure = c

	funcAndArgs := l.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nArgs)
	l.stack.pop()

	l.pushLuaStack(newStack)
	r := c.goFunc(l)
	l.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(r)
		l.stack.check(len(results))
		l.stack.pushN(results, nResults)
	}
}
