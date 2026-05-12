package evaluation

import (
	"learningLanguage/lexer"
	"learningLanguage/parser"
	"strings"
	"testing"
)

func TestCreateSetEval(test *testing.T) {
	input := `create int x;
			set x = 64;
			print(x);`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	e := Evaluator{}

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)
	expectedOutput := "64"

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s", expectedOutput, output)
	}
}

func TestIfEval(test *testing.T) {
	input := `if (1 > 2) begin; print("1 greater than 2"); end; else begin; print("1 not greater than 2"); end;`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	e := Evaluator{}

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)
	expectedOutput := "1 not greater than 2"

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s", expectedOutput, output)
	}
}

func TestWhileEval(test *testing.T) {
	input := `create int x; set x = 0; while (x < 10) begin; print(x); set x = x + 1; end;`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	e := Evaluator{}

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	parserErrors := p.Errors()
	if len(parserErrors) != 0 {
		test.Fatalf("Errors were found: %v", parserErrors)
	}
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)
	expectedOutput := "0123456789"

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s", expectedOutput, output)
	}
}

func TestCountEval(test *testing.T) {
	input := `count i from 1 to 10 begin; print(i); end;`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	e := Evaluator{}

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	parserErrors := p.Errors()
	if len(parserErrors) != 0 {
		test.Fatalf("Errors were found: %v", parserErrors)
	}
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)
	expectedOutput := "12345678910"

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s", expectedOutput, output)
	}
}

func TestStructEval(test *testing.T) {
	input := `struct myStruct (int x, bool y);
				set myStruct.x = 123;
				set myStruct.y = false;
				print(myStruct.x);
				print(myStruct.y);`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	e := Evaluator{}

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)
	expectedOutput := "123false"

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s", expectedOutput, output)
	}
}

func TestPrefixEval(test *testing.T) {
	input := `print(-123); print(!true);`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	e := Evaluator{}

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)
	expectedOutput := "-123false"

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s", expectedOutput, output)
	}
}

func TestInfixEval(test *testing.T) {
	input := `print(1+1);
			print(2-2);
			print(10/5);
			print(8*8);
			print(1>0);
			print(1>=1);
			print(1==1);
			print(1!=1);
			print(1<2);
			print(1<=1);
			print(true and false);
			print(false or true);`
	l := lexer.New(input)
	p := parser.New(l)
	e := Evaluator{}
	program := p.ParseProgram()

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)
	expectedOutput := "20264truetruetruefalsetruetruefalsetrue"

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s", expectedOutput, output)
	}
}

func TestDataTypes(test *testing.T) {
	input := `create int w;
				set w = 123;
				print(w);
				create bool x;
				set x = true;
				print(x);
				create float y;
				set y = 3.14;
				print(y);
				create string z;
				set z = "Hello World";
				print(z);`
	l := lexer.New(input)
	p := parser.New(l)
	e := Evaluator{}
	program := p.ParseProgram()

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)
	expectedOutput := "123true3.14Hello World"

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s", expectedOutput, output)
	}
}

func TestListLiteral(test *testing.T) {
	input := `create int list a;
				set a = [1, 2, 3, 4, 5];
				print(a);`

	l := lexer.New(input)
	p := parser.New(l)
	e := Evaluator{}
	program := p.ParseProgram()
	parserErrors := p.Errors()
	if len(parserErrors) != 0 {
		test.Fatalf("Parser Errors were fount: %v", parserErrors)
	}

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	expectedOutput := "[1, 2, 3, 4, 5]"

	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s.", expectedOutput, output)
	}
}

func TestListSetIndex(test *testing.T) {
	input := `create int list a;
				set a = [1, 2, 3, 4, 5];
				set a[3] = 30;
				print(a);
				print(a[3]);`

	l := lexer.New(input)
	p := parser.New(l)
	e := Evaluator{}
	program := p.ParseProgram()
	parserErrors := p.Errors()
	if len(parserErrors) != 0 {
		test.Fatalf("Parser Errors were fount: %v", parserErrors)
	}

	output := e.EvaluateProgram(program)
	errors := e.Errors()
	expectedOutput := "[1, 2, 30, 4, 5]30"

	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	if output != expectedOutput {
		test.Fatalf("Incorrect variable value, expected %s, got %s.", expectedOutput, output)
	}
}
