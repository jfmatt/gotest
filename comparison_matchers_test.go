package gotest

import (
	"testing"
)

type Username string

func TestGt(t *testing.T) {
	// Same-type integers
	ExpectThat(t, 5, Gt(3))
	ExpectThat(t, 5, Not(Gt(5)))
	ExpectThat(t, 5, Not(Gt(10)))

	// Same-type floats
	ExpectThat(t, 10.5, Gt(10.0))
	ExpectThat(t, 10.5, Not(Gt(10.5)))
	ExpectThat(t, 10.5, Not(Gt(11.0)))

	// Mixed int sizes (compare as int64)
	ExpectThat(t, int32(100), Gt(int8(50)))
	ExpectThat(t, int64(100), Gt(int(50)))

	// Mixed uint sizes (compare as uint64)
	ExpectThat(t, uint32(100), Gt(uint8(50)))
	ExpectThat(t, uint64(100), Gt(uint(50)))

	// Float vs int (compare as float64)
	ExpectThat(t, 10.5, Gt(10))
	ExpectThat(t, 10.1, Gt(10))
	ExpectThat(t, 10.0, Not(Gt(10)))

	// Int vs float (compare as float64)
	ExpectThat(t, 11, Gt(10.5))
	ExpectThat(t, 10, Not(Gt(10.5)))

	// Uint vs int (compare as uint64)
	ExpectThat(t, uint(10), Gt(5))
	ExpectThat(t, uint8(10), Gt(int32(5)))

	// Uint vs float - incompatible, should not match
	ExpectThat(t, uint(10), Not(Gt(5.0)))
	ExpectThat(t, uint8(10), Not(Gt(float64(5))))

	// Strings
	ExpectThat(t, "banana", Gt("apple"))
	ExpectThat(t, "banana", Not(Gt("banana")))
	ExpectThat(t, "banana", Not(Gt("cherry")))

	// Custom string types (compare as string)
	ExpectThat(t, Username("bob"), Gt("alice"))
	ExpectThat(t, Username("bob"), Not(Gt("charlie")))

	// Incompatible types
	ExpectThat(t, "5", Not(Gt(3)))
	ExpectThat(t, 5, Not(Gt("3")))
}

func TestLt(t *testing.T) {
	// Same-type integers
	ExpectThat(t, 3, Lt(5))
	ExpectThat(t, 5, Not(Lt(5)))
	ExpectThat(t, 10, Not(Lt(5)))

	// Same-type floats
	ExpectThat(t, 10.0, Lt(10.5))
	ExpectThat(t, 10.5, Not(Lt(10.5)))
	ExpectThat(t, 11.0, Not(Lt(10.5)))

	// Mixed int sizes (compare as int64)
	ExpectThat(t, int8(50), Lt(int32(100)))
	ExpectThat(t, int(50), Lt(int64(100)))

	// Mixed uint sizes (compare as uint64)
	ExpectThat(t, uint8(50), Lt(uint32(100)))
	ExpectThat(t, uint(50), Lt(uint64(100)))

	// Float vs int (compare as float64)
	ExpectThat(t, 10.0, Lt(11))
	ExpectThat(t, 9.9, Lt(10))
	ExpectThat(t, 10.5, Not(Lt(10)))

	// Int vs float (compare as float64)
	ExpectThat(t, 10, Lt(10.5))
	ExpectThat(t, 11, Not(Lt(10.5)))

	// Uint vs int (compare as uint64)
	ExpectThat(t, uint(5), Lt(10))
	ExpectThat(t, uint8(5), Lt(int32(10)))

	// Uint vs float - incompatible, should not match
	ExpectThat(t, uint(5), Not(Lt(10.0)))
	ExpectThat(t, uint8(5), Not(Lt(float64(10))))

	// Strings
	ExpectThat(t, "apple", Lt("banana"))
	ExpectThat(t, "banana", Not(Lt("banana")))
	ExpectThat(t, "cherry", Not(Lt("banana")))

	// Custom string types (compare as string)
	ExpectThat(t, Username("alice"), Lt("bob"))
	ExpectThat(t, Username("charlie"), Not(Lt("bob")))

	// Incompatible types
	ExpectThat(t, "3", Not(Lt(5)))
	ExpectThat(t, 3, Not(Lt("5")))
}

func TestGte(t *testing.T) {
	// Same-type integers
	ExpectThat(t, 5, Ge(3))
	ExpectThat(t, 5, Ge(5))
	ExpectThat(t, 5, Not(Ge(10)))

	// Same-type floats
	ExpectThat(t, 10.5, Ge(10.0))
	ExpectThat(t, 10.5, Ge(10.5))
	ExpectThat(t, 10.5, Not(Ge(11.0)))

	// Mixed int sizes (compare as int64)
	ExpectThat(t, int32(100), Ge(int8(50)))
	ExpectThat(t, int32(100), Ge(int8(100)))

	// Mixed uint sizes (compare as uint64)
	ExpectThat(t, uint32(100), Ge(uint8(50)))
	ExpectThat(t, uint32(100), Ge(uint8(100)))

	// Float vs int (compare as float64)
	ExpectThat(t, 10.5, Ge(10))
	ExpectThat(t, 10.0, Ge(10))
	ExpectThat(t, 9.9, Not(Ge(10)))

	// Int vs float (compare as float64)
	ExpectThat(t, 11, Ge(10.5))
	ExpectThat(t, 10, Ge(10.0))
	ExpectThat(t, 10, Not(Ge(10.5)))

	// Uint vs int (compare as uint64)
	ExpectThat(t, uint(10), Ge(5))
	ExpectThat(t, uint(10), Ge(10))

	// Uint vs float - incompatible
	ExpectThat(t, uint(10), Not(Ge(5.0)))

	// Strings
	ExpectThat(t, "banana", Ge("apple"))
	ExpectThat(t, "banana", Ge("banana"))
	ExpectThat(t, "banana", Not(Ge("cherry")))

	// Custom string types
	ExpectThat(t, Username("bob"), Ge("alice"))
	ExpectThat(t, Username("bob"), Ge("bob"))
}

func TestLte(t *testing.T) {
	// Same-type integers
	ExpectThat(t, 3, Le(5))
	ExpectThat(t, 5, Le(5))
	ExpectThat(t, 10, Not(Le(5)))

	// Same-type floats
	ExpectThat(t, 10.0, Le(10.5))
	ExpectThat(t, 10.5, Le(10.5))
	ExpectThat(t, 11.0, Not(Le(10.5)))

	// Mixed int sizes (compare as int64)
	ExpectThat(t, int8(50), Le(int32(100)))
	ExpectThat(t, int8(100), Le(int32(100)))

	// Mixed uint sizes (compare as uint64)
	ExpectThat(t, uint8(50), Le(uint32(100)))
	ExpectThat(t, uint8(100), Le(uint32(100)))

	// Float vs int (compare as float64)
	ExpectThat(t, 10.0, Le(11))
	ExpectThat(t, 10.0, Le(10))
	ExpectThat(t, 10.5, Not(Le(10)))

	// Int vs float (compare as float64)
	ExpectThat(t, 10, Le(10.5))
	ExpectThat(t, 10, Le(10.0))
	ExpectThat(t, 11, Not(Le(10.5)))

	// Uint vs int (compare as uint64)
	ExpectThat(t, uint(5), Le(10))
	ExpectThat(t, uint(5), Le(5))

	// Uint vs float - incompatible
	ExpectThat(t, uint(5), Not(Le(10.0)))

	// Strings
	ExpectThat(t, "apple", Le("banana"))
	ExpectThat(t, "banana", Le("banana"))
	ExpectThat(t, "cherry", Not(Le("banana")))

	// Custom string types
	ExpectThat(t, Username("alice"), Le("bob"))
	ExpectThat(t, Username("bob"), Le("bob"))
}
