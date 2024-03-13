package main

import (
	"os"
	"strings"
)

func defineType(expressionType string) string {
	splitted := strings.Split(expressionType, ":")
	structName := strings.TrimSpace(splitted[0])
	fields := strings.TrimSpace(splitted[1])
	fieldsSplitted := strings.Split(fields, ",")
	src := "type " + structName + " struct {\n"
	src += strings.Join(fieldsSplitted, "\n") + "\n"
	src += "}\n\n"
	src += "func New" + structName + "(" + fields + ") " + structName + " {\n"
	src += "return " + structName + "{\n"
	for _, field := range fieldsSplitted {
		fieldName := strings.Split(strings.TrimSpace(field), " ")[0]
		src += fieldName + ",\n"
	}
	src += "}\n}\n\n"
	return src
}

func DefineAst(filePath string,  expressionTypes []string) {
	src := "package main\n\n"
	src += "type Expression interface {\n}\n\n"
	// src += "  eval()\n"
	// src += "}\n"
	for _, expressionType := range expressionTypes {
		src += defineType(expressionType)
	}

	err := os.WriteFile(filePath, []byte(src), 0666)
	if err != nil {
		panic(err)
	}
}


