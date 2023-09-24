package str

import "strconv"

type Float float64

func (f Float) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

func (Float) FromString(s string) (any, bool) {
	out, e := strconv.ParseFloat(s, 64)
	return Float(out), e == nil
}