# Hooks

Hooks are used to run functions before or after:

- Specs
- Scenarios
- Steps

These are similar to the scenario-level before/after steps,
but they are _not_ skipped in the event a step fails.
However, if a before hook panics, the after hooks for 
that level will not be run.

Hooks are isolated from the test context, but they can cause test
failure by panicking. The precise behaviour depends on the type
of hook as described below.

+ Create a `hooks_example.md` file:

```markdown
# Hooks Example

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

```

+ Create a `passing_spec.md` file:

```markdown
# A Passing Spec
## Another Passing Scenario
+ Another passing step
```

+ Create step definitions:

```go
steps[`(.+) passing step`] = 
    func(t *testing.T, s string) {
        t.Log(s, "step");
    }
steps[`Skipping step`] = 
    func(t *testing.T) {
        t.Skip("skipping step");
    }
steps[`Failing step`] = 
    func(t *testing.T) {
        t.Errorf("failing step");
    }
steps[`Panicking step`] = 
    func(t *testing.T) {
        panic("panicking step");
    }
steps[`This step will be skipped`] =
    func(t *testing.T) {
        t.Errorf("This step shouldn't be called");
    }
```

## Hook Examples

+ Create a `spec_hooks_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "github.com/mpwalkerdine/elicit"
    "testing"
)

var (
    steps = elicit.Steps{}
    specNum = 0
    scenarioNum = 0
    stepNum = 0
)

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeSpecs(func() {
            specNum++
            scenarioNum = 0
            stepNum = 0
            fmt.Println("\nHook: before spec", specNum)
        }).
        BeforeScenarios(func() {
            scenarioNum++
            stepNum = 0
            fmt.Println("\nHook:     before scenario", scenarioNum)
        }).
        BeforeSteps(func() {
            stepNum++
            fmt.Print("\nHook:         before step ", stepNum, " - ")
        }).
        AfterSteps(func() {
            fmt.Print("after ", stepNum)
        }).
        AfterScenarios(func() {
            fmt.Println("\nHook:     after scenario", scenarioNum)
        }).
        AfterSpecs(func() {
            fmt.Println("\nHook: after spec", specNum)
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test

Hook: before spec 1
=== RUN   Test/hooks_example.md/Hooks_Example

Hook:     before scenario 1
=== RUN   Test/hooks_example.md/Hooks_Example/Passing_Scenario

Hook:         before step 1 - after 1
Hook:         before step 2 - after 2
Hook:     after scenario 1

Hook:     before scenario 2
=== RUN   Test/hooks_example.md/Hooks_Example/Pending_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 2

Hook:     before scenario 3
=== RUN   Test/hooks_example.md/Hooks_Example/Skipping_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 3

Hook:     before scenario 4
=== RUN   Test/hooks_example.md/Hooks_Example/Failing_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 4

Hook:     before scenario 5
=== RUN   Test/hooks_example.md/Hooks_Example/Panicking_Scenario

Hook:         before step 1 - panic during step hooks_example.md/Hooks Example/Panicking Scenario/Panicking step: panicking step
after 1
Hook:     after scenario 5

Hook: after spec 1

Hook: before spec 2
=== RUN   Test/passing_spec.md/A_Passing_Spec

Hook:     before scenario 1
=== RUN   Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 1

Hook: after spec 2


Hooks Example
=============
Passed: 1
Skipped: 1
Pending: 1
Failed: 1
Panicked: 1

Passing Scenario
----------------
Passed

    ✓ First passing step
    ✓ Second passing step

Pending Scenario
----------------
Pending

    ? Undefined step
    ⤹ This step will be skipped

Skipping Scenario
-----------------
Skipped

    ⤹ Skipping step
    ⤹ This step will be skipped

Failing Scenario
----------------
Failed

    ✘ Failing step
    ⤹ This step will be skipped

Panicking Scenario
------------------
Panicked

    ⚡ Panicking step
    ⤹ This step will be skipped


A Passing Spec
==============
Passed: 1

Another Passing Scenario
------------------------
Passed

    ✓ Another passing step

--- FAIL: Test (0.00s)
    --- FAIL: Test/hooks_example.md/Hooks_Example (0.00s)
        --- PASS: Test/hooks_example.md/Hooks_Example/Passing_Scenario (0.00s)
        	steps_test.go:12: First step
        	steps_test.go:12: Second step
        --- SKIP: Test/hooks_example.md/Hooks_Example/Pending_Scenario (0.00s)
        --- SKIP: Test/hooks_example.md/Hooks_Example/Skipping_Scenario (0.00s)
        	steps_test.go:16: skipping step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Failing_Scenario (0.00s)
        	steps_test.go:20: failing step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Panicking_Scenario (0.00s)
    --- PASS: Test/passing_spec.md/A_Passing_Spec (0.00s)
        --- PASS: Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario (0.00s)
        	steps_test.go:12: Another step
```

## Before Spec Hook Panic

Panics before a spec will prevent any of the scenario tests from running. The
spec itself will run, but will be marked as failed. All scenarios are skipped,
their associated subtests are not run.

+ Create a `before_spec_panic_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "github.com/mpwalkerdine/elicit"
    "testing"
)

var (
    steps = elicit.Steps{}
    specNum = 0
    scenarioNum = 0
    stepNum = 0
)

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeSpecs(func() {
            specNum++
            scenarioNum = 0
            stepNum = 0
            fmt.Println("\nHook: before spec", specNum)
            panic(fmt.Errorf("panic before spec %d", specNum))
        }).
        BeforeScenarios(func() {
            scenarioNum++
            stepNum = 0
            fmt.Println("\nHook:     before scenario", scenarioNum)
        }).
        BeforeSteps(func() {
            stepNum++
            fmt.Print("\nHook:         before step ", stepNum, " - ")
        }).
        AfterSteps(func() {
            fmt.Print("after ", stepNum)
        }).
        AfterScenarios(func() {
            fmt.Println("\nHook:     after scenario", scenarioNum)
        }).
        AfterSpecs(func() {
            fmt.Println("\nHook: after spec", specNum)
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test

Hook: before spec 1
panic during before spec hook: panic before spec 1
=== RUN   Test/hooks_example.md/Hooks_Example

Hook: before spec 2
panic during before spec hook: panic before spec 2
=== RUN   Test/passing_spec.md/A_Passing_Spec


Hooks Example
=============
Skipped: 5

Passing Scenario
----------------
Skipped

    ⤹ First passing step
    ⤹ Second passing step

Pending Scenario
----------------
Skipped

    ? Undefined step
    ⤹ This step will be skipped

Skipping Scenario
-----------------
Skipped

    ⤹ Skipping step
    ⤹ This step will be skipped

Failing Scenario
----------------
Skipped

    ⤹ Failing step
    ⤹ This step will be skipped

Panicking Scenario
------------------
Skipped

    ⤹ Panicking step
    ⤹ This step will be skipped


A Passing Spec
==============
Skipped: 1

Another Passing Scenario
------------------------
Skipped

    ⤹ Another passing step

--- FAIL: Test (0.00s)
    --- FAIL: Test/hooks_example.md/Hooks_Example (0.00s)
    --- FAIL: Test/passing_spec.md/A_Passing_Spec (0.00s)
```


## After Spec Hook Panic

Panics after a spec will cause the test to fail, even if all scenarios passed.

+ Create a `after_spec_panic_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "github.com/mpwalkerdine/elicit"
    "testing"
)

var (
    steps = elicit.Steps{}
    specNum = 0
    scenarioNum = 0
    stepNum = 0
)

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeSpecs(func() {
            specNum++
            scenarioNum = 0
            stepNum = 0
            fmt.Println("\nHook: before spec", specNum)
        }).
        BeforeScenarios(func() {
            scenarioNum++
            stepNum = 0
            fmt.Println("\nHook:     before scenario", scenarioNum)
        }).
        BeforeSteps(func() {
            stepNum++
            fmt.Print("\nHook:         before step ", stepNum, " - ")
        }).
        AfterSteps(func() {
            fmt.Print("after ", stepNum)
        }).
        AfterScenarios(func() {
            fmt.Println("\nHook:     after scenario", scenarioNum)
        }).
        AfterSpecs(func() {
            fmt.Println("\nHook: after spec", specNum)
            panic(fmt.Errorf("panic after spec %d", specNum))
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test

Hook: before spec 1
=== RUN   Test/hooks_example.md/Hooks_Example

Hook:     before scenario 1
=== RUN   Test/hooks_example.md/Hooks_Example/Passing_Scenario

Hook:         before step 1 - after 1
Hook:         before step 2 - after 2
Hook:     after scenario 1

Hook:     before scenario 2
=== RUN   Test/hooks_example.md/Hooks_Example/Pending_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 2

Hook:     before scenario 3
=== RUN   Test/hooks_example.md/Hooks_Example/Skipping_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 3

Hook:     before scenario 4
=== RUN   Test/hooks_example.md/Hooks_Example/Failing_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 4

Hook:     before scenario 5
=== RUN   Test/hooks_example.md/Hooks_Example/Panicking_Scenario

Hook:         before step 1 - panic during step hooks_example.md/Hooks Example/Panicking Scenario/Panicking step: panicking step
after 1
Hook:     after scenario 5

Hook: after spec 1
panic during after spec hook: panic after spec 1

Hook: before spec 2
=== RUN   Test/passing_spec.md/A_Passing_Spec

Hook:     before scenario 1
=== RUN   Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 1

Hook: after spec 2
panic during after spec hook: panic after spec 2


Hooks Example
=============
Passed: 1
Skipped: 1
Pending: 1
Failed: 1
Panicked: 1

Passing Scenario
----------------
Passed

    ✓ First passing step
    ✓ Second passing step

Pending Scenario
----------------
Pending

    ? Undefined step
    ⤹ This step will be skipped

Skipping Scenario
-----------------
Skipped

    ⤹ Skipping step
    ⤹ This step will be skipped

Failing Scenario
----------------
Failed

    ✘ Failing step
    ⤹ This step will be skipped

Panicking Scenario
------------------
Panicked

    ⚡ Panicking step
    ⤹ This step will be skipped


A Passing Spec
==============
Passed: 1

Another Passing Scenario
------------------------
Passed

    ✓ Another passing step

--- FAIL: Test (0.00s)
    --- FAIL: Test/hooks_example.md/Hooks_Example (0.00s)
        --- PASS: Test/hooks_example.md/Hooks_Example/Passing_Scenario (0.00s)
        	steps_test.go:12: First step
        	steps_test.go:12: Second step
        --- SKIP: Test/hooks_example.md/Hooks_Example/Pending_Scenario (0.00s)
        --- SKIP: Test/hooks_example.md/Hooks_Example/Skipping_Scenario (0.00s)
        	steps_test.go:16: skipping step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Failing_Scenario (0.00s)
        	steps_test.go:20: failing step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Panicking_Scenario (0.00s)
    --- FAIL: Test/passing_spec.md/A_Passing_Spec (0.00s)
        --- PASS: Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario (0.00s)
        	steps_test.go:12: Another step
```


## Before Scenario Hook Panic

Panics before a scenario will cause the test to fail, all of the steps will
be skipped.

+ Create a `before_scenario_panic_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "github.com/mpwalkerdine/elicit"
    "testing"
)

var (
    steps = elicit.Steps{}
    specNum = 0
    scenarioNum = 0
    stepNum = 0
)

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeSpecs(func() {
            specNum++
            scenarioNum = 0
            stepNum = 0
            fmt.Println("\nHook: before spec", specNum)
        }).
        BeforeScenarios(func() {
            scenarioNum++
            stepNum = 0
            fmt.Println("\nHook:     before scenario", scenarioNum)
            panic(fmt.Errorf("panic before scenario %d", scenarioNum))
        }).
        BeforeSteps(func() {
            stepNum++
            fmt.Print("\nHook:         before step ", stepNum, " - ")
        }).
        AfterSteps(func() {
            fmt.Print("after ", stepNum)
        }).
        AfterScenarios(func() {
            fmt.Println("\nHook:     after scenario", scenarioNum)
        }).
        AfterSpecs(func() {
            fmt.Println("\nHook: after spec", specNum)
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test

Hook: before spec 1
=== RUN   Test/hooks_example.md/Hooks_Example

Hook:     before scenario 1
panic during before scenario hook: panic before scenario 1
=== RUN   Test/hooks_example.md/Hooks_Example/Passing_Scenario

Hook:     before scenario 2
panic during before scenario hook: panic before scenario 2
=== RUN   Test/hooks_example.md/Hooks_Example/Pending_Scenario

Hook:     before scenario 3
panic during before scenario hook: panic before scenario 3
=== RUN   Test/hooks_example.md/Hooks_Example/Skipping_Scenario

Hook:     before scenario 4
panic during before scenario hook: panic before scenario 4
=== RUN   Test/hooks_example.md/Hooks_Example/Failing_Scenario

Hook:     before scenario 5
panic during before scenario hook: panic before scenario 5
=== RUN   Test/hooks_example.md/Hooks_Example/Panicking_Scenario

Hook: after spec 1

Hook: before spec 2
=== RUN   Test/passing_spec.md/A_Passing_Spec

Hook:     before scenario 1
panic during before scenario hook: panic before scenario 1
=== RUN   Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario

Hook: after spec 2


Hooks Example
=============
Panicked: 5

Passing Scenario
----------------
Panicked

    ⤹ First passing step
    ⤹ Second passing step

Pending Scenario
----------------
Panicked

    ? Undefined step
    ⤹ This step will be skipped

Skipping Scenario
-----------------
Panicked

    ⤹ Skipping step
    ⤹ This step will be skipped

Failing Scenario
----------------
Panicked

    ⤹ Failing step
    ⤹ This step will be skipped

Panicking Scenario
------------------
Panicked

    ⤹ Panicking step
    ⤹ This step will be skipped


A Passing Spec
==============
Panicked: 1

Another Passing Scenario
------------------------
Panicked

    ⤹ Another passing step

--- FAIL: Test (0.00s)
    --- FAIL: Test/hooks_example.md/Hooks_Example (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Passing_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Pending_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Skipping_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Failing_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Panicking_Scenario (0.00s)
    --- FAIL: Test/passing_spec.md/A_Passing_Spec (0.00s)
        --- FAIL: Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario (0.00s)
```


## After Scenario Hook Panic

Panics after a scenario will cause the test to fail, even if all the steps
ran successfully.

+ Create a `after_scenario_panic_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "github.com/mpwalkerdine/elicit"
    "testing"
)

var (
    steps = elicit.Steps{}
    specNum = 0
    scenarioNum = 0
    stepNum = 0
)

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeSpecs(func() {
            specNum++
            scenarioNum = 0
            stepNum = 0
            fmt.Println("\nHook: before spec", specNum)
        }).
        BeforeScenarios(func() {
            scenarioNum++
            stepNum = 0
            fmt.Println("\nHook:     before scenario", scenarioNum)
        }).
        BeforeSteps(func() {
            stepNum++
            fmt.Print("\nHook:         before step ", stepNum, " - ")
        }).
        AfterSteps(func() {
            fmt.Print("after ", stepNum)
        }).
        AfterScenarios(func() {
            fmt.Println("\nHook:     after scenario", scenarioNum)
            panic(fmt.Errorf("panic after scenario %d", scenarioNum))
        }).
        AfterSpecs(func() {
            fmt.Println("\nHook: after spec", specNum)
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test

Hook: before spec 1
=== RUN   Test/hooks_example.md/Hooks_Example

Hook:     before scenario 1
=== RUN   Test/hooks_example.md/Hooks_Example/Passing_Scenario

Hook:         before step 1 - after 1
Hook:         before step 2 - after 2
Hook:     after scenario 1
panic during after scenario hook: panic after scenario 1

Hook:     before scenario 2
=== RUN   Test/hooks_example.md/Hooks_Example/Pending_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 2
panic during after scenario hook: panic after scenario 2

Hook:     before scenario 3
=== RUN   Test/hooks_example.md/Hooks_Example/Skipping_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 3
panic during after scenario hook: panic after scenario 3

Hook:     before scenario 4
=== RUN   Test/hooks_example.md/Hooks_Example/Failing_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 4
panic during after scenario hook: panic after scenario 4

Hook:     before scenario 5
=== RUN   Test/hooks_example.md/Hooks_Example/Panicking_Scenario

Hook:         before step 1 - panic during step hooks_example.md/Hooks Example/Panicking Scenario/Panicking step: panicking step
after 1
Hook:     after scenario 5
panic during after scenario hook: panic after scenario 5

Hook: after spec 1

Hook: before spec 2
=== RUN   Test/passing_spec.md/A_Passing_Spec

Hook:     before scenario 1
=== RUN   Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario

Hook:         before step 1 - after 1
Hook:     after scenario 1
panic during after scenario hook: panic after scenario 1

Hook: after spec 2


Hooks Example
=============
Panicked: 5

Passing Scenario
----------------
Panicked

    ✓ First passing step
    ✓ Second passing step

Pending Scenario
----------------
Panicked

    ? Undefined step
    ⤹ This step will be skipped

Skipping Scenario
-----------------
Panicked

    ⤹ Skipping step
    ⤹ This step will be skipped

Failing Scenario
----------------
Panicked

    ✘ Failing step
    ⤹ This step will be skipped

Panicking Scenario
------------------
Panicked

    ⚡ Panicking step
    ⤹ This step will be skipped


A Passing Spec
==============
Panicked: 1

Another Passing Scenario
------------------------
Panicked

    ✓ Another passing step

--- FAIL: Test (0.00s)
    --- FAIL: Test/hooks_example.md/Hooks_Example (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Passing_Scenario (0.00s)
        	steps_test.go:12: First step
        	steps_test.go:12: Second step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Pending_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Skipping_Scenario (0.00s)
        	steps_test.go:16: skipping step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Failing_Scenario (0.00s)
        	steps_test.go:20: failing step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Panicking_Scenario (0.00s)
    --- FAIL: Test/passing_spec.md/A_Passing_Spec (0.00s)
        --- FAIL: Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario (0.00s)
        	steps_test.go:12: Another step
```


## Before Step Hook Panic

Panics before a step will cause the step to fail, as if the step itself had
panicked.

+ Create a `before_step_panic_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "github.com/mpwalkerdine/elicit"
    "testing"
)

var (
    steps = elicit.Steps{}
    specNum = 0
    scenarioNum = 0
    stepNum = 0
)

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeSpecs(func() {
            specNum++
            scenarioNum = 0
            stepNum = 0
            fmt.Println("\nHook: before spec", specNum)
        }).
        BeforeScenarios(func() {
            scenarioNum++
            stepNum = 0
            fmt.Println("\nHook:     before scenario", scenarioNum)
        }).
        BeforeSteps(func() {
            stepNum++
            fmt.Print("\nHook:         before step ", stepNum, " - ")
            panic(fmt.Errorf("panic before step %d", stepNum))
        }).
        AfterSteps(func() {
            fmt.Print("after ", stepNum)
        }).
        AfterScenarios(func() {
            fmt.Println("\nHook:     after scenario", scenarioNum)
        }).
        AfterSpecs(func() {
            fmt.Println("\nHook: after spec", specNum)
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test

Hook: before spec 1
=== RUN   Test/hooks_example.md/Hooks_Example

Hook:     before scenario 1
=== RUN   Test/hooks_example.md/Hooks_Example/Passing_Scenario

Hook:         before step 1 - panic during before step hook: panic before step 1

Hook:     after scenario 1

Hook:     before scenario 2
=== RUN   Test/hooks_example.md/Hooks_Example/Pending_Scenario

Hook:         before step 1 - panic during before step hook: panic before step 1

Hook:     after scenario 2

Hook:     before scenario 3
=== RUN   Test/hooks_example.md/Hooks_Example/Skipping_Scenario

Hook:         before step 1 - panic during before step hook: panic before step 1

Hook:     after scenario 3

Hook:     before scenario 4
=== RUN   Test/hooks_example.md/Hooks_Example/Failing_Scenario

Hook:         before step 1 - panic during before step hook: panic before step 1

Hook:     after scenario 4

Hook:     before scenario 5
=== RUN   Test/hooks_example.md/Hooks_Example/Panicking_Scenario

Hook:         before step 1 - panic during before step hook: panic before step 1

Hook:     after scenario 5

Hook: after spec 1

Hook: before spec 2
=== RUN   Test/passing_spec.md/A_Passing_Spec

Hook:     before scenario 1
=== RUN   Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario

Hook:         before step 1 - panic during before step hook: panic before step 1

Hook:     after scenario 1

Hook: after spec 2


Hooks Example
=============
Panicked: 5

Passing Scenario
----------------
Panicked

    ⚡ First passing step
    ⤹ Second passing step

Pending Scenario
----------------
Panicked

    ⚡ Undefined step
    ⤹ This step will be skipped

Skipping Scenario
-----------------
Panicked

    ⚡ Skipping step
    ⤹ This step will be skipped

Failing Scenario
----------------
Panicked

    ⚡ Failing step
    ⤹ This step will be skipped

Panicking Scenario
------------------
Panicked

    ⚡ Panicking step
    ⤹ This step will be skipped


A Passing Spec
==============
Panicked: 1

Another Passing Scenario
------------------------
Panicked

    ⚡ Another passing step

--- FAIL: Test (0.00s)
    --- FAIL: Test/hooks_example.md/Hooks_Example (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Passing_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Pending_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Skipping_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Failing_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Panicking_Scenario (0.00s)
    --- FAIL: Test/passing_spec.md/A_Passing_Spec (0.00s)
        --- FAIL: Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario (0.00s)
```


## After Step Hook Panic

Panics after a step will cause the step to fail, as if the step itself had
panicked.

+ Create a `after_step_panic_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "github.com/mpwalkerdine/elicit"
    "testing"
)

var (
    steps = elicit.Steps{}
    specNum = 0
    scenarioNum = 0
    stepNum = 0
)

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        BeforeSpecs(func() {
            specNum++
            scenarioNum = 0
            stepNum = 0
            fmt.Println("\nHook: before spec", specNum)
        }).
        BeforeScenarios(func() {
            scenarioNum++
            stepNum = 0
            fmt.Println("\nHook:     before scenario", scenarioNum)
        }).
        BeforeSteps(func() {
            stepNum++
            fmt.Print("\nHook:         before step ", stepNum, " - ")
        }).
        AfterSteps(func() {
            fmt.Print("after ", stepNum, " - ")
            panic(fmt.Errorf("panic after step %d", stepNum))
        }).
        AfterScenarios(func() {
            fmt.Println("\nHook:     after scenario", scenarioNum)
        }).
        AfterSpecs(func() {
            fmt.Println("\nHook: after spec", specNum)
        }).
        RunTests(t)
}
```

+ Running `go test -v` will output:

```
=== RUN   Test

Hook: before spec 1
=== RUN   Test/hooks_example.md/Hooks_Example

Hook:     before scenario 1
=== RUN   Test/hooks_example.md/Hooks_Example/Passing_Scenario

Hook:         before step 1 - after 1 - panic during after step hook: panic after step 1

Hook:     after scenario 1

Hook:     before scenario 2
=== RUN   Test/hooks_example.md/Hooks_Example/Pending_Scenario

Hook:         before step 1 - after 1 - panic during after step hook: panic after step 1

Hook:     after scenario 2

Hook:     before scenario 3
=== RUN   Test/hooks_example.md/Hooks_Example/Skipping_Scenario

Hook:         before step 1 - after 1 - panic during after step hook: panic after step 1

Hook:     after scenario 3

Hook:     before scenario 4
=== RUN   Test/hooks_example.md/Hooks_Example/Failing_Scenario

Hook:         before step 1 - after 1 - panic during after step hook: panic after step 1

Hook:     after scenario 4

Hook:     before scenario 5
=== RUN   Test/hooks_example.md/Hooks_Example/Panicking_Scenario

Hook:         before step 1 - panic during step hooks_example.md/Hooks Example/Panicking Scenario/Panicking step: panicking step
after 1 - panic during after step hook: panic after step 1

Hook:     after scenario 5

Hook: after spec 1

Hook: before spec 2
=== RUN   Test/passing_spec.md/A_Passing_Spec

Hook:     before scenario 1
=== RUN   Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario

Hook:         before step 1 - after 1 - panic during after step hook: panic after step 1

Hook:     after scenario 1

Hook: after spec 2


Hooks Example
=============
Panicked: 5

Passing Scenario
----------------
Panicked

    ⚡ First passing step
    ⤹ Second passing step

Pending Scenario
----------------
Panicked

    ⚡ Undefined step
    ⤹ This step will be skipped

Skipping Scenario
-----------------
Panicked

    ⚡ Skipping step
    ⤹ This step will be skipped

Failing Scenario
----------------
Panicked

    ⚡ Failing step
    ⤹ This step will be skipped

Panicking Scenario
------------------
Panicked

    ⚡ Panicking step
    ⤹ This step will be skipped


A Passing Spec
==============
Panicked: 1

Another Passing Scenario
------------------------
Panicked

    ⚡ Another passing step

--- FAIL: Test (0.00s)
    --- FAIL: Test/hooks_example.md/Hooks_Example (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Passing_Scenario (0.00s)
        	steps_test.go:12: First step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Pending_Scenario (0.00s)
        --- FAIL: Test/hooks_example.md/Hooks_Example/Skipping_Scenario (0.00s)
        	steps_test.go:16: skipping step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Failing_Scenario (0.00s)
        	steps_test.go:20: failing step
        --- FAIL: Test/hooks_example.md/Hooks_Example/Panicking_Scenario (0.00s)
    --- FAIL: Test/passing_spec.md/A_Passing_Spec (0.00s)
        --- FAIL: Test/passing_spec.md/A_Passing_Spec/Another_Passing_Scenario (0.00s)
        	steps_test.go:12: Another step
```