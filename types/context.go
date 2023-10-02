package types

import (
	"strings"
	"sync"

	expr "github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/util"
	"github.com/petersalex27/yew-packages/util/stack"
)

type ruleElement[T Type[U], U nameable.Nameable] struct {
	name string
	rule func(cxt *Context[U]) (T, error)
}

type equivalenceClassMap[T nameable.Nameable] map[string]Monotyped[T]
type expressionClassMap map[string]expr.Expression

type Context[T nameable.Nameable] struct {
	contextNumber int32
	varCounter uint32
	makeName func(string)T
	equivClasses equivalenceClassMap[T]
	exprClasses expressionClassMap
	stack *stack.Stack[Type[T]]
}

func InheritContext[T nameable.Nameable](parent *Context[T]) *Context[T] {
	child := NewContext[T]()

	child.equivClasses = util.CopyMap(parent.equivClasses)
	child.exprClasses = util.CopyMap(parent.exprClasses)
	child.varCounter = parent.varCounter
	child.stack = parent.stack

	return child
}

func NewContext[T nameable.Nameable]() *Context[T] {
	cxt := new(Context[T])
	cxt.equivClasses = make(equivalenceClassMap[T])
	cxt.exprClasses = make(expressionClassMap)
	cxt.stack = stack.NewStack[Type[T]](1 << 5 /*cap=32*/)
	return cxt
}

func (cxt *Context[T]) SetNameMaker(f func(string)T) *Context[T] {
	cxt.makeName = f
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
func rules_[T nameable.Nameable]() []ruleElement[Type[T], T] {
	return []ruleElement[Type[T], T] {
		/*application*/ {
			ApplyRule, func(cxt *Context[T]) (Type[T], error) {
				ms, e := PopTypes[Monotyped[T]](cxt, 2)
				if e != nil {
					return nil, e
				}
				
				return cxt.Apply(ms[1], ms[0])
			},
		},
		/*abstraction*/ {
			AbstractRule, func(cxt *Context[T]) (Type[T], error) {
				// get arg
				m, e := cxt.PopMonotype()
				if e != nil {
					return nil, e
				}
				return cxt.Abstract(m), nil
			},
		},
		/*construction*/ {
			ConstructRule, func(cxt *Context[T]) (Type[T], error) {
				ms, e := PopTypes[Monotyped[T]](cxt, 2)
				if e != nil {
					return nil, e
				}
				
				return cxt.ConsRule(ms[1], ms[0]), nil
			},
		},
		/*head-separation*/ {
			HeadRule, func(cxt *Context[T]) (Type[T], error) {
				m, e := cxt.PopMonotype()
				if e != nil {
					return nil, e
				}
				return cxt.Head(m)
			},
		},
		/*tail-separation*/ {
			TailRule, func(cxt *Context[T]) (Type[T], error) {
				m, e := cxt.PopMonotype()
				if e != nil {
					return nil, e
				}
				return cxt.Tail(m)
			},
		},
		/*disjunction*/ {
			DisjunctRule, func(cxt *Context[T]) (Type[T], error) {
				m, e := cxt.PopMonotype()
				if e != nil {
					return nil, e
				}

				return cxt.JoinRule(m), nil
			},
		},
		/*expansion*/ {
			ExpandRule, func(cxt *Context[T]) (Type[T], error) {
				ms, e := PopTypes[Monotyped[T]](cxt, 2)
				if e != nil {
					return nil, e
				}
				
				return cxt.Expansion(ms[1], ms[0]), nil
			},
		},
		/*realization*/ {
			RealizeRule, func(cxt *Context[T]) (Type[T], error) {
				ms, e := PopTypes[Monotyped[T]](cxt, 3)
				if e != nil {
					return nil, e
				}

				return cxt.Realization(ms[2], ms[1], ms[0]) 
			},
		},
		/*contextualization*/ {
			ContextRule, func(cxt *Context[T]) (Type[T], error) {
				ps, e := cxt.PopTypesAsPolys(2)
				if e != nil {
					return nil, e
				}
				
				return cxt.Contextualization(ps[1], ps[0]), nil
			},
		},
		/*instantiation*/ {
			InstanceRule, func(cxt *Context[T]) (Type[T], error) {
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
			GeneralRule, func(cxt *Context[T]) (Type[T], error) {
				t := cxt.Pop()
				return t.Generalize(cxt), nil
			},
		},
	}
}

/*
func AddRule[T nameable.Nameable](name string, rule func(cxt *Context[T]) (Type[T], error)) (RuleID, error) {
	ruleLock.Lock()
	defer ruleLock.Unlock()

	if _, found := ruleLookup[name]; found {
		return 0, ruleAlreadyDefined(name)
	}

	id := RuleID(len(rules_[T]())) // get new rule id
	ruleLookup[name] = id // add to lookup table

	// add rule
	newRule := ruleElement[Type[T]]{name: name, rule: rule}
	rules = append(rules, newRule)
	
	return id, nil
}

func (cxt *Context[T]) Rule(id RuleID) (Type[T], error) {
	ruleLock.Lock()
	defer ruleLock.Unlock()

	// check that rule exists
	if RuleID(len(rules)) <= id {
		return nil, ruleIdDNE(id)
	}
	// call rule
	return rules[id].rule(cxt)
}*/

func (cxt *Context[T]) FindExpression(e expr.Expression) expr.Expression {
	if v, ok := e.(expr.Variable); ok {
		if out, found := cxt.exprClasses[v.String()]; found {
			return out
		}
		return v
	}
	return e
}

// returns representative for equiv. class
func (cxt *Context[T]) Find(m Monotyped[T]) Monotyped[T] {
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
func (cxt *Context[T]) register(m Monotyped[T]) func(...Monotyped[T]) {
	return func(ts ...Monotyped[T]) {
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

func (cxt *Context[T]) registerExpression(e expr.Expression) func(...expr.Expression) {
	return func(exprs ...expr.Expression) {
		for _, exp := range exprs {
			if _, ok := exp.(expr.Variable); !ok {
				continue
			}
			cxt.exprClasses[exp.StrictString()] = e
		}
	}
}

func (cxt *Context[T]) newType(name string, m Monotyped[T]) error {
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
func (cxt *Context[T]) union(a, b Monotyped[T]) {
	if IsVariable(a) {
		cxt.register(b)(a, b)
	} else if IsVariable(b) {
		cxt.register(a)(a, b)
	} else {
		panic("tried to join two distinct types")
	}
}

func (cxt *Context[T]) expressionUnion(e1, e2 expr.Expression) {
	if _, ok := e1.(expr.Variable); ok {
		cxt.registerExpression(e1)(e1, e2)
	} else if _, ok := e2.(expr.Variable); ok {
		cxt.registerExpression(e2)(e1, e2)
	} else {
		panic("tried to join two distinct expressions")
	}
}

// see Unify for description
func (cxt *Context[T]) unify(a, b Monotyped[T]) bool {
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

func (cxt *Context[T]) tryToEquateExpressions(e1, e2 expr.Expression) bool {
	// f: [a; n+m] -> ([a; n], [a; m])
	// f $ (x: [a; 6])
	// => n+m=6 => n=6-m
	// => m=6-n
	return true
}

func (cxt *Context[T]) unifyExpression(a, b expr.Expression) bool {
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
func (cxt *Context[T]) Unify(a, b Monotyped[T]) error {
	if !cxt.unify(a, b) {
		return typeMismatch[T](a, b)
	}
	return nil
}

func (cxt *Context[T]) StringClasses() string {
	var builder strings.Builder
	for k, v := range cxt.equivClasses {
		builder.WriteString(k + " : " + v.String() + "\n")
	}
	return builder.String()
}
