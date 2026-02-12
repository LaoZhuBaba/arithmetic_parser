package parser

import (
	"fmt"
	"strconv"

	"github.com/LaoZhuBaba/arithmetic_parser/pkg/lexer"
)

func NewParserOp(operations []Operation, opGroups []OperationGroup) (p parserOp) {
	p = parserOp{Operations: operations, OperationGroups: opGroups}
	return p
}

func (p parserOp) getOperationByTokenId(t lexer.TokenId) (*Operation, error) {
	for _, op := range p.Operations {
		if op.TokenId == t {
			return &op, nil
		}
	}

	return nil, fmt.Errorf("no operation defined with TokenId %d", t)
}

// Eval accepts a list of elements representing an arithmetic expression
// and returns the result as a pointer to an int.
func (p parserOp) Eval(e lexer.ElementList) (i *int, err error) {
	// Make a copy of the slice so we can modify it without affecting the original
	// Fuzz testing requires that the slice be immutable.
	elementList := make(lexer.ElementList, len(e))
	copy(elementList, e)

	// Evaluate parenthetical expressions first
	elementList, err = p.evalParen(elementList)
	if err != nil {
		return nil, err
	}

	// operationGroups is the map of precedence levels to the operators that can be used in that level
	// Each key is a precedence level that refers to a group of operators that have the same precedence
	// and associativity.  For example, multiplication and division share the same precedence level called
	// "multiplyDivide" and they are both left associative.  It is assumed that the first precedence
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

	result, err := strconv.Atoi(elementList[0].TokenValue)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// evalParen reduces parenthetical expressions to numbers, calling Eval() for subexpressions
func (p parserOp) evalParen(elementList lexer.ElementList) (lexer.ElementList, error) {
	// Iterate until every parenthetical expression has been reduced to a number
	for {
		// It's not an error to find no left parenthesis.  Unbalanced parentheses are handled in findRParen()
		lParenIdx, tf := elementList.FindLParen()
		if !tf {
			break
		}

		rParenIdx, err := elementList.FindRParen(lParenIdx)
		if err != nil {
			return nil, err
		}
		// Submit the expression inside the parentheses for evaluation
		val, err := p.Eval(elementList[lParenIdx+1 : rParenIdx])
		if err != nil {
			return nil, err
		}
		// Build a number element with the result of the evaluation
		newElement := lexer.ElementList{{Token: lexer.Number, TokenValue: fmt.Sprintf("%d", *val)}}
		// Replace the parentheses with the evaluated expression
		elementList = append(
			elementList[:lParenIdx],
			append(newElement, elementList[rParenIdx+1:]...)...,
		)
	}

	return elementList, nil
}

// evalArithmetic reduces arithmetic expressions to a single number element, calling Eval() for subexpressions.
func (p parserOp) evalArithmetic(elementList lexer.ElementList, precedence precedence) (lexer.ElementList, error) {
	for {
		var exprVal, lVal, rVal, idx int

		var tok lexer.TokenId

		switch p.OperationGroups[precedence].Associativity {
		case RightAssociative:
			// Get the index of the next operator and the TokenId
			tok, idx = elementList.FindRightOperator(p.OperationGroups[precedence].Tokens)

		case LeftAssociative:
			// Get the index of the next operator and the TokenId
			tok, idx = elementList.FindLeftOperator(p.OperationGroups[precedence].Tokens)
		}

		if tok == lexer.NullToken {
			break
		}

		// Get the elements that make up the expression: [number, operator, number]
		subExpr, err := p.getOperatorElements(idx, elementList)
		if err != nil {
			return nil, err
		}

		lVal, err = strconv.Atoi(subExpr[0].TokenValue)
		if err != nil {
			return nil, err
		}

		rVal, err = strconv.Atoi(subExpr[2].TokenValue)
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
		remainder := make(lexer.ElementList, len(elementList[idx+2:]))
		copy(remainder, elementList[idx+2:])
		// idx-1 is the index of the left operand so we are taking everything before the
		// expression and appending the result of the expression.
		elementList = append(elementList[:idx-1], lexer.Element{Token: lexer.Number, TokenValue: fmt.Sprintf("%d", exprVal)})
		elementList = append(elementList, remainder...)
	}

	return elementList, nil
}

// getOperatorElements returns the elements that make up an operator expression: [number, operator, number]
func (p parserOp) getOperatorElements(idx int, elementList lexer.ElementList) (subExp lexer.ElementList, err error) {
	// elements[idx] should be an operator TokenId so there must be a character before and after it
	if idx < 1 || idx >= len(elementList)-1 {
		return nil, fmt.Errorf("out of range with index %d and elements: %v", idx, elementList)
	}

	tok := elementList[idx].Token

	op, err := p.getOperationByTokenId(tok)
	if err != nil {
		return nil, fmt.Errorf("invalid TokenId: %v", tok)
	}

	switch {
	case elementList[idx-1].Token != lexer.Number:
		return nil, fmt.Errorf(
			"invalid TokenId before %s: expected Number, got %v", op.Description, elementList[idx-1].Token)
	case elementList[idx+1].Token != lexer.Number:
		return nil, fmt.Errorf(
			"invalid TokenId after %s: expected Number, got %v", op.Description, elementList[idx+1].Token)
	default:
		return elementList[idx-1 : idx+2], nil
	}
}
