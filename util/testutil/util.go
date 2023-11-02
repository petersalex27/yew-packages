package testutil

import (
	"fmt"
)

// Test description.
//
// why? ..
//
// don't have to remember different function names for
// functions that do the same but w/o a description; just call the public
// functions with the same names as the public methods of Description
type Description string

type TestingFor struct {
	What string
	Description
}

// Testing creates a name for the thing being tested with an optional 
// description. If no argument is given for description, the output won't have
// one. If one or more is given, then only the first description will be used
func Testing(what string, description ...string) TestingFor {
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	return TestingFor{
		What: what,
		Description: Description(desc),
	}
}

// constant receiver for Description functions that causes description
// related parts of functions to be ignored
const NoDescription Description = Description("")

// returns a message saying a test with `description` and, optionaly, test case
// number indexes[0]+1 and sub case numbers index[1:] (+1 to value for ea.) and
// a message formated according to `format` and `args` failed
//
// example:
//
//	Testing("stuff", "my test").FailMessagef("%s%v","arr = ",[]int{1,2,2})(15,2,3)
//	return
//	`failed test: arr = [1 2 2]
//		testing: stuff
//		description: "my test"
//		case 16, sub 3:4
// 	`
func (test TestingFor) FailMessagef(format string, args ...any) func(indexes ...int) string {
	// create user message	
	userMessage := fmt.Sprintf(format, args...)

	// inner function allowing multiple variable-length args to be used
	return func(indexes ...int) string {
		format = "failed test: %s"
		// +1 for user msg, +1 for description, +1 for testing
		maxArgLen := len(indexes) + 2 
		args := make([]any, 0, maxArgLen)
		args = append(args, userMessage)

		format, args = getTesting(format, args, test.What)
		format, args = getDescription(format, args, string(test.Description))
		format, args = getMajorCase(format, args, indexes)
		format, args = getSubCases(format, args, indexes)
		format = format + "\n"
		return fmt.Sprintf(format, args...)
	}
}

// returns a message saying a test with `description` and, optionaly, test case
// number indexes[0]+1 and sub case numbers index[1:] (+1 to value for ea.) and
// a message formated according to `format` and `args` failed
//
// example:
//
//	Description("my test").FailMessagef("%s%v","arr = ",[]int{1,2,2})(15,2,3)
//	return
//	`failed test: arr = [1 2 2]
//		description: "my test"
//		case 16, sub 3:4
// 	`
func (description Description) FailMessagef(format string, args ...any) func(indexes ...int) string {
	return Testing("").FailMessagef(format, args...)
}

// returns a message saying a test with test case number indexes[0]+1 and sub-
// case numbers index[1:] (+1 to value for ea.) and a message formated
// according to `format` and `args` failed.
//
// example:
//
//	FailMessagef("%s%v", "arr = ", []int{1, 2, 2})(15, 2, 3)
//	return
//	`failed test: arr = [1 2 2]
//		case 16, sub 3:4
//	`
func FailMessagef(format string, args ...any) func(indexes ...int) string {
	return (TestingFor{"",NoDescription}).FailMessagef(format, args...)
}

// return updated format and args to account for major case number
//
// helper function for
//
//	(Description).FailMessage(any, any, ...int)
func getMajorCase(format string, args []any, indexes []int) (newFormat string, newArgs []any) {
	if len(indexes) > 0 {
		format = format + "\n\tcase %d"
		args = append(args, indexes[0]+1)
	}
	return format, args
}

// return updated format and args to account for sub cases
//
// helper function for
//
//	(Description).FailMessage(any, any, ...int)
func getSubCases(format string, args []any, indexes []int) (newFormat string, newArgs []any) {
	if len(indexes) <= 1 {
		// guard: do not enter loop below b/c that will cause an out of
		// bounds exception
		return format, args
	}

	// get sub case
	format = format + ", sub %d"
	args = append(args, indexes[1]+1)

	// get sub-sub case, sub-sub-sub case, .. numbers
	for _, index := range indexes[2:] {
		format = format + ":%d"
		args = append(args, index+1)
	}
	return format, args
}

// if description is valid, return:
//
//	newFormat = format + "\n\tdescription: \"%s\""
//	newArgs = append(args, description)
//
// else return:
//
//	newFormat = format
//	newArgs = args
func getTesting(format string, args []any, what string) (newFormat string, newArgs []any) {
	if what != "" {
		format = format + "\n\ttesting: %s"
		args = append(args, what)
	}
	return format, args
}

// if description is valid, return:
//
//	newFormat = format + "\n\tdescription: \"%s\""
//	newArgs = append(args, description)
//
// else return:
//
//	newFormat = format
//	newArgs = args
func getDescription(format string, args []any, description string) (newFormat string, newArgs []any) {
	if description != "" {
		format = format + "\n\tdescription: \"%s\""
		args = append(args, description)
	}
	return format, args
}

// returns:
//
//	newFormat = format + "\nexpected:\n%v"
//	newArgs = append(args, expected)
func getExpected(format string, args []any, expected any) (newFormat string, newArgs []any) {
	format = format + "\nexpected:\n%v"
	args = append(args, expected)
	return format, args
}

// returns:
//
//	newFormat = format + "\nactual:\n%v"
//	newArgs = append(args, actual)
func getActual(format string, args []any, actual any) (newFormat string, newArgs []any) {
	format = format + "\nactual:\n%v"
	args = append(args, actual)
	return format, args
}

// returns a message saying a test with `description` and, optionaly, test case
// number indexes[0]+1 and sub case numbers index[1:] (+1 to value for ea.);
// then gives the `expected` value and `actual` value
//
// example:
//
//	FailMessage([]int{1,2,3}, []int{1,2,2}, 15, 2, 3)
//	return`
//	failed test:
//		case 16, sub 3:4
//	expected:
//	[1 2 3]
//	actual:
//	[1 2 2]`
func FailMessage(expected any, actual any, caseIndexes ...int) string {
	return NoDescription.FailMessage(expected, actual, caseIndexes...)
}

// returns a message saying a test with `description` and, optionaly, test case
// number indexes[0]+1 and sub case numbers index[1:] (+1 to value for ea.);
// then gives the `expected` value and `actual` value
//
// example:
//
//	Testing("stuff", "some case").FailMessage([]int{1,2,3}, []int{1,2,2}, 15, 2, 3)
//	return
//	`failed test:
//		testing: stuff
//		description: "some case"
//		case 16, sub 3:4
//	expected:
//	[1 2 3]
//	actual:
//	[1 2 2]`
func (test TestingFor) FailMessage(expected any, actual any, caseIndexes ...int) string {
	// in its full form:
	//	_args:_____________index[0]___index[1:]___________expected______actual___
	//	_spec:________________v_______vvvvvvvvvv_____________v____________v______
	//	"failed test:\n\tcase %d, sub %d:%d..:%d\nexpected:\n%v\nactual:\n%v\n"
	var format string = "failed test:"
	// room for each arg: 
	//	test.What=1, test.Description=1, expected=1, actual=1, len(caseIndexes)
	maxLen := 4 + len(caseIndexes)
	// accumulated arguments
	args := make([]any, 0, maxLen)

	// build format and args from parts
	format, args = getTesting(format, args, test.What)
	format, args = getDescription(format, args, string(test.Description))
	format, args = getMajorCase(format, args, caseIndexes)
	format, args = getSubCases(format, args, caseIndexes)
	format, args = getExpected(format, args, expected)
	format, args = getActual(format, args, actual)
	format = format + "\n"

	// return string-ed result
	return fmt.Sprintf(format, args...)
}

// returns a message saying a test with `description` and, optionaly, test case
// number indexes[0]+1 and sub case numbers index[1:] (+1 to value for ea.);
// then gives the `expected` value and `actual` value
//
// example:
//
//	Description("some case").FailMessage([]int{1,2,3}, []int{1,2,2}, 15, 2, 3)
//	return
//	`failed test:
//		description: "some case"
//		case 16, sub 3:4
//	expected:
//	[1 2 3]
//	actual:
//	[1 2 2]`
func (description Description) FailMessage(expected any, actual any, caseIndexes ...int) string {
	return Testing("", string(description)).FailMessage(expected, actual, caseIndexes...)
}

// == deprecated stuff: STUFF BELOW WILL BE REMOVED! ==========================

// Returns a message saying test failed at given case indexes and then a
// message according to `format` and `args`
//
// Deprecated: replaced by more general function FailMessagef or
// (Description).FailMessagef
func TestSubMsg(index, subindex int, format string, args ...any) string {
	indexStr := fmt.Sprint(index + 1)
	if subindex >= 0 {
		indexStr = indexStr + "." + fmt.Sprint(subindex+1)
	}

	return fmt.Sprintf("failed test #%s: %s", indexStr, fmt.Sprintf(format, args...))
}

// Returns a message saying test failed at given case index and then a message
// according to `format` and `args`
//
// Deprecated: replaced by more general function FailMessagef or
// (Description).FailMessagef
func TestMsg(index int, format string, args ...any) string {
	return TestSubMsg(index, -1, format, args...)
}

// Returns a string reporting a test failure b/c expected != actual
//
// Deprecated: Use (Description).FailMessage(any, any, ...int); the
// functionality of this function has been moved to
// (Description).FailMessage(any, any, ...int). TestFail2 will eventually be
// removed. Why? The name of this function is bad, and it requires a major
// test case to be passed as an argument which is not a thing in some tests or
// just too much info sometimes
func TestFail2(title string, expected, actual any, testCaseIndex int, minorTestCaseIndexes ...int) string {
	cases := append([]int{testCaseIndex}, minorTestCaseIndexes...)
	return Description(title).FailMessage(expected, actual, cases...)
}

// Returns a string reporting a test failure b/c expected != actual
//
// Deprecated: Use FailMessage(any, any, ...int); the functionality of
// this function has been moved to FailMessage(any, any, ...int).
// TestFail will eventually be removed. Why? the name of this function is bad,
// and it requires a major test case to be passed as an argument which is not a
// thing in some tests or just too much info sometimes
func TestFail(expected, actual any, index int, subindexes ...int) string {
	cases := append([]int{index}, subindexes...)
	return FailMessage(expected, actual, cases...)
}
