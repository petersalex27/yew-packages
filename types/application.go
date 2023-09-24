package types

import (
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/str"
)

type Application struct {
	c Constant
	ts []Monotyped
}

func Apply(c Constant, ts ...Monotyped) Application {
	return Application{ c: c, ts: ts, }
}

func App(name string, ts ...Monotyped) Application {
	return Application{ c: Constant(name), ts: ts, }
}

func (a Application) String() string {
	return "(" + a.c.String() + " " + str.Join(a.ts, str.String(" ")) + ")"
}

func (a Application) Replace(v Variable, m Monotyped) Monotyped {
	f := func(mono Monotyped) Monotyped { return mono.Replace(v, m) }
	return Apply(a.c, fun.FMap(a.ts, f)...)
}

func (a Application) Split() (string, []Monotyped) { return string(a.c), a.ts }

func (a Application) ReplaceDependent(v Variable, m Monotyped) DependentTyped {
	f := func(mono Monotyped) Monotyped { return mono.Replace(v, m) }
	return Apply(a.c, fun.FMap(a.ts, f)...)
}

func (a Application) ReplaceKindVar(replacing Variable, with Monotyped) Monotyped {
	f := func(m Monotyped) Monotyped { return m.ReplaceKindVar(replacing, with) }
	return Apply(a.c, fun.FMap(a.ts, f)...)
}

func (a Application) FreeInstantiation() DependentTyped {
	f := func(m Monotyped) Monotyped { return m.FreeInstantiation().(Monotyped) }
	return Apply(a.c, fun.FMap(a.ts, f)...)
}

func (a Application) Generalize() Polytype {
	return Polytype{
		typeBinders: MakeDummyVars(1),
		bound: a,
	}
}

func (a Application) Equals(t Type) bool {
	a2, ok := t.(Application)
	if !ok {
		return false
	}

	if a.c != a2.c || len(a.ts) != len(a2.ts) {
		return false
	}

	for i := range a.ts {
		if !a.ts[i].Equals(a2.ts[i]) {
			return false
		}
	}
	return true
}