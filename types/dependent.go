package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	str "github.com/petersalex27/yew-packages/stringable"
)

type DependentTyped[T nameable.Nameable] interface {
	Type[T]
	FreeInstantiation(cxt *Context[T]) DependentTyped[T]
	ReplaceDependent(v Variable[T], with Monotyped[T]) DependentTyped[T]
}

/*
(Array a; n: Uint) =

	[]: Array a; 0
	| Cons a (Array a; n): (Array a; Succ n)

forall a . (

	typefamily FAM = { Uint, () } union { (Arr_0 a), (Arr_1 a), (Arr_2 a), .., (Arr_n a), .. }
	(Array a): Uint -> FAM
	(Array a)[0] = (Arr_0 a)
	(Array a)[1] = (Arr_1 a)
	...
	(Array a)[n] = (Arr_n a)
	...

)
forall a (n: Uint) . ((Array a)[n] = (Arr_n a))

-- forall a . Array a; 0
-- ^^^^^^^^^^^^^^^^^^^^^ this cannot be derived!!

-- Array Int; 0
-- ^^^^^^^^^^^^ this can be derived

-- forall a (n: Uint) . ((Array a)[0] = (Arr_0 a) & (Array n)[n] = (Arr_n a))
-- ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ this is derivable
-- forall a (n: N) . ([a; 0] = Ar0 a) & ([a; n] = Arn a)
-- Proof, 𝚪 = {0: N}:
-- 1. forall a (n: N) . ([a; n] = (Arn a))						[Premise]
-- 2. 	| newvar v												[Free Var Intro]
-- 3. 	| 	| forall (n: N) . ([v;n] = (Arn v))					1,2 [Instant.]
-- 4. 	| 	| 	| newvalvar k: N								[Free Val Intro]
-- 5. 	| 	| 	| 	| [v;k]=(Ark v)  							3,4 [Selection]
-- 6. 	| 	| 	| 	| 0: N										[Var]
-- 7. 	| 	| 	| 	| [v;0]=(Ar0 v)								3,6 [Selection]
-- 8.	|	|	| 	| [v;0]=(Ar0 v) & [v;k]=(Ark v)				5,7 [Construction]
-- 9. 	|	| 	| forall(n: N).([v;0]=(Ar0 v)&[v;k]=(Ark v))	5-8	[Generalization]
--10.	|	| forall(n: N).([v;0]=(Ar0 v)&[v;n]=(Arn v))		4,9 [Free Val Elim]
--11.	| forall a (n: N).([a;0]=(Ar0 a)&[a;n]=(Arn a))			3-10[Generalization]
--12. forall a (n: N).([a;0]=(Ar0 a)&[a;n]=(Arn a))				2,11[Free Var Elim]

Array = forall a (n: Uint) . Array a; n
forall a . (Array a; 0 == [])
Cons

Arr = forall a . map n: Uint . {
	[]: 0
	| Cons a (Arr a; n): n + 1
}(n)
*/

/*
Int = 0 | Succ Int | Pred Int
*/

type DependentTypeFunction[T nameable.Nameable] Application[T]

// Dependent Type: `(mapall (a: A) (b: B) ..) . (F a b ..)`
type DependentType[T nameable.Nameable] struct {
	mapall []TypeJudgement[T, expr.Variable[T]]
	//DependentTypeFunction[T]
	DependentTypeInstance[T]
}

func (d DependentType[T]) String() string {
	return "mapall " + str.Join(d.mapall, str.String(" ")) + " . " + d.DependentTypeInstance.String()
}

func kindInstantiation[T nameable.Nameable](d DependentType[T], defaultElem expr.Expression[T]) DependentTypeInstance[T] {
	index := make([]TypeJudgement[T, expr.Expression[T]], len(d.mapall))
	for i := range index {
		var elem expr.Expression[T]
		if defaultElem == nil {
			elem = d.mapall[i].expression
		} else {
			elem = defaultElem
		}
		index[i] = Judgement(elem, d.mapall[i].ty)
	}
	return DependentTypeInstance[T]{
		Application: d.DependentTypeInstance.Application,
		index:       index,
	}
}

func (d DependentType[T]) KindInstantiation() DependentTypeInstance[T] {
	return kindInstantiation(d, nil)
}

// ((mapall (a: A) (b: B) ..) . C) -> ((mapall (b: B) ..) . (C e))
func (cxt *Context[T]) InstantiateKind(d DependentType[T], e expr.Expression[T]) DependentTyped[T] {
	inst := d.DependentTypeInstance
	ty := d.mapall[0].ty 
	index := make([]TypeJudgement[T, expr.Expression[T]], len(inst.index)+1)
	copy(index, inst.index)
	index[len(inst.index)].expression = e
	index[len(inst.index)].ty = ty // type of expression should be type of variable being replaced
	
	out := DependentType[T]{
		mapall: d.mapall[1:],
		DependentTypeInstance: DependentTypeInstance[T]{
			Application: d.Application,
			index: index,
		},
	}

	if len(out.mapall) == 0 {
		return out.DependentTypeInstance
	}
	return out
}

func (d DependentType[T]) FreeInstantiation(cxt *Context[T]) DependentTyped[T] {
	v := expr.Var(cxt.makeName("_"))
	return kindInstantiation(d, expr.Expression[T](v))
}

// This test exact equality, not judgemental equality. For example,
//
//	(mapall (a: A) (b: B) . (C b)) != (mapall (b: B) . (C b))
//	despite the two being equiv. in some (probably useful) sense.
//
// Additionally, the following is not equiv. either:
//
//	(mapall (a: A) (b: B) . (C b)) != (mapall (b: B) (a: A) . (C b))
func (d DependentType[T]) Equals(t Type[T]) bool {
	d2, ok := t.(DependentType[T])
	if !ok {
		return false
	}

	if len(d.mapall) != len(d2.mapall) {
		return false
	}

	for i, judge := range d.mapall {
		if !judge.ty.Equals(d2.mapall[i].ty) {
			return false
		}
		// nil is okay here because variables don't require context object for equality
		if !judge.expression.Equals(nil, d2.mapall[i].expression) {
			return false
		}
	}

	return d.DependentTypeInstance.Equals(d2.DependentTypeInstance)
}

func (d DependentType[T]) Generalize(cxt *Context[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: cxt.MakeDummyVars(1),
		bound:       d,
	}
}

func (d DependentType[T]) Collect() []T {
	var res []T = []T{}
	if len(d.mapall) != 0 {
		res = d.mapall[0].Collect()
		for _, m := range d.mapall[1:] {
			res = append(res, m.Collect()...)
		}
	}
	res = append(res, d.DependentTypeInstance.Collect()...)
	return res
}

func (d DependentType[T]) ReplaceDependent(v Variable[T], m Monotyped[T]) DependentTyped[T] {
	mapall := make([]TypeJudgement[T, expr.Variable[T]], len(d.mapall))
	for i := range d.mapall {
		mapall[i].expression = d.mapall[i].expression
		if mono, ok := d.mapall[i].ty.(Monotyped[T]); ok {
			mapall[i].ty = mono.Replace(v, m)
		} else {
			mapall[i].ty = d.mapall[i].ty // TODO: if this branch happens, something is wrong (prob)
		}
	}

	inst, ok := d.DependentTypeInstance.ReplaceDependent(v, m).(DependentTypeInstance[T])
	if !ok {
		panic("bug: replacement should've resulted in application")
	}

	return DependentType[T]{
		mapall:                mapall,
		DependentTypeInstance: inst,
	}
}
