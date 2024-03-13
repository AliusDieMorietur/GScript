package main

type Expression interface {
}

type Ternary struct {
	left   Expression
	middle Expression
	right  Expression
}

func NewTernary(left Expression, middle Expression, right Expression) Ternary {
	return Ternary{
		left,
		middle,
		right,
	}
}

type Binary struct {
	left     Expression
	operator Token
	right    Expression
}

func NewBinary(left Expression, operator Token, right Expression) Binary {
	return Binary{
		left,
		operator,
		right,
	}
}

type Grouping struct {
	expression Expression
}

func NewGrouping(expression Expression) Grouping {
	return Grouping{
		expression,
	}
}

type Literal struct {
	value any
}

func NewLiteral(value any) Literal {
	return Literal{
		value,
	}
}

type Unary struct {
	operator Token
	right    Expression
}

func NewUnary(operator Token, right Expression) Unary {
	return Unary{
		operator,
		right,
	}
}
