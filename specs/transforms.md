# Transforms

Captured parameters are automatically converted for the following types:

- `string`
- `int`
- `[]string`
- `[]int`

The slice values should be comma-separated, see [Slices](#Slices)

Arbitrary types can be converted by supplying additional transforms during
setup. A transform is function of the form `func([]string) <type>`, where
the `[]string` parameter is the result of the pattern matching and `<type>`
in the target type.

+ Create a temporary environment

## Simple Types

+ Create a `simple_test.go` file:

```go
package elicit_test

import (
    "testing"
)

type CustomString string

func init() {
    steps[`Step with a CustomString "(.*)"`] =
        func(t *testing.T, c CustomString) {
            t.Log(c)
        }

    transforms[`.*`] =
        func(params []string) CustomString {
            return CustomString(params[0])
        }
}
```

+ Create a `simple_types.md` file:

```markdown
# Simple Type Transforms
## Renamed string
+ Step with a CustomString "param"
```

+ Running `go test -v` will output:

```
Simple Type Transforms
======================
Passed: 1

Renamed string
--------------
Passed

    ✓ Step with a CustomString "param"

--- PASS: Test (0.00s)
    --- PASS: Test/simple_types.md/Simple_Type_Transforms (0.00s)
        --- PASS: Test/simple_types.md/Simple_Type_Transforms/Renamed_string (0.00s)
        	simple_test.go:12: param
```

## Slices

+ Create a step definition:

```go
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
```

+ Create a `sum.md` file:

```markdown
# List Summation
## First Four Numbers
+ Sum of 1,2,3,4 is 10
```

+ Running `go test -v` will output:

```
List Summation
==============
Passed: 1

First Four Numbers
------------------
Passed

    ✓ Sum of 1,2,3,4 is 10
```

## Structs

+ Create a `struct_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "strconv"
    "testing"
)

type DOB [3]int

func (d DOB) String() string {
    return fmt.Sprintf("%d-%d-%d", d[0], d[1], d[2])
}

type Person struct {
    name string
    dob DOB
}

func (p Person) String() string {
    return fmt.Sprintf("Name: %s\nDOB: %s", p.name, p.dob)
}

func init() {
    steps[`Print (.*)`] =
        func(t *testing.T, p Person) {
            t.Logf("\n%s", p)
        }

    transforms[`a person named (.*), born (\d{4})-(\d{2})-(\d{2})`] =
        func(params []string) Person {
            n := params[1]
            y, _ := strconv.Atoi(params[2])
            m, _ := strconv.Atoi(params[3])
            d, _ := strconv.Atoi(params[4])

            return Person{name: n, dob: [3]int{y,m,d}}
        }
}
```

+ Create a `simple_types.md` file:

```markdown
# Struct Transforms
## A Person
+ Print a person named Bob, born 1987-01-01
```

+ Running `go test -v` will output:

```
Struct Transforms
=================
Passed: 1

A Person
--------
Passed

    ✓ Print a person named Bob, born 1987-01-01

--- PASS: Test (0.00s)
    --- PASS: Test/simple_types.md/Struct_Transforms (0.00s)
        --- PASS: Test/simple_types.md/Struct_Transforms/A_Person (0.00s)
        	struct_test.go:27: 
        		Name: Bob
        		DOB: 1987-1-1
```
