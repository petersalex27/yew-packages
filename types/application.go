package types

import (
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type ReferableType[T nameable.Nameable] interface {
	Monotyped[T]
	GetName() T
}

type Application[T nameable.Nameable] struct {
	c  ReferableType[T]
	ts []Monotyped[T]
}

func (a Application[T]) GetName() T {
	return a.c.GetName()
}

func (a Application[T]) Merge(ms ...Monotyped[T]) Application[T] {
	return Application[T]{
		c:  a.c,
		ts: append(a.ts, ms...),
	}
}

func Apply[T nameable.Nameable](c ReferableType[T], ts ...Monotyped[T]) Application[T] {
	if app, isApp := c.(Application[T]); isApp {
		return app.Merge(ts...)
	}
	return Application[T]{c: c, ts: ts}
}

func (cxt *Context[T]) App(name string, ts ...Monotyped[T]) Application[T] {
	return Application[T]{c: cxt.Con(name), ts: ts}
}

func (a Application[T]) String() string {
	left, mid, right := "", "", ""
	lclose, rclose := "(", ")"
	if ic, ok := a.c.(InfixConst[T]); ok {
		length := len(a.ts)
		if length < 2 {
			left = ic.String()
			if length == 0 {
				return left
			} // else length == 1
			mid = " "
			right = a.ts[0].String()
		} else {
			left = a.ts[0].String()
			mid = " " + Constant[T](ic).String() + " "
			right = str.Join(a.ts[1:], str.String(" "))
		}
	} else if ec, ok := a.c.(EnclosingConst[T]); ok {
		lclose, rclose = ec.SplitString()
		mid = str.Join(a.ts, str.String(" "))
	} else {
		left = a.c.String()
		mid = " "
		right = str.Join(a.ts, str.String(" "))
	}
	return lclose + left + mid + right + rclose
}

func (a Application[T]) Replace(v Variable[T], m Monotyped[T]) Monotyped[T] {
	f := func(mono Monotyped[T]) Monotyped[T] { return mono.Replace(v, m) }
	res := a.c.Replace(v, m)
	left, ok := res.(ReferableType[T])
	right := fun.FMap(a.ts, f)
	if !ok { 
		// this branch cannot be entered with current set-up since all monotypes are 
		// also referable types
		return Apply(a.c, right...)
	}
	return Apply(left, right...)
}

func (a Application[T]) Collect() []T {
	res := make([]T, 0, len(a.ts)+1)
	res = append(res, a.c.GetName())
	for _, t := range a.ts {
		res = append(res, t.Collect()...)
	}
	return res
}

func (a Application[T]) Split() (string, []Monotyped[T]) { return a.c.GetName().GetName(), a.ts }

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
		bound:       a,
	}
}

func (a Application[T]) Equals(t Type[T]) bool {
	a2, ok := t.(Application[T])
	if !ok {
		return false
	}

	if a.c.GetName().GetName() != a2.c.GetName().GetName() || len(a.ts) != len(a2.ts) {
		return false
	}

	for i := range a.ts {
		if !a.ts[i].Equals(a2.ts[i]) {
			return false
		}
	}
	return true
}
