package types

import (
	"alex.peters/yew/fun"
	"strconv"
)

type Variable struct {
	boundContext int32
	name         string
}

// ReplaceKindVar implements Monotyped.
func (v Variable) ReplaceKindVar(replacing Variable, with Monotyped) Monotyped {
	if varEquals(v, replacing) {
		return with
	}
	return v
}

func NonBindableVar(name string) Variable {
	return Variable{boundContext: -1, name: name}
}

func (v Variable) BoundIn(i int32) Variable {
	v.boundContext = i
	return v
}

func Var(name string) Variable {
	return Variable{boundContext: 0, name: name}
}

func FreeVar(name string) Variable {
	return Var(name)
}

func dummyName(Variable) Variable { return Var("_") }

// MakeDummyVars(n) = []Variable{v /*0*/, v /*1*/, .., v /*n-1*/} where
// v = Var("_")
func MakeDummyVars(n uint) []Variable {
	return fun.FMap(make([]Variable, n), dummyName)
}

// Var("a").BoundBy(x).String() = "a"
func (v Variable) String() string {
	return v.name
}

func (v Variable) String2() (string, string) {
	return v.String(), "#(" + strconv.Itoa(int(v.boundContext)) + ")"
}

func (v Variable) ExtendedString() string {
	l, r := v.String2()
	return l + r
}

// v.Capture(_) = v
func (v Variable) Capture(w Variable) Monotyped { return v }

// Replace all `w` in `v` with `m`; because v is just a variable,
// if v == w, then return m; else return v. Formally, `v [w := m]`
func (v Variable) Replace(w Variable, m Monotyped) Monotyped {
	if varEquals(v, w) {
		return m
	}
	return v
}

func (v Variable) ReplaceDependent(w Variable, m Monotyped) DependentTyped {
	return v.Replace(w, m)
}

// Var("a").Generalize() = `forall _ . a`
func (v Variable) Generalize() Polytype {
	return Polytype{
		typeBinders: MakeDummyVars(1),
		bound:       v,
	}
}

func (v Variable) FreeInstantiation() DependentTyped {
	if v.boundContext > 0 {
		return Var("_")
	}
	return v
}

func varEquals(v Variable, w Variable) bool {
	return v.boundContext == w.boundContext && v.name == w.name
}

func (v Variable) Equals(t Type) bool {
	v2, ok := t.(Variable)
	return ok && varEquals(v, v2)
}
