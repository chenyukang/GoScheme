package eval

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

type Object struct {
	Type ObjType
	Data ObjData
}

func allocObject() *Object {
	obj := &Object{Type: UNDEF}
	return obj
}

var SymbolTable map[string]*Object
var The_EmptyList *Object
var The_True *Object
var The_False *Object
var Set_Symbol *Object
var OK_Symbol *Object
var If_Symbol *Object
var Else_Symbol *Object
var Cond_Symbol *Object
var Or_Symbol *Object
var Quote_Symbol *Object
var Lambda_Symbol *Object
var Define_Symbol *Object
var Begin_Symbol *Object
var Let_Symbol *Object
var And_Symbol *Object
var Eof_Symbol *Object
var The_Empty_Env *Object
var The_Global_Env *Object

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
	if pair.Type != PAIR {
		panic("setCar")
	}
	pair.Data.car = obj
}

func setCdr(pair *Object, obj *Object) {
	if pair.Type != PAIR {
		panic("setCdr")
	}
	pair.Data.cdr = obj
}

func car(pair *Object) *Object {
	if pair.Type != PAIR {
		panic("car")
	}
	return pair.Data.car
}

func cdr(pair *Object) *Object {
	if pair.Type != PAIR {
		panic("cdr")
	}
	return pair.Data.cdr
}

func cons(car *Object, cdr *Object) *Object {
	obj := allocObject()
	obj.Type = PAIR
	obj.Data.car = car
	obj.Data.cdr = cdr
	return obj
}

type Frame Object

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
	frame.Data.car = cons(avar, frame.Data.car)
	frame.Data.cdr = cons(aval, frame.Data.cdr)
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

func makeEnv() *Object {
	env := extendEnv(The_EmptyList, The_EmptyList, The_Empty_Env)
	return env
}

func defineVar(avar *Object, aval *Object, env *Object) {
	frame := firstFrame(env)
	vars := frame.Data.car
	vals := frame.Data.cdr
	for !isEmptyList(vars) {
		if avar == vars.Data.car {
			vals.Data.cdr = aval
			return
		}
		vars = vars.Data.cdr
		vals = vals.Data.cdr
	}
	addBinding(avar, aval, env)
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

func addProcedure(name string, fun ObjFun, env *Object) {
	defineVar(makeSymbol(name),
		makePrimitiveProc(fun),
		env)
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

func setupEnv(env *Object) {
	addProcedure("+", addProc, env)
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
}
