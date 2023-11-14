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
	var ve0 expr.Variable[nameable.Testable]

	{
		// block prevents accidental use of cxt
		cxt := NewTestableContext()
		v0 = cxt.typeContext.NewVar()
		ve0 = cxt.exprContext.NewVar()
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
	n_Uint := types.Judgement(expr.Referable[nameable.Testable](n), types.Type[nameable.Testable](Uint))
	ve0_Uint := types.Judgement(expr.Referable[nameable.Testable](ve0), types.Type[nameable.Testable](Uint))
	var_n_Uint := types.Judgement[nameable.Testable, expr.Variable[nameable.Testable]](n, Uint)
	domain := []types.ExpressionJudgement[nameable.Testable, expr.Referable[nameable.Testable]]{n_Uint}
	domain2 := []types.ExpressionJudgement[nameable.Testable, expr.Referable[nameable.Testable]]{ve0_Uint}
	Array_a_n := types.Index(Array_a, domain...)       // (Array a; n)
	Array_a_ve0 := types.Index(Array_a, domain2...)    // (Array a; x0)
	mapval_n_Uint__Array_a := types.MakeDependentType( // mapval (n: Uint) . (Array a)
		[]types.TypeJudgement[nameable.Testable, expr.Variable[nameable.Testable]]{var_n_Uint},
		Array_a,
	)
	Array_v0 := types.Apply[nameable.Testable](Array, v0)
	Array_v0_ve0 := types.Index(Array_v0, domain2...)

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
			"x: mapval (n: Uint) . (Array a) => x: (Array a; $0)",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, mapval_n_Uint__Array_a),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array_a_ve0),
		},
		{
			"x: forall a . mapval (n: Uint) . (Array a) => x: (Array $0; $e0)",
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, types.Forall(a).Bind(mapval_n_Uint__Array_a)),
			bridge.Judgement[nameable.Testable, expr.Const[nameable.Testable]](x, Array_v0_ve0),
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
	bName := nameable.MakeTestable("b")

	x := expr.Const[nameable.Testable]{Name: xName}
	y := expr.Const[nameable.Testable]{Name: yName}
	Array := types.MakeConst(arrName) // Array
	a := types.Var(aName)             // a
	b := types.Var(bName)             // b

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
			"(y: b) (x: a) => (y x): $0",
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](y, b),
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](x, a),
			b, types.Apply[nameable.Testable](arrow, a, v0), // b = a -> $0
			bridge.Judgement[nameable.Testable, expr.Application[nameable.Testable]](
				expr.Apply[nameable.Testable](y, x),
				v0,
			),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.App(test.input0, test.input1)

		if cxt.HasErrors() {
			t.Fatal(
				testutil.
					Testing("errors", test.description).
					FailMessage(nil, cxt.GetReports(), i))
		}

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

		findOutActual := cxt.Find(test.findIn)
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
	Array := types.MakeConst(arrName)                   // Array
	a := types.Var(aName)                               // a
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

func TestGen(t *testing.T) {
	// var v0 types.Variable[nameable.Testable]

	// {
	// 	// block prevents accidental use of cxt
	// 	cxt := NewTestableContext()
	// 	v0 = cxt.typeContext.NewVar()
	// }

	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")
	nName := nameable.MakeTestable("n")
	uintName := nameable.MakeTestable("Uint")

	Array := types.MakeConst(arrName)                   // Array
	a := types.Var(aName)                               // a
	Array_a := types.Apply[nameable.Testable](Array, a) // Array a
	n := expr.Var(nName)                                // n
	Uint := types.MakeConst(uintName)                   // Uint
	n_Uint := types.Judgement(expr.Referable[nameable.Testable](n), types.Type[nameable.Testable](Uint))
	var_n_Uint := types.Judgement[nameable.Testable, expr.Variable[nameable.Testable]](n, Uint)
	domain := []types.ExpressionJudgement[nameable.Testable, expr.Referable[nameable.Testable]]{n_Uint}
	vs := []types.TypeJudgement[nameable.Testable, expr.Variable[nameable.Testable]]{var_n_Uint}
	Array_a_n := types.Index(Array_a, domain...) // (Array a; n)
	Array_n := types.Index(types.Apply[nameable.Testable](Array), domain...)

	tests := []struct {
		description string
		in          types.Monotyped[nameable.Testable]
		expect      types.Polytype[nameable.Testable]
	}{
		{
			"Array => forall _ . Array",
			Array,
			types.Forall[nameable.Testable]().Bind(Array),
		},
		{
			"a => forall a . a",
			a,
			types.Forall(a).Bind(a),
		},
		{
			"Array a => forall a . Array a",
			Array_a,
			types.Forall(a).Bind(Array_a),
		},
		{
			"Array; n => forall _ . mapval (n: Uint) . Array",
			Array_n,
			types.Forall[nameable.Testable]().Bind(types.MakeDependentType[nameable.Testable](vs, types.Apply[nameable.Testable](Array))),
		},
		{
			"Array a; n => forall a . mapval (n: Uint) . Array a",
			Array_a_n,
			types.Forall(a).Bind(types.MakeDependentType[nameable.Testable](vs, Array_a)),
		},
	}

	for i, test := range tests {
		cxt := NewContext[nameable.Testable]()
		actual := cxt.Gen(test.in)
		if !actual.Equals(test.expect) {
			t.Fatal(
				testutil.Testing("equality", test.description).
					FailMessage(test.expect, actual, i),
			)
		}
	}
}

func TestLet(t *testing.T) {
	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))
	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("y")
	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")
	zeroName := nameable.MakeTestable("0")
	intName := nameable.MakeTestable("Int")

	x := expr.Const[nameable.Testable]{Name: xName}       // x (constant)
	y := expr.Const[nameable.Testable]{Name: yName}       // y (constant)
	zero := expr.Const[nameable.Testable]{Name: zeroName} // 0 (constant)
	yVar := expr.Var(yName)                               // y (variable)
	idFunc := expr.Bind[nameable.Testable](yVar).In(yVar) // (\y -> y)
	Int := types.MakeConst(intName)                       // Int
	Array := types.MakeConst(arrName)                     // Array
	a := types.Var(aName)                                 // a
	aToA := types.Apply[nameable.Testable](arrow, a, a)   // a -> a
	x_0 := expr.Apply[nameable.Testable](x, zero)         // (x 0)

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
				expr.Let[nameable.Testable](x, y, x),
				Array,
			),
		},
		{
			`x, (\y -> y): a -> a => (x 0): Int => let x = (\y -> y) in x 0: Int`,
			xName,
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](idFunc, aToA),
			bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](x_0, Int),
			bridge.Judgement[nameable.Testable, expr.NameContext[nameable.Testable]](
				expr.Let[nameable.Testable](x, idFunc, x_0),
				Int,
			),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.Let(test.inputParam, test.inputAssign)(test.inputExpr)

		eq := types.JudgementEquals[nameable.Testable, expr.NameContext[nameable.Testable], types.Type[nameable.Testable]](
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

func TestFind(t *testing.T) {
	intName := nameable.MakeTestable("Int")
	myTypeName := nameable.MakeTestable("MyType")
	aName, bName := nameable.MakeTestable("a"), nameable.MakeTestable("b")

	Int := types.MakeConst(intName)
	a, b := types.Var(aName), types.Var(bName)
	MyType := types.MakeConst(myTypeName)
	MyType_a := types.Apply[nameable.Testable](MyType, a)
	MyType_b := types.Apply[nameable.Testable](MyType, b)

	tests := []struct {
		desc string
		targ types.Variable[nameable.Testable]
		sub  types.Monotyped[nameable.Testable]
	}{
		{
			"a = Int",
			a, Int,
		},
		{
			"a = b",
			a, b,
		},
		{
			"a = MyType b",
			a, MyType_b,
		},
		{
			"a = MyType a",
			a, MyType_a,
		},
	}

	for i, test := range tests {
		expectBefore, expectAfter := test.targ, test.sub
		cxt := NewContext[nameable.Testable]()

		// test find before substitution added
		beforeSub := cxt.Find(test.targ)
		if !beforeSub.Equals(expectBefore) {
			t.Fatal(
				testutil.
					Testing("find before sub. added", test.desc).
					FailMessage(expectBefore, beforeSub, i))
		}

		// now add substitution
		cxt.typeSubs.Add(test.targ.GetReferred(), test.sub)

		// test find after substitution added
		afterSub := cxt.Find(test.targ)
		if !afterSub.Equals(expectAfter) {
			t.Fatal(
				testutil.
					Testing("find after sub. added", test.desc).
					FailMessage(expectAfter, afterSub, i))
		}
	}
}

func TestUnify(t *testing.T) {
	type expected struct {
		inTable bool
		in, out types.Monotyped[nameable.Testable]
	}
	intName := nameable.MakeTestable("Int")
	myTypeName := nameable.MakeTestable("MyType")
	myOtherTypeName := nameable.MakeTestable("MyOtherType")
	aName, bName := nameable.MakeTestable("a"), nameable.MakeTestable("b")

	Int := types.MakeConst(intName)
	a, b := types.Var(aName), types.Var(bName)
	MyType := types.MakeConst(myTypeName)
	MyOtherType := types.MakeConst(myOtherTypeName)
	MyType_a_b := types.Apply[nameable.Testable](MyType, a, b)
	MyOtherType_a := types.Apply[nameable.Testable](MyOtherType, a)
	MyOtherType_b := types.Apply[nameable.Testable](MyOtherType, b)
	MyType_a := types.Apply[nameable.Testable](MyType, a)
	MyType_b := types.Apply[nameable.Testable](MyType, b)

	tests := []struct {
		desc        string
		left, right types.Monotyped[nameable.Testable]
		expectStat  Status
		expect      []expected
	}{
		{
			"Unify(a, b)",
			a, b,
			Ok,
			[]expected{
				{true, a, b},
				{false, b, b},
			},
		},
		{
			"Unify(a, Int)",
			a, Int,
			Ok,
			[]expected{
				{true, a, Int},
				{false, Int, Int},
			},
		},
		{
			"Unify(Int, Int)",
			Int, Int,
			Ok,
			[]expected{
				{false, Int, Int},
			},
		},
		{
			"Unify(a, MyType b)",
			a, MyType_b,
			Ok,
			[]expected{
				{true, a, MyType_b},
				{false, MyType_b, MyType_b},
			},
		},
		{
			"Unify(MyType b, MyType b)",
			MyType_b, MyType_b,
			Ok,
			[]expected{
				{false, MyType, MyType},
				{true, b, b},
			},
		},
		{
			"Unify(MyType a, MyType b)",
			MyType_a, MyType_b,
			Ok,
			[]expected{
				{false, MyType, MyType},
				{true, a, b},
				{false, b, b},
			},
		},
		{
			"Unify(MyOtherType b, MyType b)",
			MyOtherType_b, MyType_b,
			ConstantMismatch,
			[]expected{
				{false, MyOtherType, MyOtherType},
				{false, MyType, MyType},
				{false, b, b},
			},
		},
		{
			"Unify(MyOtherType a, MyType b)",
			MyOtherType_a, MyType_b,
			ConstantMismatch,
			[]expected{
				{false, MyOtherType, MyOtherType},
				{false, MyType, MyType},
				{false, a, a},
				{false, b, b},
			},
		},
		{
			"Unify(MyType a b, MyType b)",
			MyType_a_b, MyType_b,
			ParamLengthMismatch,
			[]expected{
				{false, MyType, MyType},
				{false, a, a},
				{false, b, b},
			},
		},
		{
			"Unify(a, MyType a)",
			a, MyType_a,
			OccursCheckFailed,
			[]expected{
				{false, MyType, MyType},
				{false, a, a},
			},
		},
		{
			"Unify(a, MyType a b)",
			a, MyType_a_b,
			OccursCheckFailed,
			[]expected{
				{false, MyType, MyType},
				{false, a, a},
				{false, b, b},
			},
		},
		{
			"Unify(b, MyType a b)",
			b, MyType_a_b,
			OccursCheckFailed,
			[]expected{
				{false, MyType, MyType},
				{false, a, a},
				{false, b, b},
			},
		},
	}

	for i, test := range tests {
		cxt := NewContext[nameable.Testable]()

		stat := cxt.Unify(test.left, test.right)
		if stat != test.expectStat {
			t.Fatal(
				testutil.
					Testing("stat", test.desc).
					FailMessage(test.expectStat, stat, i))
		}

		for j, expect := range test.expect {
			// check if expected value for whether in sub. table
			_, inTable := cxt.typeSubs.Get(expect.in.GetReferred())
			if inTable != expect.inTable {
				t.Fatal(
					testutil.
						Testing("found in sub. table", test.desc).
						FailMessage(expect.inTable, inTable, i, j))
			}

			// check if expected result for find
			out := cxt.Find(expect.in)
			if !out.Equals(expect.out) {
				t.Fatal(
					testutil.
						Testing("find return value", test.desc).
						FailMessage(expect.out, out, i, j))
			}
		}
	}
}

// some integration tests
func TestProofValidation(t *testing.T) {
	// prove:
	//	let x = (\y -> y) in x 0: Int
	//
	// full proof:
	//	𝚪 = {0: Int, (λy.y): a -> a}:
	//
	//		                           [ x: forall a. a -> a ]¹    Inst(forall a. a -> a)
	//		                           -------------------------------------------------- [Var]
	//		                   0: Int                      x: v -> v                       t0, Int = v
	//		                   ----------------------------------------------------------------------- [App]
	//		                                               x 0: t0
	//		                                               ------- [Id]
	//		  (λy.y): a -> a                               x 0: Int
	//		1 ----------------------------------------------------- [Let]
	//		               let x = (λy.y) in x 0: Int

	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))
	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("y")
	aName := nameable.MakeTestable("a")
	zeroName := nameable.MakeTestable("0")
	intName := nameable.MakeTestable("Int")

	x := expr.Const[nameable.Testable]{Name: xName}       // x (constant)
	zero := expr.Const[nameable.Testable]{Name: zeroName} // 0 (constant)
	yVar := expr.Var(yName)                               // y (variable)
	idFunc := expr.Bind[nameable.Testable](yVar).In(yVar) // (\y -> y)
	Int := types.MakeConst(intName)                       // Int
	a := types.Var(aName)                                 // a
	aToA := types.Apply[nameable.Testable](arrow, a, a)   // a -> a
	x_0 := expr.Apply[nameable.Testable](x, zero)         // (x 0)

	letExpr := expr.Let[nameable.Testable](x, idFunc, x_0)

	cxt := NewTestableContext()

	Γ := struct {
		id, zero bridge.JudgementAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
	}{
		// (λy.y): a -> a
		id: bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](idFunc, aToA),
		// 0: Int
		zero: bridge.Judgement[nameable.Testable, expr.Expression[nameable.Testable]](zero, Int),
	}

	// step 1: add context (i.e., `x: Gen(a -> a)`) and first premise for let expression
	step := 0
	discharge_x_assumption := cxt.Let(xName, Γ.id) // returns function that discharges assumption

	// step 2: get assumption for first premise of Var
	step++
	x_assumption, found := cxt.Get(x)
	if !found {
		t.Fatal(testutil.Testing("x assumption get").FailMessage(true, found, step))
	}

	// step 3: do Var rule
	step++
	var_x := cxt.Var(x_assumption)

	// step 4: do App rule
	step++
	app_x_0_conclusion := cxt.App(bridge.Judgement(var_x.GetExpressionAndType()), Γ.zero)

	if cxt.HasErrors() {
		t.Fatal(testutil.Testing("app rule errors").FailMessage(nil, cxt.GetReports(), step))
	}

	{
		actualExpr, actualType := app_x_0_conclusion.GetExpressionAndType()
		expectExpr, expectType := x_0, Int

		if !expectExpr.StrictEquals(actualExpr) {
			t.Fatal(testutil.Testing("app rule expression result").FailMessage(expectExpr, actualExpr, step))
		}

		if !expectType.Equals(actualType) {
			t.Fatal(testutil.Testing("app rule type result").FailMessage(expectType, actualType, step))
		}
	}

	// step 5: discharge assumption and introduce let expression w/ type Int
	step++
	conclusion := discharge_x_assumption(bridge.Judgement(app_x_0_conclusion.GetExpressionAndType()))
	{
		actualExpr, actualType := conclusion.GetExpressionAndType()
		expectExpr, expectType := letExpr, Int

		if !expectExpr.StrictEquals(actualExpr) {
			t.Fatal(testutil.Testing("let rule expression result").FailMessage(expectExpr, actualExpr, step))
		}

		if !expectType.Equals(actualType) {
			t.Fatal(testutil.Testing("let rule type result").FailMessage(expectType, actualType, step))
		}
	}
}
