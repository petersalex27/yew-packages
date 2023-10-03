package expr

import "github.com/petersalex27/yew-packages/nameable"

type ExternalContainer[T nameable.Nameable] struct { InvariableExpression[T] }

func Import[T nameable.Nameable](e InvariableExpression[T]) ExternalContainer[T] {
	return ExternalContainer[T]{InvariableExpression: e}
}

func (xc ExternalContainer[T]) ForceRequest() Expression[T] {
	return xc
}

func (xc ExternalContainer[T]) String() string {
	return xc.InvariableExpression.String()
}

func (xc ExternalContainer[T]) Equals(_ *Context[T], e Expression[T]) bool {
	return xc.InvariableExpression.Equals(e)
}

func (xc ExternalContainer[T]) StrictString() string { return xc.InvariableExpression.String() }

func (xc ExternalContainer[T]) StrictEquals(e Expression[T]) bool { return xc.Equals(nil, e) }

func (xc ExternalContainer[T]) Replace(Variable[T], Expression[T]) (Expression[T], bool) { return xc, false }

func (xc ExternalContainer[T]) UpdateVars(gt int, by int) Expression[T] { return xc }

func (xc ExternalContainer[T]) Again() (Expression[T], bool) { return xc, false }

func (xc ExternalContainer[T]) Bind(BindersOnly[T]) Expression[T] { return xc }

func (xc ExternalContainer[T]) Find(Variable[T]) bool { return false }

func (xc ExternalContainer[T]) PrepareAsRHS() Expression[T] { return xc }

func (xc ExternalContainer[T]) Rebind() Expression[T] { return xc }

func (xc ExternalContainer[T]) Copy() Expression[T] {
	return Import(xc.InvariableExpression.DeepCopy())
}