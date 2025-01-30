package json

import (
	"bytes"
	"testing"
)

type unQuoteTest struct {
	in  string
	out string
}

var unquotetests = []unQuoteTest{
	{`""`, ""},
	{`"a"`, "a"},
	{`"abc"`, "abc"},
	{`"☺"`, "☺"},
	{`"hello world"`, "hello world"},
	{`"\u1234"`, "\u1234"},
	{`"'"`, "'"},
}

func TestUnquote(t *testing.T) {
	for _, tt := range unquotetests {
		testUnquote(t, tt.in, tt.out)
	}
}

func testUnquote(t *testing.T, in, want string) {
	t.Helper()
	// Test Unquote.
	got, gotOk := unquoteBytes([]byte(in))
	if !bytes.Equal(got, []byte(want)) || !gotOk {
		t.Errorf("Unquote(%q) = (%q, %v), want (%q, %v)", in, got, gotOk, want, true)
	}
}
