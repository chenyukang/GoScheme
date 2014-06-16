package eval

import (
	"errors"
	"reflect"
)

const (
	UNDEF = iota
	EMPTY_LIST
	BOOLEAN
	SYMBOL
	INT
	CHARACTER
	STRING
	PAIR
	PRIMITIVE_PROC
	COMPOUND_PROC
	INPUT_PORT
	OUTPUT_PORT
	EOF_OBJECT
)

type ObjType int
type Object interface{}
type ObjFun func(args Object) (Object, error)

var SymbolTable map[string]Object
var The_EmptyList, The_True, The_False Object
var The_Empty_Env, The_Global_Env Object
var Set_Symbol, OK_Symbol, FAIL_Symbol, If_Symbol, Else_Symbol Object
var Cond_Symbol, Or_Symbol, And_Symbol, Quote_Symbol, Lambda_Symbol Object
var Define_Symbol, Begin_Symbol, Let_Symbol, Eof_Symbol Object

//Use reflect ...

type Int struct {
	Type  int
	Value int64
}

type Str struct {
	Type  int
	Value string
}

type Bool struct {
	Type  int
	Value int
}

type Pair struct {
	Type int
	Car  Object
	Cdr  Object
}

type Symbol struct {
	Type  int
	Value string
}

type Char struct {
	Type  int
	Value byte
}

type Proc struct {
	Type  int
	Value ObjFun
}

type EmptyList struct {
	Type int
}

func typeOf(obj Object) int {
	s := reflect.ValueOf(obj).Elem()
	return int(s.FieldByName("Type").Int())
}

func makeSymbol(sym string) Object {
	if obj, ok := SymbolTable[sym]; ok {
		return obj
	}
	obj := &Symbol{SYMBOL, sym}
	SymbolTable[sym] = obj
	return obj
}

func makeEmptyList() Object {
	return &EmptyList{EMPTY_LIST}
}

func makeBoolean(val int) Object {
	return &Bool{BOOLEAN, val}
}

func makeInt(val int64) Object {
	return &Int{INT, val}
}

func makeStr(val string) Object {
	return &Str{STRING, val}
}

func makeChar(val byte) Object {
	return &Char{CHARACTER, val}
}

func makePrimitiveProc(fun ObjFun) Object {
	return &Proc{PRIMITIVE_PROC, fun}
}

func isSymbol(obj Object) bool {
	return typeOf(obj) == SYMBOL
}

func isStr(obj Object) bool {
	return typeOf(obj) == STRING
}

func isBool(obj Object) bool {
	return typeOf(obj) == BOOLEAN
}

func isFalse(obj Object) bool {
	return obj == The_False
}

func isTrue(obj Object) bool {
	return !isFalse(obj)
}

func isChar(obj Object) bool {
	return typeOf(obj) == CHARACTER
}

func isEmptyList(obj Object) bool {
	return typeOf(obj) == EMPTY_LIST
}

func isInt(obj Object) bool {
	return typeOf(obj) == INT
}

func isPair(obj Object) bool {
	return typeOf(obj) == PAIR
}

func isPrimitiveProc(obj Object) bool {
	return typeOf(obj) == PRIMITIVE_PROC
}

func fieldOf(obj Object, field string) Object {
	v := reflect.ValueOf(obj).Elem()
	return v.FieldByName(field).Interface()
}

func valueOf(m Object) reflect.Value {
	s := reflect.ValueOf(m).Elem()
	return s.FieldByName("Value")
}

func asInt(m Object) int64 {
	return valueOf(m).Int()
}

func asStr(m Object) string {
	return valueOf(m).String()
}

func asSym(m Object) string {
	return valueOf(m).String()
}

func asChar(m Object) byte {
	i := valueOf(m).Interface()
	return i.(byte)
}

func asFunc(m Object) ObjFun {
	i := valueOf(m).Interface()
	return i.(ObjFun)
}

func setField(obj Object, t Object, field string) {
	v := reflect.ValueOf(obj).Elem()
	v.FieldByName(field).Set(t.(reflect.Value))
}

func setCar(obj Object, t Object) {
	setField(obj, reflect.ValueOf(t), "Car")
}

func setCdr(obj Object, t Object) {
	setField(obj, reflect.ValueOf(t), "Cdr")
}

func car(obj Object) Object {
	return fieldOf(obj, "Car")
}

func cdr(obj Object) Object {
	return fieldOf(obj, "Cdr")
}

func cons(car Object, cdr Object) Object {
	return &Pair{PAIR, car, cdr}
}

func isLast(obj Object) bool {
	if isEmptyList(cdr(obj)) {
		return true
	} else {
		return false
	}
}

func makeFrame(vars Object, vals Object) Object {
	return cons(vars, vals)
}

func frameVars(frame Object) Object {
	return car(frame)
}

func frameVals(frame Object) Object {
	return cdr(frame)
}

func extendEnv(vars Object, vals Object, baseEnv Object) Object {
	return cons(makeFrame(vars, vals), baseEnv)
}

func addBinding(avar Object, aval Object, frame Object) {
	setCar(frame, cons(avar, car(frame)))
	setCdr(frame, cons(aval, cdr(frame)))
}

func firstFrame(env Object) Object {
	return car(env)
}

func lookupVar(avar Object, env Object) (Object, error) {
	e := env
	for !isEmptyList(e) {
		frame := firstFrame(env)
		vars := frameVars(frame)
		vals := frameVals(frame)
		for !isEmptyList(vars) {
			if avar == car(vars) {
				return car(vals), nil
			}
			vars = cdr(vars)
			vals = cdr(vals)
		}
		e = cdr(e)
	}
	return nil, errors.New("undef variable")
}

func defineVar(avar Object, aval Object, env Object) {
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
