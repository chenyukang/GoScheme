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

func TestParser(t *testing.T) {
	Init()
	obj, err := parserWrapper("demo")
	if err != nil {
		t.Error(err)
	}
	if !(obj.Type == SYMBOL &&
		obj.Data.symbol == "demo") {
		t.Error("parser")
	}
}
