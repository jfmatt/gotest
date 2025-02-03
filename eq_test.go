package gotest

import (
	"github.com/jfmatt/gotest/testdata"
	"google.golang.org/protobuf/proto"
)

type y struct {
	key    string
	field1 string
}

// This type has a custom equality check that only looks at one field.
func (yVal y) Equals(other any) bool {
	typed, ok := other.(y)
	return ok && typed.key == yVal.key
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

	Proto    *testdata.SomeData
	AnyProto proto.Message
}
