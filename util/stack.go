package util

type Stack[T any] struct {
	st    uint32 // stack top/capacity
	sc    uint32 // stack counter
	elems []T    // elements
}

type StaticStack[T any] Stack[T]

func NewStack[T any](cap uint32) *Stack[T] {
	out := new(Stack[T])
	if cap == 0 {
		cap = 8
	} else {
		// make capacity a power of 2
		cap = PowerOfTwoCeil(cap)
	}
	out.st, out.sc = cap, 0
	out.elems = make([]T, 0, cap)
	return out
}

func (s *Stack[T]) grow() {
	if (s.st << 1) < s.st { // check for overflow
		panic("stack overflow")
	}

	newElems := make([]T, s.sc, s.st<<1)
	copy(newElems, s.elems)
	s.elems = newElems
	s.st = s.st << 1
}

func (s *Stack[T]) full() bool { return s.st == s.sc }

func (s *Stack[T]) MultiPush(elems ...T) {
	for _, elem := range elems {
		s.Push(elem)
	}
}

func (s *Stack[T]) MultiPop(n uint32) (elems []T) {
	if s.sc < n {
		panic("tried to access an element on an empty stack")
	}
	elems = make([]T, n)
	for i := uint32(0); i < n; i++ {
		elems[i] = s.unsafePeekOffset(i)
	}
	s.sc = s.sc - n
	return
}

func (s *Stack[T]) Push(elem T) {
	if s.full() {
		s.grow()
	}
	s.elems = append(s.elems, elem)
	s.sc++
}

func (s *Stack[T]) Peek() T {
	if s.Empty() {
		panic("tried to access an element on an empty stack")
	}
	return s.elems[s.sc-1]
}

func (s *Stack[T]) unsafePeekOffset(n uint32) T {
	return s.elems[s.sc-1-n]
}

func (s *Stack[T]) Pop() T {
	elem := s.Peek()
	s.sc--
	return elem
}

func (s *Stack[T]) GetCapacity() uint32 {
	return s.st
}

func (s *Stack[T]) GetCount() uint32 {
	return s.sc
}

func (s *Stack[T]) Empty() bool {
	return s.GetCount() == 0
}
