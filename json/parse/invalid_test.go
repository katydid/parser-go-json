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
)

func TestParseJustNumber(t *testing.T) {
	s := `1`
	p := NewParser([]byte(s))
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidObjectWithValue(t *testing.T) {
	s := `{"a":1`
	p := NewParser([]byte(s))
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "a")
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expectErr(t, p.Next)
}

func TestParseInvalidObjectWithOnlyKey(t *testing.T) {
	s := `{"a"`
	p := NewParser([]byte(s))
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "a")
	expectErr(t, p.Next)
}

func TestParseInvalidArrayWithComma(t *testing.T) {
	s := `[1,`
	p := NewParser([]byte(s))
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expectErr(t, p.Next)
}

func TestParseInvalidArrayWithoutComma(t *testing.T) {
	s := `[1`
	p := NewParser([]byte(s))
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expectErr(t, p.Next)
}

func TestParseValidArrayWithSuffixSpace(t *testing.T) {
	s := `[1] `
	p := NewParser([]byte(s))
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expect(t, p.Next, ArrayCloseKind)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidArrayWithSuffix(t *testing.T) {
	s := `[1] [`
	p := NewParser([]byte(s))
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expect(t, p.Next, ArrayCloseKind)
	expectErr(t, p.Next)
}

func TestParseValidObjectWithSuffixSpace(t *testing.T) {
	s := `{"a":1} `
	p := NewParser([]byte(s))
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "a")
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expect(t, p.Next, ObjectCloseKind)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidObjectWithSuffix(t *testing.T) {
	s := `{"a":1} {`
	p := NewParser([]byte(s))
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "a")
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expect(t, p.Next, ObjectCloseKind)
	expectErr(t, p.Next)
}

func TestParseValidNestedObject(t *testing.T) {
	s := `{"a":{"b":{"c":{"d":{"e": 1}}}}}`
	p := NewParser([]byte(s))
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "a")
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "b")
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "c")
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "d")
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "e")
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expect(t, p.Next, ObjectCloseKind)
	expect(t, p.Next, ObjectCloseKind)
	expect(t, p.Next, ObjectCloseKind)
	expect(t, p.Next, ObjectCloseKind)
	expect(t, p.Next, ObjectCloseKind)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidNestedObject(t *testing.T) {
	s := `{"a":{"b":{"c":{"d":{"e": 1}}}}`
	p := NewParser([]byte(s))
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "a")
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "b")
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "c")
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "d")
	expect(t, p.Next, ObjectOpenKind)
	expect(t, p.Next, StringKind)
	expect(t, p.String, "e")
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expect(t, p.Next, ObjectCloseKind)
	expect(t, p.Next, ObjectCloseKind)
	expect(t, p.Next, ObjectCloseKind)
	expect(t, p.Next, ObjectCloseKind)
	expectErr(t, p.Next)
}

func TestParseValidNestedArray(t *testing.T) {
	s := `[[[[1]]]]`
	p := NewParser([]byte(s))
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expect(t, p.Next, ArrayCloseKind)
	expect(t, p.Next, ArrayCloseKind)
	expect(t, p.Next, ArrayCloseKind)
	expect(t, p.Next, ArrayCloseKind)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidNestedArray(t *testing.T) {
	s := `[[[[1]]]`
	p := NewParser([]byte(s))
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, ArrayOpenKind)
	expect(t, p.Next, NumberKind)
	expect(t, p.Int, 1)
	expect(t, p.Next, ArrayCloseKind)
	expect(t, p.Next, ArrayCloseKind)
	expect(t, p.Next, ArrayCloseKind)
	expectErr(t, p.Next)
}
