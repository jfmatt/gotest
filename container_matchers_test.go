package gotest

import (
	"math/rand/v2"
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
		"  Wanted: has elements matching [" +
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
	ExpectThat(t, []string{"a"}, Not(ElementsAreUnordered("b")))
	ExpectThat(t, []string{"a"}, ElementsAreUnordered("a"))
	ExpectThat(t, []string{"a", "b"}, ElementsAreUnordered("b", "a"))
	ExpectThat(t, []string{"a", "b"}, Not(ElementsAreUnordered("b", "c")))

	// Using nested matchers
	ExpectThat(t, []string{"a", "b"}, ElementsAreUnordered(
		"b",
		Len(1),
	))
	ExpectThat(t, []string{"a", "b"}, Not(ElementsAreUnordered(
		"b",
		Len(2),
	)))

	// More values than matchers
	ExpectThat(t, []string{"a", "b"}, Not(ElementsAreUnordered("b")))
	ExpectThat(t, []string{"a", "b"}, Not(ElementsAreUnordered(Len(1))))

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

	values := []string{
		"a",
		"bqq",
		"abc",
		"ABC",
		"xyz",
	}

	matchers := []any{
		"a",
		HasSubstr("b"),
		Len(3),
		StartsWith("x"),
		StartsWith("AB"),
	}

	for range 20 {
		rand.Shuffle(len(values), func(i, j int) {
			values[i], values[j] = values[j], values[i]
		})
		ExpectThat(t, values, ElementsAreUnordered(matchers...))
	}

	r := &testReporter{}
	ExpectThat(r, []string{"a", "c", "b"}, ElementsAreUnordered("c", "a", Len(3)))

	// Error reporting when a matcher has no corresponding value
	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: has elements matching (in any order) ["+
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
		"  Wanted: has elements matching (in any order) ["+
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
		"  Wanted: has elements matching (in any order) ["+
			"is equal to a (string); "+
			"is equal to b (string)]",
		"  Got: {a b} (gotest.SomeStruct)",
		"  ...where type gotest.SomeStruct isn't iterable",
	))
}

func TestContains(t *testing.T) {
	// Empty matcher list
	ExpectThat(t, []string{}, Contains())

	// Empty value
	ExpectThat(t, []string{}, Not(Contains("a")))

	// Basic matching cases - in any order
	ExpectThat(t, []string{"a"}, Not(Contains("b")))
	ExpectThat(t, []string{"a"}, Contains("a"))
	ExpectThat(t, []string{"a", "b"}, Contains("b", "a"))
	ExpectThat(t, []string{"a", "b"}, Not(Contains("b", "c")))

	// Using nested matchers
	ExpectThat(t, []string{"a", "b"}, Contains(
		"b",
		Len(1),
	))
	ExpectThat(t, []string{"a", "b"}, Not(Contains(
		"b",
		Len(2),
	)))

	// More values than matchers.
	//
	// This is the only difference between Contains() and
	// ElementsAreUnordered() - these cases match, where they don't above.
	ExpectThat(t, []string{"a"}, Contains())
	ExpectThat(t, []string{"a", "b"}, Contains("b"))
	ExpectThat(t, []string{"a", "b"}, Contains(Len(1)))

	// More complex matching cases
	ExpectThat(t, []string{"a", "bb", "ab"}, Contains(
		Len(2),         // matches "bb" and "ab", will be assigned to "bb"
		HasSubstr("a"), // matches "a and "ab", will be assigned to "ab"
		"a",            // matches only "a"
	))

	ExpectThat(t, []string{"a", "ab", "cc", "ddd"}, Not(Contains(
		"ab",
		"cc",
		Len(2), // matches "ab" and "cc", but those are taken
		Any(),  // only option for "a" and "ddd"
	)))

	values := []string{
		"a",
		"bqq",
		"abc",
		"ABC",
		"xyz",
	}

	// Note fewer matchers than values - that's ok with Contains()
	matchers := []any{
		"a",
		HasSubstr("b"),
		Len(3),
		StartsWith("AB"),
	}

	for range 20 {
		rand.Shuffle(len(values), func(i, j int) {
			values[i], values[j] = values[j], values[i]
		})
		ExpectThat(t, values, Contains(matchers...))
	}

	r := &testReporter{}
	ExpectThat(r, []string{"a", "b"}, Contains("b", "c"))
	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: contains elements matching ["+
			"is equal to b (string); "+
			"is equal to c (string)]",
		"  Got: [a b] ([]string)",
		"  ...where matcher 1 matches no elements "+
			"(wanted is equal to c (string))",
	))

	r.Reset()
	ExpectThat(r, []string{"a", "bbb", "ccc"}, Contains(Len(1), "a"))
	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: contains elements matching ["+
			"has length which is equal to 1 (int); "+
			"is equal to a (string)]",
		"  Got: [a bbb ccc] ([]string)",
		"  ...where no permutation could satisfy all matchers, "+
			"closest match is 1/2 with value 0 -> matcher 0",
	))
}

type TestStruct struct {
	Name  string
	Value int
}

func TestMapIs(t *testing.T) {
	// Empty map
	ExpectThat(t, map[string]int{}, MapIs(map[string]int{}))
	ExpectThat(t, map[string]int{"a": 1}, Not(MapIs(map[string]int{})))

	// Basic exact matching
	ExpectThat(t, map[string]int{"a": 1}, MapIs(map[string]int{"a": 1}))
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, MapIs(map[string]int{"a": 1, "b": 2}))
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, MapIs(map[string]int{"b": 2, "a": 1}))

	// Missing key
	ExpectThat(t, map[string]int{"a": 1}, Not(MapIs(map[string]int{"a": 1, "b": 2})))

	// Extra key
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, Not(MapIs(map[string]int{"a": 1})))

	// Wrong value
	ExpectThat(t, map[string]int{"a": 1}, Not(MapIs(map[string]int{"a": 2})))

	// Using matchers as values
	ExpectThat(t, map[string]int{"a": 1, "b": 10}, MapIs(map[string]any{
		"a": 1,
		"b": Gt(5),
	}))
	ExpectThat(t, map[string]int{"a": 1, "b": 10}, MapIs(map[string]any{
		"a": Lt(5),
		"b": Gt(5),
	}))
	ExpectThat(t, map[string]int{"a": 1, "b": 10}, Not(MapIs(map[string]any{
		"a": Gt(5),
		"b": Gt(5),
	})))

	// Map of structs
	ExpectThat(t,
		map[string]TestStruct{
			"user1": {"Alice", 30},
			"user2": {"Bob", 25},
		},
		MapIs(map[string]any{
			"user1": TestStruct{"Alice", 30},
			"user2": TestStruct{"Bob", 25},
		}),
	)

	// Map of structs with exact matching
	ExpectThat(t,
		map[string]TestStruct{
			"user1": {"Alice", 30},
			"user2": {"Bob", 25},
		},
		MapIs(map[string]any{
			"user1": TestStruct{"Alice", 30},
			"user2": TestStruct{"Bob", 25},
		}),
	)

	// Map of slices
	ExpectThat(t,
		map[string][]int{
			"evens": {2, 4, 6},
			"odds":  {1, 3, 5},
		},
		MapIs(map[string]any{
			"evens": []int{2, 4, 6},
			"odds":  []int{1, 3, 5},
		}),
	)

	// Map of slices with matchers
	ExpectThat(t,
		map[string][]int{
			"evens": {2, 4, 6},
			"odds":  {1, 3, 5},
		},
		MapIs(map[string]any{
			"evens": ElementsAreUnordered(4, 2, 6),
			"odds":  Contains(1, 5),
		}),
	)

	// Map of maps
	ExpectThat(t,
		map[string]map[string]int{
			"group1": {"a": 1, "b": 2},
			"group2": {"x": 10, "y": 20},
		},
		MapIs(map[string]any{
			"group1": map[string]int{"a": 1, "b": 2},
			"group2": map[string]int{"x": 10, "y": 20},
		}),
	)

	// Map of maps with matchers
	ExpectThat(t,
		map[string]map[string]int{
			"group1": {"a": 1, "b": 2},
			"group2": {"x": 10, "y": 20},
		},
		MapIs(map[string]any{
			"group1": MapIs(map[string]any{"a": Lt(5), "b": 2}),
			"group2": MapContains(map[string]any{"x": 10}),
		}),
	)

	// Complex nested structures: map of slices of structs
	ExpectThat(t,
		map[string][]TestStruct{
			"team1": {{"Alice", 30}, {"Bob", 25}},
			"team2": {{"Charlie", 35}},
		},
		MapIs(map[string]any{
			"team1": ElementsAre(
				TestStruct{"Alice", 30},
				TestStruct{"Bob", 25},
			),
			"team2": ElementsAre(
				TestStruct{"Charlie", 35},
			),
		}),
	)

	// Error message formatting - should say "has" for exact match
	r := &testReporter{}
	ExpectThat(r, map[string]int{"a": 1}, MapIs(map[string]any{
		"a": Gt(5),
	}))
	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: has map entries [key a -> is greater than 5 (int)]",
		"  Got: map[a:1] (map[string]int)",
	))
}

func TestMapContains(t *testing.T) {
	// Empty map
	ExpectThat(t, map[string]int{}, MapContains(map[string]int{}))
	ExpectThat(t, map[string]int{"a": 1}, MapContains(map[string]int{}))

	// Basic matching with extra keys allowed
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, MapContains(map[string]int{"a": 1}))
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, MapContains(map[string]int{"b": 2}))
	ExpectThat(t, map[string]int{"a": 1, "b": 2, "c": 3}, MapContains(map[string]int{"a": 1, "c": 3}))

	// Missing key
	ExpectThat(t, map[string]int{"a": 1}, Not(MapContains(map[string]int{"b": 2})))

	// Wrong value
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, Not(MapContains(map[string]int{"a": 2})))

	// Using matchers as values
	ExpectThat(t, map[string]int{"a": 1, "b": 10, "c": 5}, MapContains(map[string]any{
		"b": Gt(5),
	}))

	// Map of structs with exact matching
	ExpectThat(t,
		map[string]TestStruct{
			"user1": {"Alice", 30},
			"user2": {"Bob", 25},
			"user3": {"Charlie", 35},
		},
		MapContains(map[string]any{
			"user1": TestStruct{"Alice", 30},
			"user3": TestStruct{"Charlie", 35},
		}),
	)

	// Map of slices with matchers
	ExpectThat(t,
		map[string][]string{
			"fruits":     {"apple", "banana"},
			"vegetables": {"carrot", "lettuce"},
			"grains":     {"rice", "wheat"},
		},
		MapContains(map[string]any{
			"fruits": Contains("apple"),
		}),
	)

	// Map of maps with matchers
	ExpectThat(t,
		map[string]map[string]int{
			"group1": {"a": 1, "b": 2},
			"group2": {"x": 10, "y": 20},
			"group3": {"m": 100, "n": 200},
		},
		MapContains(map[string]any{
			"group1": MapContains(map[string]any{"a": 1}),
			"group3": MapIs(map[string]any{"m": 100, "n": 200}),
		}),
	)

	// Error message formatting - should say "contains" for subset match
	r := &testReporter{}
	ExpectThat(r, map[string]int{"a": 1, "b": 2}, MapContains(map[string]any{
		"c": 3,
	}))
	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: contains map entries [key c -> is equal to 3 (int)]",
		"  Got: map[a:1 b:2] (map[string]int)",
	))
}

func TestMapIsKVs(t *testing.T) {
	// Empty map
	ExpectThat(t, map[string]int{}, MapIsKVs())
	ExpectThat(t, map[string]int{"a": 1}, Not(MapIsKVs()))

	// Basic exact matching with KeyVal
	ExpectThat(t, map[string]int{"a": 1}, MapIsKVs(KeyVal("a", 1)))
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, MapIsKVs(
		KeyVal("a", 1),
		KeyVal("b", 2),
	))

	// Order doesn't matter
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, MapIsKVs(
		KeyVal("b", 2),
		KeyVal("a", 1),
	))

	// Matchers on values
	ExpectThat(t, map[string]int{"a": 1, "b": 10}, MapIsKVs(
		KeyVal("a", Lt(5)),
		KeyVal("b", Gt(5)),
	))

	// Matchers on keys
	ExpectThat(t, map[string]int{"apple": 1, "banana": 2}, MapIsKVs(
		KeyVal(StartsWith("a"), 1),
		KeyVal(StartsWith("b"), 2),
	))

	// Matchers on both keys and values
	ExpectThat(t, map[string]int{"apple": 10, "banana": 20}, MapIsKVs(
		KeyVal(HasSubstr("app"), Gt(5)),
		KeyVal(HasSubstr("ban"), Gt(15)),
	))

	// Missing key-value pair
	ExpectThat(t, map[string]int{"a": 1}, Not(MapIsKVs(
		KeyVal("a", 1),
		KeyVal("b", 2),
	)))

	// Extra key-value pair
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, Not(MapIsKVs(
		KeyVal("a", 1),
	)))

	// Map of structs with key matchers
	ExpectThat(t,
		map[string]TestStruct{
			"user1": {"Alice", 30},
			"user2": {"Bob", 25},
		},
		MapIsKVs(
			KeyVal(StartsWith("user"), TestStruct{"Alice", 30}),
			KeyVal(StartsWith("user"), TestStruct{"Bob", 25}),
		),
	)

	// Map of slices with matchers
	ExpectThat(t,
		map[string][]int{
			"evens": {2, 4, 6},
			"odds":  {1, 3, 5},
		},
		MapIsKVs(
			KeyVal("evens", ElementsAre(2, 4, 6)),
			KeyVal("odds", ElementsAreUnordered(5, 3, 1)),
		),
	)

	// Complex nested: map with int keys
	ExpectThat(t,
		map[int][]string{
			1: {"a", "b"},
			2: {"x", "y"},
		},
		MapIsKVs(
			KeyVal(Lt(2), Contains("a")),
			KeyVal(Gt(1), Contains("x")),
		),
	)

	// Error message formatting - should say "has" for exact match
	r := &testReporter{}
	ExpectThat(r, map[string]int{"a": 1}, MapIsKVs(
		KeyVal("b", 2),
	))
	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: has map entries [key (is equal to b (string)) -> is equal to 2 (int)]",
		"  Got: map[a:1] (map[string]int)",
	))
}

func TestMapContainsKVs(t *testing.T) {
	// Empty matcher list
	ExpectThat(t, map[string]int{}, MapContainsKVs())
	ExpectThat(t, map[string]int{"a": 1}, MapContainsKVs())

	// Basic matching with extra pairs allowed
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, MapContainsKVs(
		KeyVal("a", 1),
	))
	ExpectThat(t, map[string]int{"a": 1, "b": 2, "c": 3}, MapContainsKVs(
		KeyVal("a", 1),
		KeyVal("c", 3),
	))

	// Matchers on keys and values
	ExpectThat(t, map[string]int{"apple": 10, "banana": 20, "cherry": 30}, MapContainsKVs(
		KeyVal(StartsWith("a"), Gt(5)),
		KeyVal(HasSubstr("err"), Gt(25)),
	))

	// Missing key-value pair
	ExpectThat(t, map[string]int{"a": 1, "b": 2}, Not(MapContainsKVs(
		KeyVal("c", 3),
	)))

	// Map of structs with key matchers
	ExpectThat(t,
		map[string]TestStruct{
			"user1": {"Alice", 30},
			"user2": {"Bob", 25},
			"user3": {"Charlie", 35},
		},
		MapContainsKVs(
			KeyVal(StartsWith("user"), TestStruct{"Alice", 30}),
			KeyVal(HasSubstr("3"), TestStruct{"Charlie", 35}),
		),
	)

	// Map of slices with matchers - extra entries allowed
	ExpectThat(t,
		map[string][]string{
			"fruits":     {"apple", "banana"},
			"vegetables": {"carrot", "lettuce"},
			"grains":     {"rice", "wheat"},
		},
		MapContainsKVs(
			KeyVal("fruits", Contains("apple")),
		),
	)

	// Map of maps with complex matchers
	ExpectThat(t,
		map[string]map[string]int{
			"group1": {"a": 1, "b": 2},
			"group2": {"x": 10, "y": 20},
			"group3": {"m": 100, "n": 200},
		},
		MapContainsKVs(
			KeyVal(HasSubstr("1"), MapContains(map[string]any{"a": 1})),
			KeyVal(HasSubstr("3"), MapIs(map[string]any{"m": 100, "n": 200})),
		),
	)

	// Map with int keys and slice values
	ExpectThat(t,
		map[int][]TestStruct{
			1: {{"Alice", 30}, {"Bob", 25}},
			2: {{"Charlie", 35}},
			3: {{"David", 40}, {"Eve", 45}},
		},
		MapContainsKVs(
			KeyVal(1, ElementsAreUnordered(
				TestStruct{"Bob", 25},
				TestStruct{"Alice", 30},
			)),
			KeyVal(Gt(2), Contains(TestStruct{"David", 40})),
		),
	)

	// Error message formatting - should say "contains" for subset match
	r := &testReporter{}
	ExpectThat(r, map[string]int{"a": 1, "b": 2}, MapContainsKVs(
		KeyVal("c", 3),
	))
	ExpectThat(t, strings.Split(r.nonFatals[0], "\n"), ElementsAre(
		"Expectation failed:",
		"  Wanted: contains map entries [key (is equal to c (string)) -> is equal to 3 (int)]",
		"  Got: map[a:1 b:2] (map[string]int)",
	))
}
