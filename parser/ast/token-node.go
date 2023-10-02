package ast

import (
	"github.com/petersalex27/yew-packages/token"
)

type Token struct{ token.Token }

func (t Token) NodeType() Type {
	return Type(t.Token.GetType())
}

func TokenNode(tok token.Token) Token {
	return Token{tok}
}

func (t Token) Equals(a Ast) bool {
	t2, ok := a.(Token)
	if !ok {
		return false
	}
	ty1, val1 := t.Token.GetType(), t.Token.GetValue()
	ty2, val2 := t2.Token.GetType(), t2.Token.GetValue()
	return ty1 == ty2 && val1 == val2
}

func (t Token) InOrderTraversal(f func(token.Token)) { f(t.Token) }