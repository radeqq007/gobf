package main

import (
	"flag"
)

var filepath string

func cli() {
	flag.Parse()
	var args []string = flag.Args()
	filepath = args[0]
}
