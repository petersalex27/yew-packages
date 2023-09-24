package expr

import "strconv"

type BoundVariable uint

func (v BoundVariable) String() string {
	return "a" + strconv.FormatInt(int64(v), 10)
}

func (v BoundVariable) Increment() BoundVariable {
	return v + 1
}

func (v BoundVariable) Decrement() BoundVariable {
	return v - 1
}