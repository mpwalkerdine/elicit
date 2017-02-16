# Specification Syntax

Specs are markdown documents which are made executable through the use of
"steps". They are then executed by a normal `go test` run.

For example, the scenarios below make use of the file created below.

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

+ Create a `spec_heading.md` file:

```markdown
# Spec Name
```

+ Running `go test -v` will output:

```
Spec Name
=========
```

The name is derived from the root method defined above, the relative path
of the `.md` file and the spec heading.


## Scenario Heading

Level 2 headings name a Scenario, e.g.

+ Create a `scenario_heading.md` file:

```markdown
# Spec Name
## Scenario Name
```

+ Running `go test -v` will output:

```
Spec Name
=========
Pending: 1

Scenario Name
-------------
Pending

```


## Steps

List items using the `+` character define executable steps, e.g.

+ Create a `simple_step.md` file:

```markdown
# Spec Name
## Scenario Name
+ A Step
```

+ Running `go test -v` will output:

```
Spec Name
=========
Pending: 1

Scenario Name
-------------
Pending

    ? A Step
```

If you wish to include bullets in the specification which aren't steps,
use `-` or `*` in the markdown.

Steps are run in order, one after the other. If a step is missing an
implementation, is skipped or failed, then subsequent steps will be
skipped automatically.

Step implementations are described in [Steps](steps.md).


## Before and After Steps

Steps defined before the first scenario are run before _every_ scenario.

A horizontal rule at the end of the file allows steps to be run after
_every_ scenario, e.g.

+ Create a `before_after_steps.md` file:

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

+ Running `go test -v` will output:

```
Spec
====
Pending: 2

First Scenario
--------------
Pending

    ? before
    ? step 1
    ? after

Last Scenario
-------------
Pending

    ? before
    ? step 2
    ? after
```

Note that, like other steps, the before and after steps will be skipped if an
earlier step is undefined, skipped or failed.


## Tables

Tables defined immediately after a step are passed into the step implementation
as a parameter. The step text for these will include the ☷ symbol for each
table in the output log to indicate that the step implementation must accept an
`elicit.Table` parameter.

All other tables are added into their parent context to be used for step
parameterisation.

Parameterised steps use the `<param name>` syntax which corresponds to a table
header. This creates multiple steps in place of the original, but with values
substituted in from the table.

Note that at present only a single table can be used for parameterisation of a
single step. The mechanism for handling parameterisation across tables hasn't
been decided.

+ Create a `tables.md` file:

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

Scenario steps can use <a>, <b> and <c>, but <d>, <e> and <f> are scoped to this
scenario.

+ print "during: a = <a>"
+ print "during: d = <d>"

---
Any tables in the footer are scoped to the specification.

+ print " after: a = <a>"

```

+ Create step definitions using "mmatt/elicit":

```go
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
        t.Log(v)
    }
```

+ Running `go test -v` will output:

```
Tables
======
Passed: 2

Step Table
----------
Passed

    ✓ print "before: a = 1, b = 2, c = 3"
    ✓ print "before: a = 4, b = 5, c = 6"
    ✓ Step with table ☷
    ✓ print " after: a = 1"
    ✓ print " after: a = 4"

Scenario Tables
---------------
Passed

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

If you need to pass a block of text to a step, you can used a fenced code block
like the majority of this file.

The output appends a ☰ symbol for each text block.

+ Create a `text_block.md` file:

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

+ Create a step definition using "mmatt/elicit", "strings":

```go
steps[`This step takes a block of text:`] =
    func(t *testing.T, text elicit.TextBlock) {
        trimmed := strings.TrimSpace(text.Content)
        lines := strings.Split(trimmed, "\n")
        t.Log("\n> " + strings.Join(lines, "\n> "))
    }
```

+ Running `go test -v` will output:

```
Text Blocks
===========
Passed: 1

Text Block
----------
Passed

    ✓ This step takes a block of text: ☰

--- PASS: Test (0.00s)
    --- PASS: Test/text_block.md/Text_Blocks (0.00s)
        --- PASS: Test/text_block.md/Text_Blocks/Text_Block (0.00s)
        	steps_test.go:15: 
        		> Multiple lines
        		> of text which
        		> are passed into
        		> the step
        		> implementation
```
