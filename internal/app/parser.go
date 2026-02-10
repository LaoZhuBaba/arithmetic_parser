package app

import (
	"fmt"
	"slices"
	"strconv"
)

// evalParen reduces parenthetical expressions to numbers, calling Eval() for subexpressions
func evalParen(elements []Element) ([]Element, error) {
	// Iterate until every parenthetical expression has been reduced to a number
	for {
		// It's not an error to find no left parenthesis.  Unbalanced parentheses are handled in findRParen()
		lParenIdx, tf := findLParen(elements)
		if !tf {
			break
		}

		rParenIdx, err := findRParen(elements, lParenIdx)
		if err != nil {
			return nil, err
		}
		// Submit the expression inside the parentheses for evaluation
		val, err := Eval(elements[lParenIdx+1 : rParenIdx])
		if err != nil {
			return nil, err
		}
		// Build a number element with the result of the evaluation
		newElement := []Element{{token: Number, tokenValue: fmt.Sprintf("%d", *val)}}
		// Replace the parentheses with the evaluated expression
		elements = append(
			elements[:lParenIdx],
			append(newElement, elements[rParenIdx+1:]...)...,
		)
	}

	return elements, nil
}

// evalArithmetic reduces arithmetic expressions to a single number element, calling Eval() for subexpressions.
func evalArithmetic(elements []Element, precedence precedence) ([]Element, error) {
	for {
		var exprVal, lVal, rVal, idx int

		var tok token

		switch operationGroups[precedence].associativity {
		case rightAssociative:
			// Get the index of the next operator and the token
			tok, idx = findRightOperator(elements, operationGroups[precedence].tokens)

		case leftAssociative:
			// Get the index of the next operator and the token
			tok, idx = findLeftOperator(elements, operationGroups[precedence].tokens)
		}

		if tok == NullToken {
			break
		}

		// Get the elements that make up the expression: [number, operator, number]
		subExpr, err := getOperatorElements(idx, elements)
		if err != nil {
			return nil, err
		}

		lVal, err = strconv.Atoi(subExpr[0].tokenValue)
		if err != nil {
			return nil, err
		}

		rVal, err = strconv.Atoi(subExpr[2].tokenValue)
		if err != nil {
			return nil, err
		}

		if _, ok := operations[tok]; !ok {
			return nil, fmt.Errorf("invalid token: %v", tok)
		}

		exprVal, err = operations[tok].fn(lVal, rVal)
		if err != nil {
			return nil, err
		}

		// remainder is the slice of elements after the expression being evaluated.
		// idx+2 is used because idx is the index of the operator token, so idx+1 is
		// the second number token, and idx+2 is the start of anything that follows.
		remainder := make([]Element, len(elements[idx+2:]))
		copy(remainder, elements[idx+2:])
		// idx-1 is the index of the left operand so we are taking everything before the
		// expression and appending the result of the expression.
		elements = append(elements[:idx-1], Element{token: Number, tokenValue: fmt.Sprintf("%d", exprVal)})
		elements = append(elements, remainder...)
	}

	return elements, nil
}

// Eval accepts a list of elements representing an arithmetic expression
// and returns the result as a pointer to an int.
func Eval(e []Element) (*int, error) {
	// Make a copy of the slice so we can modify it without affecting the original
	// Fuzz testing requires that the slice be immutable.
	elements := make([]Element, len(e))
	copy(elements, e)

	// Evaluate parenthetical expressions first
	elem, err := evalParen(elements)
	if err != nil {
		return nil, err
	}

	// This is a bit awkward but we need the keys of operationGroups sorted in ascending order
	// to ensure precedence levels are followed correctly
	sortedPrecedence := make([]precedence, 0, len(operationGroups))
	for k := range operationGroups {
		sortedPrecedence = append(sortedPrecedence, k)
	}

	slices.Sort(sortedPrecedence)

	// operationGroups is the map of precedence levels to the operators that can be used in that level
	// Each key is a precedence level that refers to a group of operators that have the same precedence
	// and associativity.  For example, multiplication and division share the same precedence level called
	// "multiplyDivide" and they are both left associative.
	for _, opGroupKey := range sortedPrecedence {
		elem, err = evalArithmetic(elem, opGroupKey)
		if err != nil {
			return nil, err
		}
	}

	// After all operation groups have been processed, there should only be one element left: the result
	if len(elem) != 1 {
		return nil, fmt.Errorf("invalid expression: %v", elem)
	}

	result, err := strconv.Atoi(elem[0].tokenValue)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// findLeftOperator returns the leftmost matching operator and its index from the elements slice or
// (NullToken, -1) if not found.  NullToken is a sentinel value, not an error.
func findLeftOperator(elements []Element, operators []token) (token, int) {
	for i, e := range elements {
		for _, tok := range operators {
			if e.token == tok {
				return tok, i
			}
		}
	}

	return NullToken, -1
}

// findRightOperator returns the rightmost matching operator and its index from the elements slice or
// (NullToken, -1) if not found.  NullToken is a sentinel value, not an error.
func findRightOperator(elements []Element, operators []token) (token, int) {
	for i := range elements {
		reverseIdx := len(elements) - i - 1
		for _, tok := range operators {
			if elements[reverseIdx].token == tok {
				return tok, reverseIdx
			}
		}
	}

	return NullToken, -1
}

// getOperatorElements returns the elements that make up an operator expression: [number, operator, number]
func getOperatorElements(idx int, elements []Element) (subExpr []Element, err error) {
	// elements[idx] should be an operator token so there must be a character before and after it
	if idx < 1 || idx >= len(elements)-1 {
		return nil, fmt.Errorf("out of range with index %d and elements: %v", idx, elements)
	}

	tok := elements[idx].token
	op, ok := operations[tok]

	if !ok {
		return nil, fmt.Errorf("invalid token: %v", tok)
	}

	switch {
	case elements[idx-1].token != Number:
		return nil, fmt.Errorf(
			"invalid token before %s: expected Number, got %v", op.description, elements[idx-1].token)
	case elements[idx+1].token != Number:
		return nil, fmt.Errorf(
			"invalid token after %s: expected Number, got %v", op.description, elements[idx+1].token)
	default:
		return elements[idx-1 : idx+2], nil
	}
}

// findLParen searches for the first left parenthesis in the slice and returns its index
// and true if found or -1 and false.
func findLParen(elements []Element) (int, bool) {
	for i, e := range elements {
		if e.token == LParen {
			return i, true
		}
	}

	return -1, false
}

// findRParen searches for right parenthesis matching the left parenthesis which
// should always been found at index lParenIdx.
func findRParen(elements []Element, lParenIdx int) (rParenIdx int, err error) {
	// Use depth to track the number of unmatched parentheses
	var depth int

	if lParenIdx < 0 || lParenIdx >= len(elements) {
		return 0, fmt.Errorf(
			"invalid index %d: out of range for elements slice",
			lParenIdx,
		)
	}

	if elements[lParenIdx].token != LParen {
		return 0, fmt.Errorf(
			"invalid token at index %d: expected LParen, got %v",
			lParenIdx,
			elements[lParenIdx].token,
		)
	}

	for i := lParenIdx + 1; i < len(elements); i++ {
		switch elements[i].token {
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
