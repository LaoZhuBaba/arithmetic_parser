package parser

import (
	"github.com/LaoZhuBaba/arithmetic_parser/pkg/lexer"
)

type precedence int8
type associativity int8

// The following values define the order in which OperationGroups are evaluated.
// Order is significant because lower numbers are evaluated before higher numbers.
const (
	PrecedenceExponent precedence = iota
	PrecedenceMultiplyDivide
	PrecedencePlusMinus
)

type Parser struct {
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

// OperationGroup defines a group of Operations that share the same precedence.
// and associativity.  Each Operation is identified by a TokenId.
type OperationGroup struct {
	Tokens        []lexer.TokenId
	Associativity associativity
	Precedence    precedence
}
