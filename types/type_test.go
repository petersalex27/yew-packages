package types

import (
	"testing"

	expr "github.com/petersalex27/yew-packages/expr"
)

type test_nameable string

func (n test_nameable) GetName() string {
	return string(n)
}

func test_nameable_fn(s string) test_nameable {
	return test_nameable(s)
}

func judgment[E expr.Expression[test_nameable]](e E, ty Type[test_nameable]) TypeJudgment[test_nameable, E] {
	return Judgment(e, ty)
}

var base = NewContext[test_nameable]().
	SetNameMaker(test_nameable_fn)

func TestString(t *testing.T) {
	tests := []struct {
		in     Type[test_nameable]
		expect string
	}{
		// monotypes
		{in: _Con("Type"), expect: "Type"},                                     // just type
		{in: _App("Type", _Con("x")), expect: "(Type x)"},                      // application
		{in: _App("Type", _App("Type", _Con("x"))), expect: "(Type (Type x))"}, // nested application
		{in: _Var("a"), expect: "a"},                                           // variable
		{in: _App("Type", _Var("a")), expect: "(Type a)"},                      // application w/ variable
		{in: Apply[test_nameable](_Var("a"), _App("Type", _Var("a"))), expect: ("(a (Type a))")},
		{in: Apply[test_nameable](_App("Type", _Var("a")), _App("Type", _Var("a"))), expect: ("(Type a (Type a))")},
		{in: Apply[test_nameable](base.Infix(_Con("Type"), "->", _Con("Type")), _App("Type", _Var("a"))), expect: ("(Type -> Type (Type a))")},
		{in: Apply[test_nameable](_App("Type", _Var("a")), base.Infix(_Con("Type"), "->", _Con("Type"))), expect: ("(Type a (Type -> Type))")},
		{in: base.Infix(_Con("Type"), "->", _Con("Type")), expect: "(Type -> Type)"},
		{in: base.Infix(_App("Color", _Con("Red")), "->", _Con("Type")), expect: "((Color Red) -> Type)"},
		{in: base.Infix(_Con("Type"), "->", _App("Color", _Con("Red"))), expect: "(Type -> (Color Red))"},
		{in: base.Infix(_Con("Type"), "->", base.Infix(_Con("Type"), "->", _Con("Type"))), expect: "(Type -> (Type -> Type))"},
		{in: base.Infix(base.Infix(_Con("Type"), "->", _Con("Type")), "->", _Con("Type")), expect: "((Type -> Type) -> Type)"},
		{in: base.Infix(_Var("a"), "->", _Var("a")), expect: "(a -> a)"},
		{in: base.Infix(_Var("a"), "->"), expect: "((->) a)"},
		{in: Application[test_nameable]{c: base.InfixCon("->"), ts: nil}, expect: "(->)"},
		{in: Application[test_nameable]{c: base.EnclosingCon(1, "[]"), ts: nil}, expect: "[]"},
		{in: Apply[test_nameable](base.EnclosingCon(1, "[]"), _Con("A")), expect: "[A]"},
		{
			in: Apply[test_nameable](
				base.EnclosingCon(1, "[]"),
				Apply[test_nameable](base.EnclosingCon(1, "{}"), _Con("A")),
			),
			expect: "[{A}]",
		},
		{in: Apply[test_nameable](base.EnclosingCon(1, "[]"), _Con("Type"), _Var("a")), expect: "[Type a]"},
		// polytypes
		{in: _Con("Type").Generalize(base), expect: "forall _ . Type"},
		{in: _Forall("a").Bind(_Con("Type")), expect: "forall a . Type"},
		{in: _Forall("a", "b").Bind(_Con("Type")), expect: "forall a b . Type"},
		{in: _Forall("a", "b").Bind(_App("Type", _App("Type2", _Var("b"), _Var("a")))), expect: "forall a b . (Type (Type2 b a))"},
		{in: _Forall("a").Bind(_Var("a")), expect: "forall a . a"},
		{in: _Forall("a").Bind(Apply[test_nameable](base.EnclosingCon(1, "[]"), _Var("a"))), expect: "forall a . [a]"},
		// dependent types
		{
			in:     Index[test_nameable](_App("Array", _Con("Int")), Judgment[test_nameable, expr.Referable[test_nameable]](expr.Var(base.makeName("n")), _Con("Uint"))),
			expect: "((Array Int); (n: Uint))",
		},
		{
			in:     Index[test_nameable](Apply[test_nameable](base.EnclosingCon(1, "[]"), _Con("Int")), Judgment[test_nameable, expr.Referable[test_nameable]](expr.Var(base.makeName("n")), _Con("Uint"))),
			expect: "[Int; (n: Uint)]",
		},
		{
			in:     Index[test_nameable](Apply[test_nameable](base.EnclosingCon(1, "[]"), _Con("Int")), FreeJudge[test_nameable, expr.Referable[test_nameable]](base, expr.Var(base.makeName("n")))),
			expect: "[Int; n]",
		},
		{
			in: Index[test_nameable](
				Apply[test_nameable](
					base.EnclosingCon(1, "[]"),
					Index[test_nameable](
						Apply[test_nameable](
							base.EnclosingCon(1, "[]"), _Con("Int"),
						),
						Judgment[test_nameable, expr.Referable[test_nameable]](
							expr.Var(base.makeName("n")),
							_Con("Uint"),
						),
					),
				),
				Judgment[test_nameable, expr.Referable[test_nameable]](
					expr.Var(base.makeName("n")),
					_Con("Uint"),
				),
			),
			expect: "[[Int; (n: Uint)]; (n: Uint)]",
		},
		{
			in: DependentType[test_nameable]{
				[]TypeJudgment[test_nameable, expr.Variable[test_nameable]]{
					judgment(expr.Var(base.makeName("n")), _Con("Uint")),
				},
				Index[test_nameable](_App("Array", _Var("a"))),
			},
			expect: "mapval (n: Uint) . (Array a)",
		},
	}

	for testIndex, test := range tests {
		if actual := test.in.String(); actual != test.expect {
			t.Fatalf("failed test #%d:\nexpected:\n%s\nactual:\n%s\n", testIndex+1, test.expect, actual)
		}
	}
}

func _Forall(vs ...string) binders[test_nameable] {
	return base.Forall(vs...)
}

func _App(name string, ms ...Monotyped[test_nameable]) Application[test_nameable] {
	return base.App(name, ms...)
}

func _Var(name string) Variable[test_nameable] {
	return base.Var(name)
}

func _Con(name string) Constant[test_nameable] {
	return base.Con(name)
}

func TestInstantiate(t *testing.T) {
	tests := []struct {
		poly   Polytype[test_nameable]
		mono   Monotyped[test_nameable]
		expect Type[test_nameable]
	}{
		{ // (forall a . Type a) $ b == Type b
			poly:   _Forall("a").Bind(_App("Type", _Var("a"))),
			mono:   _Var("b"),
			expect: _App("Type", _Var("b")),
		},
		{ // (forall a . Type a) $ Louis == Type Louis
			poly:   _Forall("a").Bind(_App("Type", _Var("a"))),
			mono:   _Con("Louis"),
			expect: _App("Type", _Con("Louis")),
		},
		{ // (forall a . Type b) $ c == Type b
			poly:   _Forall("a").Bind(_App("Type", _Var("b"))),
			mono:   _Var("c"),
			expect: _App("Type", _Var("b")),
		},
		{ // (forall a b c . Type a b c) $ x == forall b c . Type x b c
			poly:   _Forall("a", "b", "c").Bind(_App("Type", _Var("a"), _Var("b"), _Var("c"))),
			mono:   _Var("x"),
			expect: _Forall("b", "c").Bind(_App("Type", _Var("x"), _Var("b"), _Var("c"))),
		},
		{ // (forall b c . Type x b c) $ x == forall c . Type x x c
			poly:   _Forall("b", "c").Bind(_App("Type", _Var("x"), _Var("b"), _Var("c"))),
			mono:   _Var("x"),
			expect: _Forall("c").Bind(_App("Type", _Var("x"), _Var("x"), _Var("c"))),
		},
		{ // (forall c . Type x x c) $ x == Type x x x
			poly:   _Forall("c").Bind(_App("Type", _Var("x"), _Var("x"), _Var("c"))),
			mono:   _Var("x"),
			expect: _App("Type", _Var("x"), _Var("x"), _Var("x")),
		},
		{ // (forall a . (a (Type a))) $ Louis == (Louis (Type Louis))
			poly:   _Forall("a").Bind(Apply[test_nameable](_Var("a"), _App("Type", _Var("a")))),
			mono:   _Con("Louis"),
			expect: Apply[test_nameable](_Con("Louis"), _App("Type", _Con("Louis"))),
		},
		{ // (forall a . (a (Type (a -> a)))) $ Louis == (Louis (Type (Louis -> Louis)))
			poly:   _Forall("a").Bind(Apply[test_nameable](_Var("a"), _App("Type", _Function(_Var("a"), _Var("a"))))),
			mono:   _Con("Louis"),
			expect: Apply[test_nameable](_Con("Louis"), _App("Type", _Function(_Con("Louis"), _Con("Louis")))),
		},
		{ // (forall a . ((a -> a) (Type (a -> a)))) $ Louis == ((Louis -> Louis) (Type (Louis -> Louis)))
			poly:   _Forall("a").Bind(Apply[test_nameable](_Function(_Var("a"), _Var("a")), _App("Type", _Function(_Var("a"), _Var("a"))))),
			mono:   _Con("Louis"),
			expect: Apply[test_nameable](_Function(_Con("Louis"), _Con("Louis")), _App("Type", _Function(_Con("Louis"), _Con("Louis")))),
		},
		{
			poly: _Forall("a").Bind(
				Index(
					_App("Array", _Var("a")),
					ExpressionJudgment[test_nameable, expr.Referable[test_nameable]](
						Judgment(
							expr.Referable[test_nameable](expr.Var(base.makeName("n"))),
							Type[test_nameable](_Con("Uint")),
						),
					),
				),
			),
			mono: _Con("Int"),
			expect: Index(
				_App("Array", _Con("Int")),
				ExpressionJudgment[test_nameable, expr.Referable[test_nameable]](
					Judgment(
						expr.Referable[test_nameable](expr.Var(base.makeName("n"))),
						Type[test_nameable](_Con("Uint")),
					),
				),
			),
		},
	}

	for testIndex, test := range tests {
		actual := test.poly.Instantiate(test.mono)
		if !actual.Equals(test.expect) {
			t.Fatalf("failed test #%d\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect, actual)
		}
	}
}

// Let [t] = {a, b, ..} for any monotype a, b, and t. find(a) == t, find(b) == t, ..

func _Function(left Monotyped[test_nameable], right Monotyped[test_nameable]) Application[test_nameable] {
	return base.Function(left, right)
}
