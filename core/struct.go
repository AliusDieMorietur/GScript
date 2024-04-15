package main

type GSStruct struct {
	name    string
	methods map[string]*GSFunction
}

func NewGSStruct(name string, methods map[string]*GSFunction) *GSStruct {
	return &GSStruct{
		name,
		methods,
	}
}

func (g GSStruct) String() string {
	return "[struct: " + g.name + "]" 
}

func (g GSStruct) arity() int {
	return 0
}

func (g GSStruct) call(i *Interpreter, arguments []any) (error, any) {
	instance := NewGSInstance(&g)
	return nil, instance
}

func (g GSStruct) findMethod(name string) any {
	method, ok := g.methods[name]
	if ok {
		return method
	}
	return nil
}

type GSInstance struct {
	gsStruct *GSStruct
	fields   map[string]any
}

func NewGSInstance(gsStruct *GSStruct) *GSInstance {
	fields := map[string]any{}
	return &GSInstance{
		gsStruct,
		fields,
	}
}

func (g GSInstance) String() string {
	return "[instance: " + g.gsStruct.name + "]"
}

func (g GSInstance) get(name *Token) (error, any) {
	value, ok := g.fields[name.lexeme]
	if ok {
		return nil, value
	}
	method := g.gsStruct.findMethod(name.lexeme)
	if method != nil {
		return nil, method
	}
	return NewRuntimeError("Undefined property '" + name.lexeme + "'"), nil
}

func (g GSInstance) set(name *Token, value any) {
	g.fields[name.lexeme] = value
}
