package source

type StaticSource interface {
	SourceLine(line int) (string, Status)
	GetPath() string
	NumLines() int
}

type Source interface {
	// should return line as string or an empty string and one of two Stats: Eof when 
	// line > number of lines; otherwise, Ok
	GetLineChar() (line, char int)
	StaticSource
}

func CurrentLineLength(s Source) int {
	str, stat := CurrentLine(s)
	if stat.NotOk() {
		return 0
	}
	return len(str)
}

func CurrentLine(s Source) (string, Status) {
	line, _ := s.GetLineChar()
	return s.SourceLine(line)
}

func CurrentChar(s Source) (byte, Status) {
	line, char := s.GetLineChar()
	return GetSourceChar(s, line, char)
}

// returns 0 on EOF or out-of-bounds-line, panics on out-of-bounds-char
func GetSourceChar(s Source, line, char int) (byte, Status) {
	src, stat := s.SourceLine(line)
	if stat.NotOk() {
		return 0, stat
	}

	if len(src) < char {
		if len(src) == char - 1 {
			return '\n', Eol
		}
		return 0, OutOfBounds
	}

	return src[char-1], Ok
}

// returns a slice from src[char_start-1:char_end] where 
// src, _ := s.SourceLine(line). When char_end < 0, then returns a 
// slice src[char_start-1:]
func GetSourceSlice(s Source, line, char_start, char_end int) (string, Status) {
	if char_start < 1  {
		return "", BadCharStart
	} else if char_end == 0 {
		return "", BadCharEnd
	} else if char_end >= 0 && char_end < char_start {
		return "", BadCharRange
	}

	src, stat := s.SourceLine(line)
	if stat.NotOk() {
		return "", stat
	}

	if len(src) < char_start || len(src) < char_end {
		return "", OutOfBounds
	}

	if char_end < 0 {
		return src[char_start-1:], Ok
	}
	return src[char_start-1:char_end], Ok
}