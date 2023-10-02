package lexer

import (
	"regexp"
	"testing"

	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-packages/util/testutil"
)

var lexerWhitespace_test = regexp.MustCompile(` \t`)

func TestStepDirection(t *testing.T) {
	type before struct{line, char int}
	type after before
	tests := []struct{
		source []string
		stat source.Status
		step int
		before
		after
	}{
		// edge cases
		{
			[]string{``}, source.Ok, 
			0,
			before{1,1}, after{1,1},
		},
		{
			[]string{``}, source.BadLineNumber, 
			1,
			before{1,1}, after{1,1},
		},
		{
			[]string{``,``}, source.Ok, 
			1,
			before{1,1}, after{2,1},
		},
		{
			[]string{``,``}, source.Ok, 
			-1,
			before{2,1}, after{1,1},
		},
		{
			[]string{``,``,``,``,``}, source.Ok, 
			4,
			before{1,1}, after{5,1},
		},
		// base cases
		{
			[]string{`.a`}, source.Ok, 
			0,
			before{1,1}, after{1,1},
		},
		{
			[]string{`.a`}, source.Ok, 
			1,
			before{1,1}, after{1,2},
		},
		{
			[]string{`a`, `b`}, source.Ok, 
			1,
			before{1,2}, after{2,1},
		},
		{
			[]string{`.a`}, source.Ok, 
			-1,
			before{1,2}, after{1,1},
		},
		{
			[]string{`a`, `b`}, source.Ok, 
			-1,
			before{2,1}, after{1,2},
		},
		// bad cases
		{
			[]string{`.a`}, source.BadLineNumber, 
			0,
			before{0,0}, after{0,0},
		},
		{
			[]string{`.a`}, source.BadLineNumber, 
			1,
			before{-100,135}, after{-100,135},
		},
		// forwards
		{
			[]string{`.a`}, source.Ok, 
			1,
			before{1,1}, after{1,2},
		},
		{
			[]string{`.a`, `.b`}, source.Ok, 
			2,
			before{1,1}, after{1,3},
		},
		{
			[]string{`.a`, `.b`}, source.Ok, 
			3,
			before{1,1}, after{2,1},
		},
		{
			[]string{`.a`, `.b`}, source.Ok, 
			4,
			before{1,1}, after{2,2},
		},
		{
			[]string{`.a`, `.b`}, source.Ok, 
			5,
			before{1,1}, after{2,3},
		},
		{
			[]string{`.a`, `.b`}, source.BadLineNumber, 
			6,
			before{1,1}, after{1,1},
		},
		// backwards
		{
			[]string{`.a`}, source.Ok, 
			-1,
			before{1,2}, after{1,1},
		},
		{
			[]string{`.a`, `.b`}, source.Ok, 
			-2,
			before{1,3}, after{1,1},
		},
		{
			[]string{`.a`, `.b`}, source.Ok, 
			-3,
			before{2,1}, after{1,1},
		},
		{
			[]string{`.a`, `.b`}, source.Ok, 
			-4,
			before{2,2}, after{1,1},
		},
		{
			[]string{`.a`, `.b`}, source.Ok, 
			-5,
			before{2,3}, after{1,1},
		},
		{
			[]string{`.a`, `.b`}, source.BadLineNumber, 
			-6,
			before{2,3}, after{2,3},
		},
	}

	for i, test := range tests {
		lex := NewLexer(lexerWhitespace_test)
		// setup
		lex.SetSource(test.source)
		lex.SetPath("./test-lexer-unadvance-invalid.yew")
		lex.SetLineChar(test.before.line, test.before.char)

		stat := lex.stepDirection(test.step)
		if !stat.Is(test.stat) {
			t.Fatal(testutil.TestMsg(i, "expected stat=%v,  actual stat=%v", test.stat, stat))
		}

		line, char := lex.GetLineChar()
		if test.after.line != line {
			t.Fatal(testutil.TestFail2("line", test.after.line, line, i))
		}
		if test.after.char != char {
			t.Fatal(testutil.TestFail2("char", test.after.char, char, i))
		}
	}
}

func TestAdvanceChar(t *testing.T) {
	tests := []struct{
		line, char int
		expect byte
		stat source.Status
	}{
		{1,2,' ',source.Ok}, {1,3,' ',source.Ok}, {1,4,' ',source.Ok}, {2,1,'\n',source.Ok},
		{2,2,' ',source.Ok}, {2,3,' ',source.Ok}, {2,4,' ',source.Ok}, {2,4,0,source.Eof},
	}
	lex := NewLexer(lexerWhitespace_test, 0, 0, 1)
	lex.SetSource([]string{`   `,`   `})
	lex.SetPath("./test-lex-advance-char-position.yew")

	for i, test := range tests {
		actual, stat := lex.AdvanceChar()
		if !stat.Is(test.stat) {
			t.Fatal(testutil.TestFail2("stat", test.stat, stat, i))
		}
		if actual != test.expect {
			t.Fatalf(testutil.TestFail2("byte", test.expect, actual, i))
		}
		
		line, char := lex.GetLineChar()
		if test.line != line {
			t.Fatal(testutil.TestFail2("line", test.line, line, i))
		}
		if test.char != char {
			t.Fatalf(testutil.TestFail2("char", test.char, char, i))
		}
	}
}

func TestAdvanceLine(t *testing.T) {
	{
		tests := []struct{
			line, char int
			stat source.Status
		}{
			{2,1,source.Ok}, {3,1,source.Ok}, {3,1,source.Eof},
		}
		lex := NewLexer(lexerWhitespace_test, 0, 0, 1)
		lex.SetSource([]string{`   `,`   `,`   `})
		lex.SetPath("./test-lex-advance-line-position.yew")

		for i, test := range tests {
			stat := lex.AdvanceLine()
			if !stat.Is(test.stat) {
				t.Fatal(testutil.TestFail2("stat", test.stat, stat, i))
			}
			
			line, char := lex.GetLineChar()
			if test.line != line {
				t.Fatal(testutil.TestFail2("line", test.line, line, i))
			}
			if test.char != char {
				t.Fatalf(testutil.TestFail2("char", test.char, char, i))
			}
		}
	}

	// test that (*lexer.Lexer) AdvanceLine() resets char to 1
	{
		lex := NewLexer(lexerWhitespace_test, 0, 0, 1)
		lex.SetSource([]string{`   `,`   `})
		lex.SetPath("./test-lex-advance-line-position.yew")
		
		_, _ = lex.AdvanceChar()
		stat := lex.AdvanceLine()
		if !stat.Is(source.Ok) {
			t.Fatal(testutil.TestFail2("stat", source.Ok, stat, 0, 1))
		}
		
		line, char := lex.GetLineChar()
		if line != 2 {
			t.Fatal(testutil.TestFail2("line", 2, line, 0, 1))
		}
		if char != 1 {
			t.Fatalf(testutil.TestFail2("char", 1, char, 0, 1))
		}
	}
}