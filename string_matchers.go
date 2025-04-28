package gotest

import (
	"fmt"
	"regexp"
	"strings"
)

// Matches strings and byte-arrays that start with the given prefix.
//
// Examples:
//
//	ExpectThat(t, "hello", StartsWith("h"))
//	ExpectThat(t, "hello", StartsWith("hello"))
func StartsWith(s string) Matcher {
	return prefixMatcher{prefix: s}
}

type prefixMatcher struct {
	stringMatcher
	prefix string
}

func (m prefixMatcher) Matches(x any) bool {
	if asStr, ok := m.getString(x); ok {
		return strings.HasPrefix(asStr, m.prefix)
	} else {
		return false
	}
}
func (m prefixMatcher) String() string {
	return fmt.Sprintf("starts with '%s'", m.prefix)
}

// Matches strings and byte-arrays containing the given substring.
//
// Examples:
//
//	ExpectThat(t, "hello, world", HasSubstr("hello"))
//	ExpectThat(t, []byte{"a", "b", "c"}, HasSubstr("a"))
func HasSubstr(s string) Matcher {
	return substrMatcher{s: s}
}

type substrMatcher struct {
	stringMatcher
	s string
}

func (m substrMatcher) Matches(x any) bool {
	if asStr, ok := m.getString(x); ok {
		return strings.Contains(asStr, m.s)
	} else {
		return false
	}
}

func (m substrMatcher) String() string {
	return fmt.Sprintf("has substring '%s'", m.s)
}

// Matches strings and byte-arrays that match exactly the given regexp.
//
// Note that this is not the same behavior as gomock.Regex. This implementation
// is stricter by default; gomock.Regex has the behavior of our
// ContainsRegex().
//
// Examples:
//
//	ExpectThat(t, "hello", Regex("hello"))
//	ExpectThat(t, "hello", Regex("\\w+"))
//	ExpectThat(t, "hello, world", Not(Regex("\\w+")))
func Regex(r string) Matcher {
	if len(r) == 0 || r[0] != '^' {
		r = "^" + r
	}
	if r[len(r)-1] != '$' {
		r = r + "$"
	}
	return regexMatcher{r: regexp.MustCompile(r)}
}

// Matches strings and byte-arrays that contain a match for the given
// regexp. Similar behavior to gomock.Regex().
//
// This is a strictly weaker matcher than Regex() - that is, if
// ContainsRegex(r).Matches(x), then Regex(r).Matches(x) is guaranteed to be
// true.
//
// Examples:
//
//	ExpectThat(t, "hello", ContainsRegex("hello"))
//	ExpectThat(t, "hello", ContainsRegex("\\w+"))
//	ExpectThat(t, "hello, world", ContainsRegex("\\w+"))
//	ExpectThat(t, "hello, world", Not(ContainsRegex("\\d")))
func ContainsRegex(r string) Matcher {
	return regexMatcher{r: regexp.MustCompile(r)}
}

type regexMatcher struct {
	stringMatcher
	r *regexp.Regexp
}

func (r regexMatcher) Matches(x any) bool {
	if asStr, ok := x.(string); ok {
		return r.r.MatchString(asStr)
	} else if asBytes, ok := x.([]byte); ok {
		return r.r.Match(asBytes)
	} else {
		return false
	}
}

func (r regexMatcher) String() string {
	return fmt.Sprintf("matches regex '%s'", r.r)
}

// Utility mixin for string matchers. All matchers that embed this type
// should be able to support both string and []byte values.
type stringMatcher struct{}

func (stringMatcher) getString(x any) (string, bool) {
	switch v := x.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	default:
		return "", false
	}
}

func (stringMatcher) ExplainFailure(x any) (string, bool) {
	switch x.(type) {
	case string, []byte:
		return "", false
	default:
		return fmt.Sprintf("value is of type %T, not a string", x), true
	}
}
