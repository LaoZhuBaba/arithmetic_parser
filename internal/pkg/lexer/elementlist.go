package lexer

import "fmt"

// FindLeftOperator returns the leftmost matching operator and its index from the elements slice or
// (NullToken, -1) if not found.  NullToken is a sentinel value, not an error.
func (el ElementList) FindLeftOperator(operators []TokenId) (TokenId, int) {
	for i, e := range el {
		for _, tok := range operators {
			if e.Token == tok {
				return tok, i
			}
		}
	}

	return NullToken, -1
}

// FindRightOperator returns the rightmost matching operator and its index from the elements slice or
// (NullToken, -1) if not found.  NullToken is a sentinel value, not an error.
func (el ElementList) FindRightOperator(operators []TokenId) (TokenId, int) {
	for i := range el {
		reverseIdx := len(el) - i - 1
		for _, tok := range operators {
			if el[reverseIdx].Token == tok {
				return tok, reverseIdx
			}
		}
	}

	return NullToken, -1
}

// FindLParen searches for the first left parenthesis in the slice and returns its index
// and true if found or -1 and false.
func (el ElementList) FindLParen() (int, bool) {
	for i, element := range el {
		if element.Token == LParen {
			return i, true
		}
	}

	return -1, false
}

// FindRParen searches for right parenthesis matching the left parenthesis which
// should always been found at index lParenIdx.
func (el ElementList) FindRParen(lParenIdx int) (rParenIdx int, err error) {
	// Use depth to track the number of unmatched parentheses
	var depth int

	if lParenIdx < 0 || lParenIdx >= len(el) {
		return 0, fmt.Errorf(
			"invalid index %d: out of range for elements slice",
			lParenIdx,
		)
	}

	if el[lParenIdx].Token != LParen {
		return 0, fmt.Errorf(
			"invalid TokenId at index %d: expected LParen, got %v",
			lParenIdx,
			el[lParenIdx].Token,
		)
	}

	for i := lParenIdx + 1; i < len(el); i++ {
		switch el[i].Token {
		case LParen:
			depth++
		case RParen:
			if depth == 0 {
				return i, nil
			}

			depth--
		default:
			continue
		}
	}

	return 0, fmt.Errorf("unmatched parentheses")
}
