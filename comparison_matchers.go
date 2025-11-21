package gotest

import (
	"fmt"
)

// Matches values that are greater than the threshold.
//
// Works with any ordered types including:
//   - All numeric types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64)
//   - strings (lexicographic comparison)
//
// Examples:
//
//	ExpectThat(t, 5, Gt(3))
//	ExpectThat(t, 10.5, Gt(10.0))
//	ExpectThat(t, "banana", Gt("apple"))
func Gt[T ordered](threshold T) Matcher {
	return gtMatcher[T]{threshold}
}

// Matches values that are less than the threshold.
//
// Works with any ordered types including:
//   - All numeric types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64)
//   - strings (lexicographic comparison)
//
// Examples:
//
//	ExpectThat(t, 3, Lt(5))
//	ExpectThat(t, 10.0, Lt(10.5))
//	ExpectThat(t, "apple", Lt("banana"))
func Lt[T ordered](threshold T) Matcher {
	return ltMatcher[T]{threshold}
}

// Matches values that are greater than or equal to the threshold.
//
// Works with any ordered types including:
//   - All numeric types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64)
//   - strings (lexicographic comparison)
//
// Examples:
//
//	ExpectThat(t, 5, Gte(5))
//	ExpectThat(t, 10.5, Gte(10.0))
func Gte[T ordered](threshold T) Matcher {
	return gteMatcher[T]{threshold}
}

// Matches values that are less than or equal to the threshold.
//
// Works with any ordered types including:
//   - All numeric types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64)
//   - strings (lexicographic comparison)
//
// Examples:
//
//	ExpectThat(t, 5, Lte(5))
//	ExpectThat(t, 10.0, Lte(10.5))
func Lte[T ordered](threshold T) Matcher {
	return lteMatcher[T]{threshold}
}

// ordered is a constraint that permits any ordered type: any type that supports
// the operators < <= >= >.
type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

type gtMatcher[T ordered] struct {
	threshold T
}

func (g gtMatcher[T]) String() string {
	return fmt.Sprintf("is greater than %v (%T)", g.threshold, g.threshold)
}

func (g gtMatcher[T]) Matches(x any) bool {
	val, ok := x.(T)
	if !ok {
		return false
	}
	return val > g.threshold
}

func (g gtMatcher[T]) ExplainFailure(x any) (string, bool) {
	val, ok := x.(T)
	if !ok {
		return fmt.Sprintf("type %T, which is not comparable to %v (%T)", x, g.threshold, g.threshold), true
	}
	return fmt.Sprintf("is %v, which is not greater than %v", val, g.threshold), true
}

type ltMatcher[T ordered] struct {
	threshold T
}

func (l ltMatcher[T]) String() string {
	return fmt.Sprintf("is less than %v (%T)", l.threshold, l.threshold)
}

func (l ltMatcher[T]) Matches(x any) bool {
	val, ok := x.(T)
	if !ok {
		return false
	}
	return val < l.threshold
}

func (l ltMatcher[T]) ExplainFailure(x any) (string, bool) {
	val, ok := x.(T)
	if !ok {
		return fmt.Sprintf("type %T, which is not comparable to %v (%T)", x, l.threshold, l.threshold), true
	}
	return fmt.Sprintf("is %v, which is not less than %v", val, l.threshold), true
}

type gteMatcher[T ordered] struct {
	threshold T
}

func (g gteMatcher[T]) String() string {
	return fmt.Sprintf("is greater than or equal to %v (%T)", g.threshold, g.threshold)
}

func (g gteMatcher[T]) Matches(x any) bool {
	val, ok := x.(T)
	if !ok {
		return false
	}
	return val >= g.threshold
}

func (g gteMatcher[T]) ExplainFailure(x any) (string, bool) {
	val, ok := x.(T)
	if !ok {
		return fmt.Sprintf("type %T, which is not comparable to %v (%T)", x, g.threshold, g.threshold), true
	}
	return fmt.Sprintf("is %v, which is not greater than or equal to %v", val, g.threshold), true
}

type lteMatcher[T ordered] struct {
	threshold T
}

func (l lteMatcher[T]) String() string {
	return fmt.Sprintf("is less than or equal to %v (%T)", l.threshold, l.threshold)
}

func (l lteMatcher[T]) Matches(x any) bool {
	val, ok := x.(T)
	if !ok {
		return false
	}
	return val <= l.threshold
}

func (l lteMatcher[T]) ExplainFailure(x any) (string, bool) {
	val, ok := x.(T)
	if !ok {
		return fmt.Sprintf("type %T, which is not comparable to %v (%T)", x, l.threshold, l.threshold), true
	}
	return fmt.Sprintf("is %v, which is not less than or equal to %v", val, l.threshold), true
}
