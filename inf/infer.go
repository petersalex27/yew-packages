// =============================================================================
// Author-Date: Alex Peters - 2023
//
// Content: 
// contains type inference rules
//
// Notes: -
// =============================================================================
package inf

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/types"
)

// [Index] rule:
//		x: δ ∈ 𝚪    x0: t0 ∈    t = Inst(δ)    𝚪 ⊢ p1: t1 .. 𝚪 ⊢ pN: tN
//		𝚪, x: mapval (x0: t0') .. (xN: tN') . (t''; p0 p1 .. pN)    𝚪, e0: t ⊢ e1: t1
//		--------------------------------------------- [Index]
//		                𝚪 ⊢ e1: t1

// [Var] rule:
//		x: σ ∈ 𝚪    t = Inst(σ)
//    ----------------------- [Var]
//          𝚪 ⊢ x: t
func (cxt *Context[T]) Var(x bridge.JudgementAsExpression[T, expr.Const[T]]) bridge.JudgementAsExpression[T, expr.Const[T]] {
	var t types.Monotyped[T]

	tmp, xConst := x.TypeAndExpr()

	// grab polytype
	sigma, ok := tmp.(types.Polytype[T])
	if !ok { // still technically a polytype, just one w/ no zero binders, so make that explicit
		// all types that aren't polytypes, are dependent types, so assertion will pass
		dep, _ := tmp.(types.DependentTyped[T]) 
		sigma = types.Forall[T]().Bind(dep)
	}

	// replace all bound (including kind-) variables with free variables
	t = sigma.Specialize(cxt.typeContext)
	// return judgement `x: t`
	return bridge.Judgement(xConst, types.Type[T](t))
}

// [App] rule:
//		𝚪 ⊢ e0: t0    𝚪 ⊢ e1: t1    t2 = Apply(t0, t1)
//    ---------------------------------------------- [App]
//		             𝚪 ⊢ (e0 e1): t2
//
// applies j0 and j1 resulting in a type t2 and the implication that
// 		t0 = t1 -> t2
//
// the *magic* of this rule comes from the new equation which provides more
// information about type t0
//
// curry-howard: conditional elim
func (cxt *Context[T]) App(j0, j1 bridge.JudgementAsExpression[T, expr.Expression[T]]) bridge.JudgementAsExpression[T, expr.Application[T]] {
	// split judgements into types and expressions
	tmp0, e0 := j0.TypeAndExpr()
	tmp1, e1 := j1.TypeAndExpr()
	// get monotypes
	t0 := tmp0.(types.Monotyped[T])
	t1 := tmp1.(types.Monotyped[T])
	// "t2 = Apply(t0, t1)" premise of rule
	appliedType, e := cxt.typeContext.Apply(t0, t1)
	if e != nil {
		cxt.es = append(cxt.es, e) // TODO: report how and when?
	}
	// "(e0 e1)" in result of rule
	appliedExpression := expr.Apply(e0, e1)
	t2 := types.Type[T](appliedType)
	return bridge.Judgement(appliedExpression, t2)
}

// [Abs] rule:
//		t0 = newvar    𝚪, param: t0 ⊢ e: t1
//		-----------------------------------
//		    𝚪 ⊢ (λparam . e): t0 -> t1
//
// notice that the second param adds context and the third premise no longer
// has that context
//
// curry-howard: conditional intro
func (cxt *Context[T]) Abs(param T) func(bridge.JudgementAsExpression[T, expr.Expression[T]]) bridge.JudgementAsExpression[T, expr.Function[T]] {
	// first, add context (this is the first premise)
	paramConst := expr.Const[T]{Name: param}
	t0 := cxt.typeContext.NewVar()
	// grow context w/ type judgement `param: t0`
	cxt.AddWithType(paramConst, t0)

	// now, return function to allow second premise of Abs when needed
	return func(j bridge.JudgementAsExpression[T, expr.Expression[T]]) bridge.JudgementAsExpression[T, expr.Function[T]] {
		// remove context added
		cxt.Remove(paramConst)

		// split judgement
		tmp1, e := j.TypeAndExpr()
		t1 := tmp1.(types.Monotyped[T])

		// create function body by converting param-name to param-var in e
		v := cxt.exprContext.NewVar()
		e = e.BodyAbstract(v, paramConst)

		// actual function creation, finish abstraction of `e`
		f := expr.Bind(v).In(e)

		// create function type
		var fnType types.Type[T] = cxt.typeContext.Function(t0, t1)

		// last line of rule: `(λparam . e): t0 -> t1`
		return bridge.Judgement(f, fnType)
	}
}

// [Let] rule:
//		𝚪 ⊢ e0: t0    𝚪, name: Gen(t) ⊢ e1: t1
//		-------------------------------------- [Let]
//		     𝚪 ⊢ let name = e0 in e1: t1
//
// notice that the second param adds context and the third premise no longer
// has that context
//
// This rule allows for a kind of polymorphism. Here's an example given 
// 	𝚪 = {0: Int, (λx.x): a -> a}:
//
//		                           [ f: forall a. a -> a ]¹    Inst(forall a. a -> a)  
//		                           -------------------------------------------------- [Var]
//		                   0: Int                      f: x -> x                       t0, Int = x
//		                   ----------------------------------------------------------------------- [App]
//		                                               f 0: t0
//		                                               ------- [Id]
//		  (λx.x): a -> a                               f 0: Int
// 		1 ----------------------------------------------------- [Let]
//		               let f = (λx.x) in f 0: Int
func (cxt *Context[T]) Let(name T, j0 bridge.JudgementAsExpression[T, expr.Expression[T]]) func (bridge.JudgementAsExpression[T, expr.Expression[T]]) bridge.JudgementAsExpression[T, expr.NameContext[T]] {
	nameConst := expr.Const[T]{Name: name}
	tmp0, e0 := j0.TypeAndExpr()
	t0 := tmp0.(types.Monotyped[T])
	generalized_t0 := cxt.Gen(t0)
	cxt.AddWithType(nameConst, generalized_t0)

	return func(j1 bridge.JudgementAsExpression[T, expr.Expression[T]]) bridge.JudgementAsExpression[T, expr.NameContext[T]] {
		cxt.Remove(nameConst)

		t1, e1 := j1.TypeAndExpr()
		let := expr.Let(nameConst, e0, e1)
		return bridge.Judgement(let, t1)
	}
}

// generalizes a type: binds all free variables w/in monotype
func (cxt *Context[T]) Gen(ty types.Monotyped[T]) types.Polytype[T] {
	/*var res types.DependentTyped[T] 
	if dti, ok := ty.(types.DependentTypeInstance[T]); ok {
		dti.GetFreeVariables()
	}*/
	vs := ty.GetFreeVariables()
	return types.Forall(vs...).Bind(ty)
}