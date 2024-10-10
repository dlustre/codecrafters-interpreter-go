package main

import "fmt"

type Parser struct {
	Tokens  []Token
	Current int
}

func (p *Parser) Parse() Expr {
	expr, err := p.Expression()
	if err != nil {
		return nil
	}
	return expr
}

// Expression -> equality
func (p *Parser) Expression() (Expr, error) {
	return p.equality()
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
	if p.match(LEFT_PAREN) {
		expr, err := p.Expression()
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

var ErrParse = fmt.Errorf("ParseError")

func parseError(token Token, message string) error {
	TokenError(token, message)
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
