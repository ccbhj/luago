package state

func (l *luaState) GetTop() int {
	return l.stack.top
}
