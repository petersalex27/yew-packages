package types

import (
	expr "github.com/petersalex27/yew-packages/expr"
	"strings"
	"sync"
	"github.com/petersalex27/yew-packages/util/stack"
	"github.com/petersalex27/yew-packages/util"
)

type ruleElement[T Type] struct {
	name string
	rule func(cxt *Context) (T, error)
}

type equivalenceClassMap map[string]Monotyped
type expressionClassMap map[string]expr.Expression

type Context struct {
	contextNumber int32
	varCounter uint32
	equivClasses equivalenceClassMap
	exprClasses expressionClassMap
	stack *stack.Stack[Type]
}

func InheritContext(parent *Context) *Context {
	child := NewContext()

	child.equivClasses = util.CopyMap(parent.equivClasses)
	child.exprClasses = util.CopyMap(parent.exprClasses)
	child.varCounter = parent.varCounter
	child.stack = parent.stack

	return child
}

func NewContext() *Context {
	cxt := new(Context)
	cxt.equivClasses = make(equivalenceClassMap)
	cxt.exprClasses = make(expressionClassMap)
	cxt.stack = stack.NewStack[Type](1 << 5 /*cap=32*/)
	return cxt
}

type RuleID uint

var ruleLock sync.Mutex
const (
	ApplyRule string = "application"
	AbstractRule string = "abstraction"
	ConstructRule string = "construction"
	HeadRule string = "head-separation"
	TailRule string = "tail-separation"
	DisjunctRule string = "disjunction"
	ExpandRule string = "expansion"
	RealizeRule string = "realization"
	ContextRule string = "contextualization"
	InstanceRule string = "instantiation"
	GeneralRule string = "generalization"
)
const (
	ApplyID RuleID = iota
	AbstractID
	ConstructID
	HeadID
	TailID
	DisjunctID
	ExpandID
	RealizeID
	ContextID
	InstanceID
	GeneralID
)
var ruleLookup = map[string]RuleID{
	ApplyRule: ApplyID,
	AbstractRule: AbstractID,
	ConstructRule: ConstructID,
	HeadRule: HeadID,
	TailRule: TailID,
	DisjunctRule: DisjunctID,
	ExpandRule: ExpandID,
	RealizeRule: RealizeID,
	ContextRule: ContextID,
	InstanceRule: InstanceID,
	GeneralRule: GeneralID,
}
var rules = []ruleElement[Type]{
	/*application*/ {
		ApplyRule, func(cxt *Context) (Type, error) {
			ms, e := PopTypes[Monotyped](cxt, 2)
			if e != nil {
				return nil, e
			}
			
			return cxt.Apply(ms[1], ms[0])
		},
	},
	/*abstraction*/ {
		AbstractRule, func(cxt *Context) (Type, error) {
			// get arg
			m, e := cxt.PopMonotype()
			if e != nil {
				return nil, e
			}
			return cxt.Abstract(m), nil
		},
	},
	/*construction*/ {
		ConstructRule, func(cxt *Context) (Type, error) {
			ms, e := PopTypes[Monotyped](cxt, 2)
			if e != nil {
				return nil, e
			}
			
			return cxt.Cons(ms[1], ms[0]), nil
		},
	},
	/*head-separation*/ {
		HeadRule, func(cxt *Context) (Type, error) {
			m, e := cxt.PopMonotype()
			if e != nil {
				return nil, e
			}
			return cxt.Head(m)
		},
	},
	/*tail-separation*/ {
		TailRule, func(cxt *Context) (Type, error) {
			m, e := cxt.PopMonotype()
			if e != nil {
				return nil, e
			}
			return cxt.Tail(m)
		},
	},
	/*disjunction*/ {
		DisjunctRule, func(cxt *Context) (Type, error) {
			m, e := cxt.PopMonotype()
			if e != nil {
				return nil, e
			}

			return cxt.Join(m), nil
		},
	},
	/*expansion*/ {
		ExpandRule, func(cxt *Context) (Type, error) {
			ms, e := PopTypes[Monotyped](cxt, 2)
			if e != nil {
				return nil, e
			}
			
			return cxt.Expansion(ms[1], ms[0]), nil
		},
	},
	/*realization*/ {
		RealizeRule, func(cxt *Context) (Type, error) {
			ms, e := PopTypes[Monotyped](cxt, 3)
			if e != nil {
				return nil, e
			}

			return cxt.Realization(ms[2], ms[1], ms[0]) 
		},
	},
	/*contextualization*/ {
		ContextRule, func(cxt *Context) (Type, error) {
			ps, e := cxt.PopTypesAsPolys(2)
			if e != nil {
				return nil, e
			}
			
			return cxt.Contextualization(ps[1], ps[0]), nil
		},
	},
	/*instantiation*/ {
		InstanceRule, func(cxt *Context) (Type, error) {
			m, eMono := cxt.PopMonotype()
			if eMono != nil {
				return nil, eMono
			}

			p, ePoly := cxt.PopPolytype()
			if ePoly != nil {
				return nil, ePoly
			}
			
			return p.Instantiate(m), nil
		},
	},
	/*generalization*/ {
		GeneralRule, func(cxt *Context) (Type, error) {
			t := cxt.Pop()
			return t.Generalize(), nil
		},
	},
}

func AddRule(name string, rule func(cxt *Context) (Type, error)) (RuleID, error) {
	ruleLock.Lock()
	defer ruleLock.Unlock()

	if _, found := ruleLookup[name]; found {
		return 0, ruleAlreadyDefined(name)
	}

	id := RuleID(len(rules)) // get new rule id
	ruleLookup[name] = id // add to lookup table

	// add rule
	newRule := ruleElement[Type]{name: name, rule: rule}
	rules = append(rules, newRule)
	
	return id, nil
}

func (cxt *Context) Rule(id RuleID) (Type, error) {
	ruleLock.Lock()
	defer ruleLock.Unlock()

	// check that rule exists
	if RuleID(len(rules)) <= id {
		return nil, ruleIdDNE(id)
	}
	// call rule
	return rules[id].rule(cxt)
}

func (cxt *Context) FindExpression(e expr.Expression) expr.Expression {
	if v, ok := e.(expr.Variable); ok {
		if out, found := cxt.exprClasses[v.String()]; found {
			return out
		}
		return v
	}
	return e
}

// returns representative for equiv. class
func (cxt *Context) Find(m Monotyped) Monotyped {
	name, _ := Name(m)
	out, found := cxt.equivClasses[name]
	if !found {
		return m // finds itself
	}
	return out // found itself or representative of type class
}

/*
x :: a = 1
y :: b = f 1.1
f :: t -> u
(+) a b = addInt a b
z = x + y

x: a
a = find(1)
find(1) = Int
y: b
b = f 1.1
f: t -> u == -> t u

*/

// register returns a function that, for each monotype `t` given to the function, 
// creates a map `t -> m`
func (cxt *Context) register(m Monotyped) func(...Monotyped) {
	return func(ts ...Monotyped) {
		for _, t := range ts {
			/*if _, ok := t.(Variable); ok {
				continue
			}*/

			//name, _ := Name(t)
			name := t.String()
			cxt.equivClasses[name] = m
		}
	}
}

func (cxt *Context) registerExpression(e expr.Expression) func(...expr.Expression) {
	return func(exprs ...expr.Expression) {
		for _, exp := range exprs {
			if _, ok := exp.(expr.Variable); !ok {
				continue
			}
			cxt.exprClasses[exp.StrictString()] = e
		}
	}
}

func (cxt *Context) newType(name string, m Monotyped) error {
	if _, found := cxt.equivClasses[name]; found {
		return alreadyDefined(name)
	}
	cxt.equivClasses[name] = m
	return nil
} 

/*
case a and b have no reps.
 -> b reps a; a :-> b, b :-> b
case c reps a, b has no reps.
 ->
*/

// declares that `a` and `b` represent the same type--panics if at least
// one of the types is not a type variable
func (cxt *Context) union(a, b Monotyped) {
	if IsVariable(a) {
		cxt.register(b)(a, b)
	} else if IsVariable(b) {
		cxt.register(a)(a, b)
	} else {
		panic("tried to join two distinct types")
	}
}

func (cxt *Context) expressionUnion(e1, e2 expr.Expression) {
	if _, ok := e1.(expr.Variable); ok {
		cxt.registerExpression(e1)(e1, e2)
	} else if _, ok := e2.(expr.Variable); ok {
		cxt.registerExpression(e2)(e1, e2)
	} else {
		panic("tried to join two distinct expressions")
	}
}

// see Unify for description
func (cxt *Context) unify(a, b Monotyped) bool {
	ta := cxt.Find(a)
	tb := cxt.Find(b)

	if IsVariable(ta) || IsVariable(tb) {
		cxt.union(ta, tb)
		return true
	}

	aName, aParams := Split(ta)
	bName, bParams := Split(tb)
	if aName == bName && len(aParams) == len(bParams) {
		for i := range aParams {
			if !cxt.unify(aParams[i], bParams[i]) {
				return false
			}
		}
	} else {
		return false
	}

	return true
}

func IsKindVariable(e expr.Expression) bool {
	_, ok := e.(expr.Variable)
	return ok
}

func (cxt *Context) tryToEquateExpressions(e1, e2 expr.Expression) bool {
	// f: [a; n+m] -> ([a; n], [a; m])
	// f $ (x: [a; 6])
	// => n+m=6 => n=6-m
	// => m=6-n
	return true
}

func (cxt *Context) unifyExpression(a, b expr.Expression) bool {
	ta := cxt.FindExpression(a)
	tb := cxt.FindExpression(b)

	if IsKindVariable(ta) || IsKindVariable(tb) {
		cxt.expressionUnion(ta, tb)
		return true
	}

	var la, ra, lb, rb expr.Expression
	sa, okA := ta.(expr.SplitableExpression)
	sb, okB := tb.(expr.SplitableExpression)
	if !(okA && okB) {
		return cxt.tryToEquateExpressions(ta, tb)
	} 
	la, ra = sa.Split()
	lb, rb = sb.Split()
	return cxt.unifyExpression(la, lb) && cxt.unifyExpression(ra, rb)
}

// Unify declares that two monotypes, i.e., tbe arguments given for `a` and `b`, are the same
// type. On success, nil is returned; otherwise an error is returned.
func (cxt *Context) Unify(a, b Monotyped) error {
	if !cxt.unify(a, b) {
		return typeMismatch(a, b)
	}
	return nil
}

func (cxt *Context) StringClasses() string {
	var builder strings.Builder
	for k, v := range cxt.equivClasses {
		builder.WriteString(k + " : " + v.String() + "\n")
	}
	return builder.String()
}
