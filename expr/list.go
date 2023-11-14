package expr

import (
	"strings"

	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
)

type List[T nameable.Nameable] []Expression[T]

func (ls List[T]) Flatten() []Expression[T] {
	f := (Expression[T]).Flatten
	fold := func(l, r []Expression[T]) []Expression[T] {
		return append(l, r...)
	}
	return fun.FoldLeft([]Expression[T]{}, fun.FMap(ls, f), fold)
}

func (list List[T]) ToAlmostPattern() (pat AlmostPattern[T], ok bool) {
	res := fun.FMapFilter(
		list,
		func(e Expression[T]) (out matchable[T], ok bool) {
			var p Patternable[T]
			if p, ok = e.(Patternable[T]); !ok {
				return
			}
			if tmp, ok := p.ToAlmostPattern(); ok {
				out = tmp.pattern
			}
			return
		},
	)

	// check if all elems fmap-ped
	if len(res) != len(list) {
		// not all elems fmap-ped
		ok = false
		return
	}
	return MakeSequence[T](PatternSequenceList, res...).ToAlmostPattern()
}

func (ls List[T]) BodyAbstract(v Variable[T], name Const[T]) Expression[T] {
	return List[T](
		fun.FMap(
			ls,
			func(e Expression[T]) Expression[T] {
				return e.BodyAbstract(v, name)
			},
		),
	)
}

func (ls List[T]) ExtractVariables(gt int) []Variable[T] {
	vars := []Variable[T]{}
	for _, elem := range ls {
		vars = append(vars, elem.ExtractVariables(gt)...)
	}
	return vars
}

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

func (ls List[T]) rawInner(f func(Expression[T])string) []string {
	return fun.FMap(ls, f)
}

func (ls List[T]) String() string {
	sep := listSepString()
	f := (Expression[T]).String
	ls.rawInner(f)
	list := strings.Join(ls.rawInner(f), sep)
	return encloseListString(list)
}

func (ls List[T]) StrictString() string {
	f := (Expression[T]).StrictString
	strs := ls.rawInner(f)
	list := strings.Join(strs, listSepString())
	return encloseListString(list)
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