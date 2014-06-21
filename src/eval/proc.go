package eval

import (
	"errors"
	"fmt"
	"os"
)

func isNullProc(args Object) (Object, error) {
	if isEmptyList(car(args)) {
		return The_True, nil
	} else {
		return The_False, nil
	}
}

func isBooleanProc(args Object) (Object, error) {
	if isBool(car(args)) {
		return The_True, nil
	} else {
		return The_False, nil
	}
}

func isSymbolProc(args Object) (Object, error) {
	if isSymbol(car(args)) {
		return The_True, nil
	} else {
		return The_False, nil
	}
}

func isIntegerProc(args Object) (Object, error) {
	if isInt(car(args)) {
		return The_True, nil
	} else {
		return The_False, nil
	}
}

func isCharProc(args Object) (Object, error) {
	if isChar(car(args)) {
		return The_True, nil
	} else {
		return The_False, nil
	}
}

func isStringProc(args Object) (Object, error) {
	if isStr(car(args)) {
		return The_True, nil
	} else {
		return The_False, nil
	}
}

func isPairProc(args Object) (Object, error) {
	if isPair(car(args)) {
		return The_True, nil
	} else {
		return The_False, nil
	}
}

func listProce(args Object) (Object, error) {
	return args, nil
}

func addProcedure(name string, fun ObjFun, env Object) {
	defineVar(makeSymbol(name),
		makePrimitiveProc(fun),
		env)
}

func isNumEqualProc(args Object) (Object, error) {
	value := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		if value != asInt(car(args)) {
			return The_False, nil
		}
		args = cdr(args)
	}
	return The_True, nil
}

func isLessProc(args Object) (Object, error) {
	value := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		if value >= asInt(car(args)) {
			return The_False, nil
		}
		args = cdr(args)
	}
	return The_True, nil
}

func isLargerProc(args Object) (Object, error) {
	value := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		if value <= asInt(car(args)) {
			return The_False, nil
		}
		args = cdr(args)
	}
	return The_True, nil
}

func addProc(args Object) (Object, error) {
	res := int64(0)
	for {
		if isEmptyList(args) {
			break
		}
		res += asInt(car(args))
		args = cdr(args)
	}
	return makeInt(res), nil
}

func subProc(args Object) (Object, error) {
	res := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		res -= asInt(car(args))
		args = cdr(args)
	}
	return makeInt(res), nil
}

func mulProc(args Object) (Object, error) {
	res := int64(1)
	for {
		if isEmptyList(args) {
			break
		}
		res *= asInt(car(args))
		args = cdr(args)
	}
	return makeInt(res), nil
}

func divProc(args Object) (Object, error) {
	res := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		next := asInt(car(args))
		if next == 0 {
			return FAIL_Symbol, errors.New("divide zero")
		}
		res /= next
		args = cdr(args)
	}
	return makeInt(res), nil
}

func cadr(obj Object) Object {
	return car(cdr(obj))
}

func cddr(obj Object) Object {
	return cdr(cdr(obj))
}

func caar(obj Object) Object {
	return car(car(obj))
}

func cdar(obj Object) Object {
	return cdr(car(obj))
}

func consProc(args Object) (Object, error) {
	return cons(car(args), cadr(args)), nil
}

func carProc(args Object) (Object, error) {
	return caar(args), nil
}

func cdrProc(args Object) (Object, error) {
	return cdar(args), nil
}

func setcarProc(args Object) (Object, error) {
	setCar(car(args), cadr(args))
	return OK_Symbol, nil
}

func setcdrProc(args Object) (Object, error) {
	setCdr(car(args), cadr(args))
	return OK_Symbol, nil
}

func listProc(args Object) (Object, error) {
	return args, nil
}

func equal(obj1 Object, obj2 Object) bool {
	type1, type2 := typeOf(obj1), typeOf(obj2)
	if type1 != type2 {
		return false
	}
	switch type1 {
	case INT:
		if asInt(obj1) == asInt(obj2) {
			return true
		} else {
			return false
		}
	case CHARACTER:
		if asChar(obj1) == asChar(obj2) {
			return true
		} else {
			return false
		}
	case STRING:
		if asStr(obj1) == asStr(obj2) {
			return true
		} else {
			return false
		}
	default:
		if obj1 == obj2 {
			return true
		} else {
			return false
		}
	}
}

func eqProc(args Object) (Object, error) {
	obj1 := car(args)
	obj2 := cadr(args)
	if equal(obj1, obj2) {
		return The_True, nil
	} else {
		return The_False, nil
	}
}

func errorProc(args Object) (Object, error) {
	fmt.Println("Error:")
	for {
		if isEmptyList(args) {
			break
		}
		fmt.Println(car(args))
		args = cdr(args)
	}
	os.Exit(1)
	return nil, nil
}

func makeEnv() Object {
	env := extendEnv(The_EmptyList, The_EmptyList, The_Empty_Env)
	return env
}

func setupEnv(env Object) {
	addProcedure("+", addProc, env)
	addProcedure("-", subProc, env)
	addProcedure("*", mulProc, env)
	addProcedure("/", divProc, env)
	addProcedure("=", isNumEqualProc, env)
	addProcedure("<", isLessProc, env)
	addProcedure(">", isLargerProc, env)

	addProcedure("null?", isNullProc, env)
	addProcedure("boolean?", isBooleanProc, env)
	addProcedure("integer?", isIntegerProc, env)
	addProcedure("symbol?", isSymbolProc, env)
	addProcedure("string?", isStringProc, env)
	addProcedure("pair?", isPairProc, env)

	addProcedure("cons", consProc, env)
	addProcedure("car", carProc, env)
	addProcedure("cdr", cdrProc, env)
	addProcedure("set-car!", setcarProc, env)
	addProcedure("set-cdr!", setcdrProc, env)

	addProcedure("list", listProc, env)
	addProcedure("error", errorProc, env)

}

func Init() {
	SymbolTable = make(map[string]Object)
	The_EmptyList = makeEmptyList()
	The_True = makeBoolean(1)
	The_False = makeBoolean(0)
	Set_Symbol = makeSymbol("set!")
	OK_Symbol = makeSymbol("ok")
	FAIL_Symbol = makeSymbol("fail")
	If_Symbol = makeSymbol("if")
	Else_Symbol = makeSymbol("else")
	Let_Symbol = makeSymbol("let")
	And_Symbol = makeSymbol("and")
	Or_Symbol = makeSymbol("or")
	Define_Symbol = makeSymbol("define")
	Begin_Symbol = makeSymbol("begin")
	Quote_Symbol = makeSymbol("quote")
	Lambda_Symbol = makeSymbol("lambda")
	Cond_Symbol = makeSymbol("cond")
	The_Empty_Env = The_EmptyList
	The_Global_Env = makeEnv()
	setupEnv(The_Global_Env)
}
