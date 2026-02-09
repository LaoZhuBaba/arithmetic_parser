package main

import (
	"fmt"
	"os"

	"github.com/LaoZhuBaba/arithmetic_parser/internal/app"
)

var elements []app.Element

func main() {
	var expression string
	for _, arg := range os.Args[1:] {
		expression += arg
	}
	if expression == "" {
		fmt.Println("no expression provided")
		return
	}
	elements, err := app.GetElements(expression)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := app.Eval(elements)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("result: %d\n", *result)
}
