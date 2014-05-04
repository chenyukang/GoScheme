package eval

type Object interface {}
type Value  interface {}

type ObjType int
const (
        EMPTY_LIST = iota
        BOOLEAN
        SYBOL
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

type EmptyList struct {
        objType ObjType
}

type Boolean struct {
        objType ObjType
        value int
}

type Symbol struct {
        objType ObjType
        value string
}

type FixNum struct {
        objType ObjType
        value   int64
}

type Character struct {
        objType ObjType
        value   byte
}

type Pair struct {
        objType ObjType
        car     Object
        cdr     Object
}

type Proc struct {
        objType ObjType
}

var Symbol_Table   map[string]Object
var The_Empty_List Object
var The_True       Object
var The_False      Object
var Set_Symbol     Object
var OK_Symbol      Object
var If_Symbol      Object
var Else_Symbol    Object
var Cond_Symbol    Object
var Or_Symbol      Object
var Quote_Symbol   Object
var Lambda_Symbol  Object
var Define_Symbol  Object
var Begin_Symbol   Object
var Let_Symbol     Object
var And_Symbol     Object
var Eof_Symbol     Object
var The_Empty_Env  Object
var The_Global_Env Object

func MakeSymbol(sym string) Object {
        if res, ok := Symbol_Table[sym]; ok {
                return res
        }
        obj := &Symbol{ STRING, sym}
        Symbol_Table[sym] = (Object)(obj)
        return obj
}

func IsBoolean(obj Object) bool {
	if _, ok := obj.(*Boolean); ok {
		return true
	}
	return false
}

func IsFalse(obj Object) bool {
	return obj == The_False
}

func IsTrue(obj Object) bool {
	return !IsFalse(obj)
}

func IsSymbol(obj Object) bool {
	if _, ok := obj.(*Symbol); ok {
		return true;
	}
	return false
}

func IsEmptyList(obj Object) bool {
	if _, ok := obj.(*EmptyList); ok {
		return true
	}
	return false
}

func Cons(car Object, cdr Object) Object {
        res := &Pair{PAIR, car, cdr}
        return (Object)(res)
}

func IsPair(obj Object) bool {
        if _, ok := obj.(*Pair); ok {
                return true
        }
	return false
}

func Init() {
        Symbol_Table = make(map[string]Object)
        The_Empty_List = &EmptyList { EMPTY_LIST }
	The_True = &Boolean{ BOOLEAN, 1}
	The_False = &Boolean{ BOOLEAN, 0}
        Set_Symbol = MakeSymbol("set!")
        OK_Symbol = MakeSymbol("ok")
        If_Symbol = MakeSymbol("if")
	Else_Symbol = MakeSymbol("else")
	Let_Symbol = MakeSymbol("let")
	And_Symbol = MakeSymbol("and")
	Or_Symbol = MakeSymbol("or")
	Define_Symbol = MakeSymbol("define")
        Begin_Symbol = MakeSymbol("begin")
	Quote_Symbol = MakeSymbol("quote")
	Lambda_Symbol = MakeSymbol("lambda")
	Cond_Symbol = MakeSymbol("cond")
}
