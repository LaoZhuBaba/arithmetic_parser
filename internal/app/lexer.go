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

		// Handle operator characters
		if _, ok := opTokens[char]; ok {
			e = append(e, Element{opTokens[char], string(char)})
			continue
		}

		if char < '0' || char > '9' {
			return nil, fmt.Errorf("invalid operator character: %c", char)
		}

		// Handle number characters
		var numStr string

		remaining := s[idx:]
		for _, remChar := range remaining {
			if remChar < '0' || remChar > '9' {
				break
			}

			numStr += string(remChar)
		}
		// The range index naturally increments by one rune on each iteration,
		// so we only need to skip the number by which numStr exceeds 1.
		skip = len(numStr) - 1
		e = append(e, Element{Number, numStr})
	}

	return e, nil
}
