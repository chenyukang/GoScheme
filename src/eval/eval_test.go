package eval

import (
	"fmt"
	"testing"
)

func evalWrapper(buf string) (*Object, error) {
	Init()
	obj, err := parserWrapper(buf)
	return eval(obj, The_Global_Env), err
}

func TestEvalSet(t *testing.T) {
	res, _ := evalWrapper("(set! a 1)")
	if res != OK_Symbol {
		t.Error("set! error")
	}

	res, _ = evalWrapper("(cons (set! a 2) a)")
	if isPair(res) {
		val := cdr(res)
		fmt.Println(val)
		if !(isFixNum(val) &&
			val.Data.fixNum == 2) {
			t.Error("set! failed")
		}
	} else {
		t.Error("not pair")
	}
}

func TestEvalProc(t *testing.T) {
	res, _ := evalWrapper("+")
	if !(isPrimitiveProc(res) &&
		res == lookupVar(makeSymbol("+"), The_Global_Env)) {
		t.Error("+ proc")
	}
	res, _ = evalWrapper("(+ 1 2)")
	if !(isFixNum(res) &&
		res.Data.fixNum == 3) {
		t.Error("+ error")
	}
	res, _ = evalWrapper("(+ (- 1 2) 1)")
	if !(isFixNum(res) &&
		res.Data.fixNum == 0) {
		t.Error("+ error")
	}
	res, _ = evalWrapper("(* 10 (- 2 2))")
	if !(isFixNum(res) &&
		res.Data.fixNum == 0) {
		t.Error("* error")
	}
}
