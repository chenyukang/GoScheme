package main

import (
	"fmt"
	"time"
)

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


var Symbol_Table   map[string]Object
var The_Empty_List Object;
var Set_Symbol     Object;
var OK_Symbol      Object;
var If_Symbol      Object;
var Begin_Symbol   Object;


func Make_Symbol(sym string) Object {
	if res, ok := Symbol_Table[sym]; ok {
		return res
	}
	obj := &Symbol{ STRING, sym}
	Symbol_Table[sym] = (Object)(obj)
	return obj
}

func IsEmptyList(obj Object) bool {
	if obj == The_Empty_List {
		return true
	} else {
		return false
	}
}

func Init() {
	Symbol_Table = make(map[string]Object)
	The_Empty_List = &EmptyList { EMPTY_LIST }
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
	fmt.Println(IsEmptyList(The_Empty_List))
	fmt.Println(IsEmptyList(OK_Symbol))
}
