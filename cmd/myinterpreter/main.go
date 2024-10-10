package main

import (
	"fmt"
	"math"
	"os"
)

var hadError = false

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "tokenize":
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		scanner := NewScanner(string(fileContents))
		tokens := scanner.scanTokens()
		print(tokens)
		if hadError {
			return
		}
	case "parse":
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		scanner := NewScanner(string(fileContents))
		tokens := scanner.scanTokens()
		parser := Parser{tokens, 0}
		expression := parser.Parse()

		if hadError {
			return
		}

		fmt.Println(AstPrinter{}.Print(expression))
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

}

// func runFile(path []byte) {}

// func runPrompt() {}

// func run(source string) {}

func lineError(line int, message string) {
	report(line, "", message)
}

func TokenError(token Token, message string) {
	if token.Type == EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
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

// Utility function to display Lox numbers in a way that codecrafters expects.
func formatNumberLiteral(number float64) string {
	// Display numbers with at least one decimal point.
	if math.Floor(number) == number {
		return fmt.Sprintf("%v.0", number)
	}
	return fmt.Sprintf("%v", number)
}

// Utility function to comply with codecrafters' assertions.
func FormatLiteral(literal any, nilName string) string {
	switch l := literal.(type) {
	case float64:
		return formatNumberLiteral(l)
	case nil:
		return nilName
	default:
		return fmt.Sprintf("%v", l)
	}
}
