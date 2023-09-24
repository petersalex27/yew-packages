package types

import "alex.peters/yew/str"

type InfixApplication Application

func Function(left, right Monotyped) InfixApplication {
	return InfixApplication{
		c: Constant("->"),
		ts: []Monotyped{left, right},
	}
}

func Cons(left, right Monotyped) InfixApplication {
	return InfixApplication{
		c: Constant("&"),
		ts: []Monotyped{left, right},
	}
}

func Join(left, right Monotyped) InfixApplication {
	return InfixApplication{
		c: Constant("|"),
		ts: []Monotyped{left, right},
	}
}

func Infix(left Monotyped, constant string, rights ...Monotyped) InfixApplication {
	return InfixApplication{
		c: Constant(constant),
		ts: append([]Monotyped{left}, rights...),
	}
}

func (a InfixApplication) Split() (string, []Monotyped) { 
	return Application(a).Split()	
}

func (a InfixApplication) String() string {
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


func (a InfixApplication) Replace(v Variable, m Monotyped) Monotyped {
	res, _ := Application(a).Replace(v, m).(Application)
	return InfixApplication(res)
}

func (a InfixApplication) ReplaceDependent(v Variable, m Monotyped) DependentTyped {
	res, _ := Application(a).ReplaceDependent(v, m).(Application)
	return InfixApplication(res)
}

func (a InfixApplication) ReplaceKindVar(replacing Variable, with Monotyped) Monotyped {
	res, _ := Application(a).ReplaceKindVar(replacing, with).(Application)
	return InfixApplication(res)
}

func (a InfixApplication) FreeInstantiation() DependentTyped {
	res, _ := Application(a).FreeInstantiation().(Application)
	return InfixApplication(res)
}

func (a InfixApplication) Generalize() Polytype {
	return Polytype{
		typeBinders: MakeDummyVars(1),
		bound: a,
	}
}

func (a InfixApplication) Equals(t Type) bool {
	a2, ok := t.(InfixApplication)
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