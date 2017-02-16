# Logging

There are two ways of logging output:
1. Using the testing.T log API.
2. Writing to `os.Stdout` e.g. with the `"fmt"` package.

The first behaves as expected, with output written beneath the subtest result in
`go test`. Only failed tests are shown by default, but `go test -v` will show
all tests and output.

The second mechanism captures output for use in the specification execution
report, which is displayed above the `go test` detailed results. Output written
to `os.Stdout` will be displayed beneath the step text in the report. Again,
only failed tests appear in the report by default, unless `-v` is specified.

The report may optionally be written to file specified by the `-elicit.report`
flag. In this case, all results are written, regardless of `-v`.

+ Create a temporary environment

## Go Test Error Logs

This example demonstrates the `testing.T` logging API output.

+ Create a `logging_test.md` file:

```markdown
# Logging Test
## Logging Scenario
+ Logged step
```

+ Create a step definition:

```go
steps[`Logged step`] = func(t *testing.T) {
    t.Log("Logged output")
}
```

+ Running `go test -v` will output:

```
Logging Test
============
Passed: 1

Logging Scenario
----------------
Passed

    ✓ Logged step

--- PASS: Test (0.00s)
    --- PASS: Test/logging_test.md/Logging_Test (0.00s)
        --- PASS: Test/logging_test.md/Logging_Test/Logging_Scenario (0.00s)
        	steps_test.go:11: Logged output
```

## Captured Output

This example demonstrates output captured from `os.Stdout`.

+ Create a `logging_test.md` file:

```markdown
# Logging Test
## Logging Scenario
+ Logged step
```

+ Create a step definition using "fmt":

```go
steps[`Logged step`] = func(t *testing.T) {
    fmt.Println("Logged output")
}
```

+ Running `go test -v` will output:

```
Logging Test
============
Passed: 1

Logging Scenario
----------------
Passed

    ✓ Logged step
        Logged output

--- PASS: Test (0.00s)
    --- PASS: Test/logging_test.md/Logging_Test (0.00s)
        --- PASS: Test/logging_test.md/Logging_Test/Logging_Scenario (0.00s)
```

## Normal vs Chatty vs File Output

This example demonstrates the effect of the `-v` and `-elicit.report` flags on
the output.

+ Create a `logging_test.md` file:

```markdown
# Logging Test
## Undefined
+ Undefined step
## Skipped
+ Skipped step
## Failed
+ Failed step
## Panic
+ Panicked step
## Pass
+ Passing step

# Passing Spec
## Passing Scenario
+ Passing step
+ Passing step
## Another Passing Scenario
+ Passing step
+ Passing step
```

+ Create step definitions using "fmt":

```go
steps[`Skipped step`] = func(t *testing.T) {
    fmt.Println("Skipped stdout output")
    t.Skip("Skipped test output")
}
steps[`Failed step`] = func(t *testing.T) {
    fmt.Println("Failed stdout output")
    t.Errorf("Failed test output")
}
steps[`Panicked step`] = func(t *testing.T) {
    fmt.Println("Panicked stdout output")
    panic(fmt.Errorf("Panicked output"))
}
steps[`Passing step`] = func(t *testing.T) {
    fmt.Println("Passing stdout output")
    t.Log("Passing test output")
}
```

+ Running `go test` will output:

```
Logging Test
============
Passed: 1
Skipped: 1
Pending: 1
Failed: 1
Panicked: 1

Undefined
---------
Pending

    ? Undefined step

Failed
------
Failed

    ✘ Failed step
        Failed stdout output

Panic
-----
Panicked

    ⚡ Panicked step
        Panicked stdout output

--- FAIL: Test (0.00s)
    --- FAIL: Test/logging_test.md/Logging_Test (0.00s)
        --- FAIL: Test/logging_test.md/Logging_Test/Failed (0.00s)
        	steps_test.go:16: Failed test output
        --- FAIL: Test/logging_test.md/Logging_Test/Panic (0.00s)
        	step.go:40: Panicked output
```

+ Running `go test -v` will output:

```
Logging Test
============
Passed: 1
Skipped: 1
Pending: 1
Failed: 1
Panicked: 1

Undefined
---------
Pending

    ? Undefined step

Skipped
-------
Skipped

    ⤹ Skipped step
        Skipped stdout output

Failed
------
Failed

    ✘ Failed step
        Failed stdout output

Panic
-----
Panicked

    ⚡ Panicked step
        Panicked stdout output

Pass
----
Passed

    ✓ Passing step
        Passing stdout output


Passing Spec
============
Passed: 2

Passing Scenario
----------------
Passed

    ✓ Passing step
        Passing stdout output
    ✓ Passing step
        Passing stdout output

Another Passing Scenario
------------------------
Passed

    ✓ Passing step
        Passing stdout output
    ✓ Passing step
        Passing stdout output

--- FAIL: Test (0.00s)
    --- FAIL: Test/logging_test.md/Logging_Test (0.00s)
        --- SKIP: Test/logging_test.md/Logging_Test/Undefined (0.00s)
        --- SKIP: Test/logging_test.md/Logging_Test/Skipped (0.00s)
        	steps_test.go:12: Skipped test output
        --- FAIL: Test/logging_test.md/Logging_Test/Failed (0.00s)
        	steps_test.go:16: Failed test output
        --- FAIL: Test/logging_test.md/Logging_Test/Panic (0.00s)
        	step.go:40: Panicked output
        --- PASS: Test/logging_test.md/Logging_Test/Pass (0.00s)
        	steps_test.go:24: Passing test output
    --- PASS: Test/logging_test.md/Passing_Spec (0.00s)
        --- PASS: Test/logging_test.md/Passing_Spec/Passing_Scenario (0.00s)
        	steps_test.go:24: Passing test output
        	steps_test.go:24: Passing test output
        --- PASS: Test/logging_test.md/Passing_Spec/Another_Passing_Scenario (0.00s)
        	steps_test.go:24: Passing test output
        	steps_test.go:24: Passing test output
```

+ Running `go test -elicit.report ./report.md` will output:

```
Logging Test
============
Passed: 1
Skipped: 1
Pending: 1
Failed: 1
Panicked: 1

Undefined
---------
Pending

    ? Undefined step

Failed
------
Failed

    ✘ Failed step
        Failed stdout output

Panic
-----
Panicked

    ⚡ Panicked step
        Panicked stdout output

--- FAIL: Test (0.00s)
    --- FAIL: Test/logging_test.md/Logging_Test (0.00s)
        --- FAIL: Test/logging_test.md/Logging_Test/Failed (0.00s)
        	steps_test.go:16: Failed test output
        --- FAIL: Test/logging_test.md/Logging_Test/Panic (0.00s)
        	step.go:40: Panicked output
```

+ `./report.md` will contain:

```markdown
Logging Test
============
Passed: 1
Skipped: 1
Pending: 1
Failed: 1
Panicked: 1

Undefined
---------
Pending

    ? Undefined step

Skipped
-------
Skipped

    ⤹ Skipped step
        Skipped stdout output

Failed
------
Failed

    ✘ Failed step
        Failed stdout output

Panic
-----
Panicked

    ⚡ Panicked step
        Panicked stdout output

Pass
----
Passed

    ✓ Passing step
        Passing stdout output


Passing Spec
============
Passed: 2

Passing Scenario
----------------
Passed

    ✓ Passing step
        Passing stdout output
    ✓ Passing step
        Passing stdout output

Another Passing Scenario
------------------------
Passed

    ✓ Passing step
        Passing stdout output
    ✓ Passing step
        Passing stdout output
```
