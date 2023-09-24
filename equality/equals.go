package equality

type Eq[T any] interface {
	Equals(T) bool
}