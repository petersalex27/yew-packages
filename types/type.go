package types

import (
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type collectable[T nameable.Nameable] interface {
	Collect() []T
}

type Type[T nameable.Nameable] interface {
	str.Stringable
	Equals(Type[T]) bool
	Generalize(cxt *Context[T]) Polytype[T]
	collectable[T]
}

func MaybeReplace[T nameable.Nameable](ty Type[T], v Variable[T], m Monotyped[T]) Type[T] {
	if mono, ok := ty.(Monotyped[T]); ok {
		return mono.Replace(v, m)
	}
	return ty
}