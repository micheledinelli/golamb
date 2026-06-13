# GoLamb

![go](https://img.shields.io/badge/Go-1.25-00ADD8.svg?style=plain&logo=Go&logoColor=white)
[![Go Report Card](https://goreportcard.com/badge/github.com/micheledinelli/golamb)](https://goreportcard.com/report/github.com/micheledinelli/golamb)
[![Go Tests](https://github.com/micheledinelli/golamb/actions/workflows/tests.yaml/badge.svg)](https://github.com/micheledinelli/golamb/actions/workflows/tests.yaml)

**GoLamb** is an interpreter for the Untyped Lambda Calculus written in Go. It serves as both an interactive playground (REPL) and a micro-compiler for lambda-expressions.

https://github.com/user-attachments/assets/8a9efde4-ec20-4dd0-84f4-7718bb5dae20

## Examples and Syntax

### Identity

```sh
./golamb
> id = \x.x
id = λx. x
>
> id a
a
>
> id (\z.z w)
λz. z w
```

### Reduction Strategies

Under **Call By Name** or **Normal Order** the engine evaluates the outermost application first, drops the second argument, and returns w.

```sh
./golamb -s cbn
>
> (\a.\b.a) w ((\x.xx)(\x.xx))
w
```

Under **Call By Value** the engine tries to reduce the second argument before invoking the function hence getting stuck.

```sh
./golamb -s cbv
>
> (\a.\b.a) w ((\x.xx)(\x.xx))
# Stuck
```

### Macros Definition

```sh
./golamb
>
> true = \t.\f.t
true = \t.\f.t
>
> false = \t.\f.f
fasle = \t.\f.f
>
> not = \b.b false true
not = \b.b false true
>
> not true
\t.\f.f # false
```

### Import Macros From a File

Define your macros elsewhere and import them using `:import <file>`. Macros are by convention uppercase. Take a look at [std](./std) file.

```sh
./golamb
>
> :import lib.lamb
successfully imported macros from lib.lamb
>
> AND TRUE FALSE
\t.\f.f # FALSE
>
> ID
\x.x
>
> SUCC 1
\f.\x.(f (f x)) # Church numeral 2
>
> ISZERO 0
\t.\f.t # TRUE
```

### Recursive Functions

```sh
./golamb -s cbn
>
> :import lib.lamb
successfully imported macros from lib.lamb
>
> FACT = Y FACT_GEN
defined FACT
>
> FACT 3
\f.\x.(f (f (f (f (f (f x)))))) # Church numeral 6
```

## How To Run

### As a CLI

```sh
go install github.com/micheledinelli/golamb@latest

golamb --help

# Usage of GoLamb:
#
# Options:
#   -s, --strategy string
#         Evaluation strategy: cbn, cbv, normal (default "normal")
#
#   -v, --version
#         Print version information
#
# Examples:
#   ./golamb --strategy=cbn
#   ./golamb -s cbv
```

### Build and Run From Source

```sh
git clone https://github.com/micheledinelli/golamb

cd golamb

# Optional
go build -ldflags "-s -w" -o golamb

# If you built it
./golamb --help
# or simply
go run main.go --help

# Usage of GoLamb:
#
# Options:
#   -s, --strategy string
#         Evaluation strategy: cbn, cbv, normal (default "normal")
#
#   -v, --version
#         Print version information
#
# Examples:
#   ./golamb --strategy=cbn
#   ./golamb -s cbv
```

## Design & Architecture

GoLamb is structured as a classical micro-compiler frontend consisting of four phases:

1. **Lexer:** Scans the input string character-by-character, translating text into a stream of tokens.
2. **Macro Preprocessor:** Substitutes high-level macro definitions (e.g., `TRUE`, `FACT`) from environment state before the core parsing logic takes over.
3. **Parser:** Consumes tokens using a recursive-descent strategy to construct an Abstract Syntax Tree (AST). It strictly enforces **left-associative** function applications (meaning `x y z` is parsed as `((x y) z)`) and expands lambda bodies as far to the right as possible.
4. **Reduction Engine:** Executes pure $\beta$-reductions on the AST using explicit, capture-avoiding substitutions. It handles variable renames ($\alpha$-conversions) whenever an evaluation step threatens to capture a free variable.

## Future Improvements

- [ ] Add [Call-By-Push-Value](https://en.wikipedia.org/wiki/Call-by-push-value) reduction strategy.
- [ ] Add tests.
- [ ] Highlight reduction steps.
- [ ] Introduce types.
- [ ] Make FreshName generator function more robust.

## About the Name

![mascott](./res/mascotte.webp)

## Contributing

Feel free to contribute by submitting issues or pull requests.

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE License. See the [LICENSE](./LICENSE) file for details.
