package testutil

import (
	"fmt"
	"strings"
)

func TestSubMsg(index, subindex int, format string, args ...any) string {
	indexStr := fmt.Sprint(index + 1)
	if subindex >= 0 {
		indexStr = indexStr + "." + fmt.Sprint(subindex+1)
	}

	return fmt.Sprintf("failed test #%s: %s", indexStr, fmt.Sprintf(format, args...))
}

func TestMsg(index int, format string, args ...any) string {
	return TestSubMsg(index, -1, format, args...)
}

// examples of calls and corresponding strings
//	TestFail2("myTestCase", "my expected value", "my actual value", 0)
//	`
//	failed test #1 [myTestCase]:
//	expected:
//	my expected value
//	actual:
//	my actual value
//	`
//
//	TestFail2("myTestCase", "my expected value", "my actual value", 0, 0)
//	`
//	failed test #1/1 [myTestCase]:
//	expected:
//	my expected value
//	actual:
//	my actual value
//	`
//
//	TestFail2("", "my expected value", "my actual value", 0, 1, 0)
//	`
//	failed test #1/2/1:
//	expected:
//	my expected value
//	actual:
//	my actual value
//	`
func TestFail2(title string, expected, actual any, index int, subindexes ...int) string {
	indexes := append([]int{index}, subindexes...)
	indexStrs := make([]string, len(indexes))
	for i, index := range indexes {
		indexStrs[i] = fmt.Sprint(index+1)
	}
	indexStr := strings.Join(indexStrs, "/") // 1/2/1/3..

	if title != "" {
		title = fmt.Sprintf(" [%s]", title)
	}

	return fmt.Sprintf("failed test #%s%s:\nexpected:\n%v\nactual:\n%v", indexStr, title, expected, actual)
}

func TestFail(expected, actual any, index int, subindexes ...int) string {
	return TestFail2("", expected, actual, index, subindexes...)
}
