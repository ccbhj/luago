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
	return l.getTable(t, k, false)
}

func (l *luaState) getTable(t, k luaValue, ignoreMeta bool) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		if ignoreMeta || v != nil || !tbl.hasMetafield("__index") {
			l.stack.push(v)
			return typeOf(v)
		}
	}

	if !ignoreMeta {
		mf := getMetafield(t, "__index", l)
		if mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return l.getTable(x, k, false)
			case *closure:
				l.stack.pushN([]luaValue{mf, t, k}, 3)
				l.Call(2, 1)
				v := l.stack.get(-1)
				return typeOf(v)
			}
		}
	}
	panic("not a table!")
}

func (l *luaState) RawGet(idx int) LuaType {
	t := l.stack.get(idx)
	k := l.stack.pop()
	return l.getTable(t, k, true)
}

func (l *luaState) RawGetI(idx int, i int64) LuaType {
	tbl := l.stack.get(idx)
	return l.getTable(tbl, i, true)
}

func (l *luaState) GetField(idx int, k string) LuaType {
	t := l.stack.get(idx)
	return l.getTable(t, k, false)
}

func (l *luaState) GetI(idx int, i int64) LuaType {
	tbl := l.stack.get(idx)
	return l.getTable(tbl, i, false)
}

func (l *luaState) GetGlobal(name string) LuaType {
	t := l.registry.get(LUA_RIDX_GLOBALS)
	return l.getTable(t, name, false)
}

func (l *luaState) GetMetatable(idx int) bool {
	val := l.stack.get(idx)

	if mt := getMetatable(val, l); mt != nil {
		l.stack.push(mt)
		return true
	}
	return false
}
