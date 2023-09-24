package types

// Family is a collection of proper types
type Family[T Type] []T

// IndexedFamily is a collection of indexed proper types, duh
type IndexedFamily[I any, T Type] []struct{ elem T; index I }