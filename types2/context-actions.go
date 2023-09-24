package types

import (
	"errors"
	"strconv"
)

func freeVarName(n uint32) string {
	return "$" + strconv.FormatInt(int64(n), 10)
}

func (cxt *Context) NewVar() Variable {
	n := cxt.varCounter
	cxt.varCounter++
	return Variable{
		boundContext: int32(cxt.contextNumber),
		name: freeVarName(n),
	}
}

func (cxt *Context) PopMonotype() (Monotyped, error) {
	m, ok := cxt.stack.Pop().(Monotyped)
	if !ok {
		return nil, errors.New("tried to use a non-monotype as a monotype")
	}
	return m, nil
}

func (cxt *Context) PopTypesAsPolys(n uint32) ([]Polytype, error) {
	ts := cxt.stack.MultiPop(n)
	out := make([]Polytype, len(ts))
	var ok bool
	for i, t := range ts {
		var tmp Polytype
		tmp, ok = t.(Polytype)
		if !ok {
			tmp = Polytype{
				typeBinders: nil,
				bound: t.(DependentTyped),
			}
		}
		out[i] = tmp
	}
	return out, nil
}

func PopTypes[T Type](cxt *Context, n uint32) ([]T, error) {
	ts := cxt.stack.MultiPop(n)
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

func (cxt *Context) Pop() Type {
	return cxt.stack.Pop()
}

func (cxt *Context) PopPolytype() (Polytype, error) {
	p, ok := cxt.stack.Pop().(Polytype)
	if !ok {
		return p, errors.New("tried to use a non-polytype as a polytype")
	}
	return p, nil
}

func (cxt *Context) PopDependentTyped() (DependentTyped, error) {
	d, ok := cxt.stack.Pop().(DependentTyped)
	if !ok {
		return nil, errors.New("tried to use a non-dependent-type as a dependent-type")
	}
	return d, nil
}

func (cxt *Context) Push(t Type) { cxt.stack.Push(t) }