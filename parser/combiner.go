package parser

import "github.com/petersalex27/yew-packages/parser/ast"

func initRoot(lastType ast.Type) *combinerTrieRoot {
	out := new(combinerTrieRoot)
	if lastType >= ast.Root-1 {
		panic("too many types, cannot initialize")
	}

	out.counter = uint(lastType) + 1
	// makes the assumption that all types will have an entry
	out.nodes = make(map[ast.Type]*combinerTrieNode, uint(lastType))
	return out
}

type combinerTrieRoot struct {
	counter uint
	nodes   map[ast.Type]*combinerTrieNode
}

type combinerTrieNode struct {
	ast.Type
	children map[ast.Type]*combinerTrieNode
}

func (node *combinerTrieNode) get(tys ...ast.Type) (ast.Type, bool) {
	if len(tys) == 0 {
		return node.Type, true
	}

	node2, found := node.children[tys[0]]
	if !found {
		return ast.None, false
	}
	return node2.get(tys[1:]...)
}

func (root *combinerTrieRoot) get(tys ...ast.Type) (ast.Type, bool) {
	if len(tys) == 0 {
		panic("len(tys) cannot be zero")
	}

	if len(tys) == 1 {
		return tys[0], true
	}

	res, found := root.nodes[tys[0]]
	if !found {
		return ast.None, found
	}

	return res.get(tys[1:]...)
}

func (node *combinerTrieNode) set(root *combinerTrieRoot, tys ...ast.Type) ast.Type {
	if node.children == nil {
		node.children = make(map[ast.Type]*combinerTrieNode)
	}

	first := tys[0]
	node2, found := node.children[first]
	var out ast.Type
	if !found {
		index := 1
		attach := new(combinerTrieNode)
		out = ast.Type(root.counter)
		attach.Type, attach.children = out, nil
		root.counter++
		var inner *combinerTrieNode = attach
		for index < len(tys) {
			inner.children = make(map[ast.Type]*combinerTrieNode, 2)
			next := new(combinerTrieNode)
			out = ast.Type(root.counter)
			next.Type, next.children = out, nil
			inner.children[tys[index]] = next

			inner = next
			index++
			root.counter++
		}
		node.children[first] = attach
	} else {
		out = node2.Type
		if len(tys) == 1 {
			return out
		}

		return node2.set(root, tys[1:]...)
	}

	return out
}

func (root *combinerTrieRoot) set(tys ...ast.Type) ast.Type {
	if len(tys) < 1 {
		panic("must provide at least one types")
	}

	first := tys[0]

	var out ast.Type
	if node, found := root.nodes[first]; found {
		if len(tys) == 1 {
			return node.Type
		}
		return node.set(root, tys[1:]...)
	} else {
		index := 1
		attach := new(combinerTrieNode)
		out = first
		attach.Type, attach.children = out, nil
		var inner *combinerTrieNode = attach
		for index < len(tys) {
			inner.children = make(map[ast.Type]*combinerTrieNode, 2)
			next := new(combinerTrieNode)
			out = ast.Type(root.counter)
			next.Type, next.children = out, nil
			inner.children[tys[index]] = next

			inner = next
			index++
			root.counter++
		}
		root.nodes[first] = attach
	}
	return out
}
