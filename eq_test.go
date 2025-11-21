package gotest_test

import (
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	. "github.com/jfmatt/gotest"
	"github.com/jfmatt/gotest/testdata"
)

type y struct {
	key    string
	field1 string
}

// This type has a custom equality check that only looks at one field.
func (yVal y) Equal(other any) bool {
	if typed, ok := other.(y); ok {
		return typed.key == yVal.key
	} else if typedPtr, ok := other.(*y); ok {
		return typedPtr.key == yVal.key
	}
	return false
}

type z struct {
	field string
}

type x struct {
	PublicString  string
	privateString string

	Struct  z
	Struct2 z

	List    []int
	ObjList []y
	PtrList []*y

	Untyped        any
	privateUntyped any

	Recursion      *x
	RecursionByIfc any

	Proto        *testdata.SomeData
	AnyProto     proto.Message
	privateProto *testdata.SomeData

	anyList []any
}

type Name string

func TestEqual(t *testing.T) {
	// Basic case
	ExpectEq(t, x{}, x{})

	// Primitives
	ExpectEq(t, "abc", "abc")
	ExpectThat(t, "abc", Not(Eq("def")))
	ExpectEq(t, 123, 123)
	ExpectThat(t, 123, Not(Eq(456)))
	ExpectThat(t, 123, Not(Eq(int64(123))))

	// Named vs. unnamed types with same underlying type
	var n Name = "abc"
	ExpectEq(t, n, "abc")

	// With a private field
	ExpectThat(t, x{privateString: "a"}, Not(Eq(x{privateString: "b"})))

	// Pointers to equivalent structs
	ExpectEq(t, &x{PublicString: "a"}, &x{PublicString: "a"})
	x1 := &x{Recursion: &x{PublicString: "a"}}
	x2 := &x{Recursion: &x{PublicString: "a"}}
	x3 := &x{Recursion: &x{PublicString: "b"}}
	ExpectThat(t, x1, Eq(x2))
	ExpectThat(t, x1, Not(Eq(x3)))

	// Cycles
	x1.Recursion.Recursion = x1
	x2.Recursion.Recursion = x2
	x3.Recursion.Recursion = x3
	ExpectThat(t, x1, Eq(x2))
	ExpectThat(t, x1, Not(Eq(x3)))

	// Slices
	ExpectEq(t, x{List: []int{1, 2, 3}}, x{List: []int{1, 2, 3}})
	ExpectThat(t, x{List: []int{1, 2, 3}}, Not(Eq(x{List: []int{1, 2, 4}})))

	// Slices of objects with custom equality - ignores field1
	ExpectEq(t,
		x{ObjList: []y{{key: "a", field1: "x"}, {key: "b", field1: "y"}}},
		x{ObjList: []y{{key: "a", field1: "different"}, {key: "b", field1: "also different"}}},
	)
	ExpectThat(t,
		x{ObjList: []y{{key: "a", field1: "x"}, {key: "b", field1: "y"}}},
		Not(Eq(x{ObjList: []y{{key: "a", field1: "x"}, {key: "c", field1: "y"}}})),
	)

	// Slices of pointers to objects with custom equality - ignores field1
	ExpectEq(t,
		x{PtrList: []*y{{key: "a", field1: "x"}, {key: "b", field1: "y"}}},
		x{PtrList: []*y{{key: "a", field1: "different"}, {key: "b", field1: "also different"}}},
	)
	ExpectThat(t,
		x{PtrList: []*y{{key: "a", field1: "x"}, {key: "b", field1: "y"}}},
		Not(Eq(x{PtrList: []*y{{key: "a", field1: "x"}, {key: "c", field1: "y"}}})),
	)

	// Untyped fields
	ExpectEq(t,
		x{Untyped: z{field: "a"}, privateUntyped: z{field: "x"}},
		x{Untyped: z{field: "a"}, privateUntyped: z{field: "x"}},
	)
	ExpectThat(t,
		x{Untyped: z{field: "a"}, privateUntyped: z{field: "x"}},
		Not(Eq(x{Untyped: z{field: "b"}, privateUntyped: z{field: "x"}})),
	)
	ExpectThat(t,
		x{Untyped: int32(123)},
		Not(Eq(x{Untyped: uint32(123)})),
	)
	ExpectThat(t,
		x{Untyped: z{field: "a"}, privateUntyped: 123},
		Not(Eq(x{Untyped: z{field: "a"}, privateUntyped: 456})),
	)
	ExpectThat(t,
		x{Untyped: z{field: "a"}, privateUntyped: z{field: "x"}},
		Not(Eq(x{Untyped: z{field: "a"}, privateUntyped: z{field: "y"}})),
	)

	// Unexported fields from external packages should NOT be compared
	// These two values have the same public field but different private fields,
	// so they should be considered equal (private field is ignored)
	external1 := testdata.NewExternalType("public", "private1")
	external2 := testdata.NewExternalType("public", "private2")
	ExpectEq(t, external1, external2)

	// But if the public field differs, they should not be equal
	external3 := testdata.NewExternalType("different", "private1")
	ExpectThat(t, external1, Not(Eq(external3)))
}

func TestEqual_Protos(t *testing.T) {
	// We'll use two equivalent protos where one has a nil list and another has
	// an empty list - in proto semantics these are equivalent, so our Eq
	// matcher should treat them as equal, but naive reflect.DeepEqual would
	// not.
	nilListProto := &testdata.SomeData{I: 1, A: "test", L: nil}
	emptyListProto := &testdata.SomeData{I: 1, A: "test", L: []string{}}
	differentProto := &testdata.SomeData{I: 1, A: "test", L: []string{"not empty"}}
	ExpectThat(t, nilListProto, Not(gomock.Eq(emptyListProto))) // sanity check
	ExpectThat(t, nilListProto, Eq(emptyListProto))             // proto.Equal works
	ExpectThat(t, nilListProto, Not(Eq(differentProto)))        // different protos

	// Nested within our test struct
	ExpectEq(t,
		x{Proto: nilListProto},
		x{Proto: emptyListProto},
	)
	ExpectThat(t,
		x{Proto: nilListProto},
		Not(x{Proto: differentProto}),
	)

	// Nested within a private field
	ExpectEq(t,
		x{privateProto: nilListProto},
		x{privateProto: emptyListProto},
	)
	ExpectThat(t,
		x{privateProto: nilListProto},
		Not(x{privateProto: differentProto}),
	)

	// Nested within a `proto.Message` field
	ExpectEq(t,
		x{AnyProto: nilListProto},
		x{AnyProto: emptyListProto},
	)
	ExpectThat(t,
		x{AnyProto: nilListProto},
		Not(x{AnyProto: differentProto}),
	)

	// Nested within an `any` field
	ExpectEq(t,
		x{RecursionByIfc: nilListProto},
		x{RecursionByIfc: emptyListProto},
	)
	ExpectThat(t,
		x{RecursionByIfc: nilListProto},
		Not(x{RecursionByIfc: differentProto}),
	)

	// Slices of `any` containing protos
	ExpectEq(t,
		x{anyList: []any{nilListProto, "string", 123}},
		x{anyList: []any{emptyListProto, "string", 123}},
	)
	ExpectThat(t,
		x{anyList: []any{nilListProto, "string", 123}},
		Not(Eq(x{anyList: []any{differentProto, "string", 123}})),
	)
}
