package eval

import (
	"bufio"
	"fmt"
	"os"
)

type ObjType int

const (
	UNDEF = iota
	EMPTY_LIST
	BOOLEAN
	SYMBOL
	FIXNUM
	CHARACTER
	STRING
	PAIR
	PRIMITIVE_PROC
	COMPOUND_PROC
	INPUT_PORT
	OUTPUT_PORT
	EOF_OBJECT
)

type ObjFun func(args *Object) *Object

// For lack of union type in golang, ObjData will use for all type values
// Not effecient in memory now.
type ObjData struct {
	boolean int
	str     string
	symbol  string
	fixNum  int
	char    byte

	//pair data
	car *Object
	cdr *Object

	//primitive proc
	primitive ObjFun

	compond_params *Object
	compond_body   *Object
	compond_env    *Object
}

func assert(exp bool, msg string) {
	if exp == false {
		panic(msg)
	}
}

type Object struct {
	Type ObjType
	Data ObjData
}

func allocObject() *Object {
	obj := &Object{Type: UNDEF}
	return obj
}

var SymbolTable map[string]*Object

var The_EmptyList, The_True, The_False *Object
var The_Empty_Env, The_Global_Env *Object
var Set_Symbol, OK_Symbol, If_Symbol, Else_Symbol *Object
var Cond_Symbol, Or_Symbol, And_Symbol, Quote_Symbol, Lambda_Symbol *Object
var Define_Symbol, Begin_Symbol, Let_Symbol, Eof_Symbol *Object

func makeSymbol(sym string) *Object {
	if obj, ok := SymbolTable[sym]; ok {
		return obj
	}
	obj := allocObject()
	obj.Type = SYMBOL
	obj.Data.symbol = sym
	SymbolTable[sym] = obj
	return obj
}

func makeBoolean(val int) *Object {
	obj := allocObject()
	obj.Type = BOOLEAN
	obj.Data.boolean = val
	return obj
}

func makeFixNum(val int) *Object {
	obj := allocObject()
	obj.Type = FIXNUM
	obj.Data.fixNum = val
	return obj
}

func makeChar(char byte) *Object {
	obj := allocObject()
	obj.Type = CHARACTER
	obj.Data.char = char
	return obj
}

func makeString(str string) *Object {
	obj := allocObject()
	obj.Type = STRING
	obj.Data.str = str
	return obj
}

func isSymbol(obj *Object) bool {
	return obj.Type == SYMBOL
}

func isString(obj *Object) bool {
	return obj.Type == STRING
}

func isBoolean(obj *Object) bool {
	return obj.Type == BOOLEAN
}

func isFalse(obj *Object) bool {
	return obj == The_False
}

func isTrue(obj *Object) bool {
	return !isFalse(obj)
}

func isChar(obj *Object) bool {
	return obj.Type == CHARACTER
}

func isEmptyList(obj *Object) bool {
	return obj.Type == EMPTY_LIST
}

func isFixNum(obj *Object) bool {
	return obj.Type == FIXNUM
}

func isPair(obj *Object) bool {
	return obj.Type == PAIR
}

func setCar(pair *Object, obj *Object) {
	pair.Data.car = obj
}

func setCdr(pair *Object, obj *Object) {
	pair.Data.cdr = obj
}

func car(pair *Object) *Object {
	return pair.Data.car
}

func cdr(pair *Object) *Object {
	return pair.Data.cdr
}

func cons(car *Object, cdr *Object) *Object {
	obj := allocObject()
	obj.Type = PAIR
	obj.Data.car = car
	obj.Data.cdr = cdr
	return obj
}

func makeFrame(vars *Object, vals *Object) *Object {
	return cons(vars, vals)
}

func frameVars(frame *Object) *Object {
	return frame.Data.car
}

func frameVals(frame *Object) *Object {
	return frame.Data.cdr
}

func extendEnv(vars *Object, vals *Object, baseEnv *Object) *Object {
	return cons(makeFrame(vars, vals), baseEnv)
}

func addBinding(avar *Object, aval *Object, frame *Object) {
	frame.Data.car = cons(avar, car(frame))
	frame.Data.cdr = cons(aval, cdr(frame))
}

func firstFrame(env *Object) *Object {
	return env.Data.car
}

func lookupVar(avar *Object, env *Object) *Object {
	e := env
	for !isEmptyList(e) {
		frame := firstFrame(env)
		vars := frameVars(frame)
		vals := frameVals(frame)
		for !isEmptyList(vars) {
			if avar == vars.Data.car {
				return vals.Data.car
			}
			vars = vars.Data.cdr
			vals = vals.Data.cdr
		}
		e = e.Data.cdr
	}
	return nil
}

func defineVar(avar *Object, aval *Object, env *Object) {
	frame := firstFrame(env)
	vars := car(frame)
	vals := cdr(frame)
	for !isEmptyList(vars) {
		if avar == vars.Data.car {
			vals.Data.cdr = aval
			return
		}
		vars = vars.Data.cdr
		vals = vals.Data.cdr
	}
	addBinding(avar, aval, frame)
}

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

func ungetc(reader *bufio.Reader) {
	err := reader.UnreadByte()
	if err != nil {
		panic(err)
	}
}

func isSpace(val byte) bool {
	if val == '\t' || val == '\n' || val == '\r' || val == ' ' {
		return true
	} else {
		return false
	}
}

func eatWhiteSpace(reader *bufio.Reader) {
	for {
		c, err := reader.ReadByte()
		if err != nil {
			break
		}
		if isSpace(c) {
			continue
		} else if c == ';' {
			for {
				v, err := reader.ReadByte()
				if err != nil || v == '\n' {
					break
				}
			}
			continue
		}
		ungetc(reader)
		break
	}

}

func read(reader *bufio.Reader) *Object {
	return nil
}

func eval(exp *Object, env *Object) *Object {
	return nil
}

func write(obj *Object) {
}

func Run(reader *bufio.Reader) {
	fmt.Printf("Welcome to Bootstrap Scheme.\nUse ctrl-c to exit.\n")

	for {
		fmt.Printf("> ")
		exp := read(reader)
		if exp == nil {
			break
		}
		write(eval(exp, The_Global_Env))
		fmt.Println()
	}
}
