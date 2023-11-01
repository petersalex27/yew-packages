package lexer

import (
	"strings"
)

// formats comments for documentation
//
// it's very basic right now, but will expand later
// TODO: allow annotations to define how this works and create an interface
// that supports that kind of interfacing
//
// some possible fields:
//
//	isMultiLine bool
//	regexDictionary map[string]*regexp.Regexp
//	previousLines []string
type Documentor struct {
	// takes in a single line of a (multi- or single-line) comment and formats it
	// returning the formated line
	formatCommentLine func(this *Documentor, comment string) (newComment string)
}

func (documentor *Documentor) Run(comment string) (newComment string) {
	return documentor.formatCommentLine(documentor, comment)
}

// removes all trailing whitespace (or more prec., ' ' and '\t') from `comment`
// and returns it
//
// default value for
//
//	Documentor.formatCommentLine
func removeTrailingWhitespace(_ *Documentor, comment string) (_ string) {
	return strings.TrimRight(comment, " \t")
}

// MakeDocumentor creates a new documentor and assigns `formatCommentLine` to
// its field of the same name when non-nil. If `formatCommentLine` is nil, then
// the default function`removeTrailingWhitespace` is used instead
//
// SEE: removeTrailingWhitespace(*Documentor, string) string
func MakeDocumentor(formatCommentLine func(this *Documentor, comment string) (newComment string)) (this *Documentor) {
	this = new(Documentor)
	if formatCommentLine == nil {
		formatCommentLine = removeTrailingWhitespace
	}
	this.formatCommentLine = formatCommentLine
	return this
}
