package main

import (
	u "github.com/core/utils"
)

func NewParserError(message string, args ...any) error {
	return u.NewError("Syntax error: "+message, args...)
}

func NewRuntimeError(format string, args ...any) error {
	return u.NewError("Runtime error: "+format, args...)
}

type BreakError struct {
	message string
}

func (e BreakError) Error() string {
	return e.message
}

func NewBreakError() BreakError {
	return BreakError{
		message: "Break outside loop",
	}
}

type ContinueError struct {
	message string
}

func (e ContinueError) Error() string {
	return e.message
}

func NewContinueError() ContinueError {
	return ContinueError{
		message: "Continue outside loop",
	}
}

type ReturnError struct {
	message string
	value   any
}

func NewReturnError(value any) ReturnError {
	message := NewRuntimeError("Return error").Error()
	return ReturnError{
		message,
		value,
	}
}

func (r ReturnError) Error() string {
	return r.message
}
