package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestParser(t *testing.T) {
	addTok, _ := (test_token{}).SetType(uint(add_t)).SetValue("+")
	mulTok, _ := (test_token{}).SetType(uint(mul_t)).SetValue("*")
	intTok := (test_token{}).SetType(uint(integer_t))
	idTok := (test_token{}).SetType(uint(id_t))
	letTok, _ := (test_token{}).SetType(uint(let_t)).SetValue("let")
	inTok, _ := (test_token{}).SetType(uint(in_t)).SetValue("in")

	decl_fn := test_reduce_fn(decl_t).GiveName("declaration")
	expr_fn := test_reduce_fn(expr_t).GiveName("expression")
	assign_fn := test_reduce_fn(assign_t).GiveName("assignment")
	context_fn := test_reduce_fn(context_t).GiveName("context")

	set_decl := RuleSet(
		Rule(let_t, id_t).Reduce(decl_fn),
	)

	set_expr := RuleSet(
		Rule(id_t).Reduce(expr_fn), // Id -> expr
		Rule(integer_t).Reduce(expr_fn), // Id -> expr
		Rule(add_t, expr_t, expr_t).Reduce(expr_fn), // Add expr expr -> expr
		Rule(mul_t, expr_t, expr_t).Reduce(expr_fn), // Mul expr expr -> expr
	)

	set_id := RuleSet(Rule(let_t).Shift())

	set_assign := RuleSet(Rule(decl_t, expr_t).Reduce(assign_fn))
	set_decl_expr := set_decl.Union(set_expr)
	set_expr_assign := set_expr.Union(set_assign)
	set_expr_w_cxt := RuleSet(Rule(context_t, expr_t).Reduce(expr_fn))

	set_context := RuleSet(
		Rule(assign_t, in_t).Reduce(context_fn), // assign In -> context
	)

	table := 
		ForTypesThrough(lastType_t_).
		UseReductions(
			Map(let_t).Shift(),
			Map(id_t).To(set_id.Union(set_decl_expr)).ElseShift(),
			Map(integer_t).To(set_decl_expr).ElseShift(),
			Map(add_t).To(set_decl_expr).ElseShift(),
			Map(mul_t).To(set_decl_expr).ElseShift()).
		Finally(set_expr_assign)

	table_2 :=
		ForTypesThrough(lastType_t_).
		UseReductions(
			Map(let_t).To(set_context).ElseShift(),
			Map(id_t).To(set_id.Union(set_decl_expr, set_context)).ElseShift(),
			Map(integer_t).To(set_decl_expr.Union(set_context)).ElseShift(),
			Map(add_t).To(set_decl_expr.Union(set_context)).ElseShift(),
			Map(mul_t).To(set_decl_expr.Union(set_context)).ElseShift(),
			Map(in_t).To(set_expr_assign).ElseShift()).
		Finally(set_expr_assign.Union(set_context, set_expr_w_cxt))

	id_a := token.AddValue(idTok, "a")
	int_3 := token.AddValue(intTok, "3")

	tests := []struct {
		ReduceTable
		src    source.StaticSource
		stream []token.Token
		expect ast.Ast
	}{
		{
			table,
			MakeSource("test", "let a + 3 3"),
			[]token.Token{
				letTok.SetLineChar(1, 1),
				id_a.SetLineChar(1, 5),
				addTok.SetLineChar(1, 7),
				int_3.SetLineChar(1, 9),
				int_3.SetLineChar(1, 11),
			},
			mknode(assign_t,
				mknode(decl_t,
					ast.TokenNode(letTok),
					ast.TokenNode(id_a)),
				mknode(expr_t,
					ast.TokenNode(addTok),
					mknode(expr_t,
						ast.TokenNode(int_3)),
					mknode(expr_t,
						ast.TokenNode(int_3)))),
		},
		{
			table,
			MakeSource("test", "let a * 3 3"),
			[]token.Token{
				letTok.SetLineChar(1, 1),
				id_a.SetLineChar(1, 5),
				mulTok.SetLineChar(1, 7),
				int_3.SetLineChar(1, 9),
				int_3.SetLineChar(1, 11),
			},
			mknode(assign_t,
				mknode(decl_t,
					ast.TokenNode(letTok),
					ast.TokenNode(id_a)),
				mknode(expr_t,
					ast.TokenNode(mulTok),
					mknode(expr_t,
						ast.TokenNode(int_3)),
					mknode(expr_t,
						ast.TokenNode(int_3)))),
		},
		{
			table_2,
			MakeSource("test", "let a * 3 3 in * a a"),
			[]token.Token{
				letTok.SetLineChar(1, 1),
				id_a.SetLineChar(1, 5),
				mulTok.SetLineChar(1, 7),
				int_3.SetLineChar(1, 9),
				int_3.SetLineChar(1, 11),
				inTok.SetLineChar(1, 13),
				mulTok.SetLineChar(1, 16),
				id_a.SetLineChar(1, 18),
				id_a.SetLineChar(1, 20),
			},
			mknode(expr_t,
				mknode(context_t,
					mknode(assign_t,
						mknode(decl_t,
							ast.TokenNode(letTok),
							ast.TokenNode(id_a)),
						mknode(expr_t,
							ast.TokenNode(mulTok),
							mknode(expr_t,
								ast.TokenNode(int_3)),
							mknode(expr_t,
								ast.TokenNode(int_3)))),
					ast.TokenNode(inTok)),
				mknode(expr_t,
					ast.TokenNode(mulTok),
					mknode(expr_t,
						ast.TokenNode(id_a)),
					mknode(expr_t,
						ast.TokenNode(id_a)))),
		},
	}

	for i, test := range tests {
		p := New().
			Ruleset(test.ReduceTable).
			Load(test.stream, test.src, nil, nil).
			LogActions()
		actual := p.Parse()

		writeLog(p.FlushLog())

		if p.HasErrors() {
			es := p.GetErrors()
			for _, e := range es {
				fmt.Fprintf(os.Stderr, "%s\n", e.Error())
			}

			t.Fatal(testutil.TestFail2("errors", nil, p.errors, i))
		}

		if !actual.Equals(ast.AstRoot{test.expect}) {
			act := ast.GetOrderedString(actual)
			exp := ast.GetOrderedString(ast.AstRoot{test.expect})
			t.Fatal(testutil.TestFail(exp, act, i))
		}
	}
}