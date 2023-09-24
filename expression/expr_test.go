package expr

import "testing"

func TestApply(t *testing.T) {
	tests := []struct {
		l      Function
		r      Expression
		expect Expression
	}{
		{ // (λx[1] . x[1]) y[0] == y[1]
			IdFunction,
			Var("y"),
			makeVar("y", 1),
		},
		{ // (λt[2] f[1] . t[2]) z[0] == (λf[1] . z[2])
			TrueFunction,
			Var("z"),
			makeFunction([]Variable{makeVar("f", 1)}, makeVar("z", 2)),
		},
		{ // (λt[2] f[1] . f[1]) z[0] == (λf[1] . f[1])
			FalseFunction,
			Var("z"),
			makeFunction([]Variable{makeVar("f", 1)}, makeVar("f", 1)),
		},
	}

	for i, test := range tests {
		actual := test.l.Apply(test.r)
		if !actual.StrictEquals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", i+1, test.expect.StrictString(), actual.StrictString())
		}
	}
}

func TestApply2(t *testing.T) {
	tests := []struct {
		l      Function
		r      [2]Expression
		expect Expression
	}{
		{ // (λt[2] f[1] . t[2]) z[0] w[0] == (λf[1] . z[2]) w[0] == z[1]
			TrueFunction,
			[2]Expression{Var("z"), Var("w")},
			makeVar("z", 1),
		},
		{ // (λt[2] f[1] . f[1]) z[0] w[0] == (λf . f) w[0] == w[1]
			FalseFunction,
			[2]Expression{Var("z"), Var("w")},
			makeVar("w", 1),
		},
		{ // and true true == (λa b . a b (λt f . f)) (λt f . t) (λt f . t) == (λt f . t)
			AndFunction,
			[2]Expression{TrueFunction, TrueFunction},
			TrueFunction,
		},
		{ // and true false == (λa b . a b (λt f . f)) (λt f . t) (λt f . f) == (λt f . f)
			AndFunction,
			[2]Expression{TrueFunction, FalseFunction},
			FalseFunction,
		},
		{ // and false true == (λa b . a b (λt f . f)) (λt f . f) (λt f . t) == (λt f . f)
			AndFunction,
			[2]Expression{FalseFunction, TrueFunction},
			FalseFunction,
		},
		{ // and false false == (λa b . a b (λt f . f)) (λt f . f) (λt f . f) == (λt f . f)
			AndFunction,
			[2]Expression{FalseFunction, FalseFunction},
			FalseFunction,
		},
	}

	for i, test := range tests {
		actual := test.l.Apply(test.r[0])
		f, ok := actual.(Function)
		if !ok {
			t.Fatalf("failed test #%d: first application did not produce a function, found:\n%v\n", i+1, actual)
		}
		actual = f.Apply(test.r[1])
		if !actual.StrictEquals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", i+1, test.expect.StrictString(), actual.StrictString())
		}
	}
}

func TestEtaReduction(t *testing.T) {
	tests := []struct {
		in     Function
		expect Function
	}{
		{
			// (λx[1] . (λy[1] . y[1]) x[1]) == (λy[1] . y[1])
			in:     Bind(Var("x")).In(Apply(Bind(Var("y")).In(Var("y")), Var("x"))),
			expect: Bind(Var("y")).In(Var("y")),
		},
		{
			// (λa[2] x[1] . (λy[1] . y[1]) x[1]) == (λa[2] y[1] . y[1])
			in:     Bind(Var("a"), Var("x")).In(Apply(Bind(Var("y")).In(Var("y")), Var("x"))),
			expect: Bind(Var("a"), Var("y")).In(Var("y")),
		},
		{
			// (λa[2] x[1] . (λy[1] . a[3] y[1]) x[1]) == (λa[2] y[1] . a[2] y[1])
			in:     Bind(Var("a"), Var("x")).In(Apply(Bind(Var("y")).In(Apply(Var("a"), Var("y"))), Var("x"))),
			expect: Bind(Var("a"), Var("y")).In(Apply(Var("a"), Var("y"))),
		},
		{
			// (λx[1] . (λy[1] . x[2] y[1])) == (λx[1] . (λy[1] . x[2] y[1]))
			in:     Bind(Var("x")).In(Apply(Bind(Var("y")).In(Var("y")), Var("x"))),
			expect: Bind(Var("y")).In(Var("y")),
		},
		{
			// (λa[2] x[1] . (λy[1] . x[2] y[1]) x[1]) == (λa[2] x[1] . (λy[1] . x[2] y[1]) x[1])
			in:     Bind(Var("a"), Var("x")).In(Apply(Bind(Var("y")).In(Apply(Var("x"), Var("y"))), Var("x"))),
			expect: Bind(Var("a"), Var("x")).In(Apply(Bind(Var("y")).In(Apply(Var("x"), Var("y"))), Var("x"))),
		},
	}

	for i, test := range tests {
		actual := test.in.EtaReduction()
		if !actual.StrictEquals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", i+1, test.expect.StrictString(), actual.StrictString())
		}
	}
}

func TestInstruction(t *testing.T) {
	to_string := DefineInstruction("toString", 1, func(instr InstructionArgs) Expression {
		arg := instr.GetArgAtIndex(0)
		return Const(arg.String())
	})

	eq10 := DefineInstruction("eq10", 1, func(instr InstructionArgs) Expression {
		arg := instr.GetArgAtIndex(0)
		if arg.Equals(List{Const("1"), Const("0")}) {
			return TrueFunction
		} 
		return FalseFunction
	})

	succ := DefineInstruction("succ", 1, func(instr InstructionArgs) Expression {
		arg := instr.GetArgAtIndex(0)
		ls, ok := arg.(List)
		if !ok {
			panic("not a list!\n")
		}
		var inc bool = len(ls) > 0
		var newLs List
		if inc {
			newLs = ls.copy()
		} else {
			newLs = ls
		}

		for i := len(ls) - 1; i >= 0 && inc; i-- {
			inc = false
			c, ok := newLs[i].ForceRequest().(Const)
			if !ok || len(c) != 1 {
				panic("not an integer!\n")
			}
			switch c[0] {
			case '0':
				newLs[i] = Const("1")
			case '1':
				newLs[i] = Const("2")
			case '2':
				newLs[i] = Const("3")
			case '3':
				newLs[i] = Const("4")
			case '4':
				newLs[i] = Const("5")
			case '5':
				newLs[i] = Const("6")
			case '6':
				newLs[i] = Const("7")
			case '7':
				newLs[i] = Const("8")
			case '8':
				newLs[i] = Const("9")
			case '9':
				newLs[i] = Const("0")
				inc = true
			default:
				panic("not an integer!\n")
			}
		}
		if inc {
			newLs = append(List{Const("1")}, newLs...)
		}
		return newLs
	})

	tests := []struct{
		in Expression
		expect Expression
	}{
		{
			// (toString z) == Const("z")
			in: Apply(to_string.MakeInstance(), Var("z")),
			expect: Const("z"),
		},
		{
			// (λz . (toString z)) (λt f . t) == Const("(λt f . t)")
			in: Apply(
				Bind(Var("z")).In(Apply(to_string.MakeInstance(), Var("z"))), TrueFunction),
			expect: Const(TrueFunction.String()),
		},
		{
			// (succ [0]) == List{Const("1")}
			in: Apply(succ.MakeInstance(), List{Const("0")}),
			expect: List{Const("1")},
		},
		/*{
			// (Y ((λf x . (eq10 x) x (f x)) succ) [0]) == [1, 0]
			in: Apply(Y, Bind(Var("f"), Var("x")).In(Apply(eq10.MakeInstance(), Var("x"), Var("x"), Apply(Var("f"), Var("x")))), succ.MakeInstance(), List{Const("0")}),
			expect: List{Const("1"), Const("0")},
		},*/
	}

	println(Y.Apply(Bind(Var("f"), Var("x")).In(Apply(eq10.MakeInstance(), Var("x"), Var("x"), Apply(Var("f"), Var("x"))))))

	for i, test := range tests {
		actual := test.in.ForceRequest()
		if !actual.StrictEquals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", i+1, test.expect.StrictString(), actual.StrictString())
		}
	}
} 

func TestLaziness(t *testing.T) {
	// if evaluation is not lazy, then expression will never finish evaluating
	// Ω = (λx . x x)
	omega := Bind(Var("x")).In(Apply(Var("x"), Var("x")))
	// omegaSkipFunction = (λx . x z (Ω Ω))
	omegaSkipFunction := Bind(Var("x")).In(Apply(Var("x"), Var("z"), Apply(omega, omega)))
	omegaSkipFunction.Apply(TrueFunction).ForceRequest()

	// Y (λf y . g y) == (λg y . g y) (Y (λg y . g y))
	// == (λy . (Y (λg y . g y)) y)
	inFunc := Bind(Var("g"), Var("y")).In(Apply(Var("g"), Var("y")))
	actual := Y.Apply(inFunc)
	// inFunc ((λx . inFunc (x x)) (λx . inFunc (x x)))
	expect := Function{
		vars: []Variable{makeVar("y", 1)},
		e: Apply(Apply(
			Function{
				vars: []Variable{makeVar("x", 1)},
				e: Apply(inFunc, Apply(makeVar("x", 1), makeVar("x", 1))),
			},
			Function{
				vars: []Variable{makeVar("x", 1)},
				e: Apply(inFunc, Apply(makeVar("x", 1), makeVar("x", 1))),
			},
		), makeVar("y", 1)),
	}
	
	if !expect.StrictEquals(actual) {
		t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", 2, expect.StrictString(), actual.StrictString())
	}
}

func Test(t *testing.T) {
	s := Select(Var("x"), 
		Bind(Var("y")).When(Apply(Const("Just"), Var("y"))).Then(Var("y")),
		Bind(Var("w")).When(Const("Nothing")).Then(Var("w")),
	)
	println(s.String())
	println(s.StrictString())
}

/*func Test(t *testing.T) {
	ConsFunc := Bind(Var("el"), Var("els"), Var("x")).In(Apply(Var("x"), Var("el"), Var("els")))
	//HeadFunc := Bind(Var("ls")).In(Apply(Var("ls"), TrueFunction))
	//TailFunc := Bind(Var("ls")).In(Apply(Var("ls"), FalseFunction))

	s := ConsFunc.Apply(Const("1")).(Function).Apply(
		ConsFunc.Apply(Const("2")).(Function).Apply(
			ConsFunc.Apply(Const("3")).(Function).Apply(
				ConsFunc.Apply(Const("4")).(Function).Apply(
					FalseFunction)))).String()
	println(s)
}*/
