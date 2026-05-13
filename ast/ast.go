package ast

import (
	"bytes"
	"fmt"
	"learningLanguage/token"
)

// Interfaces
type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	StatementNode()
}

type Expression interface {
	Node
	ExpressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

// Variable Creation Statement
type CreateStatement struct {
	Token token.Token
	Name  *Identifier
}

func (cs *CreateStatement) StatementNode()       {}
func (cs *CreateStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *CreateStatement) String() string {
	return cs.Name.String()
}

// Variable Set Statement
type SetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ss *SetStatement) StatementNode()       {}
func (ss *SetStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SetStatement) String() string {
	var out bytes.Buffer

	out.WriteString("Identifier Name: ")
	out.WriteString(ss.Name.String())
	out.WriteString(". Expression: ")
	out.WriteString(ss.Value.String())
	out.WriteString(".")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) StatementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

type IfStatement struct {
	Token     token.Token
	Condition Expression
	IfTrue    []Statement
	Else      []Statement
}

func (is *IfStatement) StatementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("if (")
	out.WriteString(is.Condition.String())
	out.WriteString("):\n")
	for _, stmt := range is.IfTrue {
		out.WriteString(stmt.String())
	}
	out.WriteString("\nelse:\n")
	for _, stmt := range is.IfTrue {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type WhileStatement struct {
	Token          token.Token
	Condition      Expression
	LoopStatements []Statement
}

func (ws *WhileStatement) StatementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while (")
	out.WriteString(ws.Condition.String())
	out.WriteString("):\n")
	for _, stmt := range ws.LoopStatements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type CountStatement struct {
	Token          token.Token
	Counter        Identifier
	From           Expression
	To             Expression
	By             Expression
	LoopStatements []Statement
}

func (cs *CountStatement) StatementNode()       {}
func (cs *CountStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *CountStatement) String() string {
	var out bytes.Buffer

	out.WriteString("count ")
	out.WriteString(cs.Counter.Value)
	out.WriteString(" from ")
	out.WriteString(cs.From.String())
	out.WriteString(" to ")
	out.WriteString(cs.To.String())
	out.WriteString(" by ")
	out.WriteString(cs.By.String())
	out.WriteString(":\n")
	for _, stmt := range cs.LoopStatements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type StructStatement struct {
	Token       token.Token
	StructIdent Identifier
	Attributes  []Identifier
	Values      map[string]Expression
}

func (ss *StructStatement) StatementNode()       {}
func (ss *StructStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *StructStatement) String() string {
	var out bytes.Buffer
	out.WriteString("Struct: ")
	out.WriteString(ss.StructIdent.String())

	for _, attribute := range ss.Attributes {
		out.WriteString("Attribute: ")
		out.WriteString(attribute.String())
		out.WriteString(". Value: ")
		out.WriteString(ss.Values[attribute.Value].String())
		out.WriteString(".")
	}

	return out.String()
}

type PrintStatement struct {
	Token   token.Token
	Value   Expression
	NewLine bool
}

func (ps *PrintStatement) StatementNode()       {}
func (ps *PrintStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PrintStatement) String() string {
	if ps.NewLine {
		return "Println(" + ps.Value.String() + ")"
	} else {
		return "Print(" + ps.Value.String() + ")"
	}
}

type AppendStatement struct {
	Token token.Token
	List  Identifier
	Value Expression
}

func (as *AppendStatement) StatementNode()       {}
func (as *AppendStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AppendStatement) String() string {
	return fmt.Sprintf("Append %s to list %s", as.Value.String(), as.List.String())
}

type Identifier struct {
	Token     token.Token
	Value     string
	DataType  string
	Attribute string
	IsList    bool
	Index     Expression
}

func (i *Identifier) ExpressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	var out bytes.Buffer

	if i.IsList {
		out.WriteString("List ")
	}

	if i.Attribute != "" {
		out.WriteString(fmt.Sprintf("Variable %s.%s", i.Value, i.Attribute))
	} else {
		out.WriteString(fmt.Sprintf("Variable %s", i.Value))
	}

	if i.Index != nil {
		out.WriteString(fmt.Sprintf("[%s] ", i.Index.String()))
	}
	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) ExpressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) ExpressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) ExpressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) ExpressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ListLiteral struct {
	Token token.Token
	List  []Expression
}

func (ll *ListLiteral) ExpressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }
func (ll *ListLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for _, exp := range ll.List {
		out.WriteString(fmt.Sprintf("%s, ", exp.String()))
	}
	out.WriteString("]")
	return out.String()
}

type LengthExpression struct {
	Token token.Token
	List  Identifier
}

func (le *LengthExpression) ExpressionNode()      {}
func (le *LengthExpression) TokenLiteral() string { return le.Token.Literal }
func (le *LengthExpression) String() string {
	return fmt.Sprintf("len(%s)", le.List.String())
}

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) ExpressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) ExpressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
