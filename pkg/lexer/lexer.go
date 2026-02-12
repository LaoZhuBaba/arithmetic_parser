package lexer

import (
	"fmt"
)

func NewLexer(input string, tokens []Token) (expr Lexer) {
	expr.Input = input

	expr.tokens = map[string]TokenId{}
	for _, token := range tokens {
		expr.tokens[token.Value] = token.Id
	}

	return expr
}

func (e Element) String() string {
	return e.TokenValue
}

// GetElementList parses a string into a slice of Elements representing Tokens
func (l Lexer) GetElementList() (elementList ElementList, err error) {
	var skip int

	elementList = ElementList{}
	// Iterate over each rune in the string (not each byte)
	for idx, c := range l.Input {
		if c == ' ' {
			continue
		}
		// Because we are iterating over a range, we can't increment idx within the for loop.
		// So instead, if a multiple rune element is processed, skip will be a positive value.
		if skip > 0 {
			skip--
			continue
		}
		// Handle operator characters
		if _, ok := l.tokens[string(c)]; ok {
			elementList = append(elementList, Element{l.tokens[string(c)], string(c)})
			continue
		}

		if c < '0' || c > '9' {
			return nil, fmt.Errorf("invalid operator character: %c", c)
		}

		// Handle number characters
		var numStr string

		remaining := l.Input[idx:]
		for _, remChar := range remaining {
			if remChar < '0' || remChar > '9' {
				break
			}

			numStr += string(remChar)
		}
		// The range index naturally increments by one rune on each iteration,
		// so we only need to skip the number by which numStr exceeds 1.
		skip = len(numStr) - 1
		elementList = append(elementList, Element{Number, numStr})
	}

	return elementList, nil
}
