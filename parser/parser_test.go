package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
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

//  - check that precondition is run only when matched in the first place
//  - check that if precondition is true, then precondition's `call` calls 
//		wrapped production's `call`
//	- check that correct stat is returned
//
// think of test cases for expectApplied like a truth table:
//	----------------------------------------------
//	|    match    |	precondition | expectApplied |
//	|-------------|--------------|---------------|
//	|    true     |    true      |     true      |
//	|    true     |    false     |     false     | 
//	|    false    |    true      |     false     |
//	|    false    |    false     |     false     |
//	----------------------------------------------
func TestParameterized(t *testing.T) {
	divTok, _ := (test_token{}).SetType(uint(div_t)).SetValue("/")
	intTok := (test_token{}).SetType(uint(integer_t))
	int0, _ := intTok.SetValue("0")
	int1, _ := intTok.SetValue("1")

	error_fn := test_reduce_fn(error_t).GiveName("error")

	ruleDiv := Precondition(
		error_fn.From(div_t, expr_t, integer_t),
		func(nodes ...ast.Ast) bool {
			const _, _, denominatorIndex int = 0, 1, 2
			denominatorNode := nodes[denominatorIndex]
			denominator, ok := denominatorNode.(ast.Token)
			if !ok {
				panic("failed type assertion")
			}

			return denominator.GetValue() == "0"
		},
	)


	productions := Order(ruleDiv)
     
	tests := []struct {
		desc   string
		src source.StaticSource
		input []ast.Ast
		expectStatus status.Status
		expectApplied bool
		expectStackRaw []ast.Ast
	}{
		{ // true true
			"match true, condition true: / 0 0",
			MakeSource("test/parser/precondition-rule", "/ 0 0"),
			[]ast.Ast{ast.TokenNode(divTok), mknode(expr_t, ast.TokenNode(int0)), ast.TokenNode(int0)},
			status.Ok, true,
			[]ast.Ast{
				mknode(error_t, 
					ast.TokenNode(divTok), 
					mknode(expr_t, ast.TokenNode(int0)), 
					ast.TokenNode(int0),
				),
			},
		},
		{ // true false
			"match true, condition false: / 0 1",
			MakeSource("test/parser/precondition-rule", "/ 0 1"),
			[]ast.Ast{ast.TokenNode(divTok), mknode(expr_t, ast.TokenNode(int0)), ast.TokenNode(int1)},
			status.Ok, false,
			[]ast.Ast{ast.TokenNode(divTok), mknode(expr_t, ast.TokenNode(int0)), ast.TokenNode(int1)},
		},
		{ // false true (condition is tech. never checked, but would pass if it were)
			"match false, condition true: 0 0 0",
			MakeSource("test/parser/precondition-rule", "0 0 0"),
			[]ast.Ast{ast.TokenNode(int0), mknode(expr_t, ast.TokenNode(int0)), ast.TokenNode(int0)},
			status.EndAction, false,
			[]ast.Ast{ast.TokenNode(int0), mknode(expr_t, ast.TokenNode(int0)), ast.TokenNode(int0)},
		},
		{ // false false (condition is tech. never checked, but would pass if it were)
			"match false, condition false: 0 0 1",
			MakeSource("test/parser/precondition-rule", "0 0 1"),
			[]ast.Ast{ast.TokenNode(int0), mknode(expr_t, ast.TokenNode(int0)), ast.TokenNode(int1)},
			status.EndAction, false,
			[]ast.Ast{ast.TokenNode(int0), mknode(expr_t, ast.TokenNode(int0)), ast.TokenNode(int1)},
		},
	}

	for i, test := range tests {
		p := NewParser().
			LA(1).
			UsingReductionTable(ReductionTable{}). // empty reduction table
			Load([]token.Token{}, test.src, nil, nil).
			InitialStackPush(test.input...)
		
		actualStat, actualApplied := p.reduce(productions)
		
		// check for errors
		if p.HasErrors() {
			es := p.GetErrors()
			for _, e := range es {
				fmt.Fprintf(os.Stderr, "%s\n", e.Error())
			}
			t.Fatalf(
				testutil.Testing("errors", test.desc).FailMessage(nil, es, i))
		}

		// check results =======

		// stat
		if actualStat != test.expectStatus {
			t.Fatalf(testutil.Testing("stat", test.desc).
				FailMessage(test.expectStatus, actualStat, i))
		}

		// applied
		if actualApplied != test.expectApplied {
			t.Fatalf(testutil.Testing("applied", test.desc).
				FailMessage(test.expectApplied, actualApplied, i))
		}

		// stack len
		actual, _ := p.stack.MultiCheck(int(p.stack.GetCount()))
		if len(test.expectStackRaw) != len(actual) {
			t.Fatalf(testutil.Testing("stack length", test.desc).
				FailMessage(test.expectStackRaw, actual, i))
		}

		// stack elems
		for j, elem := range test.expectStackRaw {
			if !elem.Equals(actual[j]) {
				t.Fatalf(testutil.Testing("stack elem", test.desc).
					FailMessage(test.expectStackRaw, actual, i, j))
			}
		}
	}
}
