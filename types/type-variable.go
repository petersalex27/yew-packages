package types

import (
	"strconv"

	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
)

type Variable[T nameable.Nameable] struct {
	boundContext int32
	name         T
}

func (v Variable[T]) GetFreeVariables() []Variable[T] {
	return []Variable[T]{v}
}

func (v Variable[T]) GetReferred() T {
	return v.name
}

func (v Variable[T]) GetName() string {
	return v.name.GetName()
}

// ReplaceKindVar implements Monotyped[T].
func (v Variable[T]) ReplaceKindVar(replacing Variable[T], with Monotyped[T]) Monotyped[T] {
	if varEquals(v, replacing) {
		return with
	}
	return v
}

func NonBindableVar[T nameable.Nameable](name T) Variable[T] {
	return Variable[T]{boundContext: -1, name: name}
}

func (v Variable[T]) BoundIn(i int32) Variable[T] {
	v.boundContext = i
	return v
}

func Var[T nameable.Nameable](name T) Variable[T] {
	return Variable[T]{boundContext: 0, name: name}
}

func (cxt *Context[T]) Var(name string) Variable[T] {
	return Variable[T]{boundContext: 0, name: cxt.makeName(name)}
}

func (cxt *Context[T]) FreeVar(name string) Variable[T] {
	return cxt.Var(name)
}

func (cxt *Context[T]) dummyName(Variable[T]) Variable[T] { return cxt.Var("_") }

// MakeDummyVars(n) = []Variable{v /*0*/, v /*1*/, .., v /*n-1*/} where
// v = Var("_")
func (cxt *Context[T]) MakeDummyVars(n uint) []Variable[T] {
	return fun.FMap(make([]Variable[T], n), cxt.dummyName)
}

// Var("a").BoundBy(x).String() = "a"
func (v Variable[T]) String() string {
	return v.name.GetName()
}

func (v Variable[T]) String2() (string, string) {
	return v.String(), "#(" + strconv.Itoa(int(v.boundContext)) + ")"
}

func (v Variable[T]) ExtendedString() string {
	l, r := v.String2()
	return l + r
}

// v.Capture(_) = v
func (v Variable[T]) Capture(w Variable[T]) Monotyped[T] { return v }

// Replace all `w` in `v` with `m`; because v is just a variable,
// if v == w, then return m; else return v. Formally, `v [w := m]`
func (v Variable[T]) Replace(w Variable[T], m Monotyped[T]) Monotyped[T] {
	if varEquals(v, w) {
		return m
	}
	return v
}

func (v Variable[T]) ReplaceDependent(ws []Variable[T], ms []Monotyped[T]) Monotyped[T] {
	for i, w := range ws {
		if varEquals(v, w) {
			return ms[i]
		}
	}
	return v
}

func (v Variable[T]) FreeInstantiation(cxt *Context[T]) Monotyped[T] {
	if v.boundContext > 0 {
		return cxt.dummyName(Variable[T]{})
	}
	return v
}

func (v Variable[T]) Collect() []T {
	return []T{v.name}
}

func varEquals[T nameable.Nameable](v Variable[T], w Variable[T]) bool {
	return v.boundContext == w.boundContext && v.name.GetName() == w.name.GetName()
}

func (v Variable[T]) Equals(t Type[T]) bool {
	v2, ok := t.(Variable[T])
	return ok && varEquals(v, v2)
}
