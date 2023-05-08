package api

const (
	// -1001k       -1000k             0                 1000k
	// +-------------+-----------------+-------------------+  stack size consts
	// registery   min_valid                             max_valid_idx
	//  idx         idx
	LUA_MINSTACK      = 20
	LUAI_MAXSTACK     = 1_000_000
	LUA_RIDX_GLOBALS  = int64(2)
	LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000
)

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
		return "<UNKNOWN_TYPE>"
	}
	return v
}

type ArithOp int

const (
	LUA_OPADD  ArithOp = iota // +
	LUA_OPSUB                 // -
	LUA_OPMUL                 // *
	LUA_OPMOD                 // %
	LUA_OPPOW                 // ^
	LUA_OPDIV                 // /
	LUA_OPIDIV                // //
	LUA_OPBAND                // &
	LUA_OPBOR                 //  |
	LUA_OPBXOR                // ~
	LUA_OPSHL                 // <<
	LUA_OPSHR                 // >>
	LUA_OPUNM                 // -
	LUA_OPBNOT                // ~
)

var arithOpNames = [...]string{
	LUA_OPADD:  "+",
	LUA_OPSUB:  "-",
	LUA_OPMUL:  "*",
	LUA_OPMOD:  "%",
	LUA_OPPOW:  "^",
	LUA_OPDIV:  "/",
	LUA_OPIDIV: "//",
	LUA_OPBAND: "&",
	LUA_OPBOR:  " ",
	LUA_OPBXOR: "~",
	LUA_OPSHL:  "<",
	LUA_OPSHR:  ">",
	LUA_OPUNM:  "-",
	LUA_OPBNOT: "~",
}

func (o ArithOp) String() string {
	if o < 0 || int(o) >= len(arithOpNames) {
		return "<UNKNOWN_AIRTH_OP>"
	}
	return arithOpNames[o]
}

type CompareOp int

const (
	LUA_OPEQ CompareOp = iota // ==
	LUA_OPLT                  // <
	LUA_OPLE                  // <=
)

var compareOpName = [...]string{
	LUA_OPEQ: "==",
	LUA_OPLT: "<",
	LUA_OPLE: "<=",
}

func (o CompareOp) String() string {
	if o < 0 || int(o) >= len(compareOpName) {
		return "<UNKNOWN_AIRTH_OP>"
	}
	return compareOpName[o]
}
