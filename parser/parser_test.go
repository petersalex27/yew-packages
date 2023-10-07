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
		From(let_t, id_t).Reduce(decl_fn),
	)

	set_expr := RuleSet(
		From(id_t).Reduce(expr_fn),                  // Id -> expr
		From(integer_t).Reduce(expr_fn),             // Id -> expr
		From(add_t, expr_t, expr_t).Reduce(expr_fn), // Add expr expr -> expr
		From(mul_t, expr_t, expr_t).Reduce(expr_fn), // Mul expr expr -> expr
	)

	set_id := RuleSet(From(let_t).Shift())

	set_assign := RuleSet(From(decl_t, expr_t).Reduce(assign_fn))
	set_decl_expr := set_decl.Union(set_expr)
	set_expr_assign := set_expr.Union(set_assign)
	set_expr_w_cxt := RuleSet(From(context_t, expr_t).Reduce(expr_fn))

	set_context := RuleSet(
		From(assign_t, in_t).Reduce(context_fn), // assign In -> context
	)

	my_class := integer_t

	table :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LA(let_t).Shift(),
				LA(id_t).Then(set_id.Union(set_decl_expr)).ElseShift(),
				LA(integer_t).Then(set_decl_expr).ElseShift(),
				LA(add_t).Then(set_decl_expr).ElseShift(),
				LA(mul_t).Then(set_decl_expr).ElseShift()).
			Finally(set_expr_assign)

	table_withClasses :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LA(let_t).Shift(),
				LA(id_t).Then(set_id.Union(set_decl_expr)).ElseShift(),
				LA(my_class).
					ForN(3, integer_t).Or(add_t).Or(mul_t).
						Then(set_decl_expr).
						ElseShift(),
			).
			Finally(set_expr_assign)

	table_2 :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LA(let_t).Then(set_context).ElseShift(),
				LA(id_t).Then(set_id.Union(set_decl_expr, set_context)).ElseShift(),
				LA(integer_t).Then(set_decl_expr.Union(set_context)).ElseShift(),
				LA(add_t).Then(set_decl_expr.Union(set_context)).ElseShift(),
				LA(mul_t).Then(set_decl_expr.Union(set_context)).ElseShift(),
				LA(in_t).Then(set_expr_assign).ElseShift()).
			Finally(set_expr_assign.Union(set_context, set_expr_w_cxt))
	
	
	table_2_withClasses :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LA(let_t).Then(set_context).ElseShift(),
				LA(id_t).Then(set_id.Union(set_decl_expr, set_context)).ElseShift(),
				LA(my_class).
					ForN(3, integer_t).Or(add_t).Or(mul_t).
						Then(set_decl_expr.Union(set_context)).
						ElseShift(),
				LA(in_t).Then(set_expr_assign).ElseShift()).
			Finally(set_expr_assign.Union(set_context, set_expr_w_cxt))

	id_a := token.AddValue(idTok, "a")
	int_3 := token.AddValue(intTok, "3")

	tests := []struct {
		table ReduceTable
		classTable ReduceTable
		src    source.StaticSource
		stream []token.Token
		expect ast.Ast
	}{
		{
			table,
			table_withClasses,
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
			table_withClasses,
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
			table_2_withClasses,
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

	// marks whether class or not (indexes corr. to index `j` in loop below)
	notOrClass := []string{"", "-c"}

	for i, test := range tests {
		for j, table := range []ReduceTable{test.table, test.classTable} {
			p := New().
				Ruleset(table).
				Load(test.stream, test.src, nil, nil).
				LogActions().
				StringType(test_token_stringType)

			actual := p.Parse()

			writeLog(p.FlushLog())

			if p.HasErrors() {
				es := p.GetErrors()
				for _, e := range es {
					fmt.Fprintf(os.Stderr, "%s\n", e.Error())
				}

				title := fmt.Sprintf("errors%s", notOrClass[j])
				t.Fatal(testutil.TestFail2(title, nil, p.errors, i, j))
			}

			if !actual.Equals(ast.AstRoot{test.expect}) {
				act := ast.GetOrderedString(actual)
				exp := ast.GetOrderedString(ast.AstRoot{test.expect})
				t.Fatal(testutil.TestFail(exp, act, i, j))
			}
		}
	}
}

func TestWhen(t *testing.T) {
	idTok := (test_token{}).SetType(uint(id_t))
	funcTok, _ := (test_token{}).SetType(uint(func_t)).SetValue("func")

	fn_fn := test_reduce_fn(fn_t).GiveName("function")
	assign_fn := test_reduce_fn(assign_t).GiveName("assignment")

	rules := RuleSet(
		fn_fn.When(func_t).From(id_t),
		assign_fn.From(func_t, fn_t),
	)

	table :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LA(func_t).Shift(),
				LA(id_t).Shift(),
			).Finally(rules)

	fTok := token.AddValue(idTok, "f")

	tests := []struct {
		table ReduceTable
		src    source.StaticSource
		stream []token.Token
		expect ast.Ast
	}{
		{
			table,
			MakeSource("test", "func f"),
			[]token.Token{
				funcTok.SetLineChar(1,1),
				fTok.SetLineChar(1,6),
			},
			ast.AstRoot{
				mknode(assign_t,
					ast.TokenNode(funcTok),
					mknode(fn_t, 
						ast.TokenNode(fTok.SetLineChar(1,6)))),
			},
		},
	}

	for i, test := range tests { 
		p := New().
			Ruleset(table).
			Load(test.stream, test.src, nil, nil).
			LogActions().
			StringType(test_token_stringType)

		actual := p.Parse()

		writeLog(p.FlushLog())

		if p.HasErrors() {
			es := p.GetErrors()
			for _, e := range es {
				fmt.Fprintf(os.Stderr, "%s\n", e.Error())
			}
			t.Fatal(testutil.TestFail2("errors", nil, p.errors, i))
		}

		if !actual.Equals(test.expect) {
			act := ast.GetOrderedString(actual)
			exp := ast.GetOrderedString(test.expect)
			t.Fatal(testutil.TestFail(exp, act, i))
		}
	}
}
