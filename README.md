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
