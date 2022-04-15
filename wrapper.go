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
	subcmd          string
	help_idx        int
	help_arg        string
	help_flag_found bool
	original_args   []string
}

// Interpret command line arguments; return the finaal argument list, possibly
// differing from original
func (meta helpFlagMeta) reinterpretArgs() []string {
	if meta.help_flag_found {
		if meta.help_idx == 0 || meta.subcmd == "help" {
			return []string{"help"}
		} else {
			return []string{"help", meta.subcmd}
		}
	} else {
		return meta.original_args
	}
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
				subcmd:          subcmd,
				help_idx:        i,
				help_arg:        arg,
				help_flag_found: true,
				original_args:   args}
		}
	}
	return helpFlagMeta{help_flag_found: false, subcmd: subcmd, original_args: args}
}

func main() {
	base := filepath.Base(os.Args[0])

	meta := captureHelp(os.Args[1:])

	if _, debug := os.LookupEnv("GOLANG_HELP_WRAPPER_DEBUG"); debug {
		fmt.Fprintf(os.Stderr, "DEBUG: os.Args: %v\n", os.Args)
		fmt.Fprintf(os.Stderr, "DEBUG: %+#v\n", meta)
	}

	_, suppress_warn := os.LookupEnv("GOLANG_HELP_WRAPPER_WARN_SUPPRESS")

	args := meta.reinterpretArgs()

	if meta.help_flag_found && !suppress_warn {
		// reinterpret help flag
		fmt.Fprintln(os.Stderr, "@@@")
		fmt.Fprintf(os.Stderr, "@@@ WARNING: help flag %q at position %d reinterpreted by %q\n", meta.help_arg, meta.help_idx+1, base)
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
