package main

import (
	"fmt"
	"os"
)

var hadError = false

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

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
		os.Exit(65)
	}
}

// func runFile(path []byte) {}

// func runPrompt() {}

// func run(source string) {}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s", line, where, message)
	hadError = true
}

func print(tokens []Token) {
	for _, token := range tokens {
		fmt.Println(token)
	}
}
