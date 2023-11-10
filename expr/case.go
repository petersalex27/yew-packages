package expr

import (
	"strings"

	"github.com/petersalex27/yew-packages/nameable"
)

type Case[T nameable.Nameable] struct {
	binders    []Variable[T]
	pattern    Expression[T]
	expression Expression[T]
}

func (c Case[T]) BodyAbstract(v Variable[T], name Const[T]) Case[T] {
	return Case[T]{
		c.binders,
		c.pattern.BodyAbstract(v, name),
		c.expression.BodyAbstract(v, name),
	}
}

func (c Case[T]) ExtractFreeVariables(dummyVar Variable[T]) []Variable[T] {
	var when, then Expression[T] = c.pattern, c.expression
	for _, v := range c.binders {
		when, _ = when.Replace(v, dummyVar)
		then, _ = then.Replace(v, dummyVar)
	}

	return append(c.pattern.ExtractFreeVariables(dummyVar), c.expression.ExtractFreeVariables(dummyVar)...)
}

func (a Case[T]) Collect() []T {
	res := make([]T, 0, len(a.binders))
	for _, binder := range a.binders {
		res = append(res, binder.Collect()...)
	}
	res = append(res, a.pattern.Collect()...)
	res = append(res, a.expression.Collect()...)
	return res
}

type PartialCase_when[T nameable.Nameable] Case[T]

func (bs BindersOnly[T]) InCase(when Expression[T], then Expression[T]) Case[T] {
	out := Case[T]{binders: bs}
	out.pattern = when.Bind(bs)
	out.expression = then.Bind(bs)
	return out
}

func (bs BindersOnly[T]) When(e Expression[T]) PartialCase_when[T] {
	return PartialCase_when[T]{
		binders: bs,
		pattern: e.Bind(bs),
	}
}

func (pcw PartialCase_when[T]) Then(e Expression[T]) Case[T] {
	pcw.expression = e.Bind(pcw.binders)
	return Case[T](pcw)
}

func (c Case[T]) String() string {
	return groupStringed(c.pattern.String() + onMatchString() + c.expression.String())
}

func (c Case[T]) StrictString() string {
	strs := make([]string, len(c.binders))
	for i, v := range c.binders {
		strs[i] = v.StrictString()
	}

	hiddenBinders := ""
	if len(c.binders) != 0 {
		hiddenBinders = "Λ" + strings.Join(strs, " ") + " . "
	}

	return groupStringed(hiddenBinders + c.pattern.StrictString() + onMatchString() + c.expression.StrictString())
}

func (c Case[T]) Equals(cxt *Context[T], k Case[T]) bool {
	if len(c.binders) != len(k.binders) {
		return false
	}

	ok := c.pattern.Equals(cxt, k.pattern) && c.expression.Equals(cxt, k.expression)
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

	ok := c.pattern.StrictEquals(k.pattern) && c.expression.StrictEquals(k.expression)
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
		when, ok = f(c.pattern)
		if !ok {
			return selections, false
		}
		then, ok = f(c.expression)
		if !ok {
			return selections, false
		}
		out[i] = Case[T]{pattern: when, expression: then}
	}
	return out, true
}

func (c Case[T]) Find(v Variable[T]) bool {
	return c.pattern.Find(v) || c.expression.Find(v)
}
