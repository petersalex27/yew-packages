package parser

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-packages/token"
)

type loggableParser struct {
	*parser
	logger strings.Builder
	stringType func(ast.Type) string
}

type action_name string

const (
	late_log   action_name = "init(late):"
	init_log   action_name = "init:"
	shift_log  action_name = "shift:"
	reduce_log action_name = "reduce:"
	error_log  action_name = "error:"
	search_log action_name = "search:"
	la_ty_log action_name = "lookahead:"
	rule_log action_name = "rule(applied):"
	rule_not_app_log action_name = "rule(unapplied):"
	rules_found_log action_name = "rules?:"
	rules_log action_name = "rules:"
	action_end_log action_name = "action(end):"
)

func (p *loggableParser) shift() status.Status {
	tok, stat := p.lookAhead()
	p.log(shift_log, "tok=%v; stat=%v", tok, stat)
	return p.parser.shift()
}

func (p *loggableParser) actOnRule(rule rule_interface, vars []ast.Ast) (stat status.Status, appliedRule bool) {
	var act action_name = rule_log
	stat, appliedRule = p.parser.actOnRule(rule, vars)
	if !appliedRule {
		act = rule_not_app_log
	}
	p.log(act, "stat=%v; rule=%v; nodes=%v", stat, rule, vars)
	return
}

func (p *loggableParser) reportError(ty ast.Type) status.Status {
	stat := p.parser.reportError(ty)
	p.log(error_log, "\n%s", p.errors[len(p.errors)-1].Error())
	return stat
}

func (p *loggableParser) action() status.Status {
	toks := p.lookahead(p.parser)

	ty := toks.getType(p.parser)
	p.log(la_ty_log, "%s", p.stringType(ty))

	rules, found := p.table().table[ty]
	p.log(rules_found_log, "%t", found)
	p.log(rules_log, "%v", rules)

	stat, ruleApplied := forType(ty).actionLoop(p, rules, found)
	stat = forType(ty).followUpRule(p, rules, stat, ruleApplied)

	p.log(action_end_log, "stack=%s", p.stack.ElemString())
	return stat
}

func (p *loggableParser) Load(tokens []token.Token, src source.StaticSource, def DefaultErrorFunc, couldNotParse error) Parser {
	p.parser.Load(tokens, src, def, couldNotParse)
	return p
}

func (p *loggableParser) StringType(f func(ast.Type)string) *loggableParser {
	if f == nil { // don't allow // TODO: log this?
		return p
	}
	p.stringType = f
	return p
}

func (p *loggableParser) log2(init bool, action action_name, format string, args ...any) {
	format = string(action) + " " + format
	if format[len(format)-1] != '\n' {
		format = format + "\n"
	}
	in := fmt.Sprintf(format, args...)

	if p.logger.Len() == 0 && !init {
		p.LogActions()
	}

	p.logger.WriteString(in)
}

func (p *loggableParser) log(action action_name, format string, args ...any) {
	p.log2(false, action, format, args...)
}

func (p *loggableParser) LogMessage(header string, message string) {
	header, message = strings.TrimSpace(header), strings.TrimSpace(message)
	p.log(action_name(header), "%s: %s", header, message)
}

func (p *loggableParser) FlushLog() string {
	out := p.logger.String()
	p.logger.Reset()
	return out
}

func (p *loggableParser) ground() *parser { return p.parser }

func (p *loggableParser) Parse() ast.AstRoot {
	return parse(p)
}
