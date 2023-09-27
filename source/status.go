package source

type Status byte
const (
	Eof Status = iota // end of file
	Eol	// end of line
	Ok // ok
	OutOfBounds // out of bounds char

	BadLineNumber // line number is illegal
	BadCharNumber // char number is illegal

	BadCharStart // start char < 1
	BadCharEnd   // end char < 1
	BadCharRange // start char > end char

	Bad // bad
)

func (stat Status) String() string {
	switch stat {
	case Eof:
		return "Eof"
	case Eol:
		return "Eol"
	case Ok:
		return "Ok"
	case OutOfBounds:
		return "OutOfBounds"
	case BadLineNumber:
		return "BadLineNumber"
	case BadCharNumber:
		return "BadCharNumber"
	case BadCharStart:
		return "BadCharStart"
	case BadCharRange:
		return "BadCharRange"
	case Bad:
		return "Bad"
	default:
		return "StatusUndefined"
	}
}

func (stat Status) Equals(stat2 Status) bool { return stat.Is(stat2) }

func (stat Status) Is(stat2 Status) bool { return stat == stat2 }

func (stat Status) IsOk() bool { return stat.Is(Ok) }

func (stat Status) NotOk() bool { return !stat.Is(Ok) }

func (stat Status) IsOOB() bool { return stat.Is(OutOfBounds)}

func (stat Status) IsEof() bool { return stat.Is(Eof) }

func (stat Status) IsEol() bool { return stat.Is(Eol) }

