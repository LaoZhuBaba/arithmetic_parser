package app

import "fmt"

func (e Element) String() string {
	return e.tokenValue
}

// GetElements parses a string into a slice of Elements representing tokens
func GetElements(s string) (e []Element, err error) {
	var skip int

	e = []Element{}
	// Iterate over each rune in the string (not each byte)
outer:
	for idx, char := range s {
		if char == ' ' {
			continue
		}
		// Because we are iterating over a range, we can't increment idx within the for loop.
		// So instead, if a multiple rune element is processed, skip will be a positive value.
		if skip > 0 {
			skip--
			continue
		}

		for opTok, opChar := range opTokens {
			if char == opChar {
				e = append(e, Element{opTok, string(opChar)})
				// Once we've found and handled an operator that matches the rune char, we can skip to the next rune
				continue outer
			}
		}

		if char < '0' || char > '9' {
			return nil, fmt.Errorf("invalid character: %char", char)
		}

		remaining := s[idx:]

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

	return e, nil
}
