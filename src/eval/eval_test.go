package eval

import "testing"

func TestEmptyList(t *testing.T) {
	Init()
	if isEmptyList(The_Empty_List) == false {
		t.Error("empty_list failed")
	}
}

func TestSymbol(t *testing.T) {
	Init()
	res := isSymbol(OK_Symbol) && isSymbol(Set_Symbol) && isSymbol(If_Symbol)
	res = res && isSymbol(Begin_Symbol) && isSymbol(Let_Symbol)
	res = res && isSymbol(And_Symbol) && isSymbol(Or_Symbol)
	res = res && isSymbol(Lambda_Symbol)
	res = res && isSymbol(Quote_Symbol)
	if res == false {
		t.Error("Symbol failed")
	}

}

func TestFixNum(t *testing.T) {
	Init()
	fix := makeFixNum(1)
	if isFixNum(fix) == false {
		t.Error("fix type")
	}
	if fix.Data.fixNum != 1 {
		t.Error("fixNum val")
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
	if isPair(res) == false {
		t.Error("should be PAIR type")
	}
}
