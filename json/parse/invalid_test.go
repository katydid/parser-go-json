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
)

func TestParseJustNumber(t *testing.T) {
	s := `1`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectEOF(t, p)
}

func TestParseInvalidObjectWithValue(t *testing.T) {
	s := `{"a":1`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectErr(t, p.Next)
}

func TestParseInvalidObjectWithOnlyKey(t *testing.T) {
	s := `{"a"`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectErr(t, p.Next)
}

func TestParseInvalidArrayWithComma(t *testing.T) {
	s := `[1,`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectErr(t, p.Next)
}

func TestParseInvalidArrayWithoutComma(t *testing.T) {
	s := `[1`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectErr(t, p.Next)
}

func TestParseValidArrayWithSuffixSpace(t *testing.T) {
	s := `[1] `
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestParseInvalidArrayWithSuffix(t *testing.T) {
	s := `[1] [`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectHint(t, p, ArrayCloseHint)
	expectErr(t, p.Next)
}

func TestParseValidObjectWithSuffixSpace(t *testing.T) {
	s := `{"a":1} `
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectHint(t, p, ObjectCloseHint)
	expectEOF(t, p)
}

func TestParseInvalidObjectWithSuffix(t *testing.T) {
	s := `{"a":1} {`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectHint(t, p, ObjectCloseHint)
	expectErr(t, p.Next)
}

func TestParseValidNestedObject(t *testing.T) {
	s := `{"a":{"b":{"c":{"d":{"e": 1}}}}}`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "b")
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "c")
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "d")
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "e")
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectEOF(t, p)
}

func TestParseInvalidNestedObject(t *testing.T) {
	s := `{"a":{"b":{"c":{"d":{"e": 1}}}}`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "b")
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "c")
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "d")
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "e")
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectErr(t, p.Next)
}

func TestParseValidNestedArray2(t *testing.T) {
	s := `[[1]]`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestParseValidNestedArray4å(t *testing.T) {
	s := `[[[[1]]]]`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestParseInvalidNestedArray(t *testing.T) {
	s := `[[[[1]]]`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ValueHint)
	expectInt(t, p, 1)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectErr(t, p.Next)
}
