package eval

import (
	"bufio"
	"fmt"
	"os"
)

func makeEnv() *Object {
	env := extendEnv(The_EmptyList, The_EmptyList, The_Empty_Env)
	return env
}

func setupEnv(env *Object) {
	addProcedure("+", addProc, env)
	addProcedure("-", subProc, env)
	addProcedure("*", mulProc, env)
	addProcedure("/", divProc, env)
	addProcedure("=", isNumEqualProc, env)
	addProcedure("<", isLessProc, env)
	addProcedure(">", isLargerProc, env)
	addProcedure("null?", isNullProc, env)
	addProcedure("boolean?", isBooleanProc, env)
	addProcedure("integer?", isIntegerProc, env)
	addProcedure("symbol?", isSymbolProc, env)
	addProcedure("string?", isStringProc, env)
	addProcedure("pair?", isPairProc, env)

	addProcedure("cons", consProc, env)
	addProcedure("car", carProc, env)
	addProcedure("cdr", cdrProc, env)
	addProcedure("set-car!", setcarProc, env)
	addProcedure("set-cdr!", setcdrProc, env)

	addProcedure("error", errorProc, env)

}

func Init() {
	SymbolTable = make(map[string]*Object)
	The_EmptyList = allocObject()
	The_EmptyList.Type = EMPTY_LIST
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
	The_Empty_Env = The_EmptyList
	The_Global_Env = makeEnv()
	setupEnv(The_Global_Env)
}

func ungetc(reader *bufio.Reader) {
	err := reader.UnreadByte()
	if err != nil {
		panic(err)
	}
}

func isSpace(val byte) bool {
	if val == '\t' || val == '\n' || val == '\r' || val == ' ' {
		return true
	} else {
		return false
	}
}

func eatWhiteSpace(reader *bufio.Reader) {
	for {
		c, err := reader.ReadByte()
		if err != nil {
			break
		}
		if isSpace(c) {
			continue
		} else if c == ';' {
			for {
				v, err := reader.ReadByte()
				if err != nil || v == '\n' {
					break
				}
			}
			continue
		}
		ungetc(reader)
		break
	}

}

func readChar(reader *bufio.Reader) *Object {
	return nil
}

func read(reader *bufio.Reader) *Object {
	eatWhiteSpace(reader)
	c, _ := reader.ReadByte()
	if c == '#' {
		c, _ := reader.ReadByte()
		switch c {
		case 't':
			return The_True
		case 'f':
			return The_False
		case '\\':
			return readChar(reader)
		default:
			panic("unknown boolean or character literal\n")
		}
	}
	return nil
}

func eval(exp *Object, env *Object) *Object {
	return exp
}

func write(obj *Object) {
	switch obj.Type {
	case BOOLEAN:
		if isFalse(obj) {
			fmt.Fprintf(os.Stderr, "#f")
		} else {
			fmt.Fprintf(os.Stderr, "#t")
		}
	default:
		fmt.Println(obj)
	}
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
