package gotest

import (
	"strings"
	"testing"
)

func TestElementsAre(t *testing.T) {
	ExpectThat(t, []string{}, ElementsAre())
	ExpectThat(t, []string{"a"}, Not(ElementsAre()))
	ExpectThat(t, []string{"a"}, ElementsAre("a"))
	ExpectThat(t, []string{"a", "b"}, ElementsAre("a", "b"))
	ExpectThat(t, []string{"a", "b"}, Not(ElementsAre("a", "b", "c")))
	ExpectThat(t, []string{"a", "b"}, Not(ElementsAre("b", "a")))
	ExpectThat(t, []string{"a", "aaa"}, ElementsAre("a", Len(3)))
	ExpectThat(t, []string{"a", "aaa"}, ElementsAre("a", HasSubstr("aa")))

	r := &testReporter{}
	ExpectThat(r, []string{"a", "b"}, ElementsAre("a", Len(3)))
	ExpectEq(t, r.nonFatals[0], strings.Join([]string{
		"Expectation failed:",
		"  Wanted: contains elements matching [" +
			"is equal to a (string); " +
			"has length which is equal to 3 (int)]",
		"  Got: [a b] ([]string)",
		"  ...where element 1: length is 1",
	}, "\n"))
}

type SomeStruct struct {
	el0 string
	el1 string
}

func TestElementsAreUnordered(t *testing.T) {
	// Empty matcher list
	ExpectThat(t, []string{}, ElementsAreUnordered())
	ExpectThat(t, []string{"a"}, Not(ElementsAreUnordered()))

	// Empty value
	ExpectThat(t, []string{}, Not(ElementsAreUnordered("a")))

	// Basic matching cases - in any order
	ExpectThat(t, []string{}, ElementsAreUnordered())
	ExpectThat(t, []string{"a"}, Not(ElementsAreUnordered("b")))
	ExpectThat(t, []string{"a"}, ElementsAreUnordered("a"))
	ExpectThat(t, []string{"a", "b"}, ElementsAreUnordered("b", "a"))
	ExpectThat(t, []string{"a", "b"}, Not(ElementsAreUnordered("b", "c")))

	// More complex matching cases
	ExpectThat(t, []string{"a", "bb", "ab"}, ElementsAreUnordered(
		Len(2),         // matches "bb" and "ab", will be assigned to "bb"
		HasSubstr("a"), // matches "a and "ab", will be assigned to "ab"
		"a",            // matches only "a"
	))

	ExpectThat(t, []string{"a", "ab", "cc", "ddd"}, Not(ElementsAreUnordered(
		"ab",
		"cc",
		Len(2), // matches "ab" and "cc", but those are taken
		Any(),  // only option for "a" and "ddd"
	)))

	r := &testReporter{}
	ExpectThat(r, []string{"a", "c", "b"}, ElementsAreUnordered("c", "a", Len(3)))

	// Error reporting when a matcher has no corresponding value
	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: contains elements matching ["+
			"is equal to c (string); "+
			"is equal to a (string); "+
			"has length which is equal to 3 (int)]",
		"  Got: [a c b] ([]string)",
		"  ...where matcher 2 matches no elements "+
			"(wanted has length which is equal to 3 (int)); "+
			"value 2 matches no matchers",
	))

	// Error reporting when all matchers and values can be matched
	// individually, but there's no bijection.
	r.Reset()
	ExpectThat(r, []string{"a", "ab", "cc", "ddd"}, ElementsAreUnordered(
		"ab",
		"cc",
		Len(2), // matches "ab" and "cc", but those are taken
		Any(),  // only option for "a" and "ddd"
	))

	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: contains elements matching ["+
			"is equal to ab (string); "+
			"is equal to cc (string); "+
			"has length which is equal to 2 (int); "+
			"is anything]",
		"  Got: [a ab cc ddd] ([]string)",
		"  ...where no permutation could pair all matchers and values, closest match is 3/4 with "+
			"value 0 -> matcher 3; value 1 -> matcher 0; value 2 -> matcher 1",
	))

	// Error reporting on type mismatch
	r.Reset()
	ExpectThat(r, SomeStruct{"a", "b"}, ElementsAreUnordered("a", "b"))
	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: contains elements matching ["+
			"is equal to a (string); "+
			"is equal to b (string)]",
		"  Got: {a b} (gotest.SomeStruct)",
		"  ...where type gotest.SomeStruct isn't iterable",
	))
}
