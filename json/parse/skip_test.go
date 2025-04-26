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
	"io"
	"testing"

	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

// Skip returns an error if nothing has been parsed yet.
func TestSkipUnknownObject(t *testing.T) {
	str := `{}`
	p := NewParser(WithBuffer([]byte(str)))
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	expect.Hint(t, p, parse.ObjectCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestSkipUnknownArray(t *testing.T) {
	str := `[]`
	p := NewParser(WithBuffer([]byte(str)))
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	expect.Hint(t, p, parse.ArrayCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestSkipUnknownArrayOpen(t *testing.T) {
	str := `[1]`
	p := NewParser(WithBuffer([]byte(str)))
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ArrayCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestSkipUnknownString(t *testing.T) {
	str := `"abc"`
	p := NewParser(WithBuffer([]byte(str)))
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

// If the kind '[' was returned by Next, then the whole array is skipped.
func TestSkipArrayOpen(t *testing.T) {
	str := `[1,2]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over 1,2]
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestSkipArrayNestedOpen(t *testing.T) {
	str := `[[1,2]]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayOpenHint)
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over 1,2]
	expect.Hint(t, p, parse.ArrayCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

// If an array element was parsed, then the rest of the array is skipped.
func TestSkipArrayElement(t *testing.T) {
	str := `[1,2,3]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over 2,3]
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
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
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over 3,4]
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 5)
	expect.Hint(t, p, parse.ArrayCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestSkipArrayRecursiveElement(t *testing.T) {
	str := `[1,[2,3],[[4,5,6]]]`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over [2,3],[[4,5,6]]]
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

// If the kind '{' was returned by Next, then the whole object is skipped.
func TestSkipObjectOpen(t *testing.T) {
	str := `{"a":1,"b":2}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over "a":1,"b":2}
	if kind, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got kind %v, %v", kind, err)
	}
}

func TestSkipObjectNestedOpen(t *testing.T) {
	str := `{"a":{"b":1,"c":2}}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ObjectOpenHint)
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over "b":1,"c":2}
	expect.Hint(t, p, parse.ObjectCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
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
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over "b":2,"c":3}
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
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
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over "c":2,"d":3}
	expect.Hint(t, p, parse.ObjectCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

// If a object key was just parsed, then that key's value is skipped.
func TestSkipObjectValue(t *testing.T) {
	str := `{"a":1,"b":2}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over 2
	expect.Hint(t, p, parse.ObjectCloseHint)
	if kind, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v with kind %v", err, kind)
	}
}

func TestSkipObjectRecursiveValue(t *testing.T) {
	str := `{"a":1,"b":{"c":{"d":{"e":"f"},"g":[1,2]}}}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over {"c":{"d":{"e":"f"},"g":[1,2]}}
	expect.Hint(t, p, parse.ObjectCloseHint)
	if kind, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v with kind %v", err, kind)
	}
}

func TestSkipObjectDeepRecursiveValue(t *testing.T) {
	str := `{"a":1,"b":{"c":{"d":{"e":"f"},"g":[1,2]}}}`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")

	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "c")
	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}
	// skipped over {"d":{"e":"f"},"g":[1,2]}
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	if kind, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v, with kind %v", err, kind)
	}
}
