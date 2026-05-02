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

	output, errors := EvaluateProgram(program)
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)

	if output != "64" {
		test.Fatalf("Incorrect variable value, expected 64, got %s", output)
	}
}

func TestIfEval(test *testing.T) {
	input := `if (1 > 2) begin; print("1 greater than 2"); end; else begin; print("1 not greater than 2"); end;`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	output, errors := EvaluateProgram(program)
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)

	if output != "1 not greater than 2" {
		test.Fatalf("Incorrect variable value, expected false, got %s", output)
	}
}

func TestWhileEval(test *testing.T) {
	input := `create int x; set x = 0; while (x < 10) begin; print(x); set x = x + 1; end;`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	output, errors := EvaluateProgram(program)
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)

	if output != "0\n1\n2\n3\n4\n5\n6\n7\n8\n9" {
		test.Fatalf("Incorrect variable value, expected 0123456789, got %s", output)
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

	output, errors := EvaluateProgram(program)
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)

	if output != "123\nfalse" {
		test.Fatalf("Incorrect variable value, expected \n123\nfalse, got \n%s", output)
	}
}

func TestPrefixEval(test *testing.T) {
	input := `print(-123); print(!true);`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	output, errors := EvaluateProgram(program)
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)

	if output != "-123\nfalse" {
		test.Fatalf("Incorrect variable value, expected \n-123\nfalse, got \n%s", output)
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
	program := p.ParseProgram()

	output, errors := EvaluateProgram(program)
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)

	if output != "2\n0\n2\n64\ntrue\ntrue\ntrue\nfalse\ntrue\ntrue\nfalse\ntrue" {
		test.Fatalf("Incorrect variable value, expected \n2\n0\n2\n64\ntrue\ntrue\ntrue\nfalse\ntrue\ntrue\nfalse\ntrue\ngot \n%s", output)
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
	program := p.ParseProgram()

	output, errors := EvaluateProgram(program)
	if len(errors) != 0 {
		test.Fatalf("Errors were found: %v", errors)
	}

	output = strings.TrimSpace(output)

	if output != "123\ntrue\n3.14\nHello World" {
		test.Fatalf("Incorrect variable value, expected \n123\ntrue\n3.14\nHello World\ngot \n%s", output)
	}
}
