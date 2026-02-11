package app

type TokenId int8
type precedence int8
type Associativity int8

type Token struct {
	Id    TokenId
	Value string
}

const (
	NullToken TokenId = iota // NullToken is used as a sentinel
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
	LeftAssociative Associativity = iota
	RightAssociative
)

// The following values define the order in which groups of operators are evaluated.
// The first defined value must be 0 and increase consecutively
const (
	PrecedenceExponent precedence = iota
	PrecedenceMultiplyDivide
	PrecedencePlusMinus
)

type OperationFn func(int, int) (int, error)
type Operation struct {
	Description string
	TokenId     TokenId
	Fn          OperationFn
}

type OperationGroup struct {
	Tokens        []TokenId
	Associativity Associativity
	Precedence    precedence
}

type Lexer struct {
	Input  string
	tokens map[string]TokenId
}

type Element struct {
	token      TokenId
	tokenValue string
}

type ElementList []Element

type Parser struct {
	Operations      []Operation
	OperationGroups []OperationGroup
}
