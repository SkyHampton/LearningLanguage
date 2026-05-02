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

const (
	LOWEST      = iota + 1 // everything else
	EQUALS                 // ==
	LESSGREATER            // >, >=, <, <=
	ORAND                  // or/and boolean logic
	SUMDIFF                // + or -
	PRODUCTDIV             // * or /
	PREFIX                 // -X
)

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

var DATATYPES = []token.TokenType{token.BOOL, token.INT, token.STRING, token.FLOAT}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func getPrecedence(tok token.Token) int {
	precedence, ok := precedences[tok.Type]
	if ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumber)
	p.registerPrefix(token.QUOTE, p.parseStringLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)

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

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.CREATE:
		stmt := p.parseCreateStatement()
		if stmt == nil {
			return nil
		}
		return stmt

	case token.SET:
		stmt := p.parseSetStatement()
		if stmt == nil {
			return nil
		}
		return stmt

	case token.IF:
		stmt := p.parseIfStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	case token.WHILE:
		stmt := p.parseWhileStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	case token.STRUCT:
		stmt := p.parseStructStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	case token.PRINT:
		stmt := p.parsePrintStatement()
		if stmt == nil {
			return nil
		}
		return stmt

	default:
		stmt := p.parseExpressionStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	}

}

func (p *Parser) checkNextToken(tokType token.TokenType) bool {
	if p.peekToken.Type == tokType {
		p.nextToken()
		return true
	} else {
		p.peekError(tokType)
		return false
	}
}

func (p *Parser) checkMultipleNextToken(tokTypes []token.TokenType) bool {
	if slices.Contains(tokTypes, p.peekToken.Type) {
		p.nextToken()
		return true
	} else {
		p.peekMultiError(tokTypes)
		return false
	}
}

func (p *Parser) peekMultiError(tokTypes []token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %v, received %s.", tokTypes, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(tokType token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, received %s.", tokType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseCreateStatement() *ast.CreateStatement {
	//CREATE DATATYPE(INT/FLOAT/ETC.) IDENT SEMICOLON
	statement := &ast.CreateStatement{Token: p.curToken}

	//DATATYPE
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

	//IDENT
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, DataType: datatype}

	//SEMICOLON
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

func (p *Parser) parseSetStatement() *ast.SetStatement {
	//SET IDENT EQUALS EXPRESSION SEMICOLON
	statement := &ast.SetStatement{Token: p.curToken}

	//IDENT
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type == token.DOT {
		p.checkNextToken(token.DOT)
		if !p.checkNextToken(token.IDENT) {
			return nil
		}
		statement.Name.Attribute = p.curToken.Literal
	}

	if !p.checkNextToken(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.curToken}

	statement.Expression = p.parseExpression(LOWEST)

	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	statement := &ast.IfStatement{Token: p.curToken}

	if !p.checkNextToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	statement.Condition = p.parseExpression(LOWEST)

	if !p.checkNextToken(token.RPAEREN) {
		return nil
	}
	if !p.checkNextToken(token.BEGIN) {
		return nil
	}
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	for p.peekToken.Type != token.END {
		p.nextToken()
		statement.IfTrue = append(statement.IfTrue, p.parseStatement())
		if p.peekToken.Type == token.ELSE {
			p.errors = append(p.errors, "ELSE found before END;")
			return nil
		}
	}

	if !p.checkNextToken(token.END) {
		return nil
	}
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	if p.peekToken.Type != token.ELSE {
		statement.Else = nil
		return statement
	} else {
		p.nextToken()
		if !p.checkNextToken(token.BEGIN) {
			return nil
		}
		if !p.checkNextToken(token.SEMICOLON) {
			return nil
		}

		for p.peekToken.Type != token.END {
			p.nextToken()
			statement.Else = append(statement.Else, p.parseStatement())
			if p.peekToken.Type == token.EOF {
				p.errors = append(p.errors, "Missing END Token in if statement.")
				return nil
			}
		}

		if !p.checkNextToken(token.END) {
			return nil
		}
		if !p.checkNextToken(token.SEMICOLON) {
			return nil
		}
		return statement
	}
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	statement := &ast.WhileStatement{Token: p.curToken}

	if !p.checkNextToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	statement.Condition = p.parseExpression(LOWEST)

	if !p.checkNextToken(token.RPAEREN) {
		return nil
	}

	if !p.checkNextToken(token.BEGIN) {
		return nil
	}

	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	for p.peekToken.Type != token.END {
		p.nextToken()
		statement.LoopStatements = append(statement.LoopStatements, p.parseStatement())
		if p.peekToken.Type == token.EOF {
			p.errors = append(p.errors, "Missing END on while loop.")
			return nil
		}
	}

	if !p.checkNextToken(token.END) {
		return nil
	}

	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

func (p *Parser) parseStructStatement() *ast.StructStatement {
	statement := &ast.StructStatement{Token: p.curToken, Values: make(map[string]ast.Expression)}

	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	statement.StructIdent = ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, DataType: "struct"}

	if !p.checkNextToken(token.LPAREN) {
		return nil
	}
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
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, DataType: datatype}
	statement.Attributes = append(statement.Attributes, *ident)

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
	if !p.checkNextToken(token.RPAEREN) {
		return nil
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.checkNextToken(token.SEMICOLON)
		return statement
	}

	if !p.checkNextToken(token.LBRACKET) {
		return nil
	}
	if !p.checkNextToken(token.IDENT) {
		return nil
	}
	identString := p.curToken.Literal
	if !p.checkNextToken(token.COLON) {
		return nil
	}
	p.nextToken()
	statement.Values[identString] = p.parseExpression(LOWEST)

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

	if !p.checkNextToken(token.RBRACKET) {
		return nil
	}
	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	statement := &ast.PrintStatement{Token: p.curToken}

	if !p.checkNextToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	if !p.checkNextToken(token.RPAEREN) {
		return nil
	}

	if !p.checkNextToken(token.SEMICOLON) {
		return nil
	}

	return statement
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLON && precedence < getPrecedence(p.peekToken) {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekToken.Type == token.DOT {
		p.checkNextToken(token.DOT)
		if !p.checkNextToken(token.IDENT) {
			return nil
		}
		ident.Attribute = p.curToken.Literal
	}
	return ident
}

func (p *Parser) parseNumber() ast.Expression {
	if strings.Contains(p.curToken.Literal, ".") {
		return p.parseFloatLiteral()
	} else {
		return p.parseIntegerLiteral()
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.BooleanLiteral{Token: p.curToken}

	switch p.curToken.Literal {
	case "false":
		lit.Value = false
	case "true":
		lit.Value = true
	default:
		msg := fmt.Sprintf("could not parse %q as boolean", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	return lit
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
