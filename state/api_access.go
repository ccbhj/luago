package state

import (
	"fmt"
	. "luago/api"
)

func (l *luaState) TypeName(t LuaType) string {
	return t.String()
}

func (l *luaState) Type(idx int) LuaType {
	if l.stack.isValid(idx) {
		return typeOf(l.stack.get(idx))
	}
	return LUA_TNONE
}

func (l *luaState) IsNone(idx int) bool {
	return l.Type(idx) == LUA_TNONE
}

func (l *luaState) IsNil(idx int) bool {
	return l.Type(idx) == LUA_TNIL
}

func (l *luaState) IsNoneOrNil(idx int) bool {
	return l.Type(idx) <= LUA_TNIL
}

func (l *luaState) IsBoolean(idx int) bool {
	return l.Type(idx) == LUA_TBOOLEAN
}

func (l *luaState) IsString(idx int) bool {
	t := l.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

func (l *luaState) IsNumber(idx int) bool {
	_, ok := l.ToNumberX(idx)
	return ok
}

func (l *luaState) IsInteger(idx int) bool {
	val := l.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (l *luaState) ToBoolean(idx int) bool {
	val := l.stack.get(idx)
	return convertToBoolean(val)
}

func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func (l *luaState) ToNumber(idx int) float64 {
	f, _ := l.ToNumberX(idx)
	return f
}

func (l *luaState) ToNumberX(idx int) (float64, bool) {
	return convertToFloat(l.stack.get(idx))
}

func (l *luaState) ToInteger(idx int) int64 {
	i, _ := l.ToIntegerX(idx)
	return i
}

func (l *luaState) ToIntegerX(idx int) (int64, bool) {
	return convertToInteger(l.stack.get(idx))
}

// ToStringX returns the string at `idx` if there's a string
// or convert the integer or number into string and return it
func (l *luaState) ToStringX(idx int) (string, bool) {
	val := l.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		l.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

func (l *luaState) ToString(idx int) string {
	s, _ := l.ToStringX(idx)
	return s
}
