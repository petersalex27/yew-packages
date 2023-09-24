package expr

import (
	"github.com/petersalex27/yew-packages/str"
	"strings"
)

// e1 of (\x -> e2) else (\y -> e3)
type Selection struct {
	selector Expression
	selections []Case
}

func Select(selector Expression, selections ...Case) Selection {
	return Selection{ selector: selector, selections: selections, }
}

func (s Selection) String() string {
	return s.selector.String() + " of " + str.Join(s.selections, str.String(" else "))
}

func (s Selection) Equals(e Expression) bool {
	s2, ok := e.(Selection)
	if !ok || len(s.selections) != len(s2.selections) {
		return false
	}

	if !s.selector.Equals(s2.selector) {
		return false
	}

	for i, c := range s.selections {
		if !c.Equals(s2.selections[i]) {
			return false
		}
	}

	return true
}

func (s Selection) StrictString() string {
	head := s.selector.StrictString() + " of "
	tail := make([]string, len(s.selections))
	for i := range tail {
		tail[i] = s.selections[i].StrictString()
	}
	return head + strings.Join(tail, " else ")
}

func (s Selection) StrictEquals(e Expression) bool {
	s2, ok := e.(Selection)
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

func (s Selection) Replace(v Variable, e Expression) (Expression, bool) {
	f := func(x Expression) (Expression, bool) { return x.Replace(v, e) }
	selector, ok := s.selector.Replace(v, e)
	if !ok {
		return s, false
	}
	selections, allGood := selectionsMap(s.selections, f)
	if !allGood {
		return s, false
	}

	return Select(selector, selections...), true
}

func (s Selection) UpdateVars(gt int, by int) Expression {
	f := func(x Expression) (Expression, bool) { return x.UpdateVars(gt, by), true }
	selector := s.selector.UpdateVars(gt, by)
	seletions, _ := selectionsMap(s.selections, f)
	return Select(selector, seletions...)
}

func (s Selection) Again() (Expression, bool) { return s, false }

func (s Selection) Bind(bs BindersOnly) Expression {
	f := func(x Expression) (Expression, bool) { return x.Bind(bs), true }
	selector := s.selector.Bind(bs)
	seletions, _ := selectionsMap(s.selections, f)
	return Select(selector, seletions...)
}

func (s Selection) Find(v Variable) bool {
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

func (s Selection) PrepareAsRHS() Expression {
	return s
}

func (s Selection) Rebind() Expression {
	f := func(x Expression) (Expression, bool) { return x.Rebind(), true }
	selector := s.selector.Rebind()
	seletions, _ := selectionsMap(s.selections, f)
	return Select(selector, seletions...)
}

func (s Selection) Copy() Expression {
	f := func(x Expression) (Expression, bool) { return x.Copy(), true }
	selector := s.selector.Copy()
	seletions, _ := selectionsMap(s.selections, f)
	return Select(selector, seletions...)
}

func (s Selection) ForceRequest() Expression {
	return s
}