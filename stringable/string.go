package stringable

import "strings"

type String string

func (s String) String() string {
	return string(s)
}

func (String) FromString(s string) (any, bool) {
	return String(s), true
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