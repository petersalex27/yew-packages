package types

import expr "github.com/petersalex27/yew-packages/expr"


type DependentTypeConstructor struct {
	Monotyped
	index DependentTypeInstance
}

func (c DependentTypeConstructor) FreeInstantiateKinds(vs ...TypeJudgement[expr.Variable]) DependentTypeConstructor {
	return DependentTypeConstructor{
		Monotyped: c.Monotyped,
		index: c.index.FreeInstantiateKinds(vs...),
	}
}

func (c DependentTypeConstructor) Replace(v Variable, m Monotyped) DependentTypeConstructor {
	return DependentTypeConstructor{
		Monotyped: c.Monotyped.Replace(v, m),
		index: c.index.Replace(v, m).(DependentTypeInstance),
	}
}