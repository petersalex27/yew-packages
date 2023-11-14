package expr

import (
	"strconv"

	"github.com/petersalex27/yew-packages/nameable"
)

type Variable[T nameable.Nameable] struct {
	name  T
	depth int
}

func (v Variable[T]) Flatten() []Expression[T] {
	return []Expression[T]{v}
}

func (v Variable[T]) GetReferred() T {
	return v.name
}

func (v Variable[T]) ToAlmostPattern() (AlmostPattern[T], bool) {
	return MakeElem[T](PatternElementVar, v.name).ToAlmostPattern()
}

func (v Variable[T]) BodyAbstract(Variable[T], Const[T]) Expression[T] { return v }

func (v Variable[T]) ExtractVariables(gt int) []Variable[T] {
	if v.depth > gt {
		// variable is free variable; all bound variables were replaced w/ dummy
		// variable
		return []Variable[T]{v}
	}
	return []Variable[T]{}
}

func (v Variable[T]) Collect() []T {
	return []T{v.name}
}

func (v Variable[T]) copy() Variable[T] {
	return Variable[T]{
		name:  v.name,
		depth: v.depth,
	}
}

func (v Variable[T]) Copy() Expression[T] {
	return v.copy()
}

// Makes a variable
//
// Deprecated: to be removed???
func MakeVar[T nameable.Nameable](name T, depth int) Variable[T] {
	return Variable[T]{name: name, depth: depth}
}

func (cxt *Context[T]) makeVar(name string, depth int) Variable[T] {
	return Variable[T]{name: cxt.makeName(name), depth: depth}
}

func (v Variable[T]) PrepareAsRHS() Expression[T] {
	if v.depth < 1 {
		return Variable[T]{
			name:  v.name,
			depth: 1,
		}
	}
	return v
}

func (v Variable[T]) UpdateVars(gt int, by int) Expression[T] {
	if v.depth > gt {
		newVar := Var(v.name)
		newVar.depth = v.depth + by
		return newVar
	}
	return v
}

func (v Variable[T]) Rebind() Expression[T] {
	return Var(v.name)
}

func (v Variable[T]) Bind(bs BindersOnly[T]) Expression[T] {
	depth := len(bs)
	if v.depth != 0 && v.depth <= depth {
		return v
	}

	name := v.name
	out := Var(name)
	// is free Variable
	for _, b := range bs {
		if name.GetName() == b.name.GetName() {
			// Variable gets bound at b.depth
			out.depth = b.depth
			return out
		}
		// Variable does not get bound, maybe next binder..?
	}

	// Variable remains unbound

	// Set variable as a free variable; free var# + depth of binders;
	// this represents a free variable w/ number value "#" that is free w/in 
	// the "depth" enclosing binders
	out.depth = v.depth + depth
	if v.depth == 0 { // free variables should not have value 0
		// set variable number as 1
		out.depth = out.depth + 1 // look at that! +1! Variable is recognized :)
	}
	return out
}

func (cxt *Context[T]) Var(name string) Variable[T] {
	return Var[T](cxt.makeName(name))
}

func Var[T nameable.Nameable](name T) Variable[T] {
	return Variable[T]{name: name, depth: 0}
}

func (v Variable[T]) Again() (Expression[T], bool) {
	return v, false
}

func (v Variable[T]) Replace(w Variable[T], e Expression[T]) (Expression[T], bool) {
	if varEquals(v, w) {
		return e, false
	}
	return v, false
}

func (v Variable[T]) Find(w Variable[T]) bool { return varEquals(v, w) }

func varEquals[T nameable.Nameable](v, w Variable[T]) bool {
	return v.depth == w.depth && v.name.GetName() == w.name.GetName()
}

func (v Variable[T]) Equals(_ *Context[T], e Expression[T]) bool {
	v2, ok := e.ForceRequest().(Variable[T])
	if !ok {
		return false
	}
	return varEquals(v, v2)
}

func (v Variable[T]) StrictEquals(e Expression[T]) bool {
	v2, ok := e.(Variable[T])
	if !ok {
		return false
	}
	return varEquals(v, v2)
}

func (v Variable[T]) String() string {
	return v.name.GetName()
}

func (v Variable[T]) StrictString() string {
	return v.name.GetName() + "[" + strconv.Itoa(v.depth) + "]"
}

func (v Variable[T]) ForceRequest() Expression[T] { return v }
