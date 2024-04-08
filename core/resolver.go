package main

import "fmt"

type Resolver struct {
	interpreter *Interpreter
	scopes      []map[string]bool
}

func NewResolver(interpreter *Interpreter) Resolver {
	scopes := []map[string]bool{}
	return Resolver{
		interpreter,
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
		fmt.Print("Warning: empty scope")
	}
}

func (r Resolver) resolve(value any) {
	switch option := value.(type) {
	case []any:
		for _, value := range option {
			r.resolve(value)
		}
	case BlockStatement:
		r.beginScope()
		r.resolve(option.statements)
		r.endScope()
	case LetStatement:
		// r.declare(option.name)
		if option.initializer != nil {
			r.resolve(option.initializer)
		}
		// r.define(option.name)
	}
}
