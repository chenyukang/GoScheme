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
var Set_Symbol     Object
var OK_Symbol      Object
var If_Symbol      Object
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

func IsEmptyList(obj Object) bool {
	if obj == The_Empty_List {
		return true
	} else {
		return false
	}
}

func Cons(car Object, cdr Object) Object {
	res := &Pair{PAIR, car, cdr}
	return (Object)(res)
}

func IsPair(obj Object) bool {
	if _, ok := obj.(*Pair); ok {
		return true
	} else {
		return false
	}
}

func Init() {
	Symbol_Table = make(map[string]Object)
	The_Empty_List = &EmptyList { EMPTY_LIST }
	Set_Symbol = MakeSymbol("set")
	OK_Symbol = MakeSymbol("ok")
	If_Symbol = MakeSymbol("if")
	Begin_Symbol = MakeSymbol("begin")
}

func main() {
	var a byte = 'a'
	fmt.Println("byte:", a)
	fmt.Println("Welcome to GoScheme.\nUse Ctrl-C to exit.")
	Init()
	fmt.Printf("now val: %d\n", The_Empty_List)
	fmt.Println("time:", time.Now())
	fmt.Println(IsEmptyList(The_Empty_List))
	fmt.Println(IsEmptyList(OK_Symbol))
	res := Cons(OK_Symbol, If_Symbol)
	fmt.Println("res:", res)
	fmt.Println(IsPair(res))
}
