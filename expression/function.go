package expr

import (
	"github.com/petersalex27/yew-packages/fun"
	str "github.com/petersalex27/yew-packages/stringable"
	"strings"
)

// λx.e
type Function struct {
	vars []Variable
	e    Expression
}

// (λx . x)
var IdFunction = Bind(Var("x")).In(Var("x"))

// (λx . y)
var ConstFunction = Bind(Var("x")).In(Var("y"))

// (λt f . t)
var TrueFunction = Bind(Var("t"), Var("f")).In(Var("t"))

// (λt f . f)
var FalseFunction = Bind(Var("t"), Var("f")).In(Var("f"))

// (λa b . a true b)
var OrFunction = Bind(Var("a"), Var("b")).In(Apply(Var("a"), TrueFunction, Var("b")))

// (λc t e . c t e)
var IfFunction = Bind(Var("c"), Var("t"), Var("e")).In(Apply(Var("c"), Var("t"), Var("e")))

// (λa b . a b false)
var AndFunction = Bind(Var("a"), Var("b")).In(Apply(Var("a"), Var("b"), FalseFunction))

// (λa . a false true)
var NotFunction = Bind(Var("a")).In(Apply(Var("a"), FalseFunction, TrueFunction))

// (λf . (λx . f (x x)) (λx . f (x x)))
var Y = Bind(Var("f")).In(Apply(
	Bind(Var("x")).In(Apply(Var("f"), Apply(Var("x"), Var("x")))),
	Bind(Var("x")).In(Apply(Var("f"), Apply(Var("x"), Var("x"))))))

func (f Function) Find(v Variable) bool {
	// update to account for binders
	v2 := Var(v.name)
	v2.depth = v.depth + f.BindDepth()
	return f.e.Find(v2)
}

// assumes this was called from EtaReduction
func (f Function) etaReduction_absorb(g Function) Function {
	bindDepth := f.BindDepth()
	res := g.UpdateVars(f.BindDepth(), -1).(Function)
	if bindDepth == 1 {
		return res
	}

	//add := bindDepth - 1
	vs := make([]Variable, (len(f.vars)-1)+len(g.vars))
	//ex := res.e.CleanVars() // cleans function for re-binding
	//println(ex.StrictString())
	for i, v := range f.vars[:bindDepth-1] {
		vs[i] = Var(v.name)
	}
	for i, v := range g.vars {
		vs[i+bindDepth-1] = Var(v.name)
	}
	//vs = append(f.vars[:bindDepth-1], vs...)
	return mkfunc(vs...).rebinding(res.e)
	/*return Function{
		vars: vs,
		e:    ex,
	}*/
}

// (\x -> (\f -> e) x) == (\f -> e)
func (f Function) EtaReduction() Function {
	// f = (\x -> e)
	if a, ok := f.e.(Application); ok {
		lookFor := f.vars[f.BindDepth()-1]
		// e = (e1 e2) => f = (\x -> e1 e2)
		g, ok := a.left.(Function)
		if !ok {
			return f
		}
		// e = ((\y -> e3) e2) => f = (\x -> ((\y -> e3) e2))

		if a.right.Equals(lookFor) {
			if g.Find(lookFor) {
				return f // g contains instances of lookFor
			}
			// g contains no instances of lookFor, do eta reduction
			return f.etaReduction_absorb(g)
		}
	}
	return f
}

// (λx.e)
func (f Function) String() string {
	return "(" + binder_string + BindersOnly(f.vars).String() + to_string + f.e.String() + ")"
}

func (f Function) StrictString() string {
	return "(" + binder_string + BindersOnly(f.vars).StrictString() + to_string + f.e.StrictString() + ")"
}

func (f Function) Copy() Expression {
	vs := make([]Variable, len(f.vars))
	for i, v := range f.vars {
		vs[i] = v.copy()
	}
	return Function{
		vars: vs,
		e:    f.e.Copy(),
	}
}

func (f Function) Rebind() Expression {
	return BindersOnly(f.vars).Clean().In(f.e.Rebind())
}

type BindersOnly []Variable

func (bs BindersOnly) String() string {
	return str.Join(bs, str.String(" "))
}

func (bs BindersOnly) StrictString() string {
	strs := make([]string, len(bs))
	for i, b := range bs {
		strs[i] = b.StrictString()
	}
	return strings.Join(strs, " ")
}

func (bs BindersOnly) Clean() BindersOnly {
	out := make(BindersOnly, len(bs))
	for i, b := range bs {
		out[i] = Var(b.name)
	}
	return out
}

func (bs BindersOnly) rebinding(e Expression) Function {
	return bs.In(e.Rebind())
}

func mkfunc(bs ...Variable) BindersOnly {
	depth := len(bs)
	out := make(BindersOnly, len(bs))
	for i, b := range bs {
		out[i] = Var(b.name)
		out[i].depth = depth - i
	}
	return out
}

func Bind(binder Variable, more ...Variable) BindersOnly {
	if len(more) == 0 {
		return mkfunc(binder)
	}

	return mkfunc(append([]Variable{binder}, more...)...)
}

func (bs BindersOnly) Update(add int) BindersOnly {
	out := make(BindersOnly, len(bs))
	for i, b := range bs {
		out[i] = Var(b.name)
		out[i].depth = b.depth + add
	}
	return out
}

func makeFunction(vars []Variable, e Expression) Function {
	return Function{vars: vars, e: e}
}

func (f Function) Bind(bs BindersOnly) Expression {
	return Function{
		vars: f.vars,
		e:    f.e.Bind(bs.Update(len(f.vars))),
	}
}

func (bs BindersOnly) In(e Expression) Function {
	return Function{
		vars: bs,
		e:    e.Bind(bs),
	}
}

var _freeApplyVar = Var("_")

func (f Function) FreeApplyThrough() Expression {
	var e Expression
	g, ok := f, true
	for ok {
		e = g.Apply(_freeApplyVar)
		g, ok = e.(Function)
	}
	return e
}

func functionEquals(f, g Function) bool {
	return f.FreeApplyThrough().Equals(g.FreeApplyThrough())
}

func (f Function) Equals(e Expression) bool {
	f2, ok := e.ForceRequest().(Function)
	if !ok {
		return false
	}
	return functionEquals(f, f2)
}

func functionStrictEquals(f, g Function) bool {
	return fun.AndZip(true, f.vars, g.vars, varEquals) && f.e.StrictEquals(g.e)
}

func (f Function) StrictEquals(e Expression) bool {
	f2, ok := e.(Function)
	if !ok {
		return false
	}
	return functionStrictEquals(f, f2)
}

func (f Function) Replace(v Variable, e Expression) (Expression, bool) {
	v2 := v.UpdateVars(0, f.BindDepth()).(Variable)
	e2 := e.UpdateVars(0, f.BindDepth())
	res, again := f.e.Replace(v2, e2)
	return Function{
		vars: f.vars,
		e:    res,
	}, again
}

// func (f Function) Apply(e Expression) {
//	lookFor := f.nBinders
//	f.nBinders := f.nBinders - 1
//	e.updateVars(f.nBinders)
//	res := f.expr.Replace(lookFor, e).freeDecrement(lookFor)
// }

func (f Function) BindDepth() int {
	return len(f.vars)
}

// update all vars v > `gt` by `v = v + by`
func (f Function) UpdateVars(gt int, by int) Expression {
	return Function{
		vars: f.vars,
		e:    f.e.UpdateVars(f.BindDepth(), by),
	}
}

// just returns f
func (f Function) PrepareAsRHS() Expression { return f }

// (λa b . a b) (λa b c . a b z)
// ( f=(λλ 2 1) ).Apply( e=(λλλ 3 2 4) )
// lookFor = 2
// e2 = (λλλ 3 2 6)
// v = 2
// res = { ((λλλ 3 2 6) 1) => call to apply in Application
//	 ( (λλλ 3 2 6) ).Apply( e=(1) )
//	 lookFor = 3
//	 e2 = (4)
// 	 v = 3
//	 res = (4 2 6)
//   res = (3 2 5)
//	 return (λλ 3 2 5)
// } = (λλ 3 2 5)
// res = (λλ 3 2 5)

func (f Function) Apply(e Expression) Expression {
	lookFor := f.BindDepth() // variable number being replaced
	e2 := e.
		PrepareAsRHS(). // makes sure free variables have a depth > 0
		UpdateVars(0, lookFor) // updates free variables so they have a depth > f.BindDepth()
	v := f.vars[0] // variable to replace (should have same number as `lookFor`)
	res, again := f.e.Replace(v, e2) // replace variables matching `v` with `e2`
	res = res.UpdateVars(lookFor, -1) // dec. free vars to account for loss of binder
	for again { // need to apply args again? 
		res, again = res.Again()
	}

	if lookFor > 1 {
		return Function{ // function has at least one binder left
			vars: f.vars[1:],
			e:    res,
		}
	}
	return res // no binders left
}

func (f Function) DoApplication(e Expression) Expression {
	return f.Apply(e)
} 

func (f Function) ForceRequest() Expression { return f }

func (f Function) Again() (Expression, bool) {
	/*res, again := f.e.Again()
	return Function{
		vars: f.vars,
		e:    res,
	}, again*/
	return f, false
}

func (f Function) AgainApply(e Expression) (Expression, bool) {
	lookFor := f.BindDepth() // same value, bindDepth is just uint
	e2 := e.UpdateVars(0, lookFor)
	v := f.vars[0]
	res, again := f.e.Replace(v, e2)
	res = res.UpdateVars(lookFor, -1)

	if lookFor > 1 {
		return Function{
			vars: f.vars[1:],
			e:    res,
		}, again
	}
	return res, again
}
