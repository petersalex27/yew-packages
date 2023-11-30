package parser

import "github.com/petersalex27/yew-packages/source"

// for mocking
type EmptySource struct{}

func (EmptySource) SourceLine(line int) (string, source.Status) {
	return "", source.Ok
}

func (EmptySource) GetLineChar() (line, char int) {
	return 0, 0
}

func (EmptySource) NumLines() int {
	return 0
}

func (EmptySource) GetPath() string {
	return ""
}
