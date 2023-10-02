package types

import (
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type Type[T nameable.Nameable] interface {
	str.Stringable
	Equals(Type[T]) bool
	Generalize(cxt *Context[T]) Polytype[T]
}

func MaybeReplace[T nameable.Nameable](ty Type[T], v Variable[T], m Monotyped[T]) Type[T] {
	if mono, ok := ty.(Monotyped[T]); ok {
		return mono.Replace(v, m)
	}
	return ty
}