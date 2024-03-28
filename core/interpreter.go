package main

import (
	"fmt"

	u "github.com/core/utils"
)

func performBinaryNumberOperation(left any, right any, operation func(lhs float64, rhs float64) any) (error, any) {
	leftErr, lhs := u.AsFloat(left)
	if leftErr != nil {
		return leftErr, nil
	}
	rightErr, rhs := u.AsFloat(right)
	if rightErr != nil {
		return rightErr, nil
	}
	return nil, operation(lhs, rhs)
}

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() Interpreter {
	environment := NewEnvironment(nil)
	return Interpreter{
		environment: &environment,
	}
}

func (i *Interpreter) interpret(statements []Statement) error {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
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

func (i *Interpreter) executeBlock(statements []Statement, env *Environment) error {
	previous := i.environment
	i.environment = env
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			i.environment = previous
			return err
		}
	}
	i.environment = previous
	return nil
}

func (i *Interpreter) execute(statement Statement) error {
	switch option := statement.(type) {
	case BlockStatement:
		env := NewEnvironment(i.environment)
		err := i.executeBlock(option.statements, &env)
		if err != nil {
			return err
		}
	case IfElseStatement:
		err, result := i.evaluate(option.condition)
		if err != nil {
			return err
		}
		if i.isTruthy(result) {
			i.execute(option.thenBranch)
		} else {
			if option.elseBranch != nil {
				i.execute(option.elseBranch)
			}
		}
	case LetStatement:
		var value any
		if option.initializer != nil {
			// if (findToken(option.name, option.initializer)) {
			// 	return u.NewError("Access variable '%v' before declaration", option.name.lexeme)
			// }
			err, result := i.evaluate(option.initializer)
			if err != nil {
				return err
			}
			value = result
		}
		i.environment.define(option.name.lexeme, value)
	case PrintStatement:
		err, value := i.evaluate(option.expression)
		if err != nil {
			return err
		}
		fmt.Println(value)
	case ExpressionStatement:
		err, _ := i.evaluate(option.expression)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i Interpreter) evaluate(expression Expression) (error, any) {
	switch option := expression.(type) {
	case Logical:
		err, left := i.evaluate(option.left)
		if err != nil {
			return err, nil
		}
		if option.operator.tokenType == Or {
			if i.isTruthy(left) {
				return nil, left
			}
		} else {
			if !i.isTruthy(left) {
				return nil, left
			}
		}
		return i.evaluate(option.right)
	case Ternary:
		leftErr, left := i.evaluate(option.left)
		if leftErr != nil {
			return leftErr, nil
		}
		if i.isTruthy(left) {
			middleErr, middle := i.evaluate(option.middle)
			if middleErr != nil {
				return middleErr, nil
			}
			return nil, middle
		} else {
			rightErr, right := i.evaluate(option.right)
			if rightErr != nil {
				return rightErr, nil
			}
			return nil, right
		}
	case Variable:
		return i.environment.get(option.name)
	case Assignment:
		err, value := i.evaluate(option.value)
		if err != nil {
			return err, nil
		}
		assignErr := i.environment.assign(option.name.lexeme, value)
		if assignErr != nil {
			return assignErr, nil
		}
		return nil, value
	case Binary:
		leftErr, left := i.evaluate(option.left)
		if leftErr != nil {
			return leftErr, nil
		}
		rightErr, right := i.evaluate(option.right)
		if rightErr != nil {
			return rightErr, nil
		}
		switch operator := option.operator.tokenType; operator {
		case Greater:
			return performBinaryNumberOperation(left, right, func(lhs float64, rhs float64) any {
				return lhs > rhs
			})
		case GreaterEqual:
			return performBinaryNumberOperation(left, right, func(lhs float64, rhs float64) any {
				return lhs >= rhs
			})
		case Less:
			return performBinaryNumberOperation(left, right, func(lhs float64, rhs float64) any {
				return lhs < rhs
			})
		case LessEqual:
			return performBinaryNumberOperation(left, right, func(lhs float64, rhs float64) any {
				return lhs <= rhs
			})
		case BangEqual:
			return nil, !i.isEqual(left, right)
		case EqualEqual:
			return nil, i.isEqual(left, right)
		case Minus:
			return performBinaryNumberOperation(left, right, func(lhs float64, rhs float64) any {
				return lhs - rhs
			})
		case Star:
			return performBinaryNumberOperation(left, right, func(lhs float64, rhs float64) any {
				return lhs * rhs
			})
		case Slash:
			return performBinaryNumberOperation(left, right, func(lhs float64, rhs float64) any {
				return lhs / rhs
			})
		case Plus:
			if u.IsString(left) && u.IsString(right) {
				lhs := u.AsString(left)
				rhs := u.AsString(right)
				return nil, lhs + rhs
			}
			if u.IsFloat(left) && u.IsFloat(right) {
				leftErr, lhs := u.AsFloat(left)
				if leftErr != nil {
					return leftErr, nil
				}
				rightErr, rhs := u.AsFloat(right)
				if rightErr != nil {
					return rightErr, nil
				}
				return nil, lhs + rhs
			}
			return u.NewError("Unexpected plus types T1:'%T' T2:'%T'", left, right), nil
		default:
			return u.NewError("Unexpected binary operator '%v'", operator), nil
		}
	case Unary:
		rightErr, right := i.evaluate(option.right)
		if rightErr != nil {
			return rightErr, nil
		}
		switch operator := option.operator.tokenType; operator {
		case Minus:
			err, rhs := u.AsFloat(right)
			if err != nil {
				return err, nil
			}
			return nil, -rhs
		case Bang:
			return nil, i.isTruthy(right)
		default:
			return u.NewError("Unexpected unary operator '%v'", operator), nil
		}
	case Grouping:
		return i.evaluate(option.expression)
	case Literal:
		return nil, option.value
	default:
		return u.NewError("Unreachable evaluate"), nil
	}
}
