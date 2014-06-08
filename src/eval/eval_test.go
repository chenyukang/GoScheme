package eval

import "testing"

func evalWrapper(buf string) (*Object, error) {
	Init()
	obj, err := parserWrapper(buf)
	if err == nil {
		return eval(obj, The_Global_Env)
	} else {
		return nil, err
	}
}

func TestEvalDef(t *testing.T) {
	res, _ := evalWrapper("(define a 1)")
	if res != OK_Symbol {
		t.Error("define error")
	}

	res, _ = evalWrapper("(cons (define a 2) a)")
	if isPair(res) {
		val := cdr(res)
		if !equal(val, makeFixNum(2)) {
			t.Error("define failed")
		}
	} else {
		t.Error("not pair")
	}
}

func TestEvalSet(t *testing.T) {
	res, _ := evalWrapper("(set! a 1)")
	if res != OK_Symbol {
		t.Error("set! error")
	}

	res, _ = evalWrapper("(cons (set! a 2) a)")
	if isPair(res) {
		val := cdr(res)
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
	target, _ := lookupVar(makeSymbol("+"), The_Global_Env)
	if !equal(res, target) {
		t.Error("+ proc")
	}
	res, _ = evalWrapper("(+ 1 2)")
	if !equal(res, makeFixNum(3)) {
		t.Error("+ error")
	}
	res, _ = evalWrapper("(+ (- 1 2) 1)")
	if !(isFixNum(res) &&
		res.Data.fixNum == 0) {
		t.Error("+ error")
	}
	res, _ = evalWrapper("(* 10 (- 2 2))")
	if !equal(res, makeFixNum(0)) {
		t.Error("* error")
	}
	res, _ = evalWrapper("(null? ())")
	if res != The_True {
		t.Error("null? failed")
	}

	res, _ = evalWrapper("(boolean? #t)")
	if res != The_True {
		t.Error("boolean? failed")
	}
	res, _ = evalWrapper("(integer? 1)")
	if res != The_True {
		t.Error("integer? failed")
	}

	res, _ = evalWrapper("(= 1 2)")
	if res != The_False {
		t.Error("equal failed")
	}

	res, _ = evalWrapper("(list 1 2)")
	if !(isPair(res) &&
		equal(car(res), makeFixNum(1))) {
		t.Error("list failed")
	}
}

func TestEvalIf(t *testing.T) {
	res, _ := evalWrapper("(if 1 1 2)")
	if !equal(res, makeFixNum(1)) {
		t.Error("if failed")
	}
}
