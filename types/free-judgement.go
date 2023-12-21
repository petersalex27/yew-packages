package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
)

type FreeJudgment[T nameable.Nameable, E expr.Expression[T]] TypeJudgment[T, E]

func (j FreeJudgment[T, E]) AsTypeJudgment() TypeJudgment[T, E] {
	return TypeJudgment[T, E](j)
}

func (j FreeJudgment[_, _]) String() string {
	return j.expression.String()
}

func (FreeJudgment[T, E]) MakeJudgment(e E, ty Type[T]) ExpressionJudgment[T, E] {
	return FreeJudgment[T, E](Judgment(e, ty))
}

// Judgment makes the trivial type judgment `𝚪, e: ty ⊢ e: ty`
func FreeJudge[T nameable.Nameable, E expr.Expression[T]](cxt *Context[T], e E) FreeJudgment[T, E] {
	return FreeJudgment[T, E]{
		expression: e,
		ty:         cxt.NewVar(),
	}
}

func (j FreeJudgment[T, E]) Collect() []T {
	return j.expression.Collect()
}
