package main

import (
	"fmt"
	"os"
)

func main() {
	flags := os.Args
	fmt.Printf("args: %v\n", flags[1:])
	fmt.Println("Hello World :D")
}
