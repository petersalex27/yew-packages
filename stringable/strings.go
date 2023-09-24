package stringable

type Stringable interface {
	String() string
}

type Parseable interface {
	Stringable
	FromString(string) (any, bool)
}