package lexer

/*
Language feature currently being worked on:
print(<expression>)

Tokens:
PRINT
*/

import (
	"learningLanguage/token"
	"testing"
)

func TestLexerVariables(t *testing.T) {
	input := `create int a;
			create bool b;
			set b = true;
			set a = 5;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CREATE, "create"},
		{token.INT, "int"},
		{token.IDENT, "a"},
		{token.SEMICOLON, ";"},
		{token.CREATE, "create"},
		{token.BOOL, "bool"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.SET, "set"},
		{token.IDENT, "b"},
		{token.ASSIGN, "="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.SET, "set"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerOperations(t *testing.T) {
	input := `set a = a * 1;
			set b = b / 5;
			set a = a - 1;
			set b = b + 1;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.SET, "set"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "a"},
		{token.MULTIPLY, "*"},
		{token.NUMBER, "1"},
		{token.SEMICOLON, ";"},
		{token.SET, "set"},
		{token.IDENT, "b"},
		{token.ASSIGN, "="},
		{token.IDENT, "b"},
		{token.DIVIDE, "/"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.SET, "set"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "a"},
		{token.MINUS, "-"},
		{token.NUMBER, "1"},
		{token.SEMICOLON, ";"},
		{token.SET, "set"},
		{token.IDENT, "b"},
		{token.ASSIGN, "="},
		{token.IDENT, "b"},
		{token.PLUS, "+"},
		{token.NUMBER, "1"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerFlowControl(t *testing.T) {
	input := `if (a > b) begin;
			a; end;
			else begin;
			b; end;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.GT, ">"},
		{token.IDENT, "b"},
		{token.RPAEREN, ")"},
		{token.BEGIN, "begin"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "a"},
		{token.SEMICOLON, ";"},
		{token.END, "end"},
		{token.SEMICOLON, ";"},
		{token.ELSE, "else"},
		{token.BEGIN, "begin"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.END, "end"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestBooleans(t *testing.T) {
	input := `!true false != > >= < <= == or and`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.NOT, "!"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.NEQ, "!="},
		{token.GT, ">"},
		{token.GE, ">="},
		{token.LT, "<"},
		{token.LE, "<="},
		{token.EQ, "=="},
		{token.OR, "or"},
		{token.AND, "and"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestStruct(t *testing.T) {
	input := `struct , : . [ ]`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STRUCT, "struct"},
		{token.COMMA, ","},
		{token.COLON, ":"},
		{token.DOT, "."},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestFloatString(t *testing.T) {
	input := `float string "Hello World" 3.14`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.FLOAT, "float"},
		{token.STRING, "string"},
		{token.QUOTE, "\"Hello World\""},
		{token.NUMBER, "3.14"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestPrint(t *testing.T) {
	input := `print("Hello World")`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PRINT, "print"},
		{token.LPAREN, "("},
		{token.QUOTE, "\"Hello World\""},
		{token.RPAEREN, ")"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestWhile(t *testing.T) {
	input := `while`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.WHILE, "while"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
