package parser

import "github.com/petersalex27/yew-packages/parser/ast"

type PatternInterface interface {
	Match(nodes ...ast.Ast) bool
	MaxHandleLength() int
	Last() (last ast.Type, emptyPattern bool)
}

// R.H.S. of production rule
type pattern []ast.Type

// returns true iff the pattern `pat` matches the given nodes
func (pat pattern) Match(nodes ...ast.Ast) bool {
	if len(pat) != len(nodes) {
		return false
	}

	for i, node := range nodes {
		if pat[i] != node.NodeType() {
			return false
		}
	}
	return true
}

// (pattern).Last returns the last element in the pattern represented by the
// receiver along with `true`. If the represented pattern is empty, 
// ast.Type(0) along with `false` is returned
func (pat pattern) Last() (last ast.Type, isEmpty bool) {
	patternLength := len(pat)
	if isEmpty = patternLength == 0; !isEmpty {
		last = pat[patternLength-1]
	}
	return
}

// maximum handle length pattern can match with
func (pat pattern) MaxHandleLength() int {
	return len(pat)
}