package main

import "fmt"

type GSFunction struct {
	declaration *Function
	closure     *Environment
}

func NewGSFunction(declaration *Function, closure *Environment) *GSFunction {
	return &GSFunction{
		declaration,
		closure,
	}
}

func (f GSFunction) arity() int {
	return len(f.declaration.parameters)
}

func (f GSFunction) toString() string {
	return fmt.Sprintf("[fn: %v]", f.declaration.name.lexeme)
}

func (f GSFunction) call(i *Interpreter, arguments []any) (error, any) {
	environment := NewEnvironment(f.closure)
	for i := 0; i < len(f.declaration.parameters); i++ {
		environment.define(f.declaration.parameters[i].lexeme, arguments[i])
	}
	err := i.executeBlock(f.declaration.body, environment)
	if rErr, ok := err.(ReturnError); ok {
		return nil, rErr.value
	}
	return err, nil
}
