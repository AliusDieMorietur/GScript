package utils

import (
	"errors"
	"fmt"
	"strconv"
)

func ReturnFirstError(values ...any) error {
	for _, value := range values {
		if err, ok := value.(error); ok {
			return err
		}
	}
	return nil
}

func Ternary[T any](condition bool, a T, b T) T {
	if condition {
		return a
	} else {
		return b
	}
}

func NewError(format string, args ...any) error {
	return errors.New( fmt.Sprintf(format, args...))
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

func AsFloat(value any) (error, float64) {
	if stringValue, ok := value.(string); ok {
		value, err := strconv.ParseFloat(stringValue, 64)
		if (err != nil) {
			return NewError("Can't parse float"), 0.0
		}
		return nil, value
	}
	if floatValue, ok := value.(float64); ok {
		return nil, floatValue
	}
	return NewError("Unexpected value type for float cast '%T'", value), 0.0
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
