package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chenyukang/GoScheme/eval"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: GoScheme <prog>\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	args := flag.Args()

	file := os.Stdin
	if len(args) >= 1 {
		fileName := args[0]
		fileReader, _err := os.Open(fileName)
		if _err != nil {
			log.Fatal(_err)
			os.Exit(1)
		}
		file = fileReader
	}

	fmt.Println("Welcome to GoScheme.\nUse Ctrl-C to exit.")
	eval.Init()

	reader := bufio.NewReader(file)
	for {
		byte, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("finished")
			} else {
				fmt.Println("error happened")
			}
			return
		}
		fmt.Printf("Got: %s", byte)
	}
}
