package expr

import (
	"strings"
)

type Case struct {
	binders []Variable
	when Expression
	then Expression
}

type PartialCase_when Case

func (bs BindersOnly) InCase(when Expression, then Expression) Case {
	out := Case{ binders: bs, }
	out.when = when.Bind(bs)
	out.then = then.Bind(bs)
	return out
}

func (bs BindersOnly) When(e Expression) PartialCase_when {
	return PartialCase_when{
		binders: bs,
		when: e.Bind(bs),
	}
}

func (pcw PartialCase_when) Then(e Expression) Case {
	pcw.then = e.Bind(pcw.binders)
	return Case(pcw)
}

func (c Case) String() string {
	return c.when.String() + " -> " + c.then.String()
}

func (c Case) StrictString() string {
	strs := make([]string, len(c.binders))
	for i, v := range c.binders {
		strs[i] = v.StrictString()
	}
	return "Λ" + strings.Join(strs, " ") + " . " + c.when.StrictString() + " -> " + c.then.StrictString()
}

func (c Case) Equals(k Case) bool {
	if len(c.binders) != len(k.binders) {
		return false
	}

	ok := c.when.Equals(k.when) && c.then.Equals(k.then)
	if !ok {
		return false
	}

	for i, b := range c.binders {
		if !varEquals(b, k.binders[i]) {
			return false
		}
	}
	return true
}

func (c Case) StrictEquals(k Case) bool {
	if len(c.binders) != len(k.binders) {
		return false
	}

	ok := c.when.StrictEquals(k.when) && c.then.StrictEquals(k.then)
	if !ok {
		return false
	}

	for i, b := range c.binders {
		if !b.StrictEquals(k.binders[i]) {
			return false
		}
	}
	return true
}

func selectionsMap(selections []Case, f func(Expression) (Expression, bool)) ([]Case, bool) {
	out := make([]Case, len(selections))
	for i, c := range selections {
		var when, then Expression
		var ok bool
		when, ok = f(c.when)
		if !ok {
			return selections, false
		}
		then, ok = f(c.then)
		if !ok {
			return selections, false
		}
		out[i] = Case{when: when, then: then}
	}
	return out, true
}

func (c Case) Find(v Variable) bool {
	return c.when.Find(v) || c.then.Find(v)
}
