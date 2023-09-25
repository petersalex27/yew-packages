package types

import (
	expr "github.com/petersalex27/yew-packages/expression"
	"fmt"
	"os"
	"testing"
)


func TestString(t *testing.T) {
	tests := []struct {
		in     Type
		expect string
	}{
		// monotypes
		{in: Con("Type"), expect: "Type"},                                     // just type
		{in: App("Type", Con("x")), expect: "(Type x)"},                      // application
		{in: App("Type", App("Type", Con("x"))), expect: "(Type (Type x))"}, // nested application
		{in: Var("a"), expect: "a"},                                             // variable
		{in: App("Type", Var("a")), expect: "(Type a)"},                        // application w/ variable
		{in: Infix(Con("Type"), "->", Con("Type")), expect: "(Type -> Type)"},
		{in: Infix(App("Color", Con("Red")), "->", Con("Type")), expect: "((Color Red) -> Type)"},
		{in: Infix(Con("Type"), "->", App("Color", Con("Red"))), expect: "(Type -> (Color Red))"},
		{in: Infix(Con("Type"), "->", Infix(Con("Type"), "->", Con("Type"))), expect: "(Type -> (Type -> Type))"},
		{in: Infix(Infix(Con("Type"), "->", Con("Type")), "->", Con("Type")), expect: "((Type -> Type) -> Type)"},
		{in: Infix(Var("a"), "->", Var("a")), expect: "(a -> a)",},
		{in: Infix(Var("a"), "->"), expect: "((->) a)",},
		{in: InfixApplication{c: Constant("->"), ts: nil,}, expect: "(->)",},
		// polytypes
		{in: Con("Type").Generalize(), expect: "forall _ . Type"},
		{in: Forall("a").Bind(Con("Type")), expect: "forall a . Type"},
		{in: Forall("a", "b").Bind(Con("Type")), expect: "forall a b . Type"},
		{in: Forall("a", "b").Bind(App("Type", App("Type2", Var("b"), Var("a")))), expect: "forall a b . (Type (Type2 b a))"},
		{in: Forall("a").Bind(Var("a")), expect: "forall a . a"},
		// dependent types
		{
			in: Index(App("Array", Con("Int")), Judgement[expr.Expression](expr.Var("n"), Con("Uint"))),
			expect: "((Array Int); (n: Uint))",
		},
	}

	for testIndex, test := range tests {
		if actual := test.in.String(); actual != test.expect {
			t.Fatalf("failed test #%d:\nexpected:\n%s\nactual:\n%s\n", testIndex+1, test.expect, actual)
		}
	}
}

func TestInstantiate(t *testing.T) {
	tests := []struct {
		poly   Polytype
		mono   Monotyped
		expect Type
	}{
		{ // (forall a . Type a) $ b == Type b
			poly:   Forall("a").Bind(App("Type", Var("a"))),
			mono:   Var("b"),
			expect: App("Type", Var("b")),
		},
		{ // (forall a . Type a) $ Louis == Type Louis
			poly:   Forall("a").Bind(App("Type", Var("a"))),
			mono:   Con("Louis"),
			expect: App("Type", Con("Louis")),
		},
		{ // (forall a . Type b) $ c == Type b
			poly:   Forall("a").Bind(App("Type", Var("b"))),
			mono:   Var("c"),
			expect: App("Type", Var("b")),
		},
		{ // (forall a b c . Type a b c) $ x == forall b c . Type x b c
			poly:   Forall("a", "b", "c").Bind(App("Type", Var("a"), Var("b"), Var("c"))),
			mono:   Var("x"),
			expect: Forall("b", "c").Bind(App("Type", Var("x"), Var("b"), Var("c"))),
		},
		{ // (forall b c . Type x b c) $ x == forall c . Type x x c
			poly:   Forall("b", "c").Bind(App("Type", Var("x"), Var("b"), Var("c"))),
			mono:   Var("x"),
			expect: Forall("c").Bind(App("Type", Var("x"), Var("x"), Var("c"))),
		},
		{ // (forall c . Type x x c) $ x == Type x x x
			poly:   Forall("c").Bind(App("Type", Var("x"), Var("x"), Var("c"))),
			mono:   Var("x"),
			expect: App("Type", Var("x"), Var("x"), Var("x")),
		},
		{
			poly:	Forall("a").Bind(Index(App("Array", Var("a")), Judgement[expr.Expression](expr.Var("n"), Con("Uint")))),
			mono:	Con("Int"),
			expect: Index(App("Array", Con("Int")), Judgement[expr.Expression](expr.Var("n"), Con("Uint"))),
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
func TestFind_inTable(t *testing.T) {
	tests := []struct {
		representative Monotyped
		classMembers   []Monotyped
	}{
		{
			representative: Con("Type"),
			classMembers:   []Monotyped{Var("a")},
		},
		{
			representative: Var("a"),
			classMembers:   []Monotyped{Var("b")},
		},
		{
			representative: Con("Type"),
			classMembers:   []Monotyped{Con("Type2")},
		},
		{
			representative: Con("Type"),
			classMembers:   []Monotyped{Con("Type")},
		},
		{
			representative: Var("a"),
			classMembers:   []Monotyped{Var("a")},
		},
		{
			representative: Con("Type"),
			classMembers:   []Monotyped{Var("a"), Con("Type")},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()
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
		classMembers []Monotyped
	}{
		{
			classMembers: []Monotyped{Var("a")},
		},
		{
			classMembers: []Monotyped{Con("Type")},
		},
		{
			classMembers: []Monotyped{Var("a"), Con("Type")},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()
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
		representative Monotyped
		classMembers   []Monotyped
	}{
		{
			representative: Con("Type"),
			classMembers:   []Monotyped{Var("a")},
		},
		{
			representative: Var("a"),
			classMembers:   []Monotyped{Var("b")},
		},
		{ // test aliasing
			representative: Con("Type"),
			classMembers:   []Monotyped{Con("Type2")},
		},
		{
			representative: Con("Type"),
			classMembers:   []Monotyped{Con("Type")},
		},
		{
			representative: Var("a"),
			classMembers:   []Monotyped{Var("a")},
		},
		{
			representative: Con("Type"),
			classMembers:   []Monotyped{Var("a"), Con("Type")},
		},
		{
			representative: App("Type", Con("Int")),
			classMembers: []Monotyped{Var("a")},
		},
		{
			representative: Var("a"),
			classMembers: []Monotyped{App("Type", Con("Int"))},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()
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
		ms             [2]Monotyped
		representative Monotyped
	}{
		{
			ms:             [2]Monotyped{Var("a"), Var("b")},
			representative: Var("b"),
		},
		{
			ms:             [2]Monotyped{Var("a"), Con("Type")},
			representative: Con("Type"),
		},
		{
			ms:             [2]Monotyped{Con("Type"), Var("a")},
			representative: Con("Type"),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()
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
	tests := [][2]Monotyped{
		{Con("A"), Con("B")},
		{Con("A"), Con("A")},
	}

	runUnionPanicTest := func(a, b Monotyped) (passed bool) {
		cxt := NewContext()
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
	tests := [][2]Monotyped{
		{Con("Int"), Con("Bool")},
		{App("Type", Con("x")), App("Type", Con("y"))},
		{App("Type2", Con("x")), App("Type", Con("x"))},
		{App("Type", Var("a")), App("Type", Con("y"))},
	}
	cxt := NewContext()
	cxt.register(Con("x"))(Var("a")) // for fourth test

	for testIndex, test := range tests {
		expect := typeMismatch(test[0], test[1]).Error()
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

func TestUnify_valid(t *testing.T) {
	tests := []struct {
		msSequence [][2]Monotyped
		expected   [][2]Monotyped
	}{
		{
			// x: Int
			// y: a = x
			msSequence: [][2]Monotyped{
				{Con("Int"), Var("a")},
			},
			expected: [][2]Monotyped{
				{Var("a"), Con("Int")},
			},
		},
		{
			// x: Int -> Bool
			// y: a = x
			// z: b -> c -> d = (\w: Int -> (y w))
			msSequence: [][2]Monotyped{
				{
					Function(Con("Int"), Con("Bool")),
					Var("a"),
				},
				{
					Function(Var("b"), Function(Var("c"), Var("d"))),
					Function(Con("Int"), Function(Con("Int"), Con("Bool"))),
				},
			},
			expected: [][2]Monotyped{
				{Var("a"), Function(Con("Int"), Con("Bool"))},
				{Var("b"), Con("Int")},
				{Var("c"), Con("Int")},
				{Var("d"), Con("Bool")},
			},
		},
		{
			// x: a
			// y: b = x
			msSequence: [][2]Monotyped{
				{Var("a"), Var("b")},
			},
			expected: [][2]Monotyped{
				{Var("a"), Var("b")},
				{Var("b"), Var("b")},
			},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()
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
		in     Type
		expect Monotyped
	}{
		{
			in:     Var("a"),
			expect: Var("a"),
		},
		{
			in:     Con("Type"),
			expect: Con("Type"),
		},
		{
			in:     App("Type", Var("a")),
			expect: App("Type", Var("a")),
		},
		{
			in:     Forall("a").Bind(App("Type", Var("b"))),
			expect: App("Type", Var("b")),
		},
		{
			in:     Forall("a").Bind(App("Type", Var("a"))),
			expect: App("Type", Var("$0")),
		},
		{
			in:     Forall("a", "b").Bind(Function(Var("b"), Var("a"))),
			expect: Function(Var("$1"), Var("$0")),
		},
	}
	for testIndex, test := range tests {
		cxt := NewContext()
		if actual := cxt.Declare(test.in); !actual.Equals(test.expect) {
			t.Fatalf("test #%d failed:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect, actual)
		}
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		in     [2]Monotyped
		expect [][2]Monotyped
	}{
		{
			in:     [2]Monotyped{Var("a"), Con("Int")},
			expect: [][2]Monotyped{{Var("a"), Function(Con("Int"), Var("$0"))}},
		},
		{
			in:     [2]Monotyped{Var("a"), Con("Int")},
			expect: [][2]Monotyped{{Var("a"), Function(Con("Int"), Var("$0"))}},
		},
		{
			in:     [2]Monotyped{Function(Con("Int"), Con("Char")), Con("Int")},
			expect: [][2]Monotyped{{Var("$0"), Con("Char")}},
		},
		{
			in: [2]Monotyped{Function(Var("$1"), Con("Char")), Con("Int")},
			expect: [][2]Monotyped{
				{Var("$0"), Con("Char")},
				{Var("$1"), Con("Int")},
			},
		},
		{
			in: [2]Monotyped{Function(Var("$1"), Var("$2")), Con("Int")},
			expect: [][2]Monotyped{
				{Var("$2"), Var("$0")},
				{Var("$1"), Con("Int")},
			},
		},
	}

	expect := Var("$0")
	for testIndex, test := range tests {
		cxt := NewContext()
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
		in     [2]Monotyped
		expect error
	}{
		{
			in:     [2]Monotyped{Con("Int"), Con("Char")},
			expect: typeMismatch(Con("Int"), Function(Con("Char"), Var("$0"))),
		},
		{
			in:     [2]Monotyped{Function(Con("Int"), Con("Int")), Con("Char")},
			expect: typeMismatch(Function(Con("Int"), Con("Int")), Function(Con("Char"), Var("$0"))),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()
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
		in     Monotyped
		expect Monotyped
	}{
		{
			in:     Con("Int"),
			expect: Function(Var("$0"), Con("Int")),
		},
		{
			in:     Function(Var("$1"), Con("Int")),
			expect: Function(Var("$0"), Function(Var("$1"), Con("Int"))),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()
		actual := cxt.Abstract(test.in)
		if !actual.Equals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect, actual)
		}
	}
}

func TestCons(t *testing.T) {
	tests := []struct {
		in     [2]Monotyped
		expect Monotyped
	}{
		{
			in:     [2]Monotyped{Con("Int"), Con("Int")},
			expect: Cons(Con("Int"), Con("Int")),
		},
		{
			in:     [2]Monotyped{Con("Char"), Con("Int")},
			expect: Cons(Con("Char"), Con("Int")),
		},
		{
			in:     [2]Monotyped{Var("$0"), Var("$1")},
			expect: Cons(Var("$0"), Var("$1")),
		},
		{
			in:     [2]Monotyped{Function(Var("$0"), Var("$1")), Var("$1")},
			expect: Cons(Function(Var("$0"), Var("$1")), Var("$1")),
		},
	}

	cxt := NewContext()
	for testIndex, test := range tests {
		actual := cxt.Cons(test.in[0], test.in[1])
		if !actual.Equals(test.expect) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect, actual)
		}
	}
}

func _testHead(t *testing.T, tests []struct {
	in     Monotyped
	expect [][2]Monotyped
}, check headTailFunc) {
	for testIndex, test := range tests {
		cxt := NewContext()
		actual, e := cxt.Head(test.in)
		check(cxt, testIndex+1, actual, e, Var("$0"), test.expect)
	}
}

func _testTail(t *testing.T, tests []struct {
	in     Monotyped
	expect [][2]Monotyped
}, check headTailFunc) {
	for testIndex, test := range tests {
		cxt := NewContext()
		actual, e := cxt.Tail(test.in)
		check(cxt, testIndex+1, actual, e, Var("$1"), test.expect)
	}
}

type headTailFunc func(cxt *Context, tn int, actual Monotyped, e error, expect Monotyped, maps [][2]Monotyped)

func mapCheck(t *testing.T, testName string) headTailFunc {
	return func(cxt *Context, tn int, actual Monotyped, e error, expect Monotyped, maps [][2]Monotyped) {
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
		in     Monotyped
		expect [][2]Monotyped
	}{
		{
			in: Var("$8"),
			expect: [][2]Monotyped{
				{Var("$8"), Cons(Var("$0"), Var("$1"))},
			},
		},
		{
			in: Cons(Con("Int"), Con("Char")),
			expect: [][2]Monotyped{
				{Var("$0"), Con("Int")},
				{Var("$1"), Con("Char")},
			},
		},
		{
			in: Cons(App("Type", Var("$8")), Con("Char")),
			expect: [][2]Monotyped{
				{Var("$0"), App("Type", Var("$8"))},
				{Var("$1"), Con("Char")},
			},
		},
	}

	_testHead(t, tests, mapCheck(t, "TestHead"))
	_testTail(t, tests, mapCheck(t, "TestTail"))
}

func TestHeadTail_invalid(t *testing.T) {
	type test_type struct {
		in     Monotyped
		expect error
	}
	tests := []test_type{
		{
			in:     Con("Int"),
			expect: typeMismatch(Con("Int"), Cons(Var("$0"), Var("$1"))),
		},
		{
			in:     Function(Con("Int"), Con("Char")),
			expect: typeMismatch(Function(Con("Int"), Con("Char")), Cons(Var("$0"), Var("$1"))),
		},
	}

	check := func(name string, testIndex int, test test_type, res Monotyped, e error) {
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
		res, e := NewContext().Head(test.in)
		check("TestHead_invalid", testIndex, test, res, e)
		res, e = NewContext().Tail(test.in)
		check("TestTail_invalid", testIndex, test, res, e)
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		in     Monotyped
		expect Monotyped
	}{
		{
			in:     Var("$1"),
			expect: Join(Var("$1"), Var("$0")),
		},
		{
			in:     Con("Int"),
			expect: Join(Con("Int"), Var("$0")),
		},
		{
			in:     Join(Con("Int"), Var("$1")),
			expect: Join(Join(Con("Int"), Var("$1")), Var("$0")),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()
		actual := cxt.Join(test.in)
		if !actual.Equals(test.expect) {
			t.Fatalf("test #%d failed:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.expect, actual)
		}
	}
}

func TestRealization(t *testing.T) {
	tests := []struct {
		in     [3]Monotyped
		expect [][2]Monotyped
	}{
		{
			in: [3]Monotyped{
				Join(Con("Int"), Con("Char")),
				Function(Con("Int"), Con("Char")),
				Function(Con("Char"), Con("Char")),
			},
			expect: [][2]Monotyped{
				{Var("$0"), Con("Int")},
				{Var("$1"), Con("Char")},
				{Var("$2"), Con("Char")},
			},
		},
		{
			in: [3]Monotyped{
				Var("$8"),
				Function(Con("Int"), Con("Char")),
				Function(Con("Char"), Con("Char")),
			},
			expect: [][2]Monotyped{
				{Var("$0"), Con("Int")},
				{Var("$1"), Con("Char")},
				{Var("$2"), Con("Char")},
				{Var("$8"), Join(Var("$0"), Var("$1"))},
			},
		},
		{
			in: [3]Monotyped{
				Join(Var("$7"), Var("$8")),
				Function(Con("Int"), Con("Char")),
				Function(Con("Char"), Con("Char")),
			},
			expect: [][2]Monotyped{
				{Var("$0"), Con("Int")},
				{Var("$1"), Con("Char")},
				{Var("$2"), Con("Char")},
				{Var("$7"), Var("$0")},
				{Var("$8"), Var("$1")},
			},
		},
	}

	expect := Var("$2")
	check := mapCheck(t, "TestUnjoin")

	for testIndex, test := range tests {
		cxt := NewContext()
		actual, e := cxt.Realization(test.in[0], test.in[1], test.in[2])
		check(cxt, testIndex+1, actual, e, expect, test.expect)
	}
}

func TestRealization_invalid(t *testing.T) {
	tests := []struct {
		in     [3]Monotyped
		expect error
	}{
		{
			in: [3]Monotyped{
				Con("Int"),
				Function(Con("Int"), Con("Char")),
				Function(Con("Char"), Con("Char")),
			},
			expect: typeMismatch(Con("Int"), Join(Var("$0"), Var("$1"))),
		},
		{
			in: [3]Monotyped{
				Join(Con("Int"), Con("Char")),
				Con("Int"),
				Function(Con("Char"), Con("Char")),
			},
			expect: typeMismatch(Con("Int"), Function(Var("$0"), Var("$2"))),
		},
		{
			in: [3]Monotyped{
				Join(Con("Int"), Con("Char")),
				Function(Con("Int"), Con("Char")),
				Con("Char"),
			},
			expect: typeMismatch(Con("Char"), Function(Var("$1"), Var("$2"))),
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()
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
