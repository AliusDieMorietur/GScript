package main

import (
	u "github.com/core/utils"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) Environment {
	values := map[string]any{}
	return Environment{
		values,
		enclosing,
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) assign(name string, value any) error {
	_, ok := e.values[name]
	if ok {
		e.values[name] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}
	return u.NewError("Undefined variable '" + name + "'.")
}

func (e Environment) get(name Token) (error, any) {
	value, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}
		return u.NewError("Undefined variable '" + name.lexeme + "'."), nil
	}
	return nil, value
}

func (e Environment) has(name Token) bool {
	_, ok := e.values[name.lexeme]
	return ok
}
