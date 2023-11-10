package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
)

type FreeJudgement[T nameable.Nameable, E expr.Expression[T]] TypeJudgement[T, E]

func (j FreeJudgement[T, E]) AsTypeJudgement() TypeJudgement[T, E] {
	return TypeJudgement[T, E](j)
}

func (j FreeJudgement[_,_]) String() string {
	return j.expression.String()
}

func (FreeJudgement[T, E]) MakeJudgement(e E, ty Type[T]) ExpressionJudgement[T, E] {
	return FreeJudgement[T,E](Judgement(e, ty))
}

// Judgement makes the trivial type judgement `𝚪, e: ty ⊢ e: ty`
func FreeJudge[T nameable.Nameable, E expr.Expression[T]](cxt *Context[T], e E) FreeJudgement[T,E] {
	return FreeJudgement[T,E]{
		expression: e,
		ty: cxt.NewVar(),
	}
}

func (j FreeJudgement[T, E]) Collect() []T {
	return j.expression.Collect()
}