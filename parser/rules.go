package parser

import (
	"fmt"
	"sort"

	"sync"

	"github.com/petersalex27/yew-packages/errors"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
	"github.com/petersalex27/yew-packages/stringable"
)

var sortLock sync.Mutex
var doSort = false

func RuleSetSort(sortRules bool) {
	sortLock.Lock()
	defer sortLock.Unlock()
	doSort = sortRules
}

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

func (r reduction) call(p *parser, nodes ...ast.Ast) status.Status {
	// pop stack (this removes `nodes`)
	n := uint(len(nodes))
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
	return r.reduction.call(p, nodes...)
}

type error_rule struct {
	pattern
	ErrorFn
}

func (r error_rule) String() string {
	return fmt.Sprintf("rule(%v -> error)", r.pattern)
}

func (r error_rule) getPattern() pattern { return r.pattern }

func (r error_rule) call(p *parser, nodes ...ast.Ast) status.Status {
	node, _ := p.top()
	p.errors = append(p.errors, r.ErrorFn(node, nodes...))
	return status.Error
}

type warning_rule struct {
	pattern
	WarnFn
}

func (r warning_rule) String() string {
	return fmt.Sprintf("rule(%v -> (warning, %v))", r.pattern, r.Reduction)
}

func (r warning_rule) getPattern() pattern { return r.pattern }

func (r warning_rule) call(p *parser, nodes ...ast.Ast) status.Status {
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

type shift_rule pattern

func (r shift_rule) getPattern() pattern { return pattern(r) }

func (r shift_rule) call(p *parser, _ ...ast.Ast) (stat status.Status) {
	if stat = p.shift(); stat.IsOk() {
		stat = status.EndAction
	}
	return stat
}

func (r shift_rule) String() string {
	return fmt.Sprintf("rule(%v -> shift)", pattern(r))
}

type end_rule reduction

func (r end_rule) getPattern() pattern {
	return nil
}

func (r end_rule) call(p *parser, nodes ...ast.Ast) (stat status.Status) {
	return reduction(r).call(p, nodes...)
}

type pattern []ast.Type

/*
RuleSet(
	Rule(Id, TypeJudgement, Id, Assign, Value).Reduce(assignment),
	Rule(Id, Assign, Value).Reduce(assignment),
)
*/

func Rule(types ...ast.Type) pattern { return append(pattern{}, types...) }

type ruleSet struct {
	rules      []rule_interface
	shiftAtEnd bool
}

func (rs ruleSet) Less(i, j int) bool {
	return len(rs.rules[i].getPattern()) < len(rs.rules[j].getPattern())
}

func (rs ruleSet) Len() int { return len(rs.rules) }

func (rs ruleSet) Swap(i, j int) { rs.rules[i], rs.rules[j] = rs.rules[j], rs.rules[i] }

func (rs ruleSet) Union(ruleSets ...ruleSet) ruleSet {
	out := make([]rule_interface, 0, len(rs.rules))
	out = append(out, rs.rules...)
	ruleSetOut := ruleSet{rules: out, shiftAtEnd: rs.shiftAtEnd}
	for _, set := range ruleSets {
		ruleSetOut.rules = append(ruleSetOut.rules, set.rules...)
		ruleSetOut.shiftAtEnd = ruleSetOut.shiftAtEnd || set.shiftAtEnd
	}
	if doSort {
		sort.Sort(ruleSetOut)
	}
	return ruleSetOut
}

func RuleSet(rules ...rule_interface) (out ruleSet) {
	out.rules = append([]rule_interface{}, rules...)
	if doSort {
		sort.Sort(out)
	}
	return out
}

func isShiftRule(rs ruleSet) bool {
	return rs.rules == nil && rs.shiftAtEnd
}

// creates a reduce action
func (p pattern) Reduce(f Reduction) rule_interface { return rule{p, reduction{f}} }

// creates a shift action
func (p pattern) Shift() rule_interface { return shift_rule(p) }

func (p pattern) Error(e ErrorFn) rule_interface { return error_rule{p, e} }

func (p pattern) Warn(w WarnFn) rule_interface { return warning_rule{p, w} }