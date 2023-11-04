package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/errors"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
)

type WarnFn struct {
	Warn func(top ast.Ast, nodes ...ast.Ast) errors.Warning
	Productions
}

type warningRule struct {
	pattern
	WarnFn
}

func (r warningRule) String() string {
	return fmt.Sprintf("rule(%v -> (warning, %v))", r.pattern, r.Productions)
}

func (r warningRule) getPattern() PatternInterface { return r.pattern }

func (r warningRule) call(p *parser, nodes ...ast.Ast) status.Status {
	// look ahead token is for warning information; it should not be used in
	// reduction
	tok, _ := p.lookAhead()

	// pop stack (this removes `nodes`)
	n := uint(len(nodes))
	p.stack.Clear(n) // must be called before pushing reduction result

	// do warning and production action (production happens inside r.WarnFn)
	warning := r.WarnFn.Warn(ast.TokenNode(tok), nodes...)
	// add warning
	p.warnings = append(p.warnings, warning)
	// do production action
	product := r.do(nodes...)
	// save result?
	p.maybePush(product)
	return status.Ok
}

type warnNeedsReduction warningRule

type warnNeedPattern warningRule

func Warn(warnFn func(top ast.Ast, nodes ...ast.Ast) errors.Warning) warnNeedsReduction {
	return warnNeedsReduction(warningRule{WarnFn: WarnFn{Warn: warnFn}})
}

func (w warnNeedsReduction) ThenGet(f func(nodes ...ast.Ast) ast.Ast) warnNeedPattern {
	return warnNeedPattern{WarnFn: WarnFn{Warn: w.Warn, Productions: ProductionFunction(f)}}
}

func (w warnNeedPattern) From(tys ...ast.Type) productionInterface {
	return warningRule{
		pattern: pattern(tys),
		WarnFn:  w.WarnFn,
	}
}
