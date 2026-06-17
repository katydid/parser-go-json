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

package unquote

import (
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

var u4table = [256]rune{}

func init() {
	for i := range 256 {
		u4table[i] = rune(-1)
	}
	u4table['0'] = rune(0)
	u4table['1'] = rune(1)
	u4table['2'] = rune(2)
	u4table['3'] = rune(3)
	u4table['4'] = rune(4)
	u4table['5'] = rune(5)
	u4table['6'] = rune(6)
	u4table['7'] = rune(7)
	u4table['8'] = rune(8)
	u4table['9'] = rune(9)
	u4table['a'] = rune(10)
	u4table['b'] = rune(11)
	u4table['c'] = rune(12)
	u4table['d'] = rune(13)
	u4table['e'] = rune(14)
	u4table['f'] = rune(15)
	u4table['A'] = rune(10)
	u4table['B'] = rune(11)
	u4table['C'] = rune(12)
	u4table['D'] = rune(13)
	u4table['E'] = rune(14)
	u4table['F'] = rune(15)
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	var r rune
	for _, c := range s[2:6] {
		r1 := u4table[c]
		if r1 == -1 {
			return -1
		}
		r = r*16 + rune(r1)
	}
	return r
}

var unusualTable = [256]byte{}

func init() {
	for i := range 256 {
		if i < ' ' {
			unusualTable[i] = 1
		}
		if i > utf8.RuneSelf {
			unusualTable[i] = 1
		}
	}
	unusualTable['\\'] = 1
	unusualTable['"'] = 1
}

var backslashTable = [256]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'\'': '\'',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'u':  1,
}

func unquoteBytes(alloc func(size int) []byte, s []byte) ([]byte, int, bool) {
	if len(s) < 2 || s[0] != '"' {
		return nil, 0, false
	}

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 1
	for r < len(s) {
		if unusualTable[s[r]] == 0 {
			r++
			continue
		} else {
			break
		}
	}
	if r >= len(s) {
		return nil, 0, false
	}
	if s[r] == '"' {
		return s[1:r], r + 1, true
	}

	b := alloc(len(s) * utf8.UTFMax)
	w := copy(b, s[1:r])
	for r < len(s) {
		c := s[r]
		switch c {
		case '\\':
			r++
			if r >= len(s) {
				return nil, 0, false
			}
			b[w] = backslashTable[s[r]]
			switch b[w] {
			case 0:
				return nil, 0, false
			case 1:
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return nil, 0, false
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			default:
				r++
				w++
			}
		case '"':
			return b[0:w], r + 1, true

		default:
			if c < ' ' {
				// Quote, control characters are invalid.
				return nil, 0, false
			} else if c < utf8.RuneSelf {
				// ASCII
				b[w] = c
				r++
				w++
			} else {
				// Coerce to well-formed UTF-8.
				rr, size := utf8.DecodeRune(s[r:])
				r += size
				w += utf8.EncodeRune(b[w:], rr)
			}
		}
	}
	return nil, 0, false
}
