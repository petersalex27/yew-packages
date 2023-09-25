package expr

import "strconv"

type Int int64
func (Int) getPrimType() primType { return intPrim }

type Uint uint64
func (Uint) getPrimType() primType { return uintPrim }

type Float float64
func (Float) getPrimType() primType { return floatPrim }

type Char byte
func (Char) getPrimType() primType { return charPrim }

type Str string
func (Str) getPrimType() primType { return strPrim }

type primType uint8
const (
	intPrim primType = iota
	uintPrim
	floatPrim
	charPrim
	strPrim
)

type primIFace interface {
	getPrimType() primType
}

type Prim struct {
	primIFace
}

func (p Prim) String() string {
	switch p.getPrimType() {
	case intPrim:
		return strconv.FormatInt(int64(p.primIFace.(Int)), 10)
	case uintPrim:
		return strconv.FormatInt(int64(p.primIFace.(Uint)), 10)
	case floatPrim:
		return strconv.FormatFloat(float64(p.primIFace.(Float)), 'f', -1, 64)
	case charPrim:
		return string(p.primIFace.(Char))
	case strPrim:
		return string(p.primIFace.(Str))
	}
	panic("unknown prim.")
}

func (p Prim) Equals(e Expression) bool {
	p2, ok := e.(Prim)
	if !ok {
		return false
	}
	if p.getPrimType() != p2.getPrimType() {
		return false 
	}
	switch p.getPrimType() {
	case intPrim:
		return p.primIFace.(Int) == p2.primIFace.(Int)
	case uintPrim:
		return p.primIFace.(Uint) == p2.primIFace.(Uint)
	case floatPrim:
		return p.primIFace.(Float) == p2.primIFace.(Float)
	case charPrim:
		return p.primIFace.(Char) == p2.primIFace.(Char)
	case strPrim:
		return p.primIFace.(Str) == p2.primIFace.(Str)
	default:
		return false
	}
}

func (p Prim) StrictString() string {
	return p.String()
}

func (p Prim) StrictEquals(e Expression) bool {
	return p.Equals(e)
}

func (p Prim) Replace(Variable, Expression) (Expression, bool) {
	return p, false
}

func (p Prim) UpdateVars(gt int, by int) Expression { return p }

func (p Prim) Again() (Expression, bool) { return p, false}

func (p Prim) Bind(BindersOnly) Expression { return p }

func (p Prim) Find(Variable) bool { return false }

func (p Prim) PrepareAsRHS() Expression { return p }

func (p Prim) Rebind() Expression { return p }

func (p Prim) Copy() Expression { return p }

func (p Prim) ForceRequest() Expression { return p }