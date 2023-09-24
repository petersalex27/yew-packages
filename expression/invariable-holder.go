package expr

type ExternalContainer struct { InvariableExpression }

func Import(e InvariableExpression) ExternalContainer {
	return ExternalContainer{InvariableExpression: e}
}

func (xc ExternalContainer) ForceRequest() Expression {
	return xc
}

func (xc ExternalContainer) String() string {
	return xc.InvariableExpression.String()
}

func (xc ExternalContainer) Equals(e Expression) bool {
	return xc.InvariableExpression.Equals(e)
}

func (xc ExternalContainer) StrictString() string { return xc.InvariableExpression.String() }

func (xc ExternalContainer) StrictEquals(e Expression) bool { return xc.Equals(e) }

func (xc ExternalContainer) Replace(Variable, Expression) (Expression, bool) { return xc, false }

func (xc ExternalContainer) UpdateVars(gt int, by int) Expression { return xc }

func (xc ExternalContainer) Again() (Expression, bool) { return xc, false }

func (xc ExternalContainer) Bind(BindersOnly) Expression { return xc }

func (xc ExternalContainer) Find(Variable) bool { return false }

func (xc ExternalContainer) PrepareAsRHS() Expression { return xc }

func (xc ExternalContainer) Rebind() Expression { return xc }

func (xc ExternalContainer) Copy() Expression {
	return Import(xc.InvariableExpression.DeepCopy())
}