package gotest

import (
	"cmp"
	"fmt"
	"math"
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

type numClass int

const (
	numClassNonNumeric numClass = iota
	numClassNormalInt
	numClassNegativeInt
	numClassUint
	numClassBigUint // bigger than MaxInt64
	numClassFloat
	numClassNegativeFloat
	numClassBigFloat
)

// tryCompare attempts to compare two values using smart type promotion.
// Returns (canCompare bool, comparisonResult int) where:
//   - canCompare is false if types are incompatible
//   - comparisonResult is -1 if actual < threshold, 0 if equal, 1 if actual > threshold
func tryCompare[T cmp.Ordered](actual any, threshold T) (bool, int) {
	actualVal := reflect.ValueOf(actual)
	thresholdVal := reflect.ValueOf(threshold)

	// If both are strings, compare as string
	if actualVal.Kind() == reflect.String && thresholdVal.Kind() == reflect.String {
		return true, strings.Compare(actualVal.String(), thresholdVal.String())
	}

	actualClass := classify(actual)
	thresholdClass := classify(threshold)
	if actualClass == numClassNonNumeric || thresholdClass == numClassNonNumeric {
		return false, 0
	}

	// Handle comparisons based on classification
	switch {
	case actualClass == numClassBigUint && thresholdClass == numClassBigUint:
		if actualVal.Uint() < thresholdVal.Uint() {
			return true, -1
		}
		if actualVal.Uint() > thresholdVal.Uint() {
			return true, 1
		}
		return true, 0

	// Negative vs positive: negative is always less
	case (actualClass == numClassNegativeInt || actualClass == numClassNegativeFloat) &&
		(thresholdClass == numClassNormalInt || thresholdClass == numClassUint ||
			thresholdClass == numClassBigUint || thresholdClass == numClassFloat ||
			thresholdClass == numClassBigFloat):
		return true, -1

	case (thresholdClass == numClassNegativeInt || thresholdClass == numClassNegativeFloat) &&
		(actualClass == numClassNormalInt || actualClass == numClassUint ||
			actualClass == numClassBigUint || actualClass == numClassFloat ||
			actualClass == numClassBigFloat):
		return true, 1

	// BigUint comparisons with other types
	case actualClass == numClassBigUint && (thresholdClass == numClassNormalInt ||
		thresholdClass == numClassUint || thresholdClass == numClassFloat):
		return true, 1 // BigUint is always > smaller positive types

	case thresholdClass == numClassBigUint && (actualClass == numClassNormalInt ||
		actualClass == numClassUint || actualClass == numClassFloat):
		return true, -1 // Smaller positive types are always < BigUint

	case actualClass == numClassBigUint && thresholdClass == numClassBigFloat:
		// Compare as float64
		actualFloat := float64(actualVal.Uint())
		thresholdFloat := thresholdVal.Float()
		if actualFloat < thresholdFloat {
			return true, -1
		}
		if actualFloat > thresholdFloat {
			return true, 1
		}
		return true, 0

	case thresholdClass == numClassBigUint && actualClass == numClassBigFloat:
		actualFloat := actualVal.Float()
		thresholdFloat := float64(thresholdVal.Uint())
		if actualFloat < thresholdFloat {
			return true, -1
		}
		if actualFloat > thresholdFloat {
			return true, 1
		}
		return true, 0

	// BigFloat comparisons
	case actualClass == numClassBigFloat && (thresholdClass == numClassNormalInt ||
		thresholdClass == numClassUint || thresholdClass == numClassFloat):
		actualFloat := actualVal.Float()
		thresholdFloat := toFloat64(thresholdVal)
		if actualFloat < thresholdFloat {
			return true, -1
		}
		if actualFloat > thresholdFloat {
			return true, 1
		}
		return true, 0

	case thresholdClass == numClassBigFloat && (actualClass == numClassNormalInt ||
		actualClass == numClassUint || actualClass == numClassFloat):
		actualFloat := toFloat64(actualVal)
		thresholdFloat := thresholdVal.Float()
		if actualFloat < thresholdFloat {
			return true, -1
		}
		if actualFloat > thresholdFloat {
			return true, 1
		}
		return true, 0

	case actualClass == numClassBigFloat && thresholdClass == numClassBigFloat:
		if actualVal.Float() < thresholdVal.Float() {
			return true, -1
		}
		if actualVal.Float() > thresholdVal.Float() {
			return true, 1
		}
		return true, 0

	case actualClass == numClassBigFloat && (thresholdClass == numClassNegativeInt ||
		thresholdClass == numClassNegativeFloat):
		return true, 1 // BigFloat is always > negative numbers

	case thresholdClass == numClassBigFloat && (actualClass == numClassNegativeInt ||
		actualClass == numClassNegativeFloat):
		return true, -1 // Negative numbers are always < BigFloat

	// Float comparisons with smaller types
	case (actualClass == numClassFloat || actualClass == numClassNegativeFloat) &&
		(thresholdClass == numClassNormalInt || thresholdClass == numClassUint ||
			thresholdClass == numClassFloat || thresholdClass == numClassNegativeInt ||
			thresholdClass == numClassNegativeFloat):
		actualFloat := toFloat64(actualVal)
		thresholdFloat := toFloat64(thresholdVal)
		if actualFloat < thresholdFloat {
			return true, -1
		}
		if actualFloat > thresholdFloat {
			return true, 1
		}
		return true, 0

	case (thresholdClass == numClassFloat || thresholdClass == numClassNegativeFloat) &&
		(actualClass == numClassNormalInt || actualClass == numClassUint ||
			actualClass == numClassNegativeInt):
		actualFloat := toFloat64(actualVal)
		thresholdFloat := toFloat64(thresholdVal)
		if actualFloat < thresholdFloat {
			return true, -1
		}
		if actualFloat > thresholdFloat {
			return true, 1
		}
		return true, 0

	// Integer-only comparisons (both non-negative)
	case (actualClass == numClassNormalInt || actualClass == numClassUint) &&
		(thresholdClass == numClassNormalInt || thresholdClass == numClassUint):
		actualUint := toUint64(actualVal)
		thresholdUint := toUint64(thresholdVal)
		if actualUint < thresholdUint {
			return true, -1
		}
		if actualUint > thresholdUint {
			return true, 1
		}
		return true, 0

	// Both negative integers
	case actualClass == numClassNegativeInt && thresholdClass == numClassNegativeInt:
		if actualVal.Int() < thresholdVal.Int() {
			return true, -1
		}
		if actualVal.Int() > thresholdVal.Int() {
			return true, 1
		}
		return true, 0

	// Negative int vs negative float
	case actualClass == numClassNegativeInt && thresholdClass == numClassNegativeFloat:
		actualFloat := float64(actualVal.Int())
		thresholdFloat := thresholdVal.Float()
		if actualFloat < thresholdFloat {
			return true, -1
		}
		if actualFloat > thresholdFloat {
			return true, 1
		}
		return true, 0

	case actualClass == numClassNegativeFloat && thresholdClass == numClassNegativeInt:
		actualFloat := actualVal.Float()
		thresholdFloat := float64(thresholdVal.Int())
		if actualFloat < thresholdFloat {
			return true, -1
		}
		if actualFloat > thresholdFloat {
			return true, 1
		}
		return true, 0

	// Default: types cannot be compared
	default:
		return false, 0
	}
}

func classify(v any) numClass {
	val := reflect.ValueOf(v)
	switch {
	case val.CanInt() && val.Int() < 0:
		return numClassNegativeInt
	case val.CanInt():
		return numClassNormalInt
	case val.CanUint() && val.Uint() > math.MaxInt64:
		return numClassBigUint
	case val.CanUint():
		return numClassUint
	case val.CanFloat() && val.Float() < 0:
		return numClassNegativeFloat
	case val.CanFloat() && val.Float() > math.MaxInt64:
		return numClassBigFloat
	case val.CanFloat():
		return numClassFloat
	default:
		return numClassNonNumeric
	}
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
	case v.CanUint():
		return float64(v.Uint())
	default:
		return 0
	}
}
