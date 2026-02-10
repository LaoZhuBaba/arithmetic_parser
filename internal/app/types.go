package app

import (
	"fmt"
	"math"
)

type token int8
type precedence int8
type associativity int8

const (
	NullToken token = iota
	Number
	Plus
	Minus
	Multiply
	Divide
	LParen
	RParen
	Exponent
)

const (
	leftAssociative associativity = iota
	rightAssociative
)

// The order below is significant!
const (
	precedenceExponent precedence = iota
	precedenceMultiplyDivide
	precedencePlusMinus
)

var operations = map[token]struct {
	token       token
	description string
	fn          func(int, int) (int, error)
}{
	Plus:     {token: Plus, description: "Plus", fn: func(a, b int) (int, error) { return a + b, nil }},
	Minus:    {token: Minus, description: "Minus", fn: func(a, b int) (int, error) { return a - b, nil }},
	Multiply: {token: Multiply, description: "Multiply", fn: func(a, b int) (int, error) { return a * b, nil }},
	Divide: {token: Divide, description: "Divide", fn: func(a, b int) (int, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}},
	Exponent: {token: Exponent, description: "Exponent", fn: func(a, b int) (int, error) { return int(math.Pow(float64(a), float64(b))), nil }},
}

var opTokens = map[token]string{
	Plus:     "+",
	Minus:    "-",
	Multiply: "*",
	Divide:   "/",
	LParen:   "(",
	RParen:   ")",
	Exponent: "^",
}

type operationGroup struct {
	tokens        []token
	associativity associativity
}

var operationGroups = map[precedence]operationGroup{
	precedenceExponent:       {tokens: []token{Exponent}, associativity: rightAssociative},
	precedenceMultiplyDivide: {tokens: []token{Multiply, Divide}, associativity: leftAssociative},
	precedencePlusMinus:      {tokens: []token{Plus, Minus}, associativity: leftAssociative},
}

type Element struct {
	token      token
	tokenValue string
}
