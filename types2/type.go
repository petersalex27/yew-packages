package types

import (
	"alex.peters/yew/str"
)

type Type interface {
	str.Stringable
	Equals(Type) bool
	Generalize() Polytype
}

func MaybeReplace(ty Type, v Variable, m Monotyped) Type {
	if mono, ok := ty.(Monotyped); ok {
		return mono.Replace(v, m)
	}
	return ty
}