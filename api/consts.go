package api

import "strconv"

type LuaType int

const (
	LUA_TNONE LuaType = iota - 1 // -1
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

var typeNames = map[LuaType]string{
	LUA_TNONE:     "no value",
	LUA_TNIL:      "nil",
	LUA_TBOOLEAN:  "boolean",
	LUA_TNUMBER:   "number",
	LUA_TSTRING:   "string",
	LUA_TTABLE:    "table",
	LUA_TFUNCTION: "function",
	LUA_TUSERDATA: "userdata",
	LUA_TTHREAD:   "thread",
}

func (s LuaType) String() string {
	v, in := typeNames[s]
	if !in {
		panic("unknown type: " + strconv.Itoa(int(s)))
	}
	return v
}
