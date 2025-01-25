//  Copyright 2015 Walter Schulze
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

package json_test

import (
	"encoding/json"
	"testing"

	sjson "github.com/katydid/parser-go-json/json"
	"github.com/katydid/parser-go/parser/debug"
)

func TestDebug(t *testing.T) {
	p := sjson.NewJsonParser()
	data, err := json.Marshal(debug.Input)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	m := debug.Walk(p)
	if !m.Equal(debug.Output) {
		t.Fatalf("expected %s but got %s", debug.Output, m)
	}
}

func TestRandomDebug(t *testing.T) {
	p := sjson.NewJsonParser()
	data, err := json.Marshal(debug.Input)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if err := p.Init(data); err != nil {
			t.Fatal(err)
		}
		//l := debug.NewLogger(p, debug.NewLineLogger())
		debug.RandomWalk(p, debug.NewRand(), 10, 3)
		//t.Logf("original %v vs random %v", debug.Output, m)
	}
}

func TestEscapedChar(t *testing.T) {
	j := map[string][]interface{}{
		`a\"`: {1},
	}
	data, err := json.Marshal(j)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", string(data))
	parser := sjson.NewJsonParser()
	if err := parser.Init(data); err != nil {
		t.Fatal(err)
	}
	m, err := walk(parser)
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
	parser := sjson.NewJsonParser()
	if err := parser.Init([]byte(s)); err != nil {
		t.Fatal(err)
	}
	m, err := walk(parser)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", m)
}

func TestIntWithExponent(t *testing.T) {
	s := `{"A":1e+08}`
	parser := sjson.NewJsonParser()
	if err := parser.Init([]byte(s)); err != nil {
		t.Fatal(err)
	}
	if err := parser.Next(); err != nil {
		t.Fatal(err)
	}
	parser.Down()
	if err := parser.Next(); err != nil {
		t.Fatal(err)
	}
	if !parser.IsLeaf() {
		t.Fatal("incorrect walk, please adjust the path above")
	}
	if i, err := parser.Int(); err != nil {
		t.Fatalf("did not expect error %v", err)
	} else if i != 1e+08 {
		t.Fatalf("got %d", i)
	}
}

func TestTooLargeNumber(t *testing.T) {
	input := `123456789.123456789e+123456789`
	parser := sjson.NewJsonParser()
	if err := parser.Init([]byte(input)); err != nil {
		t.Fatalf("init error: %v", err)
	}
	if err := parser.Next(); err != nil {
		t.Fatalf("Next err = %v", err)
	}
	if _, err := parser.Double(); err == nil {
		t.Fatal("expected number to be too large")
	}
	bs, err := parser.Bytes()
	if err != nil {
		t.Fatalf("expected bytes to return anyway, but got error = %v", err)
	}
	if string(bs) != input {
		t.Fatalf("expected %v, but got %v", input, string(bs))
	}
}

func TestValues(t *testing.T) {
	testValue(t, "0", "0")
	testValue(t, "1", "1")
	testValue(t, "-1", "-1")
	testValue(t, "123", "123")
	testValue(t, "1.1", "1.1")
	testValue(t, "1.123", "1.123")
	testValue(t, "1.1e1", "11")
	testValue(t, "1.1e-1", "0.11")
	testValue(t, "1.1e10", "11000000000")
	testValue(t, "1.1e+10", "11000000000")
	testValue(t, `"a"`, "a")
	testValue(t, `"abc"`, "abc")
	testValue(t, `""`, "")
	testValue(t, `"\b"`, "\b")
	testValue(t, `true`, "true")
	testValue(t, `false`, "false")
	testValue(t, `null`, "null")
}

func testWalk(t *testing.T, s string) {
	t.Helper()
	parser := sjson.NewJsonParser()
	if err := parser.Init([]byte(s)); err != nil {
		t.Error(err)
		return
	}
	m, err := walk(parser)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", m)
}

func TestArray(t *testing.T) {
	testWalk(t, `[1]`)
	testWalk(t, `[1,2.3e5]`)
	testWalk(t, `[1,"a"]`)
	testWalk(t, `[true, false, null]`)
	testWalk(t, `[{"a": true, "b": [1,2]}]`)
}
