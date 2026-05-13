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

func TestLexerKeywords(t *testing.T) {
	input := `set create if else while count from to by begin end true false struct int bool float string print println and or list append len`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.SET, "set"},
		{token.CREATE, "create"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.WHILE, "while"},
		{token.COUNT, "count"},
		{token.FROM, "from"},
		{token.TO, "to"},
		{token.BY, "by"},
		{token.BEGIN, "begin"},
		{token.END, "end"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.STRUCT, "struct"},
		{token.INT, "int"},
		{token.BOOL, "bool"},
		{token.FLOAT, "float"},
		{token.STRING, "string"},
		{token.PRINT, "print"},
		{token.PRINTLN, "println"},
		{token.AND, "and"},
		{token.OR, "or"},
		{token.LIST, "list"},
		{token.APPEND, "append"},
		{token.LEN, "len"},
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

func TestLexerInfixPrefix(t *testing.T) {
	input := `+ - / * == != > >= < <= !`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.DIVIDE, "/"},
		{token.MULTIPLY, "*"},
		{token.EQ, "=="},
		{token.NEQ, "!="},
		{token.GT, ">"},
		{token.GE, ">="},
		{token.LT, "<"},
		{token.LE, "<="},
		{token.NOT, "!"},
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

func TestLexerIdentNum(t *testing.T) {
	input := `123 3.14 test test123`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.NUMBER, "123"},
		{token.NUMBER, "3.14"},
		{token.IDENT, "test"},
		{token.IDENT, "test"},
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

func TestLexerMisc(t *testing.T) {
	input := `()[];:.,"hello"`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.COLON, ":"},
		{token.DOT, "."},
		{token.COMMA, ","},
		{token.QUOTE, "\"hello\""},
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
