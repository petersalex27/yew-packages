package types

import (
	"errors"
	"strconv"

	"github.com/petersalex27/yew-packages/nameable"
)

func freeVarName(n uint32) string {
	return "$" + strconv.FormatInt(int64(n), 10)
}

func (cxt *Context[T]) NewVar() Variable[T] {
	n := cxt.varCounter
	cxt.varCounter++
	return cxt.Var(freeVarName(n)).BoundIn(int32(cxt.contextNumber))
}

func (cxt *Context[T]) PopMonotype() (Monotyped[T], error) {
	m, ok := cxt.Pop().(Monotyped[T])
	if !ok {
		return nil, errors.New("tried to use a non-monotype as a monotype")
	}
	return m, nil
}

func (cxt *Context[T]) PopTypesAsPolys(n uint) ([]Polytype[T], error) {
	ts, _ := cxt.stack.MultiPop(n)
	out := make([]Polytype[T], len(ts))
	var ok bool
	for i, t := range ts {
		var tmp Polytype[T]
		tmp, ok = t.(Polytype[T])
		if !ok {
			tmp = Polytype[T]{
				typeBinders: nil,
				bound: t.(DependentTyped[T]),
			}
		}
		out[i] = tmp
	}
	return out, nil
}

func PopTypes[T Type[U], U nameable.Nameable](cxt *Context[U], n uint) ([]T, error) {
	ts, _ := cxt.stack.MultiPop(n)
	out := make([]T, len(ts))
	var ok bool
	for i, t := range ts {
		out[i], ok = t.(T)
		if !ok {
			return nil, errors.New("failed to assert type")
		}
	}
	return out, nil
}

func (cxt *Context[T]) Pop() Type[T] {
	ty, stat := cxt.stack.Pop()
	if stat.NotOk() {
		panic("bug: empty stack")
	}
	return ty
}

func (cxt *Context[T]) PopPolytype() (Polytype[T], error) {
	p, ok := cxt.Pop().(Polytype[T])
	if !ok {
		return p, errors.New("tried to use a non-polytype as a polytype")
	}
	return p, nil
}

func (cxt *Context[T]) PopDependentTyped() (DependentTyped[T], error) {
	d, ok := cxt.Pop().(DependentTyped[T])
	if !ok {
		return nil, errors.New("tried to use a non-dependent-type as a dependent-type")
	}
	return d, nil
}

func (cxt *Context[T]) Push(t Type[T]) { cxt.stack.Push(t) }