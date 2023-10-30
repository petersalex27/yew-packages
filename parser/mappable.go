package parser

import "github.com/petersalex27/yew-packages/parser/ast"

type mappable interface {
	Shift() ReductionRuleClass
	Then(productionOrder) ReductionRuleClass
}

func class(m mappable) (class ClassMapper, isClass bool) {
	class, isClass = m.(ClassMapper)
	return
}

func justMapper(m mappable) (mapper Mapper, isMapper bool) {
	mapper, isMapper = m.(Mapper)
	return
}

type Mapper []ast.Type

type ClassMapper struct {
	rep  Mapper
	mems []members
}

/*
Mapper functions
*/

// creates a lookahead to be mapped to a rule
func LookAhead(tys ...ast.Type) Mapper { return Mapper(tys) }

func (tys Mapper) Shift() ReductionRuleClass {
	return ReductionRuleClass{tys, shiftRuleSet}
}

func (tys Mapper) Then(rs productionOrder) ReductionRuleClass {
	return ReductionRuleClass{tys, rs}
}

func (tys Mapper) IsClass() (isClass bool, class ClassMapper) {
	isClass = false
	return
}

// == ClassMapper functions ===================================================

type members []ast.Type

// like calling LA multiple times to use with the same rule
func (m Mapper) For(tys ...ast.Type) ClassMapper {
	return m.ForN(4, tys...) // 4 is just an estimate
}

func (m Mapper) ForN(n uint, tys ...ast.Type) ClassMapper {
	mems := make([]members, 0, n)
	return ClassMapper{m, append(mems, members(tys))}
}

func (c ClassMapper) Or(tys ...ast.Type) ClassMapper {
	return ClassMapper{
		rep:  c.rep,
		mems: append(c.mems, members(tys)),
	}
}

func (tys ClassMapper) Shift() ReductionRuleClass {
	return ReductionRuleClass{tys, shiftRuleSet}
}

func (tys ClassMapper) Then(rs productionOrder) ReductionRuleClass {
	return ReductionRuleClass{tys, rs}
}

func (tys ClassMapper) IsClass() (isClass bool, class ClassMapper) {
	isClass, class = true, tys
	return
}
