package test

import "fmt"

func TestSubMsg(index, subindex int, format string, args ...any) string {
	indexStr := fmt.Sprint(index+1)
	if subindex >= 0 {
		indexStr = indexStr + "." + fmt.Sprint(subindex+1)
	}

	return fmt.Sprintf("failed test #%s: %s", indexStr, fmt.Sprintf(format, args...))
}

func TestMsg(index int, format string, args ...any) string {
	return TestSubMsg(index, -1, format, args...)
}

// returns the following string
//  "failed test #<index>[.<subindex0>[.<subindex1> ..]]:\nexpected:\n<expected>\nactual:\n<actual>"
func TestFail2(title string, expected, actual any, index int, subindexes ...int) string {
	indexStr := fmt.Sprint(index+1)
	for index := range subindexes {
		indexStr = indexStr + "." + fmt.Sprint(index+1)
	}
	if title != "" {
		title = fmt.Sprintf(" (%s)", title)
	}

	return fmt.Sprintf("failed test #%s%s:\nexpected:\n%v\nactual:\n%v", indexStr, title, expected, actual)
}

func TestFail(expected, actual any, index int, subindexes ...int) string {
	return TestFail2("", expected, actual, index, subindexes...)
}