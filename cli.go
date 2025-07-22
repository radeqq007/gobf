package main

import (
	"flag"
	"log"
)

var filepath string

func cli() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		log.Fatal("Provide a file path.")
	}

	filepath = args[0]
}
