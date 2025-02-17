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
	"testing"
)

func TestSkipObjectValues(t *testing.T) {
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
	// expectEOF(t, p.Next)
	p.Up()

	assertNoErr(t, p.Next)
	expect(t, p.String, "b")

	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.String, "c")
	p.Down()
	assertNoErr(t, p.Next)
	// expect(t, p.Int, 3)
	// expectEOF(t, p.Next)
	p.Up()
	expectEOF(t, p.Next)
	p.Up()

	assertNoErr(t, p.Next)
	expect(t, p.String, "d")
	expectEOF(t, p.Next)
}

func TestSkipArrayElements(t *testing.T) {
	str := `["a", true, {"a":97}, [10,20]]`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}

	assertNoErr(t, p.Next)
	expect(t, p.Int, 0)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.String, "a")
	// expectEOF(t, p.Next)
	p.Up()

	assertNoErr(t, p.Next)
	expect(t, p.Int, 1)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Bool, true)
	expectEOF(t, p.Next)
	p.Up()

	assertNoErr(t, p.Next)
	expect(t, p.Int, 2)
	p.Down()
	assertNoErr(t, p.Next)
	// expect(t, p.String, "a")
	// p.Down()
	// assertNoErr(t, p.Next)
	// expect(t, p.Int, 97)
	// expectEOF(t, p.Next)
	// p.Up()
	// expectEOF(t, p.Next)
	p.Up()

	assertNoErr(t, p.Next)
	expect(t, p.Int, 3)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 0)
	p.Down()
	assertNoErr(t, p.Next)
	// expect(t, p.Int, 10)
	// expectEOF(t, p.Next)
	p.Up()
	assertNoErr(t, p.Next)
	// expect(t, p.Int, 1)
	// p.Down()
	// assertNoErr(t, p.Next)
	// expect(t, p.Int, 20)
	// expectEOF(t, p.Next)
	// p.Up()
	// expectEOF(t, p.Next)
	p.Up()

	expectEOF(t, p.Next)
}

func TestSkipArrayObjectElement(t *testing.T) {
	str := `{"A":[{"a":97}],"B":2}`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	expect(t, p.String, "A")
	p.Down()
	assertNoErr(t, p.Next)
	p.Up()
	assertNoErr(t, p.Next)
	expect(t, p.String, "B")
	expectEOF(t, p.Next)
}

func TestSkipUpTwiceInARow(t *testing.T) {
	str := `{
		"A": [
			"a",
			"b",
			"c"
		]
	}`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	expect(t, p.String, "A")

	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 0)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.String, "a")
	p.Up()
	p.Up()
	expectEOF(t, p.Next)
}

func TestSkipDownUp(t *testing.T) {
	str := `{
		"A": [
			"a",
			"b",
			"c"
		]
	}`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	expect(t, p.String, "A")

	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 0)
	p.Down()
	p.Up()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 1)
	p.Up()
	expectEOF(t, p.Next)
}

func TestSkipDownUpUp(t *testing.T) {
	str := `{
		"A": [
			"a",
			"b",
			"c"
		],
		"B": 1
	}`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	expect(t, p.String, "A")

	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 0)
	p.Down()
	p.Up()
	p.Up()
	assertNoErr(t, p.Next)
	expect(t, p.String, "B")
	expectEOF(t, p.Next)
}

func TestSkipUpUpUp(t *testing.T) {
	str := `{
		"A": [
			{"a": "b", "c": "d"},
			"b",
			"c"
		],
		"B": 1
	}`
	p := NewParser()
	if err := p.Init([]byte(str)); err != nil {
		t.Fatal(err)
	}
	assertNoErr(t, p.Next)
	expect(t, p.String, "A")

	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.Int, 0)
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.String, "a")
	p.Down()
	assertNoErr(t, p.Next)
	expect(t, p.String, "b")
	p.Up()
	p.Up()
	p.Up()
	assertNoErr(t, p.Next)
	expect(t, p.String, "B")
	expectEOF(t, p.Next)
}
