package iterator

type iterator[T any] struct {
	index int
	end int
	elems []T
}

func (it *iterator[T]) Next() (elem T, exists bool) {
	exists = it.index < it.end
	if exists {
		elem = it.elems[it.index]
		it.index++
	}
	return
}

func (it *iterator[T]) HasNext() bool {
	return it.index < it.end
}

func Iterator[T any](ts []T) (it *iterator[T]) {
	it = new(iterator[T])
	*it = iterator[T]{
		index: 0, end: len(ts),
		elems: ts,
	}
	return
}