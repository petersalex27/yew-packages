package bridge

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

type JudgmentAsExpression[T nameable.Nameable, E expr.Expression[T]] types.TypeJudgment[T, E]

func (judgment JudgmentAsExpression[T, E]) Flatten() []expr.Expression[T] {
	_, e := judgment.TypeAndExpr()
	if _, ok := expr.Expression[T](e).(expr.Variable[T]); ok {
		return []expr.Expression[T]{judgment}
	} else if _, ok := expr.Expression[T](e).(expr.Const[T]); ok {
		return []expr.Expression[T]{judgment}
	}
	return e.Flatten()
}

func (judgment JudgmentAsExpression[N, _]) MakeJudgment(expr.Expression[N], types.Type[N]) types.ExpressionJudgment[N, expr.Expression[N]] {
	return judgment
}

func (judgment JudgmentAsExpression[N, _]) GetExpressionAndType() (expr.Expression[N], types.Type[N]) {
	return judgment.ToTypeJudgment().GetExpressionAndType()
}

func (judgment JudgmentAsExpression[T, E]) ExtractVariables(gt int) []expr.Variable[T] {
	_, e := judgment.TypeAndExpr()
	return e.ExtractVariables(gt)
}

func (judgment JudgmentAsExpression[T, E]) AsTypeJudgment() types.TypeJudgment[T, expr.Expression[T]] {
	e, t := judgment.GetExpressionAndType()
	return types.Judgment(expr.Expression[T](e), t)
}

func (judgment JudgmentAsExpression[T, E]) ToTypeJudgment() types.TypeJudgment[T, E] {
	return types.TypeJudgment[T, E](judgment)
}

func Judgment[T nameable.Nameable, E expr.Expression[T]](e E, t types.Type[T]) JudgmentAsExpression[T, E] {
	return JudgmentAsExpression[T, E](types.Judgment(e, t))
}

func (judgment JudgmentAsExpression[T, E]) TypeAndExpr() (types.Type[T], E) {
	j := judgment.AsTypeJudgment()
	return j.GetType(), j.GetExpression().(E)
}

func (judgment JudgmentAsExpression[T, _]) String() string {
	return judgment.AsTypeJudgment().String()
}

func (judgment JudgmentAsExpression[T, E]) equalsHead(e expr.Expression[T]) (e1, e2 expr.Expression[T], ok bool) {
	judgment2, ok := e.(JudgmentAsExpression[T, E])
	if !ok {
		ok = false
		return
	}

	j1, j2 := judgment.AsTypeJudgment(), judgment2.AsTypeJudgment()
	return j1.GetExpression(), j2.GetExpression(), j1.GetType().Equals(j2.GetType())
}

// see implementations of
//
//	Type[T].Equals(Type[T]) bool
//
// and
//
//	Expression[T].Equals(Context[T], Expression[T]) bool
func (judgment JudgmentAsExpression[T, _]) Equals(cxt *expr.Context[T], e expr.Expression[T]) bool {
	e1, e2, ok := judgment.equalsHead(e)
	if !ok {
		return false
	}
	return e1.Equals(cxt, e2)
}

func (judgment JudgmentAsExpression[T, _]) StrictString() string {
	return judgment.AsTypeJudgment().GetExpression().StrictString()
}

func (judgment JudgmentAsExpression[T, _]) StrictEquals(e expr.Expression[T]) bool {
	e1, e2, ok := judgment.equalsHead(e)
	if !ok {
		return false
	}
	return e1.StrictEquals(e2)
}

func (judgment JudgmentAsExpression[T, _]) expressionAction(f func(expr.Expression[T]) expr.Expression[T]) JudgmentAsExpression[T, expr.Expression[T]] {
	ty, e := judgment.TypeAndExpr()
	return Judgment[T, expr.Expression[T]](f(e), ty)
}

func (judgment JudgmentAsExpression[T, _]) Replace(v expr.Variable[T], e expr.Expression[T]) (expr.Expression[T], bool) {
	var res = new(bool)
	return judgment.expressionAction(
		func(eIn expr.Expression[T]) (ex expr.Expression[T]) {
			ex, *res = eIn.Replace(v, e)
			return
		}), *res
}

func (judgment JudgmentAsExpression[T, _]) UpdateVars(gt int, by int) expr.Expression[T] {
	return judgment.expressionAction(
		func(e expr.Expression[T]) (ex expr.Expression[T]) { ex = e.UpdateVars(gt, by); return })
}

func (judgment JudgmentAsExpression[T, _]) Again() (expr.Expression[T], bool) {
	var res = new(bool)
	return judgment.expressionAction(
		func(e expr.Expression[T]) (ex expr.Expression[T]) {
			ex, *res = e.Again()
			return
		}), *res
}

func (judgment JudgmentAsExpression[T, _]) Bind(e expr.BindersOnly[T]) expr.Expression[T] {
	return judgment.expressionAction(
		func(ex expr.Expression[T]) expr.Expression[T] { return ex.Bind(e) },
	)
}

func (judgment JudgmentAsExpression[T, _]) Find(v expr.Variable[T]) bool {
	_, e := judgment.TypeAndExpr()
	return e.Find(v)
}

func (judgment JudgmentAsExpression[T, _]) PrepareAsRHS() expr.Expression[T] {
	return judgment.expressionAction(
		func(e expr.Expression[T]) expr.Expression[T] { return e.PrepareAsRHS() },
	)
}

func (judgment JudgmentAsExpression[T, _]) Rebind() expr.Expression[T] {
	return judgment.expressionAction(
		func(e expr.Expression[T]) expr.Expression[T] { return e.Rebind() },
	)
}

func (judgment JudgmentAsExpression[T, _]) Copy() expr.Expression[T] {
	return judgment.expressionAction(
		func(e expr.Expression[T]) expr.Expression[T] { return e.Copy() },
	)
}

func (judgment JudgmentAsExpression[T, _]) ForceRequest() expr.Expression[T] {
	return judgment.expressionAction(
		func(e expr.Expression[T]) expr.Expression[T] { return e.ForceRequest() },
	)
}

// collects all Ts
func (judgment JudgmentAsExpression[T, _]) Collect() (out []T) {
	ty, e := judgment.TypeAndExpr()
	out = append([]T{}, ty.Collect()...)
	out = append(out, e.Collect()...)
	return out
}

func (judgment JudgmentAsExpression[T, E]) BodyAbstract(v expr.Variable[T], name expr.Const[T]) expr.Expression[T] {
	ty, e := judgment.TypeAndExpr()
	return JudgmentAsExpression[T, E](types.Judgment(e, ty))
}
