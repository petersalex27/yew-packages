package types

import (
	"fmt"
	"os"
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

func judgement[E expr.Expression[test_nameable]](e E, ty Type[test_nameable]) TypeJudgement[test_nameable, E] {
	return Judgement(e, ty)
}

type refer ReferableType[test_nameable]

var base = 
		NewContext[test_nameable]().
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
		{in: _Var("a"), expect: "a"},                                             // variable
		{in: _App("Type", _Var("a")), expect: "(Type a)"},                        // application w/ variable
		{in: Apply[test_nameable](_Var("a"), _App("Type", _Var("a"))), expect: ("(a (Type a))")},
		{in: Apply[test_nameable](_App("Type", _Var("a")), _App("Type", _Var("a"))), expect: ("(Type a (Type a))")},
		{in: Apply[test_nameable](base.Infix(_Con("Type"), "->", _Con("Type")), _App("Type", _Var("a"))), expect: ("(Type -> Type (Type a))")},
		{in: Apply[test_nameable](_App("Type", _Var("a")), base.Infix(_Con("Type"), "->", _Con("Type"))), expect: ("(Type a (Type -> Type))")},
		{in: base.Infix(_Con("Type"), "->", _Con("Type")), expect: "(Type -> Type)"},
		{in: base.Infix(_App("Color", _Con("Red")), "->", _Con("Type")), expect: "((Color Red) -> Type)"},
		{in: base.Infix(_Con("Type"), "->", _App("Color", _Con("Red"))), expect: "(Type -> (Color Red))"},
		{in: base.Infix(_Con("Type"), "->", base.Infix(_Con("Type"), "->", _Con("Type"))), expect: "(Type -> (Type -> Type))"},
		{in: base.Infix(base.Infix(_Con("Type"), "->", _Con("Type")), "->", _Con("Type")), expect: "((Type -> Type) -> Type)"},
		{in: base.Infix(_Var("a"), "->", _Var("a")), expect: "(a -> a)",},
		{in: base.Infix(_Var("a"), "->"), expect: "((->) a)",},
		{in: Application[test_nameable]{c: base.InfixCon("->"), ts: nil,}, expect: "(->)",},
		{in: Application[test_nameable]{c: base.EnclosingCon(1, "[]"), ts: nil,}, expect: "[]"},
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
			in: Index[test_nameable](_App("Array", _Con("Int")), Judgement[test_nameable,expr.Expression[test_nameable]](expr.Var(base.makeName("n")), _Con("Uint"))),
			expect: "((Array Int); (n: Uint))",
		},
		{
			in: Index[test_nameable](Apply[test_nameable](base.EnclosingCon(1, "[]"), _Con("Int")), Judgement[test_nameable,expr.Expression[test_nameable]](expr.Var(base.makeName("n")), _Con("Uint"))),
			expect: "[Int; (n: Uint)]",
		},
		{
			in: Index[test_nameable](Apply[test_nameable](base.EnclosingCon(1, "[]"), _Con("Int")), FreeJudge[test_nameable,expr.Expression[test_nameable]](base, expr.Var(base.makeName("n")))),
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
						Judgement[test_nameable,expr.Expression[test_nameable]](
							expr.Var(base.makeName("n")), 
							_Con("Uint"),
						),
					),
				), 
				Judgement[test_nameable,expr.Expression[test_nameable]](
					expr.Var(base.makeName("n")), 
					_Con("Uint"),
				),
			),
			expect: "[[Int; (n: Uint)]; (n: Uint)]",
		},
		{
			in: DependentType[test_nameable]{
				[]TypeJudgement[test_nameable, expr.Variable[test_nameable]]{
					judgement(expr.Var(base.makeName("n")), _Con("Uint")),
				},
				Index[test_nameable](_App("Array", _Var("a"))),
			},
			expect: "mapall (n: Uint) . (Array a)",
		},
	}

	for testIndex, test := range tests {
		if actual := test.in.String(); actual != test.expect {
			t.Fatalf("failed test #%d:\nexpected:\n%s\nactual:\n%s\n", testIndex+1, test.expect, actual)
		}
	}
}

func _Forall(vs ...string)partialPoly[test_nameable] {
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
			poly:	_Forall("a").Bind(
				Index(
					_App("Array", _Var("a")), 
					ExpressionJudgement[test_nameable, expr.Expression[test_nameable]](
						Judgement(
							expr.Expression[test_nameable](expr.Var(base.makeName("n"))), 
							Type[test_nameable](_Con("Uint")),
						),
					),
				),
			),
			mono:	_Con("Int"),
			expect: Index(
				_App("Array", _Con("Int")),
				ExpressionJudgement[test_nameable, expr.Expression[test_nameable]](
					Judgement(
						expr.Expression[test_nameable](expr.Var(base.makeName("n"))), 
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

type mts_test []Monotyped[test_nameable]

// Let [t] = {a, b, ..} for any monotype a, b, and t. find(a) == t, find(b) == t, ..
func TestFind_inTable(t *testing.T) {
	tests := []struct {
		representative Monotyped[test_nameable]
		classMembers   []Monotyped[test_nameable]
	}{
		{
			representative: _Con("Type"),
			classMembers:   mts_test{_Var("a")},
		},
		{
			representative: _Var("a"),
			classMembers:   mts_test{_Var("b")},
		},
		{
			representative: _Con("Type"),
			classMembers:   mts_test{_Con("Type2")},
		},
		{
			representative: _Con("Type"),
			classMembers:   mts_test{_Con("Type")},
		},
		{
			representative: _Var("a"),
			classMembers:   mts_test{_Var("a")},
		},
		{
			representative: _Con("Type"),
			classMembers:   mts_test{_Var("a"), _Con("Type")},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]()
		for _, m := range test.classMembers {
			name := m.String()
			cxt.equivClasses[name] = test.representative
		}
		for _, m := range test.classMembers {
			actual := cxt.Find(m)
			if !test.representative.Equals(actual) {
				t.Fatalf("failed test #%d\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.representative, actual)
			}
		}
	}
}

// expect all calls to find(m) == m for any monotype m
func TestFind_notInTable(t *testing.T) {
	tests := []struct {
		classMembers mts_test
	}{
		{
			classMembers: mts_test{_Var("a")},
		},
		{
			classMembers: mts_test{_Con("Type")},
		},
		{
			classMembers: mts_test{_Var("a"), _Con("Type")},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		for _, m := range test.classMembers {
			actual := cxt.Find(m)
			if !m.Equals(actual) {
				t.Fatalf("failed test #%d\nexpected:\n%v\nactual:\n%v\n", testIndex+1, m, actual)
			}
		}
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		representative Monotyped[test_nameable]
		classMembers   mts_test
	}{
		{
			representative: _Con("Type"),
			classMembers:   mts_test{_Var("a")},
		},
		{
			representative: _Var("a"),
			classMembers:   mts_test{_Var("b")},
		},
		{ // test aliasing
			representative: _Con("Type"),
			classMembers:   mts_test{_Con("Type2")},
		},
		{
			representative: _Con("Type"),
			classMembers:   mts_test{_Con("Type")},
		},
		{
			representative: _Var("a"),
			classMembers:   mts_test{_Var("a")},
		},
		{
			representative: _Con("Type"),
			classMembers:   mts_test{_Var("a"), _Con("Type")},
		},
		{
			representative: _App("Type", _Con("Int")),
			classMembers: mts_test{_Var("a")},
		},
		{
			representative: _Var("a"),
			classMembers: mts_test{_App("Type", _Con("Int"))},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		if testIndex == 7 {
			print("here\n")
		}
		cxt.register(test.representative)(test.classMembers...)
		for _, m := range test.classMembers {
			name := m.String()
			actual, found := cxt.equivClasses[name]
			/*if !IsVariable(m) {
				if found {
					t.Fatalf("failed test #%d: %v was registered as a class member\n", testIndex+1, m)
				}
				continue
			}*/

			if !found {
				t.Fatalf("failed test #%d: %v wasn't registered\n", testIndex+1, m)
			}
			if !test.representative.Equals(actual) {
				t.Fatalf("failed test #%d\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.representative, actual)
			}
		}
	}
}

func TestUnion_valid(t *testing.T) {
	tests := []struct {
		ms             [2]Monotyped[test_nameable]
		representative Monotyped[test_nameable]
	}{
		{
			ms:             [2]Monotyped[test_nameable]{_Var("a"), _Var("b")},
			representative: _Var("b"),
		},
		{
			ms:             [2]Monotyped[test_nameable]{_Var("a"), _Con("Type")},
			representative: _Con("Type"),
		},
		{
			ms:             [2]Monotyped[test_nameable]{_Con("Type"), _Var("a")},
			representative: _Con("Type"),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		cxt.union(test.ms[0], test.ms[1])
		for mIndex, m := range test.ms {
			name := m.String()
			actual, found := cxt.equivClasses[name]
			/*if !IsVariable(m) {
				if found {
					t.Fatalf("failed test #%d.%d: %v was registered as a class member\n", testIndex+1, mIndex+1, m)
				}
				continue
			}*/

			if found {
				if !test.representative.Equals(actual) {
					t.Fatalf("failed test #%d.%d\nexpected:\n%v\nactual:\n%v\n", testIndex+1, mIndex+1, test.representative, actual)
				}
			} else {
				t.Fatalf("failed test #%d.%d: %v wasn't registered\n", testIndex+1, mIndex+1, m)
			}
		}
	}
}

func TestUnion_panic(t *testing.T) {
	tests := [][2]Monotyped[test_nameable]{
		{_Con("A"), _Con("B")},
		{_Con("A"), _Con("A")},
	}

	runUnionPanicTest := func(a, b Monotyped[test_nameable]) (passed bool) {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		defer func() {
			passed = recover() != nil
		}()
		cxt.union(a, b)
		return passed
	}

	for testIndex, test := range tests {
		if !runUnionPanicTest(test[0], test[1]) {
			t.Fatalf("failed test #%d: expected call to union to panic\n", testIndex+1)
		}
	}
}

func TestUnify_invalid(t *testing.T) {
	tests := [][2]Monotyped[test_nameable]{
		{_Con("Int"), _Con("Bool")},
		{_App("Type", _Con("x")), _App("Type", _Con("y"))},
		{_App("Type2", _Con("x")), _App("Type", _Con("x"))},
		{_App("Type", _Var("a")), _App("Type", _Con("y"))},
	}
	cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
	cxt.register(_Con("x"))(_Var("a")) // for fourth test

	for testIndex, test := range tests {
		expect := typeMismatch[test_nameable](test[0], test[1]).Error()
		e := cxt.Unify(test[0], test[1])
		if e == nil {
			t.Fatalf("test #%d failed: expected Unfify(%v, %v) == error(%s)\n",
				testIndex+1, test[0], test[1], expect)
		}

		if e.Error() != expect {
			t.Fatalf("test #%d failed:\nexpected:\n%s\nactual:\n%s\n",
				testIndex+1, expect, e.Error())
		}
	}
}

func _Function(left Monotyped[test_nameable], right Monotyped[test_nameable]) Application[test_nameable] {
	return base.Function(left, right)
}

func TestUnify_valid(t *testing.T) {
	tests := []struct {
		msSequence [][2]Monotyped[test_nameable]
		expected   [][2]Monotyped[test_nameable]
	}{
		{
			// x: Int
			// y: a = x
			msSequence: [][2]Monotyped[test_nameable]{
				{_Con("Int"), _Var("a")},
			},
			expected: [][2]Monotyped[test_nameable]{
				{_Var("a"), _Con("Int")},
			},
		},
		{
			// x: Int -> Bool
			// y: a = x
			// z: b -> c -> d = (\w: Int -> (y w))
			msSequence: [][2]Monotyped[test_nameable]{
				{
					_Function(_Con("Int"), _Con("Bool")),
					_Var("a"),
				},
				{
					_Function(_Var("b"), _Function(_Var("c"), _Var("d"))),
					_Function(_Con("Int"), _Function(_Con("Int"), _Con("Bool"))),
				},
			},
			expected: [][2]Monotyped[test_nameable]{
				{_Var("a"), _Function(_Con("Int"), _Con("Bool"))},
				{_Var("b"), _Con("Int")},
				{_Var("c"), _Con("Int")},
				{_Var("d"), _Con("Bool")},
			},
		},
		{
			// x: a
			// y: b = x
			msSequence: [][2]Monotyped[test_nameable]{
				{_Var("a"), _Var("b")},
			},
			expected: [][2]Monotyped[test_nameable]{
				{_Var("a"), _Var("b")},
				{_Var("b"), _Var("b")},
			},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		for unifyCallIndex, ms := range test.msSequence {
			if e := cxt.Unify(ms[0], ms[1]); e != nil {
				t.Fatalf("failed test #%d: Unify(ms[0]=%v, ms[1]=%v) call #%d returned error(%s)\n",
					testIndex+1, ms[0], ms[1], unifyCallIndex+1, e.Error())
			}
		}
		for findIndex, map_ := range test.expected {
			in, expected := map_[0], map_[1]
			name := in.String()
			if actual, found := cxt.equivClasses[name]; found {
				if !expected.Equals(actual) {
					t.Fatalf("failed test #%d (lookup #%d):\nexpected:\n%v\nactual:\n%v\n",
						testIndex+1, findIndex+1, expected, actual)
				}
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", cxt.equivClasses)
				t.Fatalf("failed test #%d (lookup #%d): lookup with %v found nothing\n%s",
					testIndex+1, findIndex+1, in, cxt.StringClasses())
			}
		}
	}
}

func TestDeclare(t *testing.T) {
	tests := []struct {
		in     Type[test_nameable]
		expect Monotyped[test_nameable]
	}{
		{
			in:     _Var("a"),
			expect: _Var("a"),
		},
		{
			in:     _Con("Type"),
			expect: _Con("Type"),
		},
		{
			in:     _App("Type", _Var("a")),
			expect: _App("Type", _Var("a")),
		},
		{
			in:     _Forall("a").Bind(_App("Type", _Var("b"))),
			expect: _App("Type", _Var("b")),
		},
		{
			in:     _Forall("a").Bind(_App("Type", _Var("a"))),
			expect: _App("Type", _Var("$0")),
		},
		{
			in:     _Forall("a", "b").Bind(_Function(_Var("b"), _Var("a"))),
			expect: _Function(_Var("$1"), _Var("$0")),
		},
	}
	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		if actual := cxt.Declare(test.in); !actual.Equals(test.expect) {
			t.Fatalf("test #%d failed:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect, actual)
		}
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		in     [2]Monotyped[test_nameable]
		expect [][2]Monotyped[test_nameable]
	}{
		{
			in:     [2]Monotyped[test_nameable]{_Var("a"), _Con("Int")},
			expect: [][2]Monotyped[test_nameable]{{_Var("a"), _Function(_Con("Int"), _Var("$0"))}},
		},
		{
			in:     [2]Monotyped[test_nameable]{_Var("a"), _Con("Int")},
			expect: [][2]Monotyped[test_nameable]{{_Var("a"), _Function(_Con("Int"), _Var("$0"))}},
		},
		{
			in:     [2]Monotyped[test_nameable]{_Function(_Con("Int"), _Con("Char")), _Con("Int")},
			expect: [][2]Monotyped[test_nameable]{{_Var("$0"), _Con("Char")}},
		},
		{
			in: [2]Monotyped[test_nameable]{_Function(_Var("$1"), _Con("Char")), _Con("Int")},
			expect: [][2]Monotyped[test_nameable]{
				{_Var("$0"), _Con("Char")},
				{_Var("$1"), _Con("Int")},
			},
		},
		{
			in: [2]Monotyped[test_nameable]{_Function(_Var("$1"), _Var("$2")), _Con("Int")},
			expect: [][2]Monotyped[test_nameable]{
				{_Var("$2"), _Var("$0")},
				{_Var("$1"), _Con("Int")},
			},
		},
	}

	expect := _Var("$0")
	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		res, e := cxt.Apply(test.in[0], test.in[1])
		if e != nil {
			t.Fatalf("test #%d failed with: %s\n", testIndex+1, e.Error())
		}
		if !res.Equals(expect) {
			t.Fatalf("test #%d failed:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, expect, res)
		}

		for mIndex, m := range test.expect {
			name := m[0].String()
			ty, found := cxt.equivClasses[name]
			if !found {
				t.Fatalf("test #%d failed on map #%d: could not find map\n", testIndex+1, mIndex+1)
			}
			if !ty.Equals(m[1]) {
				t.Fatalf("test #%d failed on map #%d:\nexpected:\n%v :-> %v\nactual:\n%v :-> %v\n",
					testIndex+1, mIndex+1,
					m[0], m[1],
					m[0], ty,
				)
			}
		}
	}
}

func TestApply_invalid(t *testing.T) {
	tests := []struct {
		in     [2]Monotyped[test_nameable]
		expect error
	}{
		{
			in:     [2]Monotyped[test_nameable]{_Con("Int"), _Con("Char")},
			expect: typeMismatch[test_nameable](_Con("Int"), _Function(_Con("Char"), _Var("$0"))),
		},
		{
			in:     [2]Monotyped[test_nameable]{_Function(_Con("Int"), _Con("Int")), _Con("Char")},
			expect: typeMismatch[test_nameable](_Function(_Con("Int"), _Con("Int")), _Function(_Con("Char"), _Var("$0"))),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		_, e := cxt.Apply(test.in[0], test.in[1])
		if e == nil {
			t.Fatalf("test #%d failed: expected cxt.Apply(%v, %v) to fail\n", testIndex+1, test.in[0], test.in[1])
		}
		if e.Error() != test.expect.Error() {
			t.Fatalf("test #%d failed:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect.Error(), e.Error())
		}
	}
}

func TestAbstract(t *testing.T) {
	tests := []struct {
		in     Monotyped[test_nameable]
		expect Monotyped[test_nameable]
	}{
		{
			in:     _Con("Int"),
			expect: _Function(_Var("$0"), _Con("Int")),
		},
		{
			in:     _Function(_Var("$1"), _Con("Int")),
			expect: _Function(_Var("$0"), _Function(_Var("$1"), _Con("Int"))),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		actual := cxt.Abstract(test.in)
		if !actual.Equals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect, actual)
		}
	}
}

func _Cons(left Monotyped[test_nameable], right Monotyped[test_nameable]) Application[test_nameable] {
	return base.Cons(left, right)
}

func TestCons(t *testing.T) {
	tests := []struct {
		in     [2]Monotyped[test_nameable]
		expect Monotyped[test_nameable]
	}{
		{
			in:     [2]Monotyped[test_nameable]{_Con("Int"), _Con("Int")},
			expect: _Cons(_Con("Int"), _Con("Int")),
		},
		{
			in:     [2]Monotyped[test_nameable]{_Con("Char"), _Con("Int")},
			expect: _Cons(_Con("Char"), _Con("Int")),
		},
		{
			in:     [2]Monotyped[test_nameable]{_Var("$0"), _Var("$1")},
			expect: _Cons(_Var("$0"), _Var("$1")),
		},
		{
			in:     [2]Monotyped[test_nameable]{_Function(_Var("$0"), _Var("$1")), _Var("$1")},
			expect: _Cons(_Function(_Var("$0"), _Var("$1")), _Var("$1")),
		},
	}

	cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
	for testIndex, test := range tests {
		actual := cxt.Cons(test.in[0], test.in[1])
		if !actual.Equals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect, actual)
		}
	}
}

func _testHead(t *testing.T, tests []struct {
	in     Monotyped[test_nameable]
	expect [][2]Monotyped[test_nameable]
}, check headTailFunc) {
	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		actual, e := cxt.Head(test.in)
		check(cxt, testIndex+1, actual, e, _Var("$0"), test.expect)
	}
}

func _testTail(t *testing.T, tests []struct {
	in     Monotyped[test_nameable]
	expect [][2]Monotyped[test_nameable]
}, check headTailFunc) {
	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		actual, e := cxt.Tail(test.in)
		check(cxt, testIndex+1, actual, e, _Var("$1"), test.expect)
	}
}

type headTailFunc func(cxt *Context[test_nameable], tn int, actual Monotyped[test_nameable], e error, expect Monotyped[test_nameable], maps [][2]Monotyped[test_nameable])

func mapCheck(t *testing.T, testName string) headTailFunc {
	return func(cxt *Context[test_nameable], tn int, actual Monotyped[test_nameable], e error, expect Monotyped[test_nameable], maps [][2]Monotyped[test_nameable]) {
		if e != nil {
			t.Fatalf("%s's test #%d failed with: %s\n", testName, tn, e.Error())
		}

		if !actual.Equals(expect) {
			t.Fatalf("%s's test #%d failed:\nexpected:\n%v\nactual:\n%v\n", testName, tn, expect, actual)
		}

		for mi, m := range maps {
			in, out := m[0], m[1]
			name, _ := Name(in)
			res, found := cxt.equivClasses[name]
			if !found {
				t.Fatalf("%s's test #%d failed on map #%d: map not found\n", testName, tn, mi+1)
			}
			if !res.Equals(out) {
				t.Fatalf("%s's test #%d failed on map #%d:\nexpected:\n%v :-> %v\nactual:\n%v :-> %v\n",
					testName,
					tn, mi+1,
					in, out,
					in, res,
				)
			}
		}
	}
}

func TestHeadAndTail(t *testing.T) {
	tests := []struct {
		in     Monotyped[test_nameable]
		expect [][2]Monotyped[test_nameable]
	}{
		{
			in: _Var("$8"),
			expect: [][2]Monotyped[test_nameable]{
				{_Var("$8"), _Cons(_Var("$0"), _Var("$1"))},
			},
		},
		{
			in: _Cons(_Con("Int"), _Con("Char")),
			expect: [][2]Monotyped[test_nameable]{
				{_Var("$0"), _Con("Int")},
				{_Var("$1"), _Con("Char")},
			},
		},
		{
			in: _Cons(_App("Type", _Var("$8")), _Con("Char")),
			expect: [][2]Monotyped[test_nameable]{
				{_Var("$0"), _App("Type", _Var("$8"))},
				{_Var("$1"), _Con("Char")},
			},
		},
	}

	_testHead(t, tests, mapCheck(t, "TestHead"))
	_testTail(t, tests, mapCheck(t, "TestTail"))
}

func TestHeadTail_invalid(t *testing.T) {
	type test_type struct {
		in     Monotyped[test_nameable]
		expect error
	}
	tests := []test_type{
		{
			in:     _Con("Int"),
			expect: typeMismatch[test_nameable](_Con("Int"), _Cons(_Var("$0"), _Var("$1"))),
		},
		{
			in:     _Function(_Con("Int"), _Con("Char")),
			expect: typeMismatch[test_nameable](_Function(_Con("Int"), _Con("Char")), _Cons(_Var("$0"), _Var("$1"))),
		},
	}

	check := func(name string, testIndex int, test test_type, res Monotyped[test_nameable], e error) {
		if res != nil {
			t.Fatalf("%s's test #%d failed: type result was not nil\n", name, testIndex+1)
		}
		if e == nil {
			t.Fatalf("%s's test #%d failed: error result was nil", name, testIndex+1)
		}
		if e.Error() != test.expect.Error() {
			t.Fatalf("%s's test #%d failed:\nexpected:\n%v\nactual:\n%v\n", name, testIndex+1, test.expect.Error(), e.Error())
		}
	}

	for testIndex, test := range tests {
		res, e := NewContext[test_nameable]().SetNameMaker(test_nameable_fn).Head(test.in)
		check("TestHead_invalid", testIndex, test, res, e)
		res, e = NewContext[test_nameable]().SetNameMaker(test_nameable_fn).Tail(test.in)
		check("TestTail_invalid", testIndex, test, res, e)
	}
}

func _Join(left Monotyped[test_nameable], right Monotyped[test_nameable]) Application[test_nameable] {
	return base.Join(left, right)
}

func TestJoin(t *testing.T) {
	tests := []struct {
		in     Monotyped[test_nameable]
		expect Monotyped[test_nameable]
	}{
		{
			in:     _Var("$1"),
			expect: _Join(_Var("$1"), _Var("$0")),
		},
		{
			in:     _Con("Int"),
			expect: _Join(_Con("Int"), _Var("$0")),
		},
		{
			in:     _Join(_Con("Int"), _Var("$1")),
			expect: _Join(_Join(_Con("Int"), _Var("$1")), _Var("$0")),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		actual := cxt.JoinRule(test.in)
		if !actual.Equals(test.expect) {
			t.Fatalf("test #%d failed:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect, actual)
		}
	}
}

func TestRealization(t *testing.T) {
	tests := []struct {
		in     [3]Monotyped[test_nameable]
		expect [][2]Monotyped[test_nameable]
	}{
		{
			in: [3]Monotyped[test_nameable]{
				_Join(_Con("Int"), _Con("Char")),
				_Function(_Con("Int"), _Con("Char")),
				_Function(_Con("Char"), _Con("Char")),
			},
			expect: [][2]Monotyped[test_nameable]{
				{_Var("$0"), _Con("Int")},
				{_Var("$1"), _Con("Char")},
				{_Var("$2"), _Con("Char")},
			},
		},
		{
			in: [3]Monotyped[test_nameable]{
				_Var("$8"),
				_Function(_Con("Int"), _Con("Char")),
				_Function(_Con("Char"), _Con("Char")),
			},
			expect: [][2]Monotyped[test_nameable]{
				{_Var("$0"), _Con("Int")},
				{_Var("$1"), _Con("Char")},
				{_Var("$2"), _Con("Char")},
				{_Var("$8"), _Join(_Var("$0"), _Var("$1"))},
			},
		},
		{
			in: [3]Monotyped[test_nameable]{
				_Join(_Var("$7"), _Var("$8")),
				_Function(_Con("Int"), _Con("Char")),
				_Function(_Con("Char"), _Con("Char")),
			},
			expect: [][2]Monotyped[test_nameable]{
				{_Var("$0"), _Con("Int")},
				{_Var("$1"), _Con("Char")},
				{_Var("$2"), _Con("Char")},
				{_Var("$7"), _Var("$0")},
				{_Var("$8"), _Var("$1")},
			},
		},
	}

	expect := _Var("$2")
	check := mapCheck(t, "TestUn_Join")

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		actual, e := cxt.Realization(test.in[0], test.in[1], test.in[2])
		check(cxt, testIndex+1, actual, e, expect, test.expect)
	}
}

func TestRealization_invalid(t *testing.T) {
	tests := []struct {
		in     [3]Monotyped[test_nameable]
		expect error
	}{
		{
			in: [3]Monotyped[test_nameable]{
				_Con("Int"),
				_Function(_Con("Int"), _Con("Char")),
				_Function(_Con("Char"), _Con("Char")),
			},
			expect: typeMismatch[test_nameable](_Con("Int"), _Join(_Var("$0"), _Var("$1"))),
		},
		{
			in: [3]Monotyped[test_nameable]{
				_Join(_Con("Int"), _Con("Char")),
				_Con("Int"),
				_Function(_Con("Char"), _Con("Char")),
			},
			expect: typeMismatch[test_nameable](_Con("Int"), _Function(_Var("$0"), _Var("$2"))),
		},
		{
			in: [3]Monotyped[test_nameable]{
				_Join(_Con("Int"), _Con("Char")),
				_Function(_Con("Int"), _Con("Char")),
				_Con("Char"),
			},
			expect: typeMismatch[test_nameable](_Con("Char"), _Function(_Var("$1"), _Var("$2"))),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_nameable]().SetNameMaker(test_nameable_fn)
		res, e := cxt.Realization(test.in[0], test.in[1], test.in[2])
		if res != nil {
			t.Fatalf("test #%d failed: res was not nil\n", testIndex+1)
		}
		if e == nil {
			t.Fatalf("test #%d failed: error returned was nil\n", testIndex+1)
		}
		if e.Error() != test.expect.Error() {
			t.Fatalf("test #%d failed:\nexpected:\n%v\nactual:\n%v\n",
				testIndex+1,
				test.expect.Error(),
				e.Error(),
			)
		}
	}
}
