# The Raiton Language

[![CI](https://github.com/rfejzic1/raiton/actions/workflows/cli.yml/badge.svg?branch=main)](https://github.com/rfejzic1/raiton/actions/workflows/cli.yml)

A toolchain for the Raiton language.

> Colorless green ideas sleep furiously. - Noam Chomsky

Raiton is a language that is the successor of the [Katon](https://github.com/rfejzic1/katon) language I implemented
with a friend as a project for a university course about programming languages. Back then, I had just started learning about
compilers and interpreters when I started the Katon project, so it's not so well implemented (plus it's a
simple tree-walking interpreter). I learned a lot from implementing the interpreter, but there's still so much to learn.

Enter *Raiton*. I wanted this to be a simple language influenced by procedural and functional langauges
such as Go, Rust, C, Haskell, Ocaml, Clojure, etc. I imagined a functional language, not pure like Haskell, but with some
procedural elements like Ocaml has. I'll probably implement a transpiler to Ocaml first, then decide whether to implement
a VM with bytecode to run on it or to just compile to native code. I prefer the latter, because I'd like to implement
a fast compiler (hopefully) and produce machine code directly. Who knows, maybe add C interop later down the road.

This is a pet project for now, but I hope other people will like the language in time as much as I do.

> Note: This `README.md` describes the current state of the project. It will be updated as I implement more features.

## Getting started

The `raiton` tool is implemented in [Go v1.21.0](https://go.dev/dl/), so you will need the `go` tool installed to build the project.
This will probably work with some older versions, but this is the one I worked with.

After cloning the repository, run `make build` to build the tool. You should get the
`raiton` binary in the `./build` directory.

> Note: This is only built for linux for now. In the future I'll probably add support for other platforms and architectrues.

To get help, run:
```
raiton help
```

## The `raiton` tool

The `raiton` tool is a toolchain that will have the option to compile code, as wall as a REPL (like [utop](https://github.com/ocaml-community/utop)

for ocaml). The REPL mode is used to evaluate the expressions and can report errors in case there are any.
There are only two built-in functions for now:
- `add` to add two integers
- `map` to map elements of arrays and slices

The plan is to extend the available built-in functions and objects.

The tool also has a command called `tokenize` to tokenize a file and print out the tokens to `stdout`. This was also useful to manually test
different cases and look at the stream of tokens produced.


## Syntax

For now, the language supports only a single file. The file itself is a `Scope`, which can contain on of the following:
- a definition
- an expression

### Definitions

Definitions are function or value definitions, where an expression is bound to a name. For example:
```bash
# These are definitions

five: 5

greeting: "Hello, Raiton"

add_two: \a b -> (add a b)
# or
fn add_two a b -> (add a b)
```

The `#` symbol denotes the start of a comment, spanning to the end of the line.

Functions can be defined in two ways, as shown above. One way is to do default definitions and use a function expression
or to use the `fn` keyword to define a function in a more explicit manner. Both have the same outcome, binding a function
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

If you notice, the block is just a scope, like the one at the file level! The colon (`:`) is omitted, because the record
literal syntax uses the curly braces as well. So for now the way to use a scope expression with a definition is to omitt the
colon. The last expression is the one to which the entire scope evaluates to, in this case a function invocation to concatinate
the string arguments.

### Expressions

Expressions are currently evaluated eagerly. The plan is to have them evaluated lazily in the future.
Here are some examples of expressions:
```bash
# number literal
5

# string literal
"John"

# function literal
\x -> (square x)

# function application
(concat "Rai" "ton")

# array literal
[3: 1 2 3]

# slice literal
[1 2 3]

# identifier
some_function

# record litaral
{ attack_power: 100, health_point: 1000 }

# selector
person.name

# array indexing via selector
my_arr.0
```
