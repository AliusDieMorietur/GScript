// expression → equality ;
// equality → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term → factor ( ( "-" | "+" ) factor )* ;
// factor → ternary ( ( "/" | "*" ) ternary )* ;
// ternary → unary ( ? ternary : ternary )
// unary → ( "!" | "-" ) unary | primary ;
// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
package main

import (
	u "github.com/core/utils"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) Parser {
	return Parser{tokens, 0}
}

func (p *Parser) expression() Expression {
	return p.equality()
}

func (p *Parser) equality() Expression {
	expression := p.comparison()
	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expression = NewBinary(expression, operator, right)
	}
	return expression
}

func (p *Parser) comparison() Expression {
	expression := p.term()
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right := p.term()
		expression = NewBinary(expression, operator, right)
	}
	return expression
}

func (p *Parser) term() Expression {
	expression := p.factor()
	for p.match(Minus, Plus) {
		operator := p.previous()
		right := p.factor()
		expression = NewBinary(expression, operator, right)
	}
	return expression
}

func (p *Parser) factor() Expression {
	expression := p.ternary()
	for p.match(Star, Slash) {
		operator := p.previous()
		right := p.ternary()
		expression = NewBinary(expression, operator, right)
	}
	return expression
}

func (p *Parser) ternary() Expression {
	left := p.unary()
	if p.match(Question) {
		middle := p.ternary()
		if p.match(Colon) {
			right := p.ternary()
			return NewTernary(left, middle, right)
		} else {
			panic("Expected ':'")
		}

	}

	return left
}

func (p *Parser) unary() Expression {

	for p.match(Minus, Plus) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}
	return p.primary()
}

func (p *Parser) primary() Expression {
	if p.match(False) {
		return NewLiteral(false)
	}
	if p.match(True) {
		return NewLiteral(true)
	}
	if p.match(Null) {
		return NewLiteral(nil)
	}
	if p.match(Number, String) {
		return NewLiteral(p.previous().literal)
	}
	if p.match(LeftBrace) {
		expression := p.expression()
		p.consume(RightBrace, "Expect ')' after expression.")
		return NewGrouping(expression)
	}
	expression := p.expression()
	return NewGrouping(expression)
	// panic("Exprected expression")
}

func (p Parser) error(token Token, message string) {
	if token.tokenType == Eof {
		u.Report(token.line, " at end", message)
	} else {
		u.Report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func (p *Parser) consume(tokenType string, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.error(p.peek(), message)
	panic("Unpredictable expression in consume")
}

func (p *Parser) match(tokenTypes ...string) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType string) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p Parser) isAtEnd() bool {
	return p.peek().tokenType == Eof
}

func (p Parser) peek() Token {
	return p.tokens[p.current]
}

func (p Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == Semicolon {
			return
		}

		switch p.peek().tokenType {
		case Struct:
		case Fn:
		case Let:
		case For:
		case If:
		case While:
		case Print:
			return
		}

		p.advance()
	}
}

// func (p *Parser) statement() Statement {

// }

func (p *Parser) parse() Expression {
	// statements := []Statement{}
	// for !p.isAtEnd() {
	// 	statements = append(statements, p.statement())
	// }
	return p.expression()
}
