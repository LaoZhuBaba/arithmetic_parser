package app

import (
	"fmt"

	"github.com/LaoZhuBaba/arithmetic_parser/internal/app/config"
	"github.com/LaoZhuBaba/arithmetic_parser/pkg/lexer"
	"github.com/LaoZhuBaba/arithmetic_parser/pkg/parser"
)

func Calculate(s string) error {
	lx := lexer.NewLexer(s, config.Tokens)

	elements, err := lx.GetElementList()
	if err != nil {
		fmt.Printf("lexer failed with error: %v", err)

		return err
	}

	pa := parser.NewParserOp(config.Operations, config.OpGroup)

	result, err := pa.Eval(elements)
	if err != nil {
		return err
	}

	fmt.Println(*result)

	return nil
}
