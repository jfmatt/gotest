package gotest

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
)

// ------------------------------------------
//
// A fake of testing.T. Used to observe whether our assertions would cause test
// failures, without actually causing test failures.

var _ gomock.TestHelper = (*testReporter)(nil)

type testReporter struct {
	nonFatals []string
	fatals    []string
}

func (c *testReporter) Reset() {
	c.nonFatals = nil
	c.fatals = nil
}

func (c *testReporter) HasFailures() bool {
	return len(c.nonFatals) > 0 || len(c.fatals) > 0
}

func (c *testReporter) HasFatals() bool {
	return len(c.fatals) > 0
}

func (c *testReporter) Errorf(format string, args ...any) {
	c.nonFatals = append(c.nonFatals, fmt.Sprintf(format, args...))
}

func (c *testReporter) Fatalf(format string, args ...any) {
	c.fatals = append(c.fatals, fmt.Sprintf(format, args...))
}

func (c *testReporter) Helper() {}

// ------------------------------------------

// Tests that ExpectThat() and friends work to fail tests. These assertions are
// used in testing all matchers in this library, so we'll write this test
// without using any of them.  Instead, we'll use the matchers from the
// standard gomock library in order to bootstrap.
func TestExpectations(t *testing.T) {
	expectNonFatal := func(t testing.TB, r *testReporter, substr string) {
		t.Helper()
		if len(r.nonFatals) != 1 {
			t.Errorf("Expected 1 non-fatal error, got: %v", r.nonFatals)
			return
		}
		if len(r.fatals) != 0 {
			t.Errorf("got unexpected fatal error: %v", r.fatals)
			return
		}
		if !strings.Contains(r.nonFatals[0], substr) {
			t.Errorf("got the wrong error; wanted %s, got %s", substr, r.nonFatals[0])
		}
	}
	expectFatal := func(t testing.TB, r *testReporter, substr string) {
		t.Helper()
		if len(r.fatals) != 1 {
			t.Errorf("Expected 1 fatal error, got: %v", r.fatals)
			return
		}
		if len(r.nonFatals) != 0 {
			t.Errorf("got unexpected non-fatal error: %v", r.nonFatals)
			return
		}
		if !strings.Contains(r.fatals[0], substr) {
			t.Errorf("got the wrong error; wanted %s, got %s", substr, r.fatals[0])
		}
	}
	t.Run("BasicExpect", func(t *testing.T) {
		r := testReporter{}

		// Should do nothing.
		r.Reset()
		ExpectThat(&r, "hello, world", gomock.Eq("hello, world"))
		if r.HasFailures() {
			t.Errorf("got unexpected failure: %v, %v", r.fatals, r.nonFatals)
		}

		// Should fail non-fatally.
		r.Reset()
		ExpectThat(&r, "hello, world", gomock.Eq("hello, mars"))
		expectNonFatal(t, &r, strings.Join([]string{
			"Expectation failed:",
			"  Wanted: is equal to hello, mars (string)",
			"  Got: hello, world (string)",
		}, "\n"))
	})

	t.Run("BasicAssert", func(t *testing.T) {
		r := testReporter{}
		// Should do nothing.
		r.Reset()
		AssertThat(&r, "hello, world", gomock.Eq("hello, world"))
		if r.HasFailures() {
			t.Errorf("got unexpected failure: %v, %v", r.fatals, r.nonFatals)
		}

		// Should fail fatally
		r.Reset()
		AssertThat(&r, "hello, world", gomock.Eq("hello, mars"))
		expectFatal(t, &r, strings.Join([]string{
			"Assertion failed:",
			"  Wanted: is equal to hello, mars (string)",
			"  Got: hello, world (string)",
		}, "\n"))
	})

	t.Run("ExpectFatal", func(t *testing.T) {
		ourError := errors.New("something bad happened")
		isOurError := gomock.Cond(func(x any) bool {
			err, ok := x.(error)
			return ok && errors.Is(err, ourError)
		})

		r := testReporter{}

		// If no fatal is thrown, the test should be marked as a failure.
		ExpectFatal(&r, isOurError, func() {
			// not panicking
		})
		expectNonFatal(t, &r, "Expected fatal error, but none occurred")

		// If the matcher succeeds, then so does ExpectFatal
		r.Reset()
		ExpectFatal(&r, isOurError, func() {
			panic(fmt.Errorf("wrapped: %w", ourError))
		})
		if r.HasFailures() {
			t.Errorf("got unexpected failure: %v, %v", r.fatals, r.nonFatals)
		}

		// If the function panics but in an unexpected way, ExpectFatal fails.
		r.Reset()
		ExpectFatal(&r, isOurError, func() {
			panic("something else")
		})
		expectNonFatal(t, &r, "Wanted: adheres to a custom condition")
	})

	t.Run("AssertFatal", func(t *testing.T) {
		ourError := errors.New("something bad happened")
		isOurError := gomock.Cond(func(x any) bool {
			err, ok := x.(error)
			return ok && errors.Is(err, ourError)
		})

		r := testReporter{}

		// If no fatal is thrown, the test should be marked as a failure.
		AssertFatal(&r, isOurError, func() {
			// not panicking
		})
		expectFatal(t, &r, "Asserted fatal error, but none occurred")

		// If the matcher succeeds, then so does AssertFatal
		r.Reset()
		AssertFatal(&r, isOurError, func() {
			panic(fmt.Errorf("wrapped: %w", ourError))
		})
		if r.HasFailures() {
			t.Errorf("got unexpected failure: %v, %v", r.fatals, r.nonFatals)
		}

		// If the function panics but in an unexpected way, AssertFatal fails.
		r.Reset()
		AssertFatal(&r, isOurError, func() {
			panic("something else")
		})
		expectFatal(t, &r, "Wanted: adheres to a custom condition")
	})
}
