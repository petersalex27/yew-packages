package parser

// ordered collection of production rules
type ProductionOrder struct {
	// production rules
	rules []productionInterface
	// Maps last type in production rule pattern to a subset of production rules.
	// The purpose of this is to reduce the number of rules that have to be
	// searched when trying to match a handle to a rule's pattern.
	//
	// NOTE: expiremental
	classes *productionClassMap
	// when there exists no handle (some number of top elements from parse stack)
	// matching any available production rule's R.H.S., there are two options
	// based on the value of elseShift:
	//  - (1), elseShift == false: throw syntax error
	//  - (2), elseShift == true: shift next token
	elseShift bool
}

// join zero or more production orders into one production order
//
// the order that each production order appears in the arguments is the order
// they are unified in
func Union(sets ...ProductionOrder) (unified ProductionOrder) {
	unified = ProductionOrder{
		rules:     []productionInterface{},
		classes:   newProductionClassifier(),
		elseShift: false,
	}
	for _, set := range sets {
		unified.classes.classifyReductions(set)
		unified.rules = append(unified.rules, set.rules...)
		unified.elseShift = unified.elseShift || set.elseShift
	}
	return
}

// define an order for `productions`
func Order(productions ...productionInterface) (out ProductionOrder) {
	out.rules = append([]productionInterface{}, productions...)
	out.classes = newProductionClassifier()
	out.classes.classifyReductions(out)
	return out
}
