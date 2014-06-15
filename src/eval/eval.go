package eval

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func isApp(exp Object) bool {
	if isPair(exp) {
		return true
	} else {
		return false
	}
}

func isSelfEval(exp Object) bool {
	if isBool(exp) ||
		isInt(exp) ||
		isChar(exp) ||
		isStr(exp) {
		return true
	}
	return false
}

func isVariable(exp Object) bool {
	if isSymbol(exp) {
		return true
	} else {
		return false
	}
}

func isTaggedWith(exp Object, tag Object) bool {
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

func isQuoted(exp Object) bool {
	if isTaggedWith(exp, Quote_Symbol) {
		return true
	} else {
		return false
	}
}

func isAssign(exp Object) bool {
	if isTaggedWith(exp, Set_Symbol) {
		return true
	} else {
		return false
	}
}

func isDef(exp Object) bool {
	if isTaggedWith(exp, Define_Symbol) {
		return true
	} else {
		return false
	}
}

func isAnd(exp Object) bool {
	if isTaggedWith(exp, And_Symbol) {
		return true
	} else {
		return false
	}
}

func isOr(exp Object) bool {
	if isTaggedWith(exp, Or_Symbol) {
		return true
	} else {
		return false
	}
}

func isCond(exp Object) bool {
	if isTaggedWith(exp, Cond_Symbol) {
		return true
	} else {
		return false
	}
}

func defVar(exp Object) (Object, error) {
	if isSymbol(cadr(exp)) {
		return cadr(exp), nil
	}
	return nil, errors.New("defvar target is not symbol")
}

func defVal(exp Object) (Object, error) {
	if isSymbol(cadr(exp)) {
		left := cdr(cdr(exp))
		return car(left), nil
	}
	return nil, errors.New("defval failed")
}

func isIf(exp Object) bool {
	if isTaggedWith(exp, If_Symbol) {
		return true
	} else {
		return false
	}
}

func listValues(exp Object, env Object) Object {
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

func evalAssign(exp Object, env Object) (Object, error) {
	_var := cadr(exp)
	_val, err := eval(car(cdr(cdr(exp))), env)
	if err != nil {
		return nil, err
	}
	defineVar(_var, _val, env)
	return OK_Symbol, nil
}

func evalDef(exp Object, env Object) (Object, error) {
	var _var, _val Object
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

func evalIf(exp Object, env Object) (Object, error) {
	pred := cadr(exp)
	ifConsT := cadr(cdr(exp))
	ifConsF := cdr(cdr(cdr(exp)))
	res, _ := eval(pred, env)
	if isTrue(res) {
		return eval(ifConsT, env)
	} else {
		if isEmptyList(ifConsF) {
			return The_False, nil
		} else {
			return eval(car(ifConsF), env)
		}
	}
}

func evalAnd(exp Object, env Object) (Object, error) {
	tests := cdr(exp)
	if isEmptyList(tests) {
		return The_True, nil
	}
	for {
		if isLast(tests) {
			break
		}
		res, _ := eval(car(tests), env)
		if isFalse(res) {
			return The_False, nil
		}
		tests = cdr(tests)
	}
	return eval(car(tests), env)
}

func evalOr(exp Object, env Object) (Object, error) {
	tests := cdr(exp)
	if isEmptyList(tests) {
		return The_True, nil
	}
	for {
		if isLast(tests) {
			break
		}
		res, _ := eval(car(tests), env)
		if isTrue(res) {
			return The_True, nil
		}
		tests = cdr(tests)
	}
	return eval(car(tests), env)
}

func evalCond(exp Object, env Object) (Object, error) {
	conds := cadr(exp)
	for {
		if isEmptyList(conds) {
			break
		}
		cur := car(conds)
		val, _ := eval(car(cur), env)
		if isTrue(val) || equal(val, Else_Symbol) {
			return eval(cadr(cur), env)
		}
		conds = cdr(conds)
	}
	return The_True, nil
}

func evalApp(exp Object, env Object) (Object, error) {
	proc, err := eval(car(exp), env)
	if err != nil {
		return nil, err
	}
	args := listValues(cdr(exp), env)
	if isPrimitiveProc(proc) {
		p := asFunc(proc)
		val := p(args)
		return val, nil
	}
	return nil, errors.New("not implemented")
}

func eval(exp Object, env Object) (Object, error) {
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
	} else if isIf(exp) {
		return evalIf(exp, env)
	} else if isAnd(exp) {
		return evalAnd(exp, env)
	} else if isOr(exp) {
		return evalOr(exp, env)
	} else if isCond(exp) {
		return evalCond(exp, env)
	} else if isApp(exp) {
		return evalApp(exp, env)
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
