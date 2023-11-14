package types

import "github.com/petersalex27/yew-packages/nameable"

type InfixConst[T nameable.Nameable] Constant[T]

func (c InfixConst[T]) GetReferred() T {
	return c.name
}

func (c InfixConst[T]) GetName() string {
	return c.name.GetName()
}

func MakeInfixConst[T nameable.Nameable](t T) InfixConst[T] {
	return InfixConst[T]{t}
}

func (cxt *Context[T]) InfixCon(name string) InfixConst[T] {
	return InfixConst[T]{cxt.makeName(name)}
}

// just returns receiver `c`
func (c InfixConst[T]) ReplaceKindVar(replacing Variable[T], with Monotyped[T]) Monotyped[T] {
	return c
}

func (c InfixConst[T]) GetFreeVariables() []Variable[T] {
	return []Variable[T]{}
}

// Constant(x).Equals(y) = true iff y.(Constant) is true and string(y.(Constant)) == x
func (c InfixConst[T]) Equals(t Type[T]) bool {
	c2, ok := t.(InfixConst[T])
	return ok && c.name.GetName() == c2.name.GetName()
}

// InfixCon("Type").String() = "(Type)"
func (c InfixConst[T]) String() string {
	return "(" + c.name.GetName() + ")"
}

// InfixCon("Type").Generalize() = `forall _ . (Type)`
func (c InfixConst[T]) Generalize(cxt *Context[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: []Variable[T]{cxt.dummyName(Variable[T]{})},
		bound:       c,
	}
}

// c.Replace(_, _) = c
func (c InfixConst[T]) Replace(v Variable[T], m Monotyped[T]) Monotyped[T] { return c }

// c.FreeInstantiation() = c
func (c InfixConst[T]) FreeInstantiation(*Context[T]) Monotyped[T] { return c }

func (c InfixConst[T]) ReplaceDependent(vs []Variable[T], ms []Monotyped[T]) Monotyped[T] {
	return c
}

func (c InfixConst[T]) Collect() []T {
	return []T{c.name}
}