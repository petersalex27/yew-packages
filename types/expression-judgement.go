package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/stringable"
)

type ExpressionJudgment[T nameable.Nameable, E expr.Expression[T]] interface {
	stringable.Stringable
	nameable.Collectable[T]
	AsTypeJudgment() TypeJudgment[T, E]
	MakeJudgment(E, Type[T]) ExpressionJudgment[T, E]
}

func GetType[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgment[T, E]) Type[T] {
	return j.AsTypeJudgment().ty
}

func GetExpression[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgment[T, E]) E {
	return j.AsTypeJudgment().expression
}

func AsJudgment[N nameable.Nameable, E expr.Expression[N], T Type[N]](ej ExpressionJudgment[N, E]) (judgment TypedJudgment[N, E, T], success bool) {
	tj := ej.AsTypeJudgment()
	e := tj.expression
	var t T
	if t, success = tj.ty.(T); success {
		judgment = TypedJudgment[N, E, T]{e, t}
	}
	return
}

func GetExpressionAndType[N nameable.Nameable, E expr.Expression[N], T Type[N]](ej ExpressionJudgment[N, E]) (e E, t T, ok bool) {
	var judgment TypedJudgment[N, E, T]
	if judgment, ok = AsJudgment[N, E, T](ej); ok {
		e, t = judgment.expression, judgment.typing
	}
	return
}

func Replace[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgment[T, E], v Variable[T], m Monotyped[T]) ExpressionJudgment[T, E] {
	ex, ty := GetExpression(j), GetType(j)
	return j.MakeJudgment(ex, MaybeReplace(ty, v, m))
}

func JudgesEquals[N nameable.Nameable, T, U expr.Expression[N]](j1 ExpressionJudgment[N, T], j2 ExpressionJudgment[N, U]) bool {
	return Equals(j1.AsTypeJudgment(), j2.AsTypeJudgment())
}

func Collect[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgment[T, E]) []T {
	return j.Collect()
}

func String[T nameable.Nameable, E expr.Expression[T]](j ExpressionJudgment[T, E]) string {
	return j.String()
}
