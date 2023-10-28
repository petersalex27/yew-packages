package table

import (
	"testing"

	"github.com/petersalex27/yew-packages/util/stack"
	"github.com/petersalex27/yew-packages/util/testutil"
)

// TODO: This test is way too complicated. Make it simple in structure or 
// change how TableStack is implemented so the test can be simple in stucture
func TestTableStackAdd(t *testing.T) {
	tests := []struct{
		description string
		initialPushes []Table[int]
		add []tableElement[int]
		topExpect []map[string]tableElement[int]
	}{
		{
			"(add element to empty table stack)",
			[]Table[int]{},
			[]tableElement[int]{{test_nameable("a"), 0}},
			[]map[string]tableElement[int]{{"a": {test_nameable("a"), 0}}},
		},
		{
			"(add element to non-empty table stack)",
			[]Table[int]{
				{
					map[string]tableElement[int]{
						"a": {test_nameable("a"), 0},
					},
				},
			},
			[]tableElement[int]{{test_nameable("b"), 1}},
			[]map[string]tableElement[int]{
				{
					"a": {test_nameable("a"), 0},
					"b": {test_nameable("b"), 1},
				},
			},
		},
		{
			"(add key-val where key already exists in top table)",
			[]Table[int]{
				{
					map[string]tableElement[int]{
						"a": {test_nameable("a"), 0},
					},
				},
			},
			[]tableElement[int]{{test_nameable("b"), 1}},
			[]map[string]tableElement[int]{
				{
					"a": {test_nameable("a"), 0},
					"b": {test_nameable("b"), 1},
				},
			},
		},
		{
			"overwrite to table stack that has the same key on multiple tables",
			[]Table[int]{
				// bottom table
				{
					map[string]tableElement[int]{
						"a": {test_nameable("a"), -1},
					},
				},
				// top table
				{
					map[string]tableElement[int]{
						"a": {test_nameable("a"), 0},
					},
				},
			},
			[]tableElement[int]{
				{test_nameable("a"), 3},
			},
			[]map[string]tableElement[int]{
				{
					"a": {test_nameable("a"), 3},
				},
			},
		},
		{
			"add to table stack that has multiple tables",
			[]Table[int]{
				// bottom table
				{
					map[string]tableElement[int]{
						"x": {test_nameable("x"), -1},
						"y": {test_nameable("y"), -2},
					},
				},
				// top table
				{
					map[string]tableElement[int]{
						"a": {test_nameable("a"), 0},
						"b": {test_nameable("b"), 1},
					},
				},
			},
			[]tableElement[int]{
				{test_nameable("c"), 2},
			},
			[]map[string]tableElement[int]{
				{
					"a": {test_nameable("a"), 0},
					"b": {test_nameable("b"), 1},
					"c": {test_nameable("c"), 2},
				},
			},
		},
		{
			"add and overwrite to table stack that has multiple tables",
			[]Table[int]{
				// bottom table
				{
					map[string]tableElement[int]{
						"x": {test_nameable("x"), -1},
						"y": {test_nameable("y"), -2},
					},
				},
				// top table
				{
					map[string]tableElement[int]{
						"a": {test_nameable("a"), 0},
						"b": {test_nameable("b"), 1},
					},
				},
			},
			[]tableElement[int]{
				{test_nameable("c"), 2},
				{test_nameable("a"), 3},
			},
			[]map[string]tableElement[int]{
				{
					"a": {test_nameable("a"), 0},
					"b": {test_nameable("b"), 1},
					"c": {test_nameable("c"), 2},
				},
				{
					"a": {test_nameable("a"), 3},
					"b": {test_nameable("b"), 1},
					"c": {test_nameable("c"), 2},
				},
			},
		},
	}

	for i, test := range tests {
		stk := NewTableStack[int]()
		// push tables
		for i := range test.initialPushes {
			stk.Push(&test.initialPushes[i])
		}

		expectedLength := stk.Len()
		if expectedLength == 0 { // check for the case when stack is empty
			expectedLength = 1 // update
		}

		for j, addition := range test.add {
			stk.Add(addition.key, addition.val)

			// check stack length
			if stk.Len() != expectedLength {
				t.Fatal(
					testutil.TestFail2("stack-length, "+test.description, 
						expectedLength, stk.Len(), i, j))
			}
			
			// attempt to get top table; stat.IsOk() == true
			table, stat := stk.Peek()
			// check top table
			if stat.IsOk() {
				actual := table.GetRawMap()
				expect := test.topExpect[j]
				expectLen, actLen := len(expect), len(actual) 
				// test length
				if expectLen != actLen {
					t.Fatal(testutil.TestFail2("length-top, "+test.description, expectLen, actLen, i, j))
				}

				// test each element
				for k, v := range expect {
					av, found := actual[k]
					if !found { // test locatability
						t.Fatal(
							testutil.TestSubMsg(i, j, "expected key=`%s` to map to val=`%v`, but key was not found in table", k, v))
					}
					if av.val != v.val { // test value
						t.Fatal(testutil.TestFail2("equality-top, "+test.description, v.val, av.val, i, j))
					}
				}
			} else {
				// failed stat test; stat is not stack.Ok
				t.Fatal(
					testutil.TestFail2("pop-status, "+test.description,
						stack.Ok, stat, i, j),
				)
			}
			
			// check remaining tables (pop tables until stat is empty)
			k := uint(1) // stack element reversed index (0 indexes stack top, i.e., indexes array end)
			// test each top table of actual stack against corr. initial test tables
			for table = stk.peekOffset(k); table != nil; table = stk.peekOffset(k) { // loop until stack is empty
				actual := table.GetRawMap()
				unchangedTableIndex := len(test.initialPushes) - (int(k) + 1)
				expect := test.initialPushes[unchangedTableIndex].GetRawMap()
				expectLen, actLen := len(expect), len(actual) 
				if expectLen != actLen {
					t.Fatal(testutil.TestFail2("length, "+test.description, expectLen, actLen, i, j, int(k-1)))
				}

				// check each element
				for key, v := range expect {
					av, found := actual[key]
					if !found {
						t.Fatal(
							testutil.TestSubMsg(i, j, "[sub-sub=%d] expected key=`%s` to map to val=`%v`, but key was not found in table", k, key, v))
					}
					if av.val != v.val {
						t.Fatal(testutil.TestFail2("equality, "+test.description, v.val, av.val, i, j, int(k-1)))
					}
				}
				k++
			}
		}
	}
}