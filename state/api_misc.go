package state

func (l *luaState) RawLen(idx int) int {
	val := l.stack.get(idx)

	if str, ok := val.(string); ok {
		i := int64(len(str))
		l.stack.push(i)
		return int(i)
	} else if t, ok := val.(*luaTable); ok {
		i := int64(t.len())
		l.stack.push(i)
		return int(i)
	} else {
		panic("length error")
	}
}

func (l *luaState) Len(idx int) {
	val := l.stack.get(idx)

	if str, ok := val.(string); ok {
		l.stack.push(int64(len(str)))
	} else if res, ok := callMetamethod(val, val, "__len", l); ok {
		l.stack.push(res)
	} else if t, ok := val.(*luaTable); ok {
		l.stack.push(int64(t.len()))
	} else {
		panic("length error")
	}
}

func (l *luaState) Concat(n int) {
	if n == 0 {
		l.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if l.IsString(-1) && l.IsString(-2) {
				s2, s1 := l.ToString(-1), l.ToString(-2)
				l.stack.pop()
				l.stack.pop()
				l.stack.push(s1 + s2)
				continue
			}
			b := l.stack.pop()
			a := l.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", l); ok {
				l.stack.push(result)
				continue
			}
			panic("concatenation error!")
		}
	}
}
