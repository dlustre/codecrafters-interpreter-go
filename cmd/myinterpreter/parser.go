package main

import "fmt"

type Parser struct {
	Tokens  []Token
	Current int
}

var ErrParse = fmt.Errorf("ParseError")

func (p *Parser) ParseToStatements() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		if statement := p.declaration(); statement != nil {
			statements = append(statements, statement)
		}
	}
	return statements
}

func (p *Parser) ParseToExpr() Expr {
	expr, err := p.expression()
	if err != nil {
		return nil
	}
	return expr
}

// expression -> assignment
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		stmt, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil
		}
		return stmt
	}
	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		return nil
	}
	return stmt
}

// statement -> exprStmt | ifStmt | printStmt | whileStmt | block
func (p *Parser) statement() (Stmt, error) {
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return Block{p.block()}, nil
	}
	return p.expressionStatement()
}

// ifStmt -> "if" "(" expression ")" statement ( "else" statement )?
func (p *Parser) ifStatement() (Stmt, error) {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return If{condition, thenBranch, elseBranch}, nil
}

// printStmt -> "print" expression ";"
func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(SEMICOLON, "Expect ';' after value.")
	return Print{value}, nil
}

// varDecl -> "var" IDENTIFIER ( "=" expression )? ";"
func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		value, err := p.expression()
		if err != nil {
			return nil, err
		}
		initializer = value
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return Var{name, initializer}, nil
}

// whileStmt -> "while" "(" expression ")" statement
func (p *Parser) whileStatement() (Stmt, error) {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(RIGHT_PAREN, "Expect ')' after while condition.")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return While{condition, body}, nil
}

// exprStmt -> expression ";"
func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return Expression{expr}, nil
}

func (p *Parser) block() []Stmt {
	statements := []Stmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		if statement := p.declaration(); statement != nil {
			statements = append(statements, statement)
		}
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

// assignment -> (IDENTIFIER "=" assignment) | logic_or
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(Variable); ok {
			return Assign{varExpr.Name, value}, nil
		}

		parseError(equals, "Invalid assignment target.")
	}

	return expr, nil
}

// logic_or -> logic_and ( "or" logic_and )*
func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = Logical{expr, operator, right}
	}

	return expr, nil
}

// logic_and -> equality ( "and" equality )*
func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = Logical{expr, operator, right}
	}

	return expr, nil
}

// equality -> comparison ( ("!=" | "==") comparison )*
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

// comparison -> term ( (">" | ">=" | "<" | "<=") term )*
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

// term -> factor ( ("-" | "+") term )*
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

// factor -> unary ( ("/" | "*") unary )*
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

// unary -> ( ("!" | "-") unary ) | primary
func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary{operator, right}, nil
	}

	return p.primary()
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return Literal{false}, nil
	}
	if p.match(TRUE) {
		return Literal{true}, nil
	}
	if p.match(NIL) {
		return Literal{nil}, nil
	}
	if p.match(NUMBER, STRING) {
		return Literal{p.previous().Literal}, nil
	}
	if p.match(IDENTIFIER) {
		return Variable{p.previous()}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return Grouping{expr}, nil
	}

	return nil, parseError(p.peek(), "Expect expression.")
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, t := range tokenTypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return p.peek(), parseError(p.peek(), message)
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Current-1]
}

func parseError(token Token, message string) error {
	tokenError(token, message)
	return ErrParse
}

// Gets parser back in sync by discarding tokens until statement boundary is found.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}
