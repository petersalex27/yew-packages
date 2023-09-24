package fun

func FMap[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, len(xs))
	for i, x := range xs {
		out[i] = f(x)
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


func UncheckedZip[T, U, V any](base V, xs []T, ys []U, f func(T, U, V) V) V {
	for i := range xs {
		base = f(xs[i], ys[i], base)
	}
	return base
}

func Zip[T, U, V any](base V, xs []T, ys []U, f func(T, U, V) V) V {
	if len(xs) != len(ys) {
		panic("length of input arrays do not match.")
	}
	return UncheckedZip(base, xs, ys, f)
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