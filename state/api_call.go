package state

import (
	"fmt"
	"luago/api"
	"luago/binchunk"
	"luago/llog"
	"luago/vm"
	"reflect"
	"runtime"
	"strings"
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
		if mf := getMetafield(val, "__call", l); mf != nil {
			if c, ok = mf.(*closure); ok {
				l.stack.push(val)
				l.Insert(-(nArgs + 2))
				nArgs++
			}
		}
	}
	if c.goFunc != nil {
		l.callGoClosure(nArgs, nResults, c)
		return
	}
	llog.Debug("call %s<%d,%d>, nArgs=%d, nResults=%d\n", c.proto.Source,
		c.proto.LineDefined, c.proto.LastLineDefined, nArgs, nResults)
	l.callLuaClosure(nArgs, nResults, c)
}

func (l *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	newStack := newLuaStack(nArgs+20, l)
	newStack.closure = c

	args := l.stack.popN(nArgs) // get the arguments from the old stack
	newStack.pushN(args, nArgs) // push the arguments to the new stack
	l.stack.pop()               // pop the closure

	l.pushLuaStack(newStack)
	l.PrintStack()
	llog.Debug("call go func nArgs=%d, nResults=%d", nArgs, nResults)
	r := c.goFunc(l)
	l.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(r)
		l.stack.check(len(results))
		l.stack.pushN(results, nResults)
	}
	l.PrintStack()
}

func (l *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	newStack := newLuaStack(nRegs+20, l)
	newStack.closure = c

	funcAndArgs := l.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs

	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	l.pushLuaStack(newStack)
	l.PrintStack()
	l.runLuaClosure()
	l.PrintStack()
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
		llog.Debug("[%02d] %s", l.PC()+1, inst.String())
		inst.Execute(l)
		l.PrintStack()
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

func (l *luaState) PrintStack() {
	top := l.GetTop()
	elems := make([]string, 0, 1<<3)
	for i := 1; i <= top; i++ {
		t := l.Type(i)
		switch t {
		case api.LUA_TBOOLEAN:
			elems = append(elems, fmt.Sprintf("[%t]", l.ToBoolean(i)))
		case api.LUA_TNUMBER:
			elems = append(elems, fmt.Sprintf("[%g]", l.ToNumber(i)))
		case api.LUA_TSTRING:
			elems = append(elems, fmt.Sprintf("[%q]", l.ToString(i)))
		case api.LUA_TFUNCTION:
			var s string
			c, ok := l.stack.get(i).(*closure)
			if ok {
				if c.goFunc != nil {
					fnName := runtime.FuncForPC(reflect.ValueOf(c.goFunc).Pointer()).Name()
					s = fmt.Sprintf("[%s]", fnName)
				} else {
					s = fmt.Sprintf("[<%d, %d>]", c.proto.LineDefined, c.proto.LastLineDefined)
				}
			}
			elems = append(elems, s)

		default:
			elems = append(elems, fmt.Sprintf("[%s]", l.TypeName(t)))
		}
	}

	var loc string
	if l.stack.closure != nil && l.stack.closure.proto != nil {
		loc = fmt.Sprintf("%d:%d",
			l.stack.closure.proto.LineDefined, l.stack.closure.proto.LastLineDefined)
	} else {
		loc = "x.x"
	}
	llog.Debug("<%s>stack[%d]: %s", loc, top, strings.Join(elems, ""))
}

func (l *luaState) printUpval() {
	sb := strings.Builder{}
	for _, v := range l.stack.closure.upvals {
		if v == nil {
			sb.Write([]byte("[nil]"))
			continue
		}
		sb.Write([]byte(fmt.Sprintf("[%v]", *v.val)))
	}
	llog.Debug("upvalues[%d]: %s\n", len(l.stack.closure.upvals), sb.String())
}
