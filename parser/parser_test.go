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

	set_decl := Order(
		decl_fn.From(let_t, id_t),
	)

	set_expr := Order(
		expr_fn.From(id_t),                  // Id -> expr
		expr_fn.From(integer_t),             // Id -> expr
		expr_fn.From(add_t, expr_t, expr_t), // Add expr expr -> expr
		expr_fn.From(mul_t, expr_t, expr_t), // Mul expr expr -> expr
	)

	set_id := Order(Shift().When(let_t))

	set_assign := Order(assign_fn.From(decl_t, expr_t))
	set_decl_expr := Union(set_decl, set_expr)
	set_expr_assign := Union(set_expr, set_assign)
	set_expr_w_cxt := Order(expr_fn.From(context_t, expr_t))

	set_context := Order(
		context_fn.From(assign_t, in_t), // assign In -> context
	)

	my_class := integer_t

	table :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LookAhead(let_t).Shift(),
				LookAhead(id_t).Then(Union(set_id, set_decl_expr)).ElseShift(),
				LookAhead(integer_t).Then(set_decl_expr).ElseShift(),
				LookAhead(add_t).Then(set_decl_expr).ElseShift(),
				LookAhead(mul_t).Then(set_decl_expr).ElseShift()).
			Finally(set_expr_assign)

	table_withClasses :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LookAhead(let_t).Shift(),
				LookAhead(id_t).Then(Union(set_id, set_decl_expr)).ElseShift(),
				LookAhead(my_class).
					ForN(3, integer_t).Or(add_t).Or(mul_t).
					Then(set_decl_expr).
					ElseShift(),
			).
			Finally(set_expr_assign)

	table_2 :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LookAhead(let_t).Then(set_context).ElseShift(),
				LookAhead(id_t).Then(Union(set_id, set_decl_expr, set_context)).ElseShift(),
				LookAhead(integer_t).Then(Union(set_decl_expr, set_context)).ElseShift(),
				LookAhead(add_t).Then(Union(set_decl_expr, set_context)).ElseShift(),
				LookAhead(mul_t).Then(Union(set_decl_expr, set_context)).ElseShift(),
				LookAhead(in_t).Then(set_expr_assign).ElseShift()).
			Finally(Union(set_expr_assign, set_context, set_expr_w_cxt))

	table_2_withClasses :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LookAhead(let_t).Then(set_context).ElseShift(),
				LookAhead(id_t).Then(Union(set_id, set_decl_expr, set_context)).ElseShift(),
				LookAhead(my_class).
					ForN(3, integer_t).Or(add_t).Or(mul_t).
					Then(Union(set_decl_expr, set_context)).
					ElseShift(),
				LookAhead(in_t).Then(set_expr_assign).ElseShift()).
			Finally(Union(set_expr_assign, set_context, set_expr_w_cxt))

	id_a := token.AddValue(idTok, "a")
	int_3 := token.AddValue(intTok, "3")

	tests := []struct {
		table      ReductionTable
		classTable ReductionTable
		src        source.StaticSource
		stream     []token.Token
		expect     ast.Ast
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
		for j, table := range []ReductionTable{test.table, test.classTable} {
			p := NewParser().
				UsingReductionTable(table).
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

func BenchmarkParser(b *testing.B) {
	b.StopTimer()

	listTail_t := lastType_t_
	list_t := listTail_t + 1
	eq_t := list_t + 1
	last__t := eq_t + 1

	addTok, _ := (test_token{}).SetType(uint(add_t)).SetValue("+")
	mulTok, _ := (test_token{}).SetType(uint(mul_t)).SetValue("*")
	intTok := (test_token{}).SetType(uint(integer_t))
	idTok := (test_token{}).SetType(uint(id_t))
	letTok, _ := (test_token{}).SetType(uint(let_t)).SetValue("let")
	inTok, _ := (test_token{}).SetType(uint(in_t)).SetValue("in")
	openTok, _ := (test_token{}).SetType(uint(open_t)).SetValue("(")
	closeTok, _ := (test_token{}).SetType(uint(open_t)).SetValue(")")
	commaTok, _ := (test_token{}).SetType(uint(comma_t)).SetValue(",")
	eqTok, _ := (test_token{}).SetType(uint(eq_t)).SetValue("=")

	unenclose_fn := ProductionFunction(func(nodes ...ast.Ast) ast.Ast {
		return nodes[1]
	}).GiveName("unenclose")

	lt_fn := test_reduce_fn(listTail_t).GiveName("list-tail")
	list_fn := test_reduce_fn(list_t).GiveName("list")
	decl_fn := test_reduce_fn(decl_t).GiveName("declaration")
	expr_fn := test_reduce_fn(expr_t).GiveName("expression")
	assign_fn := test_reduce_fn(assign_t).GiveName("assignment")
	context_fn := test_reduce_fn(context_t).GiveName("context")

	set_decl := Order(
		decl_fn.From(let_t, id_t),
	)

	set_expr := Order(
		unenclose_fn.From(expr_t),
		expr_fn.From(id_t), // Id -> expr
		unenclose_fn.From(id_t),
		expr_fn.From(integer_t), // Id -> expr
		unenclose_fn.From(integer_t),
		expr_fn.From(add_t, expr_t, expr_t), // Add expr expr -> expr
		expr_fn.From(mul_t, expr_t, expr_t), // Mul expr expr -> expr
	)

	set_list := Order(
		lt_fn.From(expr_t, close_t),
		lt_fn.From(expr_t, comma_t, close_t),
		lt_fn.From(expr_t, comma_t, listTail_t),
	)

	set_list_build := Order(
		list_fn.From(open_t, listTail_t),
	)

	set_id := Order(Shift().When(let_t))

	set_assign := Order(assign_fn.From(decl_t, eq_t, expr_t))
	set_decl_expr := Union(set_decl, set_expr)
	set_expr_w_cxt := Order(expr_fn.From(context_t, expr_t))

	set_context := Order(
		context_fn.From(assign_t, in_t), // assign In -> context
	)

	table :=
		ForTypesThrough(last__t).
			UseReductions(
				LookAhead(let_t).Then(Union(set_context)).ElseShift(),
				LookAhead(id_t).Then(Union(set_id, set_expr, set_list_build, set_context)).ElseShift(),
				LookAhead(integer_t).Then(Union(set_expr, set_list_build, set_context)).ElseShift(),
				LookAhead(add_t).Then(Union(set_expr, set_list_build, set_context)).ElseShift(),
				LookAhead(mul_t).Then(Union(set_expr, set_list_build, set_context)).ElseShift(),
				LookAhead(comma_t).Then(Union(set_expr, set_assign, set_list_build, set_context, set_expr_w_cxt)),
				LookAhead(open_t).Then(Union(set_expr, set_assign, set_list_build, set_context, set_expr_w_cxt)),
				LookAhead(close_t).Then(Union(set_expr, set_assign, set_context, set_expr_w_cxt, set_list)),
				LookAhead(in_t).Then(Union(set_expr, set_list_build, set_assign)),
				LookAhead(eq_t).Then(set_decl_expr),
			).Finally(Union(set_expr, set_context, set_expr_w_cxt))

	id_a := token.AddValue(idTok, "a")
	int_3 := token.AddValue(intTok, "3")

	tests := []struct {
		src    source.StaticSource
		stream []token.Token
	}{
		{
			MakeSource("test", "let a = * 3 3 in * a a"),
			[]token.Token{
				letTok.SetLineChar(1, 1),
				id_a.SetLineChar(1, 5),
				eqTok.SetLineChar(1, 7),
				mulTok.SetLineChar(1, 9),
				int_3.SetLineChar(1, 11),
				int_3.SetLineChar(1, 13),
				inTok.SetLineChar(1, 16),
				mulTok.SetLineChar(1, 18),
				id_a.SetLineChar(1, 20),
				id_a.SetLineChar(1, 22),
			},
		},
		{
			MakeSource("test", "let a = * 3 3 in let a = * + a a a in + a * a a"),
			[]token.Token{
				letTok.SetLineChar(1, 1),
				eqTok,
				id_a.SetLineChar(1, 5),
				mulTok.SetLineChar(1, 7),
				int_3.SetLineChar(1, 9),
				int_3.SetLineChar(1, 11),
				inTok.SetLineChar(1, 13),
				letTok,
				id_a,
				eqTok,
				mulTok,
				addTok,
				id_a,
				id_a,
				id_a,
				inTok,
				addTok,
				id_a,
				mulTok,
				id_a,
				id_a,
			},
		},
		{
			MakeSource("test",
				"let a = (3,3) in "+
					"let a = (+ (3,3,a) a) in "+
					"+ (a) a"),
			[]token.Token{
				letTok, id_a, eqTok, // let a =
				openTok, int_3, commaTok, int_3, closeTok, // (3,3)
				inTok,               // in
				letTok, id_a, eqTok, // let id =
				openTok,                                         // (
				addTok,                                          // +
				openTok, int_3, commaTok, int_3, id_a, closeTok, // (3,3,a)
				id_a,                    // a
				closeTok,                // )
				inTok,                   // in
				addTok,                  // +
				openTok, id_a, closeTok, // (a)
				id_a, // a
			},
		},
		{
			MakeSource("test",
				"let a = (3,3) in "+
					"let a = (+ (3,3,a) a) in "+
					"let a = ((((a,a,),a),a,a,(a,a,a),a),a) in + (a) a"),
			[]token.Token{
				letTok, id_a, eqTok, // let a =
				openTok, int_3, commaTok, int_3, closeTok, // (3,3)
				inTok,               // in
				letTok, id_a, eqTok, // let id =
				openTok,                                         // (
				addTok,                                          // +
				openTok, int_3, commaTok, int_3, id_a, closeTok, // (3,3,a)
				id_a,     // a
				closeTok, // )
				inTok,    // in
				openTok,
				openTok,
				openTok,
				openTok, id_a, commaTok, id_a, commaTok, closeTok,
				commaTok, id_a,
				closeTok,
				commaTok, id_a, commaTok, id_a, commaTok,
				openTok, id_a, commaTok, id_a, commaTok, id_a, closeTok,
				commaTok, id_a,
				closeTok,
				commaTok, id_a,
				closeTok,
				inTok,
				addTok,                  // +
				openTok, id_a, closeTok, // (a)
				id_a, // a
			},
		},
	}

	test := tests[3]
	for _, truthy := range []bool{true, false} {
		b.ResetTimer()
		b.Log("// optimizedReduce =", truthy, "/////////////")
		for i := 0; i < 10000; i++ {
			p := NewParser().
				UsingReductionTable(table).
				Load(test.stream, test.src, nil, nil).
				Benchmarker()

			p.ClearFlags()
			p.SetOptimizedReduce(truthy)

			b.StartTimer()
			_ = p.Parse()
			b.StopTimer()
		}
		b.Log(b.Elapsed())
	}

}

func TestWhen(t *testing.T) {
	idTok := (test_token{}).SetType(uint(id_t))
	funcTok, _ := (test_token{}).SetType(uint(func_t)).SetValue("func")

	fn_fn := test_reduce_fn(fn_t).GiveName("function")
	assign_fn := test_reduce_fn(assign_t).GiveName("assignment")

	rules := Order(
		fn_fn.When(func_t).From(id_t),
		assign_fn.From(func_t, fn_t),
	)

	table :=
		ForTypesThrough(lastType_t_).
			UseReductions(
				LookAhead(func_t).Shift(),
				LookAhead(id_t).Shift(),
			).Finally(rules)

	fTok := token.AddValue(idTok, "f")

	tests := []struct {
		table  ReductionTable
		src    source.StaticSource
		stream []token.Token
		expect ast.Ast
	}{
		{
			table,
			MakeSource("test", "func f"),
			[]token.Token{
				funcTok.SetLineChar(1, 1),
				fTok.SetLineChar(1, 6),
			},
			ast.AstRoot{
				mknode(assign_t,
					ast.TokenNode(funcTok),
					mknode(fn_t,
						ast.TokenNode(fTok.SetLineChar(1, 6)))),
			},
		},
	}

	for i, test := range tests {
		p := NewParser().
			UsingReductionTable(table).
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
