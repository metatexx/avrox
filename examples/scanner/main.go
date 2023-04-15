package main

import (
	"fmt"
	"github.com/metatexx/avrox"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: scanner <path>")
		os.Exit(0)
	}
	avrox.Scanner(os.Args[1])
}
