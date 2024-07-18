package main

import (
	"fmt"
	"os"
)

const (
	LEFT_PAREN = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
)

var tokenNames = map[int]string{
	LEFT_PAREN:  "LEFT_PAREN",
	RIGHT_PAREN: "RIGHT_PAREN",
	LEFT_BRACE:  "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",
}

func printToken(t int, s string) {
	fmt.Printf("%s %s null\n", tokenNames[t], s)
}

func scanToken(b byte) {
	switch b {
	case '(':
		printToken(LEFT_PAREN, string(b))
	case ')':
		printToken(RIGHT_PAREN, string(b))
	case '{':
		printToken(LEFT_BRACE, string(b))
	case '}':
		printToken(RIGHT_BRACE, string(b))
	}
}

func scanTokens(fileContents []byte) {
	for _, b := range fileContents {
		scanToken(b)
	}
	fmt.Println("EOF  null")
}

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

	if len(fileContents) > 0 {
		scanTokens(fileContents)
	} else {
		fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	}
}
