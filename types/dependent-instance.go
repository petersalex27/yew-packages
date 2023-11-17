package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

// picks out a monotype
type DependentTypeInstance[T nameable.Nameable] struct {
	Application Application[T] // dependent type function
	Indexes     Indexes[T]     // arguments to function
}

func (dti DependentTypeInstance[N]) Rebuild(findMono func(Monotyped[N]) Monotyped[N], findKind func(expr.Referable[N]) expr.Referable[N]) TypeFunction[N] {
	return DependentTypeInstance[N]{
		dti.Application.Rebuild(findMono, findKind).(Application[N]),
		fun.FMap(
			dti.Indexes,
			func(ej ExpressionJudgement[N, expr.Referable[N]]) ExpressionJudgement[N, expr.Referable[N]] {
				eTmp, t := ej.AsTypeJudgement().GetExpressionAndType()
				e := eTmp.(expr.Referable[N])
				return Judgement(findKind(e), t)
			},
		),
	}
}

func (dti DependentTypeInstance[T]) FunctionAndIndexes() (function Application[T], indexes Indexes[T]) {
	return dti.Application, dti.Indexes
}

func (dti DependentTypeInstance[T]) GetIndexes() []ExpressionJudgement[T, expr.Referable[T]] {
	return dti.Indexes
}

func (dti DependentTypeInstance[T]) GetFreeKindVariables() []expr.Variable[T] {
	vars := []expr.Variable[T]{}
	for _, index := range dti.Indexes {
		vars = append(vars, index.AsTypeJudgement().expression.ExtractVariables(0)...)
	}
	return vars
}

func (dti DependentTypeInstance[T]) GetFreeVariables() []Variable[T] {
	vars := dti.Application.GetFreeVariables()
	vars = append(vars, dti.Indexes.GetFreeVariables()...)
	return vars
}

func (dti DependentTypeInstance[T]) GetReferred() T {
	return dti.Application.GetReferred()
}

func Index[T nameable.Nameable](family Application[T], domain ...ExpressionJudgement[T, expr.Referable[T]]) DependentTypeInstance[T] {
	return DependentTypeInstance[T]{
		Application: family,
		Indexes:     domain,
	}
}

func (dti DependentTypeInstance[T]) AsFreeInstance(vs []TypeJudgement[T, expr.Variable[T]], replacements []expr.Referable[T]) TypeFunction[T] {
	application := dti.Application
	if len(dti.Indexes) == 0 {
		return DependentTypeInstance[T]{
			application, 
			fun.ZipWith(
				func(ref expr.Referable[T], old TypeJudgement[T, expr.Variable[T]]) ExpressionJudgement[T, expr.Referable[T]] {
					return Judgement(ref, old.ty)
				},
				replacements,
				vs,
			),
		}
	}
	if len(vs) != len(dti.Indexes) {
		panic("dependent type cannot be indexed")
	}

	
	indexed := fun.ZipWith(
		func(index ExpressionJudgement[T, expr.Referable[T]], v struct{Left TypeJudgement[T, expr.Variable[T]]; Right expr.Referable[T]}) ExpressionJudgement[T, expr.Referable[T]] {
			e := index.AsTypeJudgement().expression
			e2, _ := e.Replace(v.Left.expression, v.Right)
			r := e2.(expr.Referable[T])
			return Judgement(r, index.AsTypeJudgement().ty)
		},
		dti.Indexes,
		fun.Zip(vs, replacements),
	)
	return DependentTypeInstance[T]{application, indexed}
}

// uses new kind variables of corr. type in `vs` as arguments to dependent type function
func (dti DependentTypeInstance[T]) SubVars(preSub []TypeJudgement[T, expr.Variable[T]], postSub []expr.Referable[T]) TypeFunction[T] {
	return dti.AsFreeInstance(preSub, postSub)
}

func (dti DependentTypeInstance[T]) String() string {
	if len(dti.Indexes) == 0 {
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
	return lclose + app + dti.Indexes.String() + rclose
}

func (dti DependentTypeInstance[T]) Equals(t Type[T]) bool {
	dti2, ok := t.(DependentTypeInstance[T])
	ok = ok && dti.Application.Equals(dti2.Application) // check type assertion and application
	if !ok {
		return false
	}

	// check len of indexes
	if len(dti2.Indexes) != len(dti.Indexes) {
		return false
	}

	// check content of indexes
	for i, ind := range dti.Indexes {
		if !JudgesEquals(ind, dti2.Indexes[i]) {
			return false
		}
	}

	return true // dti == t
}

func (dti DependentTypeInstance[T]) Replace(v Variable[T], m Monotyped[T]) Monotyped[T] {
	return DependentTypeInstance[T]{
		Application: dti.Application.Replace(v, m).(Application[T]),
		Indexes:     dti.Indexes,
	}
}

func (dti DependentTypeInstance[T]) ReplaceDependent(vs []Variable[T], ms []Monotyped[T]) Monotyped[T] {
	return DependentTypeInstance[T]{
		Application: dti.Application.ReplaceDependent(vs, ms).(Application[T]),
		Indexes:     dti.Indexes,
	}
}

func (dti DependentTypeInstance[T]) Collect() []T {
	res := dti.Application.Collect()
	for _, v := range dti.Indexes {
		res = append(res, v.Collect()...)
	}
	return res
}
