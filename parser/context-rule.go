package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
)

type contextRule struct {
	pattern
	signalReduction status.Status
	clear           bool
	whenLen         int
	// should return status.Ok, else interpreted as fail
	action func(cxt *ParserContext, nodes ...ast.Ast) status.Status
}

func (r contextRule) String() string {
	sig := r.signalReduction.String()
	clr := "leave"
	if r.clear {
		clr = "clear"
	}

	return fmt.Sprintf("context[%s](%v -> action -> %s?reduction)", clr, r.pattern, sig)
}

func (r contextRule) getPattern() PatternInterface { return r.pattern }

func (r contextRule) call(p *parser, nodes ...ast.Ast) status.Status {
	stat := r.action(p.cxt, nodes...)
	if r.clear {
		clearLen := uint(len(nodes)) - uint(r.whenLen)
		p.stack.Clear(clearLen)
	}

	if stat.Is(r.signalReduction) && p.cxt.reduction != nil {
		reduced := p.cxt.reduction.do(nodes[r.whenLen:]...)
		// save result?
		p.maybePush(reduced)
	}
	return stat
}

type contextFunction func(*ParserContext, ...ast.Ast) status.Status

type contextualizeStep1 struct {
	f contextFunction
}

type contextualizeStep2 struct {
	f      contextFunction
	signal status.Status
}

type cxtWhenNeedPattern struct {
	when []ast.Type
	contextualizeStep2
}

func Contextualize(f func(*ParserContext, ...ast.Ast) status.Status) contextualizeStep1 {
	return contextualizeStep1{f}
}

func (c contextualizeStep1) ReduceOnStatus(stat int) contextualizeStep2 {
	return c.ReduceOn(status.Status(stat))
}

func (c contextualizeStep1) ReduceOn(stat status.Status) contextualizeStep2 {
	return contextualizeStep2{f: c.f, signal: stat}
}

func (c contextualizeStep2) When(tys ...ast.Type) cxtWhenNeedPattern {
	return cxtWhenNeedPattern{contextualizeStep2: c, when: tys}
}

func (p contextualizeStep2) From(tys ...ast.Type) productionInterface {
	return contextRule{
		pattern:         pattern(tys),
		signalReduction: p.signal,
		clear:           true,
		whenLen:         0,
		action:          p.f,
	}
}

func (p cxtWhenNeedPattern) From(tys ...ast.Type) productionInterface {
	out := contextRule{
		signalReduction: p.signal,
		clear:           true,
		whenLen:         len(p.when),
		action:          p.f,
	}
	clear := len(tys)
	pat := make(pattern, len(p.when)+clear)
	copy(pat, p.when)
	copy(pat[len(p.when):], tys)
	out.pattern = pat
	return out
}
