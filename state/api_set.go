package state

func (l *luaState) SetTable(idx int) {
	tbl := l.stack.get(idx)
	v := l.stack.pop()
	k := l.stack.pop()

	l.setTable(tbl, k, v)
}

func (l *luaState) setTable(tbl, k, v luaValue) {
	if t, ok := tbl.(*luaTable); ok {
		t.put(k, v)
		return
	}
	panic("not a table!")
}

func (l *luaState) SetField(idx int, k string) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	l.setTable(t, k, v)
}

func (l *luaState) SetI(idx int, i int64) {
	t := l.stack.get(idx)
	v := l.stack.pop()

	l.setTable(t, i, v)
}
