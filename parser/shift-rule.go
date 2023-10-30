package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
)

type shiftRule pattern

func (r shiftRule) getPattern() PatternInterface { return pattern(r) }

func (shiftRule) call(_ *parser, _ ...ast.Ast) status.Status {
	return status.DoShift // this will trigger shift
}

func (r shiftRule) String() string {
	return fmt.Sprintf("rule(%v -> shift)", pattern(r))
}

func Shift() (s shiftRule) { return }

func (shiftRule) When(tys ...ast.Type) productionInterface {
	return shiftRule(tys)
}
