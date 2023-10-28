package expr

import (
	"strconv"

	"github.com/petersalex27/yew-packages/nameable"
)

// prim int lit
type Int int64

// intPrim
func (Int) getPrimType() primType { return intPrim }

func MakeInt[T nameable.Nameable](i int64) Prim[T] {
	return Prim[T]{Int(i)}
}

// prim uint lit
type Uint uint64

// uintPrim
func (Uint) getPrimType() primType { return uintPrim }

func MakeUint[T nameable.Nameable](u uint64) Prim[T] {
	return Prim[T]{Uint(u)}
}

// prim float lit
type Float float64

// floatPrim
func (Float) getPrimType() primType { return floatPrim }

func MakeFloat[T nameable.Nameable](f float64) Prim[T] {
	return Prim[T]{Float(f)}
}

// prim char lit
type Char byte

// charPrim
func (Char) getPrimType() primType { return charPrim }

func MakeChar[T nameable.Nameable](c byte) Prim[T] {
	return Prim[T]{Char(c)}
}

// prim string lit
type Str string

// strPrim
func (Str) getPrimType() primType { return strPrim }

func MakeString[T nameable.Nameable](s string) Prim[T] {
	return Prim[T]{Str(s)}
}

type ConversionStatus int

const (
	// conversion success
	Ok ConversionStatus = iota
	// conversion failed for other reason
	BadConversion
	// int exceeds upper/lower bound of int64
	BadIntRange
	// tried to convert to an integer but not an integer
	BadIntSyntax
	// uint exceeds upper/lower bound of uint64
	BadUintRange
	// tried to convert to an unsigend integer but not an unsigned integer
	BadUintSyntax
	// float exceeds upper/lower bound of float64
	BadFloatRange
	// tried to convert to a float but not a float
	BadFloatSyntax
	// too many characters in char literal
	BadCharRange
	// invalid characters in char literal
	BadCharSyntax
	// just here for completeness
	BadStringRange; BadStringSyntax
)

func setStat(e error, firstBadStat ConversionStatus) ConversionStatus {
	if e == nil {
		return Ok
	}

	switch firstBadStat {
	case Ok:
		if e != nil {
			return BadConversion
		}
	case BadConversion:
		return BadConversion
	case BadIntRange:
		if e.(*strconv.NumError).Err == strconv.ErrSyntax {
			return BadIntSyntax
		}
	case BadIntSyntax:
		if e.(*strconv.NumError).Err != strconv.ErrSyntax {
			return BadIntRange
		}
	case BadUintRange:
		if e.(*strconv.NumError).Err == strconv.ErrSyntax {
			return BadUintSyntax
		}
	case BadUintSyntax:
		if e.(*strconv.NumError).Err != strconv.ErrSyntax {
			return BadUintRange
		}
	case BadFloatRange:
		if e.(*strconv.NumError).Err == strconv.ErrSyntax {
			return BadFloatSyntax
		}
	case BadFloatSyntax:
		if e.(*strconv.NumError).Err != strconv.ErrSyntax {
			return BadFloatRange
		}
	case BadCharRange:
		return BadCharRange
	case BadCharSyntax:
		return BadCharSyntax
	case BadStringRange:
		return BadStringRange
	case BadStringSyntax:
		return BadStringSyntax
	default:
		return BadConversion
	}

	return firstBadStat
}

// Make prim value from `s` based on `format`.
//
// The valid formats are
//  - 'i', Int
//  - 'u', Uint
//  - 'f', Float
//  - 'c', Char
//  - 's', String
//
// if the format is anything else a status of BadConversion is returned
//
// - conversion for Int, Uint, and Float can return either a status of `Ok` or 
// 		one of `BadXxxRange` or `BadXxxSyntax` for the respective type
// - conversion to Char can only return `Ok` or, when `len(s) > 1`, `BadCharSyntax`
// - conversion to String always returns `Ok`
func Makef[T nameable.Nameable](format byte, s string) (primitive Prim[T], stat ConversionStatus) {
	var e error = nil
	switch format {
	case 'i':
		var i int64
		i, e = strconv.ParseInt(s, 0, 64)
		primitive, stat = MakeInt[T](i), setStat(e, BadIntRange)
	case 'u':
		var u uint64
		u, e = strconv.ParseUint(s, 0, 64)
		primitive, stat = MakeUint[T](u), setStat(e, BadUintRange)
	case 'f':
		var f float64
		f, e = strconv.ParseFloat(s, 64)
		primitive, stat = MakeFloat[T](f), setStat(e, BadFloatRange)
	case 'c':
		var c byte
		if len(s) != 1 {
			stat = BadCharSyntax
		}
		c = s[0]
		primitive, stat = MakeChar[T](c), Ok
	case 's':
		primitive, stat = MakeString[T](s), Ok
	default:
		// bad conversion format
		stat = BadConversion
	}

	return
}

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

func (Prim[T]) ExtractFreeVariables(dummyVar Variable[T]) []Variable[T] {
	return []Variable[T]{}
}

func (Prim[T]) Collect() []T {
	return []T{}
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