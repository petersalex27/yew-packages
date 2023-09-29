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