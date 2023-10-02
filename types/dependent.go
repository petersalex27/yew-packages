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

type IndexerGenerator[T nameable.Nameable] func(...expr.Expression) Monotyped[T]

// Dependent Type[T]
type DependentType[T nameable.Nameable] struct {
	indexForm DependentTypeInstance[T]
	indexedBy []TypeJudgement[T,expr.Variable]
	indexConstruction []DependentTypeConstructor[T]
}

func (d DependentType[T]) String() string {
	return "map " + str.Join(d.indexedBy, str.String(" ")) + " . "
}

func (d DependentType[T]) KindInstantiation() DependentTypeInstance[T] {
	return d.indexForm
}

func (d DependentType[T]) FreeInstantiation(cxt *Context[T]) DependentTyped[T] {
	cs := make([]DependentTypeConstructor[T], len(d.indexConstruction))
	for i, c := range d.indexConstruction {
		cs[i] = c.FreeInstantiateKinds(d.indexedBy...)
	}
	return DependentType[T]{
		indexedBy: nil,
		indexConstruction: cs,
	}
}

// Allows the following (as long as B does not depend on A (in the first operand) 
// and A does not depend on B (in the second operand)):
// (map (x: A) (y: B) . W(y)) == (map (y: B) (x: A) . W(y)) == (map (y: B) . W(y))
func (d DependentType[T]) Equals(t Type[T]) bool {
	d2, ok := t.(DependentType[T])
	if !ok {
		return false
	}
	return d.indexForm.Equals(d2.indexForm)
	//return d.FreeInstantiation().Equals(d2.FreeInstantiation())
}

func (d DependentType[T]) Generalize(cxt *Context[T]) Polytype[T] {
	return Polytype[T]{
		typeBinders: cxt.MakeDummyVars(1),
		bound: d,
	}
}

func (d DependentType[T]) Collect() []T {
	res := make([]T, 0, 1 + len(d.indexConstruction) + len(d.indexedBy))
	res = append(res, d.indexForm.Collect()...)
	for _, v := range d.indexedBy {
		res = append(res, v.Collect()...)
	}
	for _, v := range d.indexConstruction {
		res = append(res, v.Collect()...)
	}
	return res
}

func (d DependentType[T]) ReplaceDependent(v Variable[T], m Monotyped[T]) DependentTyped[T] {
	out := DependentType[T]{
		indexedBy: d.indexedBy,
		indexConstruction: make([]DependentTypeConstructor[T], len(d.indexConstruction)),
	}

	for i, con := range d.indexConstruction {
		out.indexConstruction[i] = con.Replace(v, m)
	}
	return out
}
