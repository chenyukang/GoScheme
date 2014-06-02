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

func isDelimiter(val byte) bool {
	if isSpace(val) || val == '(' || val == ')' ||
		val == '"' || val == ';' {
		return true
	} else {
		return false
	}
}

func peekc(reader *bufio.Reader) byte {
	c, err := reader.Peek(1)
	if err != nil {
		panic(err)
	}
	return c[0]
}

func readc(reader *bufio.Reader) byte {
	c, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	return c
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
	c, err := reader.ReadByte()
	if err != nil {
		panic("incomplete char literal\n")
	}
	if !isDelimiter(peekc(reader)) {
		panic("character not followed by delimiter\n")
	}
	return makeChar(c)
}

func isDigit(val byte) bool {
	if val >= '0' && val <= '9' {
		return true
	} else {
		return false
	}
}

func isAlpha(val byte) bool {
	if (val >= 'a' && val <= 'z') ||
		(val >= 'A' && val <= 'Z') {
		return true
	} else {
		return false
	}
}

func isInitial(val byte) bool {
	if isAlpha(val) ||
		val == '*' || val == '/' ||
		val == '+' || val == '-' ||
		val == '>' || val == '<' ||
		val == '=' || val == '?' ||
		val == '!' {
		return true
	} else {
		return false
	}
}

func readPair(reader *bufio.Reader) *Object {
	eatWhiteSpace(reader)
	c := readc(reader)
	if c == ')' {
		return The_EmptyList
	}
	reader.UnreadByte()
	carObj := read(reader)
	eatWhiteSpace(reader)
	c = readc(reader)
	if c == '.' {
	} else {
		reader.UnreadByte()
		cdrObj := readPair(reader)
		return cons(carObj, cdrObj)
	}
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
	} else if isDigit(c) || (c == '-' && (isDigit(peekc(reader)))) {
		//make a number
		sign := 1
		if c == '-' {
			sign = -1
		} else {
			reader.UnreadByte()
		}
		num := 0
		n := c
		for {
			n = readc(reader)
			if !isDigit(n) {
				break
			}
			num = (num * 10) + (int(n) - '0')
		}
		num *= sign
		if isDelimiter(n) {
			reader.UnreadByte()
			return makeFixNum(num)
		} else {
			panic("number not followed by delimiter\n")
		}
	} else if c == '"' { //read a string
		buf := ""
		for {
			n := readc(reader)
			if n == '"' {
				break
			}
			buf += string(n)
		}
		return makeString(buf)
	} else if isInitial(c) {
		n := c
		buf := string(c)
		for {
			n = readc(reader)
			if !(isInitial(n) || isDigit(n)) {
				break
			}
			buf += string(n)

		}
		if isDelimiter(n) {
			return makeSymbol(buf)
		}
	} else if c == '\'' {
		return cons(Quote_Symbol, cons(read(reader), The_EmptyList))
	} else if c == '(' {
		return readPair(reader)
	}
	return nil
}

func eval(exp *Object, env *Object) *Object {
	return exp
}

func writePair(obj *Object) {
	carObj := car(obj)
	cdrObj := cdr(obj)
	write(carObj)
	if cdrObj.Type == PAIR {
		fmt.Fprintf(os.Stderr, " ")
		writePair(cdrObj)
	} else if cdrObj.Type == EMPTY_LIST {
		return
	} else {
		fmt.Fprintf(os.Stderr, " . ")
		write(cdrObj)
	}
}

func write(obj *Object) {
	switch obj.Type {
	case BOOLEAN:
		if isFalse(obj) {
			fmt.Fprintf(os.Stderr, "#f")
		} else {
			fmt.Fprintf(os.Stderr, "#t")
		}
	case EMPTY_LIST:
		fmt.Fprintf(os.Stderr, "()")
	case SYMBOL:
		fmt.Fprintf(os.Stderr, "%s", obj.Data.symbol)
	case CHARACTER:
		c := obj.Data.char
		fmt.Fprintf(os.Stderr, "#\\")
		switch c {
		case '\n':
			fmt.Fprintf(os.Stderr, "newline")
		case ' ':
			fmt.Fprintf(os.Stderr, "space")
		default:
			fmt.Fprintf(os.Stderr, "%c", c)
		}
	case FIXNUM:
		fmt.Fprintf(os.Stderr, "%d", obj.Data.fixNum)
	case STRING:
		fmt.Fprintf(os.Stderr, "\"%s\"", obj.Data.str)
	case PAIR:
		fmt.Fprintf(os.Stderr, "(")
		writePair(obj)
		fmt.Fprintf(os.Stderr, ")")
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
