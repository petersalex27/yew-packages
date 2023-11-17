package types

import (
	"github.com/petersalex27/yew-packages/nameable"
)

type Monotyped[T nameable.Nameable] interface {
	Type[T]
	DependentTyped[T]
	GetReferred() T
	Replace(Variable[T], Monotyped[T]) Monotyped[T]
	GetFreeVariables() []Variable[T]
}