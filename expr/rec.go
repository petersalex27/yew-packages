package expr

import (
	"strings"

	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/stringable"
)

type RecIn[T nameable.Nameable] struct {
	defs []Def[T]
	// expression in which `name` names `assignment`
	contextualized Expression[T]
}

func (rec RecIn[T]) Flatten() []Expression[T] {
	f := (Def[T]).Flatten
	fold := func(l, r []Expression[T]) []Expression[T] {
		return append(l, r...)
	}
	left := fun.FoldLeft([]Expression[T]{}, fun.FMap(rec.defs, f), fold)
	return append(left, rec.contextualized.Flatten()...)
}

func (rec RecIn[T]) BodyAbstract(v Variable[T], name Const[T]) Expression[T] {
	// first, validate
	for _, def := range rec.defs {
		if def.name.Name.GetName() == name.Name.GetName() {
			return rec // cannot bind, `name` is shadowed in `rec`
		}
	}

	// `name` is not shadowed directly by `rec`
	// bind defs
	defs := fun.FMap(
		rec.defs,
		func(def Def[T]) Def[T] {
			return Def[T]{
				name:       def.name,
				assignment: def.assignment.BodyAbstract(v, name),
			}
		},
	)

	// bind contextualized
	contextualized := rec.contextualized.BodyAbstract(v, name)

	return RecIn[T]{defs, contextualized}
}

// creates something similar to a let-in expression but that allows recursion
//
//	rec
//		name1 = assignment1 and
//		name2 = assignment2 and
//		..
//		namen = assignmentn in
//			contextualized
//
// NOTE: panics if len(defs) == 0
func Rec[T nameable.Nameable](defs ...Def[T]) func(contextualized Expression[T]) RecIn[T] {
	if len(defs) == 0 {
		panic("must have at least one name definition")
	}

	return func(contextualized Expression[T]) RecIn[T] {
		return RecIn[T]{defs, contextualized}
	}
}

// Setter method for struct member `contextualized`
//
// NOTE: panics if contextualized is nil
func (rec RecIn[T]) SetContextualized(contextualized Expression[T]) RecIn[T] {
	if contextualized == nil {
		panic("nil value error: argument passed for parameter `contextualized` cannot be nil")
	}

	rec.contextualized = contextualized
	return rec
}

// returns all top level defs
func (rec RecIn[T]) GetDefs() []Def[T] { return rec.defs }

// returns list of all names defined at top level
func (rec RecIn[T]) GetNames() []Const[T] {
	// defs[0] exist assuming `Rec` was used to create `rec`
	defs := rec.GetDefs()
	names := make([]Const[T], len(defs))
	for i, def := range defs {
		names[i] = def.GetName()
	}
	return names
}

// returns all things assigned to defs
func (rec RecIn[T]) GetAssignments() []Expression[T] {
	// defs[0] exist assuming `Rec` was used to create `rec`
	defs := rec.GetDefs()
	exprs := make([]Expression[T], len(defs))
	for i, def := range defs {
		exprs[i] = def.GetAssignment()
	}
	return exprs
}

func (rec RecIn[T]) GetContextualized() Expression[T] {
	return rec.contextualized
}

func (rec RecIn[T]) String() string {
	defs := stringable.Join(rec.defs, stringable.String(" and "))
	contextualized := rec.contextualized.String()
	return "rec " + defs + " in " + contextualized
}

func (rec RecIn[T]) Equals(cxt *Context[T], e Expression[T]) bool {
	rec2, ok := e.(RecIn[T])
	if !ok {
		return false
	}
	if !rec.contextualized.Equals(cxt, rec2.contextualized) {
		return false
	}

	if len(rec.defs) != len(rec2.defs) {
		return false
	}

	for i, def := range rec.defs {
		if !def.Equals(cxt, rec2.defs[i]) {
			return false
		}
	}
	return true
}

func (rec RecIn[T]) StrictString() string {
	defsStrs := make([]string, len(rec.defs))
	for i, def := range rec.defs {
		defsStrs[i] = def.StrictString()
	}
	defs := strings.Join(defsStrs, " and ")
	contextualized := rec.contextualized.StrictString()
	return "rec " + defs + " in " + contextualized
}

func (rec RecIn[T]) StrictEquals(e Expression[T]) bool {
	rec2, ok := e.(RecIn[T])
	if !ok {
		return false
	}
	if !rec.contextualized.StrictEquals(rec2.contextualized) {
		return false
	}

	if len(rec.defs) != len(rec2.defs) {
		return false
	}

	for i, def := range rec.defs {
		if !def.StrictEquals(rec2.defs[i]) {
			return false
		}
	}
	return true
}

// in `rec` replace variable `v` in with expression `e`
func (rec RecIn[T]) Replace(v Variable[T], e Expression[T]) (Expression[T], bool) {
	defs := fun.FMap(rec.defs, func(def Def[T]) Def[T] {
		assign, _ := def.assignment.Replace(v, e)
		return Def[T]{def.name, assign}
	})

	contextualized, _ := rec.contextualized.Replace(v, e)

	return Rec[T](defs...)(contextualized), false
}

func (rec RecIn[T]) UpdateVars(gt int, by int) Expression[T] {
	defs := fun.FMap(rec.defs, func(def Def[T]) Def[T] {
		assign := def.assignment.UpdateVars(gt, by)
		return Def[T]{def.name, assign}
	})

	contextualized := rec.contextualized.UpdateVars(gt, by)

	return Rec[T](defs...)(contextualized)
}

func (rec RecIn[T]) Again() (Expression[T], bool) {
	var again bool = false
	defs := fun.FMap(rec.defs, func(def Def[T]) Def[T] {
		assign, tmp := def.assignment.Again()
		again = again || tmp
		return Def[T]{def.name, assign}
	})

	contextualized, again2 := rec.contextualized.Again()
	return Rec[T](defs...)(contextualized), again || again2
}

func (rec RecIn[T]) Bind(binders BindersOnly[T]) Expression[T] {
	defs := fun.FMap(rec.defs, func(def Def[T]) Def[T] {
		assign := def.assignment.Bind(binders)
		return Def[T]{def.name, assign}
	})

	contextualized := rec.contextualized.Bind(binders)

	return Rec[T](defs...)(contextualized)
}

func (rec RecIn[T]) Find(v Variable[T]) bool {
	for _, def := range rec.defs {
		if def.assignment.Find(v) {
			return true
		}
	}

	return rec.contextualized.Find(v)
}

func (rec RecIn[T]) PrepareAsRHS() Expression[T] {
	defs := fun.FMap(rec.defs, func(def Def[T]) Def[T] {
		assign := def.assignment.PrepareAsRHS()
		return Def[T]{def.name, assign}
	})

	contextualized := rec.contextualized.PrepareAsRHS()

	return Rec[T](defs...)(contextualized)
}

func (rec RecIn[T]) Rebind() Expression[T] {
	defs := fun.FMap(rec.defs, func(def Def[T]) Def[T] {
		assign := def.assignment.Rebind()
		return Def[T]{def.name, assign}
	})

	contextualized := rec.contextualized.Rebind()

	return Rec[T](defs...)(contextualized)
}

func (rec RecIn[T]) Copy() Expression[T] {
	defs := fun.FMap(rec.defs, func(def Def[T]) Def[T] {
		assign := def.assignment.Copy()
		return Def[T]{def.name, assign}
	})

	contextualized := rec.contextualized.Copy()

	return Rec[T](defs...)(contextualized)
}

func (rec RecIn[T]) ForceRequest() Expression[T] {
	defs := fun.FMap(rec.defs, func(def Def[T]) Def[T] {
		assign := def.assignment.ForceRequest()
		return Def[T]{def.name, assign}
	})

	contextualized := rec.contextualized.ForceRequest()

	return Rec[T](defs...)(contextualized)
}

func (rec RecIn[T]) ExtractVariables(gt int) []Variable[T] {
	// extracts free variables from defs--creates 2d slice
	freeVars2d := fun.FMap(rec.defs, func(def Def[T]) []Variable[T] {
		return def.assignment.ExtractVariables(gt)
	})
	// flattens the extracted, 2d slice into a 1d slice
	freeVars := fun.FoldLeft([]Variable[T]{}, freeVars2d, func(xs, ys []Variable[T]) []Variable[T] {
		return append(xs, ys...)
	})

	contextualizedFreeVars := rec.contextualized.ExtractVariables(gt)

	freeVars = append(freeVars, contextualizedFreeVars...)
	return freeVars
}

func (rec RecIn[T]) Collect() []T {
	// collect from each def
	defsCollection2d := fun.FMap(rec.defs, func(def Def[T]) []T {
		assign := def.assignment.Collect()
		return assign
	})
	// flattens collected slice
	defsCollection := fun.FoldLeft([]T{}, defsCollection2d, func(xs, ys []T) []T {
		return append(xs, ys...)
	})

	// collect from contexualized
	contextualizedCollection := rec.contextualized.Collect()

	collection := append(defsCollection, contextualizedCollection...)
	return collection
}
