package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generateast.go <output directory>")
		os.Exit(1)
	}
	outputDir := os.Args[1]
	defineAst(outputDir, "Expr", []string{
		"Assign   : Name Token, Value Expr",
		"Binary   : Left Expr, Operator Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value any",
		"Logical  : Left Expr, Operator Token, Right Expr",
		"Unary    : Operator Token, Right Expr",
		"Variable : Name Token",
	})
	defineAst(outputDir, "Stmt", []string{
		"Block      : Statements []Stmt",
		"Expression : Expression Expr",
		"If         : Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
		"Print      : Expression Expr",
		"Var        : Name Token, Initializer Expr",
		"While      : Condition Expr, Body Stmt",
	})
}

func defineAst(outputDir, baseName string, types []string) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer closeFile(file)

	file.WriteString("package main\n\n")
	file.WriteString("type " + baseName + " interface {\n")
	file.WriteString("\tAccept(visitor " + baseName + "Visitor) any\n")
	file.WriteString("}\n\n")

	defineVisitor(file, baseName, types)

	for _, t := range types {
		className := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(file, baseName, className, fields)
	}
}

func defineVisitor(file *os.File, baseName string, types []string) {
	file.WriteString("type " + baseName + "Visitor interface {\n")
	for _, t := range types {
		typeName := strings.TrimSpace(strings.Split(t, ":")[0])
		file.WriteString("\tVisit" + typeName + baseName + "(" + strings.ToLower(baseName) + " " + typeName + ") any\n")
	}
	file.WriteString("}\n\n")
}

func defineType(file *os.File, baseName, className, fields string) {
	file.WriteString("type " + className + " struct {\n")
	for _, field := range strings.Split(fields, ", ") {
		file.WriteString("\t" + field + "\n")
	}
	file.WriteString("}\n\n")

	file.WriteString("func (t " + className + ") Accept(visitor " + baseName + "Visitor) any {\n")
	file.WriteString("\treturn visitor.Visit" + className + baseName + "(t)\n")
	file.WriteString("}\n\n")
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
}
