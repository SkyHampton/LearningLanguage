package main

import (
	"fmt"
	"io"
	"learningLanguage/evaluation"
	"learningLanguage/lexer"
	"learningLanguage/parser"
	"learningLanguage/repl"
	"os"
)

func main() {
	// if there are no args, start console REPL
	if len(os.Args) == 1 {
		repl.StartREPL()
	} else {
		// default to files as stdin/out
		inFile := os.Stdin
		outFile := os.Stdout
		// iterate over args
		index := 1
		for index < len(os.Args) {
			// look for args starting with -i or -o
			switch os.Args[index] {
			// input file
			case "-i":
				index++
				inFile, _ = os.Open(os.Args[index])
			// output file
			case "-o":
				index++
				outFile, _ = os.OpenFile(os.Args[index], os.O_WRONLY|os.O_CREATE, 0600)
			default:
				index++
			}
		}
		// execute program written in inFile, output result to outFile
		executeProgram(inFile, outFile)
		inFile.Close()
		outFile.Close()
	}
}

func executeProgram(in io.Reader, out io.Writer) {
	// get program as string
	text, err := io.ReadAll(in)
	if err != nil {
		panic(err)
	}

	// create lexer and parser
	lexer := lexer.New(string(text))
	parser := parser.New(lexer)
	// parse program
	program := parser.ParseProgram()
	// evaluate AST of program, get output and errors
	output, errors := evaluation.EvaluateProgram(program)
	parserErrors := parser.Errors()

	// if there are no errors, write program output
	if len(errors) == 0 && len(parserErrors) == 0 {
		fmt.Fprint(out, output)
	} else {
		if len(parserErrors) > 0 {
			for _, parseError := range parserErrors {
				fmt.Fprint(out, parseError+"\n")
			}
		}
		// otherwise, write errors
		for _, err := range errors {
			fmt.Fprint(out, err+"\n")
		}
	}
}
