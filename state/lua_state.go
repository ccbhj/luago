package state

type luaState struct {
	stack *luaStack
}

func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
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
