package main

import (
	"fmt"
)

const GROUP = "group"

func ExpressionToString(expression Expression) string {
	fmt.Println("expression", expression)
	switch value := expression.(type) {
	case Ternary:
		return ExpressionToString(value.left) + " ? " + ExpressionToString(value.middle) + " : " + ExpressionToString(value.right)
	case Binary:
		return ExpressionToString(value.left) + " " + value.operator.lexeme + " " + ExpressionToString(value.right)
	case Unary:
		return value.operator.lexeme + ExpressionToString(value.right)
	case Grouping:
		return "(" + ExpressionToString(value.expression) + ")"
	case Literal:
		if value.value == nil {
			return "null"
		}
		return fmt.Sprintf("%v", value.value)
	default:
		panic(fmt.Sprintf("I don't know about type %T!\n", value))
	}
}

func parenthesize(name string, expressions ...Expression) string {
	fmt.Println("name", name)
	i := IteratorFrom(expressions)
	si := Map(&i, func(item Expression, _i int, _s []Expression) string {
		return ExpressionToString(item)
	})
	fmt.Println(si)
	s := si.Join(" " + name + " ")
	if name == GROUP {
		s = "(" + s + ")"
	}
	return s
}

// func ExpressionToNumber(expression Expression) float64 {
// 	fmt.Println("expression", expression)
// 	switch value := expression.(type) {
// 	case Ternary:
// 		return panic(Ter)
// 	case Binary:
// 		return ExpressionToString(value.left) + " " + value.operator.lexeme + " " + ExpressionToString(value.right)
// 	case Unary:
// 		return value.operator.lexeme + ExpressionToString(value.right)
// 	case Grouping:
// 		return "(" + ExpressionToString(value.expression) + ")"
// 	case Literal:
// 		if value.value == nil {
// 			return "null"
// 		}
// 		return fmt.Sprintf("%v", value.value)
// 	default:
// 		panic( fmt.Sprintf("I don't know about type %T!\n", value))
// 	}
// }
