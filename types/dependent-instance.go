package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

// picks out a monotype
type DependentTypeInstance[T nameable.Nameable] struct {
	Application[T]
	index []TypeJudgement[T,expr.Expression]
}

func Index[T nameable.Nameable](family Application[T], domain ...TypeJudgement[T,expr.Expression]) DependentTypeInstance[T] {
	return DependentTypeInstance[T]{
		Application: family,
		index: domain,
	}
}

func (dti DependentTypeInstance[T]) FreeInstantiateKinds(vs ...TypeJudgement[T,expr.Variable]) DependentTypeInstance[T] {
	return DependentTypeInstance[T]{
		Application: dti.Application,
		index: fun.FMap(dti.index, func(i TypeJudgement[T,expr.Expression]) TypeJudgement[T,expr.Expression] {
			for _, v := range vs {
				expr, _ := i.expression.Replace(v.expression, expr.Var("_"))
				i.expression = expr
			}
			return i
		}),
	}
}

func (dti DependentTypeInstance[T]) String() string {
	return "(" + dti.Application.String() + "; " + str.Join(dti.index, str.String(" ")) + ")"
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