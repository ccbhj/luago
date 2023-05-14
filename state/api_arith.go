package state

import (
	. "luago/api"
	"luago/number"
	"math"
)

var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = number.IMod
	fmod  = number.FMod
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = number.IFloorDiv
	fidiv = number.FFloorDiv
	band  = func(a, b int64) int64 { return a & b }
	bor   = func(a, b int64) int64 { return a | b }
	bxor  = func(a, b int64) int64 { return a ^ b }
	shl   = number.ShiftLeft
	shr   = number.ShiftRight
	iunm  = func(a, _ int64) int64 { return -a }
	funm  = func(a, _ float64) float64 { return -a }
	bnot  = func(a, _ int64) int64 { return ^a }
)

type operator struct {
	metamethod string
	integerFn  func(int64, int64) int64
	floatFn    func(float64, float64) float64
}

var operators = [...]operator{
	{"__add", iadd, fadd},
	{"__sub", isub, fsub},
	{"__mul", imul, fmul},
	{"__mod", imod, fmod},
	{"__pow", nil, pow},
	{"__div", nil, div},
	{"__idiv", iidiv, fidiv},
	{"__band", band, nil},
	{"__bor", bor, nil},
	{"__bxor", bxor, nil},
	{"__shl", shl, nil},
	{"__shr", shr, nil},
	{"__unm", iunm, funm},
	{"__bnot", bnot, nil},
}

func (l *luaState) Arith(op ArithOp) {
	var a, b luaValue
	b = l.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = l.stack.pop()
	} else {
		a = b
	}

	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		l.stack.push(result)
		return
	}

	mm := operator.metamethod
	if result, ok := callMetamethod(a, b, mm, l); ok {
		l.stack.push(result)
		return
	}

	panic("arithmetic error!")
}

func _arith(a, b luaValue, op operator) luaValue {
	// bitwise
	if op.floatFn == nil {
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFn(x, y)
			}
		}
		return nil
	}

	// arith
	if op.integerFn != nil {
		if x, ok := a.(int64); ok {
			if y, ok := b.(int64); ok {
				return op.integerFn(x, y)
			}
		}
	}

	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return op.floatFn(x, y)
		}
	}

	return nil
}

func callMetamethod(a, b luaValue,
	mmName string, ls *luaState) (luaValue, bool) {

	var mm luaValue
	mm = getMetafield(a, mmName, ls)
	if mm == nil {
		if mm = getMetafield(b, mmName, ls); mm == nil {
			return nil, false
		}
	}

	ls.stack.check(4)
	ls.stack.push(mm)
	ls.stack.push(a)
	ls.stack.push(b)
	ls.Call(2, 1)

	return ls.stack.pop(), true
}
