package expr

import "github.com/petersalex27/yew-packages/nameable"

type Application[T nameable.Nameable] struct {
	left  Expression[T]
	right Expression[T]
}

func (a Application[T]) Copy() Expression[T] {
	return Apply(a.left.Copy(), a.right.Copy())
}

func (a Application[T]) PrepareAsRHS() Expression[T] {
	return Apply(a.left.PrepareAsRHS(), a.right.PrepareAsRHS())
}

// right side never gets forced 
func (a Application[T]) ForceRequest() Expression[T] {
	left := a.left.ForceRequest()
	if f, ok := left.(ApplicableExpression[T]); ok {
		return f.DoApplication(a.right)
	}
	return Apply(left, a.right)
}

func Apply[T nameable.Nameable](e1, e2 Expression[T], es ...Expression[T]) Application[T] {
	e := Application[T]{left: e1, right: e2}
	if len(es) > 0 {
		if len(es) > 1 {
			return Apply[T](e, es[0], es[1:]...)
		} else {
			return Apply[T](e, es[0])
		}
	}
	return e
}

func (a Application[T]) String() string {
	return "(" + a.left.String() + apply_string + a.right.String() + ")"
}

func (a Application[T]) StrictString() string {
	return "(" + a.left.StrictString() + apply_string + a.right.StrictString() + ")"
}

func applicationEquals[T nameable.Nameable](cxt *Context[T], a, b Application[T]) bool {
	return a.left.Equals(cxt, b.left) && a.right.Equals(cxt, b.right)
}

func (a Application[T]) Equals(cxt *Context[T], e Expression[T]) bool {
	a2, ok := e.ForceRequest().(Application[T])
	if !ok {
		return false
	}
	return applicationEquals(cxt, a, a2)
}

func strictApplicationEquals[T nameable.Nameable](a, b Application[T]) bool {
	return a.left.StrictEquals(b.left) && a.right.StrictEquals(b.right)
}

func (a Application[T]) StrictEquals(e Expression[T]) bool {
	a2, ok := e.(Application[T])
	if !ok {
		return false
	}
	return strictApplicationEquals(a, a2)
}

func (a Application[T]) Rebind() Expression[T] {
	return Apply(a.left.Rebind(), a.right.Rebind())
}

func (a Application[T]) Bind(bs BindersOnly[T]) Expression[T] {
	return Apply(a.left.Bind(bs), a.right.Bind(bs))
}

func (a Application[T]) UpdateVars(gt int, by int) Expression[T] {
	return Apply(a.left.UpdateVars(gt, by), a.right.UpdateVars(gt, by))
}

func (a Application[T]) Find(v Variable[T]) bool {
	return a.left.Find(v) || a.right.Find(v)
}

func (a Application[T]) Again() (Expression[T], bool) {
	left, lcheck := a.left.Again()
	//right, rcheck := a.right.Again()

	if lcheck {
		return Apply(left, a.right), true
	}

	f, ok := left.(ApplicableExpression[T])
	if ok {
		var res Expression[T]
		res, ok = f.AgainApply(a.right) // left side is a function[T], apply the right side to it
		// need to still return whether arg (right) applied to function[T] (left) can be simplified
		return res, ok
	}
	return a, false
}

func (a Application[T]) Split() (left, right Expression[T]) {
	return a.left, a.right
}

func (a Application[T]) Replace(v Variable[T], e Expression[T]) (Expression[T], bool) {
	// ((e1 e2) [v in e1 := e]) [v in e2 := e] => apply e1 e2
	left, lcheck := a.left.Replace(v, e)
	right, _ := a.right.Replace(v, e)
	_, ok := left.(ApplicableExpression[T])
	ok = ok || lcheck
	return Apply(left, right), ok
}
