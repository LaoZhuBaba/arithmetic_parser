package lexer

type TokenId int8

type Token struct {
	Id    TokenId
	Value string
}

type OperationFn func(int, int) (int, error)

type Lexer struct {
	Input  string
	tokens map[string]TokenId
}

type Element struct {
	Token      TokenId
	TokenValue string
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

type ElementList []Element
