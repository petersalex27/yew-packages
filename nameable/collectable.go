package nameable

type Collectable[T Nameable] interface {
	Collect() []T
}