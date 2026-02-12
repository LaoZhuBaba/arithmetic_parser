package lexer

import "testing"

func TestElementList_FindLeftOperator(t *testing.T) {
	tests := []struct {
		name      string
		elements  ElementList
		operators []TokenId
		wantTok   TokenId
		wantIdx   int
	}{
		{
			name:     "empty list => not found",
			elements: ElementList{},
			operators: []TokenId{
				Plus, Minus,
			},
			wantTok: NullToken,
			wantIdx: -1,
		},
		{
			name: "finds leftmost matching operator",
			elements: ElementList{
				{Token: Number, TokenValue: "2"},
				{Token: Plus, TokenValue: "+"},
				{Token: Number, TokenValue: "3"},
				{Token: Multiply, TokenValue: "*"},
				{Token: Number, TokenValue: "4"},
			},
			operators: []TokenId{Multiply, Plus},
			wantTok:   Plus,
			wantIdx:   1,
		},
		{
			name: "respects element scan order (not operator list order)",
			elements: ElementList{
				{Token: Number, TokenValue: "8"},
				{Token: Divide, TokenValue: "/"},
				{Token: Number, TokenValue: "2"},
				{Token: Multiply, TokenValue: "*"},
				{Token: Number, TokenValue: "3"},
			},
			operators: []TokenId{Multiply, Divide},
			wantTok:   Divide,
			wantIdx:   1,
		},
		{
			name: "no operator found => NullToken",
			elements: ElementList{
				{Token: Number, TokenValue: "1"},
				{Token: Number, TokenValue: "2"},
			},
			operators: []TokenId{Plus, Minus, Multiply, Divide},
			wantTok:   NullToken,
			wantIdx:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTok, gotIdx := tt.elements.FindLeftOperator(tt.operators)
			if gotTok != tt.wantTok || gotIdx != tt.wantIdx {
				t.Fatalf("FindLeftOperator() = (%v,%d), want (%v,%d)", gotTok, gotIdx, tt.wantTok, tt.wantIdx)
			}
		})
	}
}

func TestElementList_FindRightOperator(t *testing.T) {
	tests := []struct {
		name      string
		elements  ElementList
		operators []TokenId
		wantTok   TokenId
		wantIdx   int
	}{
		{
			name:     "empty list => not found",
			elements: ElementList{},
			operators: []TokenId{
				Exponent,
			},
			wantTok: NullToken,
			wantIdx: -1,
		},
		{
			name: "finds rightmost matching operator",
			elements: ElementList{
				{Token: Number, TokenValue: "2"},
				{Token: Exponent, TokenValue: "^"},
				{Token: Number, TokenValue: "3"},
				{Token: Exponent, TokenValue: "^"},
				{Token: Number, TokenValue: "2"},
			},
			operators: []TokenId{Exponent},
			wantTok:   Exponent,
			wantIdx:   3,
		},
		{
			name: "no operator found => NullToken",
			elements: ElementList{
				{Token: Number, TokenValue: "1"},
				{Token: Number, TokenValue: "2"},
			},
			operators: []TokenId{Plus, Minus, Multiply, Divide, Exponent},
			wantTok:   NullToken,
			wantIdx:   -1,
		},
		{
			name: "finds last among a subset",
			elements: ElementList{
				{Token: Number, TokenValue: "1"},
				{Token: Plus, TokenValue: "+"},
				{Token: Number, TokenValue: "2"},
				{Token: Minus, TokenValue: "-"},
				{Token: Number, TokenValue: "3"},
			},
			operators: []TokenId{Plus, Minus},
			wantTok:   Minus,
			wantIdx:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTok, gotIdx := tt.elements.FindRightOperator(tt.operators)
			if gotTok != tt.wantTok || gotIdx != tt.wantIdx {
				t.Fatalf("FindRightOperator() = (%v,%d), want (%v,%d)", gotTok, gotIdx, tt.wantTok, tt.wantIdx)
			}
		})
	}
}

func TestElementList_FindLParen(t *testing.T) {
	tests := []struct {
		name     string
		elements ElementList
		wantIdx  int
		wantOK   bool
	}{
		{
			name:     "empty => not found",
			elements: ElementList{},
			wantIdx:  -1,
			wantOK:   false,
		},
		{
			name: "none => not found",
			elements: ElementList{
				{Token: Number, TokenValue: "1"},
				{Token: Plus, TokenValue: "+"},
				{Token: Number, TokenValue: "2"},
			},
			wantIdx: -1,
			wantOK:  false,
		},
		{
			name: "found at start",
			elements: ElementList{
				{Token: LParen, TokenValue: "("},
				{Token: Number, TokenValue: "1"},
			},
			wantIdx: 0,
			wantOK:  true,
		},
		{
			name: "found later",
			elements: ElementList{
				{Token: Number, TokenValue: "9"},
				{Token: Plus, TokenValue: "+"},
				{Token: LParen, TokenValue: "("},
				{Token: Number, TokenValue: "1"},
			},
			wantIdx: 2,
			wantOK:  true,
		},
		{
			name: "multiple => returns first",
			elements: ElementList{
				{Token: Number, TokenValue: "9"},
				{Token: LParen, TokenValue: "("},
				{Token: LParen, TokenValue: "("},
			},
			wantIdx: 1,
			wantOK:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, ok := tt.elements.FindLParen()
			if gotIdx != tt.wantIdx || ok != tt.wantOK {
				t.Fatalf("FindLParen() = (%d,%v), want (%d,%v)", gotIdx, ok, tt.wantIdx, tt.wantOK)
			}
		})
	}
}

func TestElementList_FindRParen(t *testing.T) {
	tests := []struct {
		name     string
		elements ElementList
		lIdx     int
		wantIdx  int
		wantErr  bool
	}{
		{
			name: "simple pair",
			elements: ElementList{
				{Token: LParen, TokenValue: "("},
				{Token: Number, TokenValue: "1"},
				{Token: RParen, TokenValue: ")"},
			},
			lIdx:    0,
			wantIdx: 2,
			wantErr: false,
		},
		{
			name: "empty parentheses",
			elements: ElementList{
				{Token: LParen, TokenValue: "("},
				{Token: RParen, TokenValue: ")"},
			},
			lIdx:    0,
			wantIdx: 1,
			wantErr: false,
		},
		{
			name: "nested - outer matches last",
			elements: ElementList{
				{Token: LParen, TokenValue: "("}, // 0
				{Token: LParen, TokenValue: "("}, // 1
				{Token: Number, TokenValue: "1"}, // 2
				{Token: RParen, TokenValue: ")"}, // 3
				{Token: RParen, TokenValue: ")"}, // 4
			},
			lIdx:    0,
			wantIdx: 4,
			wantErr: false,
		},
		{
			name: "nested - inner matches inner close",
			elements: ElementList{
				{Token: LParen, TokenValue: "("}, // 0
				{Token: LParen, TokenValue: "("}, // 1
				{Token: Number, TokenValue: "1"}, // 2
				{Token: RParen, TokenValue: ")"}, // 3
				{Token: RParen, TokenValue: ")"}, // 4
			},
			lIdx:    1,
			wantIdx: 3,
			wantErr: false,
		},
		{
			name: "error - lParenIdx out of range",
			elements: ElementList{
				{Token: LParen, TokenValue: "("},
				{Token: RParen, TokenValue: ")"},
			},
			lIdx:    2,
			wantIdx: 0,
			wantErr: true,
		},
		{
			name: "error - token at index is not LParen",
			elements: ElementList{
				{Token: Number, TokenValue: "7"},
				{Token: LParen, TokenValue: "("},
				{Token: RParen, TokenValue: ")"},
			},
			lIdx:    0,
			wantIdx: 0,
			wantErr: true,
		},
		{
			name: "error - unmatched parentheses",
			elements: ElementList{
				{Token: LParen, TokenValue: "("},
				{Token: Number, TokenValue: "1"},
				{Token: Plus, TokenValue: "+"},
				{Token: Number, TokenValue: "2"},
			},
			lIdx:    0,
			wantIdx: 0,
			wantErr: true,
		},
		{
			name: "error - unmatched due to unclosed inner paren",
			elements: ElementList{
				{Token: LParen, TokenValue: "("}, // 0
				{Token: LParen, TokenValue: "("}, // 1
				{Token: Number, TokenValue: "1"}, // 2
				{Token: RParen, TokenValue: ")"}, // 3
			},
			lIdx:    0,
			wantIdx: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, err := tt.elements.FindRParen(tt.lIdx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindRParen() err=%v wantErr=%v", err, tt.wantErr)
			}
			if gotIdx != tt.wantIdx {
				t.Fatalf("FindRParen() got=%d want=%d", gotIdx, tt.wantIdx)
			}
		})
	}
}
