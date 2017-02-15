# Hooks

Hooks are used to run functions before or after:

- Specs
- Scenarios
- Steps

They are isolated from the test context, but they can cause test
failure by panicking. The precise behaviour depends on the type of hook.

+ Create a `hooks_example.spec` file:

```markdown
# Hooks Example
+ Before step

## Passing Scenario
+ First passing step
+ Second passing step

## Pending Scenario
+ Undefined step
+ This step will be skipped

## Skipping Scenario
+ Skipping step
+ This step will be skipped

## Failing Scenario
+ Failing step
+ This step will be skipped

## Panicking Scenario
+ Panicking step
+ This step will be skipped

---

+ After step
```

+ Create step definitions:

```go
steps[`(Before|After) step`] = func(t *testing.T, s string) { t.Log(s, "step"); }
steps[`(.+) passing step`] = func(t *testing.T, s string) { t.Log(s, "step"); }
steps[`Skipping step`] = func(t *testing.T) { t.Skip("skipping step"); }
steps[`Failing step`] = func(t *testing.T) { t.Errorf("failing step"); }
steps[`Panicking step`] = func(t *testing.T) { panic("panicking step"); }
steps[`This step will be skipped`] = func(t *testing.T) { t.Errorf("This step shouldn't be called"); }
```

## Spec Hooks

+ Create a `spec_hooks_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "mmatt/elicit"
    "testing"
)

var steps = elicit.Steps{}

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeSpecs(func() {
            fmt.Println("Before hook")
        }).
        AfterSpecs(func() {
            fmt.Println("After hook")
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test
Before hook
=== RUN   Test/hooks_example.spec/Hooks_Example
=== RUN   Test/hooks_example.spec/Hooks_Example/Passing_Scenario
=== RUN   Test/hooks_example.spec/Hooks_Example/Pending_Scenario
=== RUN   Test/hooks_example.spec/Hooks_Example/Skipping_Scenario
=== RUN   Test/hooks_example.spec/Hooks_Example/Failing_Scenario
=== RUN   Test/hooks_example.spec/Hooks_Example/Panicking_Scenario
After hook

```

## Spec Hook Panics

+ TODO

## Scenario Hooks

+ Create a `scenario_hooks_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "mmatt/elicit"
    "testing"
)

var steps = elicit.Steps{}

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeScenarios(func() {
            fmt.Println("Before hook")
        }).
        AfterScenarios(func() {
            fmt.Println("After hook")
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test
=== RUN   Test/hooks_example.spec/Hooks_Example
Before hook
=== RUN   Test/hooks_example.spec/Hooks_Example/Passing_Scenario
After hook
Before hook
=== RUN   Test/hooks_example.spec/Hooks_Example/Pending_Scenario
After hook
Before hook
=== RUN   Test/hooks_example.spec/Hooks_Example/Skipping_Scenario
After hook
Before hook
=== RUN   Test/hooks_example.spec/Hooks_Example/Failing_Scenario
After hook
Before hook
=== RUN   Test/hooks_example.spec/Hooks_Example/Panicking_Scenario
After hook

```

## Scenario Hook Panics

+ TODO

## Step Hooks

+ Create a `step_hooks_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "mmatt/elicit"
    "testing"
)

var steps = elicit.Steps{}

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeSteps(func() {
            fmt.Println("Before hook")
        }).
        AfterSteps(func() {
            fmt.Println("After hook")
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test
=== RUN   Test/hooks_example.spec/Hooks_Example
=== RUN   Test/hooks_example.spec/Hooks_Example/Passing_Scenario
Before hook
After hook
Before hook
After hook
Before hook
After hook
Before hook
After hook
=== RUN   Test/hooks_example.spec/Hooks_Example/Pending_Scenario
Before hook
After hook
Before hook
=== RUN   Test/hooks_example.spec/Hooks_Example/Skipping_Scenario
Before hook
After hook
Before hook
=== RUN   Test/hooks_example.spec/Hooks_Example/Failing_Scenario
Before hook
After hook
Before hook
After hook
Before hook
=== RUN   Test/hooks_example.spec/Hooks_Example/Panicking_Scenario
Before hook
After hook
Before hook
After hook
Before hook

```

## Step Hook Panics

+ TODO