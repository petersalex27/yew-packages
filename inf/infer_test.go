package inf

import (
	"testing"

	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
)

// tests variable inference rule
func TestVar(t *testing.T) {
	var v0 types.Variable[nameable.Testable]

	{
		// block prevents accidental use of cxt
		cxt := NewTestableContext()
		v0 = cxt.typeContext.NewVar()
	}

	xName := nameable.MakeTestable("x")
	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")
	nName := nameable.MakeTestable("n")
	uintName := nameable.MakeTestable("Uint")

	x := expr.Const[nameable.Testable]{Name: xName}
	Array := types.MakeConst(arrName)                   // Array
	a := types.Var(aName)                               // a
	Array_a := types.Apply[nameable.Testable](Array, a) // Array a
	n := expr.Var(nName)                                // n
	Uint := types.MakeConst(uintName)                   // Uint
	n_Uint := types.Judgement(expr.Expression[nameable.Testable](n), types.Type[nameable.Testable](Uint))
	var_n_Uint := types.Judgement[nameable.Testable, expr.Variable[nameable.Testable]](n, Uint)
	domain := []types.ExpressionJudgement[nameable.Testable, expr.Expression[nameable.Testable]]{n_Uint}
	Array_a_n := types.Index(Array_a, domain...)       // (Array a; n)
	mapval_n_Uint__Array_a := types.MakeDependentType( // mapval (n: Uint) . (Array a)
		[]types.TypeJudgement[nameable.Testable, expr.Variable[nameable.Testable]]{var_n_Uint},
		Array_a,
	)
	Array_v0 := types.Apply[nameable.Testable](Array, v0)
	Array_v0_n := types.Index(Array_v0, domain...)

	tests := []struct {
		description string
		input       bridge.JudgementAsExpression[nameable.Testable, expr.Const[nameable.Testable]]
		expect      bridge.JudgementAsExpression[nameable.Testable, expr.Const[nameable.Testable]]
	}{
		{
			"x: Array => x: Array",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array),
		},
		{
			"x: a => x: a",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, a),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, a),
		},
		{
			"x: Array a => x: Array a",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array_a),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array_a),
		},
		{
			"x: forall a . a => x: $0",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, types.Forall(a).Bind(a)),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, v0),
		},
		{
			"x: forall a . Array a => x: Array $0",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, types.Forall(a).Bind(Array_a)),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array_v0),
		},
		{
			"x: (Array a; n) => x: (Array a; n)",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array_a_n),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array_a_n),
		},
		{
			"x: mapval (n: Uint) . (Array a) => x: (Array a; n)",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, mapval_n_Uint__Array_a),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array_a_n),
		},
		{
			"x: forall a . mapval (n: Uint) . (Array a) => x: (Array $0; n)",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, types.Forall(a).Bind(mapval_n_Uint__Array_a)),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array_v0_n),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.Var(test.input)

		eq := types.JudgementEquals[nameable.Testable, expr.Const[nameable.Testable], types.Type[nameable.Testable]](
			actual.ToTypeJudgement(),
			test.expect.ToTypeJudgement(),
		)
		if !eq {
			t.Fatal(
				testutil.
					Testing("equality", test.description).
					FailMessage(test.expect, actual, i))
		}
	}
}

func TestApp(t *testing.T) {
	var v0 types.Variable[nameable.Testable]
	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))

	{
		// block prevents accidental use of cxt
		cxt := NewTestableContext()
		v0 = cxt.typeContext.NewVar()
	}

	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("x")
	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")

	x := expr.Const[nameable.Testable]{Name: xName}
	y := expr.Const[nameable.Testable]{Name: yName}
	Array := types.MakeConst(arrName) // Array
	a := types.Var(aName)             // a

	tests := []struct {
		description string
		input0      bridge.JudgementAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		input1      bridge.JudgementAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		findIn      types.Variable[nameable.Testable]
		findOut     types.Type[nameable.Testable]
		expect      bridge.JudgementAsExpression[nameable.Testable, expr.Application[nameable.Testable]]
	}{
		{
			"(x: a) (y: Array) => (x y): $0",
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](x, a),
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](y, Array),
			a, types.Apply[nameable.Testable](arrow, Array, v0), // a = Array -> $0
			bridge.Judgement[nameable.Testable, expr.Application[nameable.Testable]](
				expr.Apply[nameable.Testable](x, y),
				v0,
			),
		},
		{
			"(x: a) (x: a) => (x x): $0",
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](x, a),
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](x, a),
			a, types.Apply[nameable.Testable](arrow, a, v0), // a = a -> $0
			bridge.Judgement[nameable.Testable, expr.Application[nameable.Testable]](
				expr.Apply[nameable.Testable](x, x),
				v0,
			),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.App(test.input0, test.input1)

		eq := types.JudgementEquals[nameable.Testable, expr.Application[nameable.Testable], types.Type[nameable.Testable]](
			actual.ToTypeJudgement(),
			test.expect.ToTypeJudgement(),
		)
		if !eq {
			t.Fatal(
				testutil.
					Testing("equality", test.description).
					FailMessage(test.expect, actual, i))
		}

		findOutActual := cxt.typeContext.Find(test.findIn)
		eq = test.findOut.Equals(findOutActual)
		if !eq {
			t.Fatal(
				testutil.
					Testing("find", test.description).
					FailMessage(test.findOut, findOutActual, i))
		}
	}
}

func TestAbs(t *testing.T) {
	var v0 types.Variable[nameable.Testable]
	var ve0 expr.Variable[nameable.Testable]
	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))

	{
		// block prevents accidental use of cxt
		cxt := NewTestableContext()
		v0 = cxt.typeContext.NewVar()
		ve0 = cxt.exprContext.NewVar()
	}

	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("y")
	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")

	x := expr.Const[nameable.Testable]{Name: xName}
	y := expr.Const[nameable.Testable]{Name: yName}
	Array := types.MakeConst(arrName) // Array
	a := types.Var(aName)             // a
	Array_a := types.Apply[nameable.Testable](Array, a) // Array a

	tests := []struct {
		description string
		inputParam  nameable.Testable
		inputExpr   bridge.JudgementAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		expect      bridge.JudgementAsExpression[nameable.Testable, expr.Function[nameable.Testable]]
	}{
		{
			`x => y: Array => (\$0 -> y): $0 -> Array`,
			xName,
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](y, Array),
			bridge.Judgement[nameable.Testable, expr.Function[nameable.Testable]](
				expr.Bind[nameable.Testable](ve0).In(y),
				types.Apply[nameable.Testable](arrow, v0, Array),
			),
		},
		{
			`x => (x y): Array => (\$0 -> $0 y): $0 -> Array`,
			xName,
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](expr.Apply[nameable.Testable](x, y), Array),
			bridge.Judgement[nameable.Testable, expr.Function[nameable.Testable]](
				expr.Bind[nameable.Testable](ve0).In(expr.Apply[nameable.Testable](ve0, y)),
				types.Apply[nameable.Testable](arrow, v0, Array),
			),
		},
		{
			`x => (x y): a => (\$0 -> $0 y): $0 -> a`,
			xName,
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](expr.Apply[nameable.Testable](x, y), a),
			bridge.Judgement[nameable.Testable, expr.Function[nameable.Testable]](
				expr.Bind[nameable.Testable](ve0).In(expr.Apply[nameable.Testable](ve0, y)),
				types.Apply[nameable.Testable](arrow, v0, a),
			),
		},
		{
			`x => (x y): Array a => (\$0 -> $0 y): $0 -> Array a`,
			xName,
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](expr.Apply[nameable.Testable](x, y), Array_a),
			bridge.Judgement[nameable.Testable, expr.Function[nameable.Testable]](
				expr.Bind[nameable.Testable](ve0).In(expr.Apply[nameable.Testable](ve0, y)),
				types.Apply[nameable.Testable](arrow, v0, Array_a),
			),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.Abs(test.inputParam)(test.inputExpr)

		eq := types.JudgementEquals[nameable.Testable, expr.Function[nameable.Testable], types.Type[nameable.Testable]](
			actual.ToTypeJudgement(),
			test.expect.ToTypeJudgement(),
		)
		if !eq {
			t.Fatal(
				testutil.
					Testing("equality", test.description).
					FailMessage(test.expect, actual, i))
		}
	}
}
/*
func TestLet(t *testing.T) {
	var v0 types.Variable[nameable.Testable]
	var ve0 expr.Variable[nameable.Testable]
	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))

	{
		// block prevents accidental use of cxt
		cxt := NewTestableContext()
		v0 = cxt.typeContext.NewVar()
		ve0 = cxt.exprContext.NewVar()
	}

	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("y")
	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")

	x := expr.Const[nameable.Testable]{Name: xName}
	y := expr.Const[nameable.Testable]{Name: yName}
	Array := types.MakeConst(arrName) // Array
	a := types.Var(aName)             // a
	Array_a := types.Apply[nameable.Testable](Array, a) // Array a

	tests := []struct {
		description string
		inputParam  nameable.Testable
		inputAssign bridge.JudgementAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		inputExpr   bridge.JudgementAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		expect      bridge.JudgementAsExpression[nameable.Testable, expr.NameContext[nameable.Testable]]
	}{
		{
			`x, y: Array => x: Array => let x = y in x: Array`,
			xName,
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](y, Array),
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](x, Array),
			bridge.Judgement[nameable.Testable, expr.NameContext[nameable.Testable]](
				expr.Bind[nameable.Testable](ve0).In(y),
				types.Apply[nameable.Testable](arrow, v0, Array),
			),
		},
		{
			`x => (x y): Array => (\$0 -> $0 y): $0 -> Array`,
			xName,
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](expr.Apply[nameable.Testable](x, y), Array),
			bridge.Judgement[nameable.Testable, expr.Function[nameable.Testable]](
				expr.Bind[nameable.Testable](ve0).In(expr.Apply[nameable.Testable](ve0, y)),
				types.Apply[nameable.Testable](arrow, v0, Array),
			),
		},
		{
			`x => (x y): a => (\$0 -> $0 y): $0 -> a`,
			xName,
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](expr.Apply[nameable.Testable](x, y), a),
			bridge.Judgement[nameable.Testable, expr.Function[nameable.Testable]](
				expr.Bind[nameable.Testable](ve0).In(expr.Apply[nameable.Testable](ve0, y)),
				types.Apply[nameable.Testable](arrow, v0, a),
			),
		},
		{
			`x => (x y): Array a => (\$0 -> $0 y): $0 -> Array a`,
			xName,
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](expr.Apply[nameable.Testable](x, y), Array_a),
			bridge.Judgement[nameable.Testable, expr.Function[nameable.Testable]](
				expr.Bind[nameable.Testable](ve0).In(expr.Apply[nameable.Testable](ve0, y)),
				types.Apply[nameable.Testable](arrow, v0, Array_a),
			),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.Abs(test.inputParam)(test.inputExpr)

		eq := types.JudgementEquals[nameable.Testable, expr.Function[nameable.Testable], types.Type[nameable.Testable]](
			actual.ToTypeJudgement(),
			test.expect.ToTypeJudgement(),
		)
		if !eq {
			t.Fatal(
				testutil.
					Testing("equality", test.description).
					FailMessage(test.expect, actual, i))
		}
	}
}
*/
// some integration tests
func TestProofValidation(t *testing.T) {

}