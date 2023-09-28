package lexer

import (
	"testing"

	"github.com/petersalex27/yew-packages/source"
	testutil "github.com/petersalex27/yew-packages/util/test"
)

//SetLineChar(line, char int)
//AdvanceChar() (char byte, stat Status)
//UnadvanceChar() (stat Status)
//WhitespaceLength() int

//ReadUntil(char byte) (string, Status)
func TestReadUntil(t *testing.T) {
	tests := []struct{
		source []string
		escapeMap source.EscapeMap
		escape byte
		until byte
		expect string
		line, char int
	}{
		{[]string{`.`}, nil, 0, '.', ``, 1, 2},
		{[]string{`.`}, nil, '\\', '.', ``, 1, 2},
		{[]string{` "`}, nil, '\\', '"', ` `, 1, 3},
		{
			[]string{`this is a string"`}, 
			nil, 
			'\\', '"', `this is a string`,
			1, 18,
		},
		{
			[]string{`this is a string"`,`another line :)`}, 
			nil, 
			'\\', '"', `this is a string`,
			1, 18,
		},
		{
			[]string{`this is \na new line"`}, 
			source.EscapeMap_s{
				LookupSize: 1,
				Map: map[string]string{`n`:"\n"},
			}, 
			'\\', '"', "this is \na new line",
			1, 22,
		},
		{
			[]string{`this is %san s escape"`}, 
			source.EscapeMap_s{
				LookupSize: 1, 
				Map: map[string]string{`s`:"_s_"},}, 
			'%', '"', "this is _s_an s escape",
			1, 23,
		},
		{
			[]string{`this !b not a string"`}, 
			source.EscapeMap_s{
				LookupSize: 1,
				Map: map[string]string{`s`:"_s_", `b`:"is",},
			}, 
			'!', '"', 
			"this is not a string",
			1, 22,
		},
		{
			[]string{ `x68x65x6cx6cx6f, worx6cd!.` }, 
			source.EscapeMap_s{
				LookupSize: 2,
				Map: map[string]string{
					`65`:`e`, `66`:`f`, `67`:`g`, 
					`68`:`h`, `6c`:`l`, `6f`:`o`,
				},
			}, 
			'x', '.', 
			"hello, world!",
			1, 27,
		},
		{
			[]string{ `x68x65x6cx6cx6f, worx6cd!. Hi!.` }, 
			source.EscapeMap_s{
				LookupSize: 2,
				Map: map[string]string{
					`65`:`e`, `66`:`f`, `67`:`g`, 
					`68`:`h`, `6c`:`l`, `6f`:`o`,
				},
			}, 
			'x', '.', 
			"hello, world!",
			1, 27,
		},
	}

	for i, test := range tests {
		lex := NewLexer(lexerWhitespace_test)
		lex.SetSource(test.source)
		lex.SetPath("./test-lexer-read.yew")

		actual, stat := source.ReadUntil(test.escapeMap, lex, test.until, test.escape)
		if stat.NotOk() {
			t.Fatal(testutil.TestMsg(i, "expected stat=%v,  actual stat=%v", source.Ok, stat))
		}
		if test.expect != actual {
			t.Fatal(testutil.TestFail2("read", test.expect, actual, i))
		}

		line, char := lex.GetLineChar()
		if test.line != line {
			t.Fatal(testutil.TestFail2("line", test.line, line, i))
		}
		if test.char != char {
			t.Fatal(testutil.TestFail2("char", test.char, char, i))
		}
	}
}