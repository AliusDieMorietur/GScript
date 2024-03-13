package main

import "fmt"

type Iterator[T any] struct {
	Data []T
}

func NewIterator[T any]() Iterator[T] {
	return Iterator[T]{
		Data: []T{},
	}
}

func (i *Iterator[T]) Push(item T) {
	i.Data = append(i.Data, item)
}

func (i *Iterator[T]) Join(separator string) string {
	s := ""
	for index, value := range i.Data {
		if index != 0 {
			s += separator
		}
		s += fmt.Sprintf("%v", value)
	}
	return s
}

func IteratorFrom[T any](Data []T) Iterator[T] {
	return Iterator[T]{
		Data,
	}
}

func Map[T any, U any](i *Iterator[T], mapper func(T, int, []T) U) Iterator[U] {
	result := make([]U, len(i.Data))
	for index, value := range i.Data {
		result[index] = mapper(value, index, i.Data)
	}
	return IteratorFrom[U](result)
}

func (i *Iterator[T]) Filter(filter func(T, int, []T) bool) {
	result := []T{}
	for index, value := range i.Data {
		condition := filter(value, index, i.Data)
		if condition {
			result = append(result, value)
		}
	}
	i.Data = result
}
