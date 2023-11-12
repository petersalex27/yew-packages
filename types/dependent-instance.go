package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type indexes[T nameable.Nameable] []ExpressionJudgement[T, expr.Expression[T]]

func (idxs indexes[T]) GetFreeVariables() []Variable[T] {
	// 2d slice of free variables
	vs2d := fun.FMap(
		idxs,
		func(ej ExpressionJudgement[T, expr.Expression[T]]) []Variable[T] {
			mono := ej.AsTypeJudgement().ty.(Monotyped[T]) // should always pass b/c indexes's values must be typed by monotypes
			return mono.GetFreeVariables()
		},
	)
	// convert to 1d slice of free variables
	return fun.FoldLeft(
		[]Variable[T]{},
		vs2d,
		func(vs, us []Variable[T]) []Variable[T] {
			return append(vs, us...)
		},
	)
}

func (idxs indexes[T]) FreeInstantiation(cxt *Context[T]) indexes[T] {
	return fun.FMap(
		idxs,
		func(idx ExpressionJudgement[T, expr.Expression[T]]) ExpressionJudgement[T, expr.Expression[T]] {
			tj := idx.AsTypeJudgement()
			t, e := tj.GetType(), tj.GetExpression()
			m := t.(Monotyped[T]).FreeInstantiation(cxt)
			return Judgement[T, expr.Expression[T]](e, m)
		},
	)
}

func (idxs indexes[T]) String() string {
	if len(idxs) == 0 {
		return ""
	}
	return "; " + str.Join(idxs, str.String(" "))
}

// picks out a monotype
type DependentTypeInstance[T nameable.Nameable] struct {
	Application[T]            // dependent type function
	index          indexes[T] // arguments to function
}
/*
func (dti DependentTypeInstance[T]) GetFreeKindVariables() []expr.Variable[T] {
	vars := []expr.Variable{}
	for _, index := range dti.index {

		index.AsTypeJudgement().expression.ExtractFreeVariables()
	}
}*/

func (dti DependentTypeInstance[T]) GetFreeVariables() []Variable[T] {
	vars := dti.Application.c.GetFreeVariables()
	vars = append(vars, dti.index.GetFreeVariables()...)
	return vars
}

func (dti DependentTypeInstance[T]) GetName() T {
	return dti.Application.GetName()
}

func Index[T nameable.Nameable](family Application[T], domain ...ExpressionJudgement[T, expr.Expression[T]]) DependentTypeInstance[T] {
	return DependentTypeInstance[T]{
		Application: family,
		index:       domain,
	}
}

func (dti DependentTypeInstance[T]) AsFreeInstance(freeExprVars []expr.Variable[T], vs []TypeJudgement[T, expr.Variable[T]]) DependentTypeInstance[T] {
	application := dti.Application
	indexed := fun.ZipWith(
		// `binder` is a binder in `mapall (binder1) (binder2) .. (binderN) . dti`
		func(freeExprVar expr.Variable[T], binder TypeJudgement[T, expr.Variable[T]]) ExpressionJudgement[T, expr.Expression[T]] {
			t := binder.GetType()                  // type of index value
			var e expr.Expression[T] = freeExprVar // free variable that indexes dependent type (equiv'ly, arg. to dep. ty. func.)
			return Judgement(e, t)                 // judgement of `e: t`
		},
		freeExprVars,
		vs,
	)
	return DependentTypeInstance[T]{application, indexed}
}

// uses new kind variables of corr. type in `vs` as arguments to dependent type function
func (dti DependentTypeInstance[T]) FreeInstantiateKinds(cxt *expr.Context[T], vs ...TypeJudgement[T, expr.Variable[T]]) DependentTypeInstance[T] {
	freeExprVars := make([]expr.Variable[T], len(vs))
	for i := range freeExprVars {
		freeExprVars[i] = cxt.NewVar()
	}
	return dti.AsFreeInstance(freeExprVars, vs)
}

func (dti DependentTypeInstance[T]) String() string {
	if len(dti.index) == 0 {
		return dti.Application.String()
	}
	lclose, rclose := "(", ")"
	app := ""
	if ec, ok := dti.Application.c.(EnclosingConst[T]); ok {
		lclose, rclose = ec.SplitString()
		app = str.Join(dti.Application.ts, str.String(" "))
	} else {
		app = dti.Application.String()
	}
	return lclose + app + dti.index.String() + rclose
}

func (dti DependentTypeInstance[T]) Equals(t Type[T]) bool {
	dti2, ok := t.(DependentTypeInstance[T])
	ok = ok && dti.Application.Equals(dti2.Application) // check type assertion and application
	if !ok {
		return false
	}

	// check len of indexes
	if len(dti2.index) != len(dti.index) {
		return false
	}

	// check content of indexes
	for i, ind := range dti.index {
		if !JudgesEquals(ind, dti2.index[i]) {
			return false
		}
	}

	return true // dti == t
}

func (dti DependentTypeInstance[T]) Replace(v Variable[T], m Monotyped[T]) Monotyped[T] {
	return DependentTypeInstance[T]{
		Application: dti.Application.Replace(v, m).(Application[T]),
		index:       dti.index,
	}
}

func (dti DependentTypeInstance[T]) ReplaceDependent(vs []Variable[T], ms []Monotyped[T]) Monotyped[T] {
	return DependentTypeInstance[T]{
		Application: dti.Application.ReplaceDependent(vs, ms).(Application[T]),
		index:       dti.index,
	}
}

func (dti DependentTypeInstance[T]) FreeInstantiation(cxt *Context[T]) Monotyped[T] {
	return DependentTypeInstance[T]{
		Application: dti.Application.FreeInstantiation(cxt).(Application[T]),
		index:       dti.index.FreeInstantiation(cxt),
	}
}

func (dti DependentTypeInstance[T]) Collect() []T {
	res := dti.Application.Collect()
	for _, v := range dti.index {
		res = append(res, v.Collect()...)
	}
	return res
}
