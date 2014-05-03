package main

import (
	"fmt"
	"time"
)

type Object interface {}

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

type Object struct {
	objType ObjType
}

type Boolean struct {
	Object
	value int
}

type Symbol struct {
	Object
	value string
}

type FixNum struct {
	Object
	value   int64
}

var Symbol_Table   map[string]*Symbol
var The_Empty_List *Object;
var Set_Symbol     *Symbol;
var OK_Symbol      *Symbol;
var If_Symbol      *Symbol;
var Begin_Symbol   *Symbol;

func Alloc_Object() *Object {
	obj := &Object{ objType: EMPTY_LIST }
	return obj
}

func Make_Symbol(sym string) *Symbol {
	if res, ok := Symbol_Table[sym]; ok {
		return res
	}
	obj := &Symbol{ Object{STRING}, sym}
	Symbol_Table[sym] = obj
	return obj
}

func Init() {
	Symbol_Table = make(map[string]*Symbol)
	The_Empty_List = Alloc_Object()
	Set_Symbol = Make_Symbol("set")
	OK_Symbol = Make_Symbol("ok")
	If_Symbol = Make_Symbol("if")
	Begin_Symbol = Make_Symbol("begin")
}

func main() {
	fmt.Println("Welcome to GoScheme.\nUse Ctrl-C to exit.")
	Init()
	fmt.Printf("now val: %d\n", The_Empty_List)
	fmt.Println("time:", time.Now())
	a := Alloc_Object()
	b := Alloc_Object()
	fmt.Println("a:", a)
	fmt.Println("b:", b)
}
