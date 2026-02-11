package app

import (
	"fmt"
	"strconv"
)

func NewParser(operations []Operation, opGroups []OperationGroup) (parser Parser) {
	parser = Parser{Operations: operations, OperationGroups: opGroups}
	return parser
}

func (p Parser) getOperationByTokenId(t TokenId) (*Operation, error) {
	for _, op := range p.Operations {
		if op.TokenId == t {
			return &op, nil
		}
	}
	return nil, fmt.Errorf("no operation defined with TokenId %d", t)
}

// Eval accepts a list of elements representing an arithmetic expression
// and returns the result as a pointer to an int.
func (p Parser) Eval(e ElementList) (i *int, err error) {
	// Make a copy of the slice so we can modify it without affecting the original
	// Fuzz testing requires that the slice be immutable.
	elementList := make(ElementList, len(e))
	copy(elementList, e)

	// Evaluate parenthetical expressions first
	elementList, err = p.evalParen(elementList)
	if err != nil {
		return nil, err
	}

	// operationGroups is the map of Precedence levels to the operators that can be used in that level
	// Each key is a Precedence level that refers to a group of operators that have the same Precedence
	// and Associativity.  For example, multiplication and division share the same Precedence level called
	// "multiplyDivide" and they are both left associative.  It is assumed that the first Precedence
	// to be evaluated will have a value of 0 and others will have consecutive values.
	for _, group := range p.OperationGroups {
		elementList, err = p.evalArithmetic(elementList, group.Precedence)
		if err != nil {
			return nil, err
		}
	}

	// After all operation groups have been processed, there should only be one element left: the result
	if len(elementList) != 1 {
		return nil, fmt.Errorf("invalid expression: %v", elementList)
	}

	result, err := strconv.Atoi(elementList[0].tokenValue)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// evalParen reduces parenthetical expressions to numbers, calling Eval() for subexpressions
func (p Parser) evalParen(elementList ElementList) (ElementList, error) {
	// Iterate until every parenthetical expression has been reduced to a number
	for {
		// It's not an error to find no left parenthesis.  Unbalanced parentheses are handled in findRParen()
		lParenIdx, tf := elementList.findLParen()
		if !tf {
			break
		}

		rParenIdx, err := elementList.findRParen(lParenIdx)
		if err != nil {
			return nil, err
		}
		// Submit the expression inside the parentheses for evaluation
		val, err := p.Eval(elementList[lParenIdx+1 : rParenIdx])
		if err != nil {
			return nil, err
		}
		// Build a number element with the result of the evaluation
		newElement := ElementList{{token: Number, tokenValue: fmt.Sprintf("%d", *val)}}
		// Replace the parentheses with the evaluated expression
		elementList = append(
			elementList[:lParenIdx],
			append(newElement, elementList[rParenIdx+1:]...)...,
		)
	}
	return elementList, nil
}

// evalArithmetic reduces arithmetic expressions to a single number element, calling Eval() for subexpressions.
func (p Parser) evalArithmetic(elementList ElementList, precedence precedence) (ElementList, error) {
	for {
		var exprVal, lVal, rVal, idx int

		var tok TokenId

		switch p.OperationGroups[precedence].Associativity {
		case RightAssociative:
			// Get the index of the next operator and the TokenId
			tok, idx = elementList.findRightOperator(p.OperationGroups[precedence].Tokens)

		case LeftAssociative:
			// Get the index of the next operator and the TokenId
			tok, idx = elementList.findLeftOperator(p.OperationGroups[precedence].Tokens)
		}

		if tok == NullToken {
			break
		}

		// Get the elements that make up the expression: [number, operator, number]
		subExpr, err := p.getOperatorElements(idx, elementList)
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

		op, err := p.getOperationByTokenId(tok)
		if err != nil {
			return nil, err
		}
		exprVal, err = op.Fn(lVal, rVal)
		if err != nil {
			return nil, err
		}

		// remainder is the slice of elements after the expression being evaluated.
		// idx+2 is used because idx is the index of the operator TokenId, so idx+1 is
		// the second number TokenId, and idx+2 is the start of anything that follows.
		remainder := make(ElementList, len(elementList[idx+2:]))
		copy(remainder, elementList[idx+2:])
		// idx-1 is the index of the left operand so we are taking everything before the
		// expression and appending the result of the expression.
		elementList = append(elementList[:idx-1], Element{token: Number, tokenValue: fmt.Sprintf("%d", exprVal)})
		elementList = append(elementList, remainder...)
	}

	return elementList, nil
}

// getOperatorElements returns the elements that make up an operator expression: [number, operator, number]
func (p Parser) getOperatorElements(idx int, elementList ElementList) (subExp ElementList, err error) {
	// elements[idx] should be an operator TokenId so there must be a character before and after it
	if idx < 1 || idx >= len(elementList)-1 {
		return nil, fmt.Errorf("out of range with index %d and elements: %v", idx, elementList)
	}

	tok := elementList[idx].token
	op, err := p.getOperationByTokenId(tok)

	if err != nil {
		return nil, fmt.Errorf("invalid TokenId: %v", tok)
	}

	switch {
	case elementList[idx-1].token != Number:
		return nil, fmt.Errorf(
			"invalid TokenId before %s: expected Number, got %v", op.Description, elementList[idx-1].token)
	case elementList[idx+1].token != Number:
		return nil, fmt.Errorf(
			"invalid TokenId after %s: expected Number, got %v", op.Description, elementList[idx+1].token)
	default:
		return elementList[idx-1 : idx+2], nil
	}
}

// findLeftOperator returns the leftmost matching operator and its index from the elements slice or
// (NullToken, -1) if not found.  NullToken is a sentinel value, not an error.
func (e ElementList) findLeftOperator(operators []TokenId) (TokenId, int) {
	for i, e := range e {
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
func (el ElementList) findRightOperator(operators []TokenId) (TokenId, int) {
	for i := range el {
		reverseIdx := len(el) - i - 1
		for _, tok := range operators {
			if el[reverseIdx].token == tok {
				return tok, reverseIdx
			}
		}
	}
	return NullToken, -1
}

// findLParen searches for the first left parenthesis in the slice and returns its index
// and true if found or -1 and false.
func (el ElementList) findLParen() (int, bool) {
	for i, element := range el {
		if element.token == LParen {
			return i, true
		}
	}

	return -1, false
}

// findRParen searches for right parenthesis matching the left parenthesis which
// should always been found at index lParenIdx.
func (el ElementList) findRParen(lParenIdx int) (rParenIdx int, err error) {
	// Use depth to track the number of unmatched parentheses
	var depth int

	if lParenIdx < 0 || lParenIdx >= len(el) {
		return 0, fmt.Errorf(
			"invalid index %d: out of range for elements slice",
			lParenIdx,
		)
	}

	if el[lParenIdx].token != LParen {
		return 0, fmt.Errorf(
			"invalid TokenId at index %d: expected LParen, got %v",
			lParenIdx,
			el[lParenIdx].token,
		)
	}

	for i := lParenIdx + 1; i < len(el); i++ {
		switch el[i].token {
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

func (e Element) String() string {
	return e.tokenValue
}
