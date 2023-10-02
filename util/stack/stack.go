package stack

import (
	"fmt"

	"github.com/petersalex27/yew-packages/util"
)

type StackType[T any] interface {
	Push(T)
	Pop() (T, StackStatus)
	Peek() (T, StackStatus)
	GetCount() uint
	Status() StackStatus
}

type StackPlus[T any] interface {
	StackType[T]
	Clear(uint)
	MultiPush(...T)
	MultiPop(uint) ([]T, StackStatus)
}

type ReturnableStack[T any] interface {
	StackType[T]
	GetFullCount() uint
	Save()
	Return() StackStatus
}

type Stack[T any] struct {
	st    uint // stack top/capacity
	sc    uint // stack counter
	elems []T  // elements
}

// wraps fmt.Sprint(s.elems)
func (s *Stack[T]) ElemString() string { return fmt.Sprint(s.elems[:s.sc]) }

type StaticStack[T any] Stack[T]

func makeStack[T any](cap uint) (out Stack[T]) {
	if cap == 0 {
		cap = 8
	} else {
		// make capacity a power of 2
		cap = util.PowerOfTwoCeil(cap)
	}
	out.st, out.sc = cap, 0
	out.elems = make([]T, cap)
	return out
}

func NewStack[T any](cap uint) *Stack[T] {
	out := new(Stack[T])
	*out = makeStack[T](cap)
	return out
}

func (s *Stack[T]) getElemsCopy(newCap uint) []T {
	if newCap < uint(cap(s.elems)) {
		newCap = uint(cap(s.elems))
	}

	newElems := make([]T, newCap)
	for i := uint(0); i < s.sc; i++ {
		newElems[i] = s.elems[i]
	}
	return newElems
}

func (s *Stack[T]) grow() StackStatus {
	if (s.st << 1) < s.st { // check for overflow
		return Overflow
	}

	newCap := s.st << 1
	// don't want to use copy here because copy will take the len of elems, but
	// len might be greater than s.sc (s.sc is the effective length of the
	// stack)
	s.elems = s.getElemsCopy(newCap)
	s.st = newCap
	return Ok
}

func (s *Stack[T]) full() bool { return s.st == s.sc }

func (s *Stack[T]) MultiPush(elems ...T) {
	for _, elem := range elems {
		s.Push(elem)
	}
}

func (s *Stack[T]) MultiCheck(n int) (elems []T, stat StackStatus) {
	if s.sc < uint(n) {
		return nil, IllegalOperation
	}

	elems = make([]T, n)
	for i := n - 1; i >= 0; i-- {
		elems[n-1-i] = s.unsafePeekOffset(uint(i))
	}
	return elems, Ok
}

func (s *Stack[T]) MultiPop(n uint) (elems []T, stat StackStatus) {
	if s.sc < n {
		stat = IllegalOperation
		return
	}

	elems = make([]T, n)
	for i := uint(0); i < n; i++ {
		elems[i] = s.unsafePeekOffset(i)
	}
	s.sc = s.sc - n
	stat = Ok
	return
}

func (s *Stack[T]) Push(elem T) {
	if s.full() {
		if s.grow().IsOverflow() {
			panic("stack overflow")
		}
	}

	s.elems[s.sc] = elem
	s.sc++
}

func (s *Stack[T]) Clear(n uint) {
	n = util.UMin(uint(s.GetCount()), n)
	s.sc = s.sc - n
}

func (s *Stack[T]) Peek() (elem T, stat StackStatus) {
	if s.Empty() {
		stat = Empty
	} else {
		stat = Ok
		elem = s.elems[s.sc-1]
	}
	return
}

func (s *Stack[T]) unsafePeekOffset(n uint) T {
	return s.elems[s.sc-1-n]
}

func (s *Stack[T]) Pop() (T, StackStatus) {
	elem, stat := s.Peek()
	if stat.IsOk() {
		s.sc--
	}
	return elem, stat
}

func (s *Stack[T]) GetCapacity() uint {
	return s.st
}

func (s *Stack[T]) GetCount() uint {
	return s.sc
}

func (s *Stack[T]) Status() StackStatus {
	if s.Empty() {
		return Empty
	}
	return Ok
}

func (s *Stack[T]) Empty() bool {
	return s.GetCount() == 0
}
