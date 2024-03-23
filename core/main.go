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

func (l *Lng) run(source string) error {
	scanner := NewScanner(source)
	scanErr, tokens := scanner.scanTokens()
	if scanErr != nil {
		return scanErr
	}
	parser := NewParser(tokens)
	parseErr, statements := parser.parse()
	if parseErr != nil {
		return parseErr
	}
	interperter := NewInterpreter()
	interpretErr := interperter.interpret(statements)
	if interpretErr != nil {
		return interpretErr
	}
	return nil
}

func (l *Lng) runFile(filePath string) {
	data, _ := os.ReadFile(filePath)
	err := l.run(string(data))
	if err != nil {
		fmt.Println(err)
	}
}

func (l *Lng) runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	text, _ := reader.ReadString('\n')
	if strings.Trim(text, "\n") == "exit" {
		return
	}
	err := l.run(text)
	if err != nil {
		fmt.Println(err)
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
