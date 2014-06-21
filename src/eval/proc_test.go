package eval

import "testing"

func TestEmptyList(t *testing.T) {
	Init()
	if !isEmptyList(The_EmptyList) {
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
	fix := makeInt(1)
	if !(isInt(fix) && asInt(fix) == 1) {
		t.Error("fixNum val")
	}
}

func TestAddProc(t *testing.T) {
	Init()
	arg1 := makeInt(1)
	arg2 := makeInt(2)
	arg3 := makeInt(3)
	args := cons(arg1,
		cons(arg2,
			cons(arg3, The_EmptyList)))
	if isEmptyList(args) {
		t.Error("add proc")
	}
	res, _ := addProc(args)
	if asInt(res) != 6 {
		t.Error("add proc fail")
	}
	primitive := makePrimitiveProc(addProc)
	if isPrimitiveProc(primitive) == false {
		t.Error("primitivie")
	}
	fu := asFunc(primitive)
	res, _ = fu(args)
}

func TestSubProc(t *testing.T) {
	Init()
	arg1 := makeInt(1)
	arg2 := makeInt(2)
	arg3 := makeInt(3)
	args := cons(arg1,
		cons(arg2,
			cons(arg3, The_EmptyList)))
	if isEmptyList(args) {
		t.Error("sub proc")
	}
	res, _ := subProc(args)
	if !(isInt(res) && asInt(res) == -4) {
		t.Error("sub proc fail")
	}
	primitive := makePrimitiveProc(subProc)
	if isPrimitiveProc(primitive) == false {
		t.Error("primitivie")
	}
	fu := asFunc(primitive)
	res, _ = fu(args)
	if !(isInt(res) && asInt(res) == -4) {
		t.Error("sub proc fail")
	}
}

func TestDiv(t *testing.T) {
	Init()
	arg1 := makeInt(4)
	arg2 := makeInt(2)
	args := cons(arg1,
		cons(arg2, The_EmptyList))
	res, _ := divProc(args)
	if !(isInt(res) && asInt(res) == 2) {
		t.Error("div proc fail")
	}
	primitive := makePrimitiveProc(divProc)
	if isPrimitiveProc(primitive) == false {
		t.Error("primitivie")
	}
	fu := asFunc(primitive)
	res, _ = fu(args)
	if !(isInt(res) && asInt(res) == 2) {
		t.Error("div proc fail")
	}
}

func TestChar(t *testing.T) {
	Init()
	char := makeChar('a')
	if isChar(char) == false ||
		asChar(char) != 'a' {
		t.Error("char error")
	}
}

func TestString(t *testing.T) {
	Init()
	val := "this is good"
	str := makeStr(val)
	if isStr(str) == false ||
		asStr(str) != val {
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
	if !(isPair(res) &&
		car(res) == OK_Symbol &&
		cdr(res) == If_Symbol) {
		t.Error("should be equal on cdr")
	}
	setCar(res, If_Symbol)
	setCdr(res, OK_Symbol)
	if !(isPair(res) &&
		car(res) == If_Symbol &&
		cdr(res) == OK_Symbol) {
		t.Error("should be equal on cdr")
	}
}

func TestVarLookup(t *testing.T) {
	Init()
	sym := makeSymbol("+")
	obj, _ := lookupVar(sym, The_Global_Env)
	if !(typeOf(obj) == PRIMITIVE_PROC &&
		asFunc(obj) != nil) {
		t.Error("varlookup")
	}
	var1 := makeInt(1)
	var2 := makeInt(2)
	args := cons(var1, cons(var2, The_EmptyList))
	fu := asFunc(obj)
	res, _ := fu(args)
	if !(isInt(res) && asInt(res) == 3) {
		t.Error("varlookup error for addProc")
	}
}

func TestIsProc_Family(t *testing.T) {
	args := cons(The_True, The_EmptyList)
	res, _ := isBooleanProc(args)
	if isTrue(res) == false {
		t.Error("isBooleanProc error")
	}
	sym := makeSymbol("OK")
	args = cons(sym, The_EmptyList)
	res, _ = isSymbolProc(args)
	if isTrue(res) == false {
		t.Error("isSymbolProc error")
	}

	res, _ = isNullProc(args)
	if isTrue(res) {
		t.Error("isNullProc failed")
	}

	res, _ = isBooleanProc(args)
	if isTrue(res) {
		t.Error("isBooleanProc failed")
	}

	res, _ = isIntegerProc(args)
	if isTrue(res) {
		t.Error("isIntegerProc failed")
	}

	args = cons(makeInt(1), makeInt(2))
	res, _ = isSymbolProc(args)
	if isTrue(res) {
		t.Error("isSymbolProc failed")
	}

	res, _ = isCharProc(cons(makeChar('a'), makeChar('b')))
	if isFalse(res) {
		t.Error("isCharProc failed")
	}

	res, _ = isCharProc(cons(makeInt(1), makeInt(2)))
	if isTrue(res) {
		t.Error("isCharProc failed")
	}

}
