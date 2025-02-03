package gotest

import (
	"fmt"

	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func Eq(x any) Matcher {
	return eqMatcher{val: x}
}

type eqMatcher struct {
	val any
}

func (e eqMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", e.val, e.val)
}

func (e eqMatcher) Matches(x any) bool {
	// TODO: Implement full version based on DeepEqual to handle nested proto
	// fields.
	if asProto, ok := x.(proto.Message); ok {
		if expectedAsProto, ok := e.val.(proto.Message); ok {
			return proto.Equal(asProto, expectedAsProto)
		} else {
			return false
		}
	}

	return gomock.Eq(e.val).Matches(x)
}
