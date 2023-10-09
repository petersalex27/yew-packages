package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/errors"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
	"github.com/petersalex27/yew-packages/stringable"
)

type Reduction interface {
	stringable.Stringable
	function(nodes ...ast.Ast) ast.Ast
}

type ReductionFunction func(nodes ...ast.Ast) ast.Ast

func (f ReductionFunction) GiveName(name string) NamedReduction {
	return NamedReduction{
		Name:              name,
		ReductionFunction: f,
	}
}

func (f ReductionFunction) function(nodes ...ast.Ast) ast.Ast {
	return f(nodes...)
}

func (f ReductionFunction) String() string {
	return "reduction_function"
}

type NamedReduction struct {
	Name string
	ReductionFunction
}

func (f NamedReduction) function(nodes ...ast.Ast) ast.Ast {
	return f.ReductionFunction(nodes...)
}

func (f NamedReduction) String() string {
	return f.Name
}

type ErrorFn func(top ast.Ast, nodes ...ast.Ast) errors.Err

type WarnFn struct {
	Warn func(top ast.Ast, nodes ...ast.Ast) errors.Warning
	Reduction
}

type reduction struct{ Reduction }

func (r reduction) call(p *parser, n uint, nodes ...ast.Ast) status.Status {
	// pop stack (this removes `nodes`)
	p.stack.Clear(n) // must be called before pushing reduction result
	// do reduction action
	result := r.function(nodes...)
	p.stack.Push(result)
	return status.Ok
}

type rule struct {
	pattern
	reduction
}

func (r rule) String() string {
	return fmt.Sprintf("rule(%v -> %v)", r.pattern, r.Reduction)
}

func (r rule) getPattern() pattern { return r.pattern }

func (r rule) call(p *parser, nodes ...ast.Ast) status.Status {
	return r.reduction.call(p, uint(len(nodes)), nodes...)
}

type whenRule struct {
	pattern
	clear uint
	reduction
}

func (r whenRule) String() string {
	return fmt.Sprintf("when(%v / %d -> %v)", r.pattern, r.clear, r.Reduction)
}

func (r whenRule) getPattern() pattern { return r.pattern }

func (r whenRule) call(p *parser, nodes ...ast.Ast) status.Status {
	return r.reduction.call(p, r.clear, nodes[uint(len(nodes))-r.clear:]...)
}

type errorRule struct {
	pattern
	ErrorFn
}

func (r errorRule) String() string {
	return fmt.Sprintf("rule(%v -> error)", r.pattern)
}

func (r errorRule) getPattern() pattern { return r.pattern }

func (r errorRule) call(p *parser, nodes ...ast.Ast) status.Status {
	node, _ := p.top()
	p.errors = append(p.errors, r.ErrorFn(node, nodes...))
	return status.Error
}

type warningRule struct {
	pattern
	WarnFn
}

func (r warningRule) String() string {
	return fmt.Sprintf("rule(%v -> (warning, %v))", r.pattern, r.Reduction)
}

func (r warningRule) getPattern() pattern { return r.pattern }

func (r warningRule) call(p *parser, nodes ...ast.Ast) status.Status {
	// look ahead token is for warning information; it should not be used in
	// reduction
	tok, _ := p.lookAhead()

	// pop stack (this removes `nodes`)
	n := uint(len(nodes))
	p.stack.Clear(n) // must be called before pushing reduction result

	// do warning and reduction (reduction happens inside r.WarnFn)
	warning := r.WarnFn.Warn(ast.TokenNode(tok), nodes...)
	res := r.function(nodes...)
	p.warnings = append(p.warnings, warning)
	p.stack.Push(res)
	return status.Ok
}

type shiftRule pattern

func (r shiftRule) getPattern() pattern { return pattern(r) }

func (shiftRule) call(_ *parser, _ ...ast.Ast) status.Status {
	return status.DoShift // this will trigger shift
}

func (r shiftRule) String() string {
	return fmt.Sprintf("rule(%v -> shift)", pattern(r))
}

type pattern []ast.Type

/*
RuleSet(
	Rule(Id, TypeJudgement, Id, Assign, Value).Reduce(assignment),
	Rule(Id, Assign, Value).Reduce(assignment),
)
*/

//func Rule(types ...ast.Type) pattern { return append(pattern{}, types...) }

type ruleSet struct {
	rules     []rule_interface
	// maps last type in rule_interface pattern to a subset of rule interfaces
	ruleMap		map[ast.Type][]rule_interface // NOTE: expiremental
	elseShift bool
}

func mapSet(set ruleSet, m *map[ast.Type][]rule_interface) {
	for _, rule := range set.rules {
		pat := rule.getPattern()
		if len(pat) != 0 {
			last := pat[len(pat)-1]
			res, found := (*m)[last]
			if !found {
				res = []rule_interface{rule}
			} else {
				res = append(res, rule)
			}
			(*m)[last] = res
		}
	}
}

func Union(sets ...ruleSet) (unified ruleSet) {
	unified = ruleSet{
		rules: []rule_interface{}, 
		ruleMap: make(map[ast.Type][]rule_interface), 
		elseShift: false,
	}
	for _, set := range sets {
		mapSet(set, &unified.ruleMap)
		unified.rules = append(unified.rules, set.rules...)
		unified.elseShift = unified.elseShift || set.elseShift
	}
	return
}

func RuleSet(rules ...rule_interface) (out ruleSet) {
	out.rules = append([]rule_interface{}, rules...)
	out.ruleMap = make(map[ast.Type][]rule_interface)
	mapSet(out, &out.ruleMap)
	return out
}

type needPattern reduction

type whenNeedPattern struct {
	needPattern
	when []ast.Type
}

func (f NamedReduction) From(tys ...ast.Type) rule_interface {
	return needPattern(reduction{f}).From(tys...)
}

func (f NamedReduction) When(tys ...ast.Type) whenNeedPattern {
	return needPattern(reduction{f}).When(tys...)
}

func Get(f func(...ast.Ast) ast.Ast) needPattern {
	return needPattern(reduction{ReductionFunction(f)})
}

func (p needPattern) When(tys ...ast.Type) whenNeedPattern {
	return whenNeedPattern{p, tys}
}

func (p needPattern) From(tys ...ast.Type) rule_interface {
	return rule{pattern(tys), reduction{p}}
}

func (p whenNeedPattern) From(tys ...ast.Type) rule_interface {
	clear := len(tys)
	pat := make(pattern, len(p.when)+clear)
	copy(pat, p.when)
	copy(pat[len(p.when):], tys)

	return whenRule{
		pattern(pat), 
		uint(clear), 
		reduction{p.needPattern},
	}
}

func From(tys ...ast.Type) pattern { return pattern(tys) }

func Shift() (s shiftRule) {return}

func (shiftRule) When(tys ...ast.Type) rule_interface {
	return shiftRule(tys)
}

type errorNeedsPattern errorRule

func Error(e func(top ast.Ast, nodes ...ast.Ast) errors.Err) errorNeedsPattern {
	return errorNeedsPattern{ErrorFn: e}
}

func (e errorNeedsPattern) From(tys ...ast.Type) rule_interface {
	return errorRule{
		pattern: pattern(tys),
		ErrorFn: e.ErrorFn,
	}
}

type warnNeedsReduction warningRule

type warnNeedPattern warningRule

func Warn(warnFn func(top ast.Ast, nodes ...ast.Ast) errors.Warning) warnNeedsReduction { 
	return warnNeedsReduction(warningRule{WarnFn: WarnFn{Warn: warnFn}})
}

func (w warnNeedsReduction) ThenGet(f func(nodes ...ast.Ast) ast.Ast) warnNeedPattern {
	return warnNeedPattern{WarnFn: WarnFn{Warn: w.Warn, Reduction: ReductionFunction(f)}}
}

func (w warnNeedPattern) From(tys ...ast.Type) rule_interface {
	return warningRule{
		pattern: pattern(tys),
		WarnFn: w.WarnFn,
	}
}