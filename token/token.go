package token

type TokenType string

// A token consists of a tokenType (string), and literal text
type Token struct {
	Type    TokenType
	Literal string
}

// types of tokens available
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	NUMBER = "NUMBER"

	ASSIGN   = "ASSIGN"
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	DIVIDE   = "DIVIDE"
	MULTIPLY = "MULTIPLY"

	EQ  = "EQUALTO"
	GT  = "GREATER"
	GE  = "GREQUAL"
	LT  = "LESS"
	LE  = "LEQUAL"
	NEQ = "NOTEQUAL"
	NOT = "NOT"

	LPAREN    = "LPAREN"
	RPAEREN   = "RPAREN"
	SEMICOLON = "SEMICOLON"
	COLON     = "COLON"
	DOT       = "DOT"
	COMMA     = "COMMA"
	LBRACKET  = "LBRACKET"
	RBRACKET  = "RBRACKET"
	QUOTE     = "QUOTE"

	SET    = "SET"
	CREATE = "CREATE"
	IF     = "IF"
	ELSE   = "ELSE"
	BEGIN  = "BEGIN"
	END    = "END"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	STRUCT = "STRUCT"
	INT    = "INT"
	BOOL   = "BOOL"
	FLOAT  = "FLOAT"
	STRING = "STRING"
	PRINT  = "PRINT"
	AND    = "AND"
	OR     = "OR"
)
