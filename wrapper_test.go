package main

import (
	"reflect"
	"testing"
)

type testCaptureResult struct {
	value helpFlagMeta
	ok    bool
}
type testInterpResult struct {
	args []string
	ok   bool
}

type testCase struct {
	name            string
	input           []string
	expectedCapture testCaptureResult
	expectedNewArgs []string
}

/*
Examples of help flag reinterpretation:
(1) go help               -> go help
(2) go -h                 -> go help
(3) go help -h            -> go help
(4) go help subcmd        -> go help subcmd
(5) go -h subcmd          -> go help              # !!! not subcommand help (easier this way)
(6) go subcmd -h          -> go help subcmd
(7) go subcmd -h foo      -> go help subcmd       # basically same as (5)
(8) go subcmd foo -h      -> go help subcmd       # interpreted as 'go subcmd -h'
(9) go subcmd -- foo -h   -> go subcmd -- foo -h  # passthru, e.g case subcmd == run
*/

var testData []testCase = []testCase{
	testCase{ // case 1
		name:  "01-args-help",
		input: []string{"help"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "help",
				help_idx:        0,
				help_arg:        "",
				help_flag_found: false,
				original_args:   []string{"help"},
			},
		},
		expectedNewArgs: []string{"help"},
	},
	testCase{ // case 2
		name:  "02-args-h",
		input: []string{"-h"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "",
				help_idx:        0,
				help_arg:        "-h",
				help_flag_found: true,
				original_args:   []string{"-h"},
			},
		},
		expectedNewArgs: []string{"help"},
	},
	testCase{ // case 3
		name:  "03-args-help-h",
		input: []string{"help", "-h"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "help",
				help_idx:        1,
				help_arg:        "-h",
				help_flag_found: true,
				original_args:   []string{"help", "-h"},
			},
		},
		expectedNewArgs: []string{"help"},
	},
	testCase{ // case 4
		name:  "04-args-help-foo",
		input: []string{"help", "foo"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "help",
				help_idx:        0,
				help_arg:        "",
				help_flag_found: false,
				original_args:   []string{"help", "foo"},
			},
		},
		expectedNewArgs: []string{"help", "foo"},
	},
	testCase{ // case 5
		name:  "05-args-h-foo",
		input: []string{"-h", "foo"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "",
				help_idx:        0,
				help_arg:        "-h",
				help_flag_found: true,
				original_args:   []string{"-h", "foo"},
			},
		},
		expectedNewArgs: []string{"help"},
	},
	testCase{ // case 6
		name:  "06-args-foo-h",
		input: []string{"foo", "-h"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "foo",
				help_idx:        1,
				help_arg:        "-h",
				help_flag_found: true,
				original_args:   []string{"foo", "-h"},
			},
		},
		expectedNewArgs: []string{"help", "foo"},
	},
	testCase{ // case 7
		name:  "07-args-foo-h-bar",
		input: []string{"foo", "-h", "bar"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "foo",
				help_idx:        1,
				help_arg:        "-h",
				help_flag_found: true,
				original_args:   []string{"foo", "-h", "bar"},
			},
		},
		expectedNewArgs: []string{"help", "foo"},
	},
	testCase{ // case 8
		name:  "08-args-foo-bar-h",
		input: []string{"foo", "bar", "-h"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "foo",
				help_idx:        2,
				help_arg:        "-h",
				help_flag_found: true,
				original_args:   []string{"foo", "bar", "-h"},
			},
		},
		expectedNewArgs: []string{"help", "foo"},
	},
	testCase{ // case 9
		name:  "09-args-foo-doubledash-bar-h",
		input: []string{"foo", "--", "bar", "-h"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "foo",
				help_idx:        0,
				help_arg:        "",
				help_flag_found: false, // !!! NOTE: double dash prevents capture (-h is after --)
				original_args:   []string{"foo", "--", "bar", "-h"},
			},
		},
		expectedNewArgs: []string{"foo", "--", "bar", "-h"},
	},
	testCase{ // case 10
		name:  "10-args-foo-h-doubledash-bar-h",
		input: []string{"foo", "-h", "--", "bar", "-h"},
		expectedCapture: testCaptureResult{
			value: helpFlagMeta{
				subcmd:          "foo",
				help_idx:        1,
				help_arg:        "-h",
				help_flag_found: true, // !!! NOTE: double dash does not prevent capture (-h is before --)
				original_args:   []string{"foo", "-h", "--", "bar", "-h"},
			},
		},
		expectedNewArgs: []string{"help", "foo"},
	},
}

func TestCaptureHelp(t *testing.T) {
	var counter int
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			res := captureHelp(test.input)
			if !reflect.DeepEqual(res, test.expectedCapture.value) {
				t.Fatalf("input: %v\nexpected result:\n\t%+#v\ngot:\n\t%+#v", test.input, test.expectedCapture.value, res)
			}
		})
		counter += 1
	}
	if counter <= 0 {
		t.Fatalf("No tests were run?")
	}
}
func TestReinterpretArgs(t *testing.T) {
	var counter int
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			args := test.expectedCapture.value.reinterpretArgs()
			if !reflect.DeepEqual(args, test.expectedNewArgs) {
				t.Fatalf("expected result:\n\t%v\ngot:\n\t%v", test.expectedNewArgs, args)
			}
		})
		counter += 1
	}
	if counter <= 0 {
		t.Fatalf("No tests were run?")
	}
}
