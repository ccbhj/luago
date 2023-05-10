package state

import . "luago/api"

func (l *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	l.stack.push(t)
}

func (l *luaState) NewTable() {
	l.CreateTable(0, 0)
}

func (l *luaState) GetTable(idx int) LuaType {
	t := l.stack.get(idx)
	k := l.stack.pop()
	return l.getTable(t, k)
}

func (l *luaState) getTable(t, k luaValue) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		l.stack.push(v)
		return typeOf(v)
	}
	panic("not a table!")
}

func (l *luaState) GetField(idx int, k string) LuaType {
	t := l.stack.get(idx)
	return l.getTable(t, k)
}

func (l *luaState) GetI(idx int, i int64) LuaType {
	tbl := l.stack.get(idx)
	return l.getTable(tbl, i)
}

func (l *luaState) GetGlobal(name string) LuaType {
	t := l.registry.get(LUA_RIDX_GLOBALS)
	return l.getTable(t, name)
}
