package main

import (
	"fmt"
)

func Reverse(str string) string {
	reversed := []rune(str)
	stringLength := len(str)

	for i, letter := range str {
		reversed[stringLength - i - 1] = letter
	}

	return string(reversed)
}

func main() {
	startString := "Hello, OTUS!"

	fmt.Println(Reverse(startString))
}
