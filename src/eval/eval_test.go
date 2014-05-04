package eval

import (
	"testing"
)

func TestEmptyList(t *testing.T) {
	Init()
	if IsEmptyList(The_Empty_List) == false {
		t.Error("empty_list failed")
	}
}

func TestSymbol(t *testing.T) {
	Init()
	res := IsSymbol(OK_Symbol) && IsSymbol(Set_Symbol) && IsSymbol(If_Symbol)
	res = res && IsSymbol(Begin_Symbol) && IsSymbol(Let_Symbol)
	res = res && IsSymbol(And_Symbol) && IsSymbol(Or_Symbol)
	res = res && IsSymbol(Lambda_Symbol)
	res = res && IsSymbol(Quote_Symbol)
	if res == false {
		t.Error("Symbol failed")
	}

}

func TestTrueFalse(t *testing.T) {
	Init()
	if IsFalse(The_False) == false {
		t.Error("False failed")
	}
	if IsTrue(The_True) == false {
		t.Error("True failed")
	}
}

func TestEverything(t *testing.T) {
	Init()
	res := Cons(OK_Symbol, If_Symbol)
	if IsPair(res) == false {
		t.Error("should be PAIR type")
	}
}
