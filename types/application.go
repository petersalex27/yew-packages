// =============================================================================
// Author-Date: Alex Peters - 2023
//
// Content:
// Application type and its methods
// =============================================================================
package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type Application[T nameable.Nameable] struct {
	c  Monotyped[T]
	ts []Monotyped[T]
}

func (a Application[N]) Rebuild(findMono func(Monotyped[N]) Monotyped[N], _ func(expr.Referable[N]) expr.Referable[N]) TypeFunction[N] {
	return Application[N]{
		findMono(a.c),
		fun.FMap(a.ts, findMono),
	}
}

func (a Application[T]) FunctionAndIndexes() (function Application[T], indexes Indexes[T]) {
	return a, nil
}

func (a Application[T]) SubVars(preSub []TypeJudgment[T, expr.Variable[T]], postSub []expr.Referable[T]) TypeFunction[T] {
	return Application[T]{
		a.c,
		fun.FMap(
			a.ts,
			func(m Monotyped[T]) Monotyped[T] {
				if f, ok := m.(TypeFunction[T]); ok {
					return f.SubVars(preSub, postSub)
				}
				return m
			},
		),
	}
}

func (a Application[T]) GetFreeVariables() []Variable[T] {
	leftVars := a.c.GetFreeVariables()
	// get 2d slice of free vars from rhs
	rightVars2d := fun.FMap(
		a.ts,
		func(t Monotyped[T]) []Variable[T] {
			return t.GetFreeVariables()
		},
	)
	// flatten 2d slice
	rightVars := fun.FoldLeft(
		[]Variable[T]{},
		rightVars2d,
		func(vs, us []Variable[T]) []Variable[T] {
			return append(vs, us...)
		},
	)
	return append(leftVars, rightVars...)
}

func (a Application[T]) GetReferred() T {
	return a.c.GetReferred()
}

func (a Application[T]) Merge(ms ...Monotyped[T]) Application[T] {
	return Application[T]{
		c:  a.c,
		ts: append(a.ts, ms...),
	}
}

func Apply[T nameable.Nameable](c Monotyped[T], ts ...Monotyped[T]) Application[T] {
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
	left := a.c.Replace(v, m)
	right := fun.FMap(a.ts, f)
	return Apply(left, right...)
}

func (a Application[T]) Collect() []T {
	res := make([]T, 0, len(a.ts)+1)
	res = append(res, a.c.GetReferred())
	for _, t := range a.ts {
		res = append(res, t.Collect()...)
	}
	return res
}

func (a Application[T]) Split() (name string, params []Monotyped[T]) {
	return a.c.GetReferred().GetName(), a.ts
}

func (a Application[T]) ReplaceDependent(vs []Variable[T], ms []Monotyped[T]) Monotyped[T] {
	f := func(mono Monotyped[T]) Monotyped[T] { return mono.ReplaceDependent(vs, ms) }
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

	if a.c.GetReferred().GetName() != a2.c.GetReferred().GetName() || len(a.ts) != len(a2.ts) {
		return false
	}

	for i := range a.ts {
		if !a.ts[i].Equals(a2.ts[i]) {
			return false
		}
	}
	return true
}
