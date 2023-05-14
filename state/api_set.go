package state

import "luago/api"

func (l *luaState) SetTable(idx int) {
	tbl := l.stack.get(idx)
	v := l.stack.pop()
	k := l.stack.pop()

	l.setTable(tbl, k, v, false)
}

func (l *luaState) RawSet(idx int) {
	tbl := l.stack.get(idx)
	v := l.stack.pop()
	k := l.stack.pop()

	l.setTable(tbl, k, v, true)
}

func (l *luaState) RawSetI(idx int, i int64) {
	t := l.stack.get(idx)
	v := l.stack.pop()

	l.setTable(t, i, v, true)
}

func (l *luaState) setTable(tbl, k, v luaValue, raw bool) {
	if t, ok := tbl.(*luaTable); ok {
		if raw || t.get(k) != nil || !t.hasMetafield("__newindex") {
			t.put(k, v)
			return
		}
		t.put(k, v)
		return
	}
	if !raw {
		mf := getMetafield(tbl, "__newindex", l)
		if mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				l.setTable(x, k, v, false)
				return
			case *closure:
				l.stack.pushN([]luaValue{mf, tbl, k, v}, 4)
				l.Call(3, 0)
				return
			}
		}
	}
	panic("not a table!")
}

func (l *luaState) SetField(idx int, k string) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	l.setTable(t, k, v, false)
}

func (l *luaState) SetI(idx int, i int64) {
	t := l.stack.get(idx)
	v := l.stack.pop()

	l.setTable(t, i, v, false)
}

func (l *luaState) SetGlobal(name string) {
	t := l.registry.get(api.LUA_RIDX_GLOBALS)
	v := l.stack.pop()
	l.setTable(t, name, v, false)
}

func (l *luaState) Register(name string, f api.GoFunction) {
	l.PushGoFunction(f, 0)
	l.SetGlobal(name)
}

func (l *luaState) SetMetatable(idx int) {
	mtVal := l.stack.pop()
	val := l.stack.get(idx)

	if mtVal == nil {
		setMetatable(val, nil, l)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, l)
	} else {
		// TODO
		panic("table expected!")
	}
}
