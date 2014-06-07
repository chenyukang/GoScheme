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

func TestSubProc(t *testing.T) {
	Init()
	arg1 := makeFixNum(1)
	arg2 := makeFixNum(2)
	arg3 := makeFixNum(3)
	args := cons(arg1,
		cons(arg2,
			cons(arg3, The_EmptyList)))
	if isEmptyList(args) {
		t.Error("sub proc")
	}
	res := subProc(args)
	if !(isFixNum(res) && res.Data.fixNum == -4) {
		t.Error("sub proc fail")
	}
	primitive := makePrimitiveProc(subProc)
	if isPrimitiveProc(primitive) == false {
		t.Error("primitivie")
	}
	res = (primitive.Data.primitive)(args)
	if !(isFixNum(res) && res.Data.fixNum == -4) {
		t.Error("sub proc fail")
	}
}

func TestDiv(t *testing.T) {
	Init()
	arg1 := makeFixNum(4)
	arg2 := makeFixNum(2)
	args := cons(arg1,
		cons(arg2, The_EmptyList))
	res := divProc(args)
	if !(isFixNum(res) && res.Data.fixNum == 2) {
		t.Error("div proc fail")
	}
	primitive := makePrimitiveProc(divProc)
	if isPrimitiveProc(primitive) == false {
		t.Error("primitivie")
	}
	res = (primitive.Data.primitive)(args)
	if !(isFixNum(res) && res.Data.fixNum == 2) {
		t.Error("div proc fail")
	}
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

func TestVarLookup(t *testing.T) {
	Init()
	sym := makeSymbol("+")
	obj, _ := lookupVar(sym, The_Global_Env)
	if !(obj.Type == PRIMITIVE_PROC &&
		obj.Data.primitive != nil) {
		t.Error("varlookup")
	}
	var1 := makeFixNum(1)
	var2 := makeFixNum(2)
	args := cons(var1, cons(var2, The_EmptyList))
	res := obj.Data.primitive(args)
	if !(isFixNum(res) && res.Data.fixNum == 3) {
		t.Error("varlookup error for addProc")
	}
}

func TestIsProc_Family(t *testing.T) {
	args := cons(The_True, The_EmptyList)
	res := isBooleanProc(args)
	if isTrue(res) == false {
		t.Error("isBooleanProc error")
	}
	sym := makeSymbol("OK")
	args = cons(sym, The_EmptyList)
	res = isSymbolProc(args)
	if isTrue(res) == false {
		t.Error("isSymbolProc error")
	}
}
