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

	"github.com/katydid/parser-go/parse"
)

func TestParseJustNumber(t *testing.T) {
	s := `1`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidObjectWithValue(t *testing.T) {
	s := `{"a":1`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "a")
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expectErr(t, p.Next)
}

func TestParseInvalidObjectWithOnlyKey(t *testing.T) {
	s := `{"a"`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "a")
	expectErr(t, p.Next)
}

func TestParseInvalidArrayWithComma(t *testing.T) {
	s := `[1,`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expectErr(t, p.Next)
}

func TestParseInvalidArrayWithoutComma(t *testing.T) {
	s := `[1`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expectErr(t, p.Next)
}

func TestParseValidArrayWithSuffixSpace(t *testing.T) {
	s := `[1] `
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expect(t, p.Next, parse.ArrayCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidArrayWithSuffix(t *testing.T) {
	s := `[1] [`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expect(t, p.Next, parse.ArrayCloseHint)
	expectErr(t, p.Next)
}

func TestParseValidObjectWithSuffixSpace(t *testing.T) {
	s := `{"a":1} `
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "a")
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expect(t, p.Next, parse.ObjectCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidObjectWithSuffix(t *testing.T) {
	s := `{"a":1} {`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "a")
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expect(t, p.Next, parse.ObjectCloseHint)
	expectErr(t, p.Next)
}

func TestParseValidNestedObject(t *testing.T) {
	s := `{"a":{"b":{"c":{"d":{"e": 1}}}}}`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "a")
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "b")
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "c")
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "d")
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "e")
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expect(t, p.Next, parse.ObjectCloseHint)
	expect(t, p.Next, parse.ObjectCloseHint)
	expect(t, p.Next, parse.ObjectCloseHint)
	expect(t, p.Next, parse.ObjectCloseHint)
	expect(t, p.Next, parse.ObjectCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidNestedObject(t *testing.T) {
	s := `{"a":{"b":{"c":{"d":{"e": 1}}}}`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "a")
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "b")
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "c")
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "d")
	expect(t, p.Next, parse.ObjectOpenHint)
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "e")
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expect(t, p.Next, parse.ObjectCloseHint)
	expect(t, p.Next, parse.ObjectCloseHint)
	expect(t, p.Next, parse.ObjectCloseHint)
	expect(t, p.Next, parse.ObjectCloseHint)
	expectErr(t, p.Next)
}

func TestParseValidNestedArray2(t *testing.T) {
	s := `[[1]]`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expect(t, p.Next, parse.ArrayCloseHint)
	expect(t, p.Next, parse.ArrayCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseValidNestedArray4å(t *testing.T) {
	s := `[[[[1]]]]`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expect(t, p.Next, parse.ArrayCloseHint)
	expect(t, p.Next, parse.ArrayCloseHint)
	expect(t, p.Next, parse.ArrayCloseHint)
	expect(t, p.Next, parse.ArrayCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestParseInvalidNestedArray(t *testing.T) {
	s := `[[[[1]]]`
	p := NewParser(WithBuffer([]byte(s)))
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ArrayOpenHint)
	expect(t, p.Next, parse.ValueHint)
	expectInt(t, p, 1)
	expect(t, p.Next, parse.ArrayCloseHint)
	expect(t, p.Next, parse.ArrayCloseHint)
	expect(t, p.Next, parse.ArrayCloseHint)
	expectErr(t, p.Next)
}
