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
