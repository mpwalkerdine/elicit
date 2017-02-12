# Validation

The context runs some checks during execution to identify potential problems
with the setup. For example, if there are no steps, or the regex for a supplied step implementation 
will never be matched, or if no transform exists to satisfy one of the parameters.

In these cases a warning will be printed on stderr.

+ Create a temporary environment

## No Specs

+ Running `go test` will output the following lines:

```
warning: No specifications found. Add a folder containing *.spec files with Context.WithSpecsFolder().
```

## No Steps

+ Running `go test` will output the following lines:

```
warning: No steps registered. Add some with Context.WithSteps().
```

## Invalid Step Patterns

+ Create step definitions:

```go
steps[`No params`] = func() {}
steps[`Invalid first param`] = func(s string) {}
steps[`Extra (param)`] = func (t *testing.T) {}
steps[`Fewer params`] = func(t *testing.T, s string) {}
type custom string
steps[`Uncovertible (param)`] = func(t *testing.T, c custom) {}
```

+ Running `go test` will output the following lines:

```
warning: The step pattern "No params" has an invalid implementation. The first parameter must be of type *testing.T.
warning: The step pattern "Invalid first param" has an invalid implementation. The first parameter must be of type *testing.T.
warning: The step pattern "Extra (param)" captures 1 parameter but the supplied implementation takes 0.
warning: The step pattern "Fewer params" captures 0 parameters but the supplied implementation takes 1.
warning: The step pattern "Uncovertible (param)" has a parameter type "custom" for which no transforms exist.
```

## Unused Steps

If the `elicit.showOrphans` flag is set, any unused step implementations
will be logged to the console.

+ TODO

---

+ *Remove the temporary directory*