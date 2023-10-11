package expr

import "github.com/petersalex27/yew-packages/nameable"

type Const[T nameable.Nameable] struct {Name T}

func (Const[T]) ExtractFreeVariables(Variable[T]) []Variable[T] {
	return []Variable[T]{}
}

func (c Const[T]) Collect() []T {
	return []T{c.Name}
}

func (c Const[T]) ForceRequest() Expression[T] { return c }

func constEquals[T nameable.Nameable](c1, c2 Const[T]) bool {
	return c1.Name.GetName() == c2.Name.GetName()
}

func (c Const[T]) Equals(cxt *Context[T], e Expression[T]) bool {
	if c2, ok := e.ForceRequest().(Const[T]); ok {
		return constEquals(c, c2)
	}
	return false
}

func (c Const[T]) String() string { return c.Name.GetName() }

func (c Const[T]) StrictString() string { return c.Name.GetName() }

func (c Const[T]) Replace(Variable[T], Expression[T]) (Expression[T], bool) { return c, false }

func (c Const[T]) StrictEquals(e Expression[T]) bool { 
	if c2, ok := e.(Const[T]); ok {
		return constEquals(c, c2)
	}
	return false
}

func (c Const[T]) UpdateVars(gt int, by int) Expression[T] { return c }

func (c Const[T]) Again() (Expression[T], bool) { return c, false }

func (c Const[T]) Bind(BindersOnly[T]) Expression[T] { return c }

func (Const[T]) Find(Variable[T]) bool { return false }

func (c Const[T]) PrepareAsRHS() Expression[T] { return c }

func (c Const[T]) Rebind() Expression[T] { return c }

func (c Const[T]) Copy() Expression[T] { return c }
