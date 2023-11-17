package bridge

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

type JudgementAsExpression[T nameable.Nameable, E expr.Expression[T]] types.TypeJudgement[T, E]

func (judgement JudgementAsExpression[T, E]) Flatten() []expr.Expression[T] {
	_, e := judgement.TypeAndExpr()
	if _, ok := expr.Expression[T](e).(expr.Variable[T]); ok {
		return []expr.Expression[T]{judgement}
	} else if _, ok := expr.Expression[T](e).(expr.Const[T]); ok {
		return []expr.Expression[T]{judgement}
	}
	return e.Flatten()
}

func (judgement JudgementAsExpression[N, _]) MakeJudgement(expr.Expression[N], types.Type[N]) types.ExpressionJudgement[N, expr.Expression[N]] {
	return judgement
}

func (judgement JudgementAsExpression[N, _]) GetExpressionAndType() (expr.Expression[N], types.Type[N]) {
	return judgement.ToTypeJudgement().GetExpressionAndType()
}

func (judgement JudgementAsExpression[T, E]) ExtractVariables(gt int) []expr.Variable[T] {
	_, e := judgement.TypeAndExpr()
	return e.ExtractVariables(gt)
}

func (judgement JudgementAsExpression[T, E]) AsTypeJudgement() types.TypeJudgement[T, expr.Expression[T]] {
	e, t := judgement.GetExpressionAndType()
	return types.Judgement(expr.Expression[T](e), t)
}

func (judgement JudgementAsExpression[T, E]) ToTypeJudgement() types.TypeJudgement[T, E] {
	return types.TypeJudgement[T, E](judgement)
}

func Judgement[T nameable.Nameable, E expr.Expression[T]](e E, t types.Type[T]) JudgementAsExpression[T, E] {
	return JudgementAsExpression[T, E](types.Judgement(e, t))
}

func (judgement JudgementAsExpression[T, E]) TypeAndExpr() (types.Type[T], E) {
	j := judgement.AsTypeJudgement()
	return j.GetType(), j.GetExpression().(E)
}

func (judgement JudgementAsExpression[T, _]) String() string {
	return judgement.AsTypeJudgement().String()
}

func (judgement JudgementAsExpression[T, E]) equalsHead(e expr.Expression[T]) (e1, e2 expr.Expression[T], ok bool) {
	judgement2, ok := e.(JudgementAsExpression[T, E])
	if !ok {
		ok = false
		return
	}

	j1, j2 := judgement.AsTypeJudgement(), judgement2.AsTypeJudgement()
	return j1.GetExpression(), j2.GetExpression(), j1.GetType().Equals(j2.GetType())
}

// see implementations of
//
//	Type[T].Equals(Type[T]) bool
//
// and
//
//	Expression[T].Equals(Context[T], Expression[T]) bool
func (judgement JudgementAsExpression[T, _]) Equals(cxt *expr.Context[T], e expr.Expression[T]) bool {
	e1, e2, ok := judgement.equalsHead(e)
	if !ok {
		return false
	}
	return e1.Equals(cxt, e2)
}

func (judgement JudgementAsExpression[T, _]) StrictString() string {
	return judgement.AsTypeJudgement().GetExpression().StrictString()
}

func (judgement JudgementAsExpression[T, _]) StrictEquals(e expr.Expression[T]) bool {
	e1, e2, ok := judgement.equalsHead(e)
	if !ok {
		return false
	}
	return e1.StrictEquals(e2)
}

func (judgement JudgementAsExpression[T, _]) expressionAction(f func(expr.Expression[T]) expr.Expression[T]) JudgementAsExpression[T, expr.Expression[T]] {
	ty, e := judgement.TypeAndExpr()
	return Judgement[T, expr.Expression[T]](f(e), ty)
}

func (judgement JudgementAsExpression[T, _]) Replace(v expr.Variable[T], e expr.Expression[T]) (expr.Expression[T], bool) {
	var res = new(bool)
	return judgement.expressionAction(
		func(eIn expr.Expression[T]) (ex expr.Expression[T]) {
			ex, *res = eIn.Replace(v, e)
			return
		}), *res
}

func (judgement JudgementAsExpression[T, _]) UpdateVars(gt int, by int) expr.Expression[T] {
	return judgement.expressionAction(
		func(e expr.Expression[T]) (ex expr.Expression[T]) { ex = e.UpdateVars(gt, by); return })
}

func (judgement JudgementAsExpression[T, _]) Again() (expr.Expression[T], bool) {
	var res = new(bool)
	return judgement.expressionAction(
		func(e expr.Expression[T]) (ex expr.Expression[T]) {
			ex, *res = e.Again()
			return
		}), *res
}

func (judgement JudgementAsExpression[T, _]) Bind(e expr.BindersOnly[T]) expr.Expression[T] {
	return judgement.expressionAction(
		func(ex expr.Expression[T]) expr.Expression[T] { return ex.Bind(e) },
	)
}

func (judgement JudgementAsExpression[T, _]) Find(v expr.Variable[T]) bool {
	_, e := judgement.TypeAndExpr()
	return e.Find(v)
}

func (judgement JudgementAsExpression[T, _]) PrepareAsRHS() expr.Expression[T] {
	return judgement.expressionAction(
		func(e expr.Expression[T]) expr.Expression[T] { return e.PrepareAsRHS() },
	)
}

func (judgement JudgementAsExpression[T, _]) Rebind() expr.Expression[T] {
	return judgement.expressionAction(
		func(e expr.Expression[T]) expr.Expression[T] { return e.Rebind() },
	)
}

func (judgement JudgementAsExpression[T, _]) Copy() expr.Expression[T] {
	return judgement.expressionAction(
		func(e expr.Expression[T]) expr.Expression[T] { return e.Copy() },
	)
}

func (judgement JudgementAsExpression[T, _]) ForceRequest() expr.Expression[T] {
	return judgement.expressionAction(
		func(e expr.Expression[T]) expr.Expression[T] { return e.ForceRequest() },
	)
}

// collects all Ts
func (judgement JudgementAsExpression[T, _]) Collect() (out []T) {
	ty, e := judgement.TypeAndExpr()
	out = append([]T{}, ty.Collect()...)
	out = append(out, e.Collect()...)
	return out
}

func (judgement JudgementAsExpression[T, E]) BodyAbstract(v expr.Variable[T], name expr.Const[T]) expr.Expression[T] {
	ty, e := judgement.TypeAndExpr()
	return JudgementAsExpression[T,E](types.Judgement(e, ty))
}