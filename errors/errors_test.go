package errors

import (
	"testing"
)
// Ferr("tplcms", "Subtype", path, line, char, msg, src)
func TestGetTail(t *testing.T) {
	tests := []struct{
		in Err
		expect string
	}{
		{
			Ferr("flcms", "test", 1, 1, "msg", "source"),
			"[test:1:1] Error: msg\n" +
			"  1 | source\n" +
			"      ^",
		},
		{
			Ferr("tflcms", "Syntax", "test", 1, 9, "unexpected token", "Maybe a @ Just a | Nothing"),
			"[test:1:9] Error (Syntax): unexpected token\n" +
			"  1 | Maybe a @ Just a | Nothing\n" +
			"              ^",
		},
		{
			Ferr("flcmsr", "test", 1, 1, "msg", "source", 6),
			"[test:1:1] Error: msg\n" +
			"  1 | source\n" +
			"      ^^^^^^",
		},
		{
			Ferr("flcmsrp", "test", 1, 1, "msg", "source", 6, " here"),
			"[test:1:1] Error: msg\n" +
			"  1 | source\n" +
			"      ^^^^^^ here",
		},
		{
			Ferr("tflcmsr", "Syntax", "test", 1, 11, "expected type identifier", "Maybe a = just a | nothing", 4),
			"[test:1:11] Error (Syntax): expected type identifier\n" +
			"  1 | Maybe a = just a | nothing\n" +
			"                ^^^^",
		},
		{
			Ferr("tflcmsrp", "Syntax", "test", 1, 11, "expected type identifier", "Maybe a = just a | nothing", 4, "(11-14)"),
			"[test:1:11] Error (Syntax): expected type identifier\n" +
			"  1 | Maybe a = just a | nothing\n" +
			"                ^^^^(11-14)",
		},
		{
			Ferr(""),
			"Error: unknown error",
		},
	}

	for ti, test := range tests {
		actual := test.in.Error()
		if actual != test.expect {
			t.Fatalf("failed test #%d:\nexpected:\n%s\nactual:\n%s\n", ti+1, test.expect, actual)
		}
	}
}