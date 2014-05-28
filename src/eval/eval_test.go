package eval

import "testing"

func TestEmptyList(t *testing.T) {
	Init()
	if isEmptyList(The_EmptyList) == false {
		t.Error("empty_list failed")
	}
}

func TestSymbol(t *testing.T) {
	Init()
	res := isSymbol(OK_Symbol) &&
		isSymbol(Set_Symbol) &&
		isSymbol(If_Symbol) &&
		isSymbol(Begin_Symbol) &&
		isSymbol(Let_Symbol) &&
		isSymbol(And_Symbol) &&
		isSymbol(Or_Symbol) &&
		isSymbol(Lambda_Symbol) &&
		isSymbol(Quote_Symbol)
	if !res {
		t.Error("Symbol failed")
	}

}

func TestFixNum(t *testing.T) {
	Init()
	fix := makeFixNum(1)
	if !(isFixNum(fix) && fix.Data.fixNum == 1) {
		t.Error("fixNum val")
	}
}

func TestAddProc(t *testing.T) {
	Init()
	arg1 := makeFixNum(1)
	arg2 := makeFixNum(2)
	arg3 := makeFixNum(3)
	args := cons(arg1,
		cons(arg2,
			cons(arg3, The_EmptyList)))
	if isEmptyList(args) {
		t.Error("add proc")
	}
	res := addProc(args)
	if res.Data.fixNum != 6 {
		t.Error("add proc fail")
	}
	primitive := makePrimitiveProc(addProc)
	if isPrimitiveProc(primitive) == false {
		t.Error("primitivie")
	}
	res = (primitive.Data.primitive)(args)
}

func TestChar(t *testing.T) {
	Init()
	char := makeChar('a')
	if isChar(char) == false ||
		char.Data.char != 'a' {
		t.Error("char error")
	}
}

func TestString(t *testing.T) {
	Init()
	val := "this is good"
	str := makeString(val)
	if isString(str) == false ||
		str.Data.str != val {
		t.Error("string error")
	}
}

func TestTrueFalse(t *testing.T) {
	Init()
	if isFalse(The_False) == false {
		t.Error("False failed")
	}
	if isTrue(The_True) == false {
		t.Error("True failed")
	}
}

func TestEverything(t *testing.T) {
	Init()
	res := cons(OK_Symbol, If_Symbol)
	if !(isPair(res) && car(res) == OK_Symbol && cdr(res) == If_Symbol) {
		t.Error("should be equal on cdr")
	}
	setCar(res, If_Symbol)
	setCdr(res, OK_Symbol)
	if !(isPair(res) && car(res) == If_Symbol && cdr(res) == OK_Symbol) {
		t.Error("should be equal on cdr")
	}
}
