package str

type Bool bool

func (b Bool) String() string {
	if b {
		return "True"
	} 
	return "False"
}

func (Bool) FromString(s string) (any, bool) {
	if s == "True" {
		return Bool(true), true
	} else if s == "False" {
		return Bool(false), true
	}
	return Bool(false), false
}