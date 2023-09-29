package ast

import "math"

// Type is the type of ast node types; math.MaxUint is reserved! 
type Type uint

const None Type = math.MaxUint

type Ast interface {
	NodeType() Type
}