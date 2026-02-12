package app

import (
	"fmt"

	"github.com/LaoZhuBaba/arithmetic_parser/internal/app/config"
	"github.com/LaoZhuBaba/arithmetic_parser/internal/pkg/lexer"
	"github.com/LaoZhuBaba/arithmetic_parser/internal/pkg/parser"
)

func Calculate(s string) error {

	lx := lexer.NewLexer(s, config.Tokens)

	elements, err := lx.GetElementList()
	if err != nil {
		fmt.Println(err)
		return err
	}
	operations := config.Operations
	pa := parser.NewParser(operations, config.OpGroup)

	result, err := pa.Eval(elements)
	if err != nil {
		return err
	}

	fmt.Println(*result)

	return nil
}
