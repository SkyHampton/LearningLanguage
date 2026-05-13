package parser

import (
	"fmt"
	"learningLanguage/ast"
	"learningLanguage/lexer"
	"testing"
)

func checkParserErrors(test *testing.T, p *Parser) {
	numErrors := len(p.errors)
	if numErrors == 0 {
		return
	}

	test.Errorf("Parser has %d errors.", numErrors)
	for _, msg := range p.errors {
		test.Errorf("Parser error: %q", msg)
	}
	test.FailNow()
}

func TestCreateStatements(test *testing.T) {
	input := `
	create int w; 
	create bool x; 
	create string y;
	create float z;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(test, p)
	if program == nil {
		test.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 4 {
		test.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedDataType   string
	}{
		{"w", "int"},
		{"x", "bool"},
		{"y", "string"},
		{"z", "float"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testCreateStatement(test, stmt, tt.expectedIdentifier, tt.expectedDataType) {
			return
		}
	}
}

func testCreateStatement(test *testing.T, statement ast.Statement, name string, dataType string) bool {
	if statement.TokenLiteral() != "create" {
		test.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}

	createStmt, ok := statement.(*ast.CreateStatement)
	if !ok {
		test.Errorf("s not *ast.LetStatement. got=%T", statement)
		return false
	}

	if createStmt.Name.Value != name {
		test.Errorf("createStmt.Name.Value not '%s'. got=%s", name, createStmt.Name.Value)
		return false
	}

	if createStmt.Name.TokenLiteral() != name {
		test.Errorf("createStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, createStmt.Name.TokenLiteral())
		return false
	}

	if createStmt.Name.DataType != dataType {
		test.Errorf("createStmt.DataType not '%s'. got %s", dataType, createStmt.Name.DataType)
		return false
	}

	return true
}

func TestSetStatements(test *testing.T) {
	input := `
	set w = 3; 
	set x = true; 
	set y = "Hello World";
	set z = 3.14;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(test, p)
	if program == nil {
		test.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 4 {
		test.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier        string
		expectedExpressionLiteral string
	}{
		{"w", "3"},
		{"x", "true"},
		{"y", "\"Hello World\""},
		{"z", "3.14"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testSetStatement(test, stmt, tt.expectedIdentifier, tt.expectedExpressionLiteral) {
			return
		}
	}
}

func testSetStatement(test *testing.T, statement ast.Statement, name string, explit string) bool {
	if statement.TokenLiteral() != "set" {
		test.Errorf("s.TokenLiteral not 'set'. got=%q", statement.TokenLiteral())
		return false
	}

	setStmt, ok := statement.(*ast.SetStatement)
	if !ok {
		test.Errorf("s not *ast.SetStatement. got=%T", statement)
		return false
	}

	if setStmt.Name.Value != name {
		test.Errorf("setStmt.Name.Value not '%s'. got=%s", name, setStmt.Name.Value)
		return false
	}

	if setStmt.Name.TokenLiteral() != name {
		test.Errorf("setStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, setStmt.Name.TokenLiteral())
		return false
	}

	if setStmt.Value.TokenLiteral() != explit {
		test.Errorf("setStmt.Value.TokenLiteral() not '%s'. got=%s", explit, setStmt.Value.TokenLiteral())
		return false
	}

	return true
}

func TestIfStatement(t *testing.T) {
	input := `if (true) begin; 
				set a = 1; end; 
				else begin; 
				set a = 0; end;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if stmt.TokenLiteral() != "if" {
		t.Errorf("stmt.TokenLiteral not if. got=%s", stmt.TokenLiteral())
	}

	cond, ok := stmt.Condition.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("stmt.Condition is not ast.BooleanLiteral. got=%T", stmt.Condition)
	}

	ifTrue, ok := stmt.IfTrue[0].(*ast.SetStatement)
	if !ok {
		t.Fatalf("stmt.IfTrue is not ast.SetStatement. got=%T", stmt.IfTrue[0])
	}

	els, ok := stmt.Else[0].(*ast.SetStatement)
	if !ok {
		t.Fatalf("stmt.Else is not ast.SetStatement. got=%T", stmt.Else[0])
	}

	if !cond.Value {
		t.Fatalf("cond.Value is not true. got=%t", cond.Value)
	}

	testIfSetStatement(t, ifTrue, "a", 1)

	testIfSetStatement(t, els, "a", 0)
}

func testIfSetStatement(t *testing.T, stmt *ast.SetStatement, name string, value int64) {

	if stmt.Name.Value != name {
		t.Fatalf("stmt.Name.Value is not a. got=%s", stmt.Name.Value)
		return
	}

	expression, ok := stmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.IntegerLiteral. got=%T", stmt.Value)
		return
	}

	if expression.Value != value {
		t.Fatalf("expression.Value is not 1. got=%d", expression.Value)
		return
	}
}

func TestWhileStatement(t *testing.T) {
	input := `while (true) begin; print("Hello World"); end;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T",
			program.Statements[0])
	}

	if stmt.TokenLiteral() != "while" {
		t.Errorf("stmt.TokenLiteral not while. got=%s", stmt.TokenLiteral())
	}

	cond, ok := stmt.Condition.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("stmt.Condition is not ast.BooleanLiteral. got=%T", stmt.Condition)
	}

	_, ok = stmt.LoopStatements[0].(*ast.PrintStatement)
	if !ok {
		t.Fatalf("stmt.LoopStatements[0] is not ast.PrintStatement. got=%T", stmt.LoopStatements[0])
	}

	if !cond.Value {
		t.Fatalf("cond.Value is not true. got=%t", cond.Value)
	}
}

func TestCountStatement(t *testing.T) {
	input := `count i from 1 to 10 begin; print(i); end;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.CountStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.CountStatement. got=%T",
			program.Statements[0])
	}

	if stmt.TokenLiteral() != "count" {
		t.Errorf("stmt.TokenLiteral not while. got=%s", stmt.TokenLiteral())
	}

	if stmt.Counter.Value != "i" {
		t.Errorf("stnt.Counter.Value not i, got=%s", stmt.Counter.Value)
	}

	from, ok := stmt.From.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.From is not ast.IntegerLiteral. got=%T",
			stmt.From)
	}

	if from.Value != 1 {
		t.Fatalf("from.Value is not 1, got %d", from.Value)
	}

	to, ok := stmt.To.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.To is not ast.IntegerLiteral. got=%T",
			stmt.From)
	}

	if to.Value != 10 {
		t.Fatalf("to.Value is not 10, got %d", to.Value)
	}

	by, ok := stmt.By.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.From is not ast.IntegerLiteral. got=%T",
			stmt.From)
	}

	if by.Value != 1 {
		t.Fatalf("by.Value is not 1, got %d", by.Value)
	}

	if from.Value != 1 {
		t.Fatalf("from.Value is not 1, got %d", from.Value)
	}

	_, ok = stmt.LoopStatements[0].(*ast.PrintStatement)
	if !ok {
		t.Fatalf("stmt.LoopStatements[0] is not ast.PrintStatement. got=%T", stmt.LoopStatements[0])
	}

}

func TestStructStatement(t *testing.T) {
	input := `struct myStruct(
				int a,
				bool b
			) [a: 123, b: true];`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.StructStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	name := stmt.StructIdent
	if name.Value != "myStruct" {
		t.Fatalf("Struct name is not myStruct, got %s", name.Value)
	}

	attributes := stmt.Attributes

	if len(attributes) != 2 {
		t.Fatalf("Too few struct attributes, expected 2, got %d", len(attributes))
	}

	if attributes[0].Value != "a" && attributes[1].Value != "b" {
		t.Fatalf("Incorrect struct attribute names, expected 'a' and 'b', got %s, %s", attributes[0].Value, attributes[1].Value)
	}

	values := stmt.Values

	aInt, ok := values[attributes[0].Value].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("a's value is not an intlit, got %T", aInt)
	}

	bBool, ok := values[attributes[1].Value].(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("b's value is not an boolLit, got %T", bBool)
	}

	testIntegerLiteral(t, aInt, 123)
	testBooleanLiteral(t, bBool, true)
}

func TestAltStructStatement(t *testing.T) {
	input := `struct myStruct(int a, bool b);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.StructStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	name := stmt.StructIdent
	if name.Value != "myStruct" {
		t.Fatalf("Struct name is not myStruct, got %s", name.Value)
	}

	attributes := stmt.Attributes

	if len(attributes) != 2 {
		t.Fatalf("Too few struct attributes, expected 2, got %d", len(attributes))
	}

	if attributes[0].Value != "a" && attributes[1].Value != "b" {
		t.Fatalf("Incorrect struct attribute names, expected 'a' and 'b', got %s, %s", attributes[0].Value, attributes[1].Value)
	}
}

func TestAttributeSetStatement(t *testing.T) {
	input := "set myStruct.a = 123;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.SetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "myStruct" {
		t.Fatalf("Expected struct name myStruct, got %s", stmt.Name.Value)
	}

	if stmt.Name.Attribute != "a" {
		t.Fatalf("Expected struct attrubute a, got %s", stmt.Name.Attribute)
	}

	intLit, ok := stmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.IntLiteral. got=%T",
			stmt.Value)
	}

	testIntegerLiteral(t, intLit, 123)
}

func TestAttributeIdentifier(t *testing.T) {
	input := "myStruct.a;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "myStruct" {
		t.Fatalf("Struct identifier not myStruct, got %s", ident.Value)
	}

	if ident.Attribute != "a" {
		t.Fatalf("Struct attribute not a, got %s", ident.Attribute)
	}
}

func TestPrintStatement(t *testing.T) {
	input := "print(\"Hello World\");println(\"Hello World\");"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.PrintStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.PrintStatement. got=%T",
			program.Statements[0])
	}

	if stmt.TokenLiteral() != "print" {
		t.Errorf("stmt.TokenLiteral not %s. got=%s", "print", stmt.TokenLiteral())
	}

	strLit, ok := stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.StringLiteral. got=%T", stmt.Value)
	}

	if strLit.Value != "\"Hello World\"" {
		t.Fatalf("strLit.Value is not \"Hello World\". got=%s", strLit.Value)
	}

	stmt, ok = program.Statements[1].(*ast.PrintStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.PrintStatement. got=%T",
			program.Statements[1])
	}

	if stmt.TokenLiteral() != "println" {
		t.Errorf("stmt.TokenLiteral not %s. got=%s", "print", stmt.TokenLiteral())
	}

	strLit, ok = stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.StringLiteral. got=%T", stmt.Value)
	}

	if strLit.Value != "\"Hello World\"" {
		t.Fatalf("strLit.Value is not \"Hello World\". got=%s", strLit.Value)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;foobar.a;foobar[1];"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	attrStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}

	attrIdent, ok := attrStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("attrIdent.Expression not *ast.Identifier. got=%T", attrStmt.Expression)
	}

	if attrIdent.Value != "foobar" {
		t.Errorf("attrIdent.Value not %s. got=%s", "foobar", attrIdent.Value)
	}
	if attrIdent.Attribute != "a" {
		t.Errorf("attrIdent.Attribute not %s. got=%s", "a", attrIdent.Attribute)
	}
	if attrIdent.TokenLiteral() != "foobar" {
		t.Errorf("attrIdent.TokenLiteral not %s. got=%s", "foobar",
			attrIdent.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "2;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 2 {
		t.Errorf("literal.Value not %d. got=%d", 2, literal.Value)
	}
	if literal.TokenLiteral() != "2" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "2",
			literal.TokenLiteral())
	}
}

func TestBooleanLiteralExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("exp not *ast.BooleanLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != true {
		t.Errorf("literal.Value not %t. got=%t", true, literal.Value)
	}
	if literal.TokenLiteral() != "true" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "true",
			literal.TokenLiteral())
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := "\"Hello World\";"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "\"Hello World\"" {
		t.Errorf("literal.Value not %s. got=%s", "\"Hello World\"", literal.Value)
	}
	if literal.TokenLiteral() != "\"Hello World\"" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "\"Hello World\"",
			literal.TokenLiteral())
	}
}

func TestFloatLiteralExpression(t *testing.T) {
	input := "3.14;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FloatLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 3.14 {
		t.Errorf("literal.Value not %f. got=%f", 3.14, literal.Value)
	}
	if literal.TokenLiteral() != "3.14" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "3.14",
			literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixNegTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"-22;", "-", 22},
	}

	prefixBangTests := []struct {
		input    string
		operator string
		boolean  bool
	}{
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixNegTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				2, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}

	for _, tt := range prefixBangTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				2, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testBooleanLiteral(t, exp.Right, tt.boolean) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, il ast.Expression, value bool) bool {
	boolean, ok := il.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if boolean.Value != value {
		t.Errorf("integ.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("integ.TokenLiteral not %t. got=%s", value,
			boolean.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixNumTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"2+2;", 2, "+", 2},
		{"2-2;", 2, "-", 2},
		{"2*2;", 2, "*", 2},
		{"2/2;", 2, "/", 2},
		{"2>1;", 2, ">", 1},
		{"2>=1;", 2, ">=", 1},
		{"2<1;", 2, "<", 1},
		{"2<=1;", 2, "<=", 1},
		{"2==1;", 2, "==", 1},
		{"2!=1;", 2, "!=", 1},
	}

	for _, test := range infixNumTests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				2, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", statement.Expression)
		}

		if !testIntegerLiteral(t, expression.Left, test.leftValue) {
			return
		}

		if expression.Operator != test.operator {
			t.Fatalf("Operator is not %s. Got %s", test.operator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Right, test.rightValue) {
			return
		}
	}

	infixBoolTests := []struct {
		input      string
		leftValue  bool
		operator   string
		rightValue bool
	}{
		{"true and false;", true, "and", false},
		{"false or true;", false, "or", true},
	}

	for _, test := range infixBoolTests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				2, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", statement.Expression)
		}

		if !testBooleanLiteral(t, expression.Left, test.leftValue) {
			return
		}

		if expression.Operator != test.operator {
			t.Fatalf("Operator is not %s. Got %s", test.operator, expression.Operator)
		}

		if !testBooleanLiteral(t, expression.Right, test.rightValue) {
			return
		}
	}
}

func TestLists(t *testing.T) {
	input := "create int list myList; set myList = [1, 2, 3]; set myList[1] = 100; append 4 to myList; len(myList);"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 5 {
		t.Errorf("Expected %d statements, got %d", 5, len(program.Statements))
	}

	// create int list myList;
	createStmt, ok := program.Statements[0].(*ast.CreateStatement)
	if !ok {
		t.Errorf("TestLists first statement not CreateStatment, got %T", program.Statements[0])
	}

	if !createStmt.Name.IsList {
		t.Error("createStmt.Name is not a list")
	}

	// set myList = [1, 2, 3];
	setStmt, ok := program.Statements[1].(*ast.SetStatement)
	if !ok {
		t.Errorf("TestLists second statement not SetStatment, got %T", program.Statements[1])
	}

	listLit, ok := setStmt.Value.(*ast.ListLiteral)
	if !ok {
		t.Errorf("setStmt.Value is not ListLiteral, got %T", setStmt.Value)
	}

	values := []int64{1, 2, 3}
	if len(values) != len(listLit.List) {
		t.Errorf("Not enough list elements, expected %d, got %d", len(values), len(listLit.List))
	}
	for i, _ := range listLit.List {
		intLit, ok := listLit.List[i].(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("List element %s not IntegerLiteral, got %T", listLit.List[i].String(), listLit.List[i])
		}

		testIntegerLiteral(t, intLit, values[i])
	}

	// set myList[1] = 100;
	setStmt, ok = program.Statements[2].(*ast.SetStatement)
	if !ok {
		t.Errorf("TestLists third statement not SetStatment, got %T", program.Statements[2])
	}

	testIntegerLiteral(t, setStmt.Name.Index, 1)
	testIntegerLiteral(t, setStmt.Value, 100)

	// append 4 to myList;
	appendStmt, ok := program.Statements[3].(*ast.AppendStatement)
	if !ok {
		t.Errorf("TestLists fourth statement not AppendStatement, got %T", program.Statements[3])
	}

	testIntegerLiteral(t, appendStmt.Value, 4)

	if appendStmt.List.Value != "myList" {
		t.Errorf("appendStmt.List.Value is not myList, got %s", appendStmt.List.Value)
	}
	// len(myList);
	expStmt, ok := program.Statements[4].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("TestLists fifth statement not ExpressionStatement, got %T", program.Statements[4])
	}

	lenExp, ok := expStmt.Expression.(*ast.LengthExpression)
	if !ok {
		t.Errorf("expStmt.Expression not LengthExpression, got %T", expStmt.Expression)
	}

	if lenExp.List.Value != "myList" {
		t.Errorf("lenExp.List.Value is not myList, got %s", lenExp.List.Value)
	}
}
