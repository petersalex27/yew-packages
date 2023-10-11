package expr

import (
	"strings"

	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

// λx.e
type Function[T nameable.Nameable] struct {
	vars []Variable[T]
	e    Expression[T]
}

func (f Function[T]) ExtractFreeVariables(dummyVar Variable[T]) []Variable[T] {
	var e Expression[T] = f.e
	for _, v := range f.vars {
		e, _ = e.Replace(v, dummyVar)
	}
	
	return e.ExtractFreeVariables(dummyVar)
}

func (f Function[T]) Collect() []T {
	res := make([]T, 0, 1)
	for _, v := range f.vars {
		res = append(res, v.Collect()...)
	}
	res = append(res, f.e.Collect()...)
	return res
}

var x_, y_ = Var(test_named("x")), Var(test_named("y"))
var t_, f_ = Var(test_named("t")), Var(test_named("f"))
var a_, b_ = Var(test_named("a")), Var(test_named("b"))
var c_, e_ = Var(test_named("c")), Var(test_named("e"))

// (λx . x)
var IdFunction = Bind[test_named](x_).In(x_)

// (λx . y)
var ConstFunction = Bind[test_named](x_).In(y_)

// (λt f . t)
var TrueFunction = Bind[test_named](t_, f_).In(t_)

// (λt f . f)
var FalseFunction = Bind[test_named](t_, f_).In(f_)

// (λa b . a true b)
var OrFunction = Bind[test_named](a_, b_).In(Apply[test_named](a_, TrueFunction, b_))

// (λc t e . c t e)
var IfFunction = Bind[test_named](c_, t_, e_).In(Apply[test_named](c_, t_, e_))

// (λa b . a b false)
var AndFunction = Bind[test_named](a_, b_).In(Apply[test_named](a_, b_, FalseFunction))

// (λa . a false true)
var NotFunction = Bind[test_named](a_).In(Apply[test_named](a_, FalseFunction, TrueFunction))

// (λf . (λx . f (x x)) (λx . f (x x)))
var Y = Bind[test_named](f_).In(Apply[test_named](
	Bind[test_named](x_).In(Apply[test_named](f_, Apply[test_named](x_, x_))),
	Bind[test_named](x_).In(Apply[test_named](f_, Apply[test_named](x_, x_)))))

func (f Function[T]) Find(v Variable[T]) bool {
	// update to account for binders
	v2 := Var(v.name)
	v2.depth = v.depth + f.BindDepth()
	return f.e.Find(v2)
}

// assumes this was called from EtaReduction[T]
func (f Function[T]) etaReduction_absorb(g Function[T]) Function[T] {
	bindDepth := f.BindDepth()
	res := g.UpdateVars(f.BindDepth(), -1).(Function[T])
	if bindDepth == 1 {
		return res
	}

	//add := bindDepth - 1
	vs := make([]Variable[T], (len(f.vars)-1)+len(g.vars))
	//ex := res.e.CleanVars() // cleans function[T] for re-binding
	//println(ex.StrictString())
	for i, v := range f.vars[:bindDepth-1] {
		vs[i] = Var(v.name)
	}
	for i, v := range g.vars {
		vs[i+bindDepth-1] = Var(v.name)
	}
	//vs = append(f.vars[:bindDepth-1], vs...)
	return mkfunc(vs...).rebinding(res.e)
	/*return Function[T]{
		vars: vs,
		e:    ex,
	}*/
}

// (\x -> (\f -> e) x) == (\f -> e)
func (f Function[T]) EtaReduction() Function[T] {
	// f = (\x -> e)
	if a, ok := f.e.(Application[T]); ok {
		lookFor := f.vars[f.BindDepth()-1]
		// e = (e1 e2) => f = (\x -> e1 e2)
		g, ok := a.left.(Function[T])
		if !ok {
			return f
		}
		// e = ((\y -> e3) e2) => f = (\x -> ((\y -> e3) e2))

		if a.right.StrictEquals(lookFor) { // TODO: ?? a.right.Equals(lookFor)
			if g.Find(lookFor) {
				return f // g contains instances of lookFor
			}
			// g contains no instances of lookFor, do eta reduction[T]
			return f.etaReduction_absorb(g)
		}
	}
	return f
}

// (λx.e)
func (f Function[T]) String() string {
	return "(" + binder_string + BindersOnly[T](f.vars).String() + to_string + f.e.String() + ")"
}

func (f Function[T]) StrictString() string {
	return "(" + binder_string + BindersOnly[T](f.vars).StrictString() + to_string + f.e.StrictString() + ")"
}

func (f Function[T]) Copy() Expression[T] {
	vs := make([]Variable[T], len(f.vars))
	for i, v := range f.vars {
		vs[i] = v.copy()
	}
	return Function[T]{
		vars: vs,
		e:    f.e.Copy(),
	}
}

func (f Function[T]) Rebind() Expression[T] {
	return BindersOnly[T](f.vars).Clean().In(f.e.Rebind())
}

type BindersOnly[T nameable.Nameable] []Variable[T]

func (bs BindersOnly[T]) Collect() []T {
	res := make([]T, 0, len(bs))
	for _, v := range bs {
		res = append(res, v.Collect()...)
	}
	return res
}

func (bs BindersOnly[T]) String() string {
	return str.Join(bs, str.String(" "))
}

func (bs BindersOnly[T]) StrictString() string {
	strs := make([]string, len(bs))
	for i, b := range bs {
		strs[i] = b.StrictString()
	}
	return strings.Join(strs, " ")
}

func (bs BindersOnly[T]) Clean() BindersOnly[T] {
	out := make(BindersOnly[T], len(bs))
	for i, b := range bs {
		out[i] = Var(b.name)
	}
	return out
}

func (bs BindersOnly[T]) rebinding(e Expression[T]) Function[T] {
	return bs.In(e.Rebind())
}

func mkfunc[T nameable.Nameable](bs ...Variable[T]) BindersOnly[T] {
	depth := len(bs)
	out := make(BindersOnly[T], len(bs))
	for i, b := range bs {
		out[i] = Var(b.name)
		out[i].depth = depth - i
	}
	return out
}

func Bind[T nameable.Nameable](binder Variable[T], more ...Variable[T]) BindersOnly[T] {
	if len(more) == 0 {
		return mkfunc(binder)
	}

	return mkfunc(append([]Variable[T]{binder}, more...)...)
}

func (bs BindersOnly[T]) Update(add int) BindersOnly[T] {
	out := make(BindersOnly[T], len(bs))
	for i, b := range bs {
		out[i] = Var(b.name)
		out[i].depth = b.depth + add
	}
	return out
}

func makeFunction[T nameable.Nameable](vars []Variable[T], e Expression[T]) Function[T] {
	return Function[T]{vars: vars, e: e}
}

func (f Function[T]) Bind(bs BindersOnly[T]) Expression[T] {
	return Function[T]{
		vars: f.vars,
		e:    f.e.Bind(bs.Update(len(f.vars))),
	}
}

func (bs BindersOnly[T]) In(e Expression[T]) Function[T] {
	return Function[T]{
		vars: bs,
		e:    e.Bind(bs),
	}
}

//var _freeApplyVar = Var("_")

func (f Function[T]) FreeApplyThrough(cxt *Context[T]) Expression[T] {
	var e Expression[T]
	g, ok := f, true
	for ok {
		e = g.Apply(cxt.Var("_"))
		g, ok = e.(Function[T])
	}
	return e
}

func functionEquals[T nameable.Nameable](cxt *Context[T], f, g Function[T]) bool {
	return f.FreeApplyThrough(cxt).Equals(cxt, g.FreeApplyThrough(cxt))
}

func (f Function[T]) Equals(cxt *Context[T], e Expression[T]) bool {
	f2, ok := e.ForceRequest().(Function[T])
	if !ok {
		return false
	}
	return functionEquals(cxt, f, f2)
}

func functionStrictEquals[T nameable.Nameable](f, g Function[T]) bool {
	return fun.AndZip(true, f.vars, g.vars, varEquals[T]) && f.e.StrictEquals(g.e)
}

func (f Function[T]) StrictEquals(e Expression[T]) bool {
	f2, ok := e.(Function[T])
	if !ok {
		return false
	}
	return functionStrictEquals(f, f2)
}

func (f Function[T]) Replace(v Variable[T], e Expression[T]) (Expression[T], bool) {
	v2 := v.UpdateVars(0, f.BindDepth()).(Variable[T])
	e2 := e.UpdateVars(0, f.BindDepth())
	res, again := f.e.Replace(v2, e2)
	return Function[T]{
		vars: f.vars,
		e:    res,
	}, again
}

// func (f Function[T]) Apply(e Expression[T]) {
//	lookFor := f.nBinders
//	f.nBinders := f.nBinders - 1
//	e.updateVars(f.nBinders)
//	res := f.expr.Replace(lookFor, e).freeDecrement(lookFor)
// }

func (f Function[T]) BindDepth() int {
	return len(f.vars)
}

// update all vars v > `gt` by `v = v + by`
func (f Function[T]) UpdateVars(gt int, by int) Expression[T] {
	return Function[T]{
		vars: f.vars,
		e:    f.e.UpdateVars(f.BindDepth(), by),
	}
}

// just returns f
func (f Function[T]) PrepareAsRHS() Expression[T] { return f }

// (λa b . a b) (λa b c . a b z)
// ( f=(λλ 2 1) ).Apply( e=(λλλ 3 2 4) )
// lookFor = 2
// e2 = (λλλ 3 2 6)
// v = 2
// res = { ((λλλ 3 2 6) 1) => call to apply in Application[T]
//	 ( (λλλ 3 2 6) ).Apply( e=(1) )
//	 lookFor = 3
//	 e2 = (4)
// 	 v = 3
//	 res = (4 2 6)
//   res = (3 2 5)
//	 return (λλ 3 2 5)
// } = (λλ 3 2 5)
// res = (λλ 3 2 5)

func (f Function[T]) Apply(e Expression[T]) Expression[T] {
	lookFor := f.BindDepth() // Variable[T] number being replaced
	e2 := e.
		PrepareAsRHS(). // makes sure free variables have a depth > 0
		UpdateVars(0, lookFor) // updates free variables so they have a depth > f.BindDepth()
	v := f.vars[0] // Variable[T] to replace (should have same number as `lookFor`)
	res, again := f.e.Replace(v, e2) // replace variables matching `v` with `e2`
	res = res.UpdateVars(lookFor, -1) // dec. free vars to account for loss of binder
	for again { // need to apply args again? 
		res, again = res.Again()
	}

	if lookFor > 1 {
		return Function[T]{ // function[T] has at least one binder left
			vars: f.vars[1:],
			e:    res,
		}
	}
	return res // no binders left
}

func (f Function[T]) DoApplication(e Expression[T]) Expression[T] {
	return f.Apply(e)
} 

func (f Function[T]) ForceRequest() Expression[T] { return f }

func (f Function[T]) Again() (Expression[T], bool) {
	/*res, again := f.e.Again()
	return Function[T]{
		vars: f.vars,
		e:    res,
	}, again*/
	return f, false
}

func (f Function[T]) AgainApply(e Expression[T]) (Expression[T], bool) {
	lookFor := f.BindDepth() // same value, bindDepth is just uint
	e2 := e.UpdateVars(0, lookFor)
	v := f.vars[0]
	res, again := f.e.Replace(v, e2)
	res = res.UpdateVars(lookFor, -1)

	if lookFor > 1 {
		return Function[T]{
			vars: f.vars[1:],
			e:    res,
		}, again
	}
	return res, again
}
