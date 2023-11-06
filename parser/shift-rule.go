package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
)

type shiftRule pattern

func (r shiftRule) getPattern() PatternInterface { return pattern(r) }

func (shiftRule) call(*parser, ...ast.Ast) (stat status.Status, ruleApplied bool) {
	return status.DoShift, true // this will trigger shift
}

func (r shiftRule) String() string {
	return fmt.Sprintf("rule(%v -> shift)", pattern(r))
}

func Shift() (s shiftRule) { return }

func (shiftRule) When(tys ...ast.Type) productionInterface {
	return shiftRule(tys)
}
