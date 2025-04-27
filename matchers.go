package gotest

import (
	"go.uber.org/mock/gomock"
)

var (
	Any = gomock.Any
)

// In order to ensure compatibility with generated mocks, the
// Matcher type in this package is the same as the one from
// gomock. It's re-exported here for convenience, and to
// improve the formatting on pkg.go.dev.
type Matcher = gomock.Matcher

func AsMatcher(x any) Matcher {
	if alreadyMatcher, ok := x.(Matcher); ok {
		return alreadyMatcher
	} else {
		return Eq(x)
	}
}

// AssignableToTypeOf is a Matcher that matches if the parameter to the mock
// function is assignable to the type of the parameter to this function.
//
// Example usage:
//
//	var s fmt.Stringer = &bytes.Buffer{}
//	AssignableToTypeOf(s).Matches(time.Second) // returns true
//	AssignableToTypeOf(s).Matches(99) // returns false
//
//	var ctx = reflect.TypeOf((*context.Context)(nil)).Elem()
//	AssignableToTypeOf(ctx).Matches(context.Background()) // returns true
//
// (This matcher is the same as gomock.AssignableToTypeOf. It's
// re-exported here for convenience with `import _ ` so that users
// don't need to remember which package particular matchers come
// from.)
func AssignableToTypeOf(x any) Matcher {
	return gomock.AssignableToTypeOf(x)
}

// Nil returns a matcher that matches if the received value is nil.
//
// Example usage:
//
//	var x *bytes.Buffer
//	Nil().Matches(x) // returns true
//	x = &bytes.Buffer{}
//	Nil().Matches(x) // returns false
//
// (This matcher is the same as gomock.Nil. It's re-exported here
// for convenience with `import _ ` so that users don't need to
// remember which package particular matchers come from.)
func Nil() Matcher {
	return gomock.Nil()
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
func Not(x any) Matcher {
	return gomock.Not(AsMatcher(x))
}
