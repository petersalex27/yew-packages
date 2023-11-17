package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/stringable"
)

type Indexes[T nameable.Nameable] []ExpressionJudgement[T, expr.Referable[T]]

func (idxs Indexes[T]) GetFreeVariables() []Variable[T] {
	// 2d slice of free variables
	vs2d := fun.FMap(
		idxs,
		func(ej ExpressionJudgement[T, expr.Referable[T]]) []Variable[T] {
			mono := ej.AsTypeJudgement().ty.(Monotyped[T]) // should always pass b/c indexes's values must be typed by monotypes
			return mono.GetFreeVariables()
		},
	)
	// convert to 1d slice of free variables
	return fun.FoldLeft(
		[]Variable[T]{},
		vs2d,
		func(vs, us []Variable[T]) []Variable[T] {
			return append(vs, us...)
		},
	)
}

func (idxs Indexes[T]) String() string {
	if len(idxs) == 0 {
		return ""
	}
	return "; " + stringable.Join(idxs, stringable.String(" "))
}
