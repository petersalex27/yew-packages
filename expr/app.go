package expr

type Application struct {
	left  Expression
	right Expression
}

func (a Application) Copy() Expression {
	return Apply(a.left.Copy(), a.right.Copy())
}

func (a Application) PrepareAsRHS() Expression {
	return Apply(a.left.PrepareAsRHS(), a.right.PrepareAsRHS())
}

// right side never gets forced 
func (a Application) ForceRequest() Expression {
	left := a.left.ForceRequest()
	if f, ok := left.(ApplicableExpression); ok {
		return f.DoApplication(a.right)
	}
	return Apply(left, a.right)
}

func Apply(e1, e2 Expression, es ...Expression) Application {
	e := Application{left: e1, right: e2}
	if len(es) > 0 {
		if len(es) > 1 {
			return Apply(e, es[0], es[1:]...)
		} else {
			return Apply(e, es[0])
		}
	}
	return e
}

func (a Application) String() string {
	return "(" + a.left.String() + apply_string + a.right.String() + ")"
}

func (a Application) StrictString() string {
	return "(" + a.left.StrictString() + apply_string + a.right.StrictString() + ")"
}

func applicationEquals(a, b Application) bool {
	return a.left.Equals(b.left) && a.right.Equals(b.right)
}

func (a Application) Equals(e Expression) bool {
	a2, ok := e.ForceRequest().(Application)
	if !ok {
		return false
	}
	return applicationEquals(a, a2)
}

func strictApplicationEquals(a, b Application) bool {
	return a.left.StrictEquals(b.left) && a.right.StrictEquals(b.right)
}

func (a Application) StrictEquals(e Expression) bool {
	a2, ok := e.(Application)
	if !ok {
		return false
	}
	return strictApplicationEquals(a, a2)
}

func (a Application) Rebind() Expression {
	return Apply(a.left.Rebind(), a.right.Rebind())
}

func (a Application) Bind(bs BindersOnly) Expression {
	return Apply(a.left.Bind(bs), a.right.Bind(bs))
}

func (a Application) UpdateVars(gt int, by int) Expression {
	return Apply(a.left.UpdateVars(gt, by), a.right.UpdateVars(gt, by))
}

func (a Application) Find(v Variable) bool {
	return a.left.Find(v) || a.right.Find(v)
}

func (a Application) Again() (Expression, bool) {
	left, lcheck := a.left.Again()
	//right, rcheck := a.right.Again()

	if lcheck {
		return Apply(left, a.right), true
	}

	f, ok := left.(ApplicableExpression)
	if ok {
		var res Expression
		res, ok = f.AgainApply(a.right) // left side is a function, apply the right side to it
		// need to still return whether arg (right) applied to function (left) can be simplified
		return res, ok
	}
	return a, false
}

func (a Application) Split() (left, right Expression) {
	return a.left, a.right
}

func (a Application) Replace(v Variable, e Expression) (Expression, bool) {
	// ((e1 e2) [v in e1 := e]) [v in e2 := e] => apply e1 e2
	left, lcheck := a.left.Replace(v, e)
	right, _ := a.right.Replace(v, e)
	_, ok := left.(ApplicableExpression)
	ok = ok || lcheck
	return Apply(left, right), ok
}
