package gotest

import (
	"fmt"
	"reflect"
	"slices"
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
//
// Examples:
//
//	ExpectThat(t, []int{}, Empty())
//	ExpectThat(t, "", Empty())
//	ExpectThat(t, map[string]int{}, Empty())
//	ExpectThat(t, []int{1, 2, 3}, Not(Empty()))
func Empty() Matcher {
	return emptyMatcher{}
}

type emptyMatcher struct{}

func (e emptyMatcher) Matches(x any) bool {
	if length, ok := getLength(x); ok {
		return length == 0
	}
	return false
}

func (e emptyMatcher) String() string {
	return "is empty"
}

func (e emptyMatcher) ExplainFailure(x any) (string, bool) {
	if length, ok := getLength(x); ok {
		return fmt.Sprintf("length is %d", length), true
	}
	return fmt.Sprintf("type %T doesn't have a length", x), true
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
	matchers := make([]Matcher, len(elements))
	for i, el := range elements {
		matchers[i] = AsMatcher(el)
	}
	return orderedMatcher{matchers}
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
//
// Examples:
//
//	m := map[string]int{"a": 1, "b": 2, "c": 3}
//	ExpectThat(t, m, MapContains(map[string]int{"a": 1, "b": 2}))
//	ExpectThat(t, m, MapContains(map[string]any{"a": 1, "c": Gt(2)}))
//	ExpectThat(t, m, Not(MapContains(map[string]int{"d": 4})))
func MapContains[K comparable, V any](mapValues map[K]V) Matcher {
	matchers := make(map[K]Matcher)
	for k, v := range mapValues {
		matchers[k] = AsMatcher(v)
	}
	return mapMatcher[K]{matchers, false}
}

type KeyValT struct {
	K any
	V any
}

func KeyVal(k, v any) KeyValT {
	return KeyValT{k, v}
}

type keyValMatcher struct {
	K Matcher
	V Matcher
}

func (kv *keyValMatcher) String() string {
	return fmt.Sprintf("key (%s) -> %s", kv.K.String(), kv.V.String())
}

func (kv *keyValMatcher) Matches(x any) bool {
	// x is expected to be a key-value pair, represented as a two-element
	// array or slice.
	r := reflect.ValueOf(x)
	switch r.Kind() {
	case reflect.Array, reflect.Slice:
		if r.Len() != 2 {
			return false
		}
		return kv.K.Matches(r.Index(0).Interface()) &&
			kv.V.Matches(r.Index(1).Interface())
	default:
		return false
	}
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
//
// Examples:
//
//		ExpectThat(t, map[string]int{"a": 1, "bxy": 10}, MapIsKVs(
//	     	KeyVal(StartsWith("b"), Gt(5)),
//	 	))
func MapContainsKVs(pairs ...KeyValT) Matcher {
	pairMatchers := make([]Matcher, len(pairs))
	for i, p := range pairs {
		pairMatchers[i] = &keyValMatcher{
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
//
// Examples:
//
//		ExpectThat(t, map[string]int{"a": 1, "bxy": 10}, MapIsKVs(
//			KeyVal("a", Lt(2)),
//	     	KeyVal(StartsWith("b"), Gt(5)),
//	 	))
func MapIsKVs(pairs ...KeyValT) Matcher {
	pairMatchers := make([]Matcher, len(pairs))
	for i, p := range pairs {
		pairMatchers[i] = &keyValMatcher{
			K: AsMatcher(p.K),
			V: AsMatcher(p.V),
		}
	}
	return mapKvMatcher{pairMatchers, true}
}

type mapKvMatcher struct {
	matchers []Matcher
	matchAll bool
}

func (m mapKvMatcher) Matches(x any) bool {
	if reflect.ValueOf(x).Kind() != reflect.Map {
		return false
	}

	xAsList := make([][2]any, 0)
	iter := reflect.ValueOf(x).MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		xAsList = append(xAsList, [2]any{k.Interface(), v.Interface()})
	}

	return unorderedMatcher{
		elements: m.matchers,
		matchAll: m.matchAll,
	}.Matches(xAsList)
}

func (m mapKvMatcher) String() string {
	var exact string
	if m.matchAll {
		exact = "has"
	} else {
		exact = "contains"
	}

	elemStrings := make([]string, len(m.matchers))
	for i, el := range m.matchers {
		elemStrings[i] = el.String()
	}
	return fmt.Sprintf("%s map entries [%s]",
		exact, strings.Join(elemStrings, "; "))
}

type mapMatcher[K comparable] struct {
	matchers map[K]Matcher
	matchAll bool
}

func (m mapMatcher[K]) Matches(x any) bool {
	if reflect.ValueOf(x).Kind() != reflect.Map {
		return false
	}

	for k, matcher := range m.matchers {
		val := reflect.ValueOf(x).MapIndex(reflect.ValueOf(k))
		if !val.IsValid() {
			return false
		}
		if !matcher.Matches(val.Interface()) {
			return false
		}
	}

	return !m.matchAll || (reflect.ValueOf(x).Len() == len(m.matchers))
}

func (m mapMatcher[K]) String() string {
	var exact string
	if m.matchAll {
		exact = "has"
	} else {
		exact = "contains"
	}

	parts := make([]string, 0)
	for k, matcher := range m.matchers {
		parts = append(parts, fmt.Sprintf("key %v -> %s", k, matcher.String()))
	}

	return fmt.Sprintf("%s map entries [%s]",
		exact, strings.Join(parts, "; "),
	)
}

type lenMatcher struct {
	innerMatcher Matcher
}

type hasLength interface {
	Len() int
}

func (l lenMatcher) Matches(x any) bool {
	if length, ok := getLength(x); ok {
		return l.innerMatcher.Matches(length)
	} else {
		return false
	}
}

func (l lenMatcher) String() string {
	return fmt.Sprintf("has length which %s", l.innerMatcher.String())
}

func getLength(x any) (int, bool) {
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
	if length, ok := getLength(x); ok {
		return fmt.Sprintf("length is %d", length), true
	} else {
		return fmt.Sprintf("type %T doesn't have a length", x), true
	}
}

type unorderedMatcher struct {
	elements []Matcher

	// If true, all elements in the value must be matched by matchers. If
	// false, matchers can be a subset.
	matchAll bool
}

func (m unorderedMatcher) Matches(x any) bool {
	r := reflect.ValueOf(x)
	switch r.Kind() {
	case reflect.Array, reflect.Slice:
		if m.matchAll && r.Len() != len(m.elements) {
			return false
		} else if r.Len() < len(m.elements) {
			return false
		}

		// Initialize adjacency graph based on whether each value satisfies each
		// matcher.
		matchMatrix := make([][]bool, len(m.elements))
		for i := range m.elements {
			matchMatrix[i] = make([]bool, r.Len())
			for j := range len(matchMatrix[i]) {
				matchMatrix[i][j] = m.elements[i].Matches(r.Index(j).Interface())
			}
		}

		// Short-circuit by checking if any matchers (and values, if we need a
		// full bijection) are unmatchable.
		noMatchMatchers, noMatchValues := validateMatchMatrix(matchMatrix, r.Len())
		if len(noMatchMatchers) > 0 {
			return false
		}
		if m.matchAll && len(noMatchValues) > 0 {
			return false
		}

		g := newMatcherFlowGraph(matchMatrix)
		g.Solve()

		return g.matchersMatched == len(m.elements)
	default:
		return false
	}
}

func (m unorderedMatcher) String() string {
	elemStrings := make([]string, len(m.elements))
	for i, el := range m.elements {
		elemStrings[i] = el.String()
	}
	var prefix string
	if m.matchAll {
		prefix = "has elements matching (in any order)"
	} else {
		prefix = "contains elements matching"
	}

	return fmt.Sprintf("%s [%s]",
		prefix, strings.Join(elemStrings, "; "),
	)
}

func (m unorderedMatcher) ExplainFailure(val any) (string, bool) {
	r := reflect.ValueOf(val)
	switch r.Kind() {
	case reflect.Array, reflect.Slice:
		// For legibility reasons, this function is intentionally very similar
		// to Matches(). It will return increasingly specific error messages as
		// the matcher is closer and closer to being satisfied.

		if m.matchAll && r.Len() != len(m.elements) {
			return fmt.Sprintf("%d elements expected but got %d", len(m.elements), r.Len()), true
		} else if r.Len() < len(m.elements) {
			return fmt.Sprintf("at least %d elements expected but got %d", len(m.elements), r.Len()), true
		}

		// Initialize adjacency graph based on whether each value satisfies each
		// matcher.
		matchMatrix := make([][]bool, len(m.elements))
		for i := range m.elements {
			matchMatrix[i] = make([]bool, r.Len())
			for j := range len(matchMatrix[i]) {
				matchMatrix[i][j] = m.elements[i].Matches(r.Index(j).Interface())
			}
		}

		// Short-circuit by checking if any matchers are unmatchable.
		noMatchMatchers, noMatchValues := validateMatchMatrix(matchMatrix, r.Len())
		noMatchProblems := make([]string, 0)
		for _, badMatcher := range noMatchMatchers {
			noMatchProblems = append(
				noMatchProblems,
				fmt.Sprintf("matcher %d matches no elements (wanted %s)",
					badMatcher, m.elements[badMatcher].String()))
		}

		if m.matchAll {
			for _, badValue := range noMatchValues {
				noMatchProblems = append(
					noMatchProblems,
					fmt.Sprintf("value %d matches no matchers", badValue))
			}
		}

		if len(noMatchProblems) > 0 {
			return strings.Join(noMatchProblems, "; "), true
		}

		g := newMatcherFlowGraph(matchMatrix)
		g.Solve()

		var problem string
		if m.matchAll {
			problem = fmt.Sprintf("no permutation could pair all matchers and values, closest match is %d/%d with ", g.matchersMatched, len(m.elements))
		} else {
			problem = fmt.Sprintf("no permutation could satisfy all matchers, closest match is %d/%d with ", g.matchersMatched, len(m.elements))
		}

		matches := make([]string, 0)
		for i := range g.valToMatcher {
			if g.valToMatcher[i] != -1 {
				matches = append(matches, fmt.Sprintf("value %d -> matcher %d", i, g.valToMatcher[i]))
			}
		}
		problem = problem + strings.Join(matches, "; ")
		return problem, true

	default:
		return fmt.Sprintf("type %T isn't iterable", val), true
	}
}

func validateMatchMatrix(matchMatrix [][]bool, width int) ([]int, []int) {
	noMatchMatchers := make([]int, 0)
EACH_MATCHER:
	for i := range matchMatrix {
		for j := range width {
			if matchMatrix[i][j] {
				continue EACH_MATCHER
			}
		}
		noMatchMatchers = append(noMatchMatchers, i)
	}

	noMatchValues := make([]int, 0)
EACH_VALUE:
	for j := range width {
		for i := range matchMatrix {
			if matchMatrix[i][j] {
				continue EACH_VALUE
			}
		}
		noMatchValues = append(noMatchValues, j)
	}

	return noMatchMatchers, noMatchValues
}

type matcherFlowGraph struct {
	// The adjacency matrix of the graph.
	//
	// matchMatrix[i][j] is true if matcher i matches value j.
	matchMatrix [][]bool

	// The flow graph discovered so far.
	valToMatcher []int
	matcherToVal []int

	matchersMatched int
}

func newMatcherFlowGraph(
	matchMatrix [][]bool,
) *matcherFlowGraph {
	matcherToVal := slices.Repeat([]int{-1}, len(matchMatrix))
	var valToMatcher []int
	if len(matchMatrix) > 0 {
		valToMatcher = slices.Repeat([]int{-1}, len(matchMatrix[0]))
	}

	return &matcherFlowGraph{
		matchMatrix:  matchMatrix,
		matcherToVal: matcherToVal,
		valToMatcher: valToMatcher,
	}
}

func (g *matcherFlowGraph) Solve() {
	if len(g.matchMatrix) == 0 || len(g.matchMatrix[0]) == 0 {
		return
	}

	// Like the GoogleMock implementation
	// (https://github.com/google/googletest/blob/main/googlemock/src/gmock-matchers.cc),
	// this algorithm is based on the Ford-Fulkerson method for
	// finding maximum flow in a bipartite graph. The idea is that
	// we can represent the elements of the value and the matchers
	// as two sets of nodes in a bipartite graph, and the edges
	// between them as the possible matchings.
	for matcher := range len(g.matchMatrix) {
		// Try to find a matching for this matcher.
		//
		// 'visited' prevents cycles in this particular iteration.
		visited := make([]bool, len(g.matchMatrix))
		g.tryAssign(matcher, &visited)
	}

	// Count the number of matchings.
	g.matchersMatched = 0
	for _, valMatched := range g.matcherToVal {
		if valMatched != -1 {
			g.matchersMatched++
		}
	}
}

func (g *matcherFlowGraph) tryAssign(matcher int, visited *[]bool) bool {
	// Try to find a value that matches this matcher.

	// First, look for potential matches that are currently unassigned.
	// If we find one, assign it and return.
	for j, matches := range g.matchMatrix[matcher] {
		if matches && g.valToMatcher[j] == -1 {
			g.matcherToVal[matcher] = j
			g.valToMatcher[j] = matcher
			return true
		}
	}

	// Second pass: Look for values that are already assigned to other
	// matchers. If we find one, try to reassign it to a different matcher.
	// If we can reassign it, then we can assign this matcher to the
	// value.
	for j, matches := range g.matchMatrix[matcher] {
		if matches && !(*visited)[j] {
			// value j is a potential match for this matcher.
			(*visited)[j] = true
			if g.tryAssign(g.valToMatcher[j], visited) {
				g.valToMatcher[j] = matcher
				g.matcherToVal[matcher] = j
				return true
			}
		}
	}
	return false
}

type orderedMatcher struct {
	elements []Matcher
}

func (m orderedMatcher) Matches(x any) bool {
	r := reflect.ValueOf(x)
	switch r.Kind() {
	case reflect.Array, reflect.Slice:
		if r.Len() != len(m.elements) {
			return false
		}
		for i := range r.Len() {
			if !m.elements[i].Matches(r.Index(i).Interface()) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (m orderedMatcher) ExplainFailure(val any) (string, bool) {
	parts := []string{}
	r := reflect.ValueOf(val)
	switch r.Kind() {
	case reflect.Array, reflect.Slice:
		if r.Len() != len(m.elements) {
			return fmt.Sprintf("%d elements expected but got %d", len(m.elements), r.Len()), true
		}
		for i := range r.Len() {
			if !m.elements[i].Matches(r.Index(i).Interface()) {
				var explanation string
				var useE bool
				if explainer, ok := m.elements[i].(MismatchExplainer); ok {
					explanation, useE = explainer.ExplainFailure(r.Index(i).Interface())
				}
				if !useE {
					explanation = "doesn't match"
				}

				parts = append(parts, fmt.Sprintf("element %d: %s", i, explanation))
			}
		}
	default:
		return fmt.Sprintf("val is of type %T, which isn't iterable", val), true
	}
	if len(parts) == 0 {
		return "", false
	}
	return strings.Join(parts, "; "), true
}

func (m orderedMatcher) String() string {
	elemStrings := make([]string, len(m.elements))
	for i, el := range m.elements {
		elemStrings[i] = el.String()
	}
	return fmt.Sprintf("has elements matching [%s]", strings.Join(elemStrings, "; "))
}
