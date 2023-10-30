package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
	"github.com/petersalex27/yew-packages/stringable"
)

// == production interface ====================================================

// Interface for things that represent production rules or things that are like
// production rules.
//
// A standard production rule, in cs theory, the form:
//
//	replacement ::= handle
type Productions interface {
	stringable.Stringable
	do(nodes ...ast.Ast) ast.Ast // performs replacement
}

type Production struct{ Productions }

func (production Production) call(p *parser, n uint, handle ...ast.Ast) status.Status {
	// pop stack (this removes `nodes`)
	p.stack.Clear(n) // must be called before pushing reduction result
	// do reduction action
	result := production.do(handle...)
	p.stack.Push(result)
	return status.Ok
}

// represents production rule
type ProductionFunction func(handle ...ast.Ast) (replacement ast.Ast)

// give a name to the
func (f ProductionFunction) GiveName(name string) NamedProduction {
	return NamedProduction{
		Name:               name,
		ProductionFunction: f,
	}
}

func (f ProductionFunction) do(nodes ...ast.Ast) ast.Ast {
	return f(nodes...)
}

func (f ProductionFunction) String() string {
	return "production_function"
}

type NamedProduction struct {
	Name string
	ProductionFunction
}

func (f NamedProduction) do(nodes ...ast.Ast) ast.Ast {
	return f.ProductionFunction(nodes...)
}

func (f NamedProduction) String() string {
	return f.Name
}
