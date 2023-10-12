package expr

import "github.com/petersalex27/yew-packages/nameable"

type NameContext[T nameable.Nameable] struct {
	// whether context is at head of NameContext or tail
	// false (i.e., at head) =>
	//	let <name> = <assignment> in <contextualized>
	// true (i.e., at tail) =>
	// 	<contextualized> where <name> = <assignment>
	tailedContext bool
	// Important: `name` is NOT a "variable"! variables are bound or free--they do not
	// name specific expressions values like `name` does
	name Const[T]
	// expression value named by `name`
	assignment Expression[T]
	// expression in which `name` names `assignment`
	contextualized Expression[T]
}

func assembleNameContextString(name, assignment, contextualized string, tailedContext bool) string {
	if tailedContext {
		return contextualized + " where " + name + " = " + assignment
	}
	return "let " + name + " = " + assignment + " in " + contextualized
}

func (nameCxt NameContext[T]) String() string {
	name := nameCxt.name.String()
	assignment := nameCxt.assignment.String()
	contextualized := nameCxt.contextualized.String()
	return assembleNameContextString(name, assignment, contextualized, nameCxt.tailedContext)
}

func (nameCxt NameContext[T]) Equals(cxt *Context[T], e Expression[T]) bool {
	nameCxt2, ok := e.(NameContext[T])
	if !ok {
		return false
	}
	return nameCxt.tailedContext == nameCxt2.tailedContext &&
		nameCxt.name.Equals(cxt, nameCxt2.name) &&
		nameCxt.assignment.Equals(cxt, nameCxt2.assignment) &&
		nameCxt.contextualized.Equals(cxt, nameCxt2.contextualized)
}

func (nameCxt NameContext[T]) StrictString() string {
	name := nameCxt.name.StrictString()
	assignment := nameCxt.assignment.StrictString()
	contextualized := nameCxt.contextualized.StrictString()
	return assembleNameContextString(name, assignment, contextualized, nameCxt.tailedContext)
}

func (nameCxt NameContext[T]) StrictEquals(e Expression[T]) bool {
	nameCxt2, ok := e.(NameContext[T])
	if !ok {
		return false
	}
	return nameCxt.tailedContext == nameCxt2.tailedContext &&
		nameCxt.name.StrictEquals(nameCxt2.name) &&
		nameCxt.assignment.StrictEquals(nameCxt2.assignment) &&
		nameCxt.contextualized.StrictEquals(nameCxt2.contextualized)
}

func (nameCxt NameContext[T]) Replace(v Variable[T], e Expression[T]) (Expression[T], bool) {
	assignment, _ := nameCxt.assignment.Replace(v, e)
	contextualized, _ := nameCxt.contextualized.Replace(v, e)
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		name:           nameCxt.name,
		assignment:     assignment,
		contextualized: contextualized,
	}, false
}

func (nameCxt NameContext[T]) UpdateVars(gt int, by int) Expression[T] {
	assignment := nameCxt.assignment.UpdateVars(gt, by)
	contextualized := nameCxt.contextualized.UpdateVars(gt, by)
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		name:           nameCxt.name,
		assignment:     assignment,
		contextualized: contextualized,
	}
}

func (nameCxt NameContext[T]) Again() (Expression[T], bool) {
	assignment, again := nameCxt.assignment.Again()
	contextualized, again2 := nameCxt.contextualized.Again()
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		name:           nameCxt.name,
		assignment:     assignment,
		contextualized: contextualized,
	}, again || again2
}

func (nameCxt NameContext[T]) Bind(binders BindersOnly[T]) Expression[T] {
	assignment := nameCxt.assignment.Bind(binders)
	contextualized := nameCxt.contextualized.Bind(binders)
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		name:           nameCxt.name,
		assignment:     assignment,
		contextualized: contextualized,
	}
}

func (nameCxt NameContext[T]) Find(v Variable[T]) bool {
	return nameCxt.assignment.Find(v) || nameCxt.contextualized.Find(v)
}

func (nameCxt NameContext[T]) PrepareAsRHS() Expression[T] {
	assignment := nameCxt.assignment.PrepareAsRHS()
	contextualized := nameCxt.contextualized.PrepareAsRHS()
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		name:           nameCxt.name,
		assignment:     assignment,
		contextualized: contextualized,
	}
}

func (nameCxt NameContext[T]) Rebind() Expression[T] {
	assignment := nameCxt.assignment.Rebind()
	contextualized := nameCxt.contextualized.Rebind()
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		name:           nameCxt.name,
		assignment:     assignment,
		contextualized: contextualized,
	}
}

func (nameCxt NameContext[T]) Copy() Expression[T] {
	assignment := nameCxt.assignment.Copy()
	contextualized := nameCxt.contextualized.Copy()
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		name:           nameCxt.name,
		assignment:     assignment,
		contextualized: contextualized,
	}
}

func (nameCxt NameContext[T]) ForceRequest() Expression[T] {
	assignment := nameCxt.assignment.ForceRequest()
	contextualized := nameCxt.contextualized.ForceRequest()
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		name:           nameCxt.name,
		assignment:     assignment,
		contextualized: contextualized,
	}
}

func (nameCxt NameContext[T]) ExtractFreeVariables(dummyVar Variable[T]) []Variable[T] {
	return append(nameCxt.ExtractFreeVariables(dummyVar), nameCxt.assignment.ExtractFreeVariables(dummyVar)...)
}

func (nameCxt NameContext[T]) Collect() []T {
	return append(nameCxt.assignment.Collect(), nameCxt.contextualized.Collect()...)
}
