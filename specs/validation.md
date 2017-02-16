# Validation

The context runs some checks during execution to identify potential problems
with the setup. For example, if there are no steps, or the regex for a supplied step implementation 
will never be matched, or if no transform exists to satisfy one of the parameters.

In these cases a warning will be printed on stderr.

+ Create a temporary environment

## No Specs

+ Running `go test` will output the following lines:

```
warning: No specifications found. Add a folder containing *.md files with Context.WithSpecsFolder().
```

## No Steps

+ Running `go test` will output the following lines:

```
warning: No steps registered. Add some with Context.WithSteps().
```

## Invalid Step Implementations

+ Create step definitions:

```go
type custom string
steps[`Not a function`] = 0
steps[`bad (regex`] = func() {}
steps[`No params`] = func() {}
steps[`Invalid first param`] = func(s string) {}
steps[`Extra (param)`] = func (t *testing.T) {}
steps[`Fewer params`] = func(t *testing.T, s string) {}
steps[`Unconvertible (param)`] = func(t *testing.T, c custom) {}
```

+ Running `go test` will output the following lines:

```
warning: registered step "Not a function" => [int] must be a function.
warning: registered step "bad (regex" => [func()] has an invalid regular expression: missing closing ).
warning: registered step "No params" => [func()] has an invalid implementation. The first parameter must be of type *testing.T.
warning: registered step "Invalid first param" => [func(string)] has an invalid implementation. The first parameter must be of type *testing.T.
warning: registered step "Extra (param)" => [func(*testing.T)] captures 1 parameter but the supplied implementation takes 0.
warning: registered step "Fewer params" => [func(*testing.T, string)] captures 0 parameters but the supplied implementation takes 1.
warning: registered step "Unconvertible (param)" => [func(*testing.T, elicit_test.custom)] has a parameter type "elicit_test.custom" for which no transforms exist.
```

## Invalid Transforms

+ Create step definitions:

```go
transforms[`Not a function`] = 0
transforms[`bad [regex`] = func() {}
transforms[`Too few params`] = func() {}
transforms[`Too many params`] = func(params []string, s string) {}
transforms[`Incorrect param`] = func(t *testing.T) int { return 0 }
transforms[`No return`] = func(params []string) {}
```

+ Running `go test` will output the following lines:

```
warning: registered transform "Not a function" => [int] must be a function.
warning: registered transform "bad [regex" => [func()] has an invalid regular expression: missing closing ].
warning: registered transform "Too few params" => [func()] must take one argument of type []string.
warning: registered transform "Too many params" => [func([]string, string)] must take one argument of type []string.
warning: registered transform "Incorrect param" => [func(*testing.T) int] must take one argument of type []string.
warning: registered transform "No return" => [func([]string)] must return precisely one value.
```

## Ambiguous Steps

+ Create step definitions:

```go
steps[`(.*)`] = func(t *testing.T, s string) {}
steps[`(something)`] = func(t *testing.T, s string) {}
```

+ Create a `ambiguous_steps.md` file:

```markdown
# Ambiguous Steps
## Ambiguous Step
+ something
```

+ Running `go test` will output the following lines:

```
warning: step "something" is ambiguous:
            - "(.*)" => [func(*testing.T, string)]
            - "(something)" => [func(*testing.T, string)]
warning: registered step "(.*)" => [func(*testing.T, string)] is not used.
warning: registered step "(something)" => [func(*testing.T, string)] is not used.
```

+ Running `go test` will output:

```
Ambiguous Steps
===============
Pending: 1

Ambiguous Step
--------------
Pending

    ? something
```


## Unused Steps

+ Create step definitions:

```go
steps[`.^`] = func(t *testing.T) {}
```

+ Running `go test` will output the following lines:

```
warning: registered step ".^" => [func(*testing.T)] is not used.
```
