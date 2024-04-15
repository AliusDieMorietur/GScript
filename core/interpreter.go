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
	globals     *Environment
	environment *Environment
	locals      map[Expression]int
}

func NewInterpreter(locals map[Expression]int) *Interpreter {
	environment := NewEnvironment(nil)
	environment.define("clock", Clock{})
	globals := environment
	return &Interpreter{
		environment,
		globals,
		locals,
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

func (i Interpreter) executeFor(forStatement *ForStatement) error {
	previous := i.environment
	i.environment = NewEnvironment(i.environment)
	initializerErr := i.execute(forStatement.initializer)
	if initializerErr != nil {
		i.environment = previous
		return initializerErr
	}
	conditionErr, condition := i.evaluate(forStatement.condition)
	if conditionErr != nil {
		i.environment = previous
		return conditionErr
	}
	next := func() error {
		conditionErr, _ = i.evaluate(forStatement.increment)
		if conditionErr != nil {
			i.environment = previous
			return conditionErr
		}
		conditionErr, condition = i.evaluate(forStatement.condition)
		if conditionErr != nil {
			i.environment = previous
			return conditionErr
		}
		return nil
	}
	for i.isTruthy(condition) {
		loopErr := i.execute(forStatement.body)
		if loopErr != nil {
			if _, ok := loopErr.(ContinueError); ok {
				err := next()
				if err != nil {
					return err
				}
				continue
			}
			if _, ok := loopErr.(BreakError); ok {
				break
			}
			return loopErr
		}
		err := next()
		if err != nil {
			i.environment = previous
			return err
		}
	}
	i.environment = previous
	return nil
}

func (i Interpreter) executeWhile(whileStatement *WhileStatement) error {
	conditionErr, condition := i.evaluate(whileStatement.condition)
	if conditionErr != nil {
		return conditionErr
	}
	for i.isTruthy(condition) {
		loopErr := i.execute(whileStatement.statement)
		if loopErr != nil {
			if _, ok := loopErr.(ContinueError); ok {
				continue
			}
			if _, ok := loopErr.(BreakError); ok {
				break
			}
			return loopErr
		}
		conditionErr, condition = i.evaluate(whileStatement.condition)
		if conditionErr != nil {
			return conditionErr
		}
	}
	return nil
}

func (i *Interpreter) execute(statement Statement) error {
	switch option := (statement).(type) {
	case *StructStatment:
		i.environment.define(option.name.lexeme, nil)
		methods := map[string]*GSFunction{}
		for _, method := range option.methods {
			fn := NewGSFunction(method, i.environment)
			methods[method.name.lexeme] = fn
		}
		gStruct := NewGSStruct(option.name.lexeme, methods)
		i.environment.assign(option.name, gStruct)
		return nil
	case *ReturnStatement:
		err, value := i.evaluate(option.value)
		if err != nil {
			return err
		}
		return NewReturnError(value)
	case *BreakStatement:
		return NewBreakError()
	case *ContinueStatement:
		return NewContinueError()
	case *BlockStatement:
		return i.executeBlock(option.statements, NewEnvironment(i.environment))
	case *ForStatement:
		return i.executeFor(option)
	case *WhileStatement:
		return i.executeWhile(option)
	case *IfElseStatement:
		err, result := i.evaluate(option.condition)
		if err != nil {
			return err
		}
		if i.isTruthy(result) {
			err := i.execute(option.thenBranch)
			if err != nil {
				return err
			}
		} else {
			if option.elseBranch != nil {
				err := i.execute(option.elseBranch)
				if err != nil {
					return err
				}
			}
		}
	case *LetStatement:
		var value any
		if option.initializer != nil {
			err, result := i.evaluate(option.initializer)
			if err != nil {
				return err
			}
			value = result
		}
		i.environment.define(option.name.lexeme, value)
		return nil
	case *PrintStatement:
		err, value := i.evaluate(option.expression)
		if err != nil {
			return err
		}
		if callee, ok := value.(Callable); ok {
			fmt.Println(callee.String())
			return nil
		}
		fmt.Println(value)
		return nil
	case *ExpressionStatement:
		err, _ := i.evaluate(option.expression)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i Interpreter) lookUpVariable(name *Token, variable Expression) (error, any) {
	distance, ok := i.locals[variable]
	if ok {
		return i.environment.getAt(distance, name.lexeme)
	} else {
		return i.globals.get(name)
	}
}

func (i Interpreter) evaluate(expression Expression) (error, any) {
	switch option := (expression).(type) {
	case *Get:
		err, value := i.evaluate(option.object)
		if err != nil {
			return err, nil
		}
		if instance, ok := value.(*GSInstance); ok {
			return instance.get(option.name)
		}
		return NewRuntimeError("Only instances have property names"), nil
	case *Set:
		err, value := i.evaluate(option.object)
		if err != nil {
			return err, nil
		}
		if instance, ok := value.(*GSInstance); ok {
			err, value := i.evaluate(option.value)
			if err != nil {
				return err, nil
			}
			instance.set(option.name, value)
			return nil, value
		}
		return NewRuntimeError("Only instances have property names"), nil
	case *Function:
		f := NewGSFunction(option, i.environment)
		if option.name.lexeme != AnonymusFunction {
			i.environment.define(option.name.lexeme, f)
		}
		return nil, f
	case *Call:
		err, callee := i.evaluate(option.callee)
		if err != nil {
			return err, nil
		}
		arguments := []any{}
		for _, argument := range option.arguments {
			err, expression := i.evaluate(argument)
			if err != nil {
				return err, nil
			}
			arguments = append(arguments, expression)
		}
		if fn, ok := callee.(Callable); ok {
			argumentsQuantity := len(arguments)
			arity := fn.arity()
			if argumentsQuantity != arity {
				return u.NewError("Expected %v arguments but got %v", arity, argumentsQuantity), nil
			}
			return fn.call(&i, arguments)
		}
		return u.NewError("Can only call functions"), nil
	case *Logical:
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
	case *Ternary:
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
	case *Variable:
		return i.lookUpVariable(option.name, expression)
	case *Assignment:
		err, value := i.evaluate(option.value)
		if err != nil {
			return err, nil
		}
		distance, ok := i.locals[expression]
		if ok {
			i.environment.assignAt(distance, option.name, value)
		} else {
			assignErr := i.globals.assign(option.name, value)
			if assignErr != nil {
				return assignErr, nil
			}
		}
		return nil, value
	case *Binary:
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
	case *Unary:
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
	case *Grouping:
		return i.evaluate(option.expression)
	case *Literal:
		return nil, option.value
	default:
		return u.NewError("Unreachable evaluate"), nil
	}
}
