package errors

import (
	str "github.com/petersalex27/yew-packages/stringable"
	"github.com/petersalex27/yew-packages/util"
)

type Pointable interface {
	str.Stringable
	setPadding(int) Pointable
	setTail(string) Pointable
	getShared() pointer_shared
}

type pointer_shared struct {
	tailMsg string
	paddingLeft int
}

func (p pointer_shared) Strings() (padding, tail string) {
	padLen := util.Max(p.paddingLeft, 0)
	pad := make([]byte, padLen)
	for i := range pad {
		pad[i] = ' '
	}
	return string(pad), p.tailMsg
}

type Pointer struct {
	pointer_shared
}

func (p Pointer) getShared() pointer_shared {
	return p.pointer_shared
}

func (p Pointer) setTail(s string) Pointable {
	p.tailMsg = s
	return p
}

func (p Pointer) setPadding(n int) Pointable {
	p.pointer_shared.paddingLeft = n
	return p
}

// "[   ..][^][msg]"
func (p Pointer) String() string {
	padding, msg := p.pointer_shared.Strings()
	return padding + "^" + msg
}

type PointerRange struct {
	rngLen int
	pointer_shared
}

func (p PointerRange) getShared() pointer_shared {
	return p.pointer_shared
}

func (p PointerRange) buildRange() string {
	if p.rngLen <= 0 {
		return ""
	}
	
	ptrs := make([]byte, p.rngLen)
	for i := range ptrs {
		ptrs[i] = '^'
	}
	return string(ptrs)
}

func (p PointerRange) setTail(s string) Pointable {
	p.tailMsg = s
	return p
}

func (p PointerRange) setPadding(n int) Pointable {
	p.pointer_shared.paddingLeft = n
	return p
}

// "[   ..][^^^..][msg]"
func (p PointerRange) String() string {
	padding, msg := p.pointer_shared.Strings()
	res := padding + p.buildRange() + msg
	return res
}