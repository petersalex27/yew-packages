package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
)

type TypeJudgment[T nameable.Nameable, E expr.Expression[T]] struct {
	expression E
	ty         Type[T]
}

func (j TypeJudgment[T, E]) AsTypeJudgment() TypeJudgment[T, E] {
	return j
}

func (j TypeJudgment[T, _]) GetType() Type[T] {
	return j.ty
}

func (j TypeJudgment[_, E]) GetExpression() E {
	return j.expression
}

func (j TypeJudgment[N, _]) GetExpressionAndType() (expr.Expression[N], Type[N]) {
	return j.expression, j.ty
}

func (j TypeJudgment[T, _]) String() string {
	return "(" + j.expression.String() + ": " + j.ty.String() + ")"
}

// Judgment makes the trivial type judgment `e: ty ⊢ e: ty`
func Judgment[T nameable.Nameable, E expr.Expression[T]](e E, ty Type[T]) TypeJudgment[T, E] {
	return TypeJudgment[T, E]{
		expression: e,
		ty:         ty,
	}
}

func (TypeJudgment[T, E]) MakeJudgment(e E, ty Type[T]) ExpressionJudgment[T, E] {
	return Judgment[T, E](e, ty)
}

func (j TypeJudgment[T, E]) Replace(v Variable[T], m Monotyped[T]) TypeJudgment[T, E] {
	return Judgment(j.expression, MaybeReplace(j.ty, v, m))
}

func Equals[N nameable.Nameable, T, U expr.Expression[N]](j1 TypeJudgment[N, T], j2 TypeJudgment[N, U]) bool {
	return j1.ty.Equals(j2.ty) && j1.expression.StrictEquals(j2.expression)
	// TODO: ??? j1.expression.Equals(j2.expression) instead???
}

func (j TypeJudgment[T, E]) Collect() []T {
	return append(j.expression.Collect(), j.ty.Collect()...)
}
