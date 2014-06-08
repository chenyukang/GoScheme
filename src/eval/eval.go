package eval

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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

func isAssign(exp *Object) bool {
	if isTaggedWith(exp, Set_Symbol) {
		return true
	} else {
		return false
	}
}

func isDef(exp *Object) bool {
	if isTaggedWith(exp, Define_Symbol) {
		return true
	} else {
		return false
	}
}

func defVar(exp *Object) (*Object, error) {
	if isSymbol(cadr(exp)) {
		return cadr(exp), nil
	}
	return nil, errors.New("defvar target is not symbol")
}

func defVal(exp *Object) (*Object, error) {
	if isSymbol(cadr(exp)) {
		left := cdr(cdr(exp))
		return car(left), nil
	}
	return nil, errors.New("defval failed")
}

func listValues(exp *Object, env *Object) *Object {
	if isEmptyList(exp) {
		return The_EmptyList
	} else {
		first, err := eval(car(exp), env)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return cons(first, listValues(cdr(exp), env))
	}
}

func evalAssign(exp *Object, env *Object) (*Object, error) {
	_var := cadr(exp)
	_val, err := eval(car(cdr(cdr(exp))), env)
	if err != nil {
		return nil, err
	}
	defineVar(_var, _val, env)
	return OK_Symbol, nil
}

func evalDef(exp *Object, env *Object) (*Object, error) {
	var _var, _val *Object
	var err error
	_var, err = defVar(exp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	_val, err = defVal(exp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	_val, _ = eval(_val, env)
	defineVar(_var, _val, env)
	return OK_Symbol, nil
}

func eval(exp *Object, env *Object) (*Object, error) {
	if isSelfEval(exp) {
		return exp, nil
	} else if isVariable(exp) {
		return lookupVar(exp, env)
	} else if isQuoted(exp) {
		return cadr(exp), nil
	} else if isAssign(exp) {
		return evalAssign(exp, env)
	} else if isDef(exp) {
		return evalDef(exp, env)
	} else if isApp(exp) {
		proc, err := eval(car(exp), env)
		if err != nil {
			return nil, err
		}
		args := listValues(cdr(exp), env)
		if isPrimitiveProc(proc) {
			val := proc.Data.primitive(args)
			return val, nil
		}
	}
	return exp, nil
}

func Run(reader *bufio.Reader) {
	fmt.Printf("Welcome to Bootstrap Scheme.\nUse ctrl-c to exit.\n")
	for {
		fmt.Printf("> ")
		exp := read(reader)
		if exp == nil {
			break
		}
		res, err := eval(exp, The_Global_Env)
		if err != nil {
			fmt.Println(err)
		} else {
			write(res)
		}
		fmt.Println()
	}
}
