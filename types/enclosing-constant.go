package types

import "github.com/petersalex27/yew-packages/nameable"

type EnclosingConst[T nameable.Nameable] struct {
	splitAt uint
	Constant[T]
}

func (c EnclosingConst[T]) GetName() T {
	return c.name
}

func MakeEnclosingConst[T nameable.Nameable](at uint, t T) EnclosingConst[T] {
	return EnclosingConst[T]{at, MakeConst(t)}
}

func (cxt *Context[T]) EnclosingCon(at uint, name string) EnclosingConst[T] {
	return MakeEnclosingConst(at, cxt.makeName(name))
}

// just returns receiver `c`
func (c EnclosingConst[T]) ReplaceKindVar(replacing Variable[T], with Monotyped[T]) Monotyped[T] {
	return c
}

// Constant(x).Equals(y) = true iff y.(Constant) is true and string(y.(Constant)) == x
func (c EnclosingConst[T]) Equals(t Type[T]) bool {
	c2, ok := t.(EnclosingConst[T])
	return ok && c.splitAt == c2.splitAt && c.name.GetName() == c2.name.GetName()
}

// (EnclosingConst{1, "[]"}).String() == "[]"
func (c EnclosingConst[T]) String() string {
	return c.name.GetName()
}

// (EnclosingConst{1, "[]"}).String() == ("[", "]")
func (c EnclosingConst[T]) SplitString() (string, string) {
	name := c.name.GetName()
	return name[:c.splitAt], name[c.splitAt:]
}

func (c EnclosingConst[T]) Generalize(cxt *Context[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: []Variable[T]{cxt.dummyName(Variable[T]{})},
		bound:       c,
	}
}

// c.Replace(_, _) = c
func (c EnclosingConst[T]) Replace(v Variable[T], m Monotyped[T]) Monotyped[T] { return c }

// c.FreeInstantiation() = c
func (c EnclosingConst[T]) FreeInstantiation(*Context[T]) DependentTyped[T] { return c }

func (c EnclosingConst[T]) ReplaceDependent(v Variable[T], m Monotyped[T]) DependentTyped[T] {
	return c
}

func (c EnclosingConst[T]) Collect() []T {
	return []T{c.name}
}