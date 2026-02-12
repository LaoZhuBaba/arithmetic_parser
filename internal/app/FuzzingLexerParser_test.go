package app

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
	"testing"

	lexer "github.com/LaoZhuBaba/arithmetic_parser/pkg/lexer"
	parser "github.com/LaoZhuBaba/arithmetic_parser/pkg/parser"
)

// FuzzLexAndEval fuzzes the lexer+parser pipeline.
//
// Properties checked:
//  1. No panics for any input
//  2. If Lexer.GetElementList() succeeds, then Lexer.GetElementList(normalize(tokens)) succeeds and yields identical Elements
//  3. parser.Eval is deterministic for the same token stream (including right-associative operators like exponent)
func FuzzLexAndEval(f *testing.F) {
	// Lexer needs the operator vocabulary to recognize non-number Tokens.
	defaultTokens := []lexer.Token{
		{Id: lexer.Plus, Value: "+"},
		{Id: lexer.Minus, Value: "-"},
		{Id: lexer.Multiply, Value: "*"},
		{Id: lexer.Divide, Value: "/"},
		{Id: lexer.Exponent, Value: "^"},
		{Id: lexer.LParen, Value: "("},
		{Id: lexer.RParen, Value: ")"},
	}

	operations := []parser.Operation{
		{Description: "Plus", TokenId: lexer.Plus, Fn: func(a, b int) (int, error) { return a + b, nil }},
		{Description: "Minus", TokenId: lexer.Minus, Fn: func(a, b int) (int, error) { return a - b, nil }},
		{Description: "Multiply", TokenId: lexer.Multiply, Fn: func(a, b int) (int, error) { return a * b, nil }},
		{Description: "Divide", TokenId: lexer.Divide, Fn: func(a, b int) (int, error) {
			if b == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return a / b, nil
		}},
		{Description: "Exponent", TokenId: lexer.Exponent, Fn: func(a, b int) (int, error) { return int(math.Pow(float64(a), float64(b))), nil }},
	}
	opGroups := []parser.OperationGroup{
		{Tokens: []lexer.TokenId{lexer.Exponent}, Precedence: parser.PrecedenceExponent, Associativity: parser.RightAssociative},
		{Tokens: []lexer.TokenId{lexer.Multiply, lexer.Divide}, Precedence: parser.PrecedenceMultiplyDivide, Associativity: parser.LeftAssociative},
		{Tokens: []lexer.TokenId{lexer.Plus, lexer.Minus}, Precedence: parser.PrecedencePlusMinus, Associativity: parser.LeftAssociative},
	}
	p := parser.NewParserOp(operations, opGroups)

	// Seed corpus: valid, invalid, whitespacey, precedence, parentheses, associativity, etc.
	seeds := []string{
		"",
		"   ",
		"42",
		"2+3*4",
		"(2+3)*4",
		"((1+2)-(3-4))",
		"7/2",
		"1+",
		"+1",
		"()",
		"(1+2",
		"1+2)",
		"1/0",         // should not panic (even if it errors)
		"99999999999", // may overflow Atoi -> should error, not panic
		"3+a",         // invalid char
		"2^3",         // exponent operator
		"2^3^2",       // right associativity: 2^(3^2)
		"(2^3)^2",     // parentheses override associativity
		"2^(3^2)",     // explicit right association
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		// Keep fuzzing fast and avoid pathological memory use.
		if len(s) > 512 {
			return
		}

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic for input %q: %v", s, r)
			}
		}()

		lex := lexer.NewLexer(s, defaultTokens)
		elems, err := lex.GetElementList()
		if err != nil {
			// For arbitrary strings, lexer errors are expected.
			return
		}

		// Normalize: lexer ignores spaces, so rebuild expression from Tokens,
		// adding some random whitespace between Tokens.
		var b strings.Builder
		for _, e := range elems {
			randomInt := rand.IntN(30)
			b.WriteString(e.TokenValue)
			if randomInt < 10 {
				for i := 0; i < randomInt; i++ {
					b.WriteRune(' ')
				}
			}
		}
		normalized := b.String()

		lex2 := lexer.NewLexer(normalized, defaultTokens)
		elems2, err := lex2.GetElementList()
		if err != nil {
			t.Fatalf("GetElementList failed after normalization. input=%q normalized=%q err=%v", s, normalized, err)
		}

		// Elements should match exactly after normalization.
		if len(elems) != len(elems2) {
			t.Fatalf("token count changed after normalization. input=%q normalized=%q got=%d got2=%d",
				s, normalized, len(elems), len(elems2))
		}
		for i := range elems {
			if elems[i] != elems2[i] {
				t.Fatalf("Tokens changed after normalization at i=%d. input=%q normalized=%q got=%v got2=%v",
					i, s, normalized, elems[i], elems2[i])
			}
		}

		// p.Eval should not panic; it may legitimately return an error.
		// Use the normalized token stream for determinism checks too.
		got1, err1 := p.Eval(elems2)
		got2, err2 := p.Eval(elems2)

		// Determinism check: same Tokens, same outcome class (err vs ok) and same value if ok.
		if (err1 != nil) != (err2 != nil) {
			t.Fatalf("nondeterministic error behavior. input=%q normalized=%q err1=%v err2=%v", s, normalized, err1, err2)
		}
		if err1 == nil {
			if got1 == nil || got2 == nil {
				t.Fatalf("Eval returned nil result without error. input=%q normalized=%q got1=%v got2=%v", s, normalized, got1, got2)
			}
			if *got1 != *got2 {
				t.Fatalf("nondeterministic result. input=%q normalized=%q got1=%d got2=%d", s, normalized, *got1, *got2)
			}
		}
	})
}
