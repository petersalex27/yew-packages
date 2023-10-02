package stack

type SaveStack[T any] struct {
	Stack[T]
	bc          uint // base counter
	returnStack Stack[uint]
}

// only up to 2 arguments are used, one argument required; cap is the starting
// capacity for the stack. returnCap[0] is the starting capacity for the
// number of returns that can be saved. returnCap has a default value of 8
func NewSaveStack[T any](cap uint, returnCap ...uint) *SaveStack[T] {
	out := new(SaveStack[T])
	out.Stack = makeStack[T](cap)
	out.bc = 0
	rCap := uint(8)
	if len(returnCap) > 0 {
		rCap = returnCap[0]
	}
	out.returnStack = makeStack[uint](rCap)
	return out
}

func (stack *SaveStack[T]) Push(elem T) {
	stack.Stack.Push(elem)
}

func (stack *SaveStack[T]) Pop() (elem T, stat StackStatus) {
	if stack.bc == stack.sc {
		stat = Empty
		return
	}
	return stack.Stack.Pop()
}

func (stack *SaveStack[T]) Peek() (elem T, stat StackStatus) {
	if stack.bc == stack.sc {
		stat = Empty
		return
	}
	return stack.Stack.Peek()
}

func (stack *SaveStack[T]) GetCount() uint {
	return stack.GetFullCount() - stack.bc
}

func (stack *SaveStack[T]) GetFullCount() uint {
	return stack.Stack.GetCount()
}

func (stack *SaveStack[T]) Status() StackStatus {
	return stack.Stack.Status()
}

func (stack *SaveStack[T]) Save() {
	save := stack.bc
	stack.returnStack.Push(save)
	stack.bc = stack.sc
}

func (stack *SaveStack[T]) Return() (stat StackStatus) {
	if stack.returnStack.GetCount() == 0 {
		return IllegalReturn
	}
	old := stack.bc
	stack.bc, stat = stack.returnStack.Pop()
	if stat.IsOk() {
		stack.sc = old
	}
	return
}
