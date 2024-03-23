// expression → equality ;
// equality → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term → factor ( ( "-" | "+" ) factor )* ;
// factor → ternary ( ( "/" | "*" ) ternary )* ;
// ternary → unary ( ? ternary : ternary )
// unary → ( "!" | "-" ) unary | primary ;
// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;

// program → statement* EOF ;
// statement → exprStmt | printStmt ;
// exprStmt → expression ";" ;
// printStmt → "print" expression ";" ;

package main

import (
	u "github.com/core/utils"
)

func NewParserError(message string) error {
	return u.NewError("SyntaxError: %s", message)
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) Parser {
	return Parser{tokens, 0}
}

func (p *Parser) expression() (error, Expression) {
	return p.equality()
}

func (p *Parser) equality() (error, Expression) {
	err, expression := p.comparison()

	if err != nil {
		return err, nil
	}
	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		err, right := p.comparison()
		if err != nil {
			return err, nil
		}
		expression = NewBinary(expression, operator, right)
	}
	return nil, expression
}

func (p *Parser) comparison() (error, Expression) {
	err, expression := p.term()
	if err != nil {
		return err, nil
	}
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		err, right := p.term()
		if err != nil {
			return err, nil
		}
		expression = NewBinary(expression, operator, right)
	}
	return nil, expression
}

func (p *Parser) term() (error, Expression) {
	err, expression := p.factor()
	if err != nil {
		return err, nil
	}
	for p.match(Minus, Plus) {
		operator := p.previous()
		err, right := p.factor()
		if err != nil {
			return err, nil
		}
		expression = NewBinary(expression, operator, right)
	}
	return nil, expression
}

func (p *Parser) factor() (error, Expression) {
	err, expression := p.ternary()
	if err != nil {
		return err, nil
	}
	for p.match(Star, Slash) {
		operator := p.previous()
		err, right := p.ternary()
		if err != nil {
			return err, nil
		}
		expression = NewBinary(expression, operator, right)
	}
	return nil, expression
}

func (p *Parser) ternary() (error, Expression) {
	err, left := p.unary()
	if err != nil {
		return err, nil
	}
	if p.match(Question) {
		err, middle := p.ternary()
		if err != nil {
			return err, nil
		}
		if p.match(Colon) {
			err, right := p.ternary()
			if err != nil {
				return err, nil
			}
			return nil, NewTernary(left, middle, right)
		} else {
			return nil, NewParserError("Expected ':'")
		}

	}
	return nil, left
}

func (p *Parser) unary() (error, Expression) {
	for p.match(Minus, Plus) {
		operator := p.previous()
		err, right := p.unary()
		if err != nil {
			return err, nil
		}
		return nil, NewUnary(operator, right)
	}
	err, primary := p.primary()
	if err != nil {
		return err, nil
	}
	return nil, primary
}

func (p *Parser) primary() (error, Expression) {
	if p.match(False) {
		return nil, NewLiteral(false)
	}
	if p.match(True) {
		return nil, NewLiteral(true)
	}
	if p.match(Null) {
		return nil, NewLiteral(nil)
	}
	if p.match(Number, String) {
		return nil, NewLiteral(p.previous().literal)
	}
	if p.match(LeftBrace) {
		expressionError, expression := p.expression()
		if expressionError != nil {
			return expressionError, nil
		}
		consumeError, _ := p.consume(Semicolon, "Expect ')' after expression.")
		if consumeError != nil {
			return consumeError, nil
		}
		return nil, NewGrouping(expression)
	}
	return NewParserError("Unpredictable expression"), nil
}

func (p Parser) error(token Token, message string) {
	if token.tokenType == Eof {
		u.Report(token.line, " at end", message)
	} else {
		u.Report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func (p *Parser) consume(tokenType string, message string) (error, Token) {
	if p.check(tokenType) {
		return nil, p.advance()
	}
	token := p.peek()
	return NewParserError(message), token
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

func (p *Parser) expressionStatement() (error, Statement) {
	expressionError, expression := p.expression()
	if expressionError != nil {
		return expressionError, nil
	}
	consumeError, _ := p.consume(Semicolon, "Expect ';' after value")
	if consumeError != nil {
		return consumeError, nil
	}
	return nil, NewExpressionStatement(expression)
}

func (p *Parser) printStatement() (error, Statement) {
	expressionError, expression := p.expression()
	if expressionError != nil {
		return expressionError, nil
	}
	consumeError, _ := p.consume(Semicolon, "Expect ';' after value")
	if consumeError != nil {
		return consumeError, nil
	}
	return nil, NewPrintStatement(expression)
}

func (p *Parser) statement() (error, Statement) {
	if p.match(Print) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) parse() (error, []Statement) {
	statements := []Statement{}
	for !p.isAtEnd() {
		err, statement := p.statement()
		if err != nil {
			return err, statements
		}
		statements = append(statements, statement)
	}
	return nil, statements
}
