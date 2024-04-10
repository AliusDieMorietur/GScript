package main

import (
	"fmt"
)

type Resolver struct {
	locals map[Expression]int
	scopes []map[string]bool
}

func NewResolver() *Resolver {
	scopes := []map[string]bool{}
	locals := map[Expression]int{}
	return &Resolver{
		locals,
		scopes,
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) endScope() {
	if len(r.scopes) > 0 {
		r.scopes = r.scopes[:len(r.scopes)-1]
	} else {
		fmt.Print("Warning: empty scopes")
	}
}

func (r *Resolver) isScopesEmpty() bool {
	return len(r.scopes) == 0
}

func (r *Resolver) declare(name *Token) {
	if r.isScopesEmpty() {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = false
}

func (r *Resolver) define(name *Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = true
}

func (r *Resolver) resolveLocal(expression Expression, name *Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		scope := r.scopes[i]
		_, ok := scope[name.lexeme]
		if ok {
			r.locals[expression] = len(r.scopes) - 1 - i
			return
		}
	}
}

func (r *Resolver) resolveFunction(fn *Function) error {
	r.beginScope()
	for _, parameter := range fn.parameters {
		r.declare(parameter)
		r.define(parameter)
	}
	err, _ := r.resolve(fn.body)
	if err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) resolveExpression(expression Expression) error {
	switch option := (expression).(type) {
	case *Grouping:
		return r.resolveExpression(option.expression)
	case *Literal:
		return nil
	case *Logical:
		errLeft := r.resolveExpression(option.left)
		if errLeft != nil {
			return errLeft
		}
		errRight := r.resolveExpression(option.right)
		if errRight != nil {
			return errRight
		}
		return nil
	case *Unary:
		return r.resolveExpression(option.right)
	case *Call:
		errCalee := r.resolveExpression(option.callee)
		if errCalee != nil {
			return errCalee
		}
		for _, argument := range option.arguments {
			err := r.resolveExpression(argument)
			if err != nil {
				return err
			}
		}
		return nil
	case *Binary:
		errLeft := r.resolveExpression(option.left)
		if errLeft != nil {
			return errLeft
		}
		errRight := r.resolveExpression(option.right)
		if errRight != nil {
			return errRight
		}
		return nil
	case *Function:
		r.declare(option.name)
		r.define(option.name)
		err := r.resolveFunction(option)
		if err != nil {
			return err
		}
	case *Variable:
		if !r.isScopesEmpty() {
			scope := r.scopes[len(r.scopes)-1]
			v, ok := scope[option.name.lexeme]
			if ok && !v {
				return NewResolveError("Can't read local variable in its own initializer.")
			}
		}
		r.resolveLocal(expression, option.name)
		return nil
	case *Assignment:
		err := r.resolveExpression(option.value)
		if err != nil {
			return err
		}
		r.resolveLocal(expression, option.name)
	}
	return nil
}

func (r *Resolver) resolveStatement(statement Statement) error {
	switch option := (statement).(type) {
	case *ReturnStatement:
		if option.value == nil {
			return nil
		}
		err := r.resolveExpression(option.value)
		if err != nil {
			return err
		}
	case *BreakStatement:
		return nil
	case *ContinueStatement:
		return nil
	case *BlockStatement:
		r.beginScope()
		err, _ := r.resolve(option.statements)
		if err != nil {
			return err
		}
		r.endScope()
	case *ForStatement:
		initializerErr := r.resolveStatement(option.initializer)
		if initializerErr != nil {
			return initializerErr
		}
		conditionErr := r.resolveExpression(option.condition)
		if conditionErr != nil {
			return conditionErr
		}
		incrementErr := r.resolveExpression(option.increment)
		if incrementErr != nil {
			return incrementErr
		}
		bodyErr := r.resolveStatement(option.body)
		if bodyErr != nil {
			return bodyErr
		}
	case *WhileStatement:
		conditionErr := r.resolveExpression(option.condition)
		if conditionErr != nil {
			return conditionErr
		}
		statementErr := r.resolveStatement(option.statement)
		if statementErr != nil {
			return statementErr
		}
	case *IfElseStatement:
		conditionErr := r.resolveExpression(option.condition)
		if conditionErr != nil {
			return conditionErr
		}
		statementErr := r.resolveStatement(option.thenBranch)
		if statementErr != nil {
			return statementErr
		}
		if option.elseBranch != nil {
			statementErr := r.resolveStatement(option.elseBranch)
			if statementErr != nil {
				return statementErr
			}
		}
	case *LetStatement:
		r.declare(option.name)
		if option.initializer != nil {
			err := r.resolveExpression(option.initializer)
			if err != nil {
				return err
			}
		}
		r.define(option.name)
	case *PrintStatement:
		err := r.resolveExpression(option.expression)
		if err != nil {
			return err
		}
	case *ExpressionStatement:
		err := r.resolveExpression(option.expression)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolve(statements []Statement) (error, map[Expression]int) {
	for _, value := range statements {
		err := r.resolveStatement(value)
		if err != nil {
			return err, nil
		}
	}
	return nil, r.locals
}
