package parser

import (
	"strconv"
	"time"

	"github.com/petersalex27/yew-packages/errors"
	"github.com/petersalex27/yew-packages/errors/warning"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/stack"
)

type DefaultErrorFunc func(source.StaticSource, token.Token) error

type Parser interface {
	Parse() ast.AstRoot
	LogActions() *loggableParser
	Load([]token.Token, source.StaticSource, DefaultErrorFunc, error) Parser
	action() status.Status
	ground() *parser
	actOnRule(rule_interface, []ast.Ast) (stat status.Status, ruleApplied bool)
	reportError(ast.Type) status.Status
	shift() status.Status
	reduce(rules ruleSet) (stat status.Status, appliedRule bool)
}

type parser struct {
	loaded bool
	knowledgeable_parser
	stack  *stack.Stack[ast.Ast]
	src    source.StaticSource
	tokens []token.Token
}

type blank_parser struct {
	errors        []error
	warnings      []warning.Warning
	couldNotParse error
	lookahead     func(*parser) lookahead_payload
	defaultError  func(source.StaticSource, token.Token) error
}

type knowledgeable_parser struct {
	blank_parser
	ReduceTable
}

func none_token() token.Token {
	return test_token{0, 0, ast.None, ""}
}

type lookahead_payload []token.Token

func (lap lookahead_payload) getType(p *parser) ast.Type {
	// create buffer
	tys := make([]ast.Type, len(lap))
	for i, tok := range lap {
		tys[i] = ast.Type(tok.GetType())
	}

	// combine to make type
	ty, found := p.root.get(tys...)
	if !found {
		return ast.None
	}
	return ty
}

func lookahead(k uint) func(*parser) lookahead_payload {
	if k == 0 {
		return func(*parser) lookahead_payload {
			return lookahead_payload{none_token()}
		}
	}

	return func(p *parser) lookahead_payload {
		if uint(len(p.tokens)) < k {
			if len(p.tokens) == 0 {
				return lookahead_payload{none_token()}
			}
			return lookahead_payload(p.tokens)
		}
		return p.tokens[:k]
	}
}

func default_lookahead(p *parser) lookahead_payload {
	if len(p.tokens) < 1 {
		return []token.Token{none_token()}
	}
	return []token.Token{p.tokens[0]}
}

func New() blank_parser {
	p := blank_parser{}

	p.errors = make([]error, 0, 1)
	p.warnings = make([]warning.Warning, 0, 1)
	p.lookahead = default_lookahead

	return p
}

func (b blank_parser) LA(k uint) blank_parser {
	if k == 1 {
		b.lookahead = default_lookahead
		return b
	}

	b.lookahead = lookahead(k)
	return b
}

func (p *parser) top() (ast.Ast, status.Status) {
	node, stat := p.stack.Peek()
	if stat.NotOk() {
		node = ast.Nothing{}
		return node, status.StackEmpty
	}
	return node, status.Ok
}

func (p *parser) lookAhead() (tok token.Token, stat status.Status) {
	if len(p.tokens) == 0 {
		stat = status.EndOfTokens
		return
	}

	tok, stat = p.tokens[0], status.Ok
	return
}

func (p *parser) shift() status.Status {
	tok, stat := p.lookAhead()
	if stat.IsOk() {
		p.tokens = p.tokens[1:]
		p.stack.Push(ast.TokenNode(tok))
	}
	return stat
}

func (p *parser) HasErrors() bool { return len(p.errors) != 0 }

func (p *parser) GetErrors() []error { return p.errors }

type forType ast.Type

func (a forType) modify(stat status.Status, appliedRule bool) status.Status {
	if a == forType(ast.None) {
		return status.EndOfParse
	} else if stat.Is(status.EndAction) && appliedRule {
		return status.DoShift
	}
	return stat
}

func (p *parser) action() status.Status {
	toks := p.lookahead(p)
	ty := toks.getType(p)

	rules, found := p.ReduceTable.table[ty]

	stat, ruleApplied := forType(ty).actionLoop(p, rules, found)
	return forType(ty).followUpRule(p, rules, stat, ruleApplied)
}

func (p *parser) actOnRule(rule rule_interface, vars []ast.Ast) (stat status.Status, appliedRule bool) {
	stat, appliedRule = rule.call(p, vars...), true
	return
}

func (p *parser) matchStack(pattern pattern) (nodes []ast.Ast, matches bool) {
	var stackStat stack.StackStatus
	nodes, stackStat = p.stack.MultiCheck(len(pattern))
	if stackStat.NotOk() {
		return nil, false
	}
	matches = p.ReduceTable.Match(pattern, nodes...)
	return
}

func getSubset(ground *parser, rules ruleSet) (subSet []rule_interface) {
	subSet = nil
	// grab the top node and see if any subsets of rules exist w/ the top node as
	// the last node
	if node, stat := ground.stack.Peek(); stat.IsOk() {
		var found bool
		subSet, found = rules.ruleMap[node.NodeType()]
		if !found {
			subSet = nil
		}
	}
	return
}

func (p *parser) reduce(rules ruleSet) (stat status.Status, appliedRule bool) {
	stat, appliedRule = status.EndAction, false

	// experimental - start
	subSet := getSubset(p, rules)
	// experimental - end

	for _, rule := range subSet {
		pattern := rule.getPattern()
		//nodes, stackStat := ground.stack.MultiCheck(len(pattern))
		//matches := ground.ReduceTable.Match(pattern, nodes...)
		nodes, match := p.matchStack(pattern)
		if match {
			stat, appliedRule = p.actOnRule(rule, nodes)
			break
		}
	}
	return
}

func (p *loggableParser) reduce(rules ruleSet) (stat status.Status, appliedRule bool) {
	return p.ground().reduce(rules)
}

func (ty forType) actionLoop(p Parser, rules ruleSet, found bool) (stat status.Status, appliedRule bool) {
	stat, appliedRule = status.Ok, false
	if !found {
		return status.EndAction, false
	}

	for tmpApp := false; stat.IsOk(); {
		stat, tmpApp = p.reduce(rules)
		appliedRule = appliedRule || tmpApp
	}
	stat = ty.modify(stat, appliedRule)
	return
}

func (p *parser) reportError(ty ast.Type) status.Status {
	if ty == ast.None {
		p.errors = append(p.errors, p.couldNotParse) // push default errors.Err b/c p.tokens[0] (probably) DNE
	} else {
		// p.tokens[0] must exist since tokType != ast.None which is the empty stack return value
		if len(p.tokens) == 0 { // sanity check
			panic("bug: len(p.tokens) should not be zero")
		}
		p.errors = append(p.errors, p.defaultError(p.src, p.tokens[0]))
	}
	return status.Error
}

func (ty forType) followUpRule(p Parser, rules ruleSet, stat status.Status, ruleApplied bool) status.Status {
	if ruleApplied {
		if stat.Is(status.DoShift) {
			p.shift()
		} // else, already shifted // TODO: is this possible?
	} else if rules.elseShift && ty != forType(ast.None) { // don't allow none to be shifted!
		stat = status.Ok
		p.shift()
	} else {
		return p.reportError(ast.Type(ty))
	}

	return stat.MakeOk()
}

func (p blank_parser) Ruleset(table ReduceTable) knowledgeable_parser {
	return knowledgeable_parser{
		blank_parser: p,
		ReduceTable:  table,
	}
}

func estimateStackUse(fullElemLen int) uint {
	return uint(fullElemLen/2) + 1 // this is just an estimate, idk; +1 is so cap != 0
}

func (p *parser) load_src_def(src source.StaticSource, def func(source.StaticSource, token.Token) error, couldNotSimplify error) *parser {
	if src == nil {
		src = nilsrc{}
	}
	p.src = src

	if couldNotSimplify == nil {
		couldNotSimplify = errors.Ferr("tfm", "Syntax", src.GetPath(), "input could not be parsed")
	}
	p.couldNotParse = couldNotSimplify

	if def == nil {
		def = func(src source.StaticSource, tok token.Token) error {
			line, char := tok.GetLineChar()
			srcline, _ := src.SourceLine(line)
			return errors.Ferr("tflcms", "Syntax", src.GetPath(), line, char, "unexpected token", srcline)
		}
	}
	p.defaultError = def

	return p
}

func (kp knowledgeable_parser) Load(tokens []token.Token, src source.StaticSource, def func(source.StaticSource, token.Token) error, couldNotParse error) *parser {
	p := new(parser)
	p.knowledgeable_parser = kp

	// init stack
	cap := estimateStackUse(len(tokens))
	p.stack = stack.NewStack[ast.Ast](cap)
	p.tokens = tokens
	p.loaded = true

	return p.load_src_def(src, def, couldNotParse)
}

func (p *parser) Benchmarker() benchmarker {
	b := benchmarker{parser: p}
	return b
}

func (p *parser) Load(tokens []token.Token, src source.StaticSource, def DefaultErrorFunc, couldNotParse error) Parser {
	ratio := float64(p.stack.GetCapacity()) / float64(estimateStackUse(len(tokens)))
	if ratio < 0.90 {
		return p.knowledgeable_parser.Load(tokens, src, def, couldNotParse)
	}

	p.tokens = tokens
	p.loaded = true
	return p.load_src_def(src, def, couldNotParse)
}

func default_stringType(ty ast.Type) string {
	return strconv.FormatInt(int64(ty), 10)
}

func (p *parser) LogActions() (out *loggableParser) {
	var act action_name
	if p.loaded {
		act = init_log
	} else {
		act = late_log
	}

	out = new(loggableParser)
	out.parser = p
	out.stringType = default_stringType

	t, e := time.Now().In(time.UTC).MarshalText()
	if e != nil {
		t = []byte("error--" + e.Error())
	}
	out.log2(true, act, "%s", string(t)) // init: <utc_time>
	return out
}

func (p *parser) ground() *parser { return p }

func parse(p Parser) ast.AstRoot {
	grnd := p.ground()
	if !grnd.loaded {
		panic("parser must be re-loaded before calling (Parser) Parse() again")
	}

	grnd.loaded = false

	stat := status.Ok
	for stat.IsOk() {
		stat = p.action()
	}

	if !stat.EndParse() {
		return ast.AstRoot{}
	}

	grnd = p.ground()

	root, _ := grnd.stack.MultiCheck(int(grnd.stack.GetCount()))
	grnd.stack.Clear(grnd.stack.GetCount())
	return ast.AstRoot(root)
}

func (p *parser) Parse() ast.AstRoot {
	return parse(p)
}
