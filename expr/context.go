package expr

import "github.com/petersalex27/yew-packages/nameable"

type Context[T nameable.Nameable] struct {
	table    map[string]Expression[T]
	inverses map[string]Const[T]
	makeName func(string)T
}

func (cxt *Context[T]) SetNameMaker(f func(string)T) *Context[T] {
	cxt.makeName = f
	return cxt
}

func NewContext[T nameable.Nameable]() *Context[T] {
	cxt := new(Context[T])
	cxt.table = make(map[string]Expression[T])
	cxt.inverses = make(map[string]Const[T])
	return cxt
}

func (cxt *Context[T]) GetInverse(e Expression[T]) (out Expression[T], ok bool) {
	var c, invC Const[T]
	out = nil
	if c, ok = e.(Const[T]); !ok {
		return
	} else if invC, ok = cxt.inverses[c.String()]; !ok {
		return
	}

	out, ok = cxt.table[invC.String()]
	return
}

func (cxt *Context[T]) AddName(name Const[T], e Expression[T]) error {
	if _, found := cxt.table[name.String()]; found {
		return redefineNameInTable(name)
	}

	cxt.table[name.String()] = e
	return nil
}
//[T nameable.Nameable]
func (cxt *Context[T]) DeclareInverse(f Const[T], invF Const[T]) error {
	if _, found := cxt.table[f.String()]; !found {
		return nameNotDefined(f)
	} else if _, found = cxt.table[invF.String()]; !found {
		return nameNotDefined(invF)
	} else if _, found = cxt.inverses[f.String()]; found {
		return redefineInv(f)
	} else if _, found = cxt.inverses[invF.String()]; found {
		return redefineInv(invF)
	}

	cxt.inverses[f.String()], cxt.inverses[invF.String()] = invF, f
	return nil
}
