package types

import (
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type Application[T nameable.Nameable] struct {
	c Constant[T]
	ts []Monotyped[T]
}

func Apply[T nameable.Nameable](c Constant[T], ts ...Monotyped[T]) Application[T] {
	return Application[T]{ c: c, ts: ts, }
}

func (cxt *Context[T]) App(name string, ts ...Monotyped[T]) Application[T] {
	return Application[T]{ c: cxt.Con(name), ts: ts, }
}

func (a Application[T]) String() string {
	return "(" + a.c.String() + " " + str.Join(a.ts, str.String(" ")) + ")"
}

func (a Application[T]) Replace(v Variable[T], m Monotyped[T]) Monotyped[T] {
	f := func(mono Monotyped[T]) Monotyped[T] { return mono.Replace(v, m) }
	return Apply(a.c, fun.FMap(a.ts, f)...)
}

func (a Application[T]) Collect() []T {
	res := make([]T, 0, len(a.ts)+1)
	res = append(res, a.c.name)
	for _, t := range a.ts {
		res = append(res, t.Collect()...)
	}
	return res
}

func (a Application[T]) Split() (string, []Monotyped[T]) { return a.c.name.GetName(), a.ts }

func (a Application[T]) ReplaceDependent(v Variable[T], m Monotyped[T]) DependentTyped[T] {
	f := func(mono Monotyped[T]) Monotyped[T] { return mono.Replace(v, m) }
	return Apply(a.c, fun.FMap(a.ts, f)...)
}

func (a Application[T]) ReplaceKindVar(replacing Variable[T], with Monotyped[T]) Monotyped[T] {
	f := func(m Monotyped[T]) Monotyped[T] { return m.ReplaceKindVar(replacing, with) }
	return Apply(a.c, fun.FMap(a.ts, f)...)
}

func (a Application[T]) FreeInstantiation(cxt *Context[T]) DependentTyped[T] {
	f := func(m Monotyped[T]) Monotyped[T] { return m.FreeInstantiation(cxt).(Monotyped[T]) }
	return Apply(a.c, fun.FMap(a.ts, f)...)
}

func (a Application[T]) Generalize(cxt *Context[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: cxt.MakeDummyVars(1),
		bound: a,
	}
}

func (a Application[T]) Equals(t Type[T]) bool {
	a2, ok := t.(Application[T])
	if !ok {
		return false
	}

	if a.c.name.GetName() != a2.c.name.GetName() || len(a.ts) != len(a2.ts) {
		return false
	}

	for i := range a.ts {
		if !a.ts[i].Equals(a2.ts[i]) {
			return false
		}
	}
	return true
}