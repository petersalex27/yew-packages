package expr

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/nameable"
)

type InstructionHead[T nameable.Nameable] struct {
	name   string
	nArgs  int
	action InstructionAction[T]
}

type InstructionAction[T nameable.Nameable] func(instr InstructionArgs[T]) Expression[T]

type InstructionArgs[T nameable.Nameable] struct {
	args []Expression[T]
}

type Instruction[T nameable.Nameable] struct {
	InstructionHead[T]
	InstructionArgs[T]
}

func (instr Instruction[T]) Flatten() []Expression[T] {
	f := (Expression[T]).Flatten
	fold := func(l, r []Expression[T]) []Expression[T] {
		return append(l, r...)
	}
	return fun.FoldLeft([]Expression[T]{}, fun.FMap(instr.args, f), fold)
}

func (instr Instruction[T]) BodyAbstract(v Variable[T], name Const[T]) Expression[T] {
	args := fun.FMap(
		instr.args,
		func(e Expression[T]) Expression[T] {
			return e.BodyAbstract(v, name)
		},
	)
	return Instruction[T]{
		InstructionHead: instr.InstructionHead,
		InstructionArgs: InstructionArgs[T]{args},
	}
}

func (instr Instruction[T]) ExtractVariables(gt int) []Variable[T] {
	vars := []Variable[T]{}
	for _, arg := range instr.args {
		vars = append(vars, arg.ExtractVariables(gt)...)
	}
	return vars
}

func (instr Instruction[T]) Collect() []T {
	if len(instr.args) == 0 {
		return []T{}
	}
	res := instr.args[0].Collect()
	for i := 1; i < len(instr.args); i++ {
		res = append(res, instr.args[i].Collect()...)
	}
	return res
}

func (instr Instruction[T]) IsCallReady() bool {
	return len(instr.args) == int(instr.nArgs) && instr.action != nil
}

func (instr Instruction[T]) TryCall(cxt *Context[T], catchPanic func(any)) (e Expression[T], success bool) {
	defer func() {
		if err := recover(); err != nil {
			success = false
			if catchPanic == nil {
				e = Const[T]{cxt.makeName(fmt.Sprint(err))}
			} else {
				e = nil
				catchPanic(err)
			}
		} else {
			success = true
		}
	}()
	e = instr.call()
	return
}

func (instr Instruction[T]) call() Expression[T] {
	return instr.action(instr.InstructionArgs)
}

func (instr InstructionArgs[T]) GetArgAtIndex(index int) Expression[T] {
	if len(instr.args) <= index || index < 0 {
		panic("tried to get a non-existent argument\n")
	} else {
		instr.args[index] = instr.args[index].ForceRequest()
		return instr.args[index]
	}
}

func (instr InstructionArgs[T]) GetArgAtPosition(position int) Expression[T] {
	if position < 1 {
		panic("tried to get an argument at a position[T] < 1; positions start at 1\n")
	} else {
		return instr.GetArgAtIndex(position - 1)
	}
}

func (instr Instruction[T]) Again() (Expression[T], bool) {
	if instr.IsCallReady() {
		return instr.call(), false
	}
	return instr, false
}

func (instr Instruction[T]) AgainApply(e Expression[T]) (res Expression[T], again bool) {
	res = instr.DoApplication(e)
	again = false
	return
}

func (instr Instruction[T]) Bind(bs BindersOnly[T]) Expression[T] {
	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i] = arg.Bind(bs)
	}
	head := DefineInstruction[T](instr.name, instr.nArgs, instr.action)
	return Instruction[T]{
		InstructionHead: head,
		InstructionArgs: InstructionArgs[T]{args: newArgs},
	}
}

func (ih InstructionHead[T]) Copy() InstructionHead[T] {
	return DefineInstruction[T](ih.name, ih.nArgs, ih.action)
}

func (instr Instruction[T]) Copy() Expression[T] {
	return Instruction[T]{
		InstructionHead: instr.InstructionHead.Copy(),
		InstructionArgs: instr.InstructionArgs.Copy(),
	}
}

func (instr Instruction[T]) Equals(cxt *Context[T], e Expression[T]) bool {
	instr2, ok := e.ForceRequest().(Instruction[T])
	if !ok {
		return false
	}
	return instructionHeadEquals(instr.InstructionHead, instr2.InstructionHead)
}

func (instr Instruction[T]) Find(v Variable[T]) bool {
	for _, arg := range instr.args {
		if arg.Find(v) {
			return true
		}
	}
	return false
}

func (instr Instruction[T]) PrepareAsRHS() Expression[T] {
	if len(instr.args) == 0 {
		return instr
	}
	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i] = arg.PrepareAsRHS()
	}
	head := DefineInstruction[T](instr.name, instr.nArgs, instr.action)
	return Instruction[T]{
		InstructionHead: head,
		InstructionArgs: InstructionArgs[T]{args: newArgs},
	}
}

func (instr Instruction[T]) Rebind() Expression[T] {
	if len(instr.args) == 0 {
		return instr
	}

	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i] = arg.Rebind()
	}
	head := DefineInstruction[T](instr.name, instr.nArgs, instr.action)
	return Instruction[T]{
		InstructionHead: head,
		InstructionArgs: InstructionArgs[T]{args: newArgs},
	}
}

func (instr Instruction[T]) Replace(v Variable[T], e Expression[T]) (Expression[T], bool) {
	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i], _ = arg.Replace(v, e)
	}
	head := DefineInstruction[T](instr.name, instr.nArgs, instr.action)
	return Instruction[T]{
		InstructionHead: head,
		InstructionArgs: InstructionArgs[T]{args: newArgs},
	}, false
}

func (instr Instruction[T]) ForceRequest() Expression[T] {
	if instr.IsCallReady() {
		return instr.call()
	}
	return instr
}

func (instr Instruction[T]) DoApplication(e Expression[T]) Expression[T] {
	if instr.IsCallReady() {
		return Apply(instr.call(), e).ForceRequest()
	}

	instr.args = append(instr.args, e)
	return instr.ForceRequest()
}

func instructionHeadEquals[T nameable.Nameable](ih1, ih2 InstructionHead[T]) bool {
	return ih1.nArgs == ih2.nArgs && ih1.name == ih2.name &&
		((ih1.action == nil && ih2.action == nil) || (ih1.action != nil && ih2.action != nil))
}

func instructionArgsEquals[T nameable.Nameable](ia1, ia2 InstructionArgs[T]) bool {
	ok := len(ia1.args) == len(ia2.args) && cap(ia1.args) == cap(ia2.args)
	if !ok {
		return false
	}

	for i := range ia1.args {
		if !ia1.args[i].StrictEquals(ia2.args[i]) {
			return false
		}
	}
	return true
}

func (instr Instruction[T]) StrictEquals(e Expression[T]) bool {
	instr2, ok := e.(Instruction[T])
	if !ok {
		return false
	}
	return instructionHeadEquals(instr.InstructionHead, instr2.InstructionHead) &&
		instructionArgsEquals(instr.InstructionArgs, instr2.InstructionArgs)
}

func (instr Instruction[T]) StrictString() string {
	return instr.String()
}

func (instr Instruction[T]) String() string {
	tmp := make([]string, cap(instr.args), instr.nArgs)
	for i, e := range instr.args {
		tmp[i] = e.String()
	}

	for i := len(instr.args); i < instr.nArgs; i++ {
		tmp[i] = "_"
	}

	return "instruction[" + instr.name + " " + strings.Join(tmp, " ") + "]"
}

func (ia InstructionArgs[T]) makeBlankCopy() []Expression[T] {
	ln, c := len(ia.args), cap(ia.args)
	return make([]Expression[T], ln, c)
}

func (ia InstructionArgs[T]) Copy() InstructionArgs[T] {
	out := ia.makeBlankCopy()
	for i, arg := range ia.args {
		out[i] = arg.Copy()
	}
	return InstructionArgs[T]{args: out}
}

func (instr Instruction[T]) UpdateVars(gt int, by int) Expression[T] {
	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i] = arg.UpdateVars(gt, by)
	}
	return Instruction[T]{
		InstructionHead: instr.InstructionHead,
		InstructionArgs: InstructionArgs[T]{args: newArgs},
	}
}

var EqualityInstruction = DefineInstruction[test_named]("testEquality", 2, func(instr InstructionArgs[test_named]) Expression[test_named] {
	e1, e2 := instr.GetArgAtIndex(0), instr.GetArgAtIndex(1)
	if e1.StrictEquals(e2) {
		return TrueFunction
	}
	return FalseFunction
})

func (ih InstructionHead[T]) MakeInstance() Instruction[T] {
	return Instruction[T]{
		InstructionHead: ih,
		InstructionArgs: InstructionArgs[T]{
			args: make([]Expression[T], 0, ih.nArgs),
		},
	}
}

func DefineInstruction[T nameable.Nameable](name string, nArgs int, action InstructionAction[T]) InstructionHead[T] {
	if nArgs < 0 {
		panic("instruction[T] cannot accept have than 0 arguments.\n")
	}
	return InstructionHead[T]{name: name, nArgs: nArgs, action: action}
}
