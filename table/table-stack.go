package table

import (
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/util/stack"
)

type TableStack[T any] stack.Stack[*Table[T]]

func (ts *TableStack[T]) peekOffset(offset uint) *Table[T] {
	stk := (*stack.Stack[*Table[T]])(ts)
	size := stk.GetCount()
	if size < offset + 1 {
		return nil
	}
	stk.MultiPop()
}

func (stack *TableStack[T]) OffsetAdd(key nameable.Nameable, val T) {
	stack
}