package types

import (
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type Type[T nameable.Nameable] interface {
	str.Stringable
	Equals(Type[T]) bool
	nameable.Collectable[T]
}

func MaybeReplace[T nameable.Nameable](ty Type[T], v Variable[T], m Monotyped[T]) Type[T] {
	if mono, ok := ty.(Monotyped[T]); ok {
		return mono.Replace(v, m)
	}
	return ty
}

func Merge[T nameable.Nameable](head Application[T], tail Monotyped[T]) Application[T] {
	newTypes := make([]Monotyped[T], 0, len(head.ts)+1)
	copy(newTypes, head.ts)
	return Application[T]{
		c:  head.c,
		ts: append(newTypes, tail),
	}
}
