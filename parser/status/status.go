package status

import (
	"fmt"
	"sync"
)

type Status int

const (
	// staus signaling operation had expected behavior
	Ok Status = iota
	// status signaling a failed operation due to no more look-ahead tokens
	EndOfTokens 
	// status signaling a failed operation due to an empty stack
	StackEmpty 
	// status signaling an error occ.
	Error
	// status signaling that the parser should stop applying reductions
	// from the current rule set
	EndAction 
	// status signaling that parsing should end entirely but the operation returing
	// this status had expected behavior
	EndOfParse
	// status signaling that the parser should shift the next look-ahead token
	DoShift
)

var customStatLock sync.Mutex

var customStats = map[Status]string{
	Ok: "Ok",
	EndOfTokens: "EndOfTokens",
	StackEmpty: "StackEmpty",
	Error: "Error",
	EndAction: "EndAction",
	EndOfParse: "EndOfParse",
	DoShift: "DoShift",
}

func RegisterCustomStat(statNum int, name string) {
	customStatLock.Lock()
	
	if s, found := customStats[Status(statNum)]; found {
		customStatLock.Unlock()
		panic(fmt.Sprintf("cannot overwrite existing stat (stat#=%d: %s)", statNum, s))
	}

	customStats[Status(statNum)] = name
	customStatLock.Unlock()
} 

func (stat Status) String() string {
	customStatLock.Lock()
	defer customStatLock.Unlock()

	s, found := customStats[stat]
	if found {
		return s
	}
	return fmt.Sprintf("Status(%d)", int(stat))
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

func (stat Status) EndParse() bool {
	return stat == EndOfParse
}

// EndAction -> Ok; DoShift -> Ok; else, stat -> stat (Ok -> Ok too)
func (stat Status) MakeOk() Status {
	if stat == EndAction || stat == DoShift {
		return Ok
	}
	return stat
}
