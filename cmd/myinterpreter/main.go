package main

import (
	"fmt"
	"math"
	"os"
)

var hadError = false
var hadRuntimeError = false

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "tokenize":
		tokens := runTokenize(os.Args[2])
		print(tokens)
		if hadError {
			os.Exit(65)
		}
	case "parse":
		tokens := runTokenize(os.Args[2])
		expr := runParseToExpr(tokens)
		if hadError {
			os.Exit(65)
		}
		fmt.Println(PrintAst(expr))
	case "evaluate":
		tokens := runTokenize(os.Args[2])
		expr := runParseToExpr(tokens)
		InterpretExpr(expr)
		if hadRuntimeError {
			os.Exit(70)
		}
	case "run":
		tokens := runTokenize(os.Args[2])
		statements := runParseToStatements(tokens)
		if hadError {
			os.Exit(65)
		}
		InterpretStatements(statements)
		if hadRuntimeError {
			os.Exit(70)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runTokenize(filename string) []Token {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	scanner := &Scanner{
		Source:  string(fileContents),
		Tokens:  []Token{},
		Start:   0,
		Current: 0,
		Line:    1,
	}
	tokens := scanner.scanTokens()
	return tokens
}

func runParseToStatements(tokens []Token) []Stmt {
	parser := &Parser{tokens, 0}
	return parser.ParseToStatements()
}

func runParseToExpr(tokens []Token) Expr {
	parser := &Parser{tokens, 0}
	return parser.ParseToExpr()
}

func lineError(line int, message string) {
	report(line, "", message)
}

func tokenError(token Token, message string) {
	if token.Type == EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func runtimeError(err RuntimeError) {
	// fmt.Println(err)
	hadRuntimeError = true
}

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hadError = true
}

func print(tokens []Token) {
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func stringifyNumber(number float64, trailingZero bool) string {
	if trailingZero && math.Floor(number) == number {
		return fmt.Sprintf("%v.0", number)
	}
	return fmt.Sprintf("%v", number)
}

// Exposes additional in case we need to display things differently.
func stringify(literal any, nilName string, trailingZero bool) string {
	switch l := literal.(type) {
	case float64:
		return stringifyNumber(l, trailingZero)
	case nil:
		return nilName
	default:
		return fmt.Sprintf("%v", l)
	}
}
