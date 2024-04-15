package main

import (
	"fmt"
	"time"
)

type Clock struct {
}

func (c Clock) arity() int {
	return 0
}

func (c Clock) call(i *Interpreter, arguments []any) (error, any) {
	fmt.Println("time.Now().UnixNano()", time.Now().UnixMilli())
	return nil, float64(time.Now().UnixMilli())
}

func (c Clock) String() string {
	return "[fn: clock]"
}
