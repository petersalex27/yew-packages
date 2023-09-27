package errors

import (
	"github.com/fatih/color"
)

type subtype string

func (s subtype) IsNone() bool {
	return string(s) == ""
}

func (s subtype) subTypeString() string {
	if s.IsNone() {
		return ""
	}
	return color.MagentaString(" %s", s.String())
}

func (s subtype) String() string {
	if s.IsNone() {
		return string(s)
	}
	return "(" + string(s) + ")"
}
