package main

type Statement = any

type PrintStatement struct {
	expression Expression
}

func NewPrintStatement(expression Expression) *PrintStatement {
	return &PrintStatement{
		expression,
	}
}

type ExpressionStatement struct {
	expression Expression
}

func NewExpressionStatement(expression Expression) *ExpressionStatement {
	return &ExpressionStatement{
		expression,
	}
}

type LetStatement struct {
	name        *Token
	initializer Expression
}

func NewLetStatement(name *Token, initializer Expression) *LetStatement {
	return &LetStatement{
		name,
		initializer,
	}
}

type BlockStatement struct {
	statements []Statement
}

func NewBlockStatement(statements []Statement) *BlockStatement {
	return &BlockStatement{
		statements,
	}
}

type IfElseBranch struct {
	condition Expression
	branch    Statement
}

type IfElseStatement struct {
	condition  Expression
	thenBranch Statement
	// elseIfs []IfElseBranch
	elseBranch Statement
}

func NewIfStatement(condition Expression,
	thenBranch Statement,
	// elseIfs []IfElseBranch,
	elseBranch Statement) *IfElseStatement {
	return &IfElseStatement{
		condition,
		thenBranch,
		// elseIfs,
		elseBranch,
	}
}

type WhileStatement struct {
	condition Expression
	statement Statement
}

func NewWhileStatement(condition Expression, statement Statement) *WhileStatement {
	return &WhileStatement{
		condition,
		statement,
	}
}

type ForStatement struct {
	condition   Expression
	initializer Statement
	increment   Expression
	body        Statement
}

func NewForStatement(condition Expression,
	initializer Statement,
	increment Expression,
	body Statement) *ForStatement {
	return &ForStatement{
		condition,
		initializer,
		increment,
		body,
	}
}

type BreakStatement struct{}

func NewBreakStatement() *BreakStatement {
	return &BreakStatement{}
}

type ContinueStatement struct{}

func NewContinueStatement() *ContinueStatement {
	return &ContinueStatement{}
}

type ReturnStatement struct {
	value Expression
}

func NewReturnStatement(value Expression) *ReturnStatement {
	return &ReturnStatement{
		value,
	}
}

type StructStatment struct {
	name    *Token
	methods []*Function
}

func NewStructStatment(name *Token, methods []*Function) *StructStatment {
	return &StructStatment{
		name,
		methods,
	}
}
