package gotest

import (
	"strings"
	"testing"
)

func TestSubstr(t *testing.T) {
	// Basic operation
	ExpectThat(t, "hello, world", HasSubstr("o, w"))
	ExpectThat(t, "hello, world", HasSubstr("hello"))
	ExpectThat(t, []byte("hello"), HasSubstr("hello"))
	ExpectThat(t, "hello, world", HasSubstr("hello, world"))
	ExpectThat(t, "hello, world", Not(HasSubstr("hello, world2")))
	ExpectThat(t, "hello, world", Not(HasSubstr("x")))
	ExpectThat(t, "hello, world", Not(HasSubstr("o,w")))
	ExpectThat(t, "", Not(HasSubstr("h")))
	ExpectThat(t, 12, Not(HasSubstr("twelve")))

	// Test exact reported errors on failure
	r := testReporter{}
	ExpectThat(&r, "hello, world", HasSubstr("x"))
	ExpectEq(t, r.nonFatals[0], strings.Join([]string{
		"Expectation failed:",
		"  Wanted: has substring 'x'",
		"  Got: hello, world (string)",
	}, "\n"))

	r.Reset()
	ExpectThat(&r, 12, HasSubstr("twelve"))
	ExpectEq(t, r.nonFatals[0], strings.Join([]string{
		"Expectation failed:",
		"  Wanted: has substring 'twelve'",
		"  Got: 12 (int)",
		"  ...where value is of type int, not a string",
	}, "\n"))
}

func TestRegex(t *testing.T) {
	t.Run("ContainsRegex", func(t *testing.T) {
		ExpectThat(t, "hello, world", ContainsRegex("hello.*"))
		ExpectThat(t, "hello, world", ContainsRegex("hello"))
		ExpectThat(t, "hello, world", ContainsRegex("\\w+, \\w+"))
		ExpectThat(t, "hello, world", Not(ContainsRegex("hello$")))
		ExpectThat(t, "hello, world", Not(ContainsRegex("\\d+")))
		ExpectThat(t, []byte{'a', 'b', 'c'}, ContainsRegex("\\w{3}"))
		ExpectThat(t, []byte{'1', '2', '3'}, ContainsRegex("\\d+"))
		ExpectThat(t, 12, Not(ContainsRegex("twelve")))

		r := testReporter{}
		ExpectThat(&r, "hello, world", ContainsRegex("hello$"))
		ExpectEq(t, r.nonFatals[0], strings.Join([]string{
			"Expectation failed:",
			"  Wanted: matches regex 'hello$'",
			"  Got: hello, world (string)",
		}, "\n"))

		r.Reset()
		ExpectThat(&r, 12, ContainsRegex("\\d+"))
		ExpectEq(t, r.nonFatals[0], strings.Join([]string{
			"Expectation failed:",
			"  Wanted: matches regex '\\d+'",
			"  Got: 12 (int)",
			"  ...where value is of type int, not a string",
		}, "\n"))
	})

	t.Run("ExactRegex", func(t *testing.T) {
		// This matcher adds an implicit ^$ to enforce that the string
		// /exactly/ matches the regex, rather than /containing/ a match of
		// the regex. Note the second case here - Regex("hello") fails where
		// "ContainsRegex("hello") succeeds.
		ExpectThat(t, "hello, world", Regex("hello.*"))
		ExpectThat(t, "hello, world", Not(Regex("hello")))
		ExpectThat(t, "hello, world", Regex("\\w+, \\w+"))
		// If the regex already has ^ or $, that's ok too
		ExpectThat(t, "hello, world", Regex("^\\w+, \\w+$"))
		ExpectThat(t, "hello, world", Not(Regex("hello$")))
		ExpectThat(t, "hello, world", Not(Regex("\\d+")))
		ExpectThat(t, []byte{'a', 'b', 'c'}, Regex("\\w{3}"))
		ExpectThat(t, []byte{'1', '2', '3'}, Regex("\\d+"))
		ExpectThat(t, 12, Not(Regex("twelve")))

		r := testReporter{}
		ExpectThat(&r, "hello, world", Regex("hello$"))
		ExpectEq(t, r.nonFatals[0], strings.Join([]string{
			"Expectation failed:",
			"  Wanted: matches regex '^hello$'",
			"  Got: hello, world (string)",
		}, "\n"))

		r.Reset()
		ExpectThat(&r, 12, Regex("\\d+"))
		ExpectEq(t, r.nonFatals[0], strings.Join([]string{
			"Expectation failed:",
			"  Wanted: matches regex '^\\d+$'",
			"  Got: 12 (int)",
			"  ...where value is of type int, not a string",
		}, "\n"))
	})
}
