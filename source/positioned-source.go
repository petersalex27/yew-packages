package source

type PositionedSource interface {
	Source
	SetLineChar(line, char int)
	AdvanceLine() Status
	AdvanceChar() (char byte, stat Status)
	UnadvanceChar() (stat Status)
	WhitespaceLength() int
	ReadUntil(char byte) (string, Status)
}

/*
SetAndWrapLineChar sets the line and char number of an implementer of PositionedSource.
A line and char are set
*/
func SetAndWrapLineChar(p PositionedSource, line, char int) Status {
	str, stat := p.SourceLine(line)
	for len(str) < char && stat.IsOk() {
		line = line + 1 // char overflow: next line
		char = char - len(str) // chars remaining after overflow
		str, stat = p.SourceLine(line) // if error, will be caught at start of loop
	}

	if stat.NotOk() && !stat.IsEof() {
		return stat
	}
	// line and char are set on stat == Eof and stat == Ok
	p.SetLineChar(line, char)
	return stat
}

func SkipWhitespace(p PositionedSource) {
	wslen := p.WhitespaceLength()
	line, char := p.GetLineChar()
	p.SetLineChar(line, char + wslen)
}

func GetLeadingWhitespace(p PositionedSource) (string, Status) {
	line, char := p.GetLineChar()
	if char != 1 {
		return "", Ok
	}
	src, stat := p.SourceLine(line)
	if stat.NotOk() {
		return "", stat
	}
	wslen := p.WhitespaceLength()
	res := src[:wslen]
	p.SetLineChar(line, char + wslen)
	return res, Ok
}

func UngetChars(p PositionedSource, nchars int) Status {
	if nchars < 0 {
		panic("illegal argument: nchars < 0")
	}

	line, char := p.GetLineChar()
	for char < nchars {
		if line - 1 <= 0 {
			return OutOfBounds
		}
		line = line - 1
		nchars = nchars - char 

		//p.SetLineChar(line, 1)
		src, stat := p.SourceLine(line)
		if stat.NotOk() {
			return stat
		}
		char = len(src) - nchars
		p.SetLineChar(line, char)
	}
	return Ok
}

func ReadThrough(p PositionedSource, char byte) (string, Status) {
	res, stat := p.ReadUntil(char)
	if stat.NotOk() {
		return res, stat
	}

	var last byte
	last, stat = p.AdvanceChar()
	if stat.NotOk() {
		return res, stat
	}

	return res + string(last), Ok
}

