package gotest

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
)

// Tests whether x is "equal" to the expected value.
//
// Equality is defined as follows:
//
//   - Primitive types (ints, strings, etc) are compared using ==.
//   - All protos are compared using proto.Equal, including when nested within
//     other structs, slices, or maps.
//   - Non-proto structs are compared field-by-field. Unexported fields are
//     compared only for types that are defined in the same package as the matcher
//     is used.
//   - Any type that has a custom Equal method will use that method for comparison.
//
// This is the default matcher used by all other matchers to compare nested
// values when passed values directly instead of matchers. If you want to
// customize the comparison behavior, use Equiv().
//
// Examples:
//
//	ExpectThat(t, 42, Eq(42))
//	ExpectThat(t, "hello", Eq("hello"))
//	ExpectThat(t, Name("hello"), Eq("hello"))  // same underlying type
//	ExpectThat(t, []int{1, 2, 3}, Eq([]int{1, 2, 3}))
//	ExpectThat(t, 42, Not(Eq(43)))
//	ExpectThat(t, 42, Not(43)) // same as above
//
// You can also use ExpectEq() as a shorthand:
//
//	ExpectEq(t, 42, 42)
func Eq(x any) Matcher {
	callerPkg, ok := GetCallerPkg()
	if !ok {
		panic("Eq: unable to determine caller package")
	}
	opts := []cmp.Option{
		ExportFieldsFrom(callerPkg),
		CompareProtos(),
		IgnoreHiddenFieldsExceptFrom(callerPkg),
	}
	return eqMatcher{val: x, opts: opts}
}

// Like Eq, but allows customizing the comparison behavior using cmp.Options.
//
// Note that this matcher does not automatically inject any Export or Ignore
// options, so using it with types that have unexported fields requires setting
// some options.
//
// Examples:
//
//	type Person struct { Name string; Age int }
//	p1 := Person{"Alice", 30}
//	p2 := Person{"Alice", 31}
//	ExpectThat(t, p1, Equiv(p2, cmpopts.IgnoreFields(Person{}, "Age")))
//	ExpectThat(t, p1, Not(Equiv(p2)))
func Equiv(x any, options ...cmp.Option) Matcher {
	return eqMatcher{val: x, opts: options}
}

func ExportFieldsFrom(pkg string) cmp.Option {
	return cmp.Exporter(func(t reflect.Type) bool {
		return t.PkgPath() == pkg
	})
}

func IgnoreHiddenFieldsExceptFrom(pkg string) cmp.Option {
	return cmp.FilterPath(func(p cmp.Path) bool {
		sf, ok := p.Index(-1).(cmp.StructField)
		if !ok {
			return false
		}
		r, _ := utf8.DecodeRuneInString(sf.Name())
		hidden := !unicode.IsUpper(r)

		return hidden && p.Index(-2).Type().PkgPath() != pkg
	}, cmp.Ignore())
}

func CompareProtos() cmp.Option {
	return cmp.Comparer(func(a, b proto.Message) bool {
		return proto.Equal(a, b)
	})
}

func GetCallerPkg() (string, bool) {
	// Find the caller's package by skipping past any frames in our own package
	// (e.g., when called from ExpectEq, we want the test package, not gotest)
	skip := 1
	var callerPkg string
	eqPkg := getPackageName(getCurrentPC())

	for {
		pc, _, _, ok := runtime.Caller(skip)
		if !ok {
			return "", false
		}
		callerPkg = getPackageName(pc)
		if callerPkg != eqPkg {
			return callerPkg, true
		}
		skip++
	}
}

type eqMatcher struct {
	val  any
	opts []cmp.Option
}

func (e eqMatcher) String() string {
	return fmt.Sprintf("is equal to %+v (%T)", e.val, e.val)
}

func (e eqMatcher) Matches(x any) bool {
	return cmp.Equal(x, e.val, e.opts...)
}

func (e eqMatcher) ExplainFailure(x any) (string, bool) {
	diff := cmp.Diff(e.val, x, e.opts...)
	if diff == "" {
		return "", false
	}
	return fmt.Sprintf("doesn't match (-want +got):\n%s", diff), true
}

func getCurrentPC() uintptr {
	pc, _, _, _ := runtime.Caller(1)
	return pc
}

func getPackageName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}

	name := fn.Name()

	// The function name format is: package/path.FunctionName or package/path.(*Type).MethodName
	// We need to extract the package path, which is everything before the first dot
	// after the last slash (or before the first dot if there's no slash).

	// Find the last slash in the name
	lastSlash := strings.LastIndexByte(name, '/')

	// Find the first dot after the last slash
	start := 0
	if lastSlash >= 0 {
		start = lastSlash + 1
	}

	dot := strings.IndexByte(name[start:], '.')
	if dot < 0 {
		return ""
	}

	return name[:start+dot]
}
