package parser

import (
	"fmt"
	"math"
	"testing"

	"github.com/LaoZhuBaba/arithmetic_parser/internal/pkg/lexer"
)

func newTestParser() parserOp {
	operations := []Operation{
		{Description: "Plus", TokenId: lexer.Plus, Fn: func(a, b int) (int, error) { return a + b, nil }},
		{Description: "Minus", TokenId: lexer.Minus, Fn: func(a, b int) (int, error) { return a - b, nil }},
		{Description: "Multiply", TokenId: lexer.Multiply, Fn: func(a, b int) (int, error) { return a * b, nil }},
		{Description: "Divide", TokenId: lexer.Divide, Fn: func(a, b int) (int, error) {
			if b == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return a / b, nil
		}},
		{Description: "Exponent", TokenId: lexer.Exponent, Fn: func(a, b int) (int, error) {
			return int(math.Pow(float64(a), float64(b))), nil
		}},
	}
	opGroups := []OperationGroup{
		{Tokens: []lexer.TokenId{lexer.Exponent}, Precedence: PrecedenceExponent, Associativity: RightAssociative},
		{Tokens: []lexer.TokenId{lexer.Multiply, lexer.Divide}, Precedence: PrecedenceMultiplyDivide, Associativity: LeftAssociative},
		{Tokens: []lexer.TokenId{lexer.Plus, lexer.Minus}, Precedence: PrecedencePlusMinus, Associativity: LeftAssociative},
	}
	return NewParser(operations, opGroups)
}

func ptrInt(v int) *int { return &v }

func TestParser_getOperatorElements(t *testing.T) {
	p := newTestParser()

	tests := []struct {
		name     string
		idx      int
		elements lexer.ElementList
		want     lexer.ElementList
		wantErr  bool
	}{
		{
			name: "plus window",
			idx:  1,
			elements: lexer.ElementList{
				{Token: lexer.Number, TokenValue: "1"},
				{Token: lexer.Plus, TokenValue: "+"},
				{Token: lexer.Number, TokenValue: "2"},
			},
			want: lexer.ElementList{
				{Token: lexer.Number, TokenValue: "1"},
				{Token: lexer.Plus, TokenValue: "+"},
				{Token: lexer.Number, TokenValue: "2"},
			},
			wantErr: false,
		},
		{
			name: "error - idx at beginning (needs left operand)",
			idx:  0,
			elements: lexer.ElementList{
				{Token: lexer.Plus, TokenValue: "+"},
				{Token: lexer.Number, TokenValue: "2"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - right operand not Number",
			idx:  1,
			elements: lexer.ElementList{
				{Token: lexer.Number, TokenValue: "8"},
				{Token: lexer.Divide, TokenValue: "/"},
				{Token: lexer.RParen, TokenValue: ")"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.getOperatorElements(tt.idx, tt.elements)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getOperatorElements() err=%v wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("getOperatorElements() len=%d want=%d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("getOperatorElements() got=%#v want=%#v", got, tt.want)
				}
			}
		})
	}
}

func TestParser_Eval(t *testing.T) {
	p := newTestParser()

	tests := []struct {
		name     string
		elements lexer.ElementList
		want     *int
		wantErr  bool
	}{
		{
			name:     "empty slice is error",
			elements: lexer.ElementList{},
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "nil slice is error",
			elements: nil,
			want:     nil,
			wantErr:  true,
		},
		{
			name: "single number",
			elements: lexer.ElementList{
				{Token: lexer.Number, TokenValue: "42"},
			},
			want:    ptrInt(42),
			wantErr: false,
		},
		{
			name: "precedence: multiply before plus",
			elements: lexer.ElementList{
				{Token: lexer.Number, TokenValue: "2"},
				{Token: lexer.Plus, TokenValue: "+"},
				{Token: lexer.Number, TokenValue: "3"},
				{Token: lexer.Multiply, TokenValue: "*"},
				{Token: lexer.Number, TokenValue: "4"},
			},
			want:    ptrInt(14),
			wantErr: false,
		},
		{
			name: "right associativity: exponent",
			elements: lexer.ElementList{
				{Token: lexer.Number, TokenValue: "2"},
				{Token: lexer.Exponent, TokenValue: "^"},
				{Token: lexer.Number, TokenValue: "3"},
				{Token: lexer.Exponent, TokenValue: "^"},
				{Token: lexer.Number, TokenValue: "2"},
			},
			want:    ptrInt(512), // 2^(3^2)
			wantErr: false,
		},
		{
			name: "parentheses override: (2^3)^2",
			elements: lexer.ElementList{
				{Token: lexer.LParen, TokenValue: "("},
				{Token: lexer.Number, TokenValue: "2"},
				{Token: lexer.Exponent, TokenValue: "^"},
				{Token: lexer.Number, TokenValue: "3"},
				{Token: lexer.RParen, TokenValue: ")"},
				{Token: lexer.Exponent, TokenValue: "^"},
				{Token: lexer.Number, TokenValue: "2"},
			},
			want:    ptrInt(64),
			wantErr: false,
		},
		{
			name: "division by zero errors",
			elements: lexer.ElementList{
				{Token: lexer.Number, TokenValue: "1"},
				{Token: lexer.Divide, TokenValue: "/"},
				{Token: lexer.Number, TokenValue: "0"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unmatched parentheses errors",
			elements: lexer.ElementList{
				{Token: lexer.LParen, TokenValue: "("},
				{Token: lexer.Number, TokenValue: "1"},
				{Token: lexer.Plus, TokenValue: "+"},
				{Token: lexer.Number, TokenValue: "2"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Eval(tt.elements)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Eval() err=%v wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got == nil || tt.want == nil || *got != *tt.want {
				t.Fatalf("Eval() got=%v want=%v", got, tt.want)
			}
		})
	}
}
