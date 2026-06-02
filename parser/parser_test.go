// TODO this test file likely needs to be reworked, the tests are too complex and don't actually test well
package parser

import (
	"learningLanguage/ast"
	"learningLanguage/lexer"
	"testing"
)

func testErrors(test *testing.T, input string, expectedErrors bool) (*ast.Program, bool) {
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	errors := parser.Errors()
	numErrors := len(errors)

	if numErrors > 0 && !expectedErrors {
		test.Errorf("Input: %s\nUnexpected erors found: %v", input, errors)
		return nil, true
	}

	if numErrors == 0 && expectedErrors {
		test.Error("Expected errors but encountered none.")
		return nil, true
	}

	if numErrors > 0 {
		return nil, true
	}
	return program, false
}

func TestCreateStatements(Test *testing.T) {
	createTests := []struct {
		input              string
		expectedToken      string
		expectedErrors     bool
		expectedIdentifier string
		expectedDataType   string
	}{
		{"create int x;", "create", false, "x", "int"},
		{"create bool x;", "create", false, "x", "bool"},
		{"create float x;", "create", false, "x", "float"},
		{"create string x;", "create", false, "x", "string"},
		{"create x;", "create", true, "x", "string"},
		{"string x;", "create", true, "x", "string"},
		{"create string;", "create", true, "x", "string"},
		{"create string x", "create", true, "x", "string"},
	}

	for _, test := range createTests {
		testCreateStatement(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedIdentifier, test.expectedDataType)
	}
}

func testCreateStatement(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedIdentifier string, expectedDataType string) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	createStmt, ok := statement.(*ast.CreateStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected CreateStatement, got %T", input, statement)
		return
	}

	if createStmt.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, createStmt.TokenLiteral())
		return
	}

	if createStmt.Name.Value != expectedIdentifier {
		test.Errorf("Input: %s\nExpected identifier %s, got %s", input, expectedIdentifier, createStmt.Name.Value)
		return
	}

	if createStmt.Name.DataType != expectedDataType {
		test.Errorf("Input: %s\nExpected data type %s, got %s", input, expectedDataType, createStmt.Name.DataType)
		return
	}
}

func TestSetStatements(Test *testing.T) {
	setTests := []struct {
		input              string
		expectedToken      string
		expectedErrors     bool
		expectedIdentifier string
		expectedValue      int32
	}{
		{"set x = 1;", "set", false, "x", 1},
		{"set = 1;", "set", true, "x", 1},
		{"set x 1;", "set", true, "x", 1},
		{"set x = ;", "set", true, "x", 1},
		{"set x = 1", "set", true, "x", 1},
	}

	for _, test := range setTests {
		testSetStatement(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedIdentifier, test.expectedValue)
	}
}

func testSetStatement(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedIdentifier string, expectedValue int32) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	setStmt, ok := statement.(*ast.SetStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected SetStatement, got %T", input, statement)
		return
	}

	if setStmt.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, setStmt.TokenLiteral())
		return
	}

	if setStmt.Name.Value != expectedIdentifier {
		test.Errorf("Input: %s\nExpected identifier %s, got %s", input, expectedIdentifier, setStmt.Name.Value)
		return
	}

	intLit, ok := setStmt.Value.(*ast.IntegerLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected IntegerLiteral, got %T", input, setStmt.Value)
		return
	}
	if intLit.Value != expectedValue {
		test.Errorf("Input: %s\nExpected %d, got %d", input, expectedValue, intLit.Value)
	}
}

func TestIfStatements(Test *testing.T) {
	ifTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
	}{
		{"if (true) begin; set x = 1; end; else begin; set x = 0; end;", "if", false},
		{"if (true) begin; set x = 1; end;", "if", false},
		{"if true) begin; set x = 1; end; else begin; set x = 0; end;", "if", true},
		{"if (true begin; set x = 1; end; else begin; set x = 0; end;", "if", true},
		{"if (true) ; set x = 1; end; else begin; set x = 0; end;", "if", true},
		{"if (true) begin set x = 1; end; else begin; set x = 0; end;", "if", true},
		{"if (true) begin; set x = 1; ; else begin; set x = 0; end;", "if", true},
		{"if (true) begin; set x = 1; end else begin; set x = 0; end;", "if", true},
		{"if (true) begin; set x = 1; end; begin; set x = 0; end;", "if", true},
		{"if (true) begin; set x = 1; end else begin; set x = 0; end;", "if", true},
		{"if (true) begin; set x = 1; end; begin; set x = 0; end;", "if", true},
		{"if (true) begin; set x = 1; end; else ; set x = 0; end;", "if", true},
		{"if (true) begin; set x = 1; end; else begin set x = 0; end;", "if", true},
		{"if (true) begin; set x = 1; end; else begin; set x = 0; ;", "if", true},
		{"if (true) begin; set x = 1; end; else begin; set x = 0; end", "if", true},
	}

	for _, test := range ifTests {
		testIfStatement(Test, test.input, test.expectedToken, test.expectedErrors)
	}
}

func testIfStatement(test *testing.T, input string, expectedToken string, expectedErrors bool) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	ifStmt, ok := statement.(*ast.IfStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected IfStatement, got %T", input, statement)
		return
	}

	if ifStmt.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, ifStmt.TokenLiteral())
		return
	}
}

func TestWhileStatements(Test *testing.T) {
	whileTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
	}{
		{"while (true) begin; set x = 1; end;", "while", false},
		{"while true) begin; set x = 1; end;", "while", true},
		{"while () begin; set x = 1; end;", "while", true},
		{"while (true begin; set x = 1; end;", "while", true},
		{"while (true) ; set x = 1; end;", "while", true},
		{"while (true) begin set x = 1; end;", "while", true},
		{"while (true) begin; set x = 1; ;", "while", true},
		{"while (true) begin; set x = 1; end", "while", true},
	}

	for _, test := range whileTests {
		testWhileStatement(Test, test.input, test.expectedToken, test.expectedErrors)
	}
}

func testWhileStatement(test *testing.T, input string, expectedToken string, expectedErrors bool) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	whileStmt, ok := statement.(*ast.WhileStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected WhileStatement, got %T", input, statement)
		return
	}

	if whileStmt.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, whileStmt.TokenLiteral())
		return
	}
}

func TestCountStatements(Test *testing.T) {
	countTests := []struct {
		input              string
		expectedToken      string
		expectedErrors     bool
		expectedIdentifier string
		expectedFrom       int32
		expectedTo         int32
		expectedBy         int32
	}{
		{"count i from 1 to 10 begin; print(i); end;", "count", false, "i", 1, 10, 1},
		{"count i from 1 to 10 by 2 begin; print(i); end;", "count", false, "i", 1, 10, 2},
		{"count from 1 to 10 begin; print(i); end;", "count", true, "i", 1, 10, 1},
		{"count i 1 to 10 begin; print(i); end;", "count", true, "i", 1, 10, 1},
		{"count i from to 10 begin; print(i); end;", "count", true, "i", 1, 10, 1},
		{"count i from 1 10 begin; print(i); end;", "count", true, "i", 1, 10, 1},
		{"count i from 1 to 10; print(i); end;", "count", true, "i", 1, 10, 1},
		{"count i from 1 to 10 begin print(i); end;", "count", true, "i", 1, 10, 1},
		{"count i from 1 to 10 begin; print(i); ;", "count", true, "i", 1, 10, 1},
		{"count i from 1 to 10 begin; print(i); end", "count", true, "i", 1, 10, 1},
	}

	for _, test := range countTests {
		testCountStatement(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedIdentifier, test.expectedFrom, test.expectedTo, test.expectedBy)
	}
}

func testCountStatement(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedIdentifier string, expectedFrom int32, expectedTo int32, expectedBy int32) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	countStmt, ok := statement.(*ast.CountStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected CountStatement, got %T", input, statement)
		return
	}

	if countStmt.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, countStmt.TokenLiteral())
		return
	}

	if countStmt.Counter.Value != expectedIdentifier {
		test.Errorf("Input: %s\nExpected identifier %s, got %s", input, expectedIdentifier, countStmt.Counter.Value)
		return
	}

	fromIntLit, ok := countStmt.From.(*ast.IntegerLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected IntegerLiteral, got %T", input, countStmt.From)
		return
	}
	if fromIntLit.Value != expectedFrom {
		test.Errorf("Input: %s\nExpected %d, got %d", input, expectedFrom, fromIntLit.Value)
	}

	toIntLit, ok := countStmt.To.(*ast.IntegerLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected IntegerLiteral, got %T", input, countStmt.To)
		return
	}
	if toIntLit.Value != expectedTo {
		test.Errorf("Input: %s\nExpected %d, got %d", input, expectedTo, toIntLit.Value)
	}

	byIntLit, ok := countStmt.By.(*ast.IntegerLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected IntegerLiteral, got %T", input, countStmt.By)
		return
	}
	if byIntLit.Value != expectedBy {
		test.Errorf("Input: %s\nExpected %d, got %d", input, expectedBy, byIntLit.Value)
	}
}

func TestStructStatements(Test *testing.T) {
	structTests := []struct {
		input              string
		expectedToken      string
		expectedErrors     bool
		expectedIdentifier string
		expectedAttributes []string
		expectedDataTypes  []string
	}{
		{"struct a (int b, bool c)[b: 123, c: false];", "struct", false, "a", []string{"b", "c"}, []string{"int", "bool"}},
		{"struct a (int b, bool c);", "struct", false, "a", []string{"b", "c"}, []string{"int", "bool"}},
		{"struct a (float b, string c)[b: 123, c: false];", "struct", false, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct a (float b, string c);", "struct", false, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct (float b, string c);", "struct", true, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct a float b, string c);", "struct", true, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct a ( b, string c);", "struct", true, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct a (float , string c);", "struct", true, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct a (float b string c);", "struct", true, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct a (float b, c);", "struct", true, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct a (float b, string );", "struct", true, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct a (float b, string c;", "struct", true, "a", []string{"b", "c"}, []string{"float", "string"}},
		{"struct a (float b, string c)", "struct", true, "a", []string{"b", "c"}, []string{"float", "string"}},
	}

	for _, test := range structTests {
		testStructStatement(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedIdentifier, test.expectedAttributes, test.expectedDataTypes)
	}
}

func testStructStatement(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedIdentifier string, expectedAttributes []string, expectedDataTypes []string) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	structStmt, ok := statement.(*ast.StructStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected StructStatement, got %T", input, statement)
		return
	}

	if structStmt.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, structStmt.TokenLiteral())
		return
	}

	if structStmt.StructIdent.Value != expectedIdentifier {
		test.Errorf("Input: %s\nExpected identifier %s, got %s", input, expectedIdentifier, structStmt.StructIdent.Value)
		return
	}

	for i := range structStmt.Attributes {
		if structStmt.Attributes[i].Value != expectedAttributes[i] {
			test.Errorf("Input: %s\nExpected attribute %s, got %s", input, expectedAttributes[i], structStmt.Attributes[i].Value)
			return
		}

		if structStmt.Attributes[i].DataType != expectedDataTypes[i] {
			test.Errorf("Input: %s\nExpected datatype %s, got %s", input, expectedDataTypes[i], structStmt.Attributes[i].DataType)
			return
		}
	}
}

func TestPrintStatements(Test *testing.T) {
	printTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
	}{
		{"print(123);", "print", false},
		{"print(false);", "print", false},
		{"print(3.14);", "print", false},
		{"print(\"Hello World!\");", "print", false},
		{"print(123 - 3);", "print", false},
		{"print(!false);", "print", false},
		{"print123);", "print", true},
		{"print(123;", "print", true},
		{"print(123)", "print", true},
	}

	for _, test := range printTests {
		testPrintStatement(Test, test.input, test.expectedToken, test.expectedErrors)
	}
}

func testPrintStatement(test *testing.T, input string, expectedToken string, expectedErrors bool) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	printStmt, ok := statement.(*ast.PrintStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected PrintStatement, got %T", input, statement)
		return
	}

	if printStmt.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, printStmt.TokenLiteral())
		return
	}
}

func TestAppendStatements(Test *testing.T) {
	appendTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
	}{
		{"append 4 to myList;", "append", false},
		{"append to myList;", "append", true},
		{"append 4 myList;", "append", true},
		{"append 4 to;", "append", true},
		{"append 4 to myList", "append", true},
	}

	for _, test := range appendTests {
		testAppendStatement(Test, test.input, test.expectedToken, test.expectedErrors)
	}
}

func testAppendStatement(test *testing.T, input string, expectedToken string, expectedErrors bool) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	appendStmt, ok := statement.(*ast.AppendStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected AppendStatement, got %T", input, statement)
		return
	}

	if appendStmt.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, appendStmt.TokenLiteral())
		return
	}
}

func TestIntLiterals(Test *testing.T) {
	intLitTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
		expectedValue  int32
	}{
		{"123;", "123", false, 123},
		{"123456789;", "123456789", false, 123456789},
		{"987654321;", "987654321", false, 987654321},
	}

	for _, test := range intLitTests {
		testIntLiteral(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedValue)
	}
}

func testIntLiteral(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedValue int32) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	exprStmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected 'ExpressionStatement', got %T", input, statement)
		return
	}

	intLit, ok := exprStmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected 'IntegerLiteral', got %T", input, exprStmt.Expression)
		return
	}

	if intLit.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, intLit.TokenLiteral())
		return
	}

	if intLit.Value != expectedValue {
		test.Errorf("Input: %s\nExpected value %d, got %d", input, expectedValue, intLit.Value)
		return
	}
}

func TestBoolLiterals(Test *testing.T) {
	boolLitTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
		expectedValue  bool
	}{
		{"true;", "true", false, true},
		{"false;", "false", false, false},
	}

	for _, test := range boolLitTests {
		testBoolLiteral(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedValue)
	}
}

func testBoolLiteral(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedValue bool) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	exprStmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected 'ExpressionStatement', got %T", input, statement)
		return
	}

	boolLit, ok := exprStmt.Expression.(*ast.BooleanLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected 'BoolenLiteral', got %T", input, exprStmt.Expression)
		return
	}

	if boolLit.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, boolLit.TokenLiteral())
		return
	}

	if boolLit.Value != expectedValue {
		test.Errorf("Input: %s\nExpected value %t, got %t", input, expectedValue, boolLit.Value)
		return
	}
}

func TestStringLiterals(Test *testing.T) {
	stringLitTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
		expectedValue  string
	}{
		{"\"qwerty\";", "\"qwerty\"", false, "\"qwerty\""},
		{"\"qwerty asdf\";", "\"qwerty asdf\"", false, "\"qwerty asdf\""},
		{"\"qwerty;", "\"qwerty", true, "querty"},
	}

	for _, test := range stringLitTests {
		testStringLiteral(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedValue)
	}
}

func testStringLiteral(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedValue string) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	exprStmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected 'ExpressionStatement', got %T", input, statement)
		return
	}

	stringLit, ok := exprStmt.Expression.(*ast.StringLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected 'StringLiteral', got %T", input, exprStmt.Expression)
		return
	}

	if stringLit.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, stringLit.TokenLiteral())
		return
	}

	if stringLit.Value != expectedValue {
		test.Errorf("Input: %s\nExpected value %s, got %s", input, expectedValue, stringLit.Value)
		return
	}
}

func TestFloatLiterals(Test *testing.T) {
	floatLitTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
		expectedValue  float32
	}{
		{"3.14;", "3.14", false, 3.14},
		{"123.456;", "123.456", false, 123.456},
	}

	for _, test := range floatLitTests {
		testFloatLiteral(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedValue)
	}
}

func testFloatLiteral(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedValue float32) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	exprStmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected 'ExpressionStatement', got %T", input, statement)
		return
	}

	floatLit, ok := exprStmt.Expression.(*ast.FloatLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected 'FloatLiteral', got %T", input, exprStmt.Expression)
		return
	}

	if floatLit.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, floatLit.TokenLiteral())
		return
	}

	if floatLit.Value != expectedValue {
		test.Errorf("Input: %s\nExpected value %f, got %f", input, expectedValue, floatLit.Value)
		return
	}
}

func TestLengthExpressions(Test *testing.T) {
	lengthTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
	}{
		{"len(a);", "len", false},
		{"lena);", "len", true},
		{"len();", "len", true},
		{"len(a;", "len", true},
		{"len(a)", "len", true},
	}

	for _, test := range lengthTests {
		testLengthExpression(Test, test.input, test.expectedToken, test.expectedErrors)
	}
}

func testLengthExpression(test *testing.T, input string, expectedToken string, expectedErrors bool) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	exprStmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected 'ExpressionStatement', got %T", input, statement)
		return
	}

	lengthExpr, ok := exprStmt.Expression.(*ast.LengthExpression)
	if !ok {
		test.Errorf("Input: %s\nExpected 'LengthExpression', got %T", input, exprStmt.Expression)
		return
	}

	if lengthExpr.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, lengthExpr.TokenLiteral())
		return
	}
}

func TestPrefixExpressions(Test *testing.T) {
	prefixTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
	}{
		{"-123;", "-", false},
		{"!true;", "!", false},
	}

	for _, test := range prefixTests {
		testPrefixExpression(Test, test.input, test.expectedToken, test.expectedErrors)
	}
}

func testPrefixExpression(test *testing.T, input string, expectedToken string, expectedErrors bool) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	exprStmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected 'ExpressionStatement', got %T", input, statement)
		return
	}

	prefixExpr, ok := exprStmt.Expression.(*ast.PrefixExpression)
	if !ok {
		test.Errorf("Input: %s\nExpected 'PrefixExpression', got %T", input, exprStmt.Expression)
		return
	}

	if prefixExpr.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, prefixExpr.TokenLiteral())
		return
	}
}

func TestArithInfixExpressions(Test *testing.T) {
	infixArithTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
		expectedLeft   int32
		expectedRight  int32
	}{
		{"1+1;", "+", false, 1, 1},
		{"1-1;", "-", false, 1, 1},
		{"1/1;", "/", false, 1, 1},
		{"1*1;", "*", false, 1, 1},
		{"1==1;", "==", false, 1, 1},
		{"1>1;", ">", false, 1, 1},
		{"1>=1;", ">=", false, 1, 1},
		{"1<1;", "<", false, 1, 1},
		{"1<=1;", "<=", false, 1, 1},
	}

	for _, test := range infixArithTests {
		testArithInfixExpression(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedLeft, test.expectedRight)
	}
}

func testArithInfixExpression(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedLeft int32, expectedRight int32) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	exprStmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected 'ExpressionStatement', got %T", input, statement)
		return
	}

	infixExpr, ok := exprStmt.Expression.(*ast.InfixExpression)
	if !ok {
		test.Errorf("Input: %s\nExpected 'InfixExpression', got %T", input, exprStmt.Expression)
		return
	}

	if infixExpr.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, infixExpr.TokenLiteral())
		return
	}

	left, ok := infixExpr.Left.(*ast.IntegerLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected 'IntegerLiteral', got %T", input, infixExpr.Left)
		return
	}

	right, ok := infixExpr.Right.(*ast.IntegerLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected 'IntegerLiteral', got %T", input, infixExpr.Right)
		return
	}

	if left.Value != expectedLeft {
		test.Errorf("Input: %s\nExpected left %d, got %d", input, expectedLeft, left.Value)
		return
	}

	if right.Value != expectedRight {
		test.Errorf("Input: %s\nExpected right %d, got %d", input, expectedRight, right.Value)
		return
	}
}

func TestLogicInfixExpressions(Test *testing.T) {
	infixLogicTests := []struct {
		input          string
		expectedToken  string
		expectedErrors bool
		expectedLeft   bool
		expectedRight  bool
	}{
		{"true and false;", "and", false, true, false},
		{"true or false;", "or", false, true, false},
	}

	for _, test := range infixLogicTests {
		testLogicInfixExpression(Test, test.input, test.expectedToken, test.expectedErrors, test.expectedLeft, test.expectedRight)
	}
}

func testLogicInfixExpression(test *testing.T, input string, expectedToken string, expectedErrors bool, expectedLeft bool, expectedRight bool) {
	program, shouldReturn := testErrors(test, input, expectedErrors)
	if shouldReturn {
		return
	}

	if len(program.Statements) == 0 {
		test.Errorf("Input: %s\nNo statements found", input)
		return
	}
	statement := program.Statements[0]
	exprStmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		test.Errorf("Input: %s\nExpected 'ExpressionStatement', got %T", input, statement)
		return
	}

	infixExpr, ok := exprStmt.Expression.(*ast.InfixExpression)
	if !ok {
		test.Errorf("Input: %s\nExpected 'InfixExpression', got %T", input, exprStmt.Expression)
		return
	}

	if infixExpr.TokenLiteral() != expectedToken {
		test.Errorf("Input: %s\nExpected token literal %s, got %s", input, expectedToken, infixExpr.TokenLiteral())
		return
	}

	left, ok := infixExpr.Left.(*ast.BooleanLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected 'BooleanLiteral', got %T", input, infixExpr.Left)
		return
	}

	right, ok := infixExpr.Right.(*ast.BooleanLiteral)
	if !ok {
		test.Errorf("Input: %s\nExpected 'BooleanLiteral', got %T", input, infixExpr.Right)
		return
	}

	if left.Value != expectedLeft {
		test.Errorf("Input: %s\nExpected left %t, got %t", input, expectedLeft, left.Value)
		return
	}

	if right.Value != expectedRight {
		test.Errorf("Input: %s\nExpected right %t, got %t", input, expectedRight, right.Value)
		return
	}
}
