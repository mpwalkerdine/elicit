# Transforms

Captured parameters are automatically converted for the following types:

    string   
    int
    []string
    []int

The slice values should be comma-separated, see [Slices](#Slices)

Arbitrary types can be converted by supplying additional `StepArgumentTransforms` during setup.

```go
// StepArgumentTransform transforms captured groups in the step pattern to a function parameter type
// Note that if the actual string cannot be converted to the target type by the transform, it should return false
type StepArgumentTransform func(string, reflect.Type) (interface{}, bool)
```

The following scenarios use the same temporary environment (see [Specification Syntax](syntax.spec))

- Create a temporary environment

## Simple Types

- Create a `simple_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "reflect"
    "testing"
)

type CustomString string

func init() {
    steps[`Step with a CustomString "(.*)"`] =
        func(t *testing.T, c CustomString) {
            fmt.Println(c)
        }

    transforms[`.*`] =
        func(params []string, target reflect.Type) (interface{}, bool) {
            if target != reflect.TypeOf((*CustomString)(nil)).Elem() {
                return nil, false
            }

            return CustomString(params[0]), true
        }
}
```

- Create a `simple_types.spec` file:

```markdown
# Simple Type Transforms
## Renamed string
- Step with a CustomString "param"
```

- Running `go test` will output:

```markdown
param

Simple Type Transforms
======================

Renamed string
--------------
  ✓ Step with a CustomString "param"
```

## Slices

- Create a `slice_test.go` file:

```go
package elicit_test

import (
    "testing"
)

func init() {
    steps[`Sum of (.+) is (.+)`] =
        func(t *testing.T, ns []int, s int) {
            actual := 0
            for _, n := range ns {
                actual += n
            }

            if actual != s {
                t.Errorf("expected sum of %v to be %d, got %d", ns, s, actual)
            }
        }
}
```

- Create a `sum.spec` file:

```markdown
# List Summation
## First Four Numbers
- Sum of 1,2,3,4 is 10
```

- Running `go test` will output:

```markdown
List Summation
==============

First Four Numbers
------------------
  ✓ Sum of 1,2,3,4 is 10
```

## Structs

- Create a `struct_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "reflect"
    "strconv"
    "testing"
)

type DOB [3]int

type Person struct {
    name string
    dob DOB
}

func init() {
    steps[`Print (.*)`] =
        func(t *testing.T, p Person) {
            fmt.Printf("\nName: %s\nDOB: %d-%d-%d\n", p.name, p.dob[0], p.dob[1], p.dob[2]);
        }

    transforms[`a person named (.*), born (\d{4})-(\d{2})-(\d{2})`] =
        func(params []string, target reflect.Type) (interface{}, bool) {
            if target != reflect.TypeOf((*Person)(nil)).Elem() {
                return nil, false
            }

            n := params[1]
            y, _ := strconv.Atoi(params[2])
            m, _ := strconv.Atoi(params[3])
            d, _ := strconv.Atoi(params[4])

            return Person{name: n, dob: [3]int{y,m,d}}, true
        }
}
```

- Create a `simple_types.spec` file:

```markdown
# Struct Transforms
## A Person
- Print a person named Bob, born 1986-01-01
```

- Running `go test` will output:

```markdown
Name: Bob
DOB: 1986-1-1

Struct Transforms
=================

A Person
--------
  ✓ Print a person named Bob, born 1986-01-01
```

---

- *Remove the temporary directory*