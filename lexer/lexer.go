package lexer

import (
	"learningLanguage/token"
	"slices"
)

var keywords = []string{
	"int", "bool", "create",
	"set", "if", "else",
	"begin", "end", "true",
	"false", "struct", "float",
	"string", "print", "or",
	"and", "while", "count",
	"from", "to", "by",
	"list", "len", "append"}

type Lexer struct {
	input string
	head  int
	ch    byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.head >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.head]
	}
	l.head += 1
}

func (l *Lexer) goBack() {
	if l.head > 0 {
		l.head--
	}
}

func (l *Lexer) ignoreWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.ignoreWhitespace()

	if isLetter(l.ch) {
		tok = readString(l)
	} else if isNumber(l.ch) {
		tok = readNumber(l)
	} else {
		literal := string(l.ch)
		switch l.ch {
		case '=':
			temp := string(l.ch)
			l.readChar()
			if l.ch == '=' {
				temp += string(l.ch)
				tok = newToken(token.EQ, temp)
			} else {
				l.goBack()
				tok = newToken(token.ASSIGN, literal)
			}
		case '>':
			temp := string(l.ch)
			l.readChar()
			if l.ch == '=' {
				temp += string(l.ch)
				tok = newToken(token.GE, temp)
			} else {
				l.goBack()
				tok = newToken(token.GT, literal)
			}
		case '<':
			temp := string(l.ch)
			l.readChar()
			if l.ch == '=' {
				temp += string(l.ch)
				tok = newToken(token.LE, temp)
			} else {
				l.goBack()
				tok = newToken(token.LT, literal)
			}
		case '!':
			temp := string(l.ch)
			l.readChar()
			if l.ch == '=' {
				temp += string(l.ch)
				tok = newToken(token.NEQ, temp)
			} else {
				l.goBack()
				tok = newToken(token.NOT, literal)
			}
		case ';':
			tok = newToken(token.SEMICOLON, literal)
		case '+':
			tok = newToken(token.PLUS, literal)
		case '-':
			tok = newToken(token.MINUS, literal)
		case '*':
			tok = newToken(token.MULTIPLY, literal)
		case '/':
			tok = newToken(token.DIVIDE, literal)
		case '(':
			tok = newToken(token.LPAREN, literal)
		case ')':
			tok = newToken(token.RPAREN, literal)
		case '[':
			tok = newToken(token.LBRACKET, literal)
		case ']':
			tok = newToken(token.RBRACKET, literal)
		case ',':
			tok = newToken(token.COMMA, literal)
		case ':':
			tok = newToken(token.COLON, literal)
		case '.':
			tok = newToken(token.DOT, literal)
		case '"':
			l.readChar()
			for l.ch != '"' {
				if l.ch == 0 {
					tok.Literal = "ERROR"
					tok.Type = token.ILLEGAL
					l.readChar()
					return tok
				}
				literal += string(l.ch)
				l.readChar()
			}
			literal += string(l.ch)
			tok = newToken(token.QUOTE, literal)
		case 0:
			tok.Literal = ""
			tok.Type = token.EOF
		default:
			tok.Literal = "ERROR"
			tok.Type = token.ILLEGAL
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func isNumber(char byte) bool {
	return '0' <= char && char <= '9'
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func readNumber(l *Lexer) token.Token {
	str := string(l.ch)
	l.readChar()
	for isNumber(l.ch) || l.ch == '.' {
		str += string(l.ch)
		l.readChar()
	}
	l.goBack()
	return newToken(token.NUMBER, str)
}

func readString(l *Lexer) token.Token {
	var tok token.Token
	str := string(l.ch)
	l.readChar()
	for isLetter(l.ch) {
		str += string(l.ch)
		l.readChar()
	}
	l.goBack()
	isKeyword := checkKeyword(str)
	if isKeyword {
		tok = createKeyword(str)
	} else {
		tok = newToken(token.IDENT, str)
	}
	return tok
}

func checkKeyword(str string) bool {
	return slices.Contains(keywords, str)
}

func createKeyword(str string) token.Token {
	var tok token.Token
	switch str {
	case "create":
		tok = newToken(token.CREATE, str)
	case "set":
		tok = newToken(token.SET, str)
	case "if":
		tok = newToken(token.IF, str)
	case "else":
		tok = newToken(token.ELSE, str)
	case "begin":
		tok = newToken(token.BEGIN, str)
	case "end":
		tok = newToken(token.END, str)
	case "int":
		tok = newToken(token.INT, str)
	case "true":
		tok = newToken(token.TRUE, str)
	case "false":
		tok = newToken(token.FALSE, str)
	case "bool":
		tok = newToken(token.BOOL, str)
	case "struct":
		tok = newToken(token.STRUCT, str)
	case "float":
		tok = newToken(token.FLOAT, str)
	case "string":
		tok = newToken(token.STRING, str)
	case "print":
		tok = newToken(token.PRINT, str)
	case "or":
		tok = newToken(token.OR, str)
	case "and":
		tok = newToken(token.AND, str)
	case "while":
		tok = newToken(token.WHILE, str)
	case "count":
		tok = newToken(token.COUNT, str)
	case "from":
		tok = newToken(token.FROM, str)
	case "to":
		tok = newToken(token.TO, str)
	case "by":
		tok = newToken(token.BY, str)
	case "list":
		tok = newToken(token.LIST, str)
	case "append":
		tok = newToken(token.APPEND, str)
	case "len":
		tok = newToken(token.LEN, str)
	}
	return tok
}
