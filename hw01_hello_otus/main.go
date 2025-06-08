package main

import (
	"fmt"
	"golang.org/x/example/hello/reverse"
)

func main() {
	originalString := "Hello, OTUS!"

	fmt.Println(reverse.String(originalString))
}
