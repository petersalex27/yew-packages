package table

import (
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/util/stack"
)

type TableStack[T any] stack.Stack[*Table[T]]

// return element `offset` elements from the top, e.g., 
//		(&[w, x, y, z]).peekOffset(2) == x
func (ts *TableStack[T]) peekOffset(offset uint) *Table[T] {
	stk := (*stack.Stack[*Table[T]])(ts)
	size := stk.GetCount()
	if size < offset + 1 {
		return nil
	}

	// (&[w, x, y, z]).MultiCheck(2+1) == [x, y, z]
	res, stat := stk.MultiCheck(int(offset)+1)
	if stat.NotOk() {
		return nil
	} 
	
	// elemnt `offset` elements from the top is the first element in result of 
	// MultiCheck
	return res[0]
}

// number of elements in table
func (ts *TableStack[T]) Len() int {
	return int(ts.GetCount())
}

// (Over)writes `val` at domain `key`
func (ts *TableStack[T]) Add(key nameable.Nameable, val T) {
	t, stat := ts.Peek() 
	if stat.NotOk() { // stack is empty
		// create new table
		t = NewTable[T]()
		// defer pushing table
		defer (*stack.Stack[*Table[T]])(ts).Push(t)
	} // TODO: else no defer happens, right?

	// map key -> val
	t.Add(key, val)
}

// If `key` is not found in the table, then `_, false` is returned, else the
// value mapped to by `key` is returned and true is returned
func (ts *TableStack[T]) Get(key nameable.Nameable) (val T, ok bool) {
	t, stat := ts.Peek()
	if ok = stat.IsOk(); !ok {
		return
	}

	return t.Get(key)
}

// Removes key-value pair from table if `key` is in the table, returning the 
// removed value. Otherwise `_, false` is returned.
func (ts *TableStack[T]) Remove(key nameable.Nameable) (val T, ok bool) {
	t, stat := ts.Peek()
	if ok = stat.IsOk(); !ok {
		return
	}

	return t.Remove(key)
}

// Push top table onto stack
func (ts *TableStack[T]) Push(table *Table[T]) {
	(*stack.Stack[*Table[T]])(ts).Push(table)
}

// Pop top table from stack; stat is Empty on failure 
func (ts *TableStack[T]) Pop() (table *Table[T], stat stack.StackStatus) {
	return (*stack.Stack[*Table[T]])(ts).Pop()
}

// Get top table from stack; stat is Empty on failure
func (ts *TableStack[T]) Peek() (table *Table[T], stat stack.StackStatus) {
	return (*stack.Stack[*Table[T]])(ts).Peek()
}

// Get number of tables in stack
func (ts *TableStack[T]) GetCount() uint {
	return (*stack.Stack[*Table[T]])(ts).GetCount()
}

// Get status
func (ts *TableStack[T]) Status() stack.StackStatus {
	return (*stack.Stack[*Table[T]])(ts).Status()
}