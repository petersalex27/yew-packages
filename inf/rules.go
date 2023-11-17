package inf

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

// [Var] rule:
//
//			x: σ ∈ 𝚪    t = Inst(σ)
//	   ----------------------- [Var]
//	         𝚪 ⊢ x: t
func (cxt *Context[N]) varBody(x bridge.JudgementAsExpression[N, expr.Const[N]]) Conclusion[N, expr.Const[N], types.Monotyped[N]] {
	var t types.Monotyped[N]

	tmp, xConst := x.TypeAndExpr()

	// grab polytype
	sigma, ok := tmp.(types.Polytype[N])
	if !ok { // still technically a polytype, just one w/ no zero binders, so make that explicit
		// all types that aren't polytypes, are dependent types, so assertion will pass
		dep, _ := tmp.(types.DependentTyped[N])
		sigma = types.Forall[N]().Bind(dep)
	}

	// replace all bound (including kind-) variables with free variables
	t = cxt.Inst(sigma)
	// return judgement `x: t`
	return Conclude[N](xConst, t)
}

func (cxt *Context[N]) Var(x expr.Const[N]) Conclusion[N, expr.Const[N], types.Monotyped[N]] {
	xJudge, found := cxt.Get(x)
	if !found {
		// `x` is not in the context
		cxt.appendReport(makeNameReport("Var", NameNotInContext, x))
		return CannotConclude[N, expr.Const[N], types.Monotyped[N]](NameNotInContext)
	}

	return cxt.varBody(xJudge)
}

// [App] rule:
//
//			𝚪 ⊢ e0: t0    𝚪 ⊢ e1: t1    t2 = newvar    t0 = t1 -> t2
//	   -------------------------------------------------------- [App]
//			                     𝚪 ⊢ (e0 e1): t2
//
// applies j0 and j1 resulting in a type t2 and the implication that
//
//	t0 = t1 -> t2
//
// the *magic* of this rule comes from the new equation which provides more
// information about type t0
//
// curry-howard: conditional elim
func (cxt *Context[N]) App(j0, j1 TypeJudgement[N]) Conclusion[N, expr.Application[N], types.Monotyped[N]] {
	// split judgements into types and expressions
	e0, tmp0 := j0.GetExpressionAndType()
	e1, tmp1 := j1.GetExpressionAndType()
	// get monotypes
	t0 := tmp0.(types.Monotyped[N])
	t1 := tmp1.(types.Monotyped[N])
	// premise `t2 = newvar`
	t2 := cxt.typeContext.NewVar()
	// create monotype `t1 -> t2`
	t1_to_t2 := cxt.typeContext.Function(t1, t2)
	// premise `t0 = t1 -> t2`
	stat := cxt.Unify(t0, t1_to_t2)
	if stat.NotOk() {
		terms := []TypeJudgement[N]{j0, j1}
		report := makeReport("App", stat, terms...)
		cxt.appendReport(report)
		return CannotConclude[N, expr.Application[N], types.Monotyped[N]](stat)
	}
	// "(e0 e1)" in result of rule
	appliedExpression := expr.Apply(e0, e1)
	// (e0 e1): t2
	return Conclude[N](appliedExpression, cxt.GetSub(t2))
}

// [Abs] rule:
//
//	t0 = newvar    𝚪, param: t0 ⊢ e: t1
//	-----------------------------------
//	    𝚪 ⊢ (λparam . e): t0 -> t1
//
// notice that the second param adds context and the third premise no longer
// has that context
//
// curry-howard: conditional intro
func (cxt *Context[N]) Abs(param N) func(TypeJudgement[N]) Conclusion[N, expr.Function[N], types.Monotyped[N]] {
	// first, add context (this is the first premise)
	paramConst := expr.Const[N]{Name: param}
	t0 := cxt.typeContext.NewVar()
	// grow context w/ type judgement `param: t0`
	cxt.Add(paramConst, t0)

	// now, return function to allow second premise of Abs when needed
	return func(j TypeJudgement[N]) Conclusion[N, expr.Function[N], types.Monotyped[N]] {
		// remove context added
		cxt.Remove(paramConst)

		// split judgement
		e, tmp1 := j.GetExpressionAndType()
		t1 := tmp1.(types.Monotyped[N])

		// create function body by converting param-name to param-var in e
		v := cxt.exprContext.NewVar()
		e = e.BodyAbstract(v, paramConst)

		// actual function creation, finish abstraction of `e`
		f := expr.Bind(v).In(e)

		// create function type
		var fnType types.Monotyped[N] = cxt.typeContext.Function(t0, t1)

		// last line of rule: `(λparam . e): t0 -> t1`
		return Conclude[N](f, fnType)
	}
}

type letAssumptionDischarge[N nameable.Nameable] func(TypeJudgement[N]) Conclusion[N, expr.NameContext[N], types.Monotyped[N]]

// [Let] rule:
//
//	𝚪 ⊢ e0: t0    𝚪, name: Gen(t) ⊢ e1: t1
//	-------------------------------------- [Let]
//	     𝚪 ⊢ let name = e0 in e1: t1
//
// notice that the second param adds context and the third premise no longer
// has that context
//
// This rule allows for a kind of polymorphism. Here's an example given
//
//	𝚪 = {0: Int, (λx.x): a -> a}:
//
//		  [ x: forall a. a -> a ]¹    Inst(forall a. a -> a)
//		  -------------------------------------------------- [Var]
//		                      x: v -> v                       0: Int    t0, Int = v
//		                      ----------------------------------------------------- [App]
//		                                              x 0: t0
//		                                              -------- [Id]
//		  (λy.y): a -> a                              x 0: Int
//		1 ---------------------------------------------------- [Let]
//		               let x = (λy.y) in x 0: Int
func (cxt *Context[N]) Let(name N, j0 TypeJudgement[N]) letAssumptionDischarge[N] {
	nameConst := expr.Const[N]{Name: name}
	e0, tmp0 := j0.GetExpressionAndType()
	t0 := tmp0.(types.Monotyped[N])
	generalized_t0 := cxt.Gen(t0)
	cxt.Add(nameConst, generalized_t0)

	return func(j1 TypeJudgement[N]) Conclusion[N, expr.NameContext[N], types.Monotyped[N]] {
		cxt.Remove(nameConst)

		e1, t1 := j1.GetExpressionAndType()
		mono := t1.(types.Monotyped[N])
		let := expr.Let(nameConst, e0, e1)
		return Conclude[N](let, mono)
	}
}