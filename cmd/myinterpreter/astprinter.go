package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (ap AstPrinter) PrintAst(expr Expr) string {
	return expr.Accept(ap).(string)
}

func (AstPrinter) VisitBinaryExpr(expr Binary) any {
	return parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (AstPrinter) VisitGroupingExpr(expr Grouping) any {
	return parenthesize("group", expr.Expression)
}

func (AstPrinter) VisitLiteralExpr(expr Literal) any {
	return stringify(expr.Value, "nil", true)
}

func (AstPrinter) VisitUnaryExpr(expr Unary) any {
	return parenthesize(expr.Operator.Lexeme, expr.Right)
}

func parenthesize(name string, exprs ...Expr) string {
	var sb strings.Builder

	sb.WriteRune('(')
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteRune(' ')
		sb.WriteString(expr.Accept(AstPrinter{}).(string))
	}
	sb.WriteRune(')')

	return sb.String()
}

func Test() {
	expression := Binary{
		Unary{
			Token{
				MINUS, "-", nil, 1,
			},
			Literal{123},
		},
		Token{STAR, "*", nil, 1},
		Grouping{Literal{45.67}},
	}

	printResult := AstPrinter{}.PrintAst(expression)
	expectedResult := "(* (- 123) (group 45.67))"
	fmt.Println(printResult)
	if printResult != expectedResult {
		panic("Expected '" + expectedResult + "' as the printed result.")
	}
	fmt.Println("✔ AstPrinter is functional.")
}
