package table

import (
	"testing"

	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestAdd(t *testing.T) {
	tests := []struct{
		add []tableElement[int]
		expect []map[string]int
	}{
		{
			[]tableElement[int]{{test_nameable("a"), 0}},
			[]map[string]int{{"a": 0}},
		},
		{
			[]tableElement[int]{
				{test_nameable("a"), 0},
				{test_nameable("b"), 0},
			},
			[]map[string]int{
				{"a": 0},
				{"a": 0, "b": 0},
			},
		},
		{
			[]tableElement[int]{
				{test_nameable("a"), 0},
				{test_nameable("b"), 0},
				{test_nameable("a"), 1},
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
		remove test_nameable
		removedVal int
		expectOk bool
		expect map[string]int
	}{
		{
			map[string]tableElement[int]{},
			test_nameable("a"),
			0, false,
			map[string]int{},
		},
		{
			map[string]tableElement[int]{
				"a": {test_nameable("a"), 3},
			},
			test_nameable("x"),
			0, false,
			map[string]int{"a": 3},
		},
		{
			map[string]tableElement[int]{
				"a": {test_nameable("a"), 3},
			},
			test_nameable("a"),
			3, true,
			map[string]int{},
		},
		{
			map[string]tableElement[int]{
				"a": {test_nameable("a"), 3},
				"b": {test_nameable("b"), 4},
			},
			test_nameable("a"),
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
		key test_nameable
		gotVal int
		expectOk bool
		expect map[string]int // confirm unchanged
	}{
		{
			map[string]tableElement[int]{},
			test_nameable("a"),
			0, false,
			map[string]int{},
		},
		{
			map[string]tableElement[int]{
				"a": {test_nameable("a"), 3},
			},
			test_nameable("x"),
			0, false,
			map[string]int{"a": 3},
		},
		{
			map[string]tableElement[int]{
				"a": {test_nameable("a"), 3},
			},
			test_nameable("a"),
			3, true,
			map[string]int{"a": 3},
		},
		{
			map[string]tableElement[int]{
				"a": {test_nameable("a"), 3},
				"b": {test_nameable("b"), 4},
			},
			test_nameable("a"),
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