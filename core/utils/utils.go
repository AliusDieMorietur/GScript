package utils

import (
	"fmt"
	"strconv"
)

func Ternary[T any](condition bool, a T, b T) T {
	if condition {
		return a
	} else {
		return b
	}
}

func Error(line uint, message string) {
	Report(line, "", message)
}

func Report(line uint, place string, message string) {
	fmt.Println(fmt.Sprintf("[line %d] Error%s: %s", line, place, message))
}

func Expect( err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func AsFloat(value any) float64 {
	if stringValue, ok := value.(string); ok {
		value, err := strconv.ParseFloat(stringValue, 64)
		Expect(err, "Can't parse float")
		return value
	}
	if floatValue, ok := value.(float64); ok {
		return floatValue
	}
	panic(fmt.Sprint("Unexpected value type for float cast '%T'", value))
}

func AsString(value any) string {
	if stringValue, ok := value.(string); ok {
		return stringValue
	}
	return fmt.Sprint("%v", value)
}

func IsString(value any) bool {
	_, ok := value.(string);
		return ok
}

func IsFloat(value any) bool {
	_, ok := value.(float64);
		return ok
}
