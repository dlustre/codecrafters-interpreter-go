package main

import (
	"strings"
)

type AstPrinter struct{}

func (ap AstPrinter) Print(expr Expr) string {
	return expr.Accept(ap).(string)
}

func (ap AstPrinter) VisitBinaryExpr(expr Binary) any {
	return parenthesize(ap, expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap AstPrinter) VisitGroupingExpr(expr Grouping) any {
	return parenthesize(ap, "group", expr.Expression)
}

func (ap AstPrinter) VisitLiteralExpr(expr Literal) any {
	return FormatLiteral(expr.Value)
}

func (ap AstPrinter) VisitUnaryExpr(expr Unary) any {
	return parenthesize(ap, expr.Operator.Lexeme, expr.Right)

}

func parenthesize(ap AstPrinter, name string, exprs ...Expr) string {
	var sb strings.Builder

	sb.WriteRune('(')
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteRune(' ')
		sb.WriteString(expr.Accept(ap).(string))
	}
	sb.WriteRune(')')

	return sb.String()
}

// func main() {
// 	expression := Binary{
// 		Unary{
// 			Token{
// 				MINUS, "-", nil, 1,
// 			},
// 			Literal{123},
// 		},
// 		Token{STAR, "*", nil, 1},
// 		Grouping{Literal{45.67}},
// 	}

// 	printResult := AstPrinter{}.print(expression)
// 	expectedResult := "(* (- 123) (group 45.67))"
// 	fmt.Println(printResult)
// 	if printResult != expectedResult {
// 		panic("Expected '" + expectedResult + "' as the printed result.")
// 	}
// 	fmt.Println("✔ AstPrinter is functional.")
// }
