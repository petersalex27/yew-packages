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
	//"fmt"

	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

type TypeJudgement[N nameable.Nameable] interface {
	GetExpressionAndType() (expr.Expression[N], types.Type[N])
}

func JudgementsEqual[N nameable.Nameable](a, b TypeJudgement[N]) bool {
	ea, ta := a.GetExpressionAndType()
	eb, tb := b.GetExpressionAndType()
	return ea.StrictEquals(eb) && ta.Equals(tb)
}

// like (TypeJudgement) GetExpressionAndType() (expr.Expression[N], types.Type[N]), 
// but does type assertions to return desired E and T types.
func GetExpressionAndType[N nameable.Nameable, E expr.Expression[N], T types.Type[N]](judgement TypeJudgement[N]) (e E, t T) {
	someExpression, someType := judgement.GetExpressionAndType()
	e, _ = someExpression.(E)
	t, _ = someType.(T)
	return
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

// check if variable v occurs in expression e. If it does, return true; else, return false.
// if e = v, then v is not in e: v is e
func (cxt *Context[T]) kindOccurs(v expr.Variable[T], e expr.Referable[T]) (vOccursInE bool) {
	if types.IsKindVariable(expr.Expression[T](e)) {
		return false
	}

	for _, u := range e.ExtractVariables(0) {
		if v.StrictEquals(u) {
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

	// if d, ok := t.(types.DependentTypeInstance[T]); ok {
	// 	t = cxt.reindex(d)
	// }

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
		return types.MakeDependentType(binders, types.TypeFunction[T](dti.Application))
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

func checkStatus[T nameable.Nameable](c0, c1 string, ms0, ms1 []types.Monotyped[T], isA, isB types.Indexes[T]) Status {
	if c0 != c1 {
		return ConstantMismatch
	}

	if len(ms0) != len(ms1) {
		return ParamLengthMismatch
	}

	if len(isA) != len(isB) {
		return IndexLengthMismatch
	}

	return Ok
}

func Split[T nameable.Nameable](m types.Monotyped[T]) (c string, params []types.Monotyped[T], indexes types.Indexes[T]) {
	st, isTypeFunc := m.(types.TypeFunction[T])
	if !isTypeFunc {
		c, params, indexes = m.GetReferred().GetName(), nil, nil
	} else {
		var app types.Application[T]
		app, indexes = st.FunctionAndIndexes()
		c, params = app.Split()
	}
	return
}

func SplitKind[T nameable.Nameable](kind expr.Referable[T]) (c string, mems []bridge.JudgementAsExpression[T, expr.Expression[T]]) {
	if data, ok := kind.(bridge.Data[T]); ok {
		mems = data.Members
	}

	c = kind.GetReferred().GetName()
	return
}

func (cxt *Context[T]) kindUnion(v expr.Variable[T], e expr.Referable[T]) Status {
	if cxt.kindOccurs(v, e) {
		return OccursCheckFailed
	}

	cxt.exprSubs.Add(v.GetReferred(), e)
	return skipUnify
}

func checkKindStatus[T nameable.Nameable](ca, cb string, memsOfA, memsOfB []bridge.JudgementAsExpression[T, expr.Expression[T]]) Status {
	if ca != cb {
		return KindConstantMismatch
	}

	if len(memsOfA) != len(memsOfB) {
		return MemsLengthMismatch
	}

	return Ok
}

func (cxt otherwiseDo[T]) otherwiseUnifyKind(a, b expr.Referable[T]) Status {
	// function pre-condition: substitution already happened or occurs-check
	// failed?
	if cxt.stat.NotOk() {
		return fixSkip(cxt.stat)
	}

	// get constants and mems
	ca, memsOfA := SplitKind(a)
	cb, memsOfB := SplitKind(b)

	// check if alright to use in loop
	stat := checkKindStatus(ca, cb, memsOfA, memsOfB)

	// it. through all mems while stat is ok, unifying mems
	for i := 0; stat.IsOk() && i < len(memsOfA); i++ {
		ma, _ := memsOfA[i].GetExpressionAndType()
		mb, _ := memsOfB[i].GetExpressionAndType()
		stat = cxt.UnifyKind(ma.(expr.Referable[T]), mb.(expr.Referable[T]))
	}

	return stat
}

// tries to creates a substitution from a variable to a monotype
func (cxt *Context[T]) substituteKind(ea, eb expr.Referable[T]) otherwiseDo[T] {
	stat := Ok

	if v, ok := ea.(expr.Variable[T]); ok {
		stat = cxt.kindUnion(v, eb)
	} else if v, ok := eb.(expr.Variable[T]); ok {
		stat = cxt.kindUnion(v, ea)
	}

	return otherwiseDo[T]{stat, cxt}
}

func (cxt *Context[T]) UnifyKind(a, b expr.Referable[T]) Status {
	ea := cxt.FindKind(a)
	eb := cxt.FindKind(b)

	return cxt.substituteKind(ea, eb).otherwiseUnifyKind(ea, eb)
}

func (cxt *Context[T]) UnifyIndex(indexOfA, indexOfB types.ExpressionJudgement[T, expr.Referable[T]]) Status {
	// split type and expression
	ea, ta := indexOfA.AsTypeJudgement().GetExpressionAndType()
	eb, tb := indexOfB.AsTypeJudgement().GetExpressionAndType()

	// type assertions monotype
	ma := ta.(types.Monotyped[T])
	mb := tb.(types.Monotyped[T])
	// type assertions referable
	ra := ea.(expr.Referable[T])
	rb := eb.(expr.Referable[T])

	// unify types
	stat := cxt.Unify(ma, mb)

	// if type union Ok, then unify kinds
	if stat.IsOk() {
		stat = cxt.UnifyKind(ra, rb)
	}

	return stat
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

	// get constants, params, and indexes
	ca, paramsOfA, indexesOfA := Split(a)
	cb, paramsOfB, indexesOfB := Split(b)

	// check if alright to use in loop
	stat := checkStatus(ca, cb, paramsOfA, paramsOfB, indexesOfA, indexesOfB)

	// it. through all params while stat is ok, unifying params
	for i := 0; stat.IsOk() && i < len(paramsOfA); i++ {
		pa, pb := paramsOfA[i], paramsOfB[i]
		stat = cxt.Unify(pa, pb)
	}

	// it. through all indexes while stat is ok, unifying indexes
	for i := 0; stat.IsOk() && i < len(indexesOfA); i++ {
		ia, ib := indexesOfA[i], indexesOfB[i]
		stat = cxt.UnifyIndex(ia, ib)
	}

	return stat
}

type otherwiseDo[T nameable.Nameable] struct {
	stat Status
	*Context[T]
}

// tries to creates a substitution from a variable to a monotype
func (cxt *Context[T]) substitute(ta, tb types.Monotyped[T]) otherwiseDo[T] {
	stat := Ok

	if v, ok := ta.(types.Variable[T]); ok {
		stat = cxt.union(v, tb)
	} else if v, ok := tb.(types.Variable[T]); ok {
		stat = cxt.union(v, ta)
	}

	return otherwiseDo[T]{stat, cxt}
}

// unifies two monotypes a, b
func (cxt *Context[T]) Unify(a, b types.Monotyped[T]) Status {
	ta := cxt.Find(a)
	tb := cxt.Find(b)

	return cxt.substitute(ta, tb).otherwiseUnify(ta, tb)
}
