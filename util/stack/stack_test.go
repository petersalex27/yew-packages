package stack

import (
	"testing"

	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestPush(t *testing.T) {
	tests := []struct{
		expect string
	}{
		{""},
		{"a"},
		{"hello, world"},
		{"this is a string with length greater than 32"},
	}

	for i, test := range tests {
		stack := NewStack[byte](32)
		// push entire expected string
		for _, b := range test.expect {
			stack.Push(byte(b))
		}
		
		// check stack counter
		expectedCounter := uint(len(test.expect))
		if stack.sc != expectedCounter {
			t.Fatal(testutil.TestFail2("stack counter", expectedCounter, stack.sc, i))
		}
		// check stack elements
		actual := string(stack.elems[:stack.sc])
		if test.expect != actual {
			t.Fatal(testutil.TestFail2("elements", test.expect, actual, i))
		}
	}
}

func TestPop_empty(t *testing.T) {
	stack := NewStack[byte](8)
	_, stat := stack.Pop()
	if !stat.IsEmpty() {
		t.Fatal(testutil.TestFail(Empty.String(), stat.String(), 0))
	}
}

func TestPop(t *testing.T) {
	tests := []struct{
		put []byte
	}{
		{[]byte("a")},
		{[]byte("hello, world")},
		{[]byte("this is a string with length greater than 32")},
	}

	for i, test := range tests {
		stack := NewStack[byte](uint(len(test.put)))
		// push entire expected string
		for i, b := range test.put {
			stack.elems[i] = byte(b)
		}
		stack.sc = uint(len(test.put))

		for j := range test.put {
			actual, stat := stack.Pop()
			if stat.NotOk() {
				t.Fatal(testutil.TestFail2("pop status", Ok.String(), stat.String(), i, j))
			}
			expect := test.put[len(test.put)-1-j]
			if actual != expect {
				t.Fatal(testutil.TestFail2("pop element", expect, actual, i, j))
			}
		}
		
		// check stack counter
		expectedCounter := uint(0)
		if stack.sc != expectedCounter {
			t.Fatal(testutil.TestFail2("stack counter", expectedCounter, stack.sc, i))
		}
	}
}

func TestPeek(t *testing.T) {
	tests := []struct{
		put []byte
		stat StackStatus
	}{
		{[]byte(""), Empty},
		{[]byte("a"), Ok},
		{[]byte("hello, world"), Ok},
		{[]byte("this is a string with length greater than 32"), Ok},
	}

	for i, test := range tests {
		stack := NewStack[byte](uint(len(test.put)))
		// push entire expected string
		for i, b := range test.put {
			stack.elems[i] = byte(b)
		}
		stack.sc = uint(len(test.put))

		actual, stat := stack.Peek()
		if !stat.Is(test.stat) {
			t.Fatal(testutil.TestFail2("pop status", test.stat.String(), stat.String(), i))
		}
		
		if len(test.put) > 0 {
			expect := test.put[len(test.put)-1]
			if actual != expect {
				t.Fatal(testutil.TestFail2("pop element", expect, actual, i))
			}
		}
		
		// check stack counter
		expectedCounter := uint(len(test.put))
		if stack.sc != expectedCounter {
			t.Fatal(testutil.TestFail2("stack counter", expectedCounter, stack.sc, i))
		}
	}
}

func TestGetCount(t *testing.T) {
	tests := []struct{
		sc uint
	}{
		{0},
		{1},
		{10000},
	}

	for i, test := range tests {
		stack := NewStack[byte](8)
		stack.sc = test.sc
		actual := stack.GetCount()
		if actual != test.sc {
			t.Fatal(testutil.TestFail(test.sc, actual, i))
		}
	}
}

func TestStatus(t *testing.T) {
	tests := []struct{
		sc uint
		stat StackStatus
	}{
		{0, Empty},
		{1, Ok},
		{10000, Ok},
	}

	for i, test := range tests {
		stack := NewStack[byte](8)
		stack.sc = test.sc
		actual := stack.Status()
		if actual != test.stat {
			t.Fatal(testutil.TestFail(test.stat, actual, i))
		}
	}
}