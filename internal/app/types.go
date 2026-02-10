package app

import (
	"fmt"
	"math"
)

type token int8
type precedence int8
type associativity int8

const (
	NullToken token = iota // NullToken is used as a sentinel
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
	description string
	fn          func(int, int) (int, error)
}{
	Plus:     {description: "Plus", fn: func(a, b int) (int, error) { return a + b, nil }},
	Minus:    {description: "Minus", fn: func(a, b int) (int, error) { return a - b, nil }},
	Multiply: {description: "Multiply", fn: func(a, b int) (int, error) { return a * b, nil }},
	Divide: {description: "Divide", fn: func(a, b int) (int, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}},
	Exponent: {description: "Exponent", fn: func(a, b int) (int, error) { return int(math.Pow(float64(a), float64(b))), nil }},
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

// Define operator precedence.  The tokens field is a list that defines operations that share the same precedence.
// Operations at the same parenthesis level that have the same precedence are evaluated from left to right.
// E.g., '48 / 3 / 8 / 2'  is evaluated as '((48 / 3) / 8) / 2'
var operationGroups = map[precedence]operationGroup{
	precedenceExponent:       {tokens: []token{Exponent}, associativity: rightAssociative},
	precedenceMultiplyDivide: {tokens: []token{Multiply, Divide}, associativity: leftAssociative},
	precedencePlusMinus:      {tokens: []token{Plus, Minus}, associativity: leftAssociative},
}

type Element struct {
	token      token
	tokenValue string
}
