// =============================================================================
// Author-Date: Alex Peters - 2023
//
// Content: 
// contains common functional language functions, e.g., fmap, fold(l/r), and zip
//
// Notes: -
// =============================================================================
package fun

// takes a slice and a function that operates on the elements of that slice
// and returns an element of a pos. different type, then creates an array from
// that function's output
func FMap[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, len(xs))
	for i, x := range xs {
		out[i] = f(x)
	}
	return out
}

func FMapFilter[T, U any](xs []T, f func(T) (U, bool)) []U {
	out := make([]U, 0, len(xs))
	for _, x := range xs {
		res, add := f(x)
		if add {
			out = append(out, res)
		}
	}
	return out
}

// returns base if len(xs) == 0
func FoldLeft[T any](base T, xs []T, f func(T, T) T) T {
	for _, x := range xs {
		base = f(base, x)
	}
	return base
}

// returns base if len(xs) == 0
func FoldRight[T any](base T, xs []T, f func(T, T) T) T {
	for i := len(xs) - 1; i >= 0; i-- {
		base = f(xs[i], base)
	}
	return base
}

// calls a function w/ an element at same index in xs and ys for each index of
// the input and then returns the result of each function call in a
// slice; i.e., returns
//		[f(xs[0], ys[0]), f(xs[1], ys[1]), .., f(xs[len(xs)], ys[len(xs)])]
//
// NOTE: does not check length of xs or ys
func UncheckedZipWith[T, U, V any](f func(T, U) V, xs []T, ys []U) []V {
	vs := make([]V, len(xs))
	for i, x := range xs {
		vs[i] = f(x, ys[i])
	}
	return vs
}

// calls a function w/ an element at same index in xs and ys for each index of
// the input (see below) and then returns the result of each function call in a
// slice; i.e., returns
//		[f(xs[0], ys[0]), f(xs[1], ys[1]), .., f(xs[len(xs)], ys[len(xs)])]
//
// panics if lengths are not equal w/ "slice lengths must be equal"
func ZipWith[T, U, V any](f func(T, U) V, xs []T, ys []U) []V {
	if len(xs) != len(ys) {
		panic("slice lengths must be equal")
	}
	return UncheckedZipWith(f, xs, ys)
}

func UncheckedZip[T, U any](xs []T, ys []U) []struct{Left T; Right U} {
	out := make([]struct{Left T; Right U}, len(xs))
	for i, x := range xs {
		out[i].Left, out[i].Right = x, ys[i]
	}
	return out
}

func Zip[T, U any](xs []T, ys []U) []struct{Left T; Right U} {
	if len(xs) != len(ys) {
		panic("length of input arrays do not match.")
	}
	return UncheckedZip(xs, ys)
}

func UncheckedAndZip[T, U any](base bool, xs []T, ys []U, f func(T, U) bool) bool {
	if !base {
		return false
	}

	for i := range xs {
		if !f(xs[i], ys[i]) {
			return false
		}
	}
	return true
}

func UncheckedOrZip[T, U any](base bool, xs []T, ys []U, f func(T, U) bool) bool {
	if base {
		return true
	}

	for i := range xs {
		if f(xs[i], ys[i]) {
			return true
		}
	}
	return false
}

func AndZip[T, U any](base bool, xs []T, ys []U, f func(T, U) bool) bool {
	if len(xs) != len(ys) {
		return false
	}
	return UncheckedAndZip(base, xs, ys, f)
}

func OrZip[T, U any](base bool, xs []T, ys []U, f func(T, U) bool) bool {
	if len(xs) != len(ys) {
		return false
	}
	return UncheckedOrZip(base, xs, ys, f)
}