package eval

import (
	"bufio"
	"fmt"
	"os"
)

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

// EOF is also a delimiter
func isDelimiter(val byte) bool {
	if isSpace(val) || val == '(' || val == ')' ||
		val == '"' || val == ';' || val == 0 {
		return true
	} else {
		return false
	}
}

func peekc(reader *bufio.Reader) byte {
	c, err := reader.Peek(1)
	if err != nil {
		//EOF
		return 0
	}
	return c[0]
}

func readc(reader *bufio.Reader) byte {
	c, _ := reader.ReadByte()
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

func readChar(reader *bufio.Reader) Object {
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

func readPair(reader *bufio.Reader) Object {
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

func read(reader *bufio.Reader) Object {
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
			return makeInt(int64(num))
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
		return makeStr(buf)
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
			ungetc(reader)
			return makeSymbol(buf)
		}
	} else if c == '\'' {
		return cons(Quote_Symbol, cons(read(reader), The_EmptyList))
	} else if c == '(' {
		return readPair(reader)
	}
	return nil
}

func writePair(obj Object) {
	carObj := car(obj)
	cdrObj := cdr(obj)
	write(carObj)
	if typeOf(cdrObj) == PAIR {
		fmt.Fprintf(os.Stderr, " ")
		writePair(cdrObj)
	} else if typeOf(cdrObj) == EMPTY_LIST {
		return
	} else {
		fmt.Fprintf(os.Stderr, " . ")
		write(cdrObj)
	}
}

func write(obj Object) {
	_type := typeOf(obj)
	switch _type {
	case BOOLEAN:
		if isFalse(obj) {
			fmt.Fprintf(os.Stderr, "#f")
		} else {
			fmt.Fprintf(os.Stderr, "#t")
		}
	case EMPTY_LIST:
		fmt.Fprintf(os.Stderr, "()")
	case SYMBOL:
		fmt.Fprintf(os.Stderr, "%s", asSym(obj))
	case CHARACTER:
		c := asChar(obj)
		fmt.Fprintf(os.Stderr, "#\\")
		switch c {
		case '\n':
			fmt.Fprintf(os.Stderr, "newline")
		case ' ':
			fmt.Fprintf(os.Stderr, "space")
		default:
			fmt.Fprintf(os.Stderr, "%c", c)
		}
	case INT:
		fmt.Fprintf(os.Stderr, "%d", asInt(obj))
	case STRING:
		fmt.Fprintf(os.Stderr, "\"%s\"", asStr(obj))
	case PAIR:
		fmt.Fprintf(os.Stderr, "(")
		writePair(obj)
		fmt.Fprintf(os.Stderr, ")")
	case PRIMITIVE_PROC:
		fmt.Fprintf(os.Stderr, "#<primitive-procedure>")
	case COMPOUND_PROC:
		fmt.Fprintf(os.Stderr, "#<compound-proc>")
	default:
		fmt.Println(obj)
	}
}
