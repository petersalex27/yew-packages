package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/errors"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
)

type ErrorFn func(top ast.Ast, nodes ...ast.Ast) errors.Err

type errorRule struct {
	pattern
	ErrorFn
}

func (r errorRule) String() string {
	return fmt.Sprintf("rule(%v -> error)", r.pattern)
}

func (r errorRule) getPattern() PatternInterface { return r.pattern }

func (r errorRule) call(p *parser, nodes ...ast.Ast) (stat status.Status, ruleApplied bool) {
	node, _ := p.top()
	p.errors = append(p.errors, r.ErrorFn(node, nodes...))
	return status.Error, true
}

type errorNeedsPattern errorRule

func Error(e func(top ast.Ast, nodes ...ast.Ast) errors.Err) errorNeedsPattern {
	return errorNeedsPattern{ErrorFn: e}
}

func (e errorNeedsPattern) From(tys ...ast.Type) productionInterface {
	return errorRule{
		pattern: pattern(tys),
		ErrorFn: e.ErrorFn,
	}
}
