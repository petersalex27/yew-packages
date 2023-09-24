package types

import (
	str "github.com/petersalex27/yew-packages/stringable"
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