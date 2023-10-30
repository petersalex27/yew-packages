package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
)

type rule struct {
	pattern
	Production
}

func (r rule) String() string {
	return fmt.Sprintf("rule(%v -> %v)", r.pattern, r.Productions)
}

func (r rule) getPattern() PatternInterface { return r.pattern }

func (r rule) call(p *parser, nodes ...ast.Ast) status.Status {
	return r.Production.call(p, uint(len(nodes)), nodes...)
}

func (p needPattern) From(tys ...ast.Type) productionInterface {
	return rule{pattern(tys), Production{p}}
}
