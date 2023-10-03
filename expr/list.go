package expr

import (
	"strings"

	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type List[T nameable.Nameable] []Expression[T]

func (ls List[T]) Collect() []T {
	if len(ls) == 0 {
		return []T{}
	}
	res := ls[0].Collect()
	for i := 1; i < len(ls); i++ {
		res = append(res, ls[i].Collect()...)
	}
	return res
}

func (ls List[T]) copy() List[T] {
	out := make(List[T], len(ls))
	for i, el := range ls {
		out[i] = el.Copy()
	}
	return out
}

func (ls List[T]) Copy() Expression[T] {
	return ls.copy()
}

func (ls List[T]) Head() Expression[T] {
	if len(ls) > 0 {
		return nil
	}
	return ls[0].ForceRequest()
}

func (ls List[T]) Tail() List[T] {
	if len(ls) <= 1 {
		return List[T]{}
	}
	return ls[1:]
}

func Cons[T nameable.Nameable](head Expression[T], tail List[T]) List[T] {
	out := make(List[T], 1, len(tail)+1)
	out[0] = head.Copy()
	out = append(out, tail.copy()...)
	return out
}

func (ls List[T]) String() string {
	return list_l_enclose + str.Join(ls, str.String(list_split)) + list_r_enclose
}

func (ls List[T]) StrictString() string {
	strs := make([]string, len(ls))
	for i, l := range ls {
		strs[i] = l.StrictString()
	}
	return list_l_enclose + strings.Join(strs, list_split) + list_r_enclose
}

func (ls List[T]) Equals(cxt *Context[T], e Expression[T]) bool {
	ls2, ok := e.ForceRequest().(List[T])
	if !ok {
		return false
	}
	for i := range ls {
		if !ls[i].Equals(cxt, ls2[i]) {
			return false
		}
	}
	return true
}

func (ls List[T]) StrictEquals(e Expression[T]) bool {
	ls2, ok := e.(List[T])
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

func (ls List[T]) Replace(v Variable[T], e Expression[T]) (Expression[T], bool) {
	newLs := make(List[T], len(ls))
	for i := range ls {
		newLs[i], _ = ls[i].Replace(v, e) // again is ignored b/c this will get out of hand
	}
	return newLs, false
}

func (ls List[T]) UpdateVars(gt int, by int) Expression[T] {
	newLs := make(List[T], len(ls))
	for i := range ls {
		newLs[i] = ls[i].UpdateVars(gt, by) // again is ignored b/c this will get out of hand
	}
	return newLs
}

func (ls List[T]) Again() (Expression[T], bool) {
	return ls, false // refuse to do it again :)
}

func (ls List[T]) Bind(bs BindersOnly[T]) Expression[T] {
	newLs := make(List[T], len(ls))
	for i := range ls {
		newLs[i] = ls[i].Bind(bs)
	}
	return newLs
}

func (ls List[T]) Find(v Variable[T]) bool {
	for _, el := range ls {
		if el.Find(v) {
			return true
		}
	}
	return false
}

func (ls List[T]) PrepareAsRHS() Expression[T] {
	newLs := make(List[T], len(ls))
	for i := range ls {
		newLs[i] = ls[i].PrepareAsRHS()
	}
	return newLs
}

func (ls List[T]) Rebind() Expression[T] {
	newLs := make(List[T], len(ls))
	for i := range ls {
		newLs[i] = ls[i].Rebind()
	}
	return newLs
}

func (ls List[T]) ForceRequest() Expression[T] {
	return ls
}