package gotest

import (
	"errors"
	"fmt"
)

// ErrorMessage matches errors whose error message fulfills the innerMatcher.
//
// Examples:
//
//	err := errors.New("file not found")
//	ExpectThat(t, err, ErrorMessage("file not found"))
//	ExpectThat(t, err, ErrorMessage(HasSubstr("not found")))
//	ExpectThat(t, err, Not(ErrorMessage("success")))
func ErrorMessage(innerMatcher any) Matcher {
	return errMsgMatcher{AsMatcher(innerMatcher)}
}

type errMsgMatcher struct {
	innerMatcher Matcher
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

// ErrorIs matches errors that wrap the expected error, using errors.Is().
//
// Examples:
//
//	var ErrNotFound = errors.New("not found")
//	err := fmt.Errorf("failed: %w", ErrNotFound)
//	ExpectThat(t, err, ErrorIs(ErrNotFound))
//	ExpectThat(t, err, Not(ErrorIs(errors.New("other error"))))
func ErrorIs(err error) Matcher {
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
