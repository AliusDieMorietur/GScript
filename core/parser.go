// expression → assignment ;
// assignment → ( call "." )? IDENTIFIER "=" assignment | ternary ;
// ternary → logicOr ( ? ternary : ternary ) ;
// logicOr → logicAnd ( || logicAnd )*;
// logicAnd → equality ( && equality  )*;
// equality → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term → factor ( ( "-" | "+" ) factor )* ;
// factor → unary ( ( "/" | "*" ) unary )* ;
// unary → ( "!" | "-" ) unary | function ;
// function → "fn" IDENTIFIER ? "(" parameters? ")" block ;
// parameters → IDENTIFIER ( "," IDENTIFIER )* ;
// call → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
// arguments → expression ( "," expression )* ;
// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER  ;

// program → declaration* EOF ;
// declaration → structDecl | fnDecl | letDecl | statement ;
// statement → exprStmt | forStmt | ifStmt | printStmt | returnStmt | whileStmt | block ;
// forStmt → "for" "(" ( letDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;
// whileStmt → "while" "(" expression ")" statement ;
// block → "{" declaration* "}"
// ifStmt → "if" "(" expression ")" ( "else" "if" "(" expression ")" )* statement ( "else" statement )? ;
// exprStmt → expression ";" ;
// printStmt → "print" expression ";" ;
// returnStmt → "return" expression? ";" ;
// letDecl → "let" IDENTIFIER ( "=" expression )? ";" ;
// structDecl → "struct" IDENTIFIER "{" function* "}" ;
// fnDecl → "fn" function ;
// function → IDENTIFIER "(" parameters? ")" block ;
// parameters → IDENTIFIER ( "," IDENTIFIER )* ;

package main

const MAX_FN_ARGUMENTS_COUNT = 255

type Parser struct {
	tokens       []*Token
	current      int
	currentBlock int
	blockName    int
}

func NewParser(tokens []*Token) Parser {
	return Parser{tokens, 0, 0, 0}
}

func (p *Parser) expression() (error, Expression) {
	return p.assignment()
}

func (p *Parser) assignment() (error, Expression) {
	err, expression := p.ternary()
	if err != nil {
		return err, expression
	}
	if p.match(Equal) {
		err, value := p.assignment()
		if err != nil {
			return err, expression
		}

		if variable, ok := (expression).(Variable); ok {
			return nil, (NewAssignment(variable.name, value))
		}
		if get, ok := (expression).(Get); ok {
			return nil, (NewSet(get.name, get.object, value))
		}

		return NewParserError("Invalid assignment target"), nil
	}
	return nil, expression
}

func (p *Parser) ternary() (error, Expression) {
	err, left := p.or()
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
			return nil, (NewTernary(left, middle, right))
		} else {
			return NewParserError("Expected ':'"), nil
		}
	}
	return nil, left
}

func (p *Parser) or() (error, Expression) {
	err, expression := p.and()
	if err != nil {
		return err, nil
	}
	for p.match(Or) {
		operator := p.previous()
		err, right := p.and()
		if err != nil {
			return err, nil
		}
		expression = (NewLogical(expression, operator, right))
	}
	return nil, expression
}

func (p *Parser) and() (error, Expression) {
	err, expression := p.equality()
	if err != nil {
		return err, nil
	}
	for p.match(And) {
		operator := p.previous()
		err, right := p.equality()
		if err != nil {
			return err, nil
		}
		expression = (NewLogical(expression, operator, right))
	}
	return nil, expression
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
		expression = (NewBinary(expression, operator, right))
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
		expression = (NewBinary(expression, operator, right))
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
		expression = (NewBinary(expression, operator, right))
	}
	return nil, expression
}

func (p *Parser) factor() (error, Expression) {
	err, expression := p.unary()
	if err != nil {
		return err, nil
	}
	for p.match(Star, Slash) {
		operator := p.previous()
		err, right := p.unary()
		if err != nil {
			return err, nil
		}
		expression = (NewBinary(expression, operator, right))
	}
	return nil, expression
}

func (p *Parser) unary() (error, Expression) {
	if p.match(Minus, Plus) {
		operator := p.previous()
		err, right := p.unary()
		if err != nil {
			return err, nil
		}
		return nil, (NewUnary(operator, right))
	}
	err, call := p.function()
	if err != nil {
		return err, nil
	}
	return nil, call
}

func (p *Parser) finishCall(callee Expression) (error, Expression) {
	arguments := []Expression{}
	if !p.check(RightBrace) {
		for {
			if len(arguments) >= MAX_FN_ARGUMENTS_COUNT {
				return NewParserError("Can't have more than %v arguments", MAX_FN_ARGUMENTS_COUNT), nil
			}
			err, expression := p.expression()
			if err != nil {
				return err, nil
			}
			arguments = append(arguments, expression)
			if !p.match(Comma) {
				break
			}
		}
	}
	err, paren := p.consume(RightBrace, "Expect ')' after arguments.")
	if err != nil {
		return err, nil
	}
	return nil, (NewCall(callee, paren, arguments))
}

func (p *Parser) call() (error, Expression) {
	err, expression := p.primary()
	if err != nil {
		return err, nil
	}
	for {
		if p.match(LeftBrace) {
			err, call := p.finishCall(expression)
			if err != nil {
				return err, nil
			}
			expression = call
		} else if p.match(Dot) {
			err, name := p.consume(Identifier, "Exepcted property name after '.'")
			if err != nil {
				return err, nil
			}
			expression = NewGet(name, expression)
		} else {
			break
		}
	}
	return nil, expression
}

func (p *Parser) primary() (error, Expression) {
	if p.match(False) {
		return nil, (NewLiteral(false))
	}
	if p.match(True) {
		return nil, (NewLiteral(true))
	}
	if p.match(Null) {
		return nil, (NewLiteral(nil))
	}
	if p.match(Number, String) {
		return nil, (NewLiteral(p.previous().literal))
	}
	if p.match(Identifier) {
		return nil, (NewVariable(p.previous()))
	}
	if p.match(LeftBrace) {
		expressionError, expression := p.expression()
		if expressionError != nil {
			return expressionError, nil
		}
		consumeError, _ := p.consume(RightBrace, "Expect ')' after expression.")
		if consumeError != nil {
			return consumeError, nil
		}
		return nil, (NewGrouping(expression))
	}
	return NewParserError("Unpredictable expression"), nil
}

func (p *Parser) consume(tokenType string, message string) (error, *Token) {
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

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p Parser) isAtEnd() bool {
	return p.peek().tokenType == Eof
}

func (p Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p Parser) previous() *Token {
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
	if _, ok := (expression).(*Function); !ok {
		consumeError, _ := p.consume(Semicolon, "Expected ';' after value")
		if consumeError != nil {
			return consumeError, nil
		}
	}
	return nil, (NewExpressionStatement(expression))
}

func (p *Parser) printStatement() (error, Statement) {
	expressionError, expression := p.expression()
	if expressionError != nil {
		return expressionError, nil
	}
	consumeError, _ := p.consume(Semicolon, "Expected ';' after value")
	if consumeError != nil {
		return consumeError, nil
	}
	return nil, (NewPrintStatement(expression))
}

func (p *Parser) block() (error, []Statement) {
	statements := []Statement{}
	for !p.check(RightCurlyBracket) && !p.isAtEnd() {
		err, statement := p.declaration()
		if err != nil {
			return err, statements
		}
		statements = append(statements, statement)
	}
	p.consume(RightCurlyBracket, "Expect '}' after block")
	return nil, statements
}

func (p *Parser) ifStatement() (error, Statement) {
	consumeErrLeft, _ := p.consume(LeftBrace, "Expected '('")
	if consumeErrLeft != nil {
		return consumeErrLeft, nil
	}
	conditionErr, condition := p.expression()
	if conditionErr != nil {
		return conditionErr, nil
	}
	consumeErrRight, _ := p.consume(RightBrace, "Expected ')'")
	if consumeErrRight != nil {
		return consumeErrRight, nil
	}
	err, thenBranch := p.statement()
	if err != nil {
		return err, nil
	}
	var elseBranch Statement
	if p.match(Else) {
		err, statement := p.statement()
		if err != nil {
			return err, nil
		}
		elseBranch = statement
	}
	return nil, (NewIfStatement(condition, thenBranch, elseBranch))
}

func (p *Parser) whileStatement() (error, Statement) {
	consumeErrLeft, _ := p.consume(LeftBrace, "Expected '('")
	if consumeErrLeft != nil {
		return consumeErrLeft, nil
	}
	conditionErr, condition := p.expression()
	if conditionErr != nil {
		return conditionErr, nil
	}
	consumeErrRight, _ := p.consume(RightBrace, "Expected ')'")
	if consumeErrRight != nil {
		return consumeErrRight, nil
	}
	err, statement := p.statement()
	if err != nil {
		return err, nil
	}
	return nil, (NewWhileStatement(condition, statement))
}

func (p *Parser) forStatement() (error, Statement) {
	consumeErrLeft, _ := p.consume(LeftBrace, "Expected '('")
	if consumeErrLeft != nil {
		return consumeErrLeft, nil
	}
	if !p.match(Let) {
		return NewParserError("Let expected"), nil
	}
	err, initializer := p.letDeclaration()
	if err != nil {
		return err, nil
	}
	err, condition := p.expression()
	if err != nil {
		return err, nil
	}
	consumeErrSemicolon, _ := p.consume(Semicolon, "Expected ';'")
	if consumeErrSemicolon != nil {
		return consumeErrSemicolon, nil
	}
	err, increment := p.expression()
	if err != nil {
		return err, nil
	}
	consumeErrRight, _ := p.consume(RightBrace, "Expected ')'")
	if consumeErrRight != nil {
		return consumeErrRight, nil
	}
	err, body := p.statement()
	if err != nil {
		return err, nil
	}
	return nil, (NewForStatement(condition, initializer, increment, body))
}

func (p *Parser) returnStatement() (error, Statement) {
	var value Expression
	if !p.check(Semicolon) {
		err, expression := p.expression()
		if err != nil {
			return err, nil
		}
		value = expression
	}
	consumeError, _ := p.consume(Semicolon, "Expected ';' after value")
	if consumeError != nil {
		return consumeError, nil
	}
	return nil, (NewReturnStatement(value))
}

func (p *Parser) statement() (error, Statement) {
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(While) {
		return p.whileStatement()
	}
	if p.match(If) {
		return p.ifStatement()
	}
	if p.match(LeftCurlyBracket) {
		err, statements := p.block()
		if err != nil {
			return err, nil
		}
		return nil, (NewBlockStatement(statements))
	}
	if p.match(Print) {
		return p.printStatement()
	}
	if p.match(Return) {
		return p.returnStatement()
	}
	if p.match(Break) {
		consumeError, _ := p.consume(Semicolon, "Expected ';' after value")
		if consumeError != nil {
			return consumeError, nil
		}
		return nil, (NewBreakStatement())
	}
	if p.match(Continue) {
		consumeError, _ := p.consume(Semicolon, "Expected ';' after value")
		if consumeError != nil {
			return consumeError, nil
		}
		return nil, (NewContinueStatement())
	}
	return p.expressionStatement()
}

func (p *Parser) letDeclaration() (error, Statement) {
	identifierErr, name := p.consume(Identifier, "Variable name expected")
	if identifierErr != nil {
		return identifierErr, nil
	}
	var initializer Expression
	if p.match(Equal) {
		err, expression := p.expression()
		if err != nil {
			return err, nil
		}
		initializer = expression
	}
	err, _ := p.consume(Semicolon, "Expected ';' after variable declaration.")
	if err != nil {
		return err, nil
	}
	return nil, (NewLetStatement(name, initializer))
}

func (p *Parser) function() (error, Expression) {
	if !p.match(Fn) {
		return p.call()
	}
	name := NewToken(AnonymusFunction, AnonymusFunction, "", 0)
	if p.check(Identifier) {
		identifierErr, token := p.consume(Identifier, "Function name expected")
		if identifierErr != nil {
			return identifierErr, nil
		}
		name = token
	}
	leftBraceErr, _ := p.consume(LeftBrace, "Expected '('")
	if leftBraceErr != nil {
		return leftBraceErr, nil
	}
	parameters := []*Token{}
	if !p.check(RightBrace) {
		for {
			if len(parameters) > MAX_FN_ARGUMENTS_COUNT {
				return NewParserError("Can't have more than %v parameters", MAX_FN_ARGUMENTS_COUNT), nil
			}
			consumeErr, identifier := p.consume(Identifier, "Expected parameter name")
			if consumeErr != nil {
				return consumeErr, nil
			}
			parameters = append(parameters, identifier)
			if !p.match(Comma) {
				break
			}
		}
	}
	rightBraceErr, _ := p.consume(RightBrace, "Expected ')' after parameters")
	if rightBraceErr != nil {
		return rightBraceErr, nil
	}
	leftCurlyBracketErr, _ := p.consume(LeftCurlyBracket, "Expected '{' before body function")
	if leftCurlyBracketErr != nil {
		return leftCurlyBracketErr, nil
	}
	err, body := p.block()
	if err != nil {
		return err, nil
	}
	return nil, (NewFunction(name, parameters, body))
}

func (p *Parser) structDeclaration() (error, Statement) {
	identifierErr, name := p.consume(Identifier, "Expected struct name")
	if identifierErr != nil {
		return identifierErr, nil
	}
	leftCurlyBracketErr, _ := p.consume(LeftCurlyBracket, "Expected '{' before class body.")
	if leftCurlyBracketErr != nil {
		return leftCurlyBracketErr, nil
	}
	methods := []*Function{}
	for !p.check(RightCurlyBracket) && !p.isAtEnd() {
		err, fn := p.function()
		if err != nil {
			return err, nil
		}
		if value, ok := fn.(*Function); ok {
			methods = append(methods, value)
		} else {
			return NewParserError("Unexpected expression returned as function method"), nil
		}
	}

	rightCurlyBracketErr, _ := p.consume(RightCurlyBracket, "Expect '}' after class body.")
	if rightCurlyBracketErr != nil {
		return rightCurlyBracketErr, nil
	}
	return nil, NewStructStatment(name, methods)
}

func (p *Parser) declaration() (error, Statement) {
	if p.match(Struct) {
		return p.structDeclaration()
	}
	if p.match(Let) {
		return p.letDeclaration()
	}
	err, statement := p.statement()
	if err != nil {
		return err, nil
	}
	return nil, statement
}

func (p *Parser) parse() (error, []Statement) {
	statements := []Statement{}
	for !p.isAtEnd() {
		err, declaration := p.declaration()
		if err != nil {
			return err, statements
		}
		statements = append(statements, declaration)
	}
	return nil, statements
}
