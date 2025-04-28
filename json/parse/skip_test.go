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

package parse

import (
	"testing"

	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

// Skip returns an error if nothing has been parsed yet.
func TestSkipUnknownObject(t *testing.T) {
	str := `{}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.EOF(t, p)
}

func TestSkipUnknownArray(t *testing.T) {
	str := `[]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestSkipUnknownArrayOpen(t *testing.T) {
	str := `[1]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestSkipUnknownString(t *testing.T) {
	str := `"abc"`
	p := NewParser(WithBuffer([]byte(str)))
	expect.NoErr(t, p.Skip)
	expect.EOF(t, p)
}

// If the kind '[' was returned by Next, then the whole array is skipped.
func TestSkipArrayOpen(t *testing.T) {
	str := `[1,2]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.NoErr(t, p.Skip)
	// skipped over 1,2]
	expect.EOF(t, p)
}

func TestSkipArrayNestedOpen(t *testing.T) {
	str := `[[1,2]]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.NoErr(t, p.Skip)
	// skipped over 1,2]
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

// If an array element was parsed, then the rest of the array is skipped.
func TestSkipArrayElement(t *testing.T) {
	str := `[1,2,3]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over 2,3]
	expect.EOF(t, p)
}

func TestSkipArrayNestedElement(t *testing.T) {
	str := `[1,[2,3,4],5]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 2)
	expect.NoErr(t, p.Skip)
	// skipped over 3,4]
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 5)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestSkipArrayRecursiveElement(t *testing.T) {
	str := `[1,[2,3],[[4,5,6]]]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over [2,3],[[4,5,6]]]
	expect.EOF(t, p)
}

// If the kind '{' was returned by Next, then the whole object is skipped.
func TestSkipObjectOpen(t *testing.T) {
	str := `{"a":1,"b":2}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.NoErr(t, p.Skip)
	// skipped over "a":1,"b":2}
	expect.EOF(t, p)
}

func TestSkipObjectNestedOpen(t *testing.T) {
	str := `{"a":{"b":1,"c":2}}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.NoErr(t, p.Skip)
	// skipped over "b":1,"c":2}
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.EOF(t, p)
}

// If an object value was just parsed, then the rest of the object is skipped.
func TestSkipObjectKey(t *testing.T) {
	str := `{"a":1,"b":2,"c":3}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over "b":2,"c":3}
	expect.EOF(t, p)
}

func TestSkipObjectNestedKey(t *testing.T) {
	str := `{"a":{"b":1,"c":2,"d":3}}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over "c":2,"d":3}
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.EOF(t, p)
}

// If a object key was just parsed, then that key's value is skipped.
func TestSkipObjectValue(t *testing.T) {
	str := `{"a":1,"b":2}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")
	expect.NoErr(t, p.Skip)
	// skipped over 2
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.EOF(t, p)
}

func TestSkipObjectRecursiveValue(t *testing.T) {
	str := `{"a":1,"b":{"c":{"d":{"e":"f"},"g":[1,2]}}}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")
	expect.NoErr(t, p.Skip)
	// skipped over {"c":{"d":{"e":"f"},"g":[1,2]}}
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.EOF(t, p)
}

func TestSkipObjectDeepRecursiveValue(t *testing.T) {
	str := `{"a":1,"b":{"c":{"d":{"e":"f"},"g":[1,2]}}}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")

	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "c")
	expect.NoErr(t, p.Skip)
	// skipped over {"d":{"e":"f"},"g":[1,2]}
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.EOF(t, p)
}
