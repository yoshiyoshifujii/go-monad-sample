package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

type Maybe[T any] interface {
	bind(func(T) Maybe[T]) Maybe[T]
	getValue() (T, error)
}

type Just[T any] struct {
	value T
}

type Nothing[T any] struct{}

func justOf[T any](value T) Maybe[T] {
	return Just[T]{value: value}
}

func nothingOf[T any]() Maybe[T] {
	return Nothing[T]{}
}

func (j Just[T]) bind(f func(T) Maybe[T]) Maybe[T] {
	return f(j.value)
}

func (j Just[T]) getValue() (T, error) {
	return j.value, nil
}

func (n Nothing[T]) bind(f func(T) Maybe[T]) Maybe[T] {
	return n
}

func (n Nothing[T]) getValue() (T, error) {
	var zero T
	return zero, fmt.Errorf("no value for this type")
}

// Left Identity Law
func checkLeftIdentityLaw[T constraints.Ordered](value T, f func(T) Maybe[T]) bool {
	// unit(a) bind f
	left := justOf(value).bind(f)
	// f(a)
	right := f(value)

	leftValue, leftError := left.getValue()
	rightValue, rightError := right.getValue()
	if leftError != nil && rightError != nil {
		return true
	}
	if leftError != nil || rightError != nil {
		return false
	}
	return leftValue == rightValue
}

// Right Identity Law
func checkRightIdentityLaw[T constraints.Ordered](m Maybe[T]) bool {
	// m bind unit
	left := m.bind(justOf[T])

	leftValue, leftError := left.getValue()
	rightValue, rightError := m.getValue()
	if leftError != nil && rightError != nil {
		return true
	}
	if leftError != nil || rightError != nil {
		return false
	}
	return leftValue == rightValue
}

// Associativity Law
func checkAssociativeLaw[T constraints.Ordered](m Maybe[T], f func(T) Maybe[T], g func(T) Maybe[T]) bool {
	// (m bind f) bind g
	left := m.bind(f).bind(g)
	// m bind (x => f(x) bind g)
	right := m.bind(func(x T) Maybe[T] {
		return f(x).bind(g)
	})

	leftValue, leftError := left.getValue()
	rightValue, rightError := right.getValue()
	if leftError != nil && rightError != nil {
		return true
	}
	if leftError != nil || rightError != nil {
		return false
	}
	return leftValue == rightValue
}

func main() {
	increment := func(x int) Maybe[int] {
		return justOf(x + 1)
	}
	double := func(x int) Maybe[int] {
		return justOf(x * 2)
	}
	fail := func(x int) Maybe[int] {
		return nothingOf[int]()
	}

	fmt.Println("Left Identity Law (Just):", checkLeftIdentityLaw(10, increment))
	fmt.Println("Left Identity Law (Nothing):", checkLeftIdentityLaw(10, fail))

	fmt.Println("Right Identity Law (Just):", checkRightIdentityLaw(justOf(10)))
	fmt.Println("Right Identity Law (Nothing):", checkRightIdentityLaw(nothingOf[int]()))

	fmt.Println("Associativity Law (Just):", checkAssociativeLaw(justOf(10), increment, double))
	fmt.Println("Associativity Law (Nothing):", checkAssociativeLaw(nothingOf[int](), increment, double))
}
