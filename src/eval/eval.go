package eval

import (
	"bufio"
	"fmt"
)

func isApp(exp *Object) bool {
	if isPair(exp) {
		return true
	} else {
		return false
	}
}

func isSelfEval(exp *Object) bool {
	if isBoolean(exp) ||
		isFixNum(exp) ||
		isChar(exp) ||
		isString(exp) {
		return true
	}
	return false
}

func isVariable(exp *Object) bool {
	if isSymbol(exp) {
		return true
	} else {
		return false
	}
}

func isTaggedWith(exp *Object, tag *Object) bool {
	if isPair(exp) {
		theCar := car(exp)
		if isSymbol(theCar) && (theCar == tag) {
			return true
		} else {
			return false
		}
	}
	return false
}

func isQuoted(exp *Object) bool {
	if isTaggedWith(exp, Quote_Symbol) {
		return true
	} else {
		return false
	}
}

func eval(exp *Object, env *Object) *Object {
	if isSelfEval(exp) {
		return exp
	} else if isVariable(exp) {
		return lookupVar(exp, env)
	} else if isQuoted(exp) {
		return cadr(exp)
	}
	return exp
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
