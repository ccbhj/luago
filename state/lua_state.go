package state

import "luago/api"

type luaState struct {
	stack    *luaStack
	registry *luaTable
}

func New() *luaState {
	registry := newLuaTable(0, 0)
	registry.put(api.LUA_RIDX_GLOBALS, newLuaTable(0, 0))

	ls := &luaState{
		registry: registry,
	}
	ls.pushLuaStack(newLuaStack(api.LUA_MINSTACK))

	return ls
}

func (l *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = l.stack
	l.stack = stack
}

func (l *luaState) popLuaStack() {
	stack := l.stack
	l.stack = stack.prev
	stack.prev = nil
}
