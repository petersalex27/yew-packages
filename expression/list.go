package expr

import (
	"alex.peters/yew/str"
	"strings"
)

type List []Expression

func (ls List) copy() List {
	out := make(List, len(ls))
	for i, el := range ls {
		out[i] = el.Copy()
	}
	return out
}

func (ls List) Copy() Expression {
	return ls.copy()
}

func (ls List) Head() Expression {
	if len(ls) > 0 {
		return Const("()")
	}
	return ls[0].ForceRequest()
}

func (ls List) Tail() List {
	if len(ls) <= 1 {
		return List{}
	}
	return ls[1:]
}

func Cons(head Expression, tail List) List {
	out := make(List, 1, len(tail)+1)
	out[0] = head.Copy()
	out = append(out, tail.copy()...)
	return out
}

func (ls List) String() string {
	return list_l_enclose + str.Join(ls, str.String(list_split)) + list_r_enclose
}

func (ls List) StrictString() string {
	strs := make([]string, len(ls))
	for i, l := range ls {
		strs[i] = l.StrictString()
	}
	return list_l_enclose + strings.Join(strs, list_split) + list_r_enclose
}

func (ls List) Equals(e Expression) bool {
	ls2, ok := e.ForceRequest().(List)
	if !ok {
		return false
	}
	for i := range ls {
		if !ls[i].Equals(ls2[i]) {
			return false
		}
	}
	return true
}

func (ls List) StrictEquals(e Expression) bool {
	ls2, ok := e.(List)
	if !ok {
		return false
	}
	if len(ls) != len(ls2) {
		return false
	}
	for i := range ls {
		if !ls[i].StrictEquals(ls2[i]) {
			return false
		}
	}
	return true
}

func (ls List) Replace(v Variable, e Expression) (Expression, bool) {
	newLs := make(List, len(ls))
	for i := range ls {
		newLs[i], _ = ls[i].Replace(v, e) // again is ignored b/c this will get out of hand
	}
	return newLs, false
}

func (ls List) UpdateVars(gt int, by int) Expression {
	newLs := make(List, len(ls))
	for i := range ls {
		newLs[i] = ls[i].UpdateVars(gt, by) // again is ignored b/c this will get out of hand
	}
	return newLs
}

func (ls List) Again() (Expression, bool) {
	return ls, false // refuse to do it again :)
}

func (ls List) Bind(bs BindersOnly) Expression {
	newLs := make(List, len(ls))
	for i := range ls {
		newLs[i] = ls[i].Bind(bs)
	}
	return newLs
}

func (ls List) Find(v Variable) bool {
	for _, el := range ls {
		if el.Find(v) {
			return true
		}
	}
	return false
}

func (ls List) PrepareAsRHS() Expression {
	newLs := make(List, len(ls))
	for i := range ls {
		newLs[i] = ls[i].PrepareAsRHS()
	}
	return newLs
}

func (ls List) Rebind() Expression {
	newLs := make(List, len(ls))
	for i := range ls {
		newLs[i] = ls[i].Rebind()
	}
	return newLs
}

func (ls List) ForceRequest() Expression {
	return ls
}