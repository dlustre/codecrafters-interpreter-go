package main

import (
	"fmt"
	"os"
)

func scanTokens(fileContents []byte) {
	for _, b := range fileContents {
		switch b {
		case '(':
			fmt.Printf("LEFT_PAREN %s null", string(b))
		case ')':
			fmt.Printf("RIGHT_PAREN %s null", string(b))
		}
		fmt.Print("\n")
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
