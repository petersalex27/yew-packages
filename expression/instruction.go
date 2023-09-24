package expr

import (
	"fmt"
	"strings"
)

type InstructionHead struct {
	name   string
	nArgs  int
	action InstructionAction
}

type InstructionAction func(instr InstructionArgs) Expression

type InstructionArgs struct {
	args []Expression
}

type Instruction struct {
	InstructionHead
	InstructionArgs
}

func (instr Instruction) IsCallReady() bool {
	return len(instr.args) == int(instr.nArgs) && instr.action != nil
}

func (instr Instruction) TryCall(catchPanic func(any)) (e Expression, success bool) {
	defer func() {
		if err := recover(); err != nil {
			success = false
			if catchPanic == nil {
				e = Const(fmt.Sprint(err))
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

func (instr Instruction) call() Expression {
	return instr.action(instr.InstructionArgs)
}

func (instr InstructionArgs) GetArgAtIndex(index int) Expression {
	if len(instr.args) <= index || index < 0 {
		panic("tried to get a non-existent argument\n")
	} else {
		instr.args[index] = instr.args[index].ForceRequest()
		return instr.args[index]
	}
}

func (instr InstructionArgs) GetArgAtPosition(position int) Expression {
	if position < 1 {
		panic("tried to get an argument at a position < 1; positions start at 1\n")
	} else {
		return instr.GetArgAtIndex(position-1)
	}
}

func (instr Instruction) Again() (Expression, bool) {
	if instr.IsCallReady() {
		return instr.call(), false
	}
	return instr, false
}

func (instr Instruction) AgainApply(e Expression) (res Expression, again bool) {
	res = instr.DoApplication(e)
	again = false
	return
}

func (instr Instruction) Bind(bs BindersOnly) Expression {
	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i] = arg.Bind(bs)
	}
	head := DefineInstruction(instr.name, instr.nArgs, instr.action)
	return Instruction{
		InstructionHead: head,
		InstructionArgs: InstructionArgs{ args: newArgs, },
	}
}

func (ih InstructionHead) Copy() InstructionHead {
	return DefineInstruction(ih.name, ih.nArgs, ih.action)
}

func (instr Instruction) Copy() Expression {
	return Instruction{
		InstructionHead: instr.InstructionHead.Copy(),
		InstructionArgs: instr.InstructionArgs.Copy(),
	}
}

func (instr Instruction) Equals(e Expression) bool {
	instr2, ok := e.ForceRequest().(Instruction)
	if !ok {
		return false
	}
	return instructionHeadEquals(instr.InstructionHead, instr2.InstructionHead)
}

func (instr Instruction) Find(v Variable) bool { 
	for _, arg := range instr.args {
		if arg.Find(v) {
			return true
		}
	}
	return false
}

func (instr Instruction) PrepareAsRHS() Expression {
	if len(instr.args) == 0 {
		return instr
	}
	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i] = arg.PrepareAsRHS()
	}
	head := DefineInstruction(instr.name, instr.nArgs, instr.action)
	return Instruction{
		InstructionHead: head,
		InstructionArgs: InstructionArgs{ args: newArgs, },
	}
}

func (instr Instruction) Rebind() Expression {
	if len(instr.args) == 0 {
		return instr
	}

	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i] = arg.Rebind()
	}
	head := DefineInstruction(instr.name, instr.nArgs, instr.action)
	return Instruction{
		InstructionHead: head,
		InstructionArgs: InstructionArgs{ args: newArgs, },
	}
}

func (instr Instruction) Replace(v Variable, e Expression) (Expression, bool) {
	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i], _ = arg.Replace(v, e)
	}
	head := DefineInstruction(instr.name, instr.nArgs, instr.action)
	return Instruction{
		InstructionHead: head,
		InstructionArgs: InstructionArgs{ args: newArgs, },
	}, false
}

func (instr Instruction) ForceRequest() Expression {
	if instr.IsCallReady() {
		return instr.call()
	}
	return instr
}

func (instr Instruction) DoApplication(e Expression) Expression {
	if instr.IsCallReady() {
		return Apply(instr.call(), e).ForceRequest()
	}

	instr.args = append(instr.args, e)
	return instr.ForceRequest()
}

func instructionHeadEquals(ih1, ih2 InstructionHead) bool {
	return ih1.nArgs == ih2.nArgs && ih1.name == ih2.name && 
		((ih1.action == nil && ih2.action == nil) || (ih1.action != nil && ih2.action != nil))
}

func instructionArgsEquals(ia1, ia2 InstructionArgs) bool {
	ok := len(ia1.args) == len(ia2.args) && cap(ia1.args) == cap(ia2.args)
	if !ok {
		return false
	}

	for i := range ia1.args {
		if !ia1.args[i].Equals(ia2.args[i]) {
			return false
		}
	}
	return true
}

func (instr Instruction) StrictEquals(e Expression) bool {
	instr2, ok := e.(Instruction)
	if !ok {
		return false
	}
	return instructionHeadEquals(instr.InstructionHead, instr2.InstructionHead) &&
		instructionArgsEquals(instr.InstructionArgs, instr2.InstructionArgs)
}

func (instr Instruction) StrictString() string {
	return instr.String()
}

func (instr Instruction) String() string {
	tmp := make([]string, cap(instr.args), instr.nArgs)
	for i, e := range instr.args {
		tmp[i] = e.String()
	}

	for i := len(instr.args); i < instr.nArgs; i++ {
		tmp[i] = "_"
	}

	return "instruction[" + instr.name + " " + strings.Join(tmp, " ") + "]"
}

func (ia InstructionArgs) makeBlankCopy() []Expression {
	ln, c := len(ia.args), cap(ia.args)
	return make([]Expression, ln, c)
}

func (ia InstructionArgs) Copy() InstructionArgs {
	out := ia.makeBlankCopy()
	for i, arg := range ia.args {
		out[i] = arg.Copy()
	}
	return InstructionArgs{ args: out, }
}

func (instr Instruction) UpdateVars(gt int, by int) Expression { 
	newArgs := instr.InstructionArgs.makeBlankCopy()
	for i, arg := range instr.args {
		newArgs[i] = arg.UpdateVars(gt, by)
	}
	return Instruction{
		InstructionHead: instr.InstructionHead,
		InstructionArgs: InstructionArgs{ args: newArgs, },
	}
}

var EqualityInstruction = 
	DefineInstruction("testEquality", 2, func(instr InstructionArgs) Expression {
		e1, e2 := instr.GetArgAtIndex(0), instr.GetArgAtIndex(1)
		if e1.Equals(e2) {
			return TrueFunction
		}
		return FalseFunction
	})

func (ih InstructionHead) MakeInstance() Instruction {
	return Instruction{
		InstructionHead: ih,
		InstructionArgs: InstructionArgs{
			args: make([]Expression, 0, ih.nArgs),
		},
	}
}

func DefineInstruction(name string, nArgs int, action InstructionAction) InstructionHead {
	if nArgs < 0 {
		panic("instruction cannot accept have than 0 arguments.\n")
	}
	return InstructionHead{name: name, nArgs: nArgs, action: action}
}
