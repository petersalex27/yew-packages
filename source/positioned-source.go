package source

import (
	"strings"
)

type PositionedSource interface {
	Source
	SetLineChar(line, char int)
	AdvanceChar() (char byte, stat Status)
	UnadvanceChar() (stat Status)
	WhitespaceLength() int
	ReadUntil(byte) (string, Status)
}

type EscapeMap interface {
	GetMapped(remainingLine string) (result string, found bool, readTo int)
}

type EscapeMap_s struct {
	LookupSize uint
	Map map[string]string
}

func (em EscapeMap_s) GetMapped(remainingLine string) (result string, found bool, readTo int) {
	result, found, readTo = "", false, int(em.LookupSize)
	
	if len(remainingLine) < int(em.LookupSize) {
		return
	}

	result, found = em.Map[remainingLine[:em.LookupSize]]
	return
}

func readEscapable(em EscapeMap, line string, end byte, esc byte) (string, int, Status) {
	index := 0
	escaped := false
	var builder strings.Builder

	for ; index < len(line); {
		c := line[index]
		if escaped {
			escaped = false

			var res string
			var found bool
			var readTo int
			if em == nil {
				found = false
			} else {
				res, found, readTo = em.GetMapped(line[index:])
			}

			if !found {
				return "", index, Bad
			}
			index = index + readTo
			builder.WriteString(res)
		} else if c == end {
			return builder.String(), index, Ok
		} else if c == esc {
			escaped = true
			index = index + 1
		} else {
			index = index + 1
			builder.WriteByte(c)
		}
	}
	// `end` not found
	return "", index, Eol
}

func ReadUntil(em EscapeMap, p PositionedSource, end byte, esc byte) (string, Status) {
	line, char := p.GetLineChar()
	remainingLine, stat := CurrentLine(p)
	if stat.NotOk() {
		return "", stat
	}

	res, readLength, stat := readEscapable(em, remainingLine, end, esc)
	if stat.NotOk() {
		return "", stat
	}

	p.SetLineChar(line, char+readLength)
	if _, stat = p.AdvanceChar(); stat.NotOk() { // remove closing `end`
		return "", stat
	}
	return res, stat
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

func ReadThrough(p PositionedSource, end byte) (string, Status) {
	res, stat := p.ReadUntil(end)
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

