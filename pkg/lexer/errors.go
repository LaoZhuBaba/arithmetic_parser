package lexer

import "errors"

var (
	errInvalidOperator = errors.New("invalid operator character")
	errIndexOutOfRange = errors.New("index out of range")
	errInvalidTokenId = errors.New("invalid TokenId")
	errUnmatchedParen = errors.New("unmatched parenthesis")
)
