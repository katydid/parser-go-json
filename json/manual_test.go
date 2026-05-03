//  Copyright 2025 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package json

import (
	"encoding/json"
	"testing"

	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/parser-go/parse/debug"
)

func parseString(s string) (debug.Nodes, error) {
	parser := NewParser()
	parser.Init([]byte(s))
	return debug.Parse(parser)
}

func TestEscapedChar(t *testing.T) {
	j := map[string][]any{
		`a\"`: {1},
	}
	data, err := json.Marshal(j)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", string(data))
	m, err := parseString(string(data))
	if err != nil {
		t.Fatal(err)
	}
	name := m[0].Label
	if name != `a\"` {
		t.Fatalf("wrong escaped name %s", name)
	}
}

func TestMultiLineArray(t *testing.T) {
	s := `{
		"A":[1]
	}`
	m, err := parseString(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", m)
}

func TestIntWithExponent(t *testing.T) {
	s := `{"A":1e+08}`
	p := NewParser()
	p.Init([]byte(s))
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.ValueHint)
	expect.Float(t, p, 1e+08)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestTooLargeNumber(t *testing.T) {
	want := `123456789.123456789e+123456789`
	p := NewParser()
	p.Init([]byte(want))
	expect.Hint(t, p, parse.ValueHint)
	kind, val, err := p.Token()
	if err != nil {
		t.Fatal(err)
	}
	if kind != parse.DecimalKind {
		t.Fatalf("want DecimalKind, but got %v", kind)
	}
	got := cast.ToString(val)
	if want != got {
		t.Fatalf("want %s got %s", want, got)
	}
	expect.EOF(t, p)
}

func TestIndexedArray(t *testing.T) {
	s := `["a", "b", "c"]`
	p := NewParser()
	p.Init([]byte(s))
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "a")

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "b")

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 2)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "c")

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}
