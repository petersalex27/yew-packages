package ast

import (
	"math"
	"strings"

	"github.com/petersalex27/yew-packages/equality"
	"github.com/petersalex27/yew-packages/token"
)

// Type is the type of ast node types; math.MaxUint and math.MaxUint-1 are reserved!
type Type uint

const None Type = math.MaxUint
const Root Type = math.MaxUint - 1

type Ast interface {
	equality.Eq[Ast]
	NodeType() Type
	InOrderTraversal(func(token.Token))
}

func GetOrderedString(a Ast) string {
	var builder *strings.Builder = new(strings.Builder)
	f := func(tok token.Token) {
		builder.WriteString(tok.GetValue() + " ")
	}
	a.InOrderTraversal(f)
	return builder.String()
}

type AstRoot []Ast

func (root AstRoot) Equals(a Ast) bool {
	root2, ok := a.(AstRoot)
	if !ok {
		return false
	}

	if root.NodeType() != root2.NodeType() {
		return false
	}

	if len(root) != len(root2) {
		return false
	}

	for i, node := range root {
		if !node.Equals(root2[i]) {
			return false
		}
	}
	return true
}

func (AstRoot) NodeType() Type { return Root }

func (root AstRoot) InOrderTraversal(f func(token.Token)) {
	for _, node := range root {
		node.InOrderTraversal(f)
	}
}