package repl

import (
	"bufio"
	"fmt"
	"learningLanguage/evaluation"
	"learningLanguage/lexer"
	"learningLanguage/parser"
	"learningLanguage/token"
	"os"
)

// command line REPL program
func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	// LOOP
	for {
		fmt.Printf(">> ")
		// get input from stdin
		scanned := scanner.Scan()
		// if no input, exit
		if !scanned {
			return
		}

		// get line of input, create lexer, parser, and parse the program from the line of input
		line := scanner.Text()
		lexer := lexer.New(line)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		// evaluate code and get errors/output
		output, errors := evaluation.EvaluateProgram(program)

		// if parser errors exist, display them
		if len(parser.Errors()) > 0 {
			for _, error := range parser.Errors() {
				fmt.Printf("ERROR: %s\n", error)
			}
			// otherwise, if compilation errors exist, display them
		} else if len(errors) > 0 {
			for _, error := range errors {
				fmt.Printf("ERROR: %s\n", error)
			}
			// otherwise, display output of code
		} else {
			fmt.Printf("%s\n", output)
		}
	}
}

// testing function which reads, lexes, prints, and loops
func StartRLPL() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// get user input and tokenize it
		line := scanner.Text()
		lexer := lexer.New(line)
		tok := lexer.NextToken()

		// display token info
		for tok.Type != token.EOF {
			fmt.Printf("%+v\n", tok)
			tok = lexer.NextToken()
		}
	}
}

// testing function which reads, parses, prints, and loops
func StartRPPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// get user input and parse it as a program
		line := scanner.Text()
		lexer := lexer.New(line)
		parser := parser.New(lexer)
		program := parser.ParseProgram()

		// display parser errors if they exist
		if len(parser.Errors()) > 0 {
			for _, error := range parser.Errors() {
				fmt.Printf("ERROR: %s\n", error)
			}
			// otherwise, display string representation of AST node
		} else {
			fmt.Printf("%s\n", program.String())
		}
	}
}
