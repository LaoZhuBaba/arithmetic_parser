package parser

import "errors"

var (
	errInvalidOperation = errors.New("invalid operation")
	errInvalidExpression = errors.New("invalid expression")
	errIndexOutOfRange = errors.New("index out of range")
	errInvalidTokenId   = errors.New("invalid TokenId")
)
