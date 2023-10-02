package types

import "github.com/petersalex27/yew-packages/nameable"

type Constant[T nameable.Nameable] struct{ name T }

func MakeConst[T nameable.Nameable](t T) Constant[T] {
	return Constant[T]{}
}

func (cxt *Context[T]) Con(name string) Constant[T] {
	return Constant[T]{cxt.makeName(name)}
}

// just returns receiver `c`
func (c Constant[T]) ReplaceKindVar(replacing Variable[T], with Monotyped[T]) Monotyped[T] {
	return c
}

// Constant(x).Equals(y) = true iff y.(Constant) is true and string(y.(Constant)) == x
func (c Constant[T]) Equals(t Type[T]) bool {
	c2, ok := t.(Constant[T])
	return ok && c.name.GetName() == c2.name.GetName()
}

// Constant("Type[T]").String() = "Type[T]"
func (c Constant[T]) String() string {
	return c.name.GetName()
}

// Constant("Type[T]").Generalize() = `forall _ . Type[T]`
func (c Constant[T]) Generalize(cxt *Context[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: []Variable[T]{cxt.dummyName(Variable[T]{})},
		bound:       c,
	}
}

// c.Replace(_, _) = c
func (c Constant[T]) Replace(v Variable[T], m Monotyped[T]) Monotyped[T] { return c }

// c.FreeInstantiation() = c
func (c Constant[T]) FreeInstantiation(*Context[T]) DependentTyped[T] { return c }

func (c Constant[T]) ReplaceDependent(v Variable[T], m Monotyped[T]) DependentTyped[T] {
	return c
}

func (c Constant[T]) Collect() []T {
	return []T{c.name}
}
