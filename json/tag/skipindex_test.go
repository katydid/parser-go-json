//  Copyright 2026 Walter Schulze
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

package tag_test

import (
	"io"
	"testing"

	jsonparse "github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go-json/json/tag"
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestSkipIndexOnlyUnknownObjectOpen(t *testing.T) {
	str := `{}` // {}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipIndexOnlyUnknownObjectAfterOpen(t *testing.T) {
	str := `{}` // {}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.NoErr(t, p.Skip) // skip object close
	expect.EOF(t, p)
}

func TestSkipIndexOnlyUnknownArrayOpen(t *testing.T) {
	str := `[]` // []
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipIndexOnlyUnknownArrayAfterOpen(t *testing.T) {
	str := `[]` // []
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.NoErr(t, p.Skip)
	expect.EOF(t, p)
}

func TestSkipIndexOnlySingletonArrayAfterOpen(t *testing.T) {
	str := `[1]` // [0:1]
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.NoErr(t, p.Skip)
	expect.EOF(t, p)
}

func TestSkipIndexOnlyUnknownString(t *testing.T) {
	str := `"abc"` // "abc"
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.NoErr(t, p.Skip)
	expect.EOF(t, p)
}

// If the kind '[' was returned by Next, then the whole array is skipped.
func TestSkipIndexOnlyArrayOpen(t *testing.T) {
	str := `[1,2]` // [0:1, 1:2]
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.NoErr(t, p.Skip)
	// skipped over 0:1,1:2]
	expect.EOF(t, p)
}

func TestSkipIndexOnlyArrayFirst(t *testing.T) {
	str := `[1,2]` // [0:1, 1:2]
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.NoErr(t, p.Skip)
	// skipped over 1
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 2)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipIndexOnlyArrayNestedOpen(t *testing.T) {
	str := `[[1,2]]` // [0:[0:1, 1:2]]
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)

	expect.Hint(t, p, parse.EnterHint)

	expect.NoErr(t, p.Skip)
	// skipped over 0:1,1:2]
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

// If an array element was parsed, then the rest of the array is skipped.
func TestSkipIndexOnlyArrayElement(t *testing.T) {
	str := `[1,2,3]` // [0:1,1:2,2:3]
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over 1:2,2:3]
	expect.EOF(t, p)
}

func TestSkipIndexOnlyArrayNestedElement(t *testing.T) {
	str := `[1,[2,3,4],5]` // [0:1, 1:[0:2,1:3,2:4], 1:5]
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 2)

	expect.NoErr(t, p.Skip)
	// skipped over 1:3,2:4]

	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 2)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 5)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipIndexOnlyArrayRecursiveElement(t *testing.T) {
	str := `[1,[2,3],[[4,5,6]]]` // [0:1, 1:[0:2,1:3], 2:[0:[0:4,1:5,2:6]]]
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over 1:[0:2,1:3], 2:[0:[0:4,1:5,2:6]]]
	expect.EOF(t, p)
}

// If the kind '{' was returned by Next, then the whole object is skipped.
func TestSkipIndexOnlyObjectOpen(t *testing.T) {
	str := `{"a":1,"b":2}` // {"a": 1, "b": 2}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.NoErr(t, p.Skip)
	// skipped over "a":1,"b":2}
	expect.EOF(t, p)
}

func TestSkipIndexOnlyObjectNestedOpen(t *testing.T) {
	str := `{"a":{"b":1,"c":2}}` // {"a":{"b":1,"c":2}}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.EnterHint)
	expect.NoErr(t, p.Skip)
	// skipped over "b":1,"c":2}
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

// If an object value was just parsed, then the rest of the object is skipped.
func TestSkipIndexOnlyObjectKey(t *testing.T) {
	str := `{"a":1,"b":2,"c":3}` // {"a":1,"b":2,"c":3}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over "b":2,"c":3}
	expect.EOF(t, p)
}

func TestSkipIndexOnlyObjectNestedKey(t *testing.T) {
	str := `{"a":{"b":1,"c":2,"d":3}}` // {"a": {"b":1,"c":2,"d":3}}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over "c":2,"d":3}
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

// If a object key was just parsed, then that key's value is skipped.
func TestSkipIndexOnlyObjectValue(t *testing.T) {
	str := `{"a":1,"b":2}` // {"a":1,"b":2}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)
	// skipped over 1
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")
	expect.NoErr(t, p.Skip)
	// skipped over 2
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipIndexOnlyObjectRecursiveValue(t *testing.T) {
	str := `{"a":1,"b":{"c":{"d":{"e":"f"},"g":[1,2]}}}`
	// {"a":1, "b": {"c": {"d": {"e": "f"}, "g": [0:1, 1:2]}}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")
	expect.NoErr(t, p.Skip)
	// skipped over {"c": {"d": {"e": "f"}, "g": [0:1, 1:2]}
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipIndexOnlyObjectDeepRecursiveValue(t *testing.T) {
	str := `{"a":1,"b":{"c":{"d":{"e":"f"},"g":[1,2]}}}`
	// {"a":1,"b":{"c":{"d":{"e":"f"},"g":[0:1,1:2]}}}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(str))), tag.WithIndexes())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "c")
	expect.NoErr(t, p.Skip)
	// skipped over {"d":{"e":"f"},"g":[1,2]}
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipIndexOnlyTagMixTwoObjectsWithIndexes(t *testing.T) {
	s := `[{"mykey1":true}, {"mykey2":false}]`
	// will be parsed the same as : [0: {"mykey1":true}, 1: {"mykey2":false}]
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(s))), tag.WithIndexes())

	// 1: first array
	expect.Hint(t, p, parse.EnterHint)

	// 2: first object
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}

	// 2: second object
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "mykey2")
	expect.Hint(t, p, parse.ValueHint)
	expect.False(t, p)
	expect.Hint(t, p, parse.LeaveHint)

	// 1: first array
	expect.Hint(t, p, parse.LeaveHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}
