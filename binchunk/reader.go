package binchunk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"luago/llog"
	"math"
)

type reader struct {
	io.Reader
}

func NewReader(data []byte) *reader {
	return &reader{bytes.NewReader(data)}
}

func (r *reader) readBytes(n int) []byte {
	b := make([]byte, n)
	_, err := io.ReadFull(r.Reader, b)
	if err != nil {
		llog.Fatal("fail to read %d bytes: %s", n, err)
	}
	return b
}

func (r *reader) readByte() byte {
	b := make([]byte, 1)
	_, err := io.ReadFull(r.Reader, b)
	if err != nil {
		llog.Fatal("fail to read byte: %s", err)
	}
	return b[0]
}

func (r *reader) readUint32() uint32 {
	data := make([]byte, 4)
	_, err := io.ReadFull(r, data)
	if err != nil {
		llog.Fatal("fail to readUint32: %s", err)
	}
	return binary.LittleEndian.Uint32(data)
}

func (r *reader) readUint64() uint64 {
	data := make([]byte, 8)
	_, err := io.ReadFull(r, data)
	if err != nil {
		llog.Fatal("fail to readUint64: %s", err)
	}
	return binary.LittleEndian.Uint64(data)
}

func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

func (r *reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

func (r *reader) readString() string {
	const LongStrSize = 0xFF

	size := uint(r.readByte())
	if size == 0 {
		return ""
	}
	if size == LongStrSize {
		size = uint(r.readUint64())
	}
	bs := make([]byte, size-1)
	_, err := io.ReadFull(r.Reader, bs)
	if err != nil {
		llog.Fatal("readString expecting %d bytes, but got error: %s", size, err)
	}

	return string(bs)
}

func (r *reader) readConstant() interface{} {
	switch r.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return r.readByte() != 0
	case TAG_INTEGER:
		return r.readLuaInteger()
	case TAG_NUMBER:
		return r.readLuaNumber()
	case TAG_SHORT_STR, TAG_LONG_STR:
		return r.readString()
	default:
		llog.Fatal("corrupted")
		return nil
	}
}

func (r *reader) readConstants() []interface{} {
	constants := make([]interface{}, r.readUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

func (r *reader) readCode() []uint32 {
	code := make([]uint32, r.readUint32())
	for i := range code {
		code[i] = r.readUint32()
	}
	return code
}

func (r *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
		}
	}
	return upvalues
}

func (r *reader) readProtos(parentSrc string) []*Prototype {
	protos := make([]*Prototype, r.readUint32())
	for i := range protos {
		protos[i] = r.readProto(parentSrc)
	}
	return protos
}

func (r *reader) readLineInfo() []uint32 {
	lineInfos := make([]uint32, r.readUint32())
	for i := range lineInfos {
		lineInfos[i] = r.readUint32()
	}
	return lineInfos
}

func (r *reader) readLocVars() []LocVar {
	locvars := make([]LocVar, r.readUint32())
	for i := range locvars {
		locvars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readUint32(),
			EndPC:   r.readUint32(),
		}
	}
	return locvars
}

func (r *reader) readUpvalueNames() []string {
	upvalNames := make([]string, r.readUint32())
	for i := range upvalNames {
		upvalNames[i] = r.readString()
	}
	return upvalNames
}
func (r *reader) CheckHeader() {
	var errMsg string
	if sign := string(r.readBytes(4)); sign != LUA_SIGNATURE {
		errMsg = fmt.Sprintf("not a precompiled chunk(got %q)", sign)
	} else if r.readByte() != LUAC_VERSION {
		errMsg = "version mismatch"
	} else if r.readByte() != LUAC_FORMAT {
		errMsg = "format mismatch"
	} else if string(r.readBytes(6)) != LUAC_DATA {
		errMsg = "corrupted"
	} else if r.readByte() != CINT_SIZE {
		errMsg = "int size mismatch"
	} else if r.readByte() != CSIZET_SIZE {
		errMsg = "size_t size mismatch"
	} else if r.readByte() != INSTRUCTION_SIZE {
		errMsg = "instruction size mismatch"
	} else if r.readByte() != LUA_INTEGER_SIZE {
		errMsg = "lua_integer size mismatch"
	} else if r.readByte() != LUA_NUMBER_SIZE {
		errMsg = "lua_number size mismatch"
	} else if r.readLuaInteger() != LUAC_INT {
		errMsg = "endianness mismatch"
	} else if r.readLuaNumber() != LUAC_NUM {
		errMsg = "float format mismatch"
	}
	if errMsg != "" {
		llog.Throw("fail to check header: %s", errMsg)
	}
}

func (r *reader) readProto(parentSrc string) *Prototype {
	src := r.readString()
	if src == "" {
		src = parentSrc
	}
	return &Prototype{
		Source:          src,
		LineDefined:     r.readUint32(),
		LastLineDefined: r.readUint32(),
		NumParams:       r.readByte(),
		IsVararg:        r.readByte(),
		MaxStackSize:    r.readByte(),
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		Upvalues:        r.readUpvalues(),
		Protos:          r.readProtos(src),
		LineInfo:        r.readLineInfo(),
		LocVars:         r.readLocVars(),
		UpvalueNames:    r.readUpvalueNames(),
	}
}
