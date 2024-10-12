package main

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitAssignExpr(expr Assign) any
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
	VisitVariableExpr(expr Variable) any
}

type Assign struct {
	Name Token
	Value Expr
}

func (t Assign) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignExpr(t)
}

type Binary struct {
	Left Expr
	Operator Token
	Right Expr
}

func (t Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(t)
}

type Grouping struct {
	Expression Expr
}

func (t Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(t)
}

type Literal struct {
	Value any
}

func (t Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(t)
}

type Unary struct {
	Operator Token
	Right Expr
}

func (t Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(t)
}

type Variable struct {
	Name Token
}

func (t Variable) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpr(t)
}

