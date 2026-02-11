package app

import (
	"fmt"
	"math"
	"testing"
)

func newTestParser() Parser {
	operations := []Operation{
		{Description: "Plus", TokenId: Plus, Fn: func(a, b int) (int, error) { return a + b, nil }},
		{Description: "Minus", TokenId: Minus, Fn: func(a, b int) (int, error) { return a - b, nil }},
		{Description: "Multiply", TokenId: Multiply, Fn: func(a, b int) (int, error) { return a * b, nil }},
		{Description: "Divide", TokenId: Divide, Fn: func(a, b int) (int, error) {
			if b == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return a / b, nil
		}},
		{Description: "Exponent", TokenId: Exponent, Fn: func(a, b int) (int, error) {
			return int(math.Pow(float64(a), float64(b))), nil
		}},
	}
	opGroups := []OperationGroup{
		{Tokens: []TokenId{Exponent}, Precedence: PrecedenceExponent, Associativity: RightAssociative},
		{Tokens: []TokenId{Multiply, Divide}, Precedence: PrecedenceMultiplyDivide, Associativity: LeftAssociative},
		{Tokens: []TokenId{Plus, Minus}, Precedence: PrecedencePlusMinus, Associativity: LeftAssociative},
	}
	return NewParser(operations, opGroups)
}

func ptrInt(v int) *int { return &v }

func TestElementList_findLParen(t *testing.T) {
	tests := []struct {
		name     string
		elements ElementList
		wantIdx  int
		wantBool bool
	}{
		{
			name:     "empty",
			elements: ElementList{},
			wantIdx:  -1,
			wantBool: false,
		},
		{
			name: "none",
			elements: ElementList{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			wantIdx:  -1,
			wantBool: false,
		},
		{
			name: "at start",
			elements: ElementList{
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "1"},
			},
			wantIdx:  0,
			wantBool: true,
		},
		{
			name: "later",
			elements: ElementList{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "2"},
			},
			wantIdx:  2,
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, gotBool := tt.elements.findLParen()
			if gotIdx != tt.wantIdx || gotBool != tt.wantBool {
				t.Fatalf("findLParen() = (%d,%v), want (%d,%v)", gotIdx, gotBool, tt.wantIdx, tt.wantBool)
			}
		})
	}
}

func TestElementList_findRParen(t *testing.T) {
	tests := []struct {
		name     string
		elements ElementList
		lIdx     int
		want     int
		wantErr  bool
	}{
		{
			name:     "simple",
			elements: ElementList{{token: LParen, tokenValue: "("}, {token: Number, tokenValue: "1"}, {token: RParen, tokenValue: ")"}},
			lIdx:     0,
			want:     2,
			wantErr:  false,
		},
		{
			name:     "nested outer",
			elements: ElementList{{token: LParen, tokenValue: "("}, {token: LParen, tokenValue: "("}, {token: Number, tokenValue: "1"}, {token: RParen, tokenValue: ")"}, {token: RParen, tokenValue: ")"}},
			lIdx:     0,
			want:     4,
			wantErr:  false,
		},
		{
			name:     "nested inner",
			elements: ElementList{{token: LParen, tokenValue: "("}, {token: LParen, tokenValue: "("}, {token: Number, tokenValue: "1"}, {token: RParen, tokenValue: ")"}, {token: RParen, tokenValue: ")"}},
			lIdx:     1,
			want:     3,
			wantErr:  false,
		},
		{
			name:     "error - out of range",
			elements: ElementList{{token: LParen, tokenValue: "("}, {token: RParen, tokenValue: ")"}},
			lIdx:     2,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "error - not LParen at index",
			elements: ElementList{{token: Number, tokenValue: "7"}, {token: LParen, tokenValue: "("}, {token: RParen, tokenValue: ")"}},
			lIdx:     0,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "error - unmatched",
			elements: ElementList{{token: LParen, tokenValue: "("}, {token: Number, tokenValue: "1"}},
			lIdx:     0,
			want:     0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.elements.findRParen(tt.lIdx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("findRParen() err=%v wantErr=%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("findRParen() got=%d want=%d", got, tt.want)
			}
		})
	}
}

func TestElementList_findLeftOperator(t *testing.T) {
	tests := []struct {
		name      string
		elements  ElementList
		operators []TokenId
		wantTok   TokenId
		wantIdx   int
	}{
		{
			name: "finds first multiply when searching multiply/divide",
			elements: ElementList{
				{token: Number, tokenValue: "2"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "3"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "4"},
			},
			operators: []TokenId{Multiply, Divide},
			wantTok:   Multiply,
			wantIdx:   3,
		},
		{
			name: "returns NullToken when no operator found",
			elements: ElementList{
				{token: Number, tokenValue: "2"},
				{token: Number, tokenValue: "3"},
			},
			operators: []TokenId{Plus, Minus, Multiply, Divide},
			wantTok:   NullToken,
			wantIdx:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTok, gotIdx := tt.elements.findLeftOperator(tt.operators)
			if gotTok != tt.wantTok || gotIdx != tt.wantIdx {
				t.Fatalf("findLeftOperator() = (%v,%d), want (%v,%d)", gotTok, gotIdx, tt.wantTok, tt.wantIdx)
			}
		})
	}
}

func TestElementList_findRightOperator(t *testing.T) {
	tests := []struct {
		name      string
		elements  ElementList
		operators []TokenId
		wantTok   TokenId
		wantIdx   int
	}{
		{
			name: "finds last exponent",
			elements: ElementList{
				{token: Number, tokenValue: "2"},
				{token: Exponent, tokenValue: "^"},
				{token: Number, tokenValue: "3"},
				{token: Exponent, tokenValue: "^"},
				{token: Number, tokenValue: "2"},
			},
			operators: []TokenId{Exponent},
			wantTok:   Exponent,
			wantIdx:   3,
		},
		{
			name: "returns NullToken when no operator found",
			elements: ElementList{
				{token: Number, tokenValue: "2"},
				{token: Number, tokenValue: "3"},
			},
			operators: []TokenId{Plus, Minus, Multiply, Divide, Exponent},
			wantTok:   NullToken,
			wantIdx:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTok, gotIdx := tt.elements.findRightOperator(tt.operators)
			if gotTok != tt.wantTok || gotIdx != tt.wantIdx {
				t.Fatalf("findRightOperator() = (%v,%d), want (%v,%d)", gotTok, gotIdx, tt.wantTok, tt.wantIdx)
			}
		})
	}
}

func TestParser_getOperatorElements(t *testing.T) {
	p := newTestParser()

	tests := []struct {
		name     string
		idx      int
		elements ElementList
		want     ElementList
		wantErr  bool
	}{
		{
			name: "plus window",
			idx:  1,
			elements: ElementList{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			want: ElementList{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			wantErr: false,
		},
		{
			name: "error - idx at beginning (needs left operand)",
			idx:  0,
			elements: ElementList{
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - right operand not Number",
			idx:  1,
			elements: ElementList{
				{token: Number, tokenValue: "8"},
				{token: Divide, tokenValue: "/"},
				{token: RParen, tokenValue: ")"},
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
		elements ElementList
		want     *int
		wantErr  bool
	}{
		{
			name:     "empty slice is error",
			elements: ElementList{},
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
			elements: ElementList{
				{token: Number, tokenValue: "42"},
			},
			want:    ptrInt(42),
			wantErr: false,
		},
		{
			name: "precedence: multiply before plus",
			elements: ElementList{
				{token: Number, tokenValue: "2"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "3"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "4"},
			},
			want:    ptrInt(14),
			wantErr: false,
		},
		{
			name: "right associativity: exponent",
			elements: ElementList{
				{token: Number, tokenValue: "2"},
				{token: Exponent, tokenValue: "^"},
				{token: Number, tokenValue: "3"},
				{token: Exponent, tokenValue: "^"},
				{token: Number, tokenValue: "2"},
			},
			want:    ptrInt(512), // 2^(3^2)
			wantErr: false,
		},
		{
			name: "parentheses override: (2^3)^2",
			elements: ElementList{
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "2"},
				{token: Exponent, tokenValue: "^"},
				{token: Number, tokenValue: "3"},
				{token: RParen, tokenValue: ")"},
				{token: Exponent, tokenValue: "^"},
				{token: Number, tokenValue: "2"},
			},
			want:    ptrInt(64),
			wantErr: false,
		},
		{
			name: "division by zero errors",
			elements: ElementList{
				{token: Number, tokenValue: "1"},
				{token: Divide, tokenValue: "/"},
				{token: Number, tokenValue: "0"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unmatched parentheses errors",
			elements: ElementList{
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
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
