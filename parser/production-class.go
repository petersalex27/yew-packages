package parser

import "github.com/petersalex27/yew-packages/parser/ast"

// == production class =========================================================

// class of productions
type productionClass []productionInterface

// true iff class has no members
func (class productionClass) isEmpty() bool {
	return len(class) == 0
}

// add a production member to the production class (i.e., the receiver)
func (class productionClass) add(production productionInterface) productionClass {
	if class.isEmpty() {
		return productionClass{production} // create new class
	}
	return append(class, production) // add to existing class
}

// == production class map =====================================================

// collection of production classes that are classified based on the last term
// of their pattern
type productionClassMap map[ast.Type]productionClass

// create a new production class map
func newProductionClassifier() (rcm *productionClassMap) {
	rcm = new(productionClassMap)
	*rcm = make(productionClassMap)
	return
}

// adds `production` to a production class based on its PatternInterface
func (rcm *productionClassMap) classify(production productionInterface) {
	pat := production.getPattern()
	// grab last type in pattern
	last, isEmpty := pat.Last()
	if isEmpty { // nothing to classify
		return
	}

	// get class
	class := (*rcm)[last]
	// add production to class and (re-)map
	(*rcm)[last] = class.add(production)
}

// classify productions in production order
func (rcm *productionClassMap) classifyReductions(productions productionOrder) {
	for _, rule := range productions.rules {
		rcm.classify(rule)
	}
}

// if there exists a map b/w `ty` and some rule class `rc`, return `rc` along
// with `true`; else, return `nil` along with `false`
func (rcm *productionClassMap) getClass(ty ast.Type) (rc productionClass, exists bool) {
	rc, exists = (*rcm)[ty]
	return
}
