package ast

import (
	"github.com/petersalex27/yew-packages/token"
)

type token_node struct{ token.Token }

func (t token_node) NodeType() Type {
	return Type(t.Token.GetType())
}

func TokenNode(tok token.Token) token_node {
	return token_node{tok}
}

func (t token_node) Equals(a Ast) bool {
	t2, ok := a.(token_node)
	if !ok {
		return false
	}
	ty1, val1 := t.Token.GetType(), t.Token.GetValue()
	ty2, val2 := t2.Token.GetType(), t2.Token.GetValue()
	return ty1 == ty2 && val1 == val2
}

func (t token_node) InOrderTraversal(f func(token.Token)) { f(t.Token) }