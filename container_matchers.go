package gotest

import (
	"fmt"
	"reflect"
	"strings"
)

// Matches values whose length fulfills `innerMatcher`. Length is defined by
// calling a `Len()` method for objects that have one, or using the `len()`
// builtin when available (i.e. for arrays, slices, maps, strings, and
// channels).
//
// Examples:
//
//	ExpectThat(t, "asdf", Len(4))
//	ExpectThat(t, "asdf", Len(Gt(2)))
//	ExpectThat(t, map[string]int{"a": 1, "b": 2}, Len(Lt(10)))
//
//	type CustomLen struct {}
//	func (CustomLen) Len() int { return 2 }
//	c := CustomLen{}
//	ExpectThat(t, c, Len(2))
func Len(innerMatcher any) Matcher {
	return lenMatcher{AsMatcher(innerMatcher)}
}

// Same behavior as Len(0), but with better error-message reporting.
func Empty() Matcher {
	return Not(true)
}

// Matches slices or arrays containing all of the provided elements, in any
// order, potentially with additional elements as well.
//
// If any of `elements` could match multiple of the elements of the value, then
// there must be a mapping such that a different element of the value fulfills
// each of `elements`. (As a corollary, matching values must have length at
// least greater than or equal to the length of `elements`.)
//
// This is a weaker test than either ElemensAreUnordered() or ElemensAre().
// That is, if either ElementsAre(els...).Matches(x) or
// ElementsAreUnordered(els...).Matches(x), then Contains(els...).Matches(x) is
// guaranteed to be true.
//
// Examples:
//
//	// matches against 'bb' and 'a'
//	ExpectThat(t, []string{"a", "bb", "ccc", "dd"}, Contains("bb", "a"))
//	// matches against 'bb' and 'a'
//	ExpectThat(t, []string{"a", "bb", "ccc", "dd"}, Contains("bb", Len(1)))
//	// matches against 'bb' and 'dd'
//	ExpectThat(t, []string{"a", "bb", "ccc", "dd"}, Contains("bb", Len(2)))
//	// no match, because 'bb' is the only element that fulfills either matcher
//	ExpectThat(t, []string{"a", "bb", "ccc", "dd"}, Not(Contains("bb", StartsWith("b"))))
func Contains(elements ...any) Matcher {
	matchers := make([]Matcher, len(elements))
	for i, el := range elements {
		matchers[i] = AsMatcher(el)
	}
	return unorderedMatcher{matchers, false}
}

// Tests that a slice or array contains exactly the provided elements, in
// order, with no others.
//
// Each element can be either an exact value (tested by equality) or a matcher
// that must succeed for that element.
//
// This is a stronger test than either ElementsAreUnordered() or Contains().
// That is, ElementsAre(els...).Matches(x) implies both
// Contains(els...).Matches(x) and ElementsAreUnordered(els...).Matches(x).
//
// Examples:
//
//	ExpectThat(t, []string{"a", "b", "c"}, ElementsAre("a", "b", "c"))
//	ExpectThat(t, []string{"a", "b", "c"}, ElementsAre("a', Len(1), Any()))
//	ExpectThat(t, []string{"a", "b", "c"}, Not(ElementsAre(Any(), "c")))
//	ExpectThat(t, []string{"a", "b", "c"}, Not(ElementsAre("b", "a", "c")))
func ElementsAre(elements ...any) Matcher {
	return Not(true)
}

// Tests that a slice or array contains exactly the provided elements, in any
// order, but with no others.
//
// Specifically, there must be a possible 1:1 mapping between `elements` and
// the value's elements such that all matchers are satisfied. (As a corollary,
// only values with the same length as `elements` can possibly match.)
//
// This is a strictly weaker test than ElementsAre(), but a stronger test than
// Contains(). That is:
//   - ElementsAreUnordered(els...).Matches(x) implies Contains(els...).Matches(x)
//   - ElementsAre(els...).Matches(x) implies ElementsAreUnordered(els...).Matches(x)
//
// Examples:
//
//	ExpectThat(t, []string{"a", "b", "ccc"}, ElementsAreUnordered("b", "ccc", "a"))
//	ExpectThat(t, []string{"a", "b", "ccc"}, ElementsAreUnordered("b", Any(), "a"))
//	ExpectThat(t, []string{"a", "b", "ccc"}, Not(ElementsAreUnordered(Any(), Any())))
//	ExpectThat(t, []string{"a", "b", "ccc"}, Not(ElementsAreUnordered("a", "ccc", Len(Gt(1)))))
func ElementsAreUnordered(elements ...any) Matcher {
	matchers := make([]Matcher, len(elements))
	for i, el := range elements {
		matchers[i] = AsMatcher(el)
	}
	return unorderedMatcher{matchers, true}
}

// Tests that a map contains exactly the elements of `mapValues`, and no
// others.
//
// The map passed to this function must have keys that are convertible to the
// key type of the value that will be tested. In most normal cases, the keys
// here should be exactly the same as those in the value to be tested.
//
// The values of the map passed to this function could be exact values, or
// matchers, or a mix.
//
// This is a stronger test than MapContains(). That is, if MapIs(m).Matches(x),
// then MapContains(m).Matches(x) is guaranteed to be true.
//
// Examples:
//
//	ExpectThat(t, map[string]int{"a": 1, "b": 10}, MapIs(map[string]int{
//		"a": 1,
//		"b": 10,
//	})
//	ExpectThat(t, map[string]int{"a": 1, "b": 10}, MapIs(map[string]any{
//		"a": 1,
//		"b": Gt(5),
//	})
//	ExpectThat(t, map[string]int{"a": 1, "b": 10}, Not(MapIs(map[string]any{
//		"a": 1,
//	}))
//	ExpectThat(t, map[string]int{"a": 1, "b": 10}, Not(MapIs(map[string]any{
//		"a": 1,
//		"b": Gt(5),
//		"c": 3,
//	}))
func MapIs[K comparable, V any](mapValues map[K]V) Matcher {
	matchers := make(map[K]Matcher)
	for k, v := range mapValues {
		matchers[k] = AsMatcher(v)
	}
	return mapMatcher[K]{matchers, true}
}

// Tests that a map contains the elements in `mapValues`, and potentially
// others as well.
//
// This is a weaker test than MapIs(). That is, if MapIs(m).Matches(x),
// then MapContains(m).Matches(x) is guaranteed to be true.
func MapContains[K comparable, V any](mapValues map[K]V) Matcher {
	matchers := make(map[K]Matcher)
	for k, v := range mapValues {
		matchers[k] = AsMatcher(v)
	}
	return mapMatcher[K]{matchers, false}
}

type KeyVal[K any, V any] struct {
	K K
	V V
}

// Tests that a map contains the key-value pairs in `pairs`.
//
// This is very similar to MapContains(), but allows using fuzzy matchers on
// the keys as well. If keys in `pairs` are Matchers, then this matcher
// tests that there is some mapping between `pairs` and the key-value pairs in
// the value such that each of `pairs` is satisfied by a different element of
// the map.
//
// In other words, this is the same as `Contains()`, but for maps - if the map
// were converted into a list of key-value pairs.
func MapContainsKVs[K any, V any](pairs ...KeyVal[K, V]) Matcher {
	pairMatchers := make([]KeyVal[Matcher, Matcher], len(pairs))
	for i, p := range pairs {
		pairMatchers[i] = KeyVal[Matcher, Matcher]{
			K: AsMatcher(p.K),
			V: AsMatcher(p.V),
		}
	}
	return mapKvMatcher{pairMatchers, false}
}

// Tests that a map contains the key-value pairs in `pairs`, and no others.
//
// This is very similar to MapIs(), but allows using fuzzy matchers on
// the keys as well. If keys in `pairs` are Matchers, then this matcher
// tests that there is some 1:1 mapping between `pairs` and the key-value pairs in
// the value such that each of `pairs` is satisfied by a different element of
// the map.
//
// In other words, this is the same as `ElementsAreUnordered()`, but for maps -
// if the map were converted into a list of key-value pairs.
func MapIsKVs[K any, V any](pairs ...KeyVal[K, V]) Matcher {
	pairMatchers := make([]KeyVal[Matcher, Matcher], len(pairs))
	for i, p := range pairs {
		pairMatchers[i] = KeyVal[Matcher, Matcher]{
			K: AsMatcher(p.K),
			V: AsMatcher(p.V),
		}
	}
	return mapKvMatcher{pairMatchers, true}
}

type mapKvMatcher struct {
	matchers []KeyVal[Matcher, Matcher]
	matchAll bool
}

func (m mapKvMatcher) Matches(x any) bool {
	return false
}

func (m mapKvMatcher) String() string {
	// TODO
	return "UNIMPL"
}

type mapMatcher[K comparable] struct {
	matchers map[K]Matcher
	matchAll bool
}

func (m mapMatcher[K]) Matches(x any) bool {
	return false
}

func (m mapMatcher[K]) String() string {
	// TODO
	return "UNIMPL"
}

type lenMatcher struct {
	innerMatcher Matcher
}

type hasLength interface {
	Len() int
}

func (l lenMatcher) Matches(x any) bool {
	if length, ok := l.getLength(x); ok {
		return l.innerMatcher.Matches(length)
	} else {
		return false
	}
}

func (l lenMatcher) String() string {
	return fmt.Sprintf("has length which %s", l.innerMatcher.String())
}

func (l lenMatcher) getLength(x any) (int, bool) {
	if lennable, ok := x.(hasLength); ok {
		return lennable.Len(), true
	}
	r := reflect.ValueOf(x)
	switch r.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return r.Len(), true
	default:
		return 0, false
	}
}

func (l lenMatcher) ExplainFailure(x any) (string, bool) {
	if length, ok := l.getLength(x); ok {
		if innerExplainer, ok := l.innerMatcher.(MismatchExplainer); ok {
			return innerExplainer.ExplainFailure(length)
		} else {
			return "", false
		}
	} else {
		return fmt.Sprintf("val is of type %T, which doesn't have a length", x), true
	}
}

type unorderedMatcher struct {
	elements []Matcher
	matchAll bool
}

func (m unorderedMatcher) Matches(x any) bool {
	// TODO
	return true
}

func (m unorderedMatcher) String() string {
	elemStrings := make([]string, len(m.elements))
	for i, el := range m.elements {
		elemStrings[i] = el.String()
	}
	return fmt.Sprintf("contains elements matching [%s]", strings.Join(elemStrings, "; "))
}

func (m unorderedMatcher) ExplainFailure(val any) string {
	// TODO
	return "didn't match"
}
