package stack

import (
	"testing"

	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestPush_ret(t *testing.T) {
	tests := []struct {
		expect string
	}{
		{""},
		{"a"},
		{"hello, world"},
		{"this is a string with length greater than 32"},
	}

	for i, test := range tests {
		stack := NewSaveStack[byte](32)
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

func TestPop_empty_ret(t *testing.T) {
	stack := NewSaveStack[byte](8)
	_, stat := stack.Pop()
	if !stat.IsEmpty() {
		t.Fatal(testutil.TestFail(Empty.String(), stat.String(), 0))
	}
}

func TestPop_ret(t *testing.T) {
	tests := []struct {
		put []byte
	}{
		{[]byte("a")},
		{[]byte("hello, world")},
		{[]byte("this is a string with length greater than 32")},
	}

	for i, test := range tests {
		stack := NewSaveStack[byte](uint(len(test.put)))
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

func TestPeek_ret(t *testing.T) {
	tests := []struct {
		put  []byte
		stat StackStatus
	}{
		{[]byte(""), Empty},
		{[]byte("a"), Ok},
		{[]byte("hello, world"), Ok},
		{[]byte("this is a string with length greater than 32"), Ok},
	}

	for i, test := range tests {
		stack := NewSaveStack[byte](uint(len(test.put)))
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

func TestGetCount_ret(t *testing.T) {
	tests := []struct {
		sc uint
	}{
		{0},
		{1},
		{10000},
	}

	for i, test := range tests {
		stack := NewSaveStack[byte](8)
		stack.sc = test.sc
		actual := stack.GetCount()
		if actual != test.sc {
			t.Fatal(testutil.TestFail(test.sc, actual, i))
		}
	}
}

func TestStatus_ret(t *testing.T) {
	tests := []struct {
		sc   uint
		stat StackStatus
	}{
		{0, Empty},
		{1, Ok},
		{10000, Ok},
	}

	for i, test := range tests {
		stack := NewSaveStack[byte](8)
		stack.sc = test.sc
		actual := stack.Status()
		if actual != test.stat {
			t.Fatal(testutil.TestFail(test.stat, actual, i))
		}
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		beforeSave string
		afterSave  string
		peek       rune
		stat       StackStatus
		total      string
	}{
		{"", "", 0, Empty, ""},
		{"", "a", 'a', Ok, "a"},
		{"a", "", 0, Empty, "a"},
		{"123", "456", '6', Ok, "123456"},
	}

	for i, test := range tests {
		stack := NewSaveStack[byte](32, 4)

		// push
		for _, b := range []byte(test.beforeSave) {
			stack.Push(b)
		}

		// now save and check
		stack.Save()
		expectBC := uint(len(test.beforeSave))
		actualBC := stack.bc
		if expectBC != actualBC {
			t.Fatal(testutil.TestFail2("base counter", expectBC, actualBC, i))
		}

		// push after save
		for _, b := range []byte(test.afterSave) {
			stack.Push(b)
		}
		// check that bc hasn't changed
		if expectBC != stack.bc {
			t.Fatal(testutil.TestFail2("base counter after", expectBC, stack.bc, i))
		}
		actual, stat := stack.Peek()
		if !stat.Is(test.stat) {
			t.Fatal(testutil.TestFail2("status", test.stat.String(), stat.String(), i))
		}

		if stat.IsOk() {
			expect := test.afterSave[len(test.afterSave)-1]
			if actual != expect {
				t.Fatal(testutil.TestFail2("peek result", expect, actual, i))
			}
		}

		actualTotal := string(stack.elems[:stack.sc])
		if actualTotal != test.total {
			t.Fatal(testutil.TestFail2("total result", test.total, actualTotal, i))
		}
	}
}

func TestReturn(t *testing.T) {
	tests := []struct {
		afterReturn   string
		returnStatus StackStatus
		peek         rune
		stat         StackStatus
		total        string
	}{
		{"", Ok, 0, Empty, ""},
		{"", Ok, 0, Empty, "a"},
		{"a", Ok, 'a', Ok, "a"},
		{"123", Ok, '3', Ok, "123456"},
	}

	for i, test := range tests {
		stack := NewSaveStack[byte](32, 4)
		// place all elems
		for i := range test.total {
			stack.elems[i] = test.total[i]
		}

		// set counter
		stack.sc = uint(len(test.total))
		// set return point
		stack.returnStack.Push(0)
		// set base counter
		ret := uint(len(test.afterReturn))
		stack.bc = ret

		// now return
		retStat := stack.Return()
		if !retStat.Is(test.returnStatus) {
			t.Fatal(testutil.TestFail2("return status", test.returnStatus.String(), retStat.String(), i))
		}

		expectBC := uint(0)
		actualBC := stack.bc
		if expectBC != actualBC {
			t.Fatal(testutil.TestFail2("base counter", expectBC, actualBC, i))
		}

		actual, stat := stack.Peek()
		if !stat.Is(test.stat) {
			t.Fatal(testutil.TestFail2("status", test.stat.String(), stat.String(), i))
		}

		if stat.IsOk() {
			expect := test.afterReturn[len(test.afterReturn)-1]
			if actual != expect {
				t.Fatal(testutil.TestFail2("peek result", expect, actual, i))
			}
		}

		actualResult := string(stack.elems[:stack.sc])
		if actualResult != test.afterReturn {
			t.Fatal(testutil.TestFail2("after return result", test.afterReturn, actualResult, i))
		}
	}
}
