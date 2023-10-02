package status

type Status int

const (
	Ok Status = iota
	EndOfTokens
	StackEmpty
	Error
	EndAction
	EndOfParse

	NoAction

	ReductionOverwrite
	ReductionDNE
)

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

// EndAction -> Ok; else, stat -> stat (Ok -> Ok too)
func (stat Status) MakeOk() Status {
	if stat == EndAction {
		return Ok
	}
	return stat
}
