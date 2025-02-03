package gotest

import (
	"fmt"

	"go.uber.org/mock/gomock"
)

// Tests that `val` fulfills `expected`. If not, causes the test (`t`) to fail.
//
// Returns the result of the check - true on success, false on failure. This
// can be useful to avoid running redundant checks in failure scenarios.
//
// Normally, `expected` should be a gomock.Matcher. As a convenience, if
// `expected` is not a matcher, then the expectation will be Eq(expected).
//
// Examples:
//
//	ExpectThat(t, "ab", "ab")           // succeeds
//	ExpectThat(t, "ab", "a")            // fails
//	ExpectThat(t, "ab", HasSubstr("a")) // succeeds
//	if !ExpectThat(t, someList, Not(Empty())) {
//		// If the list is non-empty, run extra checks on the contents
//		// ...else, the test fails anyway
//	}
func ExpectThat(t gomock.TestHelper, val any, expected any) bool {
	t.Helper()

	matcher := AsMatcher(expected)

	ok := matcher.Matches(val)
	if ok {
		return true
	}

	t.Errorf(getExplanation("Expectation", matcher, val))
	return false
}

// Same as ExpectThat, but more explicitly tests values of the same type for
// equality.
func ExpectEq[T any](t gomock.TestHelper, actual T, expected T) bool {
	t.Helper()
	return ExpectThat(t, actual, Eq(expected))
}

// Tests that `f()` causes a fatal error that fulfills `errMatcher`.
//
// If the function does not panic, or if it panics with an error that doesn't
// match `errMatcher`, then the test represented by `t` will fail.
//
// Example:
//
//	ExpectFatal(t, ErrorMessage(HasSubstr("bad thing")), func() {
//	  panic("a bad thing happened")
//	})  // succeeds
//
//	ExpectFatal(t, Any(), func() {
//	  fmt.Println("a-ok!")
//	})  // fails, because the function didn't panic
func ExpectFatal(t gomock.TestHelper, errMatcher gomock.Matcher, f func()) (success bool) {
	t.Helper()

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected fatal error, but none occurred...")
			success = false
		} else {
			success = ExpectThat(t, r, errMatcher)
		}
	}()
	f()
	return true
}

// Same as ExpectThat, but causes the test to immediately terminate on failure.
//
// Useful for checking preconditions that would cause fatal errors in further
// test cases if violated.
//
// Example:
//
//	AssertThat(t, someList, Not(Empty()))
//	ExpectThat(t, someList[0], HasSubstr("val1"))
//
// In this scenario, the first test acts as a guard; if it fails, the second
// test, and any further code, shouldn't be run.
func AssertThat(t gomock.TestHelper, val any, expected any) {
	t.Helper()

	matcher := AsMatcher(expected)

	ok := matcher.Matches(val)
	if ok {
		return
	}

	t.Fatalf(getExplanation("Assertion", matcher, val))
}

// Same as ExpectEq(), but causes the test to immediately terminate on failure.
func AssertEq[T any](t gomock.TestHelper, actual T, expected T) {
	t.Helper()
	AssertThat(t, actual, Eq(expected))
}

// Same as ExpectFatal(), but causes the test to immediately terminate on failure.
func AssertFatal(t gomock.TestHelper, errMatcher gomock.Matcher, f func()) {
	t.Helper()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("Asserted fatal error, but none occurred...")
		} else {
			AssertThat(t, r, errMatcher)
		}
	}()
	f()
}

func getExplanation(context string, matcher gomock.Matcher, val any) string {
	var e string
	var useE bool
	if explainer, ok := matcher.(MismatchExplainer); ok {
		e, useE = explainer.ExplainFailure(val)
	} else {
		e, useE = "", false
	}

	if useE {
		return fmt.Sprintf("%s failed:\n  Wanted: %s\n  Got: %s\n  ...where %s",
			context, matcher.String(), formatGot(val, matcher), e)
	} else {
		return fmt.Sprintf("%s failed:\n  Wanted: %s\n  Got: %s",
			context, matcher.String(), formatGot(val, matcher))
	}
}
