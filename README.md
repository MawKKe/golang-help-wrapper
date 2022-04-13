# golang-help-wrapper

An utility wrapper program for showing Go help topics

Inspired by https://gist.github.com/MawKKe/485ad4ce21223309d2e90713f3b9b5ba -
now as a simple Go binary!

[![Go](https://github.com/MawKKe/golang-help-wrapper/actions/workflows/go.yml/badge.svg)
](https://github.com/MawKKe/golang-help-wrapper/actions/workflows/go.yml)

## What is this?
It annoys me greatly that none of the go subcommands have a help switch (`-h`
or `--help`), requiring me to use `go help <subcommand>` instead. This is just
annoying, adds more work than is necessary, and results in bad developer
experience. I use shell history (`ctrl-r`, `ctrl-p`, ...) heavily, and this
default mode of operation requires me to type much more than just adding `-h`
at the end of previously executed command.

This wrapper program attempts add the help flag to all the subcommands
by reinterpreting the command line arguments: if `-h` or `--help` is passed
in any position, the respective subcommand help is displayed (the subcommand
is assumed to be the first argument passed, if any).

So for example, this command

    $ go <subcommand> <whatever arguments> -h

will launch

    $ go help <subcommand>

(although this is not always the case; see Usage below.)

# Install
    
Install the executable:

    $ go install github.com/MawKKe/golang-help-wrapper

then add

    alias go=golang-help-wrapper

in your shell configuration (`~/.bashrc` or similar). Now each `go`
invocation will be redirected to this wrapper program, which will in turn
invoke the actual `go` executable with (possibly) reinterpreted arguments.

(Note: `go install` will place the binary in `$GOPATH/bin`; this directory
should exists and be found in your `$PATH` for the alias to work correctly)

# Usage

The program is not intended to be interacted directly; it merely captures
and reinterprets command line arguments that you would otherwise pass to the
`go` executable. However, there are two environment flags that may change
the behavior of the wrapper program:

## Environment variables
The following enviroment variables are understood:

- `GOLANG_HELP_WRAPPER_DEBUG`
    Determines whether to print debug information for the user. For example:

        $ env GOLANG_HELP_WRAPPER_DEBUG=1 go build -h

    prints the following at the beginning of the output:

        DEBUG: os.Args: [./golang-help-wrapper build -h]
        DEBUG: main.helpFlagMeta{subcmd:"build", help_idx:1, help_arg:"-h", \
                help_flag_found:true, original_args:[]string{"build", "-h"}}
    
    these lines may be helful for deducing problems with the wrapper and/or
    installation. You may set the variable to any value to enable it.

- `GOLANG_HELP_WRAPPER_WARN_SUPPRESS`
    By default the wrapper program prints a warning banner (into stderr)
    whenever the command line arguments are reinterpreted:

        $ go build -h
        @@@
        @@@ WARNING: help flag "-h" at position 2 reinterpreted by "golang-help-wrapper"
        @@@ WARNING: -> running 'go help build'
        @@@
        ...

    You may disable this warning banner by setting the variable to any value.

## Bugs
There can be bugs, due to unforeseen edge cases in the `go` command line
interface semantics. If you find such a case, please submit an issue.

Fear not: if you run into problems, you can always invoke the actual `go`
executable directly by prefixing the command name with a backslash. So, instead
of:

   $ go ...
run
   $ \go ...

Note that this is generic behavior found in most(?) shells.

## Tips
If you pass `-h` argument with `go run ...`, it will always be interpreted
as `go help run`, **unless** you use the backslash escape trick **or** you pass
arguments using the double dash idiom:

    $ \go run my/run/target -h
    $ go run -- my/run/target -h

both of these will pass `-h` to the target application, instead of being
reintepreted by the wrapper program.



# Dependencies

The program is written in Go, version 1.18. It may compile with older compiler versions.
The program does not have any third party dependencies.

# License

Copyright 2022 Markus Holmström (MawKKe)

The works under this repository are licenced under Apache License 2.0.
See file `LICENSE` for more information.

# Contributing

This project is hosted at https://github.com/MawKKe/golang-help-wrapper

You are welcome to leave bug reports, fixes and feature requests. Thanks!

