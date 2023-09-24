package types

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	str "github.com/petersalex27/yew-packages/stringable"
)

// picks out a monotype
type DependentTypeInstance struct {
	Application
	index []TypeJudgement[expr.Expression]
}

func Index(family Application, domain ...TypeJudgement[expr.Expression]) DependentTypeInstance {
	return DependentTypeInstance{
		Application: family,
		index: domain,
	}
}

func (dti DependentTypeInstance) FreeInstantiateKinds(vs ...TypeJudgement[expr.Variable]) DependentTypeInstance {
	return DependentTypeInstance{
		Application: dti.Application,
		index: fun.FMap(dti.index, func(i TypeJudgement[expr.Expression]) TypeJudgement[expr.Expression] {
			for _, v := range vs {
				expr, _ := i.expression.Replace(v.expression, expr.Var("_"))
				i.expression = expr
			}
			return i
		}),
	}
}

func (dti DependentTypeInstance) String() string {
	return "(" + dti.Application.String() + "; " + str.Join(dti.index, str.String(" ")) + ")"
}

func (dti DependentTypeInstance) Equals(t Type) bool {
	dti2, ok := t.(DependentTypeInstance)
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

func (dti DependentTypeInstance) Replace(v Variable, m Monotyped) Monotyped {
	return DependentTypeInstance{
		Application: dti.Application.Replace(v, m).(Application),
		index: dti.index,
	}
}

func (dti DependentTypeInstance) ReplaceDependent(v Variable, m Monotyped) DependentTyped {
	return dti.Replace(v, m)
}

func (dti DependentTypeInstance) FreeInstantiation() DependentTyped {
	return DependentTypeInstance{
		Application: dti.Application.FreeInstantiation().(Application),
		index: dti.index,
	}
}