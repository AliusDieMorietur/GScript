package main

type Statement interface {
}

type PrintStatement struct {
	expression Expression
}

func NewPrintStatement(expression Expression) PrintStatement {
	return PrintStatement{
		expression,
	}
}

type ExpressionStatement struct {
	expression Expression
}

func NewExpressionStatement(expression Expression) ExpressionStatement {
	return ExpressionStatement{
		expression,
	}
}

type LetStatement struct {
	name       Token
	initializer Expression
}

func NewLetStatement(name Token, initializer Expression) LetStatement {
	return LetStatement{
		name,
		initializer,
	}
}
