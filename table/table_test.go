package table

import (
	"testing"

	"github.com/petersalex27/yew-packages/util/testutil"
)

type testName string

func (n testName) GetName() string { return string(n) }

func TestAdd(t *testing.T) {
	tests := []struct{
		add []tableElement[int]
		expect []map[string]int
	}{
		{
			[]tableElement[int]{{testName("a"), 0}},
			[]map[string]int{{"a": 0}},
		},
		{
			[]tableElement[int]{
				{testName("a"), 0},
				{testName("b"), 0},
			},
			[]map[string]int{
				{"a": 0},
				{"a": 0, "b": 0},
			},
		},
		{
			[]tableElement[int]{
				{testName("a"), 0},
				{testName("b"), 0},
				{testName("a"), 1},
			},
			[]map[string]int{
				{"a": 0},
				{"a": 0, "b": 0},
				{"a": 1, "b": 0},
			},
		},
	}

	for i, test := range tests {
		table := NewTable[int]()
		for j, addition := range test.add {
			table.Add(addition.key, addition.val)
			
			// check map equiv.
			actual := table.GetRawMap() // actual map
			expectLen, actLen := len(test.expect[j]), len(actual) 
			if expectLen != actLen {
				t.Fatal(testutil.TestFail2("length", expectLen, actLen, i, j))
			}

			// check each element
			for k, v := range test.expect[j] {
				av, found := actual[k]
				if !found {
					t.Fatal(testutil.TestSubMsg(i, j, "expected key=`%s` to map to val=`%v`, but key was not found in table", k, v))
				}
				if av.val != v {
					t.Fatal(testutil.TestFail2("equality", v, av.val, i, j))
				}
			}
		}
	}
}

func TestRemove(t *testing.T) {
	tests := []struct{
		data map[string]tableElement[int]
		remove testName
		removedVal int
		expectOk bool
		expect map[string]int
	}{
		{
			map[string]tableElement[int]{},
			testName("a"),
			0, false,
			map[string]int{},
		},
		{
			map[string]tableElement[int]{
				"a": {testName("a"), 3},
			},
			testName("x"),
			0, false,
			map[string]int{"a": 3},
		},
		{
			map[string]tableElement[int]{
				"a": {testName("a"), 3},
			},
			testName("a"),
			3, true,
			map[string]int{},
		},
		{
			map[string]tableElement[int]{
				"a": {testName("a"), 3},
				"b": {testName("b"), 4},
			},
			testName("a"),
			3, true,
			map[string]int{"b": 4},
		},
	}

	for i, test := range tests {
		table := NewTable[int]()
		table.data = test.data
		val, ok := table.Remove(test.remove)

		if test.expectOk != ok {
			t.Fatal(testutil.TestFail2("ok", test.expectOk, ok, i))
		}

		if test.removedVal != val {
			t.Fatal(testutil.TestFail2("val", test.removedVal, val, i))
		}

		// check map equiv.
		actual := table.GetRawMap() // actual map
		expectLen, actLen := len(test.expect), len(actual) 
		if expectLen != actLen {
			t.Fatal(testutil.TestFail2("length", expectLen, actLen, i))
		}

		// check each element
		for k, v := range test.expect {
			av, found := actual[k]
			if !found {
				t.Fatal(testutil.TestMsg(i, "expected key=`%s` to map to val=`%v`, but key was not found in table", k, v))
			}
			if av.val != v {
				t.Fatal(testutil.TestFail2("equality", v, av.val, i))
			}
		}
	}
}

func TestGet(t *testing.T) {
	tests := []struct{
		data map[string]tableElement[int]
		key testName
		gotVal int
		expectOk bool
		expect map[string]int // confirm unchanged
	}{
		{
			map[string]tableElement[int]{},
			testName("a"),
			0, false,
			map[string]int{},
		},
		{
			map[string]tableElement[int]{
				"a": {testName("a"), 3},
			},
			testName("x"),
			0, false,
			map[string]int{"a": 3},
		},
		{
			map[string]tableElement[int]{
				"a": {testName("a"), 3},
			},
			testName("a"),
			3, true,
			map[string]int{"a": 3},
		},
		{
			map[string]tableElement[int]{
				"a": {testName("a"), 3},
				"b": {testName("b"), 4},
			},
			testName("a"),
			3, true,
			map[string]int{"a": 3, "b": 4},
		},
	}

	for i, test := range tests {
		table := NewTable[int]()
		table.data = test.data
		val, ok := table.Get(test.key)

		if test.expectOk != ok {
			t.Fatal(testutil.TestFail2("ok", test.expectOk, ok, i))
		}

		if test.gotVal != val {
			t.Fatal(testutil.TestFail2("val", test.gotVal, val, i))
		}

		// check map equiv.
		actual := table.GetRawMap() // actual map
		expectLen, actLen := len(test.expect), len(actual) 
		if expectLen != actLen {
			t.Fatal(testutil.TestFail2("length", expectLen, actLen, i))
		}

		// check each element
		for k, v := range test.expect {
			av, found := actual[k]
			if !found {
				t.Fatal(testutil.TestMsg(i, "expected key=`%s` to map to val=`%v`, but key was not found in table", k, v))
			}
			if av.val != v {
				t.Fatal(testutil.TestFail2("equality", v, av.val, i))
			}
		}
	}
}