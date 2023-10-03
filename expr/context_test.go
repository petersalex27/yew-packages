package expr

import "testing"

func nameMaker(s string) test_named {
		return test_named(s)
}

func TestDeclareInverses_invalid(t *testing.T) {
	f := _Const("f")
	fInv := _Const("fInv")

	tests := []struct {
		setup  func(*Context[test_named])
		expect string
	}{
		{
			func(cxt *Context[test_named]) {}, 
			nameNotDefined(f).Error(),
		},
		{
			func(cxt *Context[test_named]) { cxt.table[f.String()] = f }, 
			nameNotDefined(fInv).Error(),
		},
		{
			func(cxt *Context[test_named]) { 
				cxt.table[f.String()], cxt.table[fInv.String()] = f, fInv 
				cxt.inverses[f.String()] = _Const("_")
			}, 
			redefineInv(f).Error(),
		},
		{
			func(cxt *Context[test_named]) { 
				cxt.table[f.String()], cxt.table[fInv.String()] = f, fInv 
				cxt.inverses[fInv.String()] = _Const("_")
			}, 
			redefineInv(fInv).Error(),
		},
	}
	
	for testIndex, test := range tests {
		cxt := NewContext[test_named]().SetNameMaker(nameMaker)
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
		names [2]Const[test_named]
		exp [2]Expression[test_named]
	}{
		{
			[2]Const[test_named]{_Const("id"), _Const("id")},
			[2]Expression[test_named]{_Const("id"), _Const("id")},
		},
		{
			[2]Const[test_named]{_Const("succ"), _Const("pred")}, 
			[2]Expression[test_named]{
				_Bind(_Var("x")).In(_Apply(_Const("Succ"), _Var("x"))),
				_Bind(_Var("x")).In(_Apply(_Const("Pred"), _Var("x"))),
			},
		},
		{
			[2]Const[test_named]{_Const("(+)"), _Const("(-)")}, 
			[2]Expression[test_named]{
				_Bind(_Var("x"), _Var("y")).In(_Apply(_Const("addNum"), _Var("x"), _Var("y"))),
				_Bind(_Var("x"), _Var("y")).In(_Apply(_Const("subNum"), _Var("x"), _Var("y"))),
			},
		},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_named]().SetNameMaker(nameMaker)

		// add symbols
		cxt.table[test.names[0].String()] = test.exp[0]
		cxt.table[test.names[1].String()] = test.exp[1]

		// do declaration
		e := cxt.DeclareInverse(test.names[0], test.names[1])
		if e != nil {
			t.Fatalf("failed test #%d: call to cxt.DeclareInverse(%v, %v) failed with \"%s.\"\n", testIndex+1, test.names[0], test.names[1], e.Error())
		}

		// for test.names[0]
		actual, found := cxt.inverses[test.names[0].String()]
		if !found {
			t.Fatalf("failed test #%d: could not find cxt.inverses[%v].\n", testIndex+1, test.names[0])
		}
		if !test.names[1].StrictEquals(actual) {
			t.Fatalf("failed test #%d:\nexpected:\n%v\nactual:\n%v\n", testIndex+1, test.names[1], actual)
		}

		// for test.names[1]
		actual, found = cxt.inverses[test.names[1].String()]
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
		name   Const[test_named]
		expect string
	}{
		{_Const("id"), redefineNameInTable(_Const("id")).Error()},
	}
	testDummy := _Const("dummy")
	for testIndex, test := range tests {
		cxt := NewContext[test_named]().SetNameMaker(nameMaker)
		cxt.table[test.name.String()] = testDummy

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
		name Const[test_named]
		exp  Expression[test_named]
	}{
		{_Const("0"), _Const("0")},
		{_Const("1"), _Apply(_Const("Succ"), _Const("0"))},
		{_Const("id"), IdFunction},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_named]().SetNameMaker(nameMaker)

		e := cxt.AddName(test.name, test.exp)
		if e != nil {
			t.Fatalf("failed test #%d: call to cxt.AddName(%v, %v) failed with \"%s.\"\n", testIndex+1, test.name, test.exp, e.Error())
		}
		actual, found := cxt.table[test.name.String()]
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
		f    Const[test_named]
		fInv Const[test_named]
	}{
		{_Const("id"), _Const("id")},
		{_Const("f"), _Const("f^-1")},
	}

	for testIndex, test := range tests {
		cxt := NewContext[test_named]().SetNameMaker(nameMaker)

		cxt.inverses[test.f.String()] = test.fInv
		cxt.table[test.f.String()] = test.f
		cxt.inverses[test.fInv.String()] = test.f
		cxt.table[test.fInv.String()] = test.fInv
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
