package state

import "luago/api"

type luaStack struct {
	slots []luaValue
	top   int

	prev    *luaStack
	closure *closure
	varargs []luaValue
	pc      int
	state   *luaState
}

func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		state: state,
		top:   0,
	}
}

// check grows the stack's size to be at least `n` free slots
func (l *luaStack) check(n int) {
	free := len(l.slots) - l.top
	for i := free; i < n; i++ {
		l.slots = append(l.slots, nil)
	}
}

func (l *luaStack) push(val luaValue) {
	if l.top == len(l.slots) {
		panic("stack overflow!")
	}

	l.slots[l.top] = val
	l.top++
}

func (l *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}
	for i := 0; i < n; i++ {
		if i < nVals {
			l.push(vals[i])
			continue
		}
		l.push(nil)
	}
}

func (l *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = l.pop()
	}
	return vals
}

func (l *luaStack) pop() luaValue {
	if l.top < 1 {
		panic("stack underflow!")
	}
	l.top--
	val := l.slots[l.top]
	l.slots[l.top] = nil
	return val
}

func (l *luaStack) absIndex(idx int) int {
	if idx <= api.LUA_REGISTRYINDEX {
		return idx
	}
	if idx >= 0 {
		return idx
	}
	return idx + l.top + 1
}

func (l *luaStack) isValid(idx int) bool {
	if idx == api.LUA_REGISTRYINDEX {
		return true
	}
	idx = l.absIndex(idx)
	return idx > 0 && idx <= l.top
}

func (l *luaStack) isAbsIdxValid(idx int) bool {
	return idx > 0 && idx <= l.top
}

func (l *luaStack) get(idx int) luaValue {
	if idx == api.LUA_REGISTRYINDEX {
		return l.state.registry
	}
	idx = l.absIndex(idx)
	if l.isAbsIdxValid(idx) {
		return l.slots[idx-1]
	}
	return nil
}

func (l *luaStack) set(idx int, val luaValue) {
	if idx == api.LUA_REGISTRYINDEX {
		l.state.registry = val.(*luaTable)
		return
	}
	idx = l.absIndex(idx)
	if l.isAbsIdxValid(idx) {
		l.slots[idx-1] = val
		return
	}
	panic("invalid index!")
}
