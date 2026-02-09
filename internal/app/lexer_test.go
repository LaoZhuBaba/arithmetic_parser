package app

import (
	"reflect"
	"testing"
)

func TestGetElements(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   []Element
		shouldFail bool
	}{
		{
			name:  "single number",
			input: "42",
			expected: []Element{
				{token: Number, tokenValue: "42"},
			},
			shouldFail: false,
		},
		{
			name:  "single number initial spaces",
			input: "  42",
			expected: []Element{
				{token: Number, tokenValue: "42"},
			},
			shouldFail: false,
		},
		{
			name:  "single number with spaces",
			input: "  42   ",
			expected: []Element{
				{token: Number, tokenValue: "42"},
			},
			shouldFail: false,
		},
		{
			name:  "multiple numbers separated by spaces",
			input: " 3   5  7   ",
			expected: []Element{
				{token: Number, tokenValue: "3"},
				{token: Number, tokenValue: "5"},
				{token: Number, tokenValue: "7"},
			},
			shouldFail: false,
		},
		{
			name:  "numbers with parentheses and plus/minus",
			input: "(12+34)-56",
			expected: []Element{
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "12"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "34"},
				{token: RParen, tokenValue: ")"},
				{token: Minus, tokenValue: "-"},
				{token: Number, tokenValue: "56"},
			},
			shouldFail: false,
		},
		{
			name:  "multiply and divide",
			input: "8/2*3",
			expected: []Element{
				{token: Number, tokenValue: "8"},
				{token: Divide, tokenValue: "/"},
				{token: Number, tokenValue: "2"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "3"},
			},
			shouldFail: false,
		},
		{
			name:  "all operators with spaces",
			input: "  9 + 6 - 3 * 2 / 4  ",
			expected: []Element{
				{token: Number, tokenValue: "9"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "6"},
				{token: Minus, tokenValue: "-"},
				{token: Number, tokenValue: "3"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "2"},
				{token: Divide, tokenValue: "/"},
				{token: Number, tokenValue: "4"},
			},
			shouldFail: false,
		},
		{
			name:  "nested parentheses",
			input: "((1+ 2)- (3-4)) ",
			expected: []Element{
				{token: LParen, tokenValue: "("},
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
				{token: RParen, tokenValue: ")"},
				{token: Minus, tokenValue: "-"},
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "3"},
				{token: Minus, tokenValue: "-"},
				{token: Number, tokenValue: "4"},
				{token: RParen, tokenValue: ")"},
				{token: RParen, tokenValue: ")"},
			},
			shouldFail: false,
		},
		{
			name:  "long number sequences",
			input: "123456+7890",
			expected: []Element{
				{token: Number, tokenValue: "123456"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "7890"},
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
			name:       "invalid character leading space",
			input:      " j",
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
		{
			name:       "empty input",
			input:      "",
			expected:   []Element{},
			shouldFail: false,
		},
		{
			name:       "spaces only",
			input:      "    ",
			expected:   []Element{},
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetElements(tt.input)

			if !tt.shouldFail && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.shouldFail && err == nil {
				t.Fatalf("expected error but got none")
			}
			if !tt.shouldFail && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unexpected result, got: %#v, want: %#v", result, tt.expected)
			}
		})
	}
}
