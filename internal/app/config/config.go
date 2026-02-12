package config

import (
	"fmt"
	"math"

	"github.com/LaoZhuBaba/arithmetic_parser/internal/pkg/lexer"
	"github.com/LaoZhuBaba/arithmetic_parser/internal/pkg/parser"
)

var Tokens = []lexer.Token{
	{Id: lexer.Plus, Value: "+"},
	{Id: lexer.Minus, Value: "-"},
	{Id: lexer.Multiply, Value: "*"},
	{Id: lexer.Divide, Value: "/"},
	{Id: lexer.Exponent, Value: "^"},
	{Id: lexer.LParen, Value: "("},
	{Id: lexer.RParen, Value: ")"},
}

var Operations = []parser.Operation{
	{Description: "Plus", TokenId: lexer.Plus, Fn: func(a, b int) (int, error) { return a + b, nil }},
	{Description: "Minus", TokenId: lexer.Minus, Fn: func(a, b int) (int, error) { return a - b, nil }},
	{Description: "Multiply", TokenId: lexer.Multiply, Fn: func(a, b int) (int, error) { return a * b, nil }},
	{Description: "Divide", TokenId: lexer.Divide, Fn: func(a, b int) (int, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}

		return a / b, nil
	}},
	{Description: "Exponent", TokenId: lexer.Exponent, Fn: func(a, b int) (int, error) { return int(math.Pow(float64(a), float64(b))), nil }},
}

var OpGroup = []parser.OperationGroup{
	{Tokens: []lexer.TokenId{lexer.Exponent}, Precedence: parser.PrecedenceExponent, Associativity: parser.RightAssociative},
	{Tokens: []lexer.TokenId{lexer.Multiply, lexer.Divide}, Precedence: parser.PrecedenceMultiplyDivide, Associativity: parser.LeftAssociative},
	{Tokens: []lexer.TokenId{lexer.Plus, lexer.Minus}, Precedence: parser.PrecedencePlusMinus, Associativity: parser.LeftAssociative},
}
