package types

import "github.com/petersalex27/yew-packages/nameable"

type Monotyped[T nameable.Nameable] interface {
	Type[T]
	DependentTyped[T]
	GetReferred() T
	ReplaceKindVar(replacing Variable[T], with Monotyped[T]) Monotyped[T]
	Replace(Variable[T], Monotyped[T]) Monotyped[T]
	GetFreeVariables() []Variable[T]
}

type Splitable[T nameable.Nameable] interface {
	Monotyped[T]
	Split() (name string, params []Monotyped[T])
}

func Name[T nameable.Nameable](m Monotyped[T]) (name string, hasName bool) {
	var st Splitable[T]
	st, hasName = m.(Splitable[T])
	if !hasName {
		name = m.String()
	} else {
		name, _ = st.Split()
	}
	return
}

func IsVariable[T nameable.Nameable](m Monotyped[T]) bool {
	_, isVar := m.(Variable[T])
	return isVar
}

func Split[T nameable.Nameable](m Monotyped[T]) (name string, params []Monotyped[T]) {
	st, hasName := m.(Splitable[T])
	if !hasName {
		name, params = m.String(), nil
	} else {
		name, params = st.Split()
	}
	return
}