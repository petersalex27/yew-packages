package types

func (cxt *Context[T]) Apply(a, b Monotyped[T]) (Monotyped[T], error) {
	v := cxt.NewVar()
	e := cxt.Unify(a, cxt.Function(b, v))
	if e != nil {
		return nil, e
	}
	return v, nil
}

func (cxt *Context[T]) Abstract(t Monotyped[T]) Monotyped[T] {
	return cxt.Function(cxt.NewVar(), t)
}

func (cxt *Context[T]) ConsRule(a, b Monotyped[T]) Monotyped[T] {
	return cxt.Cons(a, b)
}

func (cxt *Context[T]) JoinRule(t Monotyped[T]) Monotyped[T] {
	return cxt.Join(t, cxt.NewVar())
}

func (cxt *Context[T]) Expansion(a, b Monotyped[T]) Monotyped[T] {
	return cxt.Join(a, b)
}

func (cxt *Context[T]) Head(a Monotyped[T]) (Monotyped[T], error) {
	// << Γ:-a  newvar(v1, v2)  unify(a, v1&v2)
	// >> Γ:-v1
	v1, v2 := cxt.NewVar(), cxt.NewVar()
	if e := cxt.Unify(a, cxt.ConsRule(v1, v2)); e != nil {
		return nil, e
	}
	return v1, nil
}

func (cxt *Context[T]) Tail(a Monotyped[T]) (Monotyped[T], error) {
	// << Γ:-a  newvar(v1, v2)  unify(a, v1&v2)
	// >> Γ:-v2
	v1, v2 := cxt.NewVar(), cxt.NewVar()
	if e := cxt.Unify(a, cxt.ConsRule(v1, v2)); e != nil {
		return nil, e
	}
	return v2, nil
}

func (*Context[T]) instantiateDependentTyped(t DependentTyped[T]) Monotyped[T] {
	if d, ok := t.(DependentType[T]); ok {
		return d.KindInstantiation()
	} else if m, ok := t.(Monotyped[T]); ok {
		return m
	}
	panic("tried to declare an unclassifiable Type")
}

func (cxt *Context[T]) Recursive(_ []Monotyped[T], contextualizedType Monotyped[T]) Monotyped[T] {
	return contextualizedType
}

func (cxt *Context[T]) Declare(t Type[T]) Monotyped[T] {
	ty, ok := t.(Polytype[T])
	if !ok {
		return cxt.instantiateDependentTyped(t.(DependentTyped[T]))
	}

	// given a polyType σ = forall a_1 a_2 .. a_n . D
	// given newvar(v_1, v_2, .., v_n)
	// then Declare(σ) = ((D[a_1 := v_1])[a_2 := v_2]..)[a_n := v_n]

	// create a new free variable for each Type binder
	vars := make([]Variable[T], len(ty.typeBinders))
	for i := range vars {
		vars[i] = cxt.NewVar()
	}

	// instantiate the polyType's with a different free variable for 
	// each variable bound by the polyType
	for _, v := range vars {
		ty := t.(Polytype[T])
		t = ty.Instantiate(v)
	}
	
	return cxt.instantiateDependentTyped(t.(DependentTyped[T]))
}

func (cxt *Context[T]) Realization(a, b, c Monotyped[T]) (Monotyped[T], error) {
	// << Γ:-a  Γ:-b  Γ:-c  newvar(v1, v2, v3) unify(a, v1|v2)  unify(b, v1->v3)  unify(c, v2->v3)
	// >> Γ:-v3
	v1, v2, v3 := cxt.NewVar(), cxt.NewVar(), cxt.NewVar()

	if e := cxt.Unify(a, cxt.Join(v1, v2)); e != nil {
		return nil, e
	}

	if e := cxt.Unify(b, cxt.Function(v1, v3)); e != nil {
		return nil, e
	}

	if e := cxt.Unify(c, cxt.Function(v2, v3)); e != nil {
		return nil, e
	}
	return v3, nil
}

func (*Context[T]) Contextualization(_, t Monotyped[T]) Monotyped[T] {
	return t
}