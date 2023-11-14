package expr

import "github.com/petersalex27/yew-packages/nameable"

type Def[T nameable.Nameable] struct {
	// Important: `name` is NOT a "variable"! variables are bound or free--they do not
	// name specific expressions values like `name` does
	name Const[T]
	// expression value named by `name`
	assignment Expression[T]
}

func (def Def[T]) Flatten() []Expression[T] {
	return []Expression[T]{def.name, def.assignment}
}

func (def Def[T]) Equals(cxt *Context[T], def2 Def[T]) bool {
	return def.name.Equals(cxt, def2.name) &&
		def.assignment.Equals(cxt, def2.assignment)
}

func (def Def[T]) StrictEquals(def2 Def[T]) bool {
	return def.name.StrictEquals(def2.name) &&
		def.assignment.StrictEquals(def2.assignment)
}

func (def Def[T]) String() string {
	return def.name.String() + " = " + def.assignment.String()
}

func (def Def[T]) StrictString() string {
	return def.name.StrictString() + " = " + def.assignment.StrictString()
}

func (def Def[T]) GetName() Const[T] { return def.name }

func (def Def[T]) GetAssignment() Expression[T] { return def.assignment }

// returns a definition
func Define[T nameable.Nameable](name Const[T], assignment Expression[T]) Def[T] {
	if assignment == nil {
		panic("nil value error: argument passed for parameter `assignemnt` cannot be nil")
	}
	return Def[T]{name, assignment}
}

type LetIn[T nameable.Nameable] struct {
	Def[T]
	// expression in which `name` names `assignment`
	contextualized Expression[T]
}

func (letin LetIn[T]) Flatten() []Expression[T] {
	return append(letin.Def.Flatten(), letin.contextualized.Flatten()...)
}

type NameContext[T nameable.Nameable] struct {
	// whether context is at head of NameContext or tail
	// false (i.e., at head) =>
	//	let <name> = <assignment> in <contextualized>
	// true (i.e., at tail) =>
	// 	<contextualized> where <name> = <assignment>
	tailedContext bool
	LetIn[T]
}

func (cxt NameContext[T]) Flatten() []Expression[T] {
	return cxt.LetIn.Flatten()
}

// Creates a non-tailed NameContext expression; in Haskell this would be a 
// "let-in" expression:
// 		let name = assignment in contextualized
//
// NOTE: panics if `assignment` is nil
func Let[T nameable.Nameable](name Const[T], assignment Expression[T], contextualized Expression[T]) NameContext[T] {
	if assignment == nil {
		panic("nil value error: argument passed for parameter `assignemnt` cannot be nil")
	}

	return NameContext[T]{false, LetIn[T]{Def[T]{name, assignment}, contextualized}}
}

// Creates a tailed NameContext expression; in Haskell this would be a "where"
// expression:
// 		contextualized where name = assignment
//
// NOTE: panics if `assignment` is nil
func Where[T nameable.Nameable](contextualized Expression[T], name Const[T], assignment Expression[T]) NameContext[T] {
	out := Let(name, assignment, contextualized)
	out.tailedContext = true
	return out
}

func (cxt NameContext[T]) BodyAbstract(v Variable[T], name Const[T]) Expression[T] {
	// avoid capturing let-bound names
	if cxt.Def.name.Name.GetName() == name.Name.GetName() {
		return cxt // `name` is shadowed by `x` in `let x = e0 in e1`
	}

	// name context does not bind `name`, so bind assignment and context
	return NameContext[T]{
		tailedContext: cxt.tailedContext,
		LetIn: LetIn[T]{
			Def: Def[T]{
				name: cxt.name,
				assignment: cxt.Def.assignment.BodyAbstract(v, name),
			},
			contextualized: cxt.contextualized.BodyAbstract(v, name),
		},
	}
}

// Setter method for struct member `contextualized`
//
// NOTE: panics if contextualized is nil
func (cxt NameContext[T]) SetContextualized(contextualized Expression[T]) NameContext[T] {
	if contextualized == nil {
		panic("nil value error: argument passed for parameter `contextualized` cannot be nil")
	}

	cxt.contextualized = contextualized
	return cxt
}

func (cxt NameContext[T]) GetContextualized() Expression[T] {
	return cxt.contextualized
}

func assembleNameContextString(def, contextualized string, tailedContext bool) string {
	if tailedContext {
		return contextualized + " where " + def
	}
	return "let " + def + " in " + contextualized
}

func (nameCxt NameContext[T]) String() string {
	def := nameCxt.Def.String()
	contextualized := nameCxt.contextualized.String()
	return assembleNameContextString(def, contextualized, nameCxt.tailedContext)
}

func (nameCxt NameContext[T]) Equals(cxt *Context[T], e Expression[T]) bool {
	nameCxt2, ok := e.(NameContext[T])
	if !ok {
		return false
	}
	return nameCxt.tailedContext == nameCxt2.tailedContext &&
		nameCxt.Def.Equals(cxt, nameCxt2.Def) &&
		nameCxt.contextualized.Equals(cxt, nameCxt2.contextualized)
}

func (nameCxt NameContext[T]) StrictString() string {
	def := nameCxt.Def.StrictString()
	contextualized := nameCxt.contextualized.StrictString()
	return assembleNameContextString(def, contextualized, nameCxt.tailedContext)
}

func (nameCxt NameContext[T]) StrictEquals(e Expression[T]) bool {
	nameCxt2, ok := e.(NameContext[T])
	if !ok {
		return false
	}
	return nameCxt.tailedContext == nameCxt2.tailedContext &&
		nameCxt.Def.StrictEquals(nameCxt2.Def) &&
		nameCxt.contextualized.StrictEquals(nameCxt2.contextualized)
}

func (nameCxt NameContext[T]) Replace(v Variable[T], e Expression[T]) (Expression[T], bool) {
	assignment, _ := nameCxt.assignment.Replace(v, e)
	contextualized, _ := nameCxt.contextualized.Replace(v, e)
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		LetIn: LetIn[T]{
			Def: Def[T]{
				name:           nameCxt.name,
				assignment:     assignment,
			},
			contextualized: contextualized,
		},
	}, false
}

func (nameCxt NameContext[T]) UpdateVars(gt int, by int) Expression[T] {
	assignment := nameCxt.assignment.UpdateVars(gt, by)
	contextualized := nameCxt.contextualized.UpdateVars(gt, by)
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		LetIn: LetIn[T]{
			Def: Def[T]{
				name:           nameCxt.name,
				assignment:     assignment,
			},
			contextualized: contextualized,
		},
	}
}

func (nameCxt NameContext[T]) Again() (Expression[T], bool) {
	assignment, again := nameCxt.assignment.Again()
	contextualized, again2 := nameCxt.contextualized.Again()
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		LetIn: LetIn[T]{
			Def: Def[T]{
				name:           nameCxt.name,
				assignment:     assignment,
			},
			contextualized: contextualized,
		},
	}, again || again2
}

func (nameCxt NameContext[T]) Bind(binders BindersOnly[T]) Expression[T] {
	assignment := nameCxt.assignment.Bind(binders)
	contextualized := nameCxt.contextualized.Bind(binders)
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		LetIn: LetIn[T]{
			Def: Def[T]{
				name:           nameCxt.name,
				assignment:     assignment,
			},
			contextualized: contextualized,
		},
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
		LetIn: LetIn[T]{
			Def: Def[T]{
				name:           nameCxt.name,
				assignment:     assignment,
			},
			contextualized: contextualized,
		},
	}
}

func (nameCxt NameContext[T]) Rebind() Expression[T] {
	assignment := nameCxt.assignment.Rebind()
	contextualized := nameCxt.contextualized.Rebind()
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		LetIn: LetIn[T]{
			Def: Def[T]{
				name:           nameCxt.name,
				assignment:     assignment,
			},
			contextualized: contextualized,
		},
	}
}

func (nameCxt NameContext[T]) Copy() Expression[T] {
	assignment := nameCxt.assignment.Copy()
	contextualized := nameCxt.contextualized.Copy()
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		LetIn: LetIn[T]{
			Def: Def[T]{
				name:           nameCxt.name,
				assignment:     assignment,
			},
			contextualized: contextualized,
		},
	}
}

func (nameCxt NameContext[T]) ForceRequest() Expression[T] {
	assignment := nameCxt.assignment.ForceRequest()
	contextualized := nameCxt.contextualized.ForceRequest()
	return NameContext[T]{
		tailedContext:  nameCxt.tailedContext,
		LetIn: LetIn[T]{
			Def: Def[T]{
				name:           nameCxt.name,
				assignment:     assignment,
			},
			contextualized: contextualized,
		},
	}
}

func (nameCxt NameContext[T]) ExtractVariables(gt int) []Variable[T] {
	cxtzdVars := nameCxt.contextualized.ExtractVariables(gt)
	assgnVars := nameCxt.assignment.ExtractVariables(gt)
	// the following statements ensure the order of the free vars matches the 
	// order they would appear in code
	if nameCxt.tailedContext {
		return append(cxtzdVars, assgnVars...)
	}
	return append(assgnVars, cxtzdVars...)
}

func (nameCxt NameContext[T]) Collect() []T {
	return append(nameCxt.assignment.Collect(), nameCxt.contextualized.Collect()...)
}
