package parser

import (
	"fmt"
	"learningLanguage/ast"
	"learningLanguage/lexer"
	"learningLanguage/token"
	"slices"
	"strconv"
	"strings"
)

// precedence integer values for pratt parsing of expressions
const (
	LOWEST      = iota + 1 // everything else
	EQUALS                 // ==
	LESSGREATER            // >, >=, <, <=
	ORAND                  // or/and boolean logic
	SUMDIFF                // + or -
	PRODUCTDIV             // * or /
	PREFIX                 // -X
)

// map tokens to precedence enumerations
var precedences = map[token.TokenType]int{
	token.PLUS:     SUMDIFF,
	token.MINUS:    SUMDIFF,
	token.DIVIDE:   PRODUCTDIV,
	token.MULTIPLY: PRODUCTDIV,
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.GT:       LESSGREATER,
	token.GE:       LESSGREATER,
	token.LT:       LESSGREATER,
	token.LE:       LESSGREATER,
	token.OR:       ORAND,
	token.AND:      ORAND,
}

// list of valid datatypes
var DATATYPES = []token.TokenType{token.BOOL, token.INT, token.STRING, token.FLOAT}

// function type definitions for prefix and infix parsing functions
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// parser object definition, contains a lexer, the current token, the token after it, a list of errors, and maps for prefix and infix parsing functions
type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// get the precedence for the given token
func getPrecedence(tok token.Token) int {
	precedence, ok := precedences[tok.Type]
	if ok {
		return precedence
	}
	return LOWEST
}

// helper functions to add prefix and infix functions to the parser
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// parser constructor
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// get first two tokens to assign currToken and peekToken
	p.nextToken()
	p.nextToken()

	// register all prefix functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumber)
	p.registerPrefix(token.QUOTE, p.parseStringLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.LBRACKET, p.parseListLiteral)
	p.registerPrefix(token.LEN, p.parseLengthExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)

	//register all infix functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GE, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LE, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)

	return p
}

// get function for the parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// get next token and move peek to current token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// parse a program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// parse statements until there are no more
	for p.curToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

// switch function to tie tokens to statement AST nodes
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	// create statement
	case token.CREATE:
		stmt := p.parseCreateStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	// set statement
	case token.SET:
		stmt := p.parseSetStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	// if statement
	case token.IF:
		stmt := p.parseIfStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	// while statement
	case token.WHILE:
		stmt := p.parseWhileStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	// struct statement
	case token.STRUCT:
		stmt := p.parseStructStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	// print statement
	case token.PRINT:
		stmt := p.parsePrintStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	// print statement with new line
	case token.PRINTLN:
		stmt := p.parsePrintStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	// count statement
	case token.COUNT:
		stmt := p.parseCountStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	// append statement
	case token.APPEND:
		stmt := p.parseAppendStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	// if no other statement types then it is an expression statement
	default:
		stmt := p.parseExpressionStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	}

}

// check next token type and send an error if it does not match
func (p *Parser) checkNextToken(tokType token.TokenType) bool {
	if p.peekToken.Type == tokType {
		p.nextToken()
		return true
	} else {
		p.peekError(tokType)
		return false
	}
}

// check next token for multiple types and send an error if it does not match
func (p *Parser) checkMultipleNextToken(tokTypes []token.TokenType) bool {
	if slices.Contains(tokTypes, p.peekToken.Type) {
		p.nextToken()
		return true
	} else {
		p.peekMultiError(tokTypes)
		return false
	}
}

// send an error indicating that multiple token types were expected but none of them match peek token
func (p *Parser) peekMultiError(tokTypes []token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %v, received %s.", tokTypes, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// send an error indicating that a token type was expected but doesn't match peek token
func (p *Parser) peekError(tokType token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, received %s.", tokType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

/*
For the following functions, the format will be used to define the language grammar
TOKENTYPE, (OPTIONAL TOKEN TYPE), [CHOOSE ONE TOKEN TYPE], <EXPRESSION/STATEMENT>, (ONE OR MORE)*
*/

// CREATE [INT/BOOL/STRING/FLOAT] (LIST) IDENTIFIER SEMICOLON
func (p *Parser) parseCreateStatement() *ast.CreateStatement {
	statement := &ast.CreateStatement{Token: p.curToken}

	// [INT/BOOL/STRING/FLOAT]
	if !p.checkMultipleNextToken(DATATYPES) { //TODO: add other datatypes
		return nil
	}

	var datatype string
	switch p.curToken.Type {
	case token.INT:
		datatype = "int"
	case token.BOOL:
		datatype = "bool"
	case token.STRING:
		datatype = "string"
	case token.FLOAT:
		datatype = "float"
	}

	var isList bool

	// (LIST)
	if p.peekToken.Type == token.LIST {
		isList = true
		p.nextToken()
	}

	// IDENTIFIER
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, DataType: datatype, IsList: isList}

	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

// SET IDENTIFIER (DOT IDENT)([<EXPRESSION>]) ASSIGN <EXPRESSION> SEMICOLON
func (p *Parser) parseSetStatement() *ast.SetStatement {
	statement := &ast.SetStatement{Token: p.curToken}

	// IDENTIFIER
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// (DOT IDENT)
	if p.peekToken.Type == token.DOT {
		p.checkNextToken(token.DOT)
		if !p.checkNextToken(token.IDENT) {
			return nil
		}
		statement.Name.Attribute = p.curToken.Literal
	}

	// ([<EXPRESSION>])
	if p.peekToken.Type == token.LBRACKET {
		p.nextToken()
		p.nextToken()
		statement.Name.Index = p.parseExpression(LOWEST)
		if !p.checkNextToken(token.RBRACKET) {
			return nil
		}
	}

	// ASSIGN
	if !p.checkNextToken(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	// <EXPRESSION>
	statement.Value = p.parseExpression(LOWEST)

	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

// <EXPRESSION> SEMICOLON
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.curToken}

	// <EXPRESSION>
	statement.Expression = p.parseExpression(LOWEST)

	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

// IF LPAREN <EXPRESSION> RPAREN BEGIN SEMICOLON (<STATEMENT>)* END SEMICOLON (ELSE BEGIN SEMICOLON (<STATEMENT>)* END SEMICOLON)
func (p *Parser) parseIfStatement() *ast.IfStatement {
	statement := &ast.IfStatement{Token: p.curToken}

	// LPAREN
	if !p.checkNextToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	// <EXPRESSION>
	statement.Condition = p.parseExpression(LOWEST)

	// LPAREN
	if !p.checkNextToken(token.RPAREN) {
		return nil
	}
	// BEGIN
	if !p.checkNextToken(token.BEGIN) {
		return nil
	}
	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	// (<STATEMENT>)*
	for p.peekToken.Type != token.END {
		p.nextToken()
		statement.IfTrue = append(statement.IfTrue, p.parseStatement())
		if p.peekToken.Type == token.ELSE {
			p.errors = append(p.errors, "ELSE found before END;")
			return nil
		}
		if p.peekToken.Type == token.EOF {
			p.errors = append(p.errors, "EOF found before END;")
			return nil
		}
	}

	// END
	if !p.checkNextToken(token.END) {
		return nil
	}
	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	// (ELSE BEGIN SEMICOLON (<STATEMENT>)* END SEMICOLON)
	if p.peekToken.Type != token.ELSE {
		statement.Else = nil
		return statement
	} else {
		p.nextToken()
		// BEGIN
		if !p.checkNextToken(token.BEGIN) {
			return nil
		}
		// SEMICOLON
		if !p.checkNextToken(token.SEMICOLON) {
			return nil
		}

		// (<STATEMENT>)*
		for p.peekToken.Type != token.END {
			p.nextToken()
			statement.Else = append(statement.Else, p.parseStatement())
			// check for end on else statement, if none then send error
			if p.peekToken.Type == token.EOF {
				p.errors = append(p.errors, "Missing END Token in else statement.")
				return nil
			}
		}

		// END
		if !p.checkNextToken(token.END) {
			return nil
		}
		// SEMICOLON
		if !p.checkNextToken(token.SEMICOLON) {
			return nil
		}
		return statement
	}
}

// WHILE LPAREN <EXPRESSION> RPAREN BEGIN SEMICOLON (<STATEMENT>)* END SEMICOLON
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	statement := &ast.WhileStatement{Token: p.curToken}

	// LPAREN
	if !p.checkNextToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	// <EXPRESSION>
	statement.Condition = p.parseExpression(LOWEST)

	// RPAREN
	if !p.checkNextToken(token.RPAREN) {
		return nil
	}

	// BEGIN
	if !p.checkNextToken(token.BEGIN) {
		return nil
	}

	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	// (<STATEMENT>)*
	for p.peekToken.Type != token.END {
		p.nextToken()
		statement.LoopStatements = append(statement.LoopStatements, p.parseStatement())
		// if end of file before end, send an error
		if p.peekToken.Type == token.EOF {
			p.errors = append(p.errors, "Missing END on while loop.")
			return nil
		}
	}

	// END
	if !p.checkNextToken(token.END) {
		return nil
	}

	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

// COUNT IDENTIFIER FROM <EXPRESSION> TO <EXPRESSION (BY <EXPRESSION>) BEGIN SEMICOLON <STATEMENTS> END SEMICOLON
func (p *Parser) parseCountStatement() *ast.CountStatement {
	statement := &ast.CountStatement{Token: p.curToken}

	// IDENT
	if !p.checkNextToken(token.IDENT) {
		return nil
	}

	statement.Counter = ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// FROM
	if !p.checkNextToken(token.FROM) {
		return nil
	}

	p.nextToken()
	// <EXPRESSION>
	statement.From = p.parseExpression(LOWEST)

	// TO
	if !p.checkNextToken(token.TO) {
		return nil
	}

	p.nextToken()
	// <EXPRESSION>
	statement.To = p.parseExpression(LOWEST)

	// (BY <EXPRESSION)
	if p.peekToken.Type != token.BY {
		statement.By = &ast.IntegerLiteral{Token: token.Token{Type: token.NUMBER, Literal: "1"}, Value: 1}
	} else {
		p.nextToken()
		p.nextToken()
		statement.By = p.parseExpression(LOWEST)
	}

	// BEGIN
	if !p.checkNextToken(token.BEGIN) {
		return nil
	}

	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	// (<STATEMENT>)*
	for p.peekToken.Type != token.END {
		p.nextToken()
		statement.LoopStatements = append(statement.LoopStatements, p.parseStatement())
		if p.peekToken.Type == token.EOF {
			p.errors = append(p.errors, "Missing END on while loop.")
			return nil
		}
	}

	// END
	if !p.checkNextToken(token.END) {
		return nil
	}

	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

// STRUCT IDENTIFIER LPAREN ([INT/BOOL/STRING/FLOAT] IDENTIFIER)* RPAREN (LBRACKET (IDENTIFIER COLON <EXPRESSION>)* RBRACKET) SEMICOLON
func (p *Parser) parseStructStatement() *ast.StructStatement {
	statement := &ast.StructStatement{Token: p.curToken, Values: make(map[string]ast.Expression)}

	// IDENTIFIER
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	statement.StructIdent = ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, DataType: "struct"}

	// LPAREN
	if !p.checkNextToken(token.LPAREN) {
		return nil
	}
	// [INT/BOOL/STRING/FLOAT]
	if !p.checkMultipleNextToken(DATATYPES) {
		return nil
	}
	var datatype string
	switch p.curToken.Type {
	case token.INT:
		datatype = "int"
	case token.BOOL:
		datatype = "bool"
	case token.STRING:
		datatype = "string"
	case token.FLOAT:
		datatype = "float"
	}
	// IDENT
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, DataType: datatype}
	statement.Attributes = append(statement.Attributes, *ident)

	// ([INT/BOOL/STRING/FLOAT] IDENTIFIER)*
	for p.peekToken.Type == token.COMMA {
		p.checkNextToken(token.COMMA)
		if !p.checkMultipleNextToken(DATATYPES) {
			return nil
		}
		switch p.curToken.Type {
		case token.INT:
			datatype = "int"
		case token.BOOL:
			datatype = "bool"
		case token.STRING:
			datatype = "string"
		case token.FLOAT:
			datatype = "float"
		}
		if !p.checkNextToken(token.IDENT) {
			return nil
		}
		ident = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, DataType: datatype}
		statement.Attributes = append(statement.Attributes, *ident)
	}
	// RPAREN
	if !p.checkNextToken(token.RPAREN) {
		return nil
	}

	// if no attribute assignment, then return statement
	if p.peekToken.Type == token.SEMICOLON {
		p.checkNextToken(token.SEMICOLON)
		return statement
	}

	// LBRACKET IDENTIFIER COLON <EXPRESSION>
	// LBRACKET
	if !p.checkNextToken(token.LBRACKET) {
		return nil
	}
	// IDENT
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	identString := p.curToken.Literal
	// COLON
	if !p.checkNextToken(token.COLON) {
		return nil
	}
	p.nextToken()
	// <EXPRESSION>
	statement.Values[identString] = p.parseExpression(LOWEST)

	//(IDENTIFIER COLON <EXPRESSION>)*
	for p.peekToken.Type == token.COMMA {
		p.checkNextToken(token.COMMA)
		if !p.checkNextToken(token.IDENT) {
			return nil
		}
		identString := p.curToken.Literal
		if !p.checkNextToken(token.COLON) {
			return nil
		}
		p.nextToken()
		statement.Values[identString] = p.parseExpression(LOWEST)
	}

	// RBRACKET
	if !p.checkNextToken(token.RBRACKET) {
		return nil
	}
	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

// [PRINT/PRINTLN] <EXPRESSION> SEMICOLON
func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	statement := &ast.PrintStatement{Token: p.curToken}

	// if println, then set newline to true
	if p.curToken.Type == token.PRINTLN {
		statement.NewLine = true
	}

	// LPAREN
	if !p.checkNextToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	// <EXPRESSION>
	statement.Value = p.parseExpression(LOWEST)

	// RPAREN
	if !p.checkNextToken(token.RPAREN) {
		return nil
	}

	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

// APPEND <EXPRESSION> TO IDENTIFIER SEMICOLON
func (p *Parser) parseAppendStatement() *ast.AppendStatement {
	statement := &ast.AppendStatement{Token: p.curToken}

	p.nextToken()
	// <EXPRESSION>
	statement.Value = p.parseExpression(LOWEST)

	// TO
	if !p.checkNextToken(token.TO) {
		return nil
	}

	// IDENT
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	statement.List = ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, IsList: true}

	// SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

// send an error if there was no token matching a prefix function
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// expression pratt parsing
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	// start with prefix
	leftExp := prefix()

	// form a correct infix expression tree based on precedence
	for p.peekToken.Type != token.SEMICOLON && precedence < getPrecedence(p.peekToken) {
		infix := p.infixParseFns[p.peekToken.Type]
		// if there are no infix operators, return just the prefix
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		// parse infix expression using the left expression
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// if the identifier has an attribute (for structs) assign the attribute in the ast node
	if p.peekToken.Type == token.DOT {
		p.checkNextToken(token.DOT)
		if !p.checkNextToken(token.IDENT) {
			return nil
		}
		ident.Attribute = p.curToken.Literal
	}
	// if the identifier has an index, assign the index in the ast node
	if p.peekToken.Type == token.LBRACKET {
		p.nextToken()
		p.nextToken()
		ident.Index = p.parseExpression(LOWEST)
		if !p.checkNextToken(token.RBRACKET) {
			return nil
		}
	}
	return ident
}

func (p *Parser) parseNumber() ast.Expression {
	// if the number has a dot, parse a float literal, otherwise parse an int literal
	if strings.Contains(p.curToken.Literal, ".") {
		return p.parseFloatLiteral()
	} else {
		return p.parseIntegerLiteral()
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	// turn token literal into 64 bit integer
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	// if not possible, throw error
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	// turn 64 bit int into 32 bit ast node value
	lit.Value = int32(value)

	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	// convert token literal to a 64 bit float
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	// if not possible, throw error
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	// cast to 32 bit float for ast node value
	lit.Value = float32(value)

	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.BooleanLiteral{Token: p.curToken}

	// assign value according to token literal
	switch p.curToken.Literal {
	case "false":
		lit.Value = false
	case "true":
		lit.Value = true
	// if neither true or false, send error
	default:
		msg := fmt.Sprintf("could not parse %q as boolean", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return lit
}

// LBRACKET (<EXPRESSION>)* RBRACKET
func (p *Parser) parseListLiteral() ast.Expression {
	lit := &ast.ListLiteral{Token: p.curToken}
	list := []ast.Expression{}

	// get expressions until we find a right bracket
	for p.curToken.Type != token.RBRACKET {
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
		// if the next token is neither a comma or bracket, send an error
		if !p.checkMultipleNextToken([]token.TokenType{token.COMMA, token.RBRACKET}) {
			return nil
		}
	}

	lit.List = list

	return lit
}

// parse string literal ast node from the token literal
func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	return lit
}

// LEN LPAREN IDENTIFIER RPAREN
func (p *Parser) parseLengthExpression() ast.Expression {
	exp := &ast.LengthExpression{Token: p.curToken}

	// LPAREN
	if !p.checkNextToken(token.LPAREN) {
		return nil
	}

	// IDENT
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	exp.List = ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, IsList: true}

	// RPAREN
	if !p.checkNextToken(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}
	precedence := getPrecedence(p.curToken)

	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}
