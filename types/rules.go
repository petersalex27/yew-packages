package types

func (cxt *Context) Apply(a, b Monotyped) (Monotyped, error) {
	v := cxt.NewVar()
	e := cxt.Unify(a, Function(b, v))
	if e != nil {
		return nil, e
	}
	return v, nil
}

func (cxt *Context) Abstract(t Monotyped) Monotyped {
	return Function(cxt.NewVar(), t)
}

func (*Context) Cons(a, b Monotyped) Monotyped {
	return Cons(a, b)
}

func (cxt *Context) Join(t Monotyped) Monotyped {
	return Join(t, cxt.NewVar())
}

func (cxt *Context) Expansion(a, b Monotyped) Monotyped {
	return Join(a, b)
}

func (cxt *Context) Head(a Monotyped) (Monotyped, error) {
	// << Γ:-a  newvar(v1, v2)  unify(a, v1&v2)
	// >> Γ:-v1
	v1, v2 := cxt.NewVar(), cxt.NewVar()
	if e := cxt.Unify(a, Cons(v1, v2)); e != nil {
		return nil, e
	}
	return v1, nil
}

func (cxt *Context) Tail(a Monotyped) (Monotyped, error) {
	// << Γ:-a  newvar(v1, v2)  unify(a, v1&v2)
	// >> Γ:-v2
	v1, v2 := cxt.NewVar(), cxt.NewVar()
	if e := cxt.Unify(a, Cons(v1, v2)); e != nil {
		return nil, e
	}
	return v2, nil
}

func (*Context) instantiateDependentTyped(t DependentTyped) Monotyped {
	if d, ok := t.(DependentType); ok {
		return d.KindInstantiation()
	} else if m, ok := t.(Monotyped); ok {
		return m
	}
	panic("tried to declare an unclassifiable type")
}

func (cxt *Context) Declare(t Type) Monotyped {
	ty, ok := t.(Polytype)
	if !ok {
		return cxt.instantiateDependentTyped(t.(DependentTyped))
	}

	// given a polytype σ = forall a_1 a_2 .. a_n . D
	// given newvar(v_1, v_2, .., v_n)
	// then Declare(σ) = ((D[a_1 := v_1])[a_2 := v_2]..)[a_n := v_n]

	// create a new free variable for each type binder
	vars := make([]Variable, len(ty.typeBinders))
	for i := range vars {
		vars[i] = cxt.NewVar()
	}

	// instantiate the polytype's with a different free variable for 
	// each variable bound by the polytype
	for _, v := range vars {
		ty := t.(Polytype)
		t = ty.Instantiate(v)
	}
	
	return cxt.instantiateDependentTyped(t.(DependentTyped))
}

func (cxt *Context) Realization(a, b, c Monotyped) (Monotyped, error) {
	// << Γ:-a  Γ:-b  Γ:-c  newvar(v1, v2, v3) unify(a, v1|v2)  unify(b, v1->v3)  unify(c, v2->v3)
	// >> Γ:-v3
	v1, v2, v3 := cxt.NewVar(), cxt.NewVar(), cxt.NewVar()

	if e := cxt.Unify(a, Join(v1, v2)); e != nil {
		return nil, e
	}

	if e := cxt.Unify(b, Function(v1, v3)); e != nil {
		return nil, e
	}

	if e := cxt.Unify(c, Function(v2, v3)); e != nil {
		return nil, e
	}
	return v3, nil
}

func (*Context) Contextualization(p, q Polytype) Type {
	if len(q.typeBinders) == 0 {
		return q.bound
	}
	return q
}