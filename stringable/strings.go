package stringable

import "strings"

type Stringable interface {
	String() string
}

type Parseable interface {
	Stringable
	FromString(string) (any, bool)
}

func Join[S Stringable](elems []S, sep Stringable) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0].String()
	}

	var b strings.Builder
	b.WriteString(elems[0].String())
	sep_ := sep.String()
	for _, s := range elems[1:] {
		b.WriteString(sep_)
		b.WriteString(s.String())
	}
	return b.String()
}