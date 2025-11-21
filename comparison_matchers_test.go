package gotest

import (
	"testing"
)

func TestGt(t *testing.T) {
	// Integers
	ExpectThat(t, 5, Gt(3))
	ExpectThat(t, 5, Not(Gt(5)))
	ExpectThat(t, 5, Not(Gt(10)))

	// Floats
	ExpectThat(t, 10.5, Gt(10.0))
	ExpectThat(t, 10.5, Not(Gt(10.5)))
	ExpectThat(t, 10.5, Not(Gt(11.0)))

	// Strings
	ExpectThat(t, "banana", Gt("apple"))
	ExpectThat(t, "banana", Not(Gt("banana")))
	ExpectThat(t, "banana", Not(Gt("cherry")))

	// Type mismatch
	ExpectThat(t, 5, Not(Gt(3.0)))
	ExpectThat(t, "5", Not(Gt(3)))
}

func TestLt(t *testing.T) {
	// Integers
	ExpectThat(t, 3, Lt(5))
	ExpectThat(t, 5, Not(Lt(5)))
	ExpectThat(t, 10, Not(Lt(5)))

	// Floats
	ExpectThat(t, 10.0, Lt(10.5))
	ExpectThat(t, 10.5, Not(Lt(10.5)))
	ExpectThat(t, 11.0, Not(Lt(10.5)))

	// Strings
	ExpectThat(t, "apple", Lt("banana"))
	ExpectThat(t, "banana", Not(Lt("banana")))
	ExpectThat(t, "cherry", Not(Lt("banana")))

	// Type mismatch
	ExpectThat(t, 3.0, Not(Lt(5)))
	ExpectThat(t, 3, Not(Lt("5")))
}

func TestGte(t *testing.T) {
	// Integers
	ExpectThat(t, 5, Gte(3))
	ExpectThat(t, 5, Gte(5))
	ExpectThat(t, 5, Not(Gte(10)))

	// Floats
	ExpectThat(t, 10.5, Gte(10.0))
	ExpectThat(t, 10.5, Gte(10.5))
	ExpectThat(t, 10.5, Not(Gte(11.0)))

	// Strings
	ExpectThat(t, "banana", Gte("apple"))
	ExpectThat(t, "banana", Gte("banana"))
	ExpectThat(t, "banana", Not(Gte("cherry")))
}

func TestLte(t *testing.T) {
	// Integers
	ExpectThat(t, 3, Lte(5))
	ExpectThat(t, 5, Lte(5))
	ExpectThat(t, 10, Not(Lte(5)))

	// Floats
	ExpectThat(t, 10.0, Lte(10.5))
	ExpectThat(t, 10.5, Lte(10.5))
	ExpectThat(t, 11.0, Not(Lte(10.5)))

	// Strings
	ExpectThat(t, "apple", Lte("banana"))
	ExpectThat(t, "banana", Lte("banana"))
	ExpectThat(t, "cherry", Not(Lte("banana")))
}
