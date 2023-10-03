package expr

import (
	"strconv"

	"github.com/petersalex27/yew-packages/nameable"
)

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

type Prim[T nameable.Nameable] struct {
	primIFace
}

func (p Prim[T]) String() string {
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

func (p Prim[T]) Equals(_ *Context[T], e Expression[T]) bool {
	p2, ok := e.(Prim[T])
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

func (p Prim[T]) StrictString() string {
	return p.String()
}

func (p Prim[T]) StrictEquals(e Expression[T]) bool {
	return p.Equals(nil, e)
}

func (p Prim[T]) Replace(Variable[T], Expression[T]) (Expression[T], bool) {
	return p, false
}

func (p Prim[T]) UpdateVars(gt int, by int) Expression[T] { return p }

func (p Prim[T]) Again() (Expression[T], bool) { return p, false}

func (p Prim[T]) Bind(BindersOnly[T]) Expression[T] { return p }

func (p Prim[T]) Find(Variable[T]) bool { return false }

func (p Prim[T]) PrepareAsRHS() Expression[T] { return p }

func (p Prim[T]) Rebind() Expression[T] { return p }

func (p Prim[T]) Copy() Expression[T] { return p }

func (p Prim[T]) ForceRequest() Expression[T] { return p }