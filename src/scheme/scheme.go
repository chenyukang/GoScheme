package main

import (
	"bufio"
	"eval"
	"flag"
	"fmt"
	"log"
	"os"
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

	eval.Init()

	reader := bufio.NewReader(file)
	iteractive := true
	if file != os.Stdin {
		iteractive = false
	}
	eval.Run(reader, iteractive)
}
