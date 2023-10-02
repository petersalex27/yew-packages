package types

import "github.com/petersalex27/yew-packages/nameable"

// Family is a collection of proper types
type Family[T Type[U], U nameable.Nameable] []T

// IndexedFamily is a collection of indexed proper types, duh
type IndexedFamily[I any, T Type[U], U nameable.Nameable] []struct{ elem T; index I }