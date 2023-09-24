package expr

type Const string

func (c Const) ForceRequest() Expression { return c }

func constEquals(c1, c2 Const) bool {
	return c1 == c2
}

func (c Const) Equals(e Expression) bool {
	if c2, ok := e.ForceRequest().(Const); ok {
		return constEquals(c, c2)
	}
	return false
}

func (c Const) String() string { return string(c) }

func (c Const) StrictString() string { return string(c) }

func (c Const) Replace(Variable, Expression) (Expression, bool) { return c, false }

func (c Const) StrictEquals(e Expression) bool { 
	if c2, ok := e.(Const); ok {
		return constEquals(c, c2)
	}
	return false
}

func (c Const) UpdateVars(gt int, by int) Expression { return c }

func (c Const) Again() (Expression, bool) { return c, false }

func (c Const) Bind(BindersOnly) Expression { return c }

func (Const) Find(Variable) bool { return false }

func (c Const) PrepareAsRHS() Expression { return c }

func (c Const) Rebind() Expression { return c }

func (c Const) Copy() Expression { return c }
