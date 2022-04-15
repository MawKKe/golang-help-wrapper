// golang-help-wrapper: capture and reinterpret '-h' and '--help' flags for running 'go help' topics
package main

// Author: Markus H (MawKKe) markus@mawkke.fi

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type helpFlagMeta struct {
	subcmd        string
	helpIdx       int
	helpArg       string
	helpFlagFound bool
	originalArgs  []string
}

// Interpret command line arguments; return the finaal argument list, possibly
// differing from original
func (meta helpFlagMeta) reinterpretArgs() []string {
	if meta.helpFlagFound {
		if meta.helpIdx == 0 || meta.subcmd == "help" {
			return []string{"help"}
		}
		return []string{"help", meta.subcmd}
	}
	return meta.originalArgs
}

// Parse command line arguments; lookup first occurence of "-h" or "--help" flag.
// Lookup will be terminated if a double dash ("--") is discovered.
func captureHelp(args []string) helpFlagMeta {
	var subcmd string
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") && subcmd == "" {
			subcmd = arg
		}
		if arg == "--" {
			break
		}
		if arg == "-h" || arg == "--help" {
			return helpFlagMeta{
				subcmd:        subcmd,
				helpIdx:       i,
				helpArg:       arg,
				helpFlagFound: true,
				originalArgs:  args}
		}
	}
	return helpFlagMeta{helpFlagFound: false, subcmd: subcmd, originalArgs: args}
}

func main() {
	base := filepath.Base(os.Args[0])

	meta := captureHelp(os.Args[1:])

	if _, debug := os.LookupEnv("GOLANG_HELP_WRAPPER_DEBUG"); debug {
		fmt.Fprintf(os.Stderr, "DEBUG: os.Args: %v\n", os.Args)
		fmt.Fprintf(os.Stderr, "DEBUG: %+#v\n", meta)
	}

	_, suppressWarn := os.LookupEnv("GOLANG_HELP_WRAPPER_WARN_SUPPRESS")

	args := meta.reinterpretArgs()

	if meta.helpFlagFound && !suppressWarn {
		// reinterpret help flag
		fmt.Fprintln(os.Stderr, "@@@")
		fmt.Fprintf(os.Stderr, "@@@ WARNING: help flag %q at position %d reinterpreted by %q\n", meta.helpArg, meta.helpIdx+1, base)
		fmt.Fprintf(os.Stderr, "@@@ WARNING: -> running 'go %s'\n", strings.Join(args, " "))
		fmt.Fprintln(os.Stderr, "@@@")
	}

	cmd := exec.Command("go", args...)

	// passthrough all descriptors
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()

	// propagate error code to caller
	os.Exit(cmd.ProcessState.ExitCode())
}
