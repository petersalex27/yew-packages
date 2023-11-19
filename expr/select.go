package expr

import (
	"strings"

	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

// match e in (\x1 .. -> e1) | (\x2 .. -> e2) | .. | (\xN .. -> eN)
type Selection[T nameable.Nameable] struct {
	selector   Expression[T]
	selections []Case[T]
}

func (s Selection[T]) Flatten() []Expression[T] {
	f := (Case[T]).Flatten
	fold := func(l, r []Expression[T]) []Expression[T] {
		return append(l, r...)
	}
	right := fun.FoldLeft([]Expression[T]{}, fun.FMap(s.selections, f), fold)
	return append(s.selector.Flatten(), right...)
}

func (s Selection[T]) ExtractVariables(gt int) []Variable[T] {
	vars := s.selector.ExtractVariables(gt)
	for _, c := range s.selections {
		vars = append(vars, c.ExtractVariables(gt)...)
	}
	return vars
}

func (s Selection[T]) BodyAbstract(v Variable[T], name Const[T]) Expression[T] {
	cases := fun.FMap(
		s.selections,
		func(c Case[T]) Case[T] {
			return c.BodyAbstract(v, name)
		},
	)
	selector := s.selector.BodyAbstract(v, name)
	return Selection[T]{selector, cases}
}

func (s Selection[T]) Collect() []T {
	res := s.selector.Collect()
	for _, sel := range s.selections {
		res = append(res, sel.Collect()...)
	}
	return res
}

func Select[T nameable.Nameable](selector Expression[T], selections ...Case[T]) Selection[T] {
	return Selection[T]{selector: selector, selections: selections}
}

func (s Selection[T]) Merge(selections ...Case[T]) Selection[T] {
	length_s := len(s.selections)
	newSelec := make([]Case[T], length_s+len(selections))
	copy(newSelec, s.selections)
	copy(newSelec[length_s:], selections)
	return Selection[T]{selector: s.selector, selections: selections}
}

func (s Selection[T]) String() string {
	return matchHeadString(s.selector.String(), str.Join(s.selections, str.String(caseSepString())))
}

func (s Selection[T]) StrictString() string {
	head := s.selector.StrictString()
	tail := make([]string, len(s.selections))
	for i := range tail {
		tail[i] = s.selections[i].StrictString()
	}
	return matchHeadString(head, strings.Join(tail, caseSepString()))
}

func (s Selection[T]) Equals(cxt *Context[T], e Expression[T]) bool {
	s2, ok := e.(Selection[T])
	if !ok || len(s.selections) != len(s2.selections) {
		return false
	}

	if !s.selector.Equals(cxt, s2.selector) {
		return false
	}

	for i, c := range s.selections {
		if !c.Equals(cxt, s2.selections[i]) {
			return false
		}
	}

	return true
}

func (s Selection[T]) StrictEquals(e Expression[T]) bool {
	s2, ok := e.(Selection[T])
	if !ok || len(s.selections) != len(s2.selections) {
		return false
	}

	if !s.selector.StrictEquals(s2.selector) {
		return false
	}

	for i, c := range s.selections {
		if !c.StrictEquals(s2.selections[i]) {
			return false
		}
	}

	return true
}

func (s Selection[T]) Replace(v Variable[T], e Expression[T]) (Expression[T], bool) {
	f := func(x Expression[T]) (Expression[T], bool) { return x.Replace(v, e) }
	selector, _ := s.selector.Replace(v, e)
	selections, _ := selectionsMap(s.selections, f)
	return Select(selector, selections...), false
}

func (s Selection[T]) UpdateVars(gt int, by int) Expression[T] {
	f := func(x Expression[T]) (Expression[T], bool) { return x.UpdateVars(gt, by), true }
	selector := s.selector.UpdateVars(gt, by)
	seletions, _ := selectionsMap(s.selections, f)
	return Select(selector, seletions...)
}

func (s Selection[T]) Again() (Expression[T], bool) { return s, false }

func (s Selection[T]) Bind(bs BindersOnly[T]) Expression[T] {
	f := func(x Expression[T]) (Expression[T], bool) { return x.Bind(bs), true }
	selector := s.selector.Bind(bs)
	seletions, _ := selectionsMap(s.selections, f)
	return Select(selector, seletions...)
}

func (s Selection[T]) Find(v Variable[T]) bool {
	if s.selector.Find(v) { // search selector
		return true
	}

	for _, c := range s.selections { // search selections
		if c.Find(v) {
			return true
		}
	}

	return false // v is not found
}

func (s Selection[T]) PrepareAsRHS() Expression[T] {
	return s
}

func (s Selection[T]) Rebind() Expression[T] {
	f := func(x Expression[T]) (Expression[T], bool) { return x.Rebind(), true }
	selector := s.selector.Rebind()
	seletions, _ := selectionsMap(s.selections, f)
	return Select(selector, seletions...)
}

func (s Selection[T]) Copy() Expression[T] {
	f := func(x Expression[T]) (Expression[T], bool) { return x.Copy(), true }
	selector := s.selector.Copy()
	seletions, _ := selectionsMap(s.selections, f)
	return Select(selector, seletions...)
}

func (s Selection[T]) ForceRequest() Expression[T] {
	return s
}
