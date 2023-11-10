package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/stringable"
)

type ExpressionJudgement[T nameable.Nameable, E expr.Expression[T]] interface {
	stringable.Stringable
	collectable[T]
	AsTypeJudgement() TypeJudgement[T, E]
	MakeJudgement(E, Type[T]) ExpressionJudgement[T, E]
}

func GetType[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E]) Type[T] {
	return j.AsTypeJudgement().ty
}

func GetExpression[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E]) E {
	return j.AsTypeJudgement().expression
}

func AsJudgement[N nameable.Nameable, E expr.Expression[N], T Type[N]](ej ExpressionJudgement[N, E]) (judgement TypedJudgement[N, E, T], success bool) {
	tj := ej.AsTypeJudgement()
	e := tj.expression
	var t T
	if t, success = tj.ty.(T); success {
		judgement = TypedJudgement[N, E, T]{e, t}
	}
	return
}

func GetExpressionAndType[N nameable.Nameable, E expr.Expression[N], T Type[N]](ej ExpressionJudgement[N, E]) (e E, t T, ok bool) {
	var judgement TypedJudgement[N, E, T]
	if judgement, ok = AsJudgement[N, E, T](ej); ok {
		e, t = judgement.expression, judgement.typing
	}
	return
}

func Replace[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E], v Variable[T], m Monotyped[T]) ExpressionJudgement[T, E] {
	ex, ty := GetExpression(j), GetType(j)
	return j.MakeJudgement(ex, MaybeReplace(ty, v, m))
}

func JudgesEquals[N nameable.Nameable, T, U expr.Expression[N]](j1 ExpressionJudgement[N, T], j2 ExpressionJudgement[N, U]) bool {
	return Equals(j1.AsTypeJudgement(), j2.AsTypeJudgement())
}

func Collect[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E]) []T {
	return j.Collect()
}

func String[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgement[T, E]) string {
	return j.String()
}
