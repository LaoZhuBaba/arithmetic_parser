package app

import (
	"fmt"
	"math"
)

func Calculate(s string) error {
	tokens := []Token{
		{Id: Plus, Value: "+"},
		{Id: Minus, Value: "-"},
		{Id: Multiply, Value: "*"},
		{Id: Divide, Value: "/"},
		{Id: Exponent, Value: "^"},
		{Id: LParen, Value: "("},
		{Id: RParen, Value: ")"},
	}
	expression := NewLexer(s, tokens)

	elements, err := expression.GetElementList()
	if err != nil {
		fmt.Println(err)
		return err
	}

	operations := []Operation{
		{Description: "Plus", TokenId: Plus, Fn: func(a, b int) (int, error) { return a + b, nil }},
		{Description: "Minus", TokenId: Minus, Fn: func(a, b int) (int, error) { return a - b, nil }},
		{Description: "Multiply", TokenId: Multiply, Fn: func(a, b int) (int, error) { return a * b, nil }},
		{Description: "Divide", TokenId: Divide, Fn: func(a, b int) (int, error) {
			if b == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return a / b, nil
		}},
		{Description: "Exponent", TokenId: Exponent, Fn: func(a, b int) (int, error) { return int(math.Pow(float64(a), float64(b))), nil }},
	}
	opGroup := []OperationGroup{
		{Tokens: []TokenId{Exponent}, Precedence: PrecedenceExponent, Associativity: RightAssociative},
		{Tokens: []TokenId{Multiply, Divide}, Precedence: PrecedenceMultiplyDivide, Associativity: LeftAssociative},
		{Tokens: []TokenId{Plus, Minus}, Precedence: PrecedencePlusMinus, Associativity: LeftAssociative},
	}
	parser := NewParser(operations, opGroup)

	result, err := parser.Eval(elements)
	if err != nil {
		return err
	}
	fmt.Println(*result)
	return nil
}
