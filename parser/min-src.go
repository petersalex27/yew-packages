package parser

import "github.com/petersalex27/yew-packages/source"

type MinSource struct {
	path string
	src []string
}

func MakeSource(path string, src ...string) MinSource {
	return MinSource{ path: path, src: src, }
}

func (src MinSource) SourceLine(line int) (string, source.Status) {
	if len(src.src) < line || line <= 0 {
		return "", source.BadLineNumber
	}
	return src.src[line-1], source.Ok
}

func (src MinSource) GetPath() string { return src.path }

func (src MinSource) NumLines() int { return len(src.src) }