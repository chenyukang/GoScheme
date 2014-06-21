package eval

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func isTagged(exp Object, tag Object) bool {
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

func isApp(exp Object) bool {
	return isPair(exp)
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
	return isSymbol(exp)
}

func isQuoted(exp Object) bool {
	return isTagged(exp, Quote_Symbol)
}

func isAssign(exp Object) bool {
	return isTagged(exp, Set_Symbol)
}

func isDef(exp Object) bool {
	return isTagged(exp, Define_Symbol)
}

func isAnd(exp Object) bool {
	return isTagged(exp, And_Symbol)
}

func isOr(exp Object) bool {
	return isTagged(exp, Or_Symbol)
}

func isIf(exp Object) bool {
	return isTagged(exp, If_Symbol)
}

func isLambda(exp Object) bool {
	return isTagged(exp, Lambda_Symbol)
}

func isCond(exp Object) bool {
	return isTagged(exp, Cond_Symbol)
}

func isLet(exp Object) bool {
	return isTagged(exp, Let_Symbol)
}

func isBegin(exp Object) bool {
	return isTagged(exp, Begin_Symbol)
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
		return nil, err
	}
	_val, err = defVal(exp)
	if err != nil {
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

func evalBegin(exp Object, env Object) (Object, error) {
	exp = cdr(exp)
	for !isLast(exp) {
		eval(car(exp), env)
		exp = cdr(exp)
	}
	return eval(car(exp), env)
}

func evalApp(exp Object, env Object) (Object, error) {
	proc, err := eval(car(exp), env)
	if err != nil {
		return nil, err
	}
	args := listValues(cdr(exp), env)
	if isPrimitiveProc(proc) {
		p := asFunc(proc)
		return p(args)
	} else if isCompProc(proc) {
		_env := extendEnv(
			procParams(proc),
			args,
			procEnv(proc))
		return eval(makeBegin(procBody(proc)), _env)
	}

	return nil, errors.New("not implemented")
}

func letBindings(exp Object) Object {
	return cadr(exp)
}

func letBody(exp Object) Object {
	return cddr(exp)
}

func bindingParam(bind Object) Object {
	return car(bind)
}

func bindingArgu(bind Object) Object {
	return cadr(bind)
}

func bindingParams(binds Object) Object {
	if isEmptyList(binds) {
		return The_EmptyList
	} else {
		return cons(bindingParam(car(binds)),
			bindingParams(cdr(binds)))
	}
}

func bindingArgus(binds Object) Object {
	if isEmptyList(binds) {
		return The_EmptyList
	} else {
		return cons(bindingArgu(car(binds)),
			bindingArgus(cdr(binds)))
	}
}

func letParams(exp Object) Object {
	return bindingParams(letBindings(exp))
}

func letArgus(exp Object) Object {
	return bindingArgus(letBindings(exp))
}

func makeApp(operator Object, operands Object) Object {
	return cons(operator, operands)
}

func makeBegin(exp Object) Object {
	return cons(Begin_Symbol, exp)
}

func makeLambda(params Object, body Object) Object {
	return cons(Lambda_Symbol, cons(params, body))
}

func evalLet(exp Object, env Object) (Object, error) {
	obj := makeApp(
		makeLambda(letParams(exp), letBody(exp)),
		letArgus(exp))
	return eval(obj, env)
}

func lambdaParams(exp Object) Object {
	return cadr(exp)
}

func lambdaBody(exp Object) Object {
	return cddr(exp)
}

func evalLambda(exp Object, env Object) (Object, error) {
	return makeCompProc(
		lambdaParams(exp),
		lambdaBody(exp),
		env), nil
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
	} else if isLambda(exp) {
		return evalLambda(exp, env)
	} else if isBegin(exp) {
		return evalBegin(exp, env)
	} else if isAnd(exp) {
		return evalAnd(exp, env)
	} else if isOr(exp) {
		return evalOr(exp, env)
	} else if isCond(exp) {
		return evalCond(exp, env)
	} else if isLet(exp) {
		return evalLet(exp, env)
	} else if isApp(exp) {
		return evalApp(exp, env)
	} else {
		//fmt.Println("now:", exp)
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
