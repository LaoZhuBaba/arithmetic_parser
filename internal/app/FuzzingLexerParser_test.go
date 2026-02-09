package app

import (
	"strings"
	"testing"

	"math/rand/v2"
)

// FuzzLexAndEval fuzzes the lexer+parser pipeline.
//
// Properties checked:
//  1. No panics for any input
//  2. If GetElements(s) succeeds, then GetElements(normalize(s)) succeeds and yields identical tokens
//  3. If Eval succeeds, it is deterministic for the same token stream
func FuzzLexAndEval(f *testing.F) {
	// Seed corpus: valid, invalid, whitespacey, precedence, parentheses, etc.
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

		elems, err := GetElements(s)
		if err != nil {
			// For arbitrary bytes/strings, lexer errors are expected.
			return
		}

		// Normalize: lexer ignores spaces, so rebuild expression from tokens.
		var b strings.Builder
		for _, e := range elems {
			// Added randomInt by hand to introduce random space characters
			randomInt := rand.IntN(30)
			b.WriteString(e.tokenValue)
			if randomInt < 10 {
				for i := 0; i < randomInt; i++ {
					b.WriteRune(' ')
				}
			}
		}
		normalized := b.String()

		elems2, err := GetElements(normalized)
		if err != nil {
			t.Fatalf("GetElements failed after normalization. input=%q normalized=%q err=%v", s, normalized, err)
		}

		// Tokens should match exactly after normalization.
		if len(elems) != len(elems2) {
			t.Fatalf("token count changed after normalization. input=%q normalized=%q got=%d got2=%d",
				s, normalized, len(elems), len(elems2))
		}
		for i := range elems {
			if elems[i] != elems2[i] {
				t.Fatalf("tokens changed after normalization at i=%d. input=%q normalized=%q got=%v got2=%v",
					i, s, normalized, elems[i], elems2[i])
			}
		}

		// Eval should not panic; it may legitimately return an error.
		got1, err1 := Eval(elems)
		got2, err2 := Eval(elems)

		// Determinism check: same tokens, same outcome class (err vs ok) and same value if ok.
		if (err1 != nil) != (err2 != nil) {
			t.Fatalf("nondeterministic error behavior. input=%q err1=%v err2=%v", s, err1, err2)
		}
		if err1 == nil {
			if got1 == nil || got2 == nil {
				t.Fatalf("Eval returned nil result without error. input=%q got1=%v got2=%v", s, got1, got2)
			}
			if *got1 != *got2 {
				t.Fatalf("nondeterministic result. input=%q got1=%d got2=%d", s, *got1, *got2)
			}
		}
	})
}
