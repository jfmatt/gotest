# gotest

[![Go
Reference](https://pkg.go.dev/badge/github.com/jfmatt/gotest.svg)](https://pkg.go.dev/github.com/jfmatt/gotest) [![Test status](https://github.com/jfmatt/gotest/actions/workflows/go.yml/badge.svg)](https://github.com/jfmatt/gotest/actions/workflows/go.yml)

A unit testing library for Golang

Based on [GoogleTest](https://github.com/google/googletest).
Built with, and interoperable with, [GoMock](https://github.com/uber-go/mock).

This library includes a collection of **composable** `Matcher`
implementations, which can be put together to form partial matches, match
elements of collection, and so on.  They can be used to match calls to a mock
object, or as standalone test assertions.

This package is also natively protobuf-aware. Protobufs in Go contain control
fields and mutexes that are not comparable, and so must be compared with
their own reflection library. Go-native operations such as
`reflect.DeepEqual()` do not know how to do so. The implementation of `Eq()`
included in this package can handle protos, including fields of type proto
nested inside other structs, so that matchers like `Eq()` work as intended.

## Assertions

The library provides two ways to set test expectations/assertions:

* `Expect` for non-fatal assertions (failure will continue the test)
* `Assert` for fatal assertions (failure will terminate
  the test). Useful for checking preconditions without which later test code
  will panic.

Each family has a few variants:
* `[Expect|Assert]That(t, value, matcher)`: the most general form, takes any
  `Matcher` and checks `value` against it.
* `[Expect|Assert]Eq(t, value, expected)`: shorthand for equality matching (with `Eq()`)
* `[Expect|Assert]Fatal(t, errMatcher, f)`: runs function `f` and checks that it causes
  a panic matching `errMatcher`.

## Matcher Composition

Matcher implementations are provided for common conditions - equality, string
operations, numeric comparisons.

More powerful matchers can be composed to express complex structural
expectations on container types (slices and maps). Any composable matcher
supports either sub-matchers to match elements, or raw values to be checked by
equality (wrapped in an `Eq()` matcher).

For example, validating a nested map-of-slices structure where different slices
have different ordering requirements:

```go
teams := map[string][]User{
    "engineering": {{"Alice", 30}, {"Bob", 25}},
    "product": {{"Charlie", 35}},
    "exec": {{"Eve", 40}, {"Dave", 45}},
}

ExpectThat(t, teams, MapIs(map[string]any{
    "engineering": ElementsAreUnordered(
        User{"Bob", 25},
        User{"Alice", 30},
    ),
    "product": []User{{"Charlie", 35}},
    "exec": Contains(User{"Eve", 40}),
}))
```

This validates that:
* the map has exactly these three keys,
* that "engineering" contains these two users in any order,
* that "product" contains exactly this one user (note that we did not need to wrap the
  slice in `Eq()` or `ElementsAre()`; raw values are automatically wrapped),
* that "exec" contains at least the user "Eve" (other users may be present).

As shown in the example, some matchers support partial matching. For instance,
the `MapContains` matcher checks for the presence of specific keys while
ignoring others:

```go
config := map[string]map[string]int{
    "group1": {"a": 1, "b": 2, "c": 3},
    "group2": {"x": 10, "y": 20, "z": 30},
    "group3": {"m": 100, "n": 200},
}

ExpectThat(t, config, MapContains(map[string]any{
    "group1": MapContains(map[string]any{"a": Lt(5)}),
    "group3": MapIs(map[string]any{"m": 100, "n": 200}),
}))
```

This checks that "group1" exists and contains a key "a" with a value less than
5 (other keys in "group1" are ignored), and that "group3" contains exactly the
keys "m" and "n" with those specific values. The "group2" entry is not checked
at all.

## Mocks

`gotest` matchers integrate with [GoMock](https://github.com/uber-go/mock) `EXPECT()` calls. Given an interface:

```go
//go:generate mockgen -destination=mocks/storage_mock.go -package=mocks . Storage

type Storage interface {
    Save(userID string, data map[string]any) error
    Query(filters map[string]any) ([]map[string]any, error)
}
```

Matchers can specify argument expectations that go beyond exact equality:

```go
mockStorage.EXPECT().Save(
    StartsWith("user_"),
    MapContains(map[string]any{
        "email":    ContainsRegex(".*@example\\.com"),
        "age":      Gt(18),
        "verified": true,
    }),
).Return(nil)
```

This expects any call where the first argument starts with "user_" and the
second argument is a map containing an "email" key matching the regex, an "age"
key greater than 18, and a "verified" key equal to true. Other keys in the map
are ignored.

A longer example, validating a full protocol:

```go
gomock.InOrder(
    // First call: save batch metadata
    mockStorage.EXPECT().Save(
        Regex("batch_[0-9]+"),
        MapIs(map[string]any{
            "status":     "pending",
            "item_count": Gt(0),
            "started_at": Any(),
        }),
    ).Return(nil),

    // Multiple calls: save items
    mockStorage.EXPECT().Save(
        StartsWith("item_"),
        MapContains(map[string]any{
            "batch_id": StartsWith("batch_"),
        }),
    ).MinTimes(1),

    // Final call: mark complete
    mockStorage.EXPECT().Save(
        Regex("batch_[0-9]+"),
        MapContains(map[string]any{"status": "completed"}),
    ).Return(nil),
)
```
