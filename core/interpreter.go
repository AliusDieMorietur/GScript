// program → statement* EOF ;
// statement → exprStmt | printStmt ;
// exprStmt → expression ";" ;
// printStmt → "print" expression ";" ;

package main

import (
	"fmt"

	u "github.com/core/utils"
)

type Interpreter struct {
}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i Interpreter) interpret(statements []Statement) {
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i Interpreter) isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if boolValue, ok := value.(bool); ok {
		return boolValue
	}
	return true
}

func (i Interpreter) isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a == b
}

func (i Interpreter) execute(statement Statement) {
	switch option := statement.(type) {
	case PrintStatement:
		value := i.evaluate(option.expression)
		fmt.Println(value)
	case ExpressionStatement:
		i.evaluate(option.expression)
	}
}

func (i Interpreter) evaluate(expression Expression) any {
	switch option := expression.(type) {
	// case Ternary:
	// 	return ExpressionToString(value.left) + " ? " + ExpressionToString(value.middle) + " : " + ExpressionToString(value.right)
	case Binary:
		left := i.evaluate(option.left)
		right := i.evaluate(option.right)
		switch operator := option.operator.tokenType; operator {
		case Greater:
			return u.AsFloat(left) > u.AsFloat(right)
		case GreaterEqual:
			return u.AsFloat(left) >= u.AsFloat(right)
		case Less:
			return u.AsFloat(left) < u.AsFloat(right)
		case LessEqual:
			return u.AsFloat(left) <= u.AsFloat(right)
		case BangEqual:
			return !i.isEqual(left, right)
		case EqualEqual:
			return i.isEqual(left, right)
		case Minus:
			return u.AsFloat(left) - u.AsFloat(right)
		case Star:
			return u.AsFloat(left) * u.AsFloat(right)
		case Slash:
			return u.AsFloat(left) / u.AsFloat(right)
		case Plus:
			if u.IsString(left) && u.IsString(right) {
				return u.AsString(left) + u.AsString(right)
			}
			if u.IsFloat(left) && u.IsFloat(right) {
				return u.AsFloat(left) + u.AsFloat(right)
			}
			panic(fmt.Sprintf("Unexpected plus types T1:'%T' T2:'%T'", left, right))
		default:
			panic(fmt.Sprintf("Unexpected binary operator '%v'", operator))
		}
	case Unary:
		right := i.evaluate(option.right)
		switch operator := option.operator.tokenType; operator {
		case Minus:
			rhs := u.AsFloat(right)
			return -rhs
		case Bang:
			return i.isTruthy(right)
		default:
			panic(fmt.Sprintf("Unexpected unary operator '%v'", operator))
		}
	case Grouping:
		return i.evaluate(option.expression)
	case Literal:
		return option.value
	default:
		panic("Unreachable evaluate")
	}
}
