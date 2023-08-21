# The Raiton Language

A toolchain for the Raiton language.

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

The `raiton` tool is implemented in [Go v2.21.0](https://go.dev/dl/), so you will need the `go` tool installed to build the project.
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

for ocaml). The REPL mode (for now) is used to check whether the parser constructs the syntax tree properly and returns an error if there is one.
This is only a temporary implementation for me to manually test the parser and lexer. 

The tool also has a command called `tokenize` to tokenize a file and print out the tokens to `stdout`. This was also useful to manually test
different cases and look at the stream of tokens produced.


## Syntax

For now, the language supports only a single file. The file itself is a `Scope`, which can contain on of the following:
- a definition
- a type definition
- an expression

### Definitions

Definitions are function or value definitions, where an expression is bound to a name. For example:
```bash
# These are definitions

<number>
five: 5

<string>
greeting: "Hello, Raiton"

<number -> number -> number>
add_two a b: (add a b) 
```

The `#` symbol denotes the start of a comment, spanning to the end of the line.

The `<...>` syntax before the identifier is a type declaration for the definition. The first one is a `number`, the next one is of type `string`
and the final definition is a function that takes two `number`s and returns a `number`, represented by `number -> number -> number`.
If you're comming from Ocaml or Haskell, this function type expression should seem familiar.

> Note that the explicit type declaration is temporary for now, but I plan to implement a decent inference system, like the one Ocaml has.

After the identifier, the parser expects a optional list parameters, like shown above. In the case that a definition has parameters,
it's considered a function. If it doesn't contain parameters, it's considered a value bind, but it's still lazily evaluated, which
we'll talk about more later.

The definitions support mutliple expressions if denoted by a block, like this:
```bash
greet name {
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

### Type Definitions

Type definitions are, as the name suggests, definition of user types. They start with the `type` keyword, followed by an identifier
and a type expression after the colon (`:`). For example:
```bash
type person: {
  name: string
  age: number
}

<person>
john: { name: "John" age: 24 }
```
This is an example of a record type named `person` and a definition named `john` of the `person` type. Alternatevly, we could cast the
record literal to the person type like this:
```bash
jane: (person { name: "Jane", age: 24 })
```

This is basically a function call in LISP style `(func arg1 arg2 etc)`, calling the type name as a function, which expects a record of
it's own type and returns it. Because of inference, this is basically a cast. Additionally, to get the value of a record field,
you would do something like:
```bash
(person.name jane) # evaluates to "Jane"
```

### Expressions

Expressions are evaluated lazily. Here are some examples of expressions:
```bash
# number literal
5

# string literal
"John"

# lambda literal
\x: (square x)

# invocation
(concat "Rai" "ton")

# array literal
[1 2 3]

# identifier
some_function

# record litaral
{ attack_power: 100, health_point: 1000 }
```
