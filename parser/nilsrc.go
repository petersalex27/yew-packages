package parser

import "github.com/petersalex27/yew-packages/source"

type nilsrc struct {}

func (nilsrc) SourceLine(line int) (string, source.Status) {
	return "", source.Ok
}

func (nilsrc) GetLineChar() (line, char int) {
	return 0, 0
}

func (nilsrc) NumLines() int {
	return 0
}

func (nilsrc) GetPath() string {
	return ""
}