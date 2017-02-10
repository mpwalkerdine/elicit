# Specification Syntax

Specs are markdown documents which are made executable through the use of "steps".
They are then executed by a normal go test function.

For example, the scenarios below make use of the file created below.

+ Create a temporary directory
+ Create a `spec_test.go` file:

```go
package elicit_test

import (
    "mmatt/elicit"
    "testing"
)

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        RunTests(t)
}

var steps = map[string]interface{}{}
```

This specification describes (and tests!) the elicit markdown syntax.

## Spec Heading

Level 1 headings name a Spec e.g.

+ Create a `spec_heading.spec` file:

```markdown
# Spec Name
```

+ Running `go test` will output:

```markdown
Spec Name
=========
```

The name is derived from the root method defined above, the relative path
of the `.spec` file and the spec heading.


## Scenario Heading

Level 2 heading name a Scenario, e.g.

+ Create a `scenario_heading.spec` file:

```markdown
# Spec Name
## Scenario Name
```

+ Running `go test` will output:

```markdown
Spec Name
=========

Scenario Name
-------------
```


## Steps

List items using the `+` character define executable steps, e.g.

+ Create a `simple_step.spec` file:

```markdown
# Spec Name
## Scenario Name
+ A Step
```

+ Running `go test` will output:

```markdown
Spec Name
=========

Scenario Name
-------------
  ? A Step
```

If you wish to include bullets in the specification which aren't steps,
use `-` or `*` in the markdown.

Steps are run in order, one after the other. If a step is missing an
implementation, is skipped or failed, then subsequent steps will be
skipped automatically.

Step implementations are described below.


## Before and After Steps

Steps defined before the first scenario are run before _every_ scenario.

A horizontal rule at the end of the file allows steps to be run after
_every_ scenario, e.g.

+ Create a `before_after_steps.spec` file:

```markdown
# Spec
+ before

## First Scenario
+ step 1

## Last Scenario
+ step 2

---
+ after
```

+ Running `go test` will output:

```markdown
Spec
====

First Scenario
--------------
  ? before
  ? step 1
  ? after

Last Scenario
-------------
  ? before
  ? step 2
  ? after
```

Note that, like other steps, the before and after steps will be skipped if an earlier step is undefined, skipped or failed, unless forced with emphasis.

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
    "fmt"
    "testing"
)

func init() {
    steps[`Simple Step`] =
        func(t *testing.T) {
            fmt.Print("simple step, ")
        }

    steps[`Step with "(.*)" parameter`] =
        func(t *testing.T, s string) {
            fmt.Printf("param: %s, ", s)
        }

    steps[`Step with an int parameter (-?\d+)`] =
        func(t *testing.T, i int) {
            fmt.Printf("%d, ", i)
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

+ Running `go test` will output:

```markdown
simple step, param: hello, param: world, 42, -1, 

Step Execution
==============

No Parameters
-------------
  ✓ Simple Step

String parameters
-----------------
  ✓ Step with "hello" parameter
  ✓ Step with "world" parameter

Int parameters
--------------
  ✓ Step with an int parameter 42
  ✓ Step with an int parameter -1

Multiple Parameters
-------------------
  ✓ 1 + 1 = 2
  ✓ 2 + 3 = 5
  ✘ 0 + 1 = 0

--- FAIL: Test (0.00s)
    --- FAIL: Test/step_execution.spec/Step_Execution (0.00s)
        --- FAIL: Test/step_execution.spec/Step_Execution/Multiple_Parameters (0.00s)
            --- FAIL: Test/step_execution.spec/Step_Execution/Multiple_Parameters/3_0_+_1_=_0 (0.00s)
            	steps_test.go:28: expected 0 + 1 = 0, got 1
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

+ Create a `steps_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "testing"
)

func init() {
    steps[`Forced step`] =
        func(t *testing.T) {
            fmt.Print("forced step")
        }
}
```

+ Running `go test` will output:

```markdown
Forcing Steps to Run
====================

Forced Step
-----------
  ? This step is skipped
  ✔ Forced step
```

## Tables

Tables defined immediately after a step are passed into the step implementation as a parameter.
The step text for these will include the ☷ symbol for each table in the output log to indicate
that the step implementation must accept an `elicit.Table` parameter.

All other tables are added into their parent context to be used for step parameterisation.

Parameterised steps use the `<param name>` syntax which corresponds to a table header.
This creates multiple steps in place of the original, but with values substituted in from the table.

Note that at present only a single table can be used for parameterisation of a single step.
The mechanism for handling parameterisation across tables hasn't been decided.

+ Create a `tables.spec` file:

```markdown
# Tables

This table is available in all scenarios, and in before/after steps.

 a | b | c 
---|---|---
 1 | 2 | 3
 4 | 5 | 6

+ print "before: a = <a>, b = <b>, c = <c>"


## Step Table

+ Step with table

 A  | Table  | Here
----|--------|------
 Is | Passed | In

## Scenario Tables 

This is a scenario-specific table

d  | e  | f 
---|----|----
7  | 8  | 9
10 | 11 | 12

Scenario steps can use <a>, <b> and <c>, but <d>, <e> and <f> are scoped to this scenario.

+ print "during: a = <a>"
+ print "during: d = <d>"

---
Any tables in the footer are scoped to the specification.

+ print " after: a = <a>"

```

+ Create a `steps_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "mmatt/elicit"
    "testing"
)

func init() {
    steps[`Step with table`] =
        func(t *testing.T, table elicit.Table) {
            if len(table.Columns) == 0 {
                t.Error("No columns")
            }
            if len(table.Rows) == 0 {
                t.Error("No rows")
            }
        } 

    steps[`print "(.*)"`] =
        func(t *testing.T, v string) {
            fmt.Print(v, "\n")
        }
}
```

+ Running `go test` will output:

```markdown
before: a = 1, b = 2, c = 3
before: a = 4, b = 5, c = 6
 after: a = 1
 after: a = 4
before: a = 1, b = 2, c = 3
before: a = 4, b = 5, c = 6
during: a = 1
during: a = 4
during: d = 7
during: d = 10
 after: a = 1
 after: a = 4


Tables
======

Step Table
----------
  ✓ print "before: a = 1, b = 2, c = 3"
  ✓ print "before: a = 4, b = 5, c = 6"
  ✓ Step with table ☷
  ✓ print " after: a = 1"
  ✓ print " after: a = 4"

Scenario Tables
---------------
  ✓ print "before: a = 1, b = 2, c = 3"
  ✓ print "before: a = 4, b = 5, c = 6"
  ✓ print "during: a = 1"
  ✓ print "during: a = 4"
  ✓ print "during: d = 7"
  ✓ print "during: d = 10"
  ✓ print " after: a = 1"
  ✓ print " after: a = 4"
```

## Text Blocks

If you need to pass a block of text to a step, you can used a fenced code block like the majority of this file.

The output appends a ☰ symbol for each text block.

+ Create a `text_block.spec` file:

````markdown
# Text Blocks
## Text Block
+ This step takes a block of text:

```title
Multiple lines
of text which
are passed into
the step
implementation
```
````

+ Create a `steps_test.go` file:

```go
package elicit_test

import (
    "fmt"
    "mmatt/elicit"
    "strings"
    "testing"
)

func init() {
    steps[`This step takes a block of text:`] =
        func(t *testing.T, text elicit.TextBlock) {
            fmt.Println("> " + strings.Join(strings.Split(strings.TrimSpace(text.Content), "\n"), "\n> "))
        }
}
```

+ Running `go test` will output:

```markdown
> Multiple lines
> of text which
> are passed into
> the step
> implementation


Text Blocks
===========

Text Block
----------
  ✓ This step takes a block of text: ☰
```

---

+ *Remove the temporary directory*