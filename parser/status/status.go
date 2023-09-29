package status

type Status int

const (
	Ok Status = iota
	EndOfTokens
	StackEmpty

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