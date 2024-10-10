package main

type Expr interface {
	Accept(visitor Visitor) any
}

type Visitor interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
}

type Binary struct {
	Left Expr
	Operator Token
	Right Expr
}

func (t Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(t)
}

type Grouping struct {
	Expression Expr
}

func (t Grouping) Accept(visitor Visitor) any {
	return visitor.VisitGroupingExpr(t)
}

type Literal struct {
	Value any
}

func (t Literal) Accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(t)
}

type Unary struct {
	Operator Token
	Right Expr
}

func (t Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(t)
}

