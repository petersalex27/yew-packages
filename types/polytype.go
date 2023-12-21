package types

import (
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

// binds zero or more variables in a dependent type. Written in its
// most general form, polytypes have the form
//
//	(forall t1 t2 ...) . (mapval (a1: A1) (a2: A2) ...) . T
//
// where; for i, j in Uint; ti is an arbitrary type variable; aj is an
// arbitrary kind variable; Aj is an arbitrary monotype; and T is an
// arbitrary monotype.
type Polytype[T nameable.Nameable] struct {
	typeBinders []Variable[T]
	bound       DependentTyped[T]
}

// returns the same slice of variables that `p` has access to; it is NOT safe
// to modify the slice returned
func (p Polytype[T]) GetBinders() []Variable[T] {
	return p.typeBinders
}

// returns type bound by polytype
func (p Polytype[T]) GetBound() DependentTyped[T] { return p.bound }

type binders[T nameable.Nameable] Polytype[T]

// See types.Forall[T](...Variable[T]) for description
// Deprecated: used for testing purposes and will be un-exported or removed
// soon
func (cxt *Context[T]) Forall(vs ...string) binders[T] {
	out := binders[T]{
		typeBinders: make([]Variable[T], len(vs)),
	}
	for i, v := range vs {
		out.typeBinders[i] = cxt.Var(v)
	}
	return out
}

// prepares type variables to act as binders in a polytype
func Forall[T nameable.Nameable](vs ...Variable[T]) binders[T] {
	return binders[T]{
		typeBinders: vs,
	}
}

func (p binders[T]) Bind(t DependentTyped[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: p.typeBinders,
		bound:       t,
	}
}

// Forall("x", "y").Bind(Apply("Type", "x")).String()
// 	== "forall x y . (Type x)"
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

func (p Polytype[T]) Collect() []T {
	res := make([]T, 0, len(p.typeBinders))
	for _, v := range p.typeBinders {
		res = append(res, v.Collect()...)
	}
	res = append(res, p.bound.Collect()...)
	return res
}

// test **syntactic** equality! I.e., two types are equal when
// they only contain symbols shared between the two types and those
// symbols appear in exactly the same order with exactly the same structure.
// for example:
//	 (forall x1 . x1) != (forall x2 x1 . x1)
// despite the two types being able to be used in similar ways.
// additionally:
//   (forall x1 x2 . x1) != (forall y1 y2 . y1)
// despite always being able to be used in the same way.
func (p Polytype[T]) Equals(t Type[T]) bool {
	q, ok := t.(Polytype[T])
	if !ok || len(p.typeBinders) != len(q.typeBinders) {
		return false
	}

	for i, binder := range p.typeBinders {
		if !binder.Equals(q.typeBinders[i]) {
			return false
		}
	}

	return p.bound.Equals(q.bound)
}

// replaces all variables bound directly by the polytype with new
// free variables, then returns the resulting dependent type
func (p Polytype[T]) Specialize(cxt *Context[T]) Monotyped[T] {
	var t DependentTyped[T] = p.bound
	// create new type variables
	vs := fun.FMap(
		p.typeBinders,
		func(v Variable[T]) Monotyped[T] {
			return cxt.NewVar()
		},
	)
	// replace all bound variables w/ newly created type variables
	return t.ReplaceDependent(p.typeBinders, vs)
}

// replaces the first variable bound by the polytype with a 
// the monotype `m`; after, if there are no more variables bound by the 
// polytype, the resulting dependent type is returned, else the 
// instantiated polytype is returned
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
