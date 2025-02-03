package gotest

import (
	"go.uber.org/mock/gomock"
)

var (
	// Re-exporting some gomock matchers that already do everything we want
	// them to. This allows users to always use `import _ gotest` instead of
	// remembering which matchers are here vs. in gomock.
	AssignableToTypeOf = gomock.AssignableToTypeOf
	Nil                = gomock.Nil
)

func AsMatcher(x any) gomock.Matcher {
	if alreadyMatcher, ok := x.(gomock.Matcher); ok {
		return alreadyMatcher
	} else {
		return Eq(x)
	}
}

// Negates the inner condition. If `x` is a matcher, then Not(x) will match
// conditions where x doesn't match. If `x` is a value, then Not(x) will match
// conditions where the value is equal.
//
// Examples:
//
//	ExpectThat(t, 4, Not(Eq(5)))
//	ExpectThat(t, 4, Not(5))
//	ExpectThat(t, 4, Not(Gt(5)))
//
// This is exactly the same as gomock.Not, except that it uses our
// implementation of Eq() when `x` is a value.
func Not(x any) gomock.Matcher {
	return gomock.Not(AsMatcher(x))
}
