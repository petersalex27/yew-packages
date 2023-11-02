package testutil

import (
	"fmt"
	"os"
	"testing"
)

func Test_getTesting(t *testing.T) {
	tests := []struct{
		expectedFormat string
		expectedArgs []any
		format string
		args []any
		input string
	}{
		// no input
		{"", []any{}, "", []any{}, ""},
		// no input; existing format
		{"TEST_FORMAT", []any{}, "TEST_FORMAT", []any{}, ""},
		// no input; existing args
		{"", []any{"already here"}, "", []any{"already here"}, ""},
		// input
		{"\n\ttesting: %s", []any{"equality"}, "", []any{}, "equality"},
		// input; existing format
		{"TEST_FORMAT\n\ttesting: %s", []any{"equality"}, "TEST_FORMAT", []any{}, "equality"},

		// The next case would cause issues, were it to be how `args` and format
		// actually looked, when combining it in the output string wherever these 
		// helper functions are used, but that's why they are private functions

		// input; existing args
		{"\n\ttesting: %s", []any{"already here", "equality"}, "", []any{"already here"}, "equality"},
	}

	for _, test := range tests {
		actualFormat, actualArgs := getTesting(test.format, test.args, test.input)

		if actualFormat != test.expectedFormat {
			t.Fatalf("formats are not equiv.\nexp=`%s`\nact=`%s`\n", test.expectedFormat, actualFormat)
		}
		if len(actualArgs) != len(test.expectedArgs) {
			t.Fatalf("args are not equiv.\nexp=`%v`\nact=`%v`\n", test.expectedArgs, actualArgs)
		}
		for i, arg := range test.expectedArgs {
			expected := arg.(string)
			actual, ok := actualArgs[i].(string)
			if !ok {
				t.Fatalf("type assertion actualArgs[%d].(string) failed\n", i)
			}
			if expected != actual {
				t.Fatalf("args at index %d are not equiv.\nexp=`%v`\nact=`%v`\n", i, expected, actual)
			}
		}
	}
}

func Test_getMajorCase(t *testing.T) {
	tests := []struct {
		expectedFormat string
		expectedArgs   []any
		format         string
		args           []any
		input          []int
	}{
		// == tests w/ no input ===================================================

		// no elems in case indexes
		{"", []any{}, "", []any{}, []int{}},
		// no elems in case indexes; existing args
		{"", []any{100}, "", []any{100}, []int{}},
		// format is returned when nothing happens
		{"TEST_FORMAT", []any{}, "TEST_FORMAT", []any{}, []int{}},

		// == tests w/ input ======================================================

		// one elem in case indexes
		{"\n\tcase %d", []any{1}, "", []any{}, []int{0}},
		// one elem in case indexes; existing args 
		{"\n\tcase %d", []any{100, 1}, "", []any{100}, []int{0}},

		// multiple elems in case indexes
		{"\n\tcase %d", []any{1}, "", []any{}, []int{0, 1}},
		// multiple elems in case indexes; existing args
		{"\n\tcase %d", []any{100, 1}, "", []any{100}, []int{0, 1}},

		// format is prepended
		{"TEST_FORMAT\n\tcase %d", []any{1}, "TEST_FORMAT", []any{}, []int{0}},
		// format is prepended; existing args
		{"TEST_FORMAT\n\tcase %d", []any{100, 1}, "TEST_FORMAT", []any{100}, []int{0}},

		// negative input
		{"\n\tcase %d", []any{-1}, "", []any{}, []int{-2}},
	}

	for _, test := range tests {
		actualFormat, actualArgs := getMajorCase(test.format, test.args, test.input)
		if actualFormat != test.expectedFormat {
			t.Fatalf("formats are not equiv.\nexp=`%s`\nact=`%s`\n", test.expectedFormat, actualFormat)
		}
		if len(actualArgs) != len(test.expectedArgs) {
			t.Fatalf("args are not equiv.\nexp=`%v`\nact=`%v`\n", test.expectedArgs, actualArgs)
		}
		for i, arg := range test.expectedArgs {
			expected := arg.(int)
			actual, ok := actualArgs[i].(int)
			if !ok {
				t.Fatalf("type assertion actualArgs[%d].(int) failed\n", i)
			}
			if expected != actual {
				t.Fatalf("args at index %d are not equiv.\nexp=`%v`\nact=`%v`\n", i, expected, actual)
			}
		}
	}
}

func Test_getActual(t *testing.T) {
	tests := []struct {
		expectedFormat string
		expectedArgs   []any
		format         string
		args           []any
		input          int
	}{
		// base case
		{"\nactual:\n%v", []any{2}, "", []any{}, 2},
		// existing args
		{"\nactual:\n%v", []any{1, 2}, "", []any{1}, 2},
		// format prepended
		{"TEST_FORMAT\nactual:\n%v", []any{2}, "TEST_FORMAT", []any{}, 2},
	}

	for _, test := range tests {
		actualFormat, actualArgs := getActual(test.format, test.args, test.input)

		if actualFormat != test.expectedFormat {
			t.Fatalf("formats are not equiv.\nexp=`%s`\nact=`%s`\n", test.expectedFormat, actualFormat)
		}
		if len(actualArgs) != len(test.expectedArgs) {
			t.Fatalf("args are not equiv.\nexp=`%v`\nact=`%v`\n", test.expectedArgs, actualArgs)
		}
		for i, arg := range test.expectedArgs {
			expected := arg.(int)
			actual, ok := actualArgs[i].(int)
			if !ok {
				t.Fatalf("type assertion actualArgs[%d].(int) failed\n", i)
			}
			if expected != actual {
				t.Fatalf("args at index %d are not equiv.\nexp=`%v`\nact=`%v`\n", i, expected, actual)
			}
		}
	}
}

func Test_getExpected(t *testing.T) {
	tests := []struct {
		expectedFormat string
		expectedArgs   []any
		format         string
		args           []any
		input          int
	}{
		// base case
		{"\nexpected:\n%v", []any{2}, "", []any{}, 2},
		// existing args
		{"\nexpected:\n%v", []any{1, 2}, "", []any{1}, 2},
		// existing args; existing format
		{"TEST_FORMAT\nexpected:\n%v", []any{1, 2}, "TEST_FORMAT", []any{1}, 2},
	}

	for _, test := range tests {
		actualFormat, actualArgs := getExpected(test.format, test.args, test.input)

		if actualFormat != test.expectedFormat {
			t.Fatalf("formats are not equiv.\nexp=`%s`\nact=`%s`\n", test.expectedFormat, actualFormat)
		}
		if len(actualArgs) != len(test.expectedArgs) {
			t.Fatalf("args are not equiv.\nexp=`%v`\nact=`%v`\n", test.expectedArgs, actualArgs)
		}
		for i, arg := range test.expectedArgs {
			expected := arg.(int)
			actual, ok := actualArgs[i].(int)
			if !ok {
				t.Fatalf("type assertion actualArgs[%d].(int) failed\n", i)
			}
			if expected != actual {
				t.Fatalf("args at index %d are not equiv.\nexp=`%v`\nact=`%v`\n", i, expected, actual)
			}
		}
	}
}

func Test_getSubCases(t *testing.T) {
	tests := []struct {
		expectedFormat string
		expectedArgs   []any
		format         string
		args           []any
		input          []int
	}{
		// == tests w/ no (sub-case) input ========================================

		// BASE CASE: no input
		{"", []any{}, "", []any{}, []int{0}},
		// no sub case; existing args
		{"", []any{100}, "", []any{100}, []int{0}},
		// keep format, no input
		{"TEST_FORMAT", []any{}, "TEST_FORMAT", []any{}, []int{0}},

		// == tests w/ (sub-case) input ===========================================

		// no entry of sub case loop
		{", sub %d", []any{2}, "", []any{}, []int{0, 1}},
		// no entry of sub case loop; existing args
		{", sub %d", []any{100, 2}, "", []any{100}, []int{0, 1}},
		// entry of sub-sub case loop
		{", sub %d:%d", []any{2, 3}, "", []any{}, []int{0, 1, 2}},
		// entry of sub-sub case loop; existing args
		{", sub %d:%d", []any{100, 2, 3}, "", []any{100}, []int{0, 1, 2}},
		// multiple iterations of sub-sub case loop
		{", sub %d:%d:%d", []any{2, 3, 4}, "", []any{}, []int{0, 1, 2, 3}},
		// multiple iterations of sub-sub case loop; existing args
		{", sub %d:%d:%d", []any{100, 2, 3, 4}, "", []any{100}, []int{0, 1, 2, 3}},
		// prepend format, no entry of sub-sub case loop
		{"TEST_FORMAT, sub %d", []any{2}, "TEST_FORMAT", []any{}, []int{0, 1}},
		// prepend format, entry of sub-sub case loop
		{"TEST_FORMAT, sub %d:%d", []any{2, 3}, "TEST_FORMAT", []any{}, []int{0, 1, 2}},
		// negative, no entry of sub-sub case loop
		{", sub %d", []any{-1}, "", []any{}, []int{0, -2}},
		// negative, entry of sub-sub case loop
		{", sub %d:%d", []any{-1, -2}, "", []any{}, []int{0, -2, -3}},
	}

	for _, test := range tests {
		actualFormat, actualArgs := getSubCases(test.format, test.args, test.input)

		if actualFormat != test.expectedFormat {
			t.Fatalf("formats are not equiv.\nexp=`%s`\nact=`%s`\n", test.expectedFormat, actualFormat)
		}
		if len(actualArgs) != len(test.expectedArgs) {
			t.Fatalf("args are not equiv.\nexp=`%v`\nact=`%v`\n", test.expectedArgs, actualArgs)
		}
		for i, arg := range test.expectedArgs {
			expected := arg.(int)
			actual, ok := actualArgs[i].(int)
			if !ok {
				t.Fatalf("type assertion actualArgs[%d].(int) failed\n", i)
			}
			if expected != actual {
				t.Fatalf("args at index %d are not equiv.\nexp=`%v`\nact=`%v`\n", i, expected, actual)
			}
		}
	}
}

func Test_getDescription(t *testing.T) {
	tests := []struct {
		expectedFormat string
		expectedArgs   []any
		format         string
		args           []any
		input          string
	}{
		// base; no input
		{"", []any{}, "", []any{}, ""},
		// base; input
		{
			"\n\tdescription: \"%s\"", 
			[]any{"testing helper functions"}, 
			"", 
			[]any{}, 
			"testing helper functions",
		},
		// input and existing args
		{
			"\n\tdescription: \"%s\"", 
			[]any{"hello, world!", "testing helper functions"}, 
			"", 
			[]any{"hello, world!"}, 
			"testing helper functions",
		},
		// input and existing format
		{
			"TEST_FORMAT\n\tdescription: \"%s\"", 
			[]any{"testing helper functions"}, 
			"TEST_FORMAT", 
			[]any{}, 
			"testing helper functions",
		},
	}

	for _, test := range tests {
		actualFormat, actualArgs := getDescription(test.format, test.args, test.input)

		if actualFormat != test.expectedFormat {
			t.Fatalf("formats are not equiv.\nexp=`%s`\nact=`%s`\n", test.expectedFormat, actualFormat)
		}
		if len(actualArgs) != len(test.expectedArgs) {
			t.Fatalf("args are not equiv.\nexp=`%v`\nact=`%v`\n", test.expectedArgs, actualArgs)
		}
		for i, arg := range test.expectedArgs {
			expected := arg.(string)
			actual, ok := actualArgs[i].(string)
			if !ok {
				t.Fatalf("type assertion actualArgs[%d].(string) failed\n", i)
			}
			if expected != actual {
				t.Fatalf("args at index %d are not equiv.\nexp=`%v`\nact=`%v`\n", i, expected, actual)
			}
		}
	}
}

// some integration tests
func TestFailMessagef(t *testing.T) {
	tests := []struct {
		expected string
		TestingFor
		format string
		args []any
		indexes  []int
	}{
		{ // no cases
			"failed test: this is my message [1 2 3]\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n",
			Testing("something", "this is a test test"),
			"this is my message %v",
			[]any{[]int{1,2,3}},
			[]int{},
		},
		{ // one case
			"failed test: this is my message [1 2 3]\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n" + 
				"\tcase 1\n",
			Testing("something", "this is a test test"),
			"this is my message %v",
			[]any{[]int{1,2,3}},
			[]int{0},
		},
		{ // two cases; no sub case loop iterations
			"failed test: this is my message [1 2 3]\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n" + 
				"\tcase 1, sub 2\n",
			Testing("something", "this is a test test"),
			"this is my message %v",
			[]any{[]int{1,2,3}},
			[]int{0, 1},
		},
		{ // three cases; n=1, n sub case loop iteration
			"failed test: this is my message [1 2 3]\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n" + 
				"\tcase 1, sub 2:3\n",
			Testing("something", "this is a test test"),
			"this is my message %v",
			[]any{[]int{1,2,3}},
			[]int{0, 1, 2},
		},
		{ // three cases; n+1 sub case loop iterations
			"failed test: this is my message [1 2 3]\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n" + 
				"\tcase 1, sub 2:3:4\n",
			Testing("something", "this is a test test"),
			"this is my message %v",
			[]any{[]int{1,2,3}},
			[]int{0, 1, 2, 3},
		},
		{ // no description
			"failed test: this is my message [1 2 3]\n" +
				"\ttesting: something\n" +
				"\tcase 1, sub 2\n",
			Testing("something"),
			"this is my message %v",
			[]any{[]int{1,2,3}},
			[]int{0, 1},
		},
		{ // no test name
			"failed test: this is my message [1 2 3]\n" +
				"\tdescription: \"this is a test test\"\n" + 
				"\tcase 1, sub 2:3:4\n",
			Testing("", "this is a test test"),
			"this is my message %v",
			[]any{[]int{1,2,3}},
			[]int{0, 1, 2, 3},
		},
		{ // no test name and no description
			"failed test: this is my message [1 2 3]\n" +
				"\tcase 1, sub 2:3:4\n",
			Testing(""),
			"this is my message %v",
			[]any{[]int{1,2,3}},
			[]int{0, 1, 2, 3},
		},
	}

	for i, test := range tests {
		actual := test.TestingFor.
			FailMessagef(test.format, test.args...)(test.indexes...)

		if test.expected != actual {
			fmt.Fprintf(os.Stderr,
				"test #%d failed:\n== expected ========\n`%s`\n== actual ==========\n`%s`\n",
				i+1, test.expected, actual,
			)
			t.Fatal()
		}
	}
}

func TestFailMessage(t *testing.T) {
	expectedArgument, actualArgument := "expected value", "actual value"
	testingFor := Testing("something", "this is a test test")
	common :=
			"expected:\n" +
			expectedArgument + "\n" +
			"actual:\n" +
			actualArgument + "\n"

	tests := []struct {
		TestingFor
		indexes  []int
		expected string
	}{
		{ // no cases
			testingFor,
			[]int{},
			"failed test:\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n" +
				common,
		},
		{ // one case
			testingFor,
			[]int{0},
			"failed test:\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n" +
				"\tcase 1\n" +
				common,
		},
		{ // two cases; no sub case loop iterations
			testingFor,
			[]int{0, 1},
			"failed test:\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n" +
				"\tcase 1, sub 2\n" +
				common,
		},
		{ // three cases; n=1, n sub case loop iteration
			testingFor,
			[]int{0, 1, 2},
			"failed test:\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n" +
				"\tcase 1, sub 2:3\n" +
				common,
		},
		{ // three cases; n+1 sub case loop iterations
			testingFor,
			[]int{0, 1, 2, 3},
			"failed test:\n" +
				"\ttesting: something\n" +
				"\tdescription: \"this is a test test\"\n" +
				"\tcase 1, sub 2:3:4\n" +
				common,
		},
		{ // no description
			Testing("something"),
			[]int{0, 1},
			"failed test:\n" +
				"\ttesting: something\n" +
				"\tcase 1, sub 2\n" +
				common,
		},
		{ // no test name
			Testing("", "this is a test test"),
			[]int{0, 1},
			"failed test:\n" +
				"\tdescription: \"this is a test test\"\n" +
				"\tcase 1, sub 2\n" +
				common,
		},
		{ // no test name and no description
			Testing(""),
			[]int{0, 1},
			"failed test:\n" +
				"\tcase 1, sub 2\n" +
				common,
		},
	}

	for i, test := range tests {
		actual := test.TestingFor.
			FailMessage(expectedArgument, actualArgument, test.indexes...)

		if test.expected != actual {
			fmt.Fprintf(os.Stderr,
				"test #%d failed:\n== expected ========\n`%s`\n== actual ==========\n`%s`\n",
				i+1, test.expected, actual,
			)
			t.Fatal()
		}
	}
}