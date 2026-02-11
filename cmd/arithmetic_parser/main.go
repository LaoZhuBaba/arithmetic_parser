package main

import (
	"fmt"
	"os"

	"github.com/LaoZhuBaba/arithmetic_parser/internal/app"
)

func main() {
	var input string
	for _, arg := range os.Args[1:] {
		input += arg
	}

	if input == "" {
		fmt.Println("no expression provided")
		return
	}

	err := app.Calculate(input)
	if err != nil {
		fmt.Printf("calculation failed with error: %v", err)
		return
	}
}
