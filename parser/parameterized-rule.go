package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
)

// ============================================================================
// parameterized rule struct
// ============================================================================

// calls production rule function only if `onlyIf` returns true, else call will
// return status Ok allowing rules to be continued to be searched.
//
// NOTE: This will pass all handle elements to `onlyIf`, including "when" elems
type parameterizedRule struct {
	production productionInterface
	onlyIf     func(nodes ...ast.Ast) bool
}

func (r parameterizedRule) getPattern() PatternInterface {
	return r.production.getPattern()
}

func (r parameterizedRule) String() string {
	return "parameterized " + r.production.String()
}

// calls production rule function only if `onlyIf` returns true, else call will
// return status Ok allowing rules to be continued to be searched.
//
// NOTE: This will pass all handle elements to `onlyIf`, including "when" elems
func (r parameterizedRule) call(p *parser, nodes ...ast.Ast) (stat status.Status, ruleApplied bool) {
	if r.onlyIf == nil || r.onlyIf(nodes...) {
		return r.production.call(p, nodes...)
	}
	return status.Ok, false
}

// gives a rule an additional condition for being called
func Precondition(pi productionInterface, condition func(nodes ...ast.Ast) bool) parameterizedRule {
	return parameterizedRule{pi, condition}
}
