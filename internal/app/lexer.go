package app

import "fmt"

type token int

const (
	NullToken token = iota
	Number
	Plus
	Minus
	Multiply
	Divide
	LParen
	RParen
)

var operators = map[token]string{
	Plus:     "Plus",
	Minus:    "Minus",
	Multiply: "Multiply",
	Divide:   "Divide",
}

type Element struct {
	token      token
	tokenValue string
}

func (e Element) String() string {
	return fmt.Sprintf("%s", e.tokenValue)
}

// GetElements parses a string into a slice of Elements representing tokens
func GetElements(s string) (e []Element, err error) {
	var skip int
	e = []Element{}
	// Iterate over each rune in the string (not each byte)
	for i, c := range s {
		if c == ' ' {
			continue
		}
		// Because we are iterating over a range we can't increment i within the for loop.
		// So instead, if a multiple rune element is processed, we set skip to a positive value
		if skip > 0 {
			skip--
			continue
		}
		switch c {
		case '+':
			e = append(e, Element{Plus, string(c)})
		case '-':
			e = append(e, Element{Minus, string(c)})
		case '*':
			e = append(e, Element{Multiply, string(c)})
		case '/':
			e = append(e, Element{Divide, string(c)})
		case '(':
			e = append(e, Element{LParen, string(c)})
		case ')':
			e = append(e, Element{RParen, string(c)})
		// The only other valid token is a number which may be multiple runes
		default:
			if c < '0' || c > '9' {
				return nil, fmt.Errorf("invalid character: %c", c)
			}
			remaining := s[i:]
			var numStr string
			for _, c2 := range remaining {
				if c2 < '0' || c2 > '9' {
					break
				}
				numStr += string(c2)
			}
			// The range index naturally increments by one rune on each iteration,
			// so we only need to skip the number by which numStr exceeds 1.
			skip = len(numStr) - 1
			e = append(e, Element{Number, numStr})
		}
	}
	return e, nil
}
