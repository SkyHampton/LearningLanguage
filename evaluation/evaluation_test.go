package evaluation

import (
	"learningLanguage/lexer"
	"learningLanguage/parser"
	"testing"
)

func TestCorrectCode(Test *testing.T) {
	tests := []struct {
		input          string
		expectedErrors bool
		expectedOutput string
	}{
		{"create int x; set x = 3; print(x);", false, "3"},
		{"create bool x; set x = true; print(x);", false, "true"},
		{"create float x; set x = 3.14; print(x);", false, "3.14"},
		{"create string x; set x = \"Hello!\"; print(x);", false, "Hello!"},
		{"if (1 > 2) begin; print(\"1 greater than 2\"); end; else begin; print(\"1 not greater than 2\"); end;", false, "1 not greater than 2"},
		{"if (2 > 1) begin; print(\"2 greater than 1\"); end; else begin; print(\"2 not greater than 1\"); end;", false, "2 greater than 1"},
		{"create int x; set x = 0; while (x < 10) begin; print(x); set x = x + 1; end;", false, "0123456789"},
		{"count i from 1 to 10 begin; print(i); end;", false, "12345678910"},
		{"count i from 1 to 10 by 2 begin; print(i); end;", false, "13579"},
		{"struct myStruct (int x, bool y);set myStruct.x = 123;set myStruct.y = false;print(myStruct.x);print(myStruct.y);", false, "123false"},
		{"struct myStruct (int x, bool y)[x:123,y:false];print(myStruct.x);print(myStruct.y);", false, "123false"},
		{"print(-123 + 3); print(!true);", false, "-120false"},
		{"print(1+1);print(2-2);print(10/5);print(8*8);", false, "20264"},
		{"print(3.14+1);print(3.14-2);print(3/5);print(0.2*4);", false, "4.141.140.60.8"},
		{"print(1>0);print(1>=1);print(1==1);print(1!=1);print(1<2);print(1<=1);", false, "truetruetruefalsetruetrue"},
		{"print(true and false);print(false or true);", false, "falsetrue"},
		{"create int list a;set a = [1, 2, 3, 4, 5];print(a);", false, "[1, 2, 3, 4, 5]"},
		{"create int list a;set a = [1, 2, 3, 4, 5];print(a);set a[1] = 100; print(a);", false, "[1, 2, 3, 4, 5][100, 2, 3, 4, 5]"},
		{"create int list a;set a = [1, 2, 3, 4, 5];print(a); append 6 to a; print(a); print(len(a));", false, "[1, 2, 3, 4, 5][1, 2, 3, 4, 5, 6]6"},
	}

	for _, test := range tests {
		testEvaluation(Test, test.input, test.expectedErrors, test.expectedOutput)
	}
}

func TestIncorrectCode(Test *testing.T) {
	tests := []struct {
		input          string
		expectedErrors bool
		expectedOutput string
	}{
		{"create int x; set x = 3.14; print(x);", true, "Errors were found, no output."},
		{"create bool x; set x = 1; print(x);", true, "Errors were found, no output."},
		{"set x = 3.14; print(x);", true, "Errors were found, no output."},
		{"create string x; set x = false; print(x);", true, "Errors were found, no output."},
		{"if (1) begin; print(\"1 greater than 2\"); end; else begin; print(\"1 not greater than 2\"); end;", true, "Errors were found, no output."},
		{"if (3.14) begin; print(\"2 greater than 1\"); end; else begin; print(\"2 not greater than 1\"); end;", true, "Errors were found, no output."},
		{"create int x; set x = 0; while (x) begin; print(x); set x = x + 1; end;", true, "Errors were found, no output."},
		{"count i from 1 to 3.14 begin; print(i); end;", true, "Errors were found, no output."},
		{"count i from true to 10 by 2 begin; print(i); end;", true, "Errors were found, no output."},
		{"struct myStruct (int x, bool y);set myStruct.x = 3.14;set myStruct.y = 1;print(myStruct.x);print(myStruct.y);", true, "Errors were found, no output."},
		{"set myStruct.x = 3.14;print(myStruct.x);", true, "Errors were found, no output."},
		{"print(-true + 3); print(!123);", true, "Errors were found, no output."},
		{"print(1+true);print(2-2);print(10/5);print(8*8);", true, "Errors were found, no output."},
		{"print(3/0);", true, "Errors were found, no output."},
		{"print(1>true);", true, "Errors were found, no output."},
		{"print(true and 1);", true, "Errors were found, no output."},
		{"create int list a;set a = [1, 2, 3.14, 4, 5];print(a);", true, "Errors were found, no output."},
		{"create int list a;set a = [1, 2, 3, 4, 5];print(a);set a[-1] = 100; print(a);", true, "Errors were found, no output."},
		{"create int list a;set a = [1, 2, 3, 4, 5];print(a);set a[6] = 100; print(a);", true, "Errors were found, no output."},
		{"create int list a;set a = [1, 2, 3, 4, 5];print(a); append 3.14 to a; print(a); print(len(a));", true, "Errors were found, no output."},
	}

	for _, test := range tests {
		testEvaluation(Test, test.input, test.expectedErrors, test.expectedOutput)
	}
}

func testEvaluation(Test *testing.T, input string, expectedErrors bool, expectedOutput string) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	e := New()
	output := e.EvaluateProgram(program)
	errors := e.Errors()

	if (len(errors) > 0) != expectedErrors {
		Test.Errorf("Input: %s\nUnexpected erors: %v", input, errors)
	}

	if output != expectedOutput {
		Test.Errorf("Input: %s\nExpected output %s, got %s.", input, expectedOutput, output)
		return
	}
}
