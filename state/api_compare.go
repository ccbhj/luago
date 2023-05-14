package state

import (
	. "luago/api"
	"luago/llog"
)

func (l *luaState) RawEqual(idx1, idx2 int) bool {
	a, b := l.stack.get(idx1), l.stack.get(idx2)
	return _eq(a, b, l, true)
}

func (l *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a, b := l.stack.get(idx1), l.stack.get(idx2)
	switch op {
	case LUA_OPEQ:
		return _eq(a, b, l, false)
	case LUA_OPLT:
		return _lt(a, b, l)
	case LUA_OPLE:
		return _le(a, b, l)
	default:
		panic("invalid compare op")
	}
}

func _eq(a, b luaValue, ls *luaState, raw bool) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return float64(y) == x
		case float64:
			return x == y
		default:
			return false
		}
	case *luaTable:
		if raw {
			return a == b
		}
		if y, ok := b.(*luaTable); ok && x != y && ls != nil {
			if res, ok := callMetamethod(x, y, "__eq", ls); ok {
				return convertToBoolean(res)
			}
			return a == b
		}
	}
	return a == b
}

func _lt(a, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x < float64(y)
		case float64:
			return x < y
		}
	}
	if res, ok := callMetamethod(a, b, "__lt", ls); ok {
		return convertToBoolean(res)
	}
	llog.Fatal("comparison error!")
	return false
}

func _le(a, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x <= float64(y)
		case float64:
			return x <= y
		}
	}
	if res, ok := callMetamethod(a, b, "__le", ls); ok {
		return convertToBoolean(res)
	} else if res, ok := callMetamethod(b, a, "__lt", ls); ok {
		return !convertToBoolean(res)
	}
	panic("comparison error!")
}
