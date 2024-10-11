package main

import (
	"errors"
	"fmt"
)

type EvalResult struct {
	Value any
	Err   error
}

type Interpreter struct{}

func InterpretAst(expression Expr) {
	evalResult := evaluate(expression)
	var err RuntimeError
	if errors.As(evalResult.Err, &err) {
		runtimeError(err)
	}
	fmt.Println(stringify(evalResult.Value, "nil", false))
}

func evaluate(expr Expr) EvalResult {
	return expr.Accept(Interpreter{}).(EvalResult)
}

func (Interpreter) VisitBinaryExpr(expr Binary) any {
	leftResult := evaluate(expr.Left)
	if leftResult.Err != nil {
		return leftResult
	}
	rightResult := evaluate(expr.Right)
	if rightResult.Err != nil {
		return rightResult
	}

	left := leftResult.Value
	right := rightResult.Value

	switch expr.Operator.Type {
	case BANG_EQUAL:
		return EvalResult{!isEqual(left, right), nil}
	case EQUAL_EQUAL:
		return EvalResult{isEqual(left, right), nil}
	case GREATER:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return EvalResult{nil, err}
		}
		return EvalResult{left.(float64) > right.(float64), nil}
	case GREATER_EQUAL:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return EvalResult{nil, err}
		}
		return EvalResult{left.(float64) >= right.(float64), nil}
	case LESS:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return EvalResult{nil, err}
		}
		return EvalResult{left.(float64) < right.(float64), nil}
	case LESS_EQUAL:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return EvalResult{nil, err}
		}
		return EvalResult{left.(float64) <= right.(float64), nil}
	case MINUS:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return EvalResult{nil, err}
		}
		return EvalResult{left.(float64) - right.(float64), nil}
	case SLASH:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return EvalResult{nil, err}
		}
		return EvalResult{left.(float64) / right.(float64), nil}
	case STAR:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return EvalResult{nil, err}
		}
		return EvalResult{left.(float64) * right.(float64), nil}
	case PLUS:
		numberLeft, leftIsNumber := left.(float64)
		numberRight, rightIsNumber := right.(float64)

		if leftIsNumber && rightIsNumber {
			return EvalResult{numberLeft + numberRight, nil}
		}

		stringLeft, leftIsString := left.(string)
		stringRight, rightIsString := right.(string)

		if leftIsString && rightIsString {
			return EvalResult{stringLeft + stringRight, nil}
		}

		return EvalResult{nil, &RuntimeError{expr.Operator, "Operands must be two numbers or two strings."}}
	}

	// Unreachable.
	return EvalResult{}
}

func (Interpreter) VisitGroupingExpr(expr Grouping) any {
	return evaluate(expr.Expression)
}

func (Interpreter) VisitLiteralExpr(expr Literal) any {
	return EvalResult{expr.Value, nil}
}

func (Interpreter) VisitUnaryExpr(expr Unary) any {
	rightResult := evaluate(expr.Right)
	if rightResult.Err != nil {
		return rightResult
	}

	right := rightResult.Value

	switch expr.Operator.Type {
	case MINUS:
		err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return EvalResult{nil, err}
		}
		return EvalResult{-right.(float64), nil}
	case BANG:
		return EvalResult{!isTruthy(right), nil}
	}

	// Unreachable.
	return EvalResult{}
}

func checkNumberOperand(operator Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return &RuntimeError{operator, "Operand must be a number."}
}

func checkNumberOperands(operator Token, left, right any) error {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)

	if leftOk && rightOk {
		return nil
	}

	return &RuntimeError{operator, "Operands must be numbers."}
}

func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if boolObject, ok := object.(bool); ok {
		return boolObject
	}
	return true
}

/*
Not sure about this one.
I cargoculted the java logic but golang nil has different semantics than Java null.
And they use a.equals(b) instead of a == b.
*/
func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}
