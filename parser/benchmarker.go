package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-packages/token"
)

type benchmarker struct {
	*parser
	optimizedReduce bool
}

func (p *benchmarker) ClearFlags() {
	p.optimizedReduce = false
}

func (p *benchmarker) SetOptimizedReduce(set bool) {
	p.optimizedReduce = set
}

func (p benchmarker) Parse() ast.AstRoot {
	return parse(p.parser)
}

func (p benchmarker) LogActions() *loggableParser {
	panic("do not call on benchmarker")
}

func (p benchmarker) Load(toks []token.Token, src source.StaticSource, def DefaultErrorFunc, e error) Parser {
	res := p.parser.Load(toks, src, def, e).(*parser)
	return res.Benchmarker()
}

func (p benchmarker) action() status.Status {
	return p.parser.action()
}

func (p benchmarker) ground() *parser {
	return p.parser
}

func (p benchmarker) actOnRule(ri productionInterface, ns []ast.Ast) (stat status.Status, ruleApplied bool) {
	return p.parser.actOnRule(ri, ns)
}

func (p benchmarker) reportError(t ast.Type) status.Status {
	return p.parser.reportError(t)
}

func (p benchmarker) shift() status.Status {
	return p.parser.shift()
}

func (p benchmarker) reduce(rules productionOrder) (stat status.Status, appliedRule bool) {
	if p.optimizedReduce {
		return p.ground().reduce(rules)
	}

	stat, appliedRule = status.EndAction, false
	ground := p.ground()

	for _, rule := range rules.rules {
		pattern := rule.getPattern()
		//nodes, stackStat := ground.stack.MultiCheck(len(pattern))
		//matches := ground.ReduceTable.Match(pattern, nodes...)
		nodes, match := ground.matchStack(pattern)
		if match {
			stat, appliedRule = p.actOnRule(rule, nodes)
			break
		}
	}
	return
}
