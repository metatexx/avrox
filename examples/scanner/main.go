package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: scanner <path>")
		os.Exit(0)
	}
	Scanner(os.Args[1])
}
