package expr

import "testing"

var base = NewContext[test_named]().SetNameMaker(func(s string) test_named {
	return test_named(s)
})

func _Const(s string) Const[test_named] {
	return Const[test_named]{test_named(s)}
}

func _Apply(e1, e2 Expression[test_named], es ...Expression[test_named]) Application[test_named] {
	return Apply[test_named](e1, e2, es...)
}

func _Bind(binder Variable[test_named], more ...Variable[test_named]) BindersOnly[test_named] {
	return Bind[test_named](binder, more...)
}

func _Var(s string) Variable[test_named] {
	return Var[test_named](test_named(s))
}

func _makeVar(s string, depth int) Variable[test_named] {
	return base.makeVar(s, depth)
}

func TestApply(t *testing.T) {
	tests := []struct {
		l      Function[test_named]
		r      Expression[test_named]
		expect Expression[test_named]
	}{
		{ // (λx[1] . x[1]) y[0] == y[1]
			IdFunction,
			_Var("y"),
			_makeVar("y", 1),
		},
		{ // (λt[2] f[1] . t[2]) z[0] == (λf[1] . z[2])
			TrueFunction,
			_Var("z"),
			makeFunction[test_named]([]Variable[test_named]{_makeVar("f", 1)}, _makeVar("z", 2)),
		},
		{ // (λt[2] f[1] . f[1]) z[0] == (λf[1] . f[1])
			FalseFunction,
			_Var("z"),
			makeFunction[test_named]([]Variable[test_named]{_makeVar("f", 1)}, _makeVar("f", 1)),
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
		l      Function[test_named]
		r      [2]Expression[test_named]
		expect Expression[test_named]
	}{
		{ // (λt[2] f[1] . t[2]) z[0] w[0] == (λf[1] . z[2]) w[0] == z[1]
			TrueFunction,
			[2]Expression[test_named]{_Var("z"), _Var("w")},
			_makeVar("z", 1),
		},
		{ // (λt[2] f[1] . f[1]) z[0] w[0] == (λf . f) w[0] == w[1]
			FalseFunction,
			[2]Expression[test_named]{_Var("z"), _Var("w")},
			_makeVar("w", 1),
		},
		{ // and true true == (λa b . a b (λt f . f)) (λt f . t) (λt f . t) == (λt f . t)
			AndFunction,
			[2]Expression[test_named]{TrueFunction, TrueFunction},
			TrueFunction,
		},
		{ // and true false == (λa b . a b (λt f . f)) (λt f . t) (λt f . f) == (λt f . f)
			AndFunction,
			[2]Expression[test_named]{TrueFunction, FalseFunction},
			FalseFunction,
		},
		{ // and false true == (λa b . a b (λt f . f)) (λt f . f) (λt f . t) == (λt f . f)
			AndFunction,
			[2]Expression[test_named]{FalseFunction, TrueFunction},
			FalseFunction,
		},
		{ // and false false == (λa b . a b (λt f . f)) (λt f . f) (λt f . f) == (λt f . f)
			AndFunction,
			[2]Expression[test_named]{FalseFunction, FalseFunction},
			FalseFunction,
		},
	}

	for i, test := range tests {
		actual := test.l.Apply(test.r[0])
		f, ok := actual.(Function[test_named])
		if !ok {
			t.Fatalf("failed test #%d: first application[test_named] did not produce a function[test_named], found:\n%v\n", i+1, actual)
		}
		actual = f.Apply(test.r[1])
		if !actual.StrictEquals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", i+1, test.expect.StrictString(), actual.StrictString())
		}
	}
}

func TestEtaReduction(t *testing.T) {
	tests := []struct {
		in     Function[test_named]
		expect Function[test_named]
	}{
		{
			// (λx[1] . (λy[1] . y[1]) x[1]) == (λy[1] . y[1])
			in:     _Bind(_Var("x")).In(_Apply(_Bind(_Var("y")).In(_Var("y")), _Var("x"))),
			expect: _Bind(_Var("y")).In(_Var("y")),
		},
		{
			// (λa[2] x[1] . (λy[1] . y[1]) x[1]) == (λa[2] y[1] . y[1])
			in:     _Bind(_Var("a"), _Var("x")).In(_Apply(_Bind(_Var("y")).In(_Var("y")), _Var("x"))),
			expect: _Bind(_Var("a"), _Var("y")).In(_Var("y")),
		},
		{
			// (λa[2] x[1] . (λy[1] . a[3] y[1]) x[1]) == (λa[2] y[1] . a[2] y[1])
			in:     _Bind(_Var("a"), _Var("x")).In(_Apply(_Bind(_Var("y")).In(_Apply(_Var("a"), _Var("y"))), _Var("x"))),
			expect: _Bind(_Var("a"), _Var("y")).In(_Apply(_Var("a"), _Var("y"))),
		},
		{
			// (λx[1] . (λy[1] . x[2] y[1])) == (λx[1] . (λy[1] . x[2] y[1]))
			in:     _Bind(_Var("x")).In(_Apply(_Bind(_Var("y")).In(_Var("y")), _Var("x"))),
			expect: _Bind(_Var("y")).In(_Var("y")),
		},
		{
			// (λa[2] x[1] . (λy[1] . x[2] y[1]) x[1]) == (λa[2] x[1] . (λy[1] . x[2] y[1]) x[1])
			in:     _Bind(_Var("a"), _Var("x")).In(_Apply(_Bind(_Var("y")).In(_Apply(_Var("x"), _Var("y"))), _Var("x"))),
			expect: _Bind(_Var("a"), _Var("x")).In(_Apply(_Bind(_Var("y")).In(_Apply(_Var("x"), _Var("y"))), _Var("x"))),
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
	to_string := DefineInstruction[test_named]("toString", 1, func(instr InstructionArgs[test_named]) Expression[test_named] {
		arg := instr.GetArgAtIndex(0)
		return _Const(arg.String())
	})

	eq10 := DefineInstruction[test_named]("eq10", 1, func(instr InstructionArgs[test_named]) Expression[test_named] {
		arg := instr.GetArgAtIndex(0)
		if arg.Equals(base, List[test_named]{_Const("1"), _Const("0")}) {
			return TrueFunction
		} 
		return FalseFunction
	})

	succ := DefineInstruction[test_named]("succ", 1, func(instr InstructionArgs[test_named]) Expression[test_named] {
		arg := instr.GetArgAtIndex(0)
		ls, ok := arg.(List[test_named])
		if !ok {
			panic("not a list!\n")
		}
		var inc bool = len(ls) > 0
		var newLs List[test_named]
		if inc {
			newLs = ls.copy()
		} else {
			newLs = ls
		}

		for i := len(ls) - 1; i >= 0 && inc; i-- {
			inc = false
			c, ok := newLs[i].ForceRequest().(Const[test_named])
			if !ok || len(c.String()) != 1 {
				panic("not an integer!\n")
			}
			switch c.String()[0] {
			case '0':
				newLs[i] = _Const("1")
			case '1':
				newLs[i] = _Const("2")
			case '2':
				newLs[i] = _Const("3")
			case '3':
				newLs[i] = _Const("4")
			case '4':
				newLs[i] = _Const("5")
			case '5':
				newLs[i] = _Const("6")
			case '6':
				newLs[i] = _Const("7")
			case '7':
				newLs[i] = _Const("8")
			case '8':
				newLs[i] = _Const("9")
			case '9':
				newLs[i] = _Const("0")
				inc = true
			default:
				panic("not an integer!\n")
			}
		}
		if inc {
			newLs = append(List[test_named]{_Const("1")}, newLs...)
		}
		return newLs
	})

	tests := []struct{
		in Expression[test_named]
		expect Expression[test_named]
	}{
		{
			// (toString z) == _Const("z")
			in: _Apply(to_string.MakeInstance(), _Var("z")),
			expect: _Const("z"),
		},
		{
			// (λz . (toString z)) (λt f . t) == _Const("(λt f . t)")
			in: _Apply(
				_Bind(_Var("z")).In(_Apply(to_string.MakeInstance(), _Var("z"))), TrueFunction),
			expect: _Const(TrueFunction.String()),
		},
		{
			// (succ [0]) == List{_Const("1")}
			in: _Apply(succ.MakeInstance(), List[test_named]{_Const("0")}),
			expect: List[test_named]{_Const("1")},
		},
		/*{
			// (Y ((λf x . (eq10 x) x (f x)) succ) [0]) == [1, 0]
			in: _Apply(Y, _Bind(_Var("f"), _Var("x")).In(_Apply(eq10.MakeInstance(), _Var("x"), _Var("x"), _Apply(_Var("f"), _Var("x")))), succ.MakeInstance(), List{_Const("0")}),
			expect: List{_Const("1"), _Const("0")},
		},*/
	}

	println(Y.Apply(_Bind(_Var("f"), _Var("x")).In(_Apply(eq10.MakeInstance(), _Var("x"), _Var("x"), _Apply(_Var("f"), _Var("x"))))))

	for i, test := range tests {
		actual := test.in.ForceRequest()
		if !actual.StrictEquals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", i+1, test.expect.StrictString(), actual.StrictString())
		}
	}
} 

func TestLaziness(t *testing.T) {
	// if evaluation[test_named] is not lazy, then expression[test_named] will never finish evaluating
	// Ω = (λx . x x)
	omega := _Bind(_Var("x")).In(_Apply(_Var("x"), _Var("x")))
	// omegaSkipFunction[test_named] = (λx . x z (Ω Ω))
	omegaSkipFunction := _Bind(_Var("x")).In(_Apply(_Var("x"), _Var("z"), _Apply(omega, omega)))
	omegaSkipFunction.Apply(TrueFunction).ForceRequest()

	// Y (λf y . g y) == (λg y . g y) (Y (λg y . g y))
	// == (λy . (Y (λg y . g y)) y)
	inFunc := _Bind(_Var("g"), _Var("y")).In(_Apply(_Var("g"), _Var("y")))
	actual := Y.Apply(inFunc)
	// inFunc ((λx . inFunc (x x)) (λx . inFunc (x x)))
	expect := Function[test_named]{
		vars: []Variable[test_named]{_makeVar("y", 1)},
		e: _Apply(_Apply(
			Function[test_named]{
				vars: []Variable[test_named]{_makeVar("x", 1)},
				e: _Apply(inFunc, _Apply(_makeVar("x", 1), _makeVar("x", 1))),
			},
			Function[test_named]{
				vars: []Variable[test_named]{_makeVar("x", 1)},
				e: _Apply(inFunc, _Apply(_makeVar("x", 1), _makeVar("x", 1))),
			},
		), _makeVar("y", 1)),
	}
	
	if !expect.StrictEquals(actual) {
		t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", 2, expect.StrictString(), actual.StrictString())
	}
}

func Test(t *testing.T) {
	s := Select[test_named](_Var("x"), 
		_Bind(_Var("y")).When(_Apply(_Const("Just"), _Var("y"))).Then(_Var("y")),
		_Bind(_Var("w")).When(_Const("Nothing")).Then(_Var("w")),
	)
	println(s.String())
	println(s.StrictString())
}

/*func Test(t *testing.T) {
	ConsFunc := _Bind(_Var("el"), _Var("els"), _Var("x")).In(_Apply(_Var("x"), _Var("el"), _Var("els")))
	//HeadFunc := _Bind(_Var("ls")).In(_Apply(_Var("ls"), TrueFunction[test_named]))
	//TailFunc := _Bind(_Var("ls")).In(_Apply(_Var("ls"), FalseFunction[test_named]))

	s := ConsFunc._Apply(_Const("1")).(Function[test_named])._Apply(
		ConsFunc._Apply(_Const("2")).(Function[test_named])._Apply(
			ConsFunc._Apply(_Const("3")).(Function[test_named])._Apply(
				ConsFunc._Apply(_Const("4")).(Function[test_named])._Apply(
					FalseFunction[test_named])))).String()
	println(s)
}*/
