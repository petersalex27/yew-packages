package expr

type Context struct {
	table    map[Const]Expression
	inverses map[Const]Const
}

func NewContext() *Context {
	cxt := new(Context)
	cxt.table = make(map[Const]Expression)
	cxt.inverses = make(map[Const]Const)
	return cxt
}

func (cxt *Context) GetInverse(e Expression) (out Expression, ok bool) {
	var c, invC Const
	out = nil
	if c, ok = e.(Const); !ok {
		return
	} else if invC, ok = cxt.inverses[c]; !ok {
		return
	}

	out, ok = cxt.table[invC]
	return
}

func (cxt *Context) AddName(name Const, e Expression) error {
	if _, found := cxt.table[name]; found {
		return redefineNameInTable(name)
	}

	cxt.table[name] = e
	return nil
}

func (cxt *Context) DeclareInverse(f Const, invF Const) error {
	if _, found := cxt.table[f]; !found {
		return nameNotDefined(f)
	} else if _, found = cxt.table[invF]; !found {
		return nameNotDefined(invF)
	} else if _, found = cxt.inverses[f]; found {
		return redefineInv(f)
	} else if _, found = cxt.inverses[invF]; found {
		return redefineInv(invF)
	}

	cxt.inverses[f], cxt.inverses[invF] = invF, f
	return nil
}
