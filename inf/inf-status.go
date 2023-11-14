package inf

import "fmt"

type Status byte

const (
	// default status
	Ok Status = iota
	// constants at same positions in unification did not match
	ConstantMismatch
	// application monotype kind did not have same number of type params as other
	// monotype being unified
	ParamLengthMismatch
	// unification of variable and monotype failed because variable occured w/in
	// monotype
	OccursCheckFailed
	// unification of variables succeeded, so signals that there is nothing left 
	// to unify
	skipUnify
)

func (stat Status) String() string {
	switch stat {
	case Ok:
		return "Ok"
	case ConstantMismatch:
		return "ConstantMismatch"
	case ParamLengthMismatch:
		return "ParamLengthMismatch"
	case OccursCheckFailed:
		return "OccursCheckFailed"
	case skipUnify:
		return "skipUnify"
	default:
		return fmt.Sprintf("Status(%d)", stat)
	}
}

func (stat Status) IsOk() bool {
	return stat == Ok
}

func (stat Status) NotOk() bool {
	return stat != Ok
}

func (stat Status) Is(stat2 Status) bool {
	return stat == stat2
}