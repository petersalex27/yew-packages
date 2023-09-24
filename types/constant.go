package types

type Constant string

func Con(name string) Constant {
	return Constant(name)
}

// just returns receiver `c`
func (c Constant) ReplaceKindVar(replacing Variable, with Monotyped) Monotyped {
	return c 
}

// Constant(x).Equals(y) = true iff y.(Constant) is true and string(y.(Constant)) == x
func (c Constant) Equals(t Type) bool {
	c2, ok := t.(Constant)
	return ok && c == c2
}

// Constant("Type").String() = "Type"
func (c Constant) String() string {
	return string(c)
}

// Constant("Type").Generalize() = `forall _ . Type`
func (c Constant) Generalize() Polytype {
	return Polytype{
		typeBinders: []Variable{Var("_")},
		bound: c,
	}
}

// c.Replace(_, _) = c
func (c Constant) Replace(v Variable, m Monotyped) Monotyped { return c }

// c.FreeInstantiation() = c
func (c Constant) FreeInstantiation() DependentTyped { return c }

func (c Constant) ReplaceDependent(v Variable, m Monotyped) DependentTyped {
	return c
}
