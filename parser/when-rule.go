package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
)

type whenRule struct {
	pattern
	clear uint
	Production
}

func (r whenRule) String() string {
	return fmt.Sprintf("when(%v / %d -> %v)", r.pattern, r.clear, r.Productions)
}

func (r whenRule) getPattern() PatternInterface { return r.pattern }

func (r whenRule) call(p *parser, nodes ...ast.Ast) status.Status {
	return r.Production.call(p, r.clear, nodes[uint(len(nodes))-r.clear:]...)
}

// intermediate struct
type whenNeedPattern struct {
	needPattern
	when []ast.Type
}

func (f NamedProduction) When(tys ...ast.Type) whenNeedPattern {
	return needPattern(Production{f}).When(tys...)
}

func (p needPattern) When(tys ...ast.Type) whenNeedPattern {
	return whenNeedPattern{p, tys}
}

func (p whenNeedPattern) From(tys ...ast.Type) productionInterface {
	clear := len(tys)
	pat := make(pattern, len(p.when)+clear)
	copy(pat, p.when)
	copy(pat[len(p.when):], tys)

	return whenRule{
		pattern(pat),
		uint(clear),
		Production{p.needPattern},
	}
}
