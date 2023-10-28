package table

import "github.com/petersalex27/yew-packages/nameable"

// element in Table type
// 
// see Table
type tableElement[T any] struct {
	key nameable.Nameable
	val T
}

// Table that maps instances of nameable.Nameable to a value of type T
type Table[T any] struct {
	data map[string]tableElement[T]
}

// number of elements in table
func (table *Table[T]) Len() int {
	return len(table.data)
}

// (Over)writes `val` at domain `key`
func (table *Table[T]) Add(key nameable.Nameable, val T) {
	table.data[key.GetName()] = tableElement[T]{key, val}
}

// If `key` is not found in the table, then `_, false` is returned, else the
// value mapped to by `key` is returned and true is returned
func (table *Table[T]) Get(key nameable.Nameable) (val T, ok bool) {
	var tmp tableElement[T]
	tmp, ok = table.data[key.GetName()]
	if ok {
		val = tmp.val
	}
	return
}

// Removes key-value pair from table if `key` is in the table, returning the 
// removed value. Otherwise `_, false` is returned.
func (table *Table[T]) Remove(key nameable.Nameable) (val T, ok bool) {
	if val, ok = table.Get(key); ok {
		delete(table.data, key.GetName())
	}
	return
}

// return underlying data used for table
func (table *Table[T]) GetRawMap() map[string]tableElement[T] { return table.data }