package str

import "strconv"

type Int int

func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (Int) FromString(s string) (any, bool) {
	out, e := strconv.ParseInt(s, 0, 64)
	return Int(out), e == nil
}