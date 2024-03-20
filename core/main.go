package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Lng struct {
	hadError bool
}

func NewLng() Lng {
	return Lng{hadError: false}
}

func (l *Lng) run(source string) {
	scanner := NewScanner(source, func() {
		l.hadError = true
	})
	tokens := scanner.scanTokens()
	// for _, token := range tokens {
	// 	fmt.Println("Token", token.ToString())
	// }
	if l.hadError {
		fmt.Println("LNG error end")
		return
	}
	parser := NewParser(tokens)
	statements := parser.parse()

	interperter := NewInterpreter()
	interperter.interpret(statements)

	// fmt.Println(ExpressionToString(expression))

}

func (l *Lng) runFile(filePath string) {
	data, _ := os.ReadFile(filePath)
	l.run(string(data))
}

func (l *Lng) runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	text, _ := reader.ReadString('\n')
	l.run(text)
	if strings.Trim(text, "\n") == "exit" {
		return
	}
	l.runPrompt()
}

func main() {
	lng := NewLng()
	if len(os.Args) > 1 {
		lng.runFile(os.Args[1])
	} else {
		lng.runPrompt()
	}
	// expression := NewBinary(
	// 	NewUnary(
	// 		NewToken(Minus, "-", nil, 1),
	// 		NewLiteral(123)),
	// 	NewToken(Star, "*", nil, 1),
	// 	NewGrouping(
	// 		NewLiteral(45.67)))

	// expression := NewBinary(
	// 	NewGrouping(
	// 		NewBinary(NewLiteral(1), NewToken(Minus, "-", nil, 1), NewLiteral(2)),
	// 	),
	// 	NewToken(Star, "*", nil, 1),
	// 	NewGrouping(
	// 		NewBinary(NewLiteral(3), NewToken(Plus, "+", nil, 1), NewLiteral(4)),
	// 	))

	// fmt.Println(ExpressionToString(expression))
}
