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
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

func (cxt *Context[T]) Inst(sigma types.Polytype[T]) types.Monotyped[T] {
	var t types.DependentTyped[T] = sigma.GetBound()
	typeVars := sigma.GetBinders()

	// create new type variables
	vs := fun.FMap(
		typeVars,
		func(v types.Variable[T]) types.Monotyped[T] {
			return cxt.typeContext.NewVar()
		},
	)

	if d, ok := t.(types.DependentType[T]); ok {
		// replace all bound expression variables w/ new expression variables
		t = d.FreeIndex(cxt.exprContext)
	}

	// replace all bound variables w/ newly created type variables
	return t.ReplaceDependent(typeVars, vs)
}

// [Var] rule:
//
//			x: σ ∈ 𝚪    t = Inst(σ)
//	   ----------------------- [Var]
//	         𝚪 ⊢ x: t
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
	t = cxt.Inst(sigma)
	// return judgement `x: t`
	return bridge.Judgement(xConst, types.Type[T](t))
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
func (cxt *Context[T]) App(j0, j1 bridge.JudgementAsExpression[T, expr.Expression[T]]) bridge.JudgementAsExpression[T, expr.Application[T]] {
	// split judgements into types and expressions
	tmp0, e0 := j0.TypeAndExpr()
	tmp1, e1 := j1.TypeAndExpr()
	// get monotypes
	t0 := tmp0.(types.Monotyped[T])
	t1 := tmp1.(types.Monotyped[T])
	// premise `t2 = newvar`
	t2 := cxt.typeContext.NewVar()
	// create monotype `t1 -> t2`
	t1_to_t2 := cxt.typeContext.Function(t1, t2)
	// premise `t0 = t1 -> t2`
	stat := cxt.Unify(t0, t1_to_t2)
	if stat.NotOk() {
		terms := []TypeJudgement[T]{j0, j1}
		report := makeReport("App", stat, terms...)
		cxt.appendReport(report) // TODO: signal failure how?
	}
	// "(e0 e1)" in result of rule
	appliedExpression := expr.Apply(e0, e1)
	v := types.Type[T](cxt.Find(t2)) // find t2's substitution
	// (e0 e1): t2
	return bridge.Judgement(appliedExpression, v)
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
//		                           [ f: forall a. a -> a ]¹    Inst(forall a. a -> a)
//		                           -------------------------------------------------- [Var]
//		                   0: Int                      f: x -> x                       t0, Int = x
//		                   ----------------------------------------------------------------------- [App]
//		                                               f 0: t0
//		                                               ------- [Id]
//		  (λx.x): a -> a                               f 0: Int
//		1 ----------------------------------------------------- [Let]
//		               let f = (λx.x) in f 0: Int
func (cxt *Context[T]) Let(name T, j0 bridge.JudgementAsExpression[T, expr.Expression[T]]) func(bridge.JudgementAsExpression[T, expr.Expression[T]]) bridge.JudgementAsExpression[T, expr.NameContext[T]] {
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

type nameableMonotype[T nameable.Nameable] interface {
	types.Monotyped[T]
	nameable.Nameable
}

type TypeJudgement[N nameable.Nameable] interface {
	GetExpressionAndType() (expr.Expression[N], types.Type[N])
}

// check if variable v occurs in monotype t. If it does, return true; else, return false.
// if t = v, then v is not in t, v is t
func (cxt *Context[T]) occurs(v types.Variable[T], t types.Monotyped[T]) (vOccursInT bool) {
	if IsVariable(t) {
		return false
	}

	for _, u := range t.GetFreeVariables() {
		if v.Equals(u) {
			return true
		}
	}

	return false
}

// declares, for types variable v, monotype t:
//
//	v = t
//
// if v ∈ t, then union returns `OccursCheckFailed`; else, skipUnify is returned
func (cxt *Context[T]) union(v types.Variable[T], t types.Monotyped[T]) Status {
	if cxt.occurs(v, t) {
		return OccursCheckFailed
	}

	cxt.typeSubs.Add(v, t)
	return skipUnify
}

func filter[N nameable.Nameable](e expr.Expression[N]) (out types.TypeJudgement[N, expr.Variable[N]], isVarJudge bool) {
	j, isJudge := e.(TypeJudgement[N])
	if !isJudge {
		return
	}
	ex, ty := j.GetExpressionAndType()
	v, isVar := ex.(expr.Variable[N])
	if !isVar {
		return
	}
	isVarJudge = isVar
	return types.Judgement[N, expr.Variable[N]](v, ty), isVarJudge
}

func filterJudgements[N nameable.Nameable](j types.ExpressionJudgement[N, expr.Referable[N]]) []types.TypeJudgement[N, expr.Variable[N]] {
	es := bridge.JudgementAsExpression[N, expr.Referable[N]](j.AsTypeJudgement()).Flatten()
	return fun.FMapFilter(es, filter[N])
}

func appendJudgements[T nameable.Nameable](left, right []types.TypeJudgement[T, expr.Variable[T]]) []types.TypeJudgement[T, expr.Variable[T]] {
	return append(left, right...)
}

func makeBaseJudgement[T nameable.Nameable]() []types.TypeJudgement[T, expr.Variable[T]] {
	return []types.TypeJudgement[T, expr.Variable[T]]{}
}

func flattenFilteredJudgements[T nameable.Nameable](filtered [][]types.TypeJudgement[T, expr.Variable[T]]) (flattened []types.TypeJudgement[T, expr.Variable[T]]) {
	return fun.FoldLeft(makeBaseJudgement[T](), filtered, appendJudgements[T])
}

func extractJudgements[T nameable.Nameable](dti types.DependentTypeInstance[T]) (binders []types.TypeJudgement[T, expr.Variable[T]]) {
	indexes := dti.GetIndexes()
	filtered := fun.FMap(indexes, filterJudgements[T])
	return flattenFilteredJudgements(filtered)
}

// generalizes a monotype into a dependent type
func DependentGeneralization[T nameable.Nameable](ty types.Monotyped[T]) types.DependentTyped[T] {
	if dti, ok := ty.(types.DependentTypeInstance[T]); ok {
		// get all kind vars in need of binding
		binders := extractJudgements(dti)
		ty = types.MakeDependentType(binders, dti.Application)
	}
	return ty
}

// generalizes a type: binds all free variables w/in monotype
func (cxt *Context[T]) Gen(t types.Monotyped[T]) types.Polytype[T] {
	// DependentGeneralization(`(t a0 .. aK; x0 .. xN)`) = `mapval (x0: X0) .. (xN: XN) . (t a0 .. aK)`
	dep := DependentGeneralization(t)
	// (t a0 .. aK; x0 .. xN) -> a0 .. aK
	vs := t.GetFreeVariables()
	// forall a0 .. aK . mapval (x0: X0) .. (xN: XN) . (t a0 .. aK)
	return types.Forall(vs...).Bind(dep)
}

func IsVariable[T nameable.Nameable](ty types.Monotyped[T]) bool {
	_, ok := ty.(types.Variable[T])
	return ok
}

func checkStatus[T nameable.Nameable](c0, c1 string, ms0, ms1 []types.Monotyped[T]) Status {
	if c0 != c1 {
		return ConstantMismatch
	}

	if len(ms0) != len(ms1) {
		return ParamLengthMismatch
	}

	return Ok
}

func Split[T nameable.Nameable](m types.Monotyped[T]) (c string, params []types.Monotyped[T]) {
	st, splittable := m.(types.Splitable[T])
	if !splittable {
		c, params = m.GetReferred().GetName(), nil
	} else {
		c, params = st.Split()
	}
	return
}

// returns Ok iff (stat.IsOk() || stat.Is(skipUnify))
func fixSkip(stat Status) Status {
	if stat.Is(skipUnify) {
		return Ok
	}
	return stat
}

func (cxt otherwiseDo[T]) otherwiseUnify(a, b types.Monotyped[T]) Status {
	// function pre-condition: substitution already happened or occurs-check
	// failed?
	if cxt.stat.NotOk() {
		return fixSkip(cxt.stat)
	}

	// get constants and params
	ca, paramsOfA := Split(a)
	cb, paramsOfB := Split(b)

	// check if alright to use in loop
	stat := checkStatus(ca, cb, paramsOfA, paramsOfB)

	// it. through all params while stat is ok, unifying params
	for i := 0; stat.IsOk() && i < len(paramsOfA); i++ {
		pa, pb := paramsOfA[i], paramsOfB[i]
		stat = cxt.Unify(pa, pb)
	}

	return stat
}

type otherwiseDo[T nameable.Nameable] struct {
	stat Status
	*Context[T]
}

func (cxt *Context[T]) substitute(ta, tb types.Monotyped[T]) otherwiseDo[T] {
	stat := Ok

	if v, ok := ta.(types.Variable[T]); ok {
		stat = cxt.union(v, tb)
	} else if v, ok := tb.(types.Variable[T]); ok {
		stat = cxt.union(v, ta)
	}

	return otherwiseDo[T]{stat, cxt}
}

func (cxt *Context[T]) Unify(a, b types.Monotyped[T]) Status {
	ta := cxt.Find(a)
	tb := cxt.Find(b)

	return cxt.substitute(ta, tb).otherwiseUnify(ta, tb)
}
