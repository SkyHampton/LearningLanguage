package lexer

import (
	"learningLanguage/token"
	"slices"
)

// slice of keywords, used to differentiate keywords from identifiers
var keywords = []string{
	"int", "bool", "create",
	"set", "if", "else",
	"begin", "end", "true",
	"false", "struct", "float",
	"string", "print", "or",
	"and", "while", "count",
	"from", "to", "by", "println",
	"list", "len", "append"}

// lexer class, includes an input string, the current character read, and the index of the next character (head)
type Lexer struct {
	input string
	head  int
	ch    byte
}

// lexer constructor
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// read next character and increment head pointer
func (l *Lexer) readChar() {
	if l.head >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.head]
	}
	l.head += 1
}

// helper function, go back one character by decrementing head
func (l *Lexer) goBack() {
	if l.head > 0 {
		l.head--
	}
}

// function to ignore whitespace
func (l *Lexer) ignoreWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// get next token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.ignoreWhitespace()

	// if the character is a letter (a-z, A-Z), its either an identifier or a keyword
	if isLetter(l.ch) {
		tok = readString(l)
		// if the character is a number, read the whole number
	} else if isNumber(l.ch) {
		tok = readNumber(l)
		// non-alphanumeric characters
	} else {
		literal := string(l.ch)
		// switch character and create corresponding token
		switch l.ch {
		case '=':
			temp := string(l.ch)
			l.readChar()
			// check for ==
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
			// check for >=
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
			// check for <=
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
			// check for !=
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
		// quote tokens consist of ", any number of characters, and another "
		case '"':
			l.readChar()
			for l.ch != '"' {
				// if we do not find an end quote, return an illegal token
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
		// error and EOF tokens
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

// create a token object with the specified token type and literal
func newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

// check if given character is a number
func isNumber(char byte) bool {
	return '0' <= char && char <= '9'
}

// check if given character is a letter
func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

// read whole number, including floating point
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

// read whole string of characters
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
	// if its in the keyword slice, create keyword token
	if isKeyword {
		tok = createKeyword(str)
		// otherwise, create ident token
	} else {
		tok = newToken(token.IDENT, str)
	}
	return tok
}

// check if string is a keyword using the slice at the top of the script
func checkKeyword(str string) bool {
	return slices.Contains(keywords, str)
}

// check which keyword it is and create a corresponding token, TODO replace with a map for greater efficiency
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
	case "println":
		tok = newToken(token.PRINTLN, str)
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
