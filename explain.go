package gotest

import (
	"fmt"

	"go.uber.org/mock/gomock"
)

// An optional interface that matchers can implement to provide more
// information about failures. Useful for compound matchers - e.g., those that
// match many elements of a container - to aid in debugging.
type MismatchExplainer interface {
	// Returns an explanation if one is useful, or ("", false) if no further
	// explanation should be provided
	ExplainFailure(val any) (string, bool)
}

func formatGot(val any, matcher gomock.Matcher) string {
	if asFormatter, ok := matcher.(gomock.GotFormatter); ok {
		return asFormatter.Got(val)
	} else {
		return fmt.Sprintf("%v (%[1]T)", val)
	}
}
