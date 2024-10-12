package main

type Stmt interface {
	Accept(visitor StmtVisitor) any
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt Expression) any
	VisitPrintStmt(stmt Print) any
}

type Expression struct {
	Expression Expr
}

func (t Expression) Accept(visitor StmtVisitor) any {
	return visitor.VisitExpressionStmt(t)
}

type Print struct {
	Expression Expr
}

func (t Print) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrintStmt(t)
}

