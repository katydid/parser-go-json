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

	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestSkipObjectValues(t *testing.T) {
	str := `{"a":1, "b":{"c":3}, "d":4}`
	p := NewParser()
	p.Init([]byte(str))

	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")

	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "c")
	expect.Hint(t, p, parse.ValueHint)
	// do not tokenize 3
	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "d")
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipArrayElements(t *testing.T) {
	str := `["a", true, {"a":97}, [10,20]]`
	p := NewParser()
	p.Init([]byte(str))

	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "a")

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ValueHint)
	expect.True(t, p)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 2)
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 3)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.NoErr(t, p.Skip)
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipArrayObjectElement(t *testing.T) {
	str := `{"A":[{"a":97}],"B":2}`
	p := NewParser()
	p.Init([]byte(str))
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.EnterHint)
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
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
	p.Init([]byte(str))
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")

	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
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
	p.Init([]byte(str))
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")

	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
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
	p.Init([]byte(str))
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.NoErr(t, p.Skip)
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipUpUpUpTwoFields(t *testing.T) {
	str := `{
		"A": [
			{"a": "b", "c": "d"},
			"b",
			"c"
		],
		"B": 1
	}`
	p := NewParser()
	p.Init([]byte(str))
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "b")

	expect.NoErr(t, p.Skip) // skip over , `"c": "d"` and `}`,
	expect.NoErr(t, p.Skip) // skip over `1:b, 2:c` and `]`

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipUpUpUpEndOfObject(t *testing.T) {
	str := `{
		"A": [
			{"a": "b"},
			"b",
			"c"
		],
		"B": 1
	}`
	p := NewParser()
	p.Init([]byte(str))
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "b")

	expect.NoErr(t, p.Skip) // skip over `}`,
	expect.NoErr(t, p.Skip) // skip over `1:b, 2:c` and `]`

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipUpUpUpOnlyObjects(t *testing.T) {
	str := `{
		"A": {
			"0": {"a": "b", "c": "d"},
			"1": "b",
			"2": "c"
		},
		"B": 1
	}`
	p := NewParser()
	p.Init([]byte(str))
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "A")
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "0")
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "b")

	expect.NoErr(t, p.Skip) // skip over , `"c": "d"` and `}`,
	expect.NoErr(t, p.Skip) // skip over `"1":"b", "2":"c"` and `}`

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "B")
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}
