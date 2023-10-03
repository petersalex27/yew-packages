package expr

import (
	"strings"

	"github.com/petersalex27/yew-packages/nameable"
)

type Case[T nameable.Nameable] struct {
	binders []Variable[T]
	when Expression[T]
	then Expression[T]
}

type PartialCase_when[T nameable.Nameable] Case[T]

func (bs BindersOnly[T]) InCase(when Expression[T], then Expression[T]) Case[T] {
	out := Case[T]{ binders: bs, }
	out.when = when.Bind(bs)
	out.then = then.Bind(bs)
	return out
}

func (bs BindersOnly[T]) When(e Expression[T]) PartialCase_when[T] {
	return PartialCase_when[T]{
		binders: bs,
		when: e.Bind(bs),
	}
}

func (pcw PartialCase_when[T]) Then(e Expression[T]) Case[T] {
	pcw.then = e.Bind(pcw.binders)
	return Case[T](pcw)
}

func (c Case[T]) String() string {
	return c.when.String() + " -> " + c.then.String()
}

func (c Case[T]) StrictString() string {
	strs := make([]string, len(c.binders))
	for i, v := range c.binders {
		strs[i] = v.StrictString()
	}
	return "Λ" + strings.Join(strs, " ") + " . " + c.when.StrictString() + " -> " + c.then.StrictString()
}

func (c Case[T]) Equals(cxt *Context[T], k Case[T]) bool {
	if len(c.binders) != len(k.binders) {
		return false
	}

	ok := c.when.Equals(cxt, k.when) && c.then.Equals(cxt, k.then)
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

func (c Case[T]) StrictEquals(k Case[T]) bool {
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

func selectionsMap[T nameable.Nameable](selections []Case[T], f func(Expression[T]) (Expression[T], bool)) ([]Case[T], bool) {
	out := make([]Case[T], len(selections))
	for i, c := range selections {
		var when, then Expression[T]
		var ok bool
		when, ok = f(c.when)
		if !ok {
			return selections, false
		}
		then, ok = f(c.then)
		if !ok {
			return selections, false
		}
		out[i] = Case[T]{when: when, then: then}
	}
	return out, true
}

func (c Case[T]) Find(v Variable[T]) bool {
	return c.when.Find(v) || c.then.Find(v)
}
