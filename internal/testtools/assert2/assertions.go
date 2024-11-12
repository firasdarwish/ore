// assert2 package add missing assertions from testify/assert package
package assert2

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/stretchr/testify/assert"
)

type tHelper interface {
	Helper()
}

type StringMatcher = func(s string) bool

// PanicsWithError asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// StringMatcher.
//
//	assert.PanicsWithError(t, ErrorStartsWith("crazy error"), func(){ GoCrazy() })
func PanicsWithError(t assert.TestingT, errStringMatcher StringMatcher, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	funcDidPanic, panicValue, panickedStack := didPanic(f)
	if !funcDidPanic {
		return assert.Fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...)
	}
	panicErr, ok := panicValue.(error)
	if !ok || !errStringMatcher(panicErr.Error()) {
		return assert.Fail(t, fmt.Sprintf("func %#v panic with unexpected Panic value:\t%#v\n\tPanic stack:\t%s", f, panicValue, panickedStack), msgAndArgs...)
	}

	return true
}

func ErrorStartsWith(prefix string) StringMatcher {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

func ErrorEndsWith(suffix string) StringMatcher {
	return func(s string) bool {
		return strings.HasSuffix(s, suffix)
	}
}

func ErrorContains(substr string) StringMatcher {
	return func(s string) bool {
		return strings.Contains(s, substr)
	}
}

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f assert.PanicTestFunc) (didPanic bool, message interface{}, stack string) {
	didPanic = true

	defer func() {
		message = recover()
		if didPanic {
			stack = string(debug.Stack())
		}
	}()

	// call the target function
	f()
	didPanic = false

	return
}
