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

type ObjData struct {
	boolean int
	strVal  string
	symbol  string
	fixNum  int
	char    byte

	//pair data
	car *Object
	cdr *Object
}

type Object struct {
	Type ObjType
	Data ObjData
}

func allocObject() *Object {
	obj := &Object{}
	obj.Type = UNDEF
	return obj
}

var SymbolTable map[string]*Object
var The_Empty_List *Object
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

func isSymbol(obj *Object) bool {
	return obj.Type == SYMBOL
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

func isEmptyList(obj *Object) bool {
	return obj.Type == EMPTY_LIST
}

func isPair(obj *Object) bool {
	return obj.Type == PAIR
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

func Init() {
	SymbolTable = make(map[string]*Object)
	The_Empty_List = allocObject()
	The_Empty_List.Type = EMPTY_LIST
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
