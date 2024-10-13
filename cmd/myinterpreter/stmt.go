package main

type Stmt interface {
	Accept(visitor StmtVisitor) any
}

type StmtVisitor interface {
	VisitBlockStmt(stmt Block) any
	VisitExpressionStmt(stmt Expression) any
	VisitPrintStmt(stmt Print) any
	VisitVarStmt(stmt Var) any
}

type Block struct {
	Statements []Stmt
}

func (t Block) Accept(visitor StmtVisitor) any {
	return visitor.VisitBlockStmt(t)
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

type Var struct {
	Name Token
	Initializer Expr
}

func (t Var) Accept(visitor StmtVisitor) any {
	return visitor.VisitVarStmt(t)
}

