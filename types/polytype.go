package types

import (
	str "github.com/petersalex27/yew-packages/stringable"
)

type Polytype struct {
	typeBinders []Variable
	bound DependentTyped
}

type partialPoly Polytype

func Forall(vs ...string) partialPoly {
	out := partialPoly{
		typeBinders: make([]Variable, len(vs)),
	}
	for i, v := range vs {
		out.typeBinders[i] = Var(v)
	}
	return out
}

func (p partialPoly) Bind(t DependentTyped) Polytype {
	return Polytype{
		typeBinders: p.typeBinders,
		bound: t,
	}
}


// (Polytype{[]Variable{{0,"x"},{0,"y"}},nil,Application{Constant("Type"),[]Monotyped{Variable{0,"x"}}}}).String()
// == "forall x y . (Type x)"
func (p Polytype) String() string {
	if len(p.typeBinders) == 0 {
		return p.bound.String()
	}

	j := str.String(" ")
	return "forall " +
		str.Join(p.typeBinders, j) +
		" . " +
		p.bound.String()
}

func (p Polytype) Generalize() Polytype {
	var out Polytype
	out.bound = p.bound
	out.typeBinders = make([]Variable, len(p.typeBinders)+1)
	copy(out.typeBinders, p.typeBinders)
	return out
}

func (p Polytype) Equals(t Type) bool {
	q, ok := t.(Polytype)
	if !ok {
		return false
	}
	return p.freeInstantiate().Equals(q.freeInstantiate())
}

// func (p Polytype) Specialize() Type

func (p Polytype) freeInstantiate() DependentTyped {
	var t DependentTyped = p.bound
	dummyVar := NonBindableVar("_")
	for _, v := range p.typeBinders {
		t = MaybeReplace(t, v, dummyVar).(DependentTyped)
	}
	return t
}

func (p Polytype) Instantiate(m Monotyped) Type {
	var t DependentTyped = p.bound
	
	binderLength := len(p.typeBinders)
	if binderLength == 0 {
		return t
	}

	if p.typeBinders[0].name != "_" { // if not non-binding binder
		t = MaybeReplace(t, p.typeBinders[0], m).(DependentTyped)
	}

	if binderLength == 1 {
		return t
	}

	binders := make([]Variable, binderLength-1)
	copy(binders, p.typeBinders[1:])
	return Polytype{
		typeBinders: binders,
		bound: t,
	}
}