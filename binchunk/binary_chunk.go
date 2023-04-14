package binchunk

// constant in header
const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

// constant tag
const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type header struct {
	signature       [4]byte // magic number OxlB4C7561 ('\x1Lua')
	version         byte    // lua version => ((Major << 4) + Minor)
	format          byte    //
	luacData        [6]byte
	cintSize        byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         int64   // constant number LUAC_INT, to validate the endian of the machaine
	luacNum         float64 // constant float LUAC_NUM, to check the float format of the machaine
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

type Prototype struct {
	// function source, could be:
	//  +--------+------------------------------------------------------+
	//  | prefix |  description                                         |
	//  +--------+------------------------------------------------------+
	//  | '@'    | file name if compiled from file and its main function|
	//  | '='    | special case like '=stdin' from stdin                |
	//  | ''     | the whole lua string that we compiled from           |
	//  +--------+------------------------------------------------------+
	Source string
	// for main function, these two line number should be 0, non-zero for other files
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte     // is variadic function ? 0 for FALSE, 1 for TRUE
	MaxStackSize    byte     // how many register(based on stack) at least do we need
	Code            []uint32 // instruction table, index by pc
	// constant table(tag[1 byte] + value)
	// +------+------+-----------+
	// | tag  | type | value     |
	// +------+------+-----------+
	// | 0x00 | nil  |  -        |
	// | 0x01 | bool | byte(0/1) |
	// | 0x02 | num  | lua float |
	// | 0x13 | int  | lua int   |
	// | 0x04 | str  | short str |
	// | 0x14 | str  | long str  |
	// +------+------+-----------+
	Constants []interface{}
	Upvalues  []Upvalue
	Protos    []*Prototype // prototype for sub-function

	// attributes for debugging
	LineInfo     []uint32 // line number for each instruction, index by pc
	LocVars      []LocVar // local variables table
	UpvalueNames []string
}

type binaryChunk struct {
	header                  // 头
	sizeUpvalues byte       // 主函数upvalue数量
	mainFunc     *Prototype // 主函数原型
}

func Undump(data []byte) *Prototype {
	reader := NewReader(data)
	reader.CheckHeader()
	reader.readByte() // skip reading upvalues size
	return reader.readProto("")
}
