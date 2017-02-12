# Steps

+ Create a temporary environment

## Step implementations

Step implementations are functions defined in go code with an associated regex
which is used to match step text in the specifcation with the correct implementation.

The regex is used to identify the correct implementation and to capture any parameters
from the step text which need to be passed to it.

Implementations must be registered with the elicit context during setup.
This seems cumbersome, but the following syntax is a succinct way to write it,
keeping the regex next to the function. Of course, you're free to construct the
map in any way you see fit. They may be organised into whatever packages you like,
but it is convenient to keep them in a single package.

+ Create a `steps_test.go` file:

```go
package elicit_test

import (
    "testing"
)

func init() {
    steps[`Simple Step`] =
        func(t *testing.T) {
            t.Logf("simple step")
        }

    steps[`Step with "(.*)" parameter`] =
        func(t *testing.T, s string) {
            t.Logf("param: %s", s)
        }

    steps[`Step with an int parameter (-?\d+)`] =
        func(t *testing.T, i int) {
            t.Logf("%d", i)
        }

    steps[`(\d+) \+ (\d+) = (\d+)`] =
        func(t *testing.T, a, b, c int) {
            r := a + b
            if r != c {
                t.Errorf("expected %d + %d = %d, got %d", a, b, c, r)
            }
        }
}
```

Note that `steps` has already been defined in the `specs_test.go` file in the spec context.
If you don't have many steps, you could put them all in the same file with the test method.

+ Create a `step_execution.spec` file:

```markdown
# Step Execution

## No Parameters
+ Simple Step

## String parameters
+ Step with "hello" parameter
+ Step with "world" parameter

## Int parameters
+ Step with an int parameter 42
+ Step with an int parameter -1

## Multiple Parameters
+ 1 + 1 = 2
+ 2 + 3 = 5
+ 0 + 1 = 0
```

+ Running `go test -v` will output:

```
Step Execution
==============
FAIL: 1
PASS: 3

No Parameters
-------------
PASS:
    ✓ Simple Step

String parameters
-----------------
PASS:
    ✓ Step with "hello" parameter
    ✓ Step with "world" parameter

Int parameters
--------------
PASS:
    ✓ Step with an int parameter 42
    ✓ Step with an int parameter -1

Multiple Parameters
-------------------
FAIL:
    ✓ 1 + 1 = 2
    ✓ 2 + 3 = 5
    ✘ 0 + 1 = 0

--- FAIL: Test (0.00s)
    --- FAIL: Test/step_execution.spec/Step_Execution (0.00s)
        --- PASS: Test/step_execution.spec/Step_Execution/No_Parameters (0.00s)
            --- PASS: Test/step_execution.spec/Step_Execution/No_Parameters/#00 (0.00s)
            	steps_test.go:10: simple step
        --- PASS: Test/step_execution.spec/Step_Execution/String_parameters (0.00s)
            --- PASS: Test/step_execution.spec/Step_Execution/String_parameters/#00 (0.00s)
            	steps_test.go:15: param: hello
            --- PASS: Test/step_execution.spec/Step_Execution/String_parameters/#01 (0.00s)
            	steps_test.go:15: param: world
        --- PASS: Test/step_execution.spec/Step_Execution/Int_parameters (0.00s)
            --- PASS: Test/step_execution.spec/Step_Execution/Int_parameters/#00 (0.00s)
            	steps_test.go:20: 42
            --- PASS: Test/step_execution.spec/Step_Execution/Int_parameters/#01 (0.00s)
            	steps_test.go:20: -1
        --- FAIL: Test/step_execution.spec/Step_Execution/Multiple_Parameters (0.00s)
            --- PASS: Test/step_execution.spec/Step_Execution/Multiple_Parameters/#00 (0.00s)
            --- PASS: Test/step_execution.spec/Step_Execution/Multiple_Parameters/#01 (0.00s)
            --- FAIL: Test/step_execution.spec/Step_Execution/Multiple_Parameters/#02 (0.00s)
            	steps_test.go:27: expected 0 + 1 = 0, got 1
```


## Errors

When steps error (either as a result of a panic or a call to `t.Fail()` or `t.Errorf()`)
then the remaining steps will be skipped (but still logged).

+ Create a `failed_steps.spec` file:

```markdown
# Failing Steps

## Fail
+ This step fails
+ This step will be skipped

## Panic
+ This step panics
+ This step will be skipped
```

+ Create step definitions:

```go
steps[`This step fails`] =
    func(t *testing.T) {
        t.Fail()
    }

steps[`This step panics`] =
    func(t *testing.T) {
        s := []int{}
        s[0] = 0
    }

steps[`This step will be skipped`] =
    func(t *testing.T) {
        t.Errorf("Step should not have been called")
    }
```

+ Running `go test -v` will output:

```
Failing Steps
=============
FAIL: 1
PANIC: 1

Fail
----
FAIL:
    ✘ This step fails
    ⤹ This step will be skipped

Panic
-----
PANIC:
    ⚡ This step panics
    ⤹ This step will be skipped
```


## Skipping

When steps are skipped (either as a result of being undefined or a call to `t.Skip()`)
then the remaining steps will also be skipped (but still logged)

+ Create a `skipped_steps.spec` file:

```markdown
# Skipping Steps

## Undefined
+ This step has no implementation
+ This step will be skipped

## Skipped
+ This step skips
+ This step will be skipped
```

+ Create step definitions:

```go
steps[`This step will be skipped`] =
    func(t *testing.T) {
        t.Errorf("Step should not have been called")
    }

steps[`This step skips`] =
    func(t *testing.T) {
        t.Skip("skipping...")
    }
```

+ Running `go test -v` will output:

```
Skipping Steps
==============
SKIP: 2

Undefined
---------
SKIP:
    ? This step has no implementation
    ⤹ This step will be skipped

Skipped
-------
SKIP:
    ⤹ This step skips
    ⤹ This step will be skipped

--- SKIP: Test (0.00s)
    --- SKIP: Test/skipped_steps.spec/Skipping_Steps (0.00s)
        --- SKIP: Test/skipped_steps.spec/Skipping_Steps/Undefined (0.00s)
            --- SKIP: Test/skipped_steps.spec/Skipping_Steps/Undefined/#00 (0.00s)
            	step.go:55: no matching step implementation
            --- SKIP: Test/skipped_steps.spec/Skipping_Steps/Undefined/#01 (0.00s)
        --- SKIP: Test/skipped_steps.spec/Skipping_Steps/Skipped (0.00s)
            --- SKIP: Test/skipped_steps.spec/Skipping_Steps/Skipped/#00 (0.00s)
            	steps_test.go:17: skipping...
            --- SKIP: Test/skipped_steps.spec/Skipping_Steps/Skipped/#01 (0.00s)
```


## Forcing Step Execution

If there are steps that need to be run regardless of the success or failure of previous steps,
add emphasis to the whole step text.

+ Create a `forced_steps.spec` file:

```markdown
# Forcing Steps to Run

## Forced Step
+ This step is skipped
+ *Forced step*
```

+ Create step definitions:

```go
steps[`Forced step`] =
    func(t *testing.T) {
        t.Logf("forced step")
    }
```

+ Running `go test -v` will output:

```
Forcing Steps to Run
====================
PENDING: 1

Forced Step
-----------
PENDING:
    ? This step is skipped
    ✔ Forced step
```

---

+ *Remove the temporary directory*