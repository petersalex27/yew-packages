package types

func (cxt *Context[T]) Function(left, right Monotyped[T]) Application[T] {
	return Application[T]{
		c:  cxt.InfixCon("->"),
		ts: []Monotyped[T]{left, right},
	}
}

func (cxt *Context[T]) Cons(left, right Monotyped[T]) Application[T] {
	return Application[T]{
		c:  cxt.InfixCon("&"),
		ts: []Monotyped[T]{left, right},
	}
}

func (cxt *Context[T]) Join(left, right Monotyped[T]) Application[T] {
	return Application[T]{
		c:  cxt.InfixCon("|"),
		ts: []Monotyped[T]{left, right},
	}
}

func (cxt *Context[T]) Infix(left Monotyped[T], constant string, rights ...Monotyped[T]) Application[T] {
	return Application[T]{
		c:  cxt.InfixCon(constant),
		ts: append([]Monotyped[T]{left}, rights...),
	}
}
