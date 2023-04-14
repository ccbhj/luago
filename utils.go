package main

import (
	"fmt"
	"luago/binchunk"

	. "luago/vm"
)

func list(proto *binchunk.Prototype) {
	printHeader(proto)
	printCode(proto)
	printDetail(proto)
	for _, p := range proto.Protos {
		list(p)
	}
}

func printHeader(proto *binchunk.Prototype) {
	funcType := "main"
	if proto.LineDefined > 0 {
		funcType = "function"
	}
	varargFlag := ""
	if proto.IsVararg > 0 {
		varargFlag = "+"
	}
	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n", funcType,
		proto.Source, proto.LineDefined, proto.LastLineDefined, len(proto.Code))
	fmt.Printf("%d%s params, %d slots, %d upvalues, \n",
		proto.NumParams, varargFlag, proto.MaxStackSize, len(proto.Upvalues))
	fmt.Printf("%d locals, %d constants, %d functions\n",
		len(proto.LocVars), len(proto.Constants), len(proto.Protos))
}

func printCode(proto *binchunk.Prototype) {
	for pc, c := range proto.Code {
		line := "-"
		if len(proto.LineInfo) > 0 {
			line = fmt.Sprintf("%d", proto.LineInfo[pc])
		}
		i := Instruction(c)
		fmt.Printf("\t%4d\t[%s]\t%-10s \t", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Println()
	}
}

func printOperands(i Instruction) {
	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()
		fmt.Printf("%5d", a)
		if OpArgMode(i.BMode()) != OpArgN {
			if b > 0xFF {
				//  highest bit is 1 means a constant table index
				fmt.Printf(" %5d", -1-b&0xFF)
			} else {
				fmt.Printf(" %5d", b)
			}
		}
		if OpArgMode(i.CMode()) != OpArgN {
			if c > 0xFF {
				fmt.Printf(" %5d", -1-c&0xFF)
			} else {
				fmt.Printf(" %5d", c)
			}
		}
	case IABx:
		a, bx := i.ABx()
		fmt.Printf("%5d", a)
		switch OpArgMode(i.BMode()) {
		case OpArgK:
			// OpArgK means a constant table index
			fmt.Printf(" %5d", -1-bx)
		case OpArgU:
			fmt.Printf(" %5d", bx)
		}
	case IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%5d %5d", a, sbx)
	case IAx:
		ax := i.Ax()
		fmt.Printf("%5d", -1-ax)
	}
}

func printDetail(proto *binchunk.Prototype) {
	fmt.Printf("constants (%d):\n", len(proto.Constants))
	for i, k := range proto.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToStr(k))
	}

	fmt.Printf("locals (%d):\n", len(proto.LocVars))
	for i, locVar := range proto.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, locVar.VarName, locVar.StartPC, locVar.EndPC)
	}

	fmt.Printf("upvalues (%d):\n", len(proto.Upvalues))
	for i, upval := range proto.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, upvalName(proto, i), upval.Instack, upval.Idx)
	}
}

func constantToStr(i interface{}) string {
	switch v := i.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case string:
		return fmt.Sprintf("%q", v)
	}
	return "?"
}

func upvalName(proto *binchunk.Prototype, i int) string {
	if len(proto.UpvalueNames) > 0 {
		return proto.UpvalueNames[i]
	}
	return "-"
}
