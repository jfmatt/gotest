package gotest

import (
	"errors"
	"fmt"

	"go.uber.org/mock/gomock"
)

func ErrorMessage(innerMatcher any) gomock.Matcher {
	return errMsgMatcher{AsMatcher(innerMatcher)}
}

type errMsgMatcher struct {
	innerMatcher gomock.Matcher
}

func (e errMsgMatcher) Matches(x any) bool {
	if asErr, ok := x.(error); ok {
		if asErr == nil {
			return false
		}
		return e.innerMatcher.Matches(asErr.Error())
	} else {
		return false
	}
}

func (e errMsgMatcher) String() string {
	return fmt.Sprintf("is an error with message that %s", e.innerMatcher.String())
}

func ErrorIs(err error) gomock.Matcher {
	return errIsMatcher{err}
}

type errIsMatcher struct {
	err error
}

func (e errIsMatcher) Matches(x any) bool {
	if asErr, ok := x.(error); ok {
		if asErr == nil {
			return false
		}
		return errors.Is(asErr, e.err)
	} else {
		return false
	}
}

func (e errIsMatcher) String() string {
	return fmt.Sprintf("is an error wrapping %s", e.err)
}
