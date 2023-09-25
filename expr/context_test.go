package expr

import "testing"

func TestDeclareInverses_invalid(t *testing.T) {
	f := Const("f")
	fInv := Const("fInv")

	tests := []struct {
		setup  func(*Context)
		expect string
	}{
		{
			func(cxt *Context) {}, 
			nameNotDefined(f).Error(),
		},
		{
			func(cxt *Context) { cxt.table[f] = f }, 
			nameNotDefined(fInv).Error(),
		},
		{
			func(cxt *Context) { 
				cxt.table[f], cxt.table[fInv] = f, fInv 
				cxt.inverses[f] = Const("_")
			}, 
			redefineInv(f).Error(),
		},
		{
			func(cxt *Context) { 
				cxt.table[f], cxt.table[fInv] = f, fInv 
				cxt.inverses[fInv] = Const("_")
			}, 
			redefineInv(fInv).Error(),
		},
	}
	
	for testIndex, test := range tests {
		cxt := NewContext()
		test.setup(cxt)
		e := cxt.DeclareInverse(f, fInv)
		
		if e == nil {
			t.Fatalf("failed test #%d: call to cxt.DeclareInverse(%v, %v) succeeded but should have failed with \"%s.\"\n", testIndex+1, f, fInv, test.expect)
		}
		actual := e.Error()
		if actual != test.expect {
			t.Fatalf("failed test #%d:\nexpected:\n%s\nactual:\n%s\n", testIndex+1, test.expect, actual)
		}
	}
}

func TestDeclareInverses(t *testing.T) {
	tests := []struct{
		names [2]Const
		exp [2]Expression
	}{
		{
			[2]Const{Const("id"), Const("id")},
			[2]Expression{Const("id"), Const("id")},
		},
		{
			[2]Const{Const("succ"), Const("pred")}, 
			[2]Expression{
				Bind(Var("x")).In(Apply(Const("Succ"), Var("x"))),
				Bind(Var("x")).In(Apply(Const("Pred"), Var("x"))),
			},
		},
		{
			[2]Const{Const("(+)"), Const("(-)")}, 
			[2]Expression{
				Bind(Var("x"), Var("y")).In(Apply(Const("addNum"), Var("x"), Var("y"))),
				Bind(Var("x"), Var("y")).In(Apply(Const("subNum"), Var("x"), Var("y"))),
			},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext()

		// add symbols
		cxt.table[test.names[0]] = test.exp[0]
		cxt.table[test.names[1]] = test.exp[1]

		// do declaration
		e := cxt.DeclareInverse(test.names[0], test.names[1])
		if e != nil {
			t.Fatalf("failed test #%d: call to cxt.DeclareInverse(%v, %v) failed with \"%s.\"\n", testIndex+1, test.names[0], test.names[1], e.Error())
		}

		// for test.names[0]
		actual, found := cxt.inverses[test.names[0]]
		if !found {
			t.Fatalf("failed test #%d: could not find cxt.inverses[%v].\n", testIndex+1, test.names[0])
		}
		if !test.names[1].StrictEquals(actual) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.names[1], actual)
		}

		// for test.names[1]
		actual, found = cxt.inverses[test.names[1]]
		if !found {
			t.Fatalf("failed test #%d: could not find cxt.inverses[%v].\n", testIndex+1, test.names[1])
		}
		if !test.names[0].StrictEquals(actual) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.names[0], actual)
		}
	}
}

func TestAddName_invalid(t *testing.T) {
	tests := []struct {
		name   Const
		expect string
	}{
		{Const("id"), redefineNameInTable(Const("id")).Error()},
	}
	testDummy := Const("dummy")
	for testIndex, test := range tests {
		cxt := NewContext()
		cxt.table[test.name] = testDummy

		e := cxt.AddName(test.name, testDummy)
		if e == nil {
			t.Fatalf("failed test #%d: call to cxt.AddName(%v, %v) succeeded but should have failed with \"%s.\"\n", testIndex+1, test.name, testDummy, test.expect)
		}
		actual := e.Error()
		if actual != test.expect {
			t.Fatalf("failed test #%d:\nexpected:\n%s\nactual:\n%s\n", testIndex+1, test.expect, actual)
		}
	}
}

func TestAddName(t *testing.T) {
	tests := []struct {
		name Const
		exp  Expression
	}{
		{Const("0"), Const("0")},
		{Const("1"), Apply(Const("Succ"), Const("0"))},
		{Const("id"), IdFunction},
	}

	for testIndex, test := range tests {
		cxt := NewContext()

		e := cxt.AddName(test.name, test.exp)
		if e != nil {
			t.Fatalf("failed test #%d: call to cxt.AddName(%v, %v) failed with \"%s.\"\n", testIndex+1, test.name, test.exp, e.Error())
		}
		actual, found := cxt.table[test.name]
		if !found {
			t.Fatalf("failed test #%d: could not find cxt.table[%v].\n", testIndex+1, test.name)
		}
		if !test.exp.StrictEquals(actual) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.exp, actual)
		}
	}
}

func TestGetInverse(t *testing.T) {
	tests := []struct {
		f    Const
		fInv Const
	}{
		{Const("id"), Const("id")},
		{Const("f"), Const("f^-1")},
	}

	for testIndex, test := range tests {
		cxt := NewContext()

		cxt.inverses[test.f] = test.fInv
		cxt.table[test.f] = test.f
		cxt.inverses[test.fInv] = test.f
		cxt.table[test.fInv] = test.fInv
		actual, ok := cxt.GetInverse(test.f)
		if !ok {
			t.Fatalf("failed test #%d: call to cxt.GetInverse(%v) failed.\n", testIndex+1, test.f)
		}
		if !test.fInv.StrictEquals(actual) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.fInv, actual)
		}

		actual, ok = cxt.GetInverse(test.fInv)
		if !ok {
			t.Fatalf("failed test #%d: call to cxt.GetInverse(%v) failed.\n", testIndex+1, test.fInv)
		}
		if !test.f.StrictEquals(actual) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.f, actual)
		}
	}
}
