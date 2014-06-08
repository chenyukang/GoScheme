package eval

import "errors"

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

func isLast(obj *Object) bool {
	if isEmptyList(cdr(obj)) {
		return true
	} else {
		return false
	}
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

func lookupVar(avar *Object, env *Object) (*Object, error) {
	e := env
	for !isEmptyList(e) {
		frame := firstFrame(env)
		vars := frameVars(frame)
		vals := frameVals(frame)
		for !isEmptyList(vars) {
			if avar == vars.Data.car {
				return vals.Data.car, nil
			}
			vars = vars.Data.cdr
			vals = vals.Data.cdr
		}
		e = e.Data.cdr
	}
	return nil, errors.New("undef variable")
}

func defineVar(avar *Object, aval *Object, env *Object) {
	frame := firstFrame(env)
	vars := car(frame)
	vals := cdr(frame)
	for !isEmptyList(vars) {
		if avar == car(vars) {
			setCar(vals, aval)
			return
		}
		vars = cdr(vars)
		vals = cdr(vals)
	}
	addBinding(avar, aval, frame)
}
