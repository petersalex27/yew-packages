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
	"github.com/petersalex27/yew-packages/util/iterator"
	"github.com/petersalex27/yew-packages/util/stack"
)

type DefaultErrorFunc func(source.StaticSource, token.Token) error

type RunnableParser interface {
	Parse() ast.AstRoot
	LogActions() *loggableParser
	Load([]token.Token, source.StaticSource, DefaultErrorFunc, error) RunnableParser
	Parser
}

type MakeableNode interface {
	ast.Ast
	MakeNode(token.Token) ast.Ast
}

type Parser interface {
	action() status.Status
	ground() *parser
	actOnRule(productionInterface, []ast.Ast) (stat status.Status, ruleApplied bool)
	reportError(ast.Type) status.Status
	Shift() status.Status
	SetShift(MakeableNode) status.Status
	Reduce(rules ProductionOrder) (stat status.Status, appliedRule bool)
	GetSource() source.StaticSource
}

type parser struct {
	loaded bool
	knowledgeable_parser
	stack   *stack.Stack[ast.Ast]
	src     source.StaticSource
	tokens  []token.Token
	actions actionRequester
}

func (p *parser) GetSource() source.StaticSource {
	return p.src
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
	cxt *ParserContext
}

func LookAheadTokens(p Parser) []token.Token {
	return p.ground().lookahead(p.ground())
}

func (p *parser) dropNext() (dropped bool) {
	dropped = len(p.tokens) >= 1
	if dropped {
		p.tokens = p.tokens[1:]
	}
	return
}

func DropNext(p Parser) {
	_ = p.ground().dropNext()
}

// returns root of combiner trie
//
// panics if parser has no context attached
func (kp knowledgeable_parser) root() *combinerTrieRoot {
	if kp.cxt == nil {
		panic("missing parse context")
	}
	return kp.cxt.currentTable.root
}

// returns current parsing table
//
// panics if parser has no context attached
func (kp knowledgeable_parser) table() *ReductionTable {
	if kp.cxt == nil {
		panic("missing parse context")
	}
	return &kp.cxt.currentTable
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
	ty, found := p.root().get(tys...)
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

func NewParser() blank_parser {
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

func (p *parser) SetShift(node MakeableNode) status.Status {
	tok, stat := p.lookAhead()
	if stat.IsOk() {
		_ = p.dropNext()
		in := node.MakeNode(tok)
		p.stack.Push(in)
	}
	return stat
}

func (p *parser) Shift() status.Status {
	tok, stat := p.lookAhead()
	if stat.IsOk() {
		_ = p.dropNext()
		p.stack.Push(ast.TokenNode(tok))
	}
	return stat
}

func (p *parser) HasErrors() bool { return len(p.errors) != 0 }

func (p *parser) GetErrors() []error { return p.errors }

func HasErrors(p Parser) bool {
	return p.ground().HasErrors()
}

func GetErrors(p Parser) []error {
	return p.ground().GetErrors()
}

// return top element of parse stack
//
// second return value is trye iff status of peeking at stack is Ok
func ParseStackPeek(p Parser) (ast.Ast, bool) {
	elem, stat := p.ground().stack.Peek()
	return elem, stat.IsOk()
}

// pushes ast node onto parse stack only when it isn't a None node
func (p *parser) maybePush(node ast.Ast) {
	if node.NodeType() != ast.None {
		p.stack.Push(node)
	}
}

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

	rules, found := p.table().table[ty]

	stat, ruleApplied := forType(ty).actionLoop(p, rules, found)
	return forType(ty).followUpRule(p, rules, stat, ruleApplied)
}

func (p *parser) actOnRule(rule productionInterface, handle []ast.Ast) (stat status.Status, appliedRule bool) {
	stat, appliedRule = rule.call(p, handle...)
	return
}

// tries to match
func (p *parser) matchStack(pattern PatternInterface) (nodes []ast.Ast, matches bool) {
	var stackStat stack.StackStatus
	handleLength := pattern.MaxHandleLength()
	nodes, stackStat = p.stack.MultiCheck(handleLength)
	if stackStat.NotOk() {
		return nil, false
	}
	matches = pattern.Match(nodes...)
	return
}

func getSubset(ground *parser, rules ProductionOrder) (subSet []productionInterface) {
	subSet = nil
	// grab the top node and see if any subsets of rules exist w/ the top node as
	// the last node
	if node, stat := ground.stack.Peek(); stat.IsOk() && rules.classes != nil {
		subSet, _ = rules.classes.getClass(node.NodeType()) // subSet = nil if not found
	}
	return
}

// return true
//
//	precondition check failed: appliedRule=false, stat=Ok
//
// return ?: (depends on rule held by precondition)
//
//	precondition check success: appliedRule=?, stat=?
//
// return false
//
//	error rule: appliedRule=true, stat=Ok
//	shift rule: appliedRule=true, stat=Shift
//	...       : appliedRule=true, stat=...
func continueLoop(stat status.Status, appliedRule bool) bool {
	return !appliedRule && stat.IsOk()
}

func initialStatAndApplied() (status.Status, bool) {
	return status.EndAction, false
}

func (p *parser) matchThen(
	pattern PatternInterface,
	rule productionInterface,
) (status.Status, bool, bool) {
	stat, appliedRule := initialStatAndApplied()
	loop := true
	if nodes, match := p.matchStack(pattern); match {
		stat, appliedRule = p.actOnRule(rule, nodes)
		loop = continueLoop(stat, appliedRule)
	}
	return stat, appliedRule, loop
}

// Parser Reduce action: replaces parse-stack handle with reduction rule
// replacement. Returns reduction status along with the truthy-ness of whether
// an actual rule was applied
func (p *parser) Reduce(rules ProductionOrder) (stat status.Status, appliedRule bool) {
	stat, appliedRule = initialStatAndApplied()

	subSet := getSubset(p, rules)

	it := iterator.Iterator(subSet)
	rule, ok := it.Next()

	for loop := true; ok && loop; rule, ok = it.Next() {
		pattern := rule.getPattern()
		stat, appliedRule, loop = p.matchThen(pattern, rule)
	}
	return
}

func (p *loggableParser) Reduce(rules ProductionOrder) (stat status.Status, appliedRule bool) {
	return p.ground().Reduce(rules)
}

func (ty forType) actionLoop(p Parser, rules ProductionOrder, found bool) (stat status.Status, appliedRule bool) {
	stat, appliedRule = status.Ok, false
	if !found {
		return status.EndAction, false
	}

	for tmpApp := false; stat.IsOk(); {
		stat, tmpApp = p.Reduce(rules)
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

func (ty forType) followUpRule(p Parser, rules ProductionOrder, stat status.Status, ruleApplied bool) status.Status {
	if ruleApplied {
		if stat.Is(status.DoShift) {
			p.Shift()
		} // else, already shifted // TODO: is this possible?
	} else if rules.elseShift && ty != forType(ast.None) { // don't allow none to be shifted!
		stat = status.Ok
		p.Shift()
	} else {
		return p.reportError(ast.Type(ty))
	}

	return stat.MakeOk()
}

func (p blank_parser) UsingReductionTable(table ReductionTable) knowledgeable_parser {
	kp := knowledgeable_parser{
		blank_parser: p,
		cxt:          new(ParserContext),
	}
	*kp.cxt = makeParserContext(ast.None, table, 1)
	return kp
}

func SetDefaultErrorGen(p Parser, defaultErrorGen func(source.StaticSource, token.Token) error) Parser {
	p.ground().loadDefaultErrorGen(p.ground().src, defaultErrorGen)
	return p
}

func SetCouldNotParseError(p Parser, couldNotParse error) Parser {
	p.ground().loadCouldNotParse(p.ground().src.GetPath(), couldNotParse)
	return p
}

// use when parser does not have a parse table
func (bp blank_parser) Load(tokenStream []token.Token, src source.StaticSource, defaultErrorGen func(source.StaticSource, token.Token) error, couldNotParse error) *parser {
	kp := knowledgeable_parser{
		blank_parser: bp,
		cxt:          nil,
	}
	return kp.Load(tokenStream, src, defaultErrorGen, couldNotParse)
}

func (p knowledgeable_parser) Mapping(ty ast.Type, table ReductionTable) knowledgeable_parser {
	p.cxt.MapTable(ty, table)
	return p
}

func estimateStackUse(fullElemLen int) uint {
	return uint(fullElemLen/2) + 1 // this is just an estimate, idk; +1 is so cap != 0
}

func (p *parser) loadSource(src source.StaticSource) {
	if src == nil {
		src = EmptySource{}
	}
	p.src = src
}

func (p *parser) loadCouldNotParse(path string, couldNotParse error) {
	if couldNotParse == nil {
		couldNotParse = errors.Ferr("tfm", "Syntax", path, "input could not be parsed")
	}
	p.couldNotParse = couldNotParse
}

func (p *parser) loadDefaultErrorGen(src source.StaticSource, defaultErrorGen func(source.StaticSource, token.Token) error) {
	if defaultErrorGen == nil {
		defaultErrorGen = func(src source.StaticSource, tok token.Token) error {
			line, char := tok.GetLineChar()
			srcline, _ := src.SourceLine(line)
			return errors.Ferr("tflcms", "Syntax", src.GetPath(), line, char, "unexpected token", srcline)
		}
	}
	p.defaultError = defaultErrorGen
}

func (kp knowledgeable_parser) Load(tokens []token.Token, src source.StaticSource, defaultErrorGen func(source.StaticSource, token.Token) error, couldNotParse error) *parser {
	p := new(parser)
	p.knowledgeable_parser = kp

	// init stack
	cap := estimateStackUse(len(tokens))
	p.stack = stack.NewStack[ast.Ast](cap)

	p.tokens = tokens
	p.loaded = true

	p.loadSource(src)

	p.loadCouldNotParse(src.GetPath(), couldNotParse)
	p.loadDefaultErrorGen(src, defaultErrorGen)

	return p
}

// InitialStackPush pushes ast nodes onto the parse stack. This function panics
// if the parser has already started parsing
func (p *parser) InitialStackPush(nodes ...ast.Ast) *parser {
	// check if parser has already started parsing (marks itself as unloaded at
	// start of parse)
	if !p.loaded {
		panic("illegal operation: stack cannot be initialized once parser has started parsing")
	}

	for _, node := range nodes {
		p.stack.Push(node)
	}
	return p
}

func InitialStackPush(p Parser, nodes ...ast.Ast) {
	p.ground().InitialStackPush(nodes...)
} 

func (p *parser) Benchmarker() benchmarker {
	b := benchmarker{parser: p}
	return b
}

func (p *parser) Load(tokens []token.Token, src source.StaticSource, defaultErrorGen DefaultErrorFunc, couldNotParse error) Parser {
	ratio := float64(p.stack.GetCapacity()) / float64(estimateStackUse(len(tokens)))
	if ratio < 0.90 { // TODO: huh? what's going on here
		return p.knowledgeable_parser.Load(tokens, src, defaultErrorGen, couldNotParse)
	}

	p.tokens = tokens
	p.loaded = true

	p.loadSource(src)

	p.loadCouldNotParse(src.GetPath(), couldNotParse)
	p.loadDefaultErrorGen(src, defaultErrorGen)

	return p
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

// runs parser to completion or error: whichever happens first
//
// returns resulting AST root (which is empty on error)
//
// panics if parse context is not set
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
