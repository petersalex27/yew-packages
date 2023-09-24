package util

func CopyMap[K comparable, V any](old map[K]V) map[K]V {
	new := make(map[K]V, len(old))
	for k, v := range old {
		new[k] = v
	}
	return new
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func PowerOfTwoCeil(n uint32) uint32 {
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
