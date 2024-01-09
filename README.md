# The Raiton Language

[![CI](https://github.com/rfejzic1/raiton/actions/workflows/cli.yml/badge.svg?branch=main)](https://github.com/rfejzic1/raiton/actions/workflows/cli.yml)

A toolchain for the Raiton language.

> Colorless green ideas sleep furiously. - Noam Chomsky

I imagined Raiton to be a simple language, leaning more to the functional side,
but having procedural elements and side-effects.
I'll decide whether to implement a VM with bytecode to run on it or to just compile
to native code. I prefer the latter, because I'd like to implement a fast compiler (hopefully) and
produce machine code directly. Maybe I'll add C interop later down the road. (Maybe too ambitious, but we'll see)

This is a pet project for now, but I hope other people will like the language in time as much as I do.

> Note: This `README.md` describes the current state of the project. It will be updated as I implement more features.

## Getting started

The `raiton` tool is implemented in [Go v1.21.0](https://go.dev/dl/), so you will need the `go` tool installed to build the project.
This will probably work with some older versions, but this is the one I worked with.

After cloning the repository, run `make build` to build the tool. You should get the
`raiton` binary in the `./build` directory.

> Note: This tool is only built for x86 linux for now. In the future I'll probably add support for other platforms and architectrues.

To get help, run:
```
raiton help
```

## The `raiton` tool

The `raiton` tool is a toolchain that will have the option to compile code, as wall as a REPL (like [utop](https://github.com/ocaml-community/utop)
for ocaml). The REPL mode is used to evaluate the expressions and can report errors in case there are any.

There are only two built-in functions for now:
- `add` to add two integers
- `eq` to compare two values for equality
- `map` to map elements of arrays and slices

The plan is to extend the available built-in functions and objects.

The tool also has a command called `tokenize` to tokenize a file and print out the tokens to `stdout`. This was also useful to manually test
different cases and look at the stream of tokens produced.


## Syntax

For now, the language supports only a single file. The file itself is a `Scope`, which can have *definitions* and *expressions*.

### Definitions

Definitions are a part of a scope which simply binds a name to an expression. Definitions themselves are not expressions,
meaning they don't yield a value. They are intended to simplify the expressions that are evaluated from a scope,
by binding parts of the expression to names.

```bash
# These are definitions

five: 5

greeting: "Hello, Raiton"

add_two: \a b: (add a b)
# or
fn add_two a b: (add a b)
```

The `#` symbol denotes the start of a comment, spanning to the end of the line.

Functions can be defined in two ways, as shown above. One way is to do definitions as is and use a function expression
or to use the `fn` keyword to define a function in a simpler manner. Both have the same outcome, binding a function
expression to a name.

The definitions support mutliple expressions if denoted by a block, like this:
```bash
fn greet name {
  # other nested definitions are supported
  greeting: "Hello"
  
  # the last expression is the one which the entire scope evaluates to
  (concat greeting ", " name)
}
```

If you notice, the block is just a scope, like the one at the file level. Because scopes are *not* expressions, they're
not constructable anywhere in the code. Except for the file-level scope, scopes are expected as blocks in function expressions
and in the consequence and altenative parts of if-expressions. Sintacticaly scopes are wrapped in curly-braces
like this `{ (add 1 2) }`, but if you have a single expression where you need to provide a scope, you can just use the
colon followed by a single expression, like this `: (add 1 2)`. The colon syntax still denotes a scope, but it's just
syntax sugar to not wrap the single expression with braces.

Here is a side-by-side example of this:
```bash
fn greet name { (concat "Hello, " name) }
fn greet name: (concat "Hello, " name)

if condition { "yes" } else { "no" }
if condition: "yes" else: "no"
```


### Expressions

Here are some examples of expressions:
```bash
# number literal
5

# string literal
"John"

# function literal
\x: (square x)

# function application
(concat "Rai" "ton")

# array literal (like arrays in c/c++)
[3: 1 2 3]

# list literal (linked-list)
[1 2 3]

# identifier
some_function

# record litaral
{ attack_power: 100 health_point: 1000 }

# selector
person.name

# array or list indexing via selector
my_arr.0
```
