package lexer

import (
	"reflect"
	"testing"
)

func TestLexer_GetElementList(t *testing.T) {
	// Lexer needs the operator vocabulary to recognize non-number tokens.
	defaultTokens := []Token{
		{Id: Plus, Value: "+"},
		{Id: Minus, Value: "-"},
		{Id: Multiply, Value: "*"},
		{Id: Divide, Value: "/"},
		{Id: Exponent, Value: "^"},
		{Id: LParen, Value: "("},
		{Id: RParen, Value: ")"},
	}

	tests := []struct {
		name       string
		input      string
		expected   ElementList
		shouldFail bool
	}{
		{
			name:  "empty input",
			input: "",
			expected: ElementList{
				// empty
			},
			shouldFail: false,
		},
		{
			name:  "spaces only",
			input: "    ",
			expected: ElementList{
				// empty
			},
			shouldFail: false,
		},
		{
			name:  "single number",
			input: "42",
			expected: ElementList{
				{Token: Number, TokenValue: "42"},
			},
			shouldFail: false,
		},
		{
			name:  "single number with spaces",
			input: "  42   ",
			expected: ElementList{
				{Token: Number, TokenValue: "42"},
			},
			shouldFail: false,
		},
		{
			name:  "multiple numbers separated by spaces",
			input: " 3   5  7   ",
			expected: ElementList{
				{Token: Number, TokenValue: "3"},
				{Token: Number, TokenValue: "5"},
				{Token: Number, TokenValue: "7"},
			},
			shouldFail: false,
		},
		{
			name:  "numbers with parentheses and plus/minus",
			input: "(12+34)-56",
			expected: ElementList{
				{Token: LParen, TokenValue: "("},
				{Token: Number, TokenValue: "12"},
				{Token: Plus, TokenValue: "+"},
				{Token: Number, TokenValue: "34"},
				{Token: RParen, TokenValue: ")"},
				{Token: Minus, TokenValue: "-"},
				{Token: Number, TokenValue: "56"},
			},
			shouldFail: false,
		},
		{
			name:  "multiply and divide",
			input: "8/2*3",
			expected: ElementList{
				{Token: Number, TokenValue: "8"},
				{Token: Divide, TokenValue: "/"},
				{Token: Number, TokenValue: "2"},
				{Token: Multiply, TokenValue: "*"},
				{Token: Number, TokenValue: "3"},
			},
			shouldFail: false,
		},
		{
			name:  "exponent tokens",
			input: "2^3^2",
			expected: ElementList{
				{Token: Number, TokenValue: "2"},
				{Token: Exponent, TokenValue: "^"},
				{Token: Number, TokenValue: "3"},
				{Token: Exponent, TokenValue: "^"},
				{Token: Number, TokenValue: "2"},
			},
			shouldFail: false,
		},
		{
			name:       "invalid character",
			input:      "j",
			expected:   nil,
			shouldFail: true,
		},
		{
			name:       "invalid characters in expression",
			input:      "3+a",
			expected:   nil,
			shouldFail: true,
		},
		{
			name:       "invalid input with special characters",
			input:      "5+@",
			expected:   nil,
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := NewLexer(tt.input, defaultTokens)
			got, err := lex.GetElementList()

			if (err != nil) != tt.shouldFail {
				t.Fatalf("GetElementList() err=%v shouldFail=%v", err, tt.shouldFail)
			}
			if tt.shouldFail {
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("GetElementList() got=%#v want=%#v", got, tt.expected)
			}
		})
	}
}
