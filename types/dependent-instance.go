package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type indexes[T nameable.Nameable] []TypeJudgement[T,expr.Expression[T]]

func (idxs indexes[T]) String() string {
	if len(idxs) == 0 {
		return ""
	}
	return "; " + str.Join(idxs, str.String(" "))
}

// picks out a monotype
type DependentTypeInstance[T nameable.Nameable] struct {
	Application[T]		// dependent type function
	index indexes[T]	// arguments to function
}

func (dti DependentTypeInstance[T]) GetName() T {
	return dti.Application.GetName()
}

func Index[T nameable.Nameable](family Application[T], domain ...TypeJudgement[T,expr.Expression[T]]) DependentTypeInstance[T] {
	return DependentTypeInstance[T]{
		Application: family,
		index: domain,
	}
}

func (dti DependentTypeInstance[T]) FreeInstantiateKinds(cxt *Context[T], vs ...TypeJudgement[T,expr.Variable[T]]) DependentTypeInstance[T] {
	return DependentTypeInstance[T]{
		Application: dti.Application,
		index: fun.FMap(dti.index, func(i TypeJudgement[T,expr.Expression[T]]) TypeJudgement[T,expr.Expression[T]] {
			for _, v := range vs {
				expr, _ := i.expression.Replace(v.expression, expr.Var(cxt.makeName("_")))
				i.expression = expr
			}
			return i
		}),
	}
}

func (dti DependentTypeInstance[T]) String() string {
	if len(dti.index) == 0 {
		return dti.Application.String()
	}
	return "(" + dti.Application.String() + dti.index.String() + ")"
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
		if !Equals(ind, dti2.index[i]) {
			return false
		}
	} 

	return true // dti == t
}

func (dti DependentTypeInstance[T]) Replace(v Variable[T], m Monotyped[T]) Monotyped[T] {
	return DependentTypeInstance[T]{
		Application: dti.Application.Replace(v, m).(Application[T]),
		index: dti.index,
	}
}

func (dti DependentTypeInstance[T]) ReplaceDependent(v Variable[T], m Monotyped[T]) DependentTyped[T] {
	return dti.Replace(v, m)
}

func (dti DependentTypeInstance[T]) FreeInstantiation(cxt *Context[T]) DependentTyped[T] {
	return DependentTypeInstance[T]{
		Application: dti.Application.FreeInstantiation(cxt).(Application[T]),
		index: dti.index,
	}
}

func (dti DependentTypeInstance[T]) Collect() []T {
	res := dti.Application.Collect()
	for _, v := range dti.index {
		res = append(res, v.Collect()...)
	}
	return res
}