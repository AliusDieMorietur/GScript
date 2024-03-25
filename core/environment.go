package main

import (
	u "github.com/core/utils"
)

type Environment struct {
	values map[string]any
}

func NewEnvironment() Environment{
	return Environment{
		values: map[string]any{},
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
	// fmt.Println("e.values", e.values);
}

func (e Environment) get(name Token) (error, any){
	// fmt.Println("e.values", e.values);
	// fmt.Println("name.lexeme", name.lexeme);
	value, ok := e.values[name.lexeme]
	// fmt.Println("value", value);
	if (!ok) {
		return u.NewError("Undefined variable '" + name.lexeme + "'."), nil
	}
	return  nil, value
}
