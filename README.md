

# Elicit [![CircleCI](https://img.shields.io/circleci/project/github/mpwalkerdine/elicit.svg)](https://circleci.com/gh/mpwalkerdine/elicit) [![Go Report Card](https://goreportcard.com/badge/github.com/mpwalkerdine/elicit)](https://goreportcard.com/report/github.com/mpwalkerdine/elicit) [![Godoc](https://img.shields.io/badge/godoc-elicit-blue.svg)](https://godoc.org/github.com/mpwalkerdine/elicit)

Elicit is a [Specification by Example] framework for [Go] inspired by similar
frameworks such as [Cucumber] and [Gauge].

**Note: Go already has excellent support for [testing] and [documentation] in
its standard library, so you might not need this framework**.

That said, if you're using Go as a general-purpose programming language and find
yourself needing close collaboration with non-technical stakeholders, then the
skills and knowledge required to use Go's existing tools may be a barrier to
this. Elicit provides the ability to write executable specifications in plain
English, which then form part of the documentation and regression test suite
for the system.

Features and aspirations:
  1. Allows executable specifications to be written in markdown.
  1. Follows Go principles, especially in terms of simplicity.
  1. Uses the Go testing API wherever possible.
  1. Doesn't impose restrictions on repository layout.
  1. Requires minimal configuration and no additional dot files.
  1. Integrates with the `go` toolchain, `go test` in particular.
     No additional tooling is required.

## Getting started

1. If you don't already have one, create a project in your [workspace].

2. Create a specification, following the guidance found in the 
   [syntax specification](./specs/syntax.md).
   It's convenient to place these in a folder somewhere, we'll assume `./specs`.

3. Create a test file, e.g. `my_test.go`. Note the `_test.go` suffix which is a 
   Go [testing] convention. You only need one of these because it will run all
   your specifications as [subtests].

    ```go
    package mypackage

    import (
        "github.com/mpwalkerdine/elicit"
        "testing"
    )

    func Test(t *testing.T) {
        elicit.New().
            WithSpecsFolder("./specs").
            WithSteps(steps).
            RunTests(t)
    }

    var steps = elicit.Steps{}
    ```

4. Run `go test`. All the steps will show as "Pending".
5. Provide some step implementations (see [Steps](./specs/steps.md)).
6. Run `go test -v` to see a complete report (see [Logging](./specs/logging.md)
   for more details on `-v` and related options).


## More Information

The [specifications](./specs) for Elicit contain more details. They also serve
as a demonstration of its capabilities, because they are also tests for the
framework itself.

- [Transforms](./specs/transforms.md):
  Use arbitrary types as parameters in step implementations.
- [Hooks](./specs/hooks.md):
  Register functions to run at particular points in the test cycle.

## Dependencies

[Blackfriday] is used for markdown parsing.

[Specification by Example]:
https://www.manning.com/books/specification-by-example

[Go]:
http://golang.org/

[Cucumber]:
https://github.com/DATA-DOG/godog

[Gauge]:
http://getgauge.io/

[workspace]:
https://golang.org/doc/code.html#Workspaces

[testing]:
https://golang.org/pkg/testing/

[subtests]:
https://golang.org/pkg/testing/#hdr-Subtests_and_Sub_benchmarks

[documentation]:
https://blog.golang.org/godoc-documenting-go-code

[Blackfriday]:
https://github.com/russross/blackfriday