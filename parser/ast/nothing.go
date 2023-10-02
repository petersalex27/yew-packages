package ast

import "github.com/petersalex27/yew-packages/token"

type Nothing struct{}

func (Nothing) NodeType() Type {
	return None
}

func (Nothing) Equals(a Ast) bool {
	_, ok := a.(Nothing)
	return ok
} 

func (Nothing) InOrderTraversal(func(token.Token)) {}