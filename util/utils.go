package util

// reverses the slice `arr` in-place.
//
// example: 
//	arr = []int{1,2,3,4}
//	ReverseInPlace(arr)
//	arr == []int{4,3,2,1}
func ReverseInPlace[T any](arr []T) {
	arrLen := len(arr)
	for i := 0; i < arrLen / 2; i++ {
		tmp := arr[i]
		arr[i] = arr[arrLen-i]
		arr[arrLen-i] = tmp
	}
}

// makes a perfect-fit copy of slice `arr`, and then reverses the elements of 
// the new slice copy
//
// see `ReverseInPlace[T any]([]T)` for an example
func Reverse[T any](arr []T) (reversed []T) {
	reversed = CopySlice(arr)
	ReverseInPlace(reversed)
	return
}

func CopySlice[T any](arr []T) (out []T) {
	out = make([]T, len(arr))
	copy(out, arr)
	return
}

func CopyMap[K comparable, V any](old map[K]V) map[K]V {
	new := make(map[K]V, len(old))
	for k, v := range old {
		new[k] = v
	}
	return new
}

// returns the absolute value of some integer type
func Abs[T ~int](a T) T {
	if a > 0 {
		return a
	}
	return -a
}

func Max[T ~int](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T ~int](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func UMax[T ~uint](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func UMin[T ~uint](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func PowerOfTwoCeil(n uint) uint {
	/*
		n = 01001101011100001100110110110110 // initial
		n = 01001101011100001100110110110101 // a
		n =
			01001101011100001100110110110101 |
			00100110101110000110011011011010
		  = 01101111111110001110111111111111 // b
		n = 01101111111110001110111111111111 |
			00011011111111100011101111111111
		  = 01111111111111101111111111111111 // c
		n = 01111111111111101111111111111111 |
			00000111111111111110111111111111
		  = 01111111111111111111111111111111 // d
		n = 01111111111111111111111111111111 |
			00000000011111111111111111111111
		  = 01111111111111111111111111111111 // e
		n = 01111111111111111111111111111111 |
			00000000000000000111111111111111
		  = 01111111111111111111111111111111 // f
		n = 10000000000000000000000000000000 // g
	*/
	/*
		n = 00000000000000000000000000100000 // initial
		n = 00000000000000000000000000011111 // a
		n =
			00000000000000000000000000011111 |
			00000000000000000000000000001111
		  = 00000000000000000000000000011111 // b
		n = 00000000000000000000000000011111 |
			00000000000000000000000000000111
		  = 00000000000000000000000000011111 // c
		n = 00000000000000000000000000011111 |
			00000000000000000000000000000001
		  = 00000000000000000000000000011111 // d
		n = 00000000000000000000000000011111 |
			00000000000000000000000000000000
		  = 00000000000000000000000000011111 // e
		n = 00000000000000000000000000011111 |
			00000000000000000000000000000000
		  = 00000000000000000000000000011111 // f
		n = 00000000000000000000000000100000 // g
	*/
	n = n - 1         // a
	n = n | (n >> 1)  // b
	n = n | (n >> 2)  // c
	n = n | (n >> 4)  // d
	n = n | (n >> 8)  // e
	n = n | (n >> 16) // f
	return n + 1      // g
}
