package types

import (
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type Polytype[T nameable.Nameable] struct {
	typeBinders []Variable[T]
	bound       DependentTyped[T]
}

type partialPoly[T nameable.Nameable] Polytype[T]

func (cxt *Context[T]) Forall(vs ...string) partialPoly[T] {
	out := partialPoly[T]{
		typeBinders: make([]Variable[T], len(vs)),
	}
	for i, v := range vs {
		out.typeBinders[i] = cxt.Var(v)
	}
	return out
}

func (p partialPoly[T]) Bind(t DependentTyped[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: p.typeBinders,
		bound:       t,
	}
}

// (Polytype[T]{[]Variable{{0,"x"},{0,"y"}},nil,Application{Constant("Type[T]"),[]Monotyped[T]{Variable{0,"x"}}}}).String()
// == "forall x y . (Type[T] x)"
func (p Polytype[T]) String() string {
	if len(p.typeBinders) == 0 {
		return p.bound.String()
	}

	j := str.String(" ")
	return "forall " +
		str.Join(p.typeBinders, j) +
		" . " +
		p.bound.String()
}

func (p Polytype[T]) Generalize(*Context[T]) Polytype[T] {
	var out Polytype[T]
	out.bound = p.bound
	out.typeBinders = make([]Variable[T], len(p.typeBinders)+1)
	copy(out.typeBinders, p.typeBinders)
	return out
}

func (p Polytype[T]) Equals(t Type[T]) bool {
	q, ok := t.(Polytype[T])
	if !ok {
		return false
	}
	return p.bound.Equals(q.bound)
}

// func (p Polytype[T]) Specialize() Type[T]

func (p Polytype[T]) freeInstantiate(cxt *Context[T]) DependentTyped[T] {
	var t DependentTyped[T] = p.bound
	dummyVar := cxt.dummyName(Variable[T]{})
	for _, v := range p.typeBinders {
		t = MaybeReplace[T](t, v, dummyVar).(DependentTyped[T])
	}
	return t
}

func (p Polytype[T]) Instantiate(m Monotyped[T]) Type[T] {
	var t DependentTyped[T] = p.bound

	binderLength := len(p.typeBinders)
	if binderLength == 0 {
		return t
	}

	if p.typeBinders[0].name.GetName() != "_" { // if not non-binding binder
		t = MaybeReplace[T](t, p.typeBinders[0], m).(DependentTyped[T])
	}

	if binderLength == 1 {
		return t
	}

	binders := make([]Variable[T], binderLength-1)
	copy(binders, p.typeBinders[1:])
	return Polytype[T]{
		typeBinders: binders,
		bound:       t,
	}
}
