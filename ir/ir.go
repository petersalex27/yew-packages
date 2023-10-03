package ir

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

// will be removed soon
func HelloWorld() string {
	mod := ir.NewModule()
	hello := constant.NewCharArrayFromString("Hello, world!\n\x00")
	str := mod.NewGlobalDef("str", hello)
	// Add external function declaration of puts.
	puts := mod.NewFunc("puts", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	main := mod.NewFunc("main", types.I32)
	entry := main.NewBlock("")
	// Cast *[15]i8 to *i8.
	zero := constant.NewInt(types.I64, 0)
	gep := constant.NewGetElementPtr(hello.Typ, str, zero, zero)
	entry.NewCall(puts, gep)
	entry.NewRet(constant.NewInt(types.I32, 0))
	return fmt.Sprint(mod)
}