package app

import (
	"reflect"
	"testing"
)

func TestFindOperatorElements(t *testing.T) {
	tests := []struct {
		name     string
		tok      token
		elements []Element
		want     []Element
		wantErr  bool
	}{
		{
			name: "plus - valid in middle",
			tok:  Plus,
			elements: []Element{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			want: []Element{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			wantErr: false,
		},
		{
			name: "multiply - valid in middle",
			tok:  Multiply,
			elements: []Element{
				{token: Number, tokenValue: "3"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "4"},
			},
			want: []Element{
				{token: Number, tokenValue: "3"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "4"},
			},
			wantErr: false,
		},
		{
			name: "divide - valid in longer expression (returns first window)",
			tok:  Divide,
			elements: []Element{
				{token: Number, tokenValue: "8"},
				{token: Divide, tokenValue: "/"},
				{token: Number, tokenValue: "2"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "1"},
			},
			want: []Element{
				{token: Number, tokenValue: "8"},
				{token: Divide, tokenValue: "/"},
				{token: Number, tokenValue: "2"},
			},
			wantErr: false,
		},
		{
			name: "error - invalid operator token (Number is not an operator)",
			tok:  Number,
			elements: []Element{
				{token: Number, tokenValue: "1"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - operator at beginning",
			tok:  Minus,
			elements: []Element{
				{token: Minus, tokenValue: "-"},
				{token: Number, tokenValue: "1"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - operator at end",
			tok:  Plus,
			elements: []Element{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - token before operator is not Number",
			tok:  Multiply,
			elements: []Element{
				{token: LParen, tokenValue: "("},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "2"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - token after operator is not Number",
			tok:  Divide,
			elements: []Element{
				{token: Number, tokenValue: "8"},
				{token: Divide, tokenValue: "/"},
				{token: RParen, tokenValue: ")"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - operator not found in elements",
			tok:  Multiply,
			elements: []Element{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "multiple operators of same type - returns first occurrence",
			tok:  Plus,
			elements: []Element{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "3"},
			},
			want: []Element{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findOperatorElements(tt.tok, tt.elements)
			if (err != nil) != tt.wantErr {
				t.Fatalf("findOperatorElements() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findOperatorElements() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestFindLParen(t *testing.T) {
	tests := []struct {
		name     string
		elements []Element
		wantIdx  int
		wantBool bool
	}{
		{
			name:     "empty slice",
			elements: []Element{},
			wantIdx:  -1,
			wantBool: false,
		},
		{
			name: "no left parentheses",
			elements: []Element{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			wantIdx:  -1,
			wantBool: false,
		},
		{
			name: "left paren at start",
			elements: []Element{
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "1"},
			},
			wantIdx:  0,
			wantBool: true,
		},
		{
			name: "left paren later",
			elements: []Element{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "2"},
			},
			wantIdx:  2,
			wantBool: true,
		},
		{
			name: "multiple left parens - returns first",
			elements: []Element{
				{token: Number, tokenValue: "1"},
				{token: LParen, tokenValue: "("},
				{token: LParen, tokenValue: "("},
			},
			wantIdx:  1,
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, tf := findLParen(tt.elements)
			if tf != tt.wantBool {
				t.Fatalf("findLParen() tf = %v, wantBool %v", tf, tt.wantBool)
			}
			if gotIdx != tt.wantIdx {
				t.Errorf("findLParen() gotIdx = %v, want %v", gotIdx, tt.wantIdx)
			}
		})
	}
}

func TestFindNextOperator(t *testing.T) {
	tests := []struct {
		name      string
		elements  []Element
		operators []token
		wantTok   token
		wantIdx   int
	}{
		{
			name: "finds first multiply when searching multiply/divide",
			elements: []Element{
				{token: Number, tokenValue: "2"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "3"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "4"},
			},
			operators: []token{Multiply, Divide},
			wantTok:   Multiply,
			wantIdx:   3,
		},
		{
			name: "finds first plus when searching plus/minus",
			elements: []Element{
				{token: Number, tokenValue: "2"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "3"},
				{token: Minus, tokenValue: "-"},
				{token: Number, tokenValue: "1"},
			},
			operators: []token{Plus, Minus},
			wantTok:   Plus,
			wantIdx:   1,
		},
		{
			name: "returns NullToken when no operator found",
			elements: []Element{
				{token: Number, tokenValue: "2"},
				{token: Number, tokenValue: "3"},
			},
			operators: []token{Plus, Minus, Multiply, Divide},
			wantTok:   NullToken,
			wantIdx:   -1,
		},
		{
			name: "returns first matching token based on slice scan (not operator list order)",
			elements: []Element{
				{token: Number, tokenValue: "8"},
				{token: Divide, tokenValue: "/"},
				{token: Number, tokenValue: "2"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "3"},
			},
			operators: []token{Multiply, Divide},
			wantTok:   Divide,
			wantIdx:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTok, gotIdx := findNextOperator(tt.elements, tt.operators)
			if gotTok != tt.wantTok || gotIdx != tt.wantIdx {
				t.Errorf("findNextOperator() = (%v, %v), want (%v, %v)", gotTok, gotIdx, tt.wantTok, tt.wantIdx)
			}
		})
	}
}

func TestEval(t *testing.T) {
	tests := []struct {
		name     string
		elements []Element
		want     *int
		wantErr  bool
	}{
		{
			name:     "empty slice",
			elements: []Element{},
			want:     ptrInt(0),
			wantErr:  true,
		},
		{
			name:     "nil slice",
			elements: nil,
			want:     ptrInt(0),
			wantErr:  true,
		},
		{
			name: "single number",
			elements: []Element{
				{token: Number, tokenValue: "42"},
			},
			want:    ptrInt(42),
			wantErr: false,
		},
		{
			name: "multiply before plus (precedence)",
			elements: []Element{
				{token: Number, tokenValue: "2"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "3"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "4"},
			},
			want:    ptrInt(14), // 2 + (3*4)
			wantErr: false,
		},
		{
			name: "left associative subtraction",
			elements: []Element{
				{token: Number, tokenValue: "10"},
				{token: Minus, tokenValue: "-"},
				{token: Number, tokenValue: "3"},
				{token: Minus, tokenValue: "-"},
				{token: Number, tokenValue: "4"},
			},
			want:    ptrInt(3), // (10-3)-4
			wantErr: false,
		},
		{
			name: "integer division",
			elements: []Element{
				{token: Number, tokenValue: "7"},
				{token: Divide, tokenValue: "/"},
				{token: Number, tokenValue: "2"},
			},
			want:    ptrInt(3),
			wantErr: false,
		},
		{
			name: "parentheses change precedence",
			elements: []Element{
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "2"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "3"},
				{token: RParen, tokenValue: ")"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "4"},
			},
			want:    ptrInt(20), // (2+3)*4
			wantErr: false,
		},
		{
			name: "nested parentheses",
			elements: []Element{
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "2"},
				{token: Multiply, tokenValue: "*"},
				{token: Number, tokenValue: "3"},
				{token: RParen, tokenValue: ")"},
				{token: RParen, tokenValue: ")"},
				{token: Minus, tokenValue: "-"},
				{token: Number, tokenValue: "4"},
			},
			want:    ptrInt(3), // (1+(2*3)) - 4
			wantErr: false,
		},
		{
			name: "error - invalid expression leaves multiple elements",
			elements: []Element{
				{token: Number, tokenValue: "1"},
				{token: Number, tokenValue: "2"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - unmatched parentheses",
			elements: []Element{
				{token: LParen, tokenValue: "("},
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
				{token: Number, tokenValue: "2"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - operator missing right operand",
			elements: []Element{
				{token: Number, tokenValue: "1"},
				{token: Plus, tokenValue: "+"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - invalid number tokenValue",
			elements: []Element{
				{token: Number, tokenValue: "not-a-number"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Eval(tt.elements)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Eval() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got == nil || tt.want == nil || *got != *tt.want {
				t.Fatalf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptrInt(v int) *int { return &v }
func TestFindRParen(t *testing.T) {
	type args struct {
		elements  []Element
		lParenIdx int
	}
	tests := []struct {
		name          string
		args          args
		wantRParenIdx int
		wantErr       bool
	}{
		{
			name: "simple pair",
			args: args{
				elements: []Element{
					{token: LParen, tokenValue: "("},
					{token: Number, tokenValue: "1"},
					{token: RParen, tokenValue: ")"},
				},
				lParenIdx: 0,
			},
			wantRParenIdx: 2,
			wantErr:       false,
		},
		{
			name: "empty parentheses",
			args: args{
				elements: []Element{
					{token: LParen, tokenValue: "("},
					{token: RParen, tokenValue: ")"},
				},
				lParenIdx: 0,
			},
			wantRParenIdx: 1,
			wantErr:       false,
		},
		{
			name: "pair not at beginning",
			args: args{
				elements: []Element{
					{token: Number, tokenValue: "9"},
					{token: Plus, tokenValue: "+"},
					{token: LParen, tokenValue: "("},
					{token: Number, tokenValue: "1"},
					{token: RParen, tokenValue: ")"},
				},
				lParenIdx: 2,
			},
			wantRParenIdx: 4,
			wantErr:       false,
		},
		{
			name: "nested - outer matches last",
			args: args{
				elements: []Element{
					{token: LParen, tokenValue: "("}, // 0
					{token: LParen, tokenValue: "("}, // 1
					{token: Number, tokenValue: "1"}, // 2
					{token: RParen, tokenValue: ")"}, // 3
					{token: RParen, tokenValue: ")"}, // 4
				},
				lParenIdx: 0,
			},
			wantRParenIdx: 4,
			wantErr:       false,
		},
		{
			name: "nested - inner matches inner close",
			args: args{
				elements: []Element{
					{token: LParen, tokenValue: "("}, // 0
					{token: LParen, tokenValue: "("}, // 1
					{token: Number, tokenValue: "1"}, // 2
					{token: RParen, tokenValue: ")"}, // 3
					{token: RParen, tokenValue: ")"}, // 4
				},
				lParenIdx: 1,
			},
			wantRParenIdx: 3,
			wantErr:       false,
		},
		{
			name: "nested with other tokens in between",
			args: args{
				elements: []Element{
					{token: LParen, tokenValue: "("}, // 0
					{token: Number, tokenValue: "1"}, // 1
					{token: Plus, tokenValue: "+"},   // 2
					{token: LParen, tokenValue: "("}, // 3
					{token: Number, tokenValue: "2"}, // 4
					{token: RParen, tokenValue: ")"}, // 5
					{token: Minus, tokenValue: "-"},  // 6
					{token: Number, tokenValue: "3"}, // 7
					{token: RParen, tokenValue: ")"}, // 8
				},
				lParenIdx: 0,
			},
			wantRParenIdx: 8,
			wantErr:       false,
		},
		{
			name: "error - lParenIdx out of range (empty slice)",
			args: args{
				elements:  []Element{},
				lParenIdx: 0,
			},
			wantRParenIdx: 0,
			wantErr:       true,
		},
		{
			name: "error - lParenIdx out of range (too large)",
			args: args{
				elements: []Element{
					{token: LParen, tokenValue: "("},
					{token: RParen, tokenValue: ")"},
				},
				lParenIdx: 2,
			},
			wantRParenIdx: 0,
			wantErr:       true,
		},
		{
			name: "error - lParenIdx is negative",
			args: args{
				elements: []Element{
					{token: LParen, tokenValue: "("},
					{token: RParen, tokenValue: ")"},
				},
				lParenIdx: -1,
			},
			wantRParenIdx: 0,
			wantErr:       true,
		},
		{
			name: "error - token at index is not LParen",
			args: args{
				elements: []Element{
					{token: Number, tokenValue: "7"},
					{token: LParen, tokenValue: "("},
					{token: Number, tokenValue: "1"},
					{token: RParen, tokenValue: ")"},
				},
				lParenIdx: 0,
			},
			wantRParenIdx: 0,
			wantErr:       true,
		},
		{
			name: "error - unmatched parentheses (no closing)",
			args: args{
				elements: []Element{
					{token: LParen, tokenValue: "("},
					{token: Number, tokenValue: "1"},
					{token: Plus, tokenValue: "+"},
					{token: Number, tokenValue: "2"},
				},
				lParenIdx: 0,
			},
			wantRParenIdx: 0,
			wantErr:       true,
		},
		{
			name: "error - unmatched parentheses (inner opened not closed)",
			args: args{
				elements: []Element{
					{token: LParen, tokenValue: "("}, // 0
					{token: LParen, tokenValue: "("}, // 1
					{token: Number, tokenValue: "1"}, // 2
					{token: RParen, tokenValue: ")"}, // 3
				},
				lParenIdx: 0,
			},
			wantRParenIdx: 0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRParenIdx, err := findRParen(tt.args.elements, tt.args.lParenIdx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("findRParen() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotRParenIdx != tt.wantRParenIdx {
				t.Errorf("findRParen() gotRParenIdx = %v, want %v", gotRParenIdx, tt.wantRParenIdx)
			}
		})
	}
}
