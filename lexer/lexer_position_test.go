package lexer

import (
	"regexp"
	"testing"

	"github.com/petersalex27/yew-packages/source"
	testutil "github.com/petersalex27/yew-packages/util/test"
)

var lexerWhitespace_test = regexp.MustCompile(` \t`)

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