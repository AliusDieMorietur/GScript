package main

type Expression = any

type Ternary struct {
	left   Expression
	middle Expression
	right  Expression
}

func NewTernary(left Expression, middle Expression, right Expression) *Ternary {
	return &Ternary{
		left,
		middle,
		right,
	}
}

type Binary struct {
	left     Expression
	operator *Token
	right    Expression
}

func NewBinary(left Expression, operator *Token, right Expression) *Binary {
	return &Binary{
		left,
		operator,
		right,
	}
}

type Grouping struct {
	expression Expression
}

func NewGrouping(expression Expression) *Grouping {
	return &Grouping{
		expression,
	}
}

type Literal struct {
	value any
}

func NewLiteral(value any) *Literal {
	return &Literal{
		value,
	}
}

type Unary struct {
	operator *Token
	right    Expression
}

func NewUnary(operator *Token, right Expression) *Unary {
	return &Unary{
		operator,
		right,
	}
}

type Variable struct {
	name *Token
}

func NewVariable(name *Token) *Variable {
	return &Variable{
		name,
	}
}

type Assignment struct {
	name  *Token
	value Expression
}

func NewAssignment(name *Token, value Expression) *Assignment {
	return &Assignment{
		name,
		value,
	}
}

type Logical struct {
	left     Expression
	operator *Token
	right    Expression
}

func NewLogical(left Expression, operator *Token, right Expression) *Logical {
	return &Logical{
		left,
		operator,
		right,
	}
}

type Call struct {
	callee    Expression
	paren     *Token
	arguments []Expression
}

func NewCall(callee Expression, paren *Token, arguments []Expression) *Call {
	return &Call{
		callee,
		paren,
		arguments,
	}
}

type Function struct {
	name       *Token
	parameters []*Token
	body       []Statement
}

func NewFunction(name *Token, parameters []*Token, body []Statement) *Function {
	return &Function{
		name,
		parameters,
		body,
	}
}

type Callable interface {
	arity() int
	call(i *Interpreter, arguments []any) (error, any)
	String() string
}

type Get struct {
	name   *Token
	object any
}

func NewGet(name *Token, object any) *Get {
	return &Get{
		name,
		object,
	}
}

type Set struct {
	name   *Token
	object Expression
	value  Expression
}

func NewSet(name *Token, object Expression,
	value Expression) *Set {
	return &Set{
		name,
		object,
		value,
	}
}

type ThisExpression struct {
	keyword *Token
}

func NewThisExpression(keyword *Token) *ThisExpression {
	return &ThisExpression{
		keyword,
	}
}
