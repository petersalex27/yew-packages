package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
	"github.com/petersalex27/yew-packages/stringable"
)

type callable interface {
	/*
		call does a rule's action based off the state of the parser and the nodes
		passed as arguments
	*/
	call(p *parser, nodes ...ast.Ast) (stat status.Status)
}

type rule_interface interface {
	stringable.Stringable
	callable
	getPattern() pattern
}
