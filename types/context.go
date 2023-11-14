package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/util/stack"
)

type Context[T nameable.Nameable] struct {
	contextNumber int32
	varCounter uint32
	makeName func(string)T
	stack *stack.Stack[Type[T]]
}

func InheritContext[T nameable.Nameable](parent *Context[T]) *Context[T] {
	child := NewContext[T]()
	child.varCounter = parent.varCounter
	child.stack = parent.stack
	child.makeName = parent.makeName

	return child
}

func NewContext[T nameable.Nameable]() *Context[T] {
	cxt := new(Context[T])
	cxt.stack = stack.NewStack[Type[T]](1 << 5 /*cap=32*/)
	return cxt
}

func NewTestableContext() *Context[nameable.Testable] {
	cxt := NewContext[nameable.Testable]()
	return cxt.SetNameMaker(nameable.MakeTestable)
}

func (cxt *Context[T]) SetNameMaker(f func(string)T) *Context[T] {
	cxt.makeName = f
	return cxt
}

func IsKindVariable[T nameable.Nameable](e expr.Expression[T]) bool {
	_, ok := e.(expr.Variable[T])
	return ok
}