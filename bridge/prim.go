package bridge

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

// does nothing
func (p Prim[T]) BodyAbstract(expr.Variable[T], expr.Const[T]) expr.Expression[T] { return p }

type PrimInterface[T nameable.Nameable] interface {
	FromString(string)
	Equals(PrimInterface[T]) bool
	Val() T
	GetType() types.Monotyped[T]
}

func (prim Prim[T]) Flatten() []expr.Expression[T] {
	return []expr.Expression[T]{prim}
}

func (prim Prim[T]) ToPattern() Pattern[T] {
	val := prim.Val.Val()
	ty := prim.Val.GetType()
	almost, _ := expr.MakeElem[T](expr.PatternElementLiteral, val).ToAlmostPattern()
	pat, _ := ToPattern[T](almost, ty)
	return pat
}

type Prim[T nameable.Nameable] struct {
	Val PrimInterface[T]
}

func (Prim[T]) ExtractVariables(int) []expr.Variable[T] {
	return []expr.Variable[T]{}
}

func (Prim[T]) Collect() []T {
	return []T{}
}

func (p Prim[T]) String() string {
	return p.Val.Val().GetName()
}

func (p Prim[T]) Equals(_ *expr.Context[T], e expr.Expression[T]) bool {
	p2, ok := e.(Prim[T])
	if !ok {
		return false
	}
	return p.Val.Equals(p2.Val)
}

func (p Prim[T]) StrictString() string {
	return p.String()
}

func (p Prim[T]) StrictEquals(e expr.Expression[T]) bool {
	return p.Equals(nil, e)
}

func (p Prim[T]) Replace(expr.Variable[T], expr.Expression[T]) (expr.Expression[T], bool) {
	return p, false
}

func (p Prim[T]) UpdateVars(gt int, by int) expr.Expression[T] { return p }

func (p Prim[T]) Again() (expr.Expression[T], bool) { return p, false}

func (p Prim[T]) Bind(expr.BindersOnly[T]) expr.Expression[T] { return p }

func (p Prim[T]) Find(expr.Variable[T]) bool { return false }

func (p Prim[T]) PrepareAsRHS() expr.Expression[T] { return p }

func (p Prim[T]) Rebind() expr.Expression[T] { return p }

func (p Prim[T]) Copy() expr.Expression[T] { return p }

func (p Prim[T]) ForceRequest() expr.Expression[T] { return p }