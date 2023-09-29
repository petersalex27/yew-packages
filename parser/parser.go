package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
	"github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/stack"
)

type Parser struct {
	stack *stack.Stack[ast.Ast]
	reductions *ReduceTable
	tokens []token.Token
}

func New(tokens []token.Token) *Parser {
	p := new(Parser)

	// init stack
	cap := uint32(len(tokens)/2)+1 // this is just an estimate, idk; +1 is so cap != 0
	p.stack = stack.NewStack[ast.Ast](cap)

	p.tokens = tokens
	return p
}

func (p *Parser) Top() (ast.Ast, status.Status) {
	node, stat := p.stack.Peek()
	if stat.NotOk() {
		return node, status.StackEmpty
	}
	return node, status.Ok
}

func (p *Parser) LookAhead() (tok token.Token, stat status.Status) {
	if len(p.tokens) == 0 {
		stat = status.EndOfTokens
		return
	}

	tok, stat = p.tokens[0], status.Ok
	return
}

func (p *Parser) LookAheadTokType() (ast.Ast, ast.Type) {
	tok, stat := p.LookAhead()
	if stat.NotOk() {
		return ast.Nothing{}, ast.None
	}
	a := ast.TokenNode(tok)
	return a, a.NodeType()
}

func (p *Parser) TopTokType() (ast.Ast, ast.Type) {
	node, stat := p.stack.Pop()
	if stat.IsEmpty() {
		return ast.Nothing{}, ast.None
	}
	return node, node.NodeType()
}

func (p *Parser) shift() status.Status {
	tok, stat := p.LookAhead()
	if stat.IsOk() {
		p.tokens = p.tokens[1:]
		p.stack.Push(ast.TokenNode(tok))
	}
	return stat
}

func (p *Parser) Action() status.Status {
	_, tokType := p.LookAheadTokType()
	rules, found := p.reductions.table[tokType]
	ruleApplied := false
	for found {
		found = false
		for _, rule := range rules {
			vars, stat, set := p.stack.MultiCheck(len(rule.pattern))
			if stat.IsOk() && rule.pattern.equals_len_known(vars) {
				ruleApplied = true
				set()
				res := rule.Reduction(vars...)
				p.stack.Push(res)
				found = true
				break
			}
		}
	}

	if ruleApplied {
		p.shift()
		return status.Ok
	}
	return status.NoAction
}