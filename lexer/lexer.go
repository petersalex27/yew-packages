package lexer

import (
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-packages/token"
)

func (lex *Lexer) AddError(e error) {
	lex.errors = append(lex.errors, e)
}

func (lex *Lexer) GetErrors() []error {
	return lex.errors
}

func NewLexer(whitespace *regexp.Regexp, caps ...int) *Lexer {
	lex := new(Lexer)
	*lex = Lexer{
		whitespace: whitespace,
		line:       1,
		char:       1,
		source: nil,
		tokens: nil,
		errors:	 	make([]error, 0),
	}
	if len(caps) > 2 {
		lex.errors = make([]error, 0, caps[2])
	}
	if len(caps) > 1 {
		lex.tokens = make([]token.Token, 0, caps[1])
	}
	if len(caps) > 0 {
		lex.source = make([]string, 0, caps[0])
	}
	return lex
}

// panics if lex.source is already set
func (lex *Lexer) SetSource(src []string) {
	if len(lex.source) != 0 {
		panic("cannot reset lexer's source")
	}
	lex.source = src
}

// panics if lex.source is already set
func (lex *Lexer) SetPath(path string) {
	if lex.path != "" {
		panic("cannot reset lexer's path")
	}
	lex.path = path
}

func Lex(path string, whitespace *regexp.Regexp) (*Lexer, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	bytes, readError := io.ReadAll(f)
	f.Close()
	if readError != nil {
		return nil, readError
	}

	avgTokenLen, avgLineLen := 5, 35 // this is just an estimate

	sourceCap := len(bytes)/avgLineLen
	tokensCap := len(bytes)/avgTokenLen
	errorsCap := 1

	lex := NewLexer(whitespace, sourceCap, tokensCap, errorsCap)
	lex.path = path

	oldI := 0
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == '\n' {
			lex.source = append(lex.source, string(bytes[oldI:i+1]))
			i++
			oldI = i
		}
	}
	lex.source = append(lex.source, string(bytes[oldI:]))
	return lex, nil
}

func (lex *Lexer) GetPath() string {
	return lex.path
}

func (lex *Lexer) PositionStatus() source.Status {
	if lex.line < 1 || lex.line > len(lex.source) {
		return source.BadLineNumber
	} else if lex.char < 1 {
		return source.BadCharNumber
	} else if lex.char > len(lex.source[lex.line-1]) {
		if lex.char != 1 || len(lex.source[lex.line-1]) > 0 {
			return source.BadCharNumber
		}
	}
	return source.Ok
}

func (lex *Lexer) SetLineChar(line, char int) {
	lex.line = line
	lex.char = char
}

func (lex *Lexer) AdvanceLine() source.Status {
	if lex.line + 1 > lex.NumLines() {
		return source.Eof
	} 
	lex.line = lex.line + 1
	lex.char = 1
	return source.Ok
}

func (lex *Lexer) UnadvanceChar() (stat source.Status) {
	return lex.stepDirection(-1)
}

func (lex *Lexer) AdvanceChar() (char byte, stat source.Status) {
	char, stat = source.GetSourceChar(lex, lex.line, lex.char)
	if stat.IsEol() {
		stat = lex.AdvanceLine()
		if stat.NotOk() {
			return 0, stat
		}
	} else if stat.IsOk() {
		lex.char = lex.char + 1
	} else {
		return 0, stat
	}
	return
}

func (lex *Lexer) ReadUntil(char byte) (string, source.Status) {
	var builder strings.Builder
	c, stat := source.CurrentChar(lex)
	for ; c != char && stat.IsOk(); c, stat = source.CurrentChar(lex) {
		in, _ := lex.AdvanceChar()
		builder.WriteByte(in)
	}
	return builder.String(), stat
}

func RegexMatch(r *regexp.Regexp, s string) (string, int) {
	loc := r.FindStringIndex(s)
	if loc == nil { // no match
		return "", 0
	}
	if loc[0] != 0 {
		return "", 0
	}
	return s[:loc[1]], loc[1]
}

func (lex *Lexer) WhitespaceLength() int {
	line, stat := source.GetSourceSlice(lex, lex.line, lex.char, -1)
	if stat.NotOk() {
		return 0
	}
	_, len := RegexMatch(lex.whitespace, string(line))
	return len
}

func (lex *Lexer) GetLineChar() (line, char int) {
	return lex.line, lex.char
}

func (lex *Lexer) SourceLine(line int) (string, source.Status) {
	if len(lex.source) < line || len(lex.source) == 0 {
		return "", source.Eof
	}
	return lex.source[line-1], source.Ok
}

func (lex *Lexer) PushToken(t token.Token) {
	lex.tokens = append(lex.tokens, t)
}

func (lex *Lexer) RemainingLine() (line string, isEndOfLine bool) {
	line = lex.currentLine()
	if line == "" {
		isEndOfLine = true
		return
	}

	if len(line) < lex.char {
		isEndOfLine = true
		line = ""
	} else {
		isEndOfLine = false
		line = line[lex.char-1:]
	}

	return
}

func (lex *Lexer) currentLine() string {
	if len(lex.source) < lex.line {
		return ""
	}
	return lex.source[lex.line-1]
}

func (lex *Lexer) Peek() (nextChar byte, isEof bool) {
	var stat source.Status
	nextChar, stat = source.CurrentChar(lex)
	isEof = stat.IsEof()
	return
}

func (lex *Lexer) IsEndOf() (isEndOfLine bool, isEndOfFile bool) {
	isEndOfFile = lex.IsEndOfFile()
	if !isEndOfFile {
		isEndOfLine = lex.char > len(lex.source[lex.line-1])
		if isEndOfLine {
			isEndOfFile = lex.line+1 > len(lex.source)
		}
	} else {
		isEndOfLine = true
	}
	return
}

func (lex *Lexer) IsEndOfFile() bool {
	return lex.line > len(lex.source)
}

func (lex *Lexer) advanceLine() {
	lex.line = lex.line + 1
	lex.char = 1
}

func (lex *Lexer) step() { lex.stepn(1) }

func (lex *Lexer) stepn(n uint32) {
	lex.char = lex.char + int(n)
	for lex.char > len(lex.source[lex.line-1]) {
		overflow := lex.char - len(lex.source[lex.line-1])
		lex.char = overflow
		lex.line = lex.line + 1
		if lex.line > len(lex.source) {
			return
		}
	}
}

func (lex *Lexer) Status() source.Status {
	if lex.line < 1 || lex.line > len(lex.source) {
		return source.BadLineNumber
	}
	if lex.char < 1 || lex.char - 1 > len(lex.source[lex.line-1]) {
		return source.BadCharNumber
	}
	if lex.char + 1 == len(lex.source[lex.line-1]) {
		if lex.line == len(lex.source) {
			return source.Eof
		}
		return source.Eol
	}
	return source.Ok
}

func (lex *Lexer) stepDirection(n int) source.Status {
	if stat := lex.Status(); stat.IsInvalid() {
		return stat
	}

	char := lex.char + n
	line := lex.line
	if char <= 0 {
		for char <= 0 {
			if line-1 <= 0 {
				return source.BadLineNumber
			}
			line = line - 1
			underflow := -char - 1
			char = len(lex.source[line-1]) - underflow
		}
	} else if char > len(lex.source[line-1]) + 1 {
		for char > len(lex.source[line-1]) + 1 {
			if line+1 > len(lex.source) {
				return source.BadLineNumber
			}
			char = (char-1) - len(lex.source[line-1])
			line = line + 1
		}
	}
	lex.line, lex.char = line, char
	return source.Ok
}

func (lex *Lexer) resetCurr() string {
	out := string(lex.curr.data[:lex.curr.len])
	lex.curr.len = 0
	return out
}

type Lexer struct {
	whitespace *regexp.Regexp
	errors 	   []error
	path       string
	line       int
	char       int
	source     []string
	tokens     []token.Token
	curr       struct {
		data []byte
		cap  int
		len  int
	}
}

func (lex *Lexer) GetTokens() []token.Token {
	return lex.tokens
}

func (lex *Lexer) HasErrors() bool {
	return len(lex.errors) != 0
}

// Number of lines in source
func (lex *Lexer) NumLines() int {
	return len(lex.source)
}
