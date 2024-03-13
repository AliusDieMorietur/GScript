package main

import (
	"fmt"
	"os"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		return
	} else {
		filePath := os.Args[1]
		types := []string{
			"Binary : left Expression, operator Token, right Expression",
			"Grouping : expression Expression",
			"Literal : value any",
			"Unary : operator Token, right Expression",
		}
		fmt.Println("filePath", filePath)
		DefineAst(filePath,  types)
	}
}
