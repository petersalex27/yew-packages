package expr

import "strconv"

type Variable struct {
	name  string
	depth int
}

func (v Variable) copy() Variable {
	return Variable{
		name:  v.name,
		depth: v.depth,
	}
}

func (v Variable) Copy() Expression {
	return v.copy()
}

func makeVar(name string, depth int) Variable {
	return Variable{name: name, depth: depth}
}

func (v Variable) PrepareAsRHS() Expression {
	if v.depth < 1 {
		return Variable{
			name:  v.name,
			depth: 1,
		}
	}
	return v
}

func (v Variable) UpdateVars(gt int, by int) Expression {
	if v.depth > gt {
		newVar := Var(v.name)
		newVar.depth = v.depth + by
		return newVar
	}
	return v
}

func (v Variable) Rebind() Expression {
	return Var(v.name)
}

func (v Variable) Bind(bs BindersOnly) Expression {
	depth := len(bs)
	if v.depth != 0 && v.depth <= depth {
		return v
	}

	name := v.name
	out := Var(name)
	// is free variable
	for _, b := range bs {
		if name == b.name {
			// variable gets bound at b.depth
			out.depth = b.depth
			return out
		}
		// variable does not get bound, maybe next binder..?
	}

	// variable remains unbound
	out.depth = v.depth + depth
	if v.depth == 0 {
		// variable is free but unrecognized as free
		out.depth = out.depth + 1 // look at that! +1! variable is recognized :)
	}
	return out
}

func Var(name string) Variable {
	return Variable{name: name, depth: 0}
}

func (v Variable) Again() (Expression, bool) {
	return v, false
}

func (v Variable) Replace(w Variable, e Expression) (Expression, bool) {
	if varEquals(v, w) {
		return e, false
	}
	return v, false
}

func (v Variable) Find(w Variable) bool { return varEquals(v, w) }

func varEquals(v, w Variable) bool {
	return v.depth == w.depth && v.name == w.name
}

func (v Variable) Equals(e Expression) bool {
	v2, ok := e.ForceRequest().(Variable)
	if !ok {
		return false
	}
	return varEquals(v, v2)
}

func (v Variable) StrictEquals(e Expression) bool {
	v2, ok := e.(Variable)
	if !ok {
		return false
	}
	return varEquals(v, v2)
}

func (v Variable) String() string {
	return v.name
}

func (v Variable) StrictString() string {
	return v.name + "[" + strconv.Itoa(v.depth) + "]"
}

func (v Variable) ForceRequest() Expression { return v }
