package parser

import (
	"github.com/LaoZhuBaba/arithmetic_parser/pkg/lexer"
)

type precedence int8
type associativity int8

// The following values define the order in which groups of operators are evaluated.
// The order is significant because lower numbers are evaluated before higher numbers.
const (
	PrecedenceExponent precedence = iota
	PrecedenceMultiplyDivide
	PrecedencePlusMinus
)

type parserOp struct {
	Operations      []Operation
	OperationGroups []OperationGroup
}
type Operation struct {
	Description string
	TokenId     lexer.TokenId
	Fn          lexer.OperationFn
}

const (
	LeftAssociative associativity = iota
	RightAssociative
)

type OperationGroup struct {
	Tokens        []lexer.TokenId
	Associativity associativity
	Precedence    precedence
}
