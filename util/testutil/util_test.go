package testutil

import (
	"fmt"
	"os"
	"testing"
)

func Test_getMajorCase(t *testing.T) {
	tests := []struct {
		expectedFormat string
		expectedArgs   []any
		format         string
		args           []any
		input          []int
	}{
		// no elems in case indexes
		{"", []any{}, "", []any{}, []int{}},
		// one elem in case indexes
		{"\n\tcase %d", []any{1}, "", []any{}, []int{0}},
		// multiple elems in case indexes
		{"\n\tcase %d", []any{1}, "", []any{}, []int{0, 1}},
		// format is returned when nothing happens
		{"TEST_FORMAT", []any{}, "TEST_FORMAT", []any{}, []int{}},
		// format is prepended
		{"TEST_FORMAT\n\tcase %d", []any{1}, "TEST_FORMAT", []any{}, []int{0}},
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
		{"\nactual:\n%v", []any{1, 2}, "", []any{1}, 2},
		{"TEST_FORMAT\nactual:\n%v", []any{1, 2}, "TEST_FORMAT", []any{1}, 2},
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
		{"\nexpected:\n%v", []any{1, 2}, "", []any{1}, 2},
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
		// no sub case
		{"", []any{}, "", []any{}, []int{0}},
		// no entry of sub case loop
		{", sub %d", []any{2}, "", []any{}, []int{0, 1}},
		// entry of sub-sub case loop
		{", sub %d:%d", []any{2, 3}, "", []any{}, []int{0, 1, 2}},
		// multiple iterations of sub-sub case loop
		{", sub %d:%d:%d", []any{2, 3, 4}, "", []any{}, []int{0, 1, 2, 3}},
		// keep format, no sub case
		{"TEST_FORMAT", []any{}, "TEST_FORMAT", []any{}, []int{0}},
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
		{"\n\tdescription: \"%s\"", []any{"testing helper functions"}, "", []any{}, "testing helper functions"},
		{"TEST_FORMAT\n\tdescription: \"%s\"", []any{"testing helper functions"}, "TEST_FORMAT", []any{}, "testing helper functions"},
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
func TestFailMessage(t *testing.T) {
	expectedArgument, actualArgument := "expected value", "actual value"
	common :=
		"expected:\n" +
			expectedArgument + "\n" +
			"actual:\n" +
			actualArgument + "\n"

	tests := []struct {
		Description
		indexes  []int
		expected string
	}{
		{ // no cases
			Description("this is a test test"),
			[]int{},
			"failed test:\n" +
				"\tdescription: \"this is a test test\"\n" +
				common,
		},
		{ // one case
			Description("this is a test test"),
			[]int{0},
			"failed test:\n" +
				"\tdescription: \"this is a test test\"\n" +
				"\tcase 1\n" +
				common,
		},
		{ // two cases; no sub case loop iterations
			Description("this is a test test"),
			[]int{0, 1},
			"failed test:\n" +
				"\tdescription: \"this is a test test\"\n" +
				"\tcase 1, sub 2\n" +
				common,
		},
		{ // three cases; n=1, n sub case loop iteration
			Description("this is a test test"),
			[]int{0, 1, 2},
			"failed test:\n" +
				"\tdescription: \"this is a test test\"\n" +
				"\tcase 1, sub 2:3\n" +
				common,
		},
		{ // three cases; n+1 sub case loop iterations
			Description("this is a test test"),
			[]int{0, 1, 2, 3},
			"failed test:\n" +
				"\tdescription: \"this is a test test\"\n" +
				"\tcase 1, sub 2:3:4\n" +
				common,
		},
		{ // no description
			NoDescription,
			[]int{0, 1},
			"failed test:\n" +
				"\tcase 1, sub 2\n" +
				common,
		},
	}

	for i, test := range tests {
		actual := test.Description.
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
