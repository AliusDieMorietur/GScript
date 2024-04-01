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

type Variable struct {
	name Token
}

func NewVariable(name Token) Variable {
	return Variable{
		name,
	}
}

type Assignment struct {
	name  Token
	value Expression
}

func NewAssignment(name Token, value Expression) Assignment {
	return Assignment{
		name,
		value,
	}
}

type Logical struct {
	left     Expression
	operator Token
	right    Expression
}

func NewLogical(left Expression, operator Token, right Expression) Logical {
	return Logical{
		left,
		operator,
		right,
	}
}

type Call struct {
	callee    Expression
	paren     Token
	arguments []Expression
}

func NewCall(callee Expression, paren Token, arguments []Expression) Call {
	return Call{
		callee,
		paren,
		arguments,
	}
}



type Callable interface {
	arity() int
	call(i *Interpreter, arguments []any) (error, any)
	toString() string
}

func findToken(token Token, expression Expression) bool {
	switch option := expression.(type) {
	case Assignment:
		if option.name.lexeme != token.lexeme {
			return findToken(token, option.value)
		}
		return true
	case Variable:
		return option.name.lexeme == token.lexeme
	case Unary:
		return findToken(token, option.right)
	case Binary:
		return findToken(token, option.left) || findToken(token, option.right)
	case Grouping:
		return findToken(token, option.expression)
	default:
		return false
	}
}
