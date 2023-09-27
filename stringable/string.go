package stringable

type String string

func (s String) String() string {
	return string(s)
}

func (String) FromString(s string) (any, bool) {
	return String(s), true
}