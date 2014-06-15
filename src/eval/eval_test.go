package eval

import "testing"

func evalWrapper(buf string) (Object, error) {
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
		if !equal(val, makeInt(2)) {
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
		if !(isInt(val) &&
			asInt(val) == 2) {
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
	if !equal(res, makeInt(3)) {
		t.Error("+ error")
	}
	res, _ = evalWrapper("(+ (- 1 2) 1)")
	if !(isInt(res) &&
		asInt(res) == 0) {
		t.Error("+ error")
	}
	res, _ = evalWrapper("(* 10 (- 2 2))")
	if !equal(res, makeInt(0)) {
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
		equal(car(res), makeInt(1))) {
		t.Error("list failed")
	}
}

func TestEvalIf(t *testing.T) {
	res, _ := evalWrapper("(if 1 1 2)")
	if !equal(res, makeInt(1)) {
		t.Error("if failed")
	}
}

func TestEvalLogical(t *testing.T) {
	res, _ := evalWrapper("(or #t #t #f)")
	if !equal(res, The_True) {
		t.Error("logical error")
	}
	res, _ = evalWrapper("(or #f #f #f #f)")
	if !equal(res, The_False) {
		t.Error("logical error")
	}
	res, _ = evalWrapper("(and #f #t #t)")
	if !equal(res, The_False) {
		t.Error("logical error")
	}
	res, _ = evalWrapper("(and #t #t #t)")
	if !equal(res, The_True) {
		t.Error("logical error")
	}

	res, _ = evalWrapper("(and (= 1 1) (= 2 2)")
	if !equal(res, The_True) {
		t.Error("logical error")
	}
}

func TestEvalCond(t *testing.T) {
	res, _ := evalWrapper("(cond ((#t 1) (#t 2)))")
	if !equal(res, makeInt(1)) {
		t.Error("cond failed")
	}
	res, _ = evalWrapper("(cond ((#f 1) (#t 2)))")
	if !equal(res, makeInt(2)) {
		t.Error("cond failed")
	}
	res, _ = evalWrapper("(cond ((#f 1) (#f 2) (else 3)))")
	if !equal(res, makeInt(3)) {
		t.Error("cond failed")
	}
}
