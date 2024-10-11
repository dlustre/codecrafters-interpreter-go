package main

import "fmt"

type RuntimeError struct {
	Token   Token
	Message string
}

func (err RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", err.Message, err.Token.Line)
}
