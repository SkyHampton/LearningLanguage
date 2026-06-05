package main

import (
	"encoding/json"
	"fmt"
	"learningLanguage/evaluation"
	"learningLanguage/lexer"
	"learningLanguage/parser"
	"net/http"
)

type CompileResult struct {
	RuntimeErrors []string `json:"runtime_errors"`
	ParserErrors  []string `json:"parser_errors"`
	Output        string   `json:"output"`
}

func runCode(resWriter http.ResponseWriter, req *http.Request) {
	header := resWriter.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Access-Control-Allow-Origin", "http://localhost")

	code := req.PathValue("code")

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()
	parseErrors := p.Errors()
	eval := evaluation.New()
	output := eval.EvaluateProgram(program)
	evalErrors := eval.Errors()

	result := CompileResult{ParserErrors: parseErrors, RuntimeErrors: evalErrors, Output: output}

	resWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(resWriter).Encode(result)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ll/{code}", runCode)

	fmt.Println("Server running at http://localhost:8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Failed to start web server: %s\n", err.Error())
	}
}
