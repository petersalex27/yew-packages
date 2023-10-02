package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
)


type DependentTypeConstructor[T nameable.Nameable] struct {
	Monotyped[T]
	index DependentTypeInstance[T]
}

func (c DependentTypeConstructor[T]) FreeInstantiateKinds(vs ...TypeJudgement[T,expr.Variable]) DependentTypeConstructor[T] {
	return DependentTypeConstructor[T]{
		Monotyped: c.Monotyped,
		index: c.index.FreeInstantiateKinds(vs...),
	}
}

func (c DependentTypeConstructor[T]) Replace(v Variable[T], m Monotyped[T]) DependentTypeConstructor[T] {
	return DependentTypeConstructor[T]{
		Monotyped: c.Monotyped.Replace(v, m),
		index: c.index.Replace(v, m).(DependentTypeInstance[T]),
	}
}

func (c DependentTypeConstructor[T]) Collect() []T {
	res := c.Monotyped.Collect()
	res = append(res, c.index.Collect()...)
	return res
}