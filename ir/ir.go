package ir
/*
import (
	"fmt"
	"sync"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/symbol"
	"github.com/petersalex27/yew-packages/table"
	"github.com/petersalex27/yew-packages/types"
)

// create a translation/replacement `irty ::= constType`
func (translator *Translator[T]) CreateTypeTranslation(constType types.Constant[T], irty irtypes.Type) Status {
	name := constType.GetName()
	// check if adding will overwrite
	_, found := translator.typeTable.Get(name)
	if found { // prevent overwrite
		return RedefinedTypeTranslation
	}

	// add translation to table
	translator.typeTable.Add(name, irty)
	return Ok
}

func (translator *Translator[T]) AssembleTypeTranslation(types.Type[T]) (irtypes.Type, Status) {
	translator.
}

func TranslateFunction[T nameable.Nameable](fn expr.Function[T]) *ir.Func {
	fn := ir.NewFunc("")
	blk := fn.NewBlock("")
}

// Data related to translating source language to LLVM-IR
type Translator[T nameable.Nameable] struct {

	masterTableLock sync.Mutex
	masterTable *table.Table[T]
	typeTable *table.Table[irtypes.Type] // stack of type tables
	current *ir.Module // LLVM Module for current source file 
	undecidedTables *table.Table[T] // stack of tables for symbols with not-yet-known classifications
	//tables *table.TableStack[TranslatedSymbol[T]]
	tables *table.TableStack[symbol.Symbol[T]] // stack of symbol tables
	sourceFiles map[string]*ir.Module // map from source file path to LLVM Module
}

func (translator *Translator[T]) Demangle(format string, name T, args ...any) {
	translator.masterTableLock.Lock()
	defer translator.masterTableLock.Unlock()

	ln := translator.masterTable.Len()
	named := fmt.Sprintf(format + name.GetName() + "_%d", append(args, ln)...)
}

type TranslatedSymbol[T nameable.Nameable] struct {

	irtypes.Type
}

// Creates new Translator and inits applicable members w/ small default values
//
// TODO: allow more control over member inits
func NewTranslator[T any]() *Translator[T] {
	translator := new(Translator[T])
	translator.tables = table.NewTableStack[T]()
	translator.sourceFiles = make(map[string]*ir.Module)
	return translator
}

// will be removed soon--example llvm "hello, world!" program for reference 
func HelloWorld() string {
	mod := ir.NewModule()
	hello := constant.NewCharArrayFromString("Hello, world!\n\x00")
	str := mod.NewGlobalDef("str", hello)
	// Add external function declaration of puts.
	puts := mod.NewFunc("puts", irtypes.I32, ir.NewParam("", irtypes.NewPointer(irtypes.I8)))
	main := mod.NewFunc("main", irtypes.I32)
	entry := main.NewBlock("")
	// Cast *[15]i8 to *i8.
	zero := constant.NewInt(irtypes.I64, 0)
	gep := constant.NewGetElementPtr(hello.Typ, str, zero, zero)
	entry.NewCall(puts, gep)
	entry.NewRet(constant.NewInt(irtypes.I32, 0))
	return fmt.Sprint(mod)
}

func (translator *Translator[T]) Function() {

}
*/