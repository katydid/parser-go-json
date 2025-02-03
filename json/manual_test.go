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
	"io"
	"testing"

	"github.com/katydid/parser-go/parser/debug"
)

func parse(s string) (debug.Nodes, error) {
	parser := NewParser()
	if err := parser.Init([]byte(s)); err != nil {
		return nil, err
	}
	return debug.Parse(parser)
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
	m, err := parse(string(data))
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
	m, err := parse(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", m)
}

func TestIntWithExponent(t *testing.T) {
	s := `{"A":1e+08}`
	parser := NewParser()
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
	parser := NewParser()
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

func TestIndexedArray(t *testing.T) {
	s := `["a", "b", "c"]`
	parser := NewParser()
	if err := parser.Init([]byte(s)); err != nil {
		t.Fatal(err)
	}
	if err := parser.Next(); err != nil {
		t.Fatal(err)
	}

	if parser.IsLeaf() {
		t.Fatal("expected index not leaf")
	}
	index, err := parser.Int()
	if err != nil {
		t.Fatal(err)
	}
	if index != 0 {
		t.Fatalf("expected index = 0, but got %d", index)
	}
	parser.Down()
	if err := parser.Next(); err != nil {
		t.Fatal(err)
	}
	if !parser.IsLeaf() {
		t.Fatal("expected leaf")
	}
	if s, err := parser.String(); err != nil {
		t.Fatal(err)
	} else if s != "a" {
		t.Fatalf("want a, but got %s", s)
	}
	parser.Up()
	if err := parser.Next(); err != nil {
		t.Fatal(err)
	}

	if parser.IsLeaf() {
		t.Fatal("expected index not leaf")
	}
	index, err = parser.Int()
	if err != nil {
		t.Fatal(err)
	}
	if index != 1 {
		t.Fatalf("expected index = 1, but got %d", index)
	}
	parser.Down()
	if err := parser.Next(); err != nil {
		t.Fatal(err)
	}
	if !parser.IsLeaf() {
		t.Fatal("expected leaf")
	}
	if s, err := parser.String(); err != nil {
		t.Fatal(err)
	} else if s != "b" {
		t.Fatalf("want b, but got %s", s)
	}
	parser.Up()
	if err := parser.Next(); err != nil {
		t.Fatal(err)
	}

	if parser.IsLeaf() {
		t.Fatal("expected index not leaf")
	}
	index, err = parser.Int()
	if err != nil {
		t.Fatal(err)
	}
	if index != 2 {
		t.Fatalf("expected index = 2, but got %d", index)
	}
	parser.Down()
	if err := parser.Next(); err != nil {
		t.Fatal(err)
	}
	if !parser.IsLeaf() {
		t.Fatal("expected leaf")
	}
	if s, err := parser.String(); err != nil {
		t.Fatal(err)
	} else if s != "c" {
		t.Fatalf("want c, but got %s", s)
	}
	parser.Up()
	if err := parser.Next(); err != io.EOF {
		t.Fatalf("expected EOF got err = %v", err)
	}
}
