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

	"github.com/katydid/parser-go-json/json/internal/expect"
	"github.com/katydid/parser-go/parse"
)

func TestParseJustNumber(t *testing.T) {
	s := `1`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.EOF(t, p)
}

func TestParseInvalidObjectWithValue(t *testing.T) {
	s := `{"a":1`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Err(t, p.Next)
}

func TestParseInvalidObjectWithOnlyKey(t *testing.T) {
	s := `{"a"`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Err(t, p.Next)
}

func TestParseInvalidArrayWithComma(t *testing.T) {
	s := `[1,`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Err(t, p.Next)
}

func TestParseInvalidArrayWithoutComma(t *testing.T) {
	s := `[1`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Err(t, p.Next)
}

func TestParseValidArrayWithSuffixSpace(t *testing.T) {
	s := `[1] `
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestParseInvalidArrayWithSuffix(t *testing.T) {
	s := `[1] [`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Err(t, p.Next)
}

func TestParseValidObjectWithSuffixSpace(t *testing.T) {
	s := `{"a":1} `
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.EOF(t, p)
}

func TestParseInvalidObjectWithSuffix(t *testing.T) {
	s := `{"a":1} {`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Err(t, p.Next)
}

func TestParseValidNestedObject(t *testing.T) {
	s := `{"a":{"b":{"c":{"d":{"e": 1}}}}}`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "c")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "d")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "e")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.EOF(t, p)
}

func TestParseInvalidNestedObject(t *testing.T) {
	s := `{"a":{"b":{"c":{"d":{"e": 1}}}}`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "c")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "d")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "e")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Err(t, p.Next)
}

func TestParseValidNestedArray2(t *testing.T) {
	s := `[[1]]`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestParseValidNestedArray4å(t *testing.T) {
	s := `[[[[1]]]]`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestParseInvalidNestedArray(t *testing.T) {
	s := `[[[[1]]]`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Err(t, p.Next)
}
