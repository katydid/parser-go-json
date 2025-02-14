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

// Package json contains the implementation of a JSON parser.
package json

import (
	"fmt"
	"testing"
)

func TestParseString(t *testing.T) {
	str := `"a"`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	if !p.IsLeaf() {
		t.Fatalf("expected leaf")
	}
	expect(t, p.String, "a")
	expectEOF(t, p.Next)
}

func TestParseInt(t *testing.T) {
	str := `-1`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	if !p.IsLeaf() {
		t.Fatalf("expected leaf")
	}
	expect(t, p.Int, -1)
	expectEOF(t, p.Next)
}

func TestParseDouble(t *testing.T) {
	str := `1.1`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	if !p.IsLeaf() {
		t.Fatalf("expected leaf")
	}
	expect(t, p.Double, 1.1)
	expectEOF(t, p.Next)
}

func TestParseObjectKeys(t *testing.T) {
	str := `{"a":1, "b":{"c":3}, "d":4}`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	expect(t, p.String, "a")
	assertNoErr(t, p.Next)
	expect(t, p.String, "b")
	assertNoErr(t, p.Next)
	expect(t, p.String, "d")
	expectEOF(t, p.Next)
}

func TestParseObjectKeysWithArray(t *testing.T) {
	str := `{"a":1, "b":{"c":2}, "d":[3,4], "e":"f"}`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	expect(t, p.String, "a")
	assertNoErr(t, p.Next)
	expect(t, p.String, "b")
	assertNoErr(t, p.Next)
	expect(t, p.String, "d")
	assertNoErr(t, p.Next)
	expect(t, p.String, "e")
	expectEOF(t, p.Next)
}

func TestParseObjectValues(t *testing.T) {
	str := `{"a":1, "b":{"c":3}, "d":4}`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	expect(t, p.String, "a")

	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 1)
	expectEOF(t, p.Next)
	p.Up()

	assertNoErr(t, p.Next)
	expect(t, p.String, "b")

	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.String, "c")
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 3)
	expectEOF(t, p.Next)
	p.Up()
	expectEOF(t, p.Next)
	p.Up()

	assertNoErr(t, p.Next)
	expect(t, p.String, "d")
	expectEOF(t, p.Next)
}

func TestParseArrayIndexes(t *testing.T) {
	str := `["a", true, [1,2], {"a":1}]`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	expect(t, p.Int, 0)
	assertNoErr(t, p.Next)
	expect(t, p.Int, 1)
	assertNoErr(t, p.Next)
	expect(t, p.Int, 2)
	assertNoErr(t, p.Next)
	expect(t, p.Int, 3)
	expectEOF(t, p.Next)
}

func TestParseArrayElements(t *testing.T) {
	str := `["a", true, {"a":97}, [10,20]]`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("first element\n")
	assertNoErr(t, p.Next)
	expect(t, p.Int, 0)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.String, "a")
	expectEOF(t, p.Next)
	p.Up()

	fmt.Printf("second element\n")
	assertNoErr(t, p.Next)
	expect(t, p.Int, 1)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Bool, true)
	expectEOF(t, p.Next)
	p.Up()

	fmt.Printf("third element\n")
	assertNoErr(t, p.Next)
	expect(t, p.Int, 2)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.String, "a")
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 97)
	expectEOF(t, p.Next)
	p.Up()
	expectEOF(t, p.Next)
	p.Up()

	fmt.Printf("fourth element\n")
	assertNoErr(t, p.Next)
	expect(t, p.Int, 3)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 0)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 10)
	expectEOF(t, p.Next)
	p.Up()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 1)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 20)
	expectEOF(t, p.Next)
	p.Up()
	expectEOF(t, p.Next)
	p.Up()

	expectEOF(t, p.Next)
}
