package app

import (
	"fmt"
	"strconv"
)

type precedence int

const (
	precedenceMultiplyDivide precedence = iota
	precedencePlusMinus
)

// evalParen reduces parenthetical expressions to numbers, calling Eval() for subexpressions
func evalParen(elements []Element) ([]Element, error) {
	// Iterate until every parenthetical expression has been reduced to a number
	for {
		// It's not an error to find no left parenthesis.  Unbalanced parentheses are handled in findRParen()
		lParenIdx, tf := findLParen(elements)
		if tf == false {
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

// evalArithmetic reduces arithmetic expressions to a single number element, calling Eval() for
// subexpressions.  This implementation could easily be extended to support any left associate
// infix operator.  E.g., '10 - 3 - 4' is evaluated as (10 - 3) - 4.  For a right associative
// operator you could enhance the findNextOperator() function to take an associativity parameter.
func evalArithmetic(elements []Element, precedence precedence) ([]Element, error) {
	// We need the outer label because switch/case has its own break scope.
outer:
	for {
		var exprVal, lVal, rVal, idx int
		var tok token

		// Get the index of the next multiply or divide operator and the token
		switch precedence {
		case precedenceMultiplyDivide:
			tok, idx = findNextOperator(elements, []token{Multiply, Divide})
			if tok == NullToken {
				break outer
			}
		case precedencePlusMinus:
			tok, idx = findNextOperator(elements, []token{Plus, Minus})
			if tok == NullToken {
				break outer
			}
		}
		// Get the elements that make up the expression: [number, operator, number]
		subExpr, err := findOperatorElements(tok, elements)
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
		switch precedence {
		case precedenceMultiplyDivide:
			switch tok {
			case Multiply:
				exprVal = lVal * rVal
			case Divide:
				if rVal == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				exprVal = lVal / rVal
			default:
				return nil, fmt.Errorf("invalid operator in multiply/divide expression: %v", tok)
			}
		case precedencePlusMinus:
			switch tok {
			case Plus:
				exprVal = lVal + rVal
			case Minus:
				exprVal = lVal - rVal
			default:
				return nil, fmt.Errorf("invalid operator in plus/minus expression: %v", tok)
			}
		default:
			return nil, fmt.Errorf("invalid precedence passed to evalArithmetic(): %v", tok)
		}
		// remainder is the slice of elements after the expression being evaluated.
		// idx+2 is used because idx is the index of the operator token so idx+1 is
		// the second number token and idx+2 is the start of anything that follows.
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

	// Next comes multiply & divide, in left to right order
	elem, err = evalArithmetic(elem, precedenceMultiplyDivide)
	if err != nil {
		return nil, err
	}
	// Now comes plus & minus in left to right order
	elem, err = evalArithmetic(elem, precedencePlusMinus)
	if err != nil {
		return nil, err
	}
	// After that, there should only be one element left: the result
	if len(elem) != 1 {
		return nil, fmt.Errorf("invalid expression: %v", elem)
	}
	result, err := strconv.Atoi(elem[0].tokenValue)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// findNextOperator returns the first matching operator and its index from the elements slice or
// (NullToken, -1) if not found.  NullToken is a sentinel value not an error.
func findNextOperator(elements []Element, operators []token) (token, int) {
	for i, e := range elements {
		for _, tok := range operators {
			if e.token == tok {
				return tok, i
			}
		}
	}
	return NullToken, -1
}

// findOperatorElements returns the elements that make up an arithmetic expression, including the
// operator token.
func findOperatorElements(tok token, elements []Element) (subExpr []Element, err error) {
	operationStr, ok := operators[tok]
	if !ok {
		return nil, fmt.Errorf("invalid token: %v", tok)
	}
	for i, e := range elements {
		if e.token == tok {
			switch {
			case i == 0:
				return nil, fmt.Errorf(
					"expression cannot begin with %s", operationStr)
			case i >= len(elements)-1:
				return nil, fmt.Errorf(
					"expression cannot end with %s", operationStr)
			case elements[i-1].token != Number:
				return nil, fmt.Errorf(
					"invalid token before %s: expected Number, got %v", operationStr, elements[i-1].token)
			case elements[i+1].token != Number:
				return nil, fmt.Errorf(
					"invalid token after %s: expected Number, got %v", operationStr, elements[i+1].token)
			default:
				return elements[i-1 : i+2], nil
			}
		}
	}
	return nil, fmt.Errorf("operator: %s not found", operationStr)
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
		if elements[i].token == LParen {
			depth++
		} else if elements[i].token == RParen {
			if depth == 0 {
				return i, nil
			}
			depth--
		}
	}
	return 0, fmt.Errorf("unmatched parentheses")
}
