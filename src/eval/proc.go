package eval

import (
	"fmt"
	"os"
)

func isNullProc(args Object) Object {
	if isEmptyList(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isBooleanProc(args Object) Object {
	if isBool(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isSymbolProc(args Object) Object {
	if isSymbol(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isIntegerProc(args Object) Object {
	if isInt(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isCharProc(args Object) Object {
	if isChar(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isStringProc(args Object) Object {
	if isStr(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isPairProc(args Object) Object {
	if isPair(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func listProce(args Object) Object {
	return args
}

func addProcedure(name string, fun ObjFun, env Object) {
	defineVar(makeSymbol(name),
		makePrimitiveProc(fun),
		env)
}

func isNumEqualProc(args Object) Object {
	value := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		if value != asInt(car(args)) {
			return The_False
		}
		args = cdr(args)
	}
	return The_True
}

func isLessProc(args Object) Object {
	value := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		if value >= asInt(car(args)) {
			return The_False
		}
		args = cdr(args)
	}
	return The_True
}

func isLargerProc(args Object) Object {
	value := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		if value <= asInt(car(args)) {
			return The_False
		}
		args = cdr(args)
	}
	return The_True
}

func addProc(args Object) Object {
	res := int64(0)
	for {
		if isEmptyList(args) {
			break
		}
		res += asInt(car(args))
		args = cdr(args)
	}
	return makeInt(res)
}

func subProc(args Object) Object {
	res := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		res -= asInt(car(args))
		args = cdr(args)
	}
	return makeInt(res)
}

func mulProc(args Object) Object {
	res := int64(1)
	for {
		if isEmptyList(args) {
			break
		}
		res *= asInt(car(args))
		args = cdr(args)
	}
	return makeInt(res)
}

func divProc(args Object) Object {
	res := asInt(car(args))
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		next := asInt(car(args))
		if next == 0 {
			panic("divide zero")
		}
		res /= next
		args = cdr(args)
	}
	return makeInt(res)
}

func cadr(obj Object) Object {
	return car(cdr(obj))
}

func caar(obj Object) Object {
	return car(car(obj))
}

func cdar(obj Object) Object {
	return cdr(car(obj))
}

func consProc(args Object) Object {
	return cons(car(args), cadr(args))
}

func carProc(args Object) Object {
	return caar(args)
}

func cdrProc(args Object) Object {
	return cdar(args)
}

func setcarProc(args Object) Object {
	setCar(car(args), cadr(args))
	return OK_Symbol
}

func setcdrProc(args Object) Object {
	setCdr(car(args), cadr(args))
	return OK_Symbol
}

func listProc(args Object) Object {
	return args
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

func eqProc(args Object) Object {
	obj1 := car(args)
	obj2 := cadr(args)
	if equal(obj1, obj2) {
		return The_True
	} else {
		return The_False
	}
}

func errorProc(args Object) Object {
	fmt.Println("Error:")
	for {
		if isEmptyList(args) {
			break
		}
		fmt.Println(car(args))
		args = cdr(args)
	}
	os.Exit(1)
	return nil
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
