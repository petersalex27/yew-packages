package ir

type Status int

const (
	Ok Status = iota
	RedefinedTypeTranslation
)

func (stat Status) IsOk() bool { return stat == Ok }

func (stat Status) NotOk() bool { return !stat.IsOk() }

func (stat Status) Is(stat2 Status) bool { return stat == stat2 }