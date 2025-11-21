package gotest

import (
	"cmp"
	"fmt"
	"reflect"
	"strings"
)

// Matches values that are greater than the threshold.
//
// Works with any ordered types including:
//   - All numeric types
//   - strings (lexicographic comparison)
//
// Examples:
//
//	ExpectThat(t, 5, Gt(3))
//	ExpectThat(t, 10.5, Gt(10.0))
//	ExpectThat(t, "banana", Gt("apple"))
func Gt[T cmp.Ordered](threshold T) Matcher {
	return gtMatcher[T]{threshold}
}

// Matches values that are less than the threshold.
//
// Works with any ordered types including:
//   - All numeric types
//   - strings (lexicographic comparison)
//
// Examples:
//
//	ExpectThat(t, 3, Lt(5))
//	ExpectThat(t, 10.0, Lt(10.5))
//	ExpectThat(t, "apple", Lt("banana"))
func Lt[T cmp.Ordered](threshold T) Matcher {
	return ltMatcher[T]{threshold}
}

// Matches values that are greater than or equal to the threshold.
//
// Works with any ordered types including:
//   - All numeric types
//   - strings (lexicographic comparison)
//
// Examples:
//
//	ExpectThat(t, 5, Ge(5))
//	ExpectThat(t, 10.5, Ge(10.0))
func Ge[T cmp.Ordered](threshold T) Matcher {
	return geMatcher[T]{threshold}
}

// Matches values that are less than or equal to the threshold.
//
// Works with any ordered types including:
//   - All numeric types
//   - strings (lexicographic comparison)
//
// Examples:
//
//	ExpectThat(t, 5, Le(5))
//	ExpectThat(t, 10.0, Le(10.5))
func Le[T cmp.Ordered](threshold T) Matcher {
	return leMatcher[T]{threshold}
}

type gtMatcher[T cmp.Ordered] struct {
	threshold T
}

func (g gtMatcher[T]) String() string {
	return fmt.Sprintf("is greater than %v (%T)", g.threshold, g.threshold)
}

func (g gtMatcher[T]) Matches(x any) bool {
	canCompare, cmpResult := tryCompare(x, g.threshold)
	if !canCompare {
		return false
	}
	return cmpResult > 0 // x > threshold
}

func (g gtMatcher[T]) ExplainFailure(x any) (string, bool) {
	val, ok := x.(T)
	if !ok {
		return fmt.Sprintf("type %T, which is not comparable to %v (%T)", x, g.threshold, g.threshold), true
	}
	return fmt.Sprintf("is %v, which is not greater than %v", val, g.threshold), true
}

type ltMatcher[T cmp.Ordered] struct {
	threshold T
}

func (l ltMatcher[T]) String() string {
	return fmt.Sprintf("is less than %v (%T)", l.threshold, l.threshold)
}

func (l ltMatcher[T]) Matches(x any) bool {
	canCompare, cmpResult := tryCompare(x, l.threshold)
	if !canCompare {
		return false
	}
	return cmpResult < 0 // x < threshold
}

func (l ltMatcher[T]) ExplainFailure(x any) (string, bool) {
	val, ok := x.(T)
	if !ok {
		return fmt.Sprintf("type %T, which is not comparable to %v (%T)", x, l.threshold, l.threshold), true
	}
	return fmt.Sprintf("is %v, which is not less than %v", val, l.threshold), true
}

type geMatcher[T cmp.Ordered] struct {
	threshold T
}

func (g geMatcher[T]) String() string {
	return fmt.Sprintf("is greater than or equal to %v (%T)", g.threshold, g.threshold)
}

func (g geMatcher[T]) Matches(x any) bool {
	canCompare, cmpResult := tryCompare(x, g.threshold)
	if !canCompare {
		return false
	}
	return cmpResult >= 0 // x >= threshold
}

func (g geMatcher[T]) ExplainFailure(x any) (string, bool) {
	val, ok := x.(T)
	if !ok {
		return fmt.Sprintf("type %T, which is not comparable to %v (%T)", x, g.threshold, g.threshold), true
	}
	return fmt.Sprintf("is %v, which is not greater than or equal to %v", val, g.threshold), true
}

type leMatcher[T cmp.Ordered] struct {
	threshold T
}

func (l leMatcher[T]) String() string {
	return fmt.Sprintf("is less than or equal to %v (%T)", l.threshold, l.threshold)
}

func (l leMatcher[T]) Matches(x any) bool {
	canCompare, cmpResult := tryCompare(x, l.threshold)
	if !canCompare {
		return false
	}
	return cmpResult <= 0 // x <= threshold
}

func (l leMatcher[T]) ExplainFailure(x any) (string, bool) {
	val, ok := x.(T)
	if !ok {
		return fmt.Sprintf("type %T, which is not comparable to %v (%T)", x, l.threshold, l.threshold), true
	}
	return fmt.Sprintf("is %v, which is not less than or equal to %v", val, l.threshold), true
}

// tryCompare attempts to compare two values using smart type promotion.
// Returns (canCompare bool, comparisonResult int) where:
//   - canCompare is false if types are incompatible
//   - comparisonResult is -1 if actual < threshold, 0 if equal, 1 if actual > threshold
func tryCompare[T cmp.Ordered](actual any, threshold T) (bool, int) {
	actualVal := reflect.ValueOf(actual)
	thresholdVal := reflect.ValueOf(threshold)

	// Check what types we're dealing with
	hasUint := actualVal.CanUint() || thresholdVal.CanUint()
	hasFloat := actualVal.CanFloat() || thresholdVal.CanFloat()

	// uint + float = incompatible
	if hasUint && hasFloat {
		return false, 0
	}

	// If either is uint (and no floats), compare as uint64
	if hasUint {
		a := toUint64(actualVal)
		b := toUint64(thresholdVal)
		if a < b {
			return true, -1
		} else if a > b {
			return true, 1
		}
		return true, 0
	}

	// If either is float (and no uints), compare as float64
	if hasFloat {
		a := toFloat64(actualVal)
		b := toFloat64(thresholdVal)
		if a < b {
			return true, -1
		} else if a > b {
			return true, 1
		}
		return true, 0
	}

	// If both are signed ints, compare as int64
	if actualVal.CanInt() && thresholdVal.CanInt() {
		a := actualVal.Int()
		b := thresholdVal.Int()
		if a < b {
			return true, -1
		} else if a > b {
			return true, 1
		}
		return true, 0
	}

	// If both are strings, compare as string
	if actualVal.Kind() == reflect.String && thresholdVal.Kind() == reflect.String {
		return true, strings.Compare(actualVal.String(), thresholdVal.String())
	}

	// Incompatible types
	return false, 0
}

// Conversion helpers
func toUint64(v reflect.Value) uint64 {
	switch {
	case v.CanUint():
		return v.Uint()
	case v.CanInt():
		return uint64(v.Int())
	default:
		return 0
	}
}

func toFloat64(v reflect.Value) float64 {
	switch {
	case v.CanFloat():
		return v.Float()
	case v.CanInt():
		return float64(v.Int())
	default:
		return 0
	}
}
