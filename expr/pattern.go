package expr

import "github.com/petersalex27/yew-packages/nameable"

// things that can be converted to `Pattern`s
type Patternable[N nameable.Nameable] interface {
	// attempts to convert receiver into a pattern. The second return value is
	// false when receiver (or some child of receiver) can't be converted.
	// Otherwise, the second return value is true and a pattern is returned
	ToAlmostPattern() (AlmostPattern[N], bool)
}

type AlmostPattern[N nameable.Nameable] struct{ pattern matchable[N] }

func (p AlmostPattern[N]) Match(p2 AlmostPattern[N]) bool {
	return p.pattern.match(p2.pattern)
}

// returns wrapped pattern
func (pat AlmostPattern[N]) GetPattern() matchable[N] { return pat.pattern }

// things that can be matched
type matchable[N nameable.Nameable] interface {
	match(matchable[N]) bool
}

// sequence of patterns
type PatternSequence[N nameable.Nameable] struct {
	ty       patternSequenceType
	sequence []matchable[N]
}

// match a sequence to another sequence
func (seq PatternSequence[N]) match(m matchable[N]) bool {
	if seq.ty == PatternSequenceWildcard { // this is possible only at the top level
		return true
	}

	// also a sequence?
	seq2, ok := m.(PatternSequence[N])
	if !ok {
		return false
	}

	// same type of sequence?
	if seq2.ty != seq.ty {
		return false
	}

	// same sequence length?
	if len(seq.sequence) != len(seq2.sequence) {
		return false
	}

	// same sequence element pattern?
	for i, e := range seq.sequence {
		if !e.match(seq2.sequence[i]) {
			return false
		}
	}

	return true
}

// types of basic pattern elements
type patternElementType byte

const (
	// represents any element
	PatternElementWildcard patternElementType = iota
	// represents a variable
	PatternElementVar
	// represents a constant
	PatternElementConst
	// represents a literal
	PatternElementLiteral
)

// types of sequences of patterns
type patternSequenceType byte

const (
	// represents any sequence
	PatternSequenceWildcard patternSequenceType = iota
	// represents a list, e.g., [1, 2, a]
	PatternSequenceList
	// represents a tuple, e.g., (a, 1)
	PatternSequenceTuple
	// represents an application, e.g., (A x [a, a])
	PatternSequenceApplication
)

// most basic pattern
type PatternElement[N nameable.Nameable] struct {
	element N
	ty      patternElementType
}

// match an element to another element
func (elem PatternElement[N]) match(m matchable[N]) bool {
	if elem.ty == PatternElementWildcard {
		return true
	}

	elem2, ok := m.(PatternElement[N])
	if !ok {
		return false
	}

	if elem.ty != elem2.ty {
		return false
	}

	return elem.element.GetName() == elem2.element.GetName()
}

// creates a pattern element--these are the most basic parts of a pattern
func MakeElem[N nameable.Nameable](ty patternElementType, element N) PatternElement[N] {
	return PatternElement[N]{element, ty}
}

// creates a sequence of patterns
func MakeSequence[N nameable.Nameable](ty patternSequenceType, elems ...matchable[N]) PatternSequence[N] {
	return PatternSequence[N]{ty, elems}
}

func (p PatternElement[N]) ToAlmostPattern() (AlmostPattern[N], bool) {
	return AlmostPattern[N]{pattern: p}, true
}

func (p PatternSequence[N]) ToAlmostPattern() (AlmostPattern[N], bool) {
	return AlmostPattern[N]{pattern: p}, true
}
