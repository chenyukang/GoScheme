package eval

import "testing"

func evalWrapper(buf string) (*Object, error) {
	obj, err := parserWrapper(buf)
	return eval(obj, The_Global_Env), err
}

func TestEval(t *testing.T) {
	Init()
	res, _ := evalWrapper("+")
	if !(isPrimitiveProc(res) &&
		res == lookupVar(makeSymbol("+"), The_Global_Env)) {
		t.Error("+ proc")
	}
}
