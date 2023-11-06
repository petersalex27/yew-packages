package parser

import (
	"fmt"
	"os"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/token"
)

type testType_t ast.Type

func (t testType_t) String() string {
	switch ast.Type(t) {
	case expr_t:
		return "expr_t"
	case assign_t:
		return "assign_t"
	case decl_t:
		return "decl_t"
	case fn_t:
		return "fn_t"
	case mul_t:
		return "mul_t"
	case add_t:
		return "add_t"
	case integer_t:
		return "integer_t"
	case id_t:
		return "id_t"
	case let_t:
		return "let_t"
	case in_t:
		return "in_t"
	case context_t:
		return "context_t"
	case func_t:
		return "func_t"
	case open_t:
		return "open_t"
	case close_t:
		return "close_t"
	case comma_t:
		return "comma_t"
	case ast.None:
		return "None"
	case ast.Root:
		return "Root"
	}
	return "testType_" + default_stringType(ast.Type(t))
}

func test_token_stringType(ty ast.Type) string {
	return testType_t(ty).String() + "(" + default_stringType(ty) + ")"
}

type test_token struct {
	line, char int
	ty         ast.Type
	val        string
}

const (
	expr_t ast.Type = iota
	assign_t
	decl_t
	fn_t
	mul_t
	add_t
	integer_t
	id_t
	let_t
	in_t
	context_t
	func_t
	open_t
	close_t
	comma_t
	div_t
	assertion_t
	error_t
	lastType_t_
)

func (tok test_token) String() string {
	return fmt.Sprintf("test_token@[%d:%d]:%s=%s",
		tok.line, tok.char, testType_t(tok.ty).String(), tok.val)
}

func (tok test_token) SetLineChar(line, char int) token.Token {
	return test_token{
		line: line,
		char: char,
		ty:   tok.ty,
		val:  tok.val,
	}
}

func (tok test_token) GetLineChar() (line, char int) { return tok.line, tok.char }

func (tok test_token) SetType(ty uint) token.Token {
	return test_token{
		line: tok.line,
		char: tok.char,
		ty:   ast.Type(ty),
		val:  tok.val,
	}
}

func (tok test_token) GetType() uint { return uint(tok.ty) }

func (tok test_token) SetValue(value string) (token.Token, error) {
	return test_token{
		line: tok.line,
		char: tok.char,
		ty:   tok.ty,
		val:  value,
	}, nil
}

func (tok test_token) GetValue() string { return tok.val }

type ast_test_node struct {
	ast.Type
	data []ast.Ast
}

func (node ast_test_node) String() string {
	return fmt.Sprintf(
		"ast_test_node:%s=%v",
		testType_t(node.Type).String(), node.data)
}

func mknode(ty ast.Type, nodes ...ast.Ast) (out ast_test_node) {
	out = ast_test_node{Type: ty, data: make([]ast.Ast, len(nodes))}
	for i := range out.data {
		out.data[i] = nodes[i]
	}
	return
}

func (a ast_test_node) Equals(b ast.Ast) bool {
	a2, ok := b.(ast_test_node)
	if !ok {
		return false
	}

	if a.Type != a2.Type {
		return false
	}

	if len(a.data) != len(a2.data) {
		return false
	}

	for i := range a.data {
		if !a.data[i].Equals(a2.data[i]) {
			return false
		}
	}
	return true
}

func (a ast_test_node) InOrderTraversal(f func(token.Token)) {
	for _, d := range a.data {
		d.InOrderTraversal(f)
	}
}

func test_reduce_fn(ty ast.Type) ProductionFunction {
	return func(nodes ...ast.Ast) ast.Ast { return mknode(ty, nodes...) }
}

func (atn ast_test_node) NodeType() ast.Type { return atn.Type }

func writeLog(log string) {
	f, e := os.Create("./test-logs/test-parser-log.txt")
	if e == nil {
		defer f.Close()
		_, e = f.WriteString(log)
	}

	if e != nil {
		print("failed to write logs to file with: ", e.Error(), "\nlog:\n", log)
		return
	}
}
