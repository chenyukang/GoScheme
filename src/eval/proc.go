package eval

import (
	"fmt"
	"os"
)

func makePrimitiveProc(fun ObjFun) *Object {
	obj := allocObject()
	obj.Type = PRIMITIVE_PROC
	obj.Data.primitive = fun
	return obj
}

func isPrimitiveProc(obj *Object) bool {
	return obj.Type == PRIMITIVE_PROC
}

func isNullProc(args *Object) *Object {
	if isEmptyList(args) {
		return The_True
	} else {
		return The_False
	}
}

func isBooleanProc(args *Object) *Object {
	if isBoolean(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isSymbolProc(args *Object) *Object {
	if isSymbol(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isIntegerProc(args *Object) *Object {
	if isFixNum(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isCharProc(args *Object) *Object {
	if isChar(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isStringProc(args *Object) *Object {
	if isString(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func isPairProc(args *Object) *Object {
	if isPair(car(args)) {
		return The_True
	} else {
		return The_False
	}
}

func addProcedure(name string, fun ObjFun, env *Object) {
	defineVar(makeSymbol(name),
		makePrimitiveProc(fun),
		env)
}

func isNumEqualProc(args *Object) *Object {
	value := (car(args)).Data.fixNum
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		if value != (car(args)).Data.fixNum {
			return The_False
		}
		args = cdr(args)
	}
	return The_True
}

func isLessProc(args *Object) *Object {
	value := (car(args)).Data.fixNum
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		if value >= (car(args)).Data.fixNum {
			return The_False
		}
		args = cdr(args)
	}
	return The_True
}

func isLargerProc(args *Object) *Object {
	value := (car(args)).Data.fixNum
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		if value <= (car(args)).Data.fixNum {
			return The_False
		}
		args = cdr(args)
	}
	return The_True
}

func addProc(args *Object) *Object {
	res := 0
	for {
		if isEmptyList(args) {
			break
		}
		res += (car(args)).Data.fixNum
		args = cdr(args)
	}
	return makeFixNum(res)
}

func subProc(args *Object) *Object {
	res := (car(args)).Data.fixNum
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		res -= (car(args)).Data.fixNum
		args = cdr(args)
	}
	return makeFixNum(res)
}

func mulProc(args *Object) *Object {
	res := 1
	for {
		if isEmptyList(args) {
			break
		}
		res *= (car(args)).Data.fixNum
		args = cdr(args)
	}
	return makeFixNum(res)
}

func divProc(args *Object) *Object {
	res := (car(args)).Data.fixNum
	args = cdr(args)
	for {
		if isEmptyList(args) {
			break
		}
		next := (car(args)).Data.fixNum
		if next == 0 {
			panic("divide zero")
		}
		res /= next
		args = cdr(args)
	}
	return makeFixNum(res)
}

func cadr(obj *Object) *Object {
	return car(cdr(obj))
}

func caar(obj *Object) *Object {
	return car(car(obj))
}

func cdar(obj *Object) *Object {
	return cdr(car(obj))
}

func consProc(args *Object) *Object {
	return cons(car(args), cadr(args))
}

func carProc(args *Object) *Object {
	return caar(args)
}

func cdrProc(args *Object) *Object {
	return cdar(args)
}

func setcarProc(args *Object) *Object {
	setCar(car(args), cadr(args))
	return OK_Symbol
}

func setcdrProc(args *Object) *Object {
	setCdr(car(args), cadr(args))
	return OK_Symbol
}

func listProc(args *Object) *Object {
	return args
}

func eqProc(args *Object) *Object {
	obj1 := car(args)
	obj2 := cadr(args)
	if obj1.Type != obj2.Type {
		return The_False
	}
	switch obj1.Type {
	case FIXNUM:
		if obj1.Data.fixNum == obj2.Data.fixNum {
			return The_True
		} else {
			return The_False
		}
	case CHARACTER:
		if obj1.Data.char == obj2.Data.char {
			return The_True
		} else {
			return The_False
		}
	case STRING:
		if obj1.Data.str == obj2.Data.str {
			return The_True
		} else {
			return The_False
		}
	default:
		if obj1 == obj2 {
			return The_True
		} else {
			return The_False
		}
	}
}

func errorProc(args *Object) *Object {
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

func makeEnv() *Object {
	env := extendEnv(The_EmptyList, The_EmptyList, The_Empty_Env)
	return env
}

func setupEnv(env *Object) {
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

	addProcedure("error", errorProc, env)

}

func Init() {
	SymbolTable = make(map[string]*Object)
	The_EmptyList = allocObject()
	The_EmptyList.Type = EMPTY_LIST
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
