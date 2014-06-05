package eval

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

// write buf to a tmporary file,
// use bufio.Reader as a parameter for parser
func parserWrapper(buf string) (*Object, error) {
	fp, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	fp.WriteString(buf)
	fp.Close()
	path := fp.Name()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	return read(reader), nil
}

func TestParserSymbol(t *testing.T) {
	Init()
	obj, _ := parserWrapper("demo")
	if !(obj.Type == SYMBOL &&
		obj.Data.symbol == "demo") {
		t.Error("parser")
	}
}

func TestParserInt(t *testing.T) {
	Init()
	obj, _ := parserWrapper("10")
	if !(obj.Type == FIXNUM &&
		obj.Data.fixNum == 10) {
		t.Error("parser")
	}
}

func TestParserChar(t *testing.T) {
	Init()
	obj, _ := parserWrapper("#\\h")
	if !(obj.Type == CHARACTER &&
		obj.Data.char == 'h') {
		t.Error("parser character failed")
	}
}

func TestParserList(t *testing.T) {
	Init()
	obj, _ := parserWrapper("()")
	if obj != The_EmptyList {
		t.Error("parser list")
	}

	obj, _ = parserWrapper("(1)")
	if !(obj.Type == PAIR &&
		car(obj).Type == FIXNUM &&
		car(obj).Data.fixNum == 1 &&
		cdr(obj) == The_EmptyList) {
		t.Error("parser list")
	}

	obj, _ = parserWrapper("(1 2)")
	if !(obj.Type == PAIR &&
		cadr(obj).Type == FIXNUM &&
		cadr(obj).Data.fixNum == 2) {
		t.Error("parser list")
	}

	obj, _ = parserWrapper("(1 2 3)")
	obj = cdr(obj)
	obj = cdr(obj)
	obj = car(obj)
	if !(obj.Type == FIXNUM &&
		obj.Data.fixNum == 3) {
		t.Error("parser list")
	}

	obj, _ = parserWrapper("(#t #f)")
	tt := car(obj)
	ff := cadr(obj)
	if !(tt.Type == BOOLEAN && tt == The_True &&
		ff.Type == BOOLEAN && ff == The_False) {
		t.Error("parser list : boolean")
	}
}
