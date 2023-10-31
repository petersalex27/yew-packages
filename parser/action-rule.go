package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
)

// production rule that allows registered functions to be called
type actionRule struct {
	pattern
	ProductionWith
}

func (ar actionRule) String() string {
	return fmt.Sprintf("actionRule(%v -> %v)", ar.pattern, ar.ProductionWith)
}

func (ar actionRule) getPattern() PatternInterface { return ar.pattern }

func (ar actionRule) call(p *parser, handle ...ast.Ast) status.Status {
	return ar.ProductionWith.call(p, uint(len(handle)), handle...)
}

type actionRuleNeedsReduceTo struct{}

type actionRuleNeedsPattern struct{ ProductionWith }

// declares an action production rule
func ActionRule() actionRuleNeedsReduceTo { return actionRuleNeedsReduceTo{} }

func (w actionRuleNeedsReduceTo) Get(
	production func(
		action func(name string) func(any),
		handle ...ast.Ast,
	) ast.Ast,
) actionRuleNeedsPattern {
	return actionRuleNeedsPattern{production}
}

// adds handle to in-progress production rule then returns it
func (w actionRuleNeedsPattern) From(ty ...ast.Type) productionInterface {
	return actionRule{
		pattern(ty),
		w.ProductionWith,
	}
}

// A production rule function with an action included. The intention of the
// attached action is to allow for users of this package to have their own
// parsing state inside the action function; the call to action would then
// change the state, print errors, whatever
type ProductionWith func(call func(name string) func(any), handle ...ast.Ast) ast.Ast

func (f ProductionWith) call(p *parser, n uint, handle ...ast.Ast) status.Status {
	// pop stack (this removes `nodes`)
	p.stack.Clear(n) // must be called before pushing reduction result

	// create closure, capturing state of parser
	callRequester := func(name string) func(any) {
		return p.actions.get(name) // call action
	}
	// do reduction action
	reduced := f(callRequester, handle...)

	// save result
	p.stack.Push(reduced)

	return status.Ok
}
