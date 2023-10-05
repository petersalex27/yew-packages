package types

import (
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type InfixApplication[T nameable.Nameable] Application[T]

func (cxt *Context[T]) Function(left, right Monotyped[T]) InfixApplication[T] {
	return InfixApplication[T]{
		c: cxt.Con("->"),
		ts: []Monotyped[T]{left, right},
	}
}

func (cxt *Context[T]) Cons(left, right Monotyped[T]) InfixApplication[T] {
	return InfixApplication[T]{
		c: cxt.Con("&"),
		ts: []Monotyped[T]{left, right},
	}
}

func (cxt *Context[T]) Join(left, right Monotyped[T]) InfixApplication[T] {
	return InfixApplication[T]{
		c: cxt.Con("|"),
		ts: []Monotyped[T]{left, right},
	}
}

func (cxt *Context[T]) Infix(left Monotyped[T], constant string, rights ...Monotyped[T]) InfixApplication[T] {
	return InfixApplication[T]{
		c: cxt.Con(constant),
		ts: append([]Monotyped[T]{left}, rights...),
	}
}

func (a InfixApplication[T]) Split() (string, []Monotyped[T]) { 
	return Application[T](a).Split()	
}

func (a InfixApplication[T]) String() string {
	length := len(a.ts)
	if length < 2 {
		name := "(" + a.c.String() + ")"
		if length == 0 {
			return name
		} // else length == 1
		return "(" + name + " " + a.ts[0].String() + ")"
	}
	return "(" + a.ts[0].String() + " " + a.c.String() + " " + str.Join(a.ts[1:], str.String(" ")) + ")"
}


func (a InfixApplication[T]) Replace(v Variable[T], m Monotyped[T]) Monotyped[T] {
	res, _ := Application[T](a).Replace(v, m).(Application[T])
	return InfixApplication[T](res)
}

func (a InfixApplication[T]) ReplaceDependent(v Variable[T], m Monotyped[T]) DependentTyped[T] {
	res, _ := Application[T](a).ReplaceDependent(v, m).(Application[T])
	return InfixApplication[T](res)
}

func (a InfixApplication[T]) ReplaceKindVar(replacing Variable[T], with Monotyped[T]) Monotyped[T] {
	res, _ := Application[T](a).ReplaceKindVar(replacing, with).(Application[T])
	return InfixApplication[T](res)
}

func (a InfixApplication[T]) FreeInstantiation(cxt *Context[T]) DependentTyped[T] {
	res, _ := Application[T](a).FreeInstantiation(cxt).(Application[T])
	return InfixApplication[T](res)
}

func (a InfixApplication[T]) Generalize(cxt *Context[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: cxt.MakeDummyVars(1),
		bound: a,
	}
}

func (a InfixApplication[T]) Collect() []T {
	res := make([]T, 0, len(a.ts) + 1)
	res = append(res, a.c.getName())
	for _, t := range a.ts {
		res = append(res, t.Collect()...)
	}
	return res
}

func (a InfixApplication[T]) Equals(t Type[T]) bool {
	a2, ok := t.(InfixApplication[T])
	if !ok {
		return false
	}

	if a.c.getName().GetName() != a2.c.getName().GetName() || len(a.ts) != len(a2.ts) {
		return false
	}

	for i := range a.ts {
		if !a.ts[i].Equals(a2.ts[i]) {
			return false
		}
	}
	return true
}