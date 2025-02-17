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

package rand

import (
	"fmt"
	"strings"
)

// Value returns a string representing random json value.
func Value(r Rand, opts ...Option) string {
	c := newConfig(opts...)
	return randValue(r, c)
}

// value BNF:
// value := object | array | string | number | "true" | "false" | "null"
func randValue(r Rand, c *config) string {
	if c.maxDepth <= 0 {
		// do not generate arrays or objects,
		// since we have generated a deep enough structure and
		// we do not want to endlessly recurse.
		return randTerminator(r, c)
	}
	switch r.Intn(3) {
	case 0:
		return randTerminator(r, c)
	case 1:
		c.maxDepth = c.maxDepth - 1
		return randArray(r, c)
	case 2:
		c.maxDepth = c.maxDepth - 1
		return randObject(r, c)
	}
	panic("unreachable")
}

// randTerminator generates a random json value that cannot recurse into generating more json values.
// More specifically, it generates either null, false, true, number or a string and does not generate an array or object.
func randTerminator(r Rand, c *config) string {
	switch r.Intn(5) {
	case 0:
		return "null"
	case 1:
		return "false"
	case 2:
		return "true"
	case 3:
		return randNumber(r, c)
	case 4:
		return randString(r, c)
	}
	panic("unreachable")
}

// Object returns a string that represents a random JSON object.
func Object(r Rand, opts ...Option) string {
	c := newConfig(opts...)
	return randObject(r, c)
}

// object BNF:
// object := '{' ws '}' | '{' members '}'
// members := member | member ',' members
// member := ws string ws ':' element
func randObject(r Rand, c *config) string {
	l := r.Intn(c.maxObjectFields)
	if l == 0 {
		return "{" + randWs(r) + "}"
	}
	ss := make([]string, l)
	for i := 0; i < l; i++ {
		ss[i] = randWs(r) + String(r) + randWs(r) + ":" + randElement(r, c)
	}
	return "{" + strings.Join(ss, ",") + "}"
}

// Array returns a string that represents a random JSON array.
func Array(r Rand, opts ...Option) string {
	c := newConfig(opts...)
	return randArray(r, c)
}

// array := '[' ws ']' | '[' elements ']'
// elements := element | element ',' elements
func randArray(r Rand, c *config) string {
	l := r.Intn(c.maxArrayLength)
	if l == 0 {
		return "[" + randWs(r) + "]"
	}
	ss := make([]string, l)
	for i := 0; i < l; i++ {
		ss[i] = randElement(r, c)
	}
	return "[" + strings.Join(ss, ",") + "]"
}

// element := ws value ws
func randElement(r Rand, c *config) string {
	return randWs(r) + randValue(r, c) + randWs(r)
}

// String returns a string that represents a random JSON string.
func String(r Rand, opts ...Option) string {
	c := newConfig(opts...)
	return randString(r, c)
}

// String BNF:
// string := '"' characters '"'
// characters := "" | character characters
func randString(r Rand, c *config) string {
	ss := make([]string, int(r.Intn(c.maxStringLength)))
	for i := range ss {
		ss[i] = randChar(r)
	}
	s := "\"" + strings.Join(ss, "") + "\""
	return s
}

// character := '0020' . '10FFFF' - '"' - '\' | '\' escape
func randChar(r Rand) string {
	switch r.Intn(2) {
	case 0:
		min := int('\u0020')
		max := int('\U0010FFFF') + 1
		ran := int((max - min) - 2)
		random := rune(r.Intn(ran) + min)
		if random != '"' && random != '\\' {
			return string([]rune{random})
		}
		return randChar(r)
	case 1:
		return "\\" + randEscape(r)
	}
	panic("unreachable")
}

// escape := '"' | '\' | '/' | 'b' | 'f' | 'n' | 'r' | 't' | 'u' hex hex hex hex
func randEscape(r Rand) string {
	switch r.Intn(9) {
	case 0:
		return "\""
	case 1:
		return "\\"
	case 2:
		return "/"
	case 3:
		return "b"
	case 4:
		return "f"
	case 5:
		return "n"
	case 6:
		return "r"
	case 7:
		return "t"
	case 8:
		return "u" + randHex(r) + randHex(r) + randHex(r) + randHex(r)
	}
	panic("unreachable")
}

// Number returns a string that represents a random JSON number.
func Number(r Rand, opts ...Option) string {
	c := newConfig(opts...)
	return randNumber(r, c)
}

// number BNF:
// number := integer fraction exponent
func randNumber(r Rand, c *config) string {
	// Sometimes generate an edge case
	if r.Intn(c.numberEdgeCaseOdds) == 0 {
		return c.numberEdgeCases[r.Intn(len(c.numberEdgeCases))]
	}
	return randInteger(r) + randFraction(r) + randExponent(r)
}

// integer := digit | onenine digits | '-' digit | '-' onenine digits
func randInteger(r Rand) string {
	switch r.Intn(4) {
	case 0:
		return randDigit(r)
	case 1:
		return randOneNine(r) + randDigits(r)
	case 2:
		return "-" + randDigit(r)
	case 3:
		return "-" + randOneNine(r) + randDigits(r)
	}
	panic("unreachable")
}

// exponent := "" | 'E' sign digits | 'e' sign digits
func randExponent(r Rand) string {
	switch r.Intn(3) {
	case 0:
		return ""
	case 1:
		return "E" + randSign(r) + randDigits(r)
	case 2:
		return "e" + randSign(r) + randDigits(r)
	}
	panic("unreachable")
}

// fraction := "" | '.' digits
func randFraction(r Rand) string {
	switch r.Intn(2) {
	case 0:
		return ""
	case 1:
		return "." + randDigits(r)
	}
	panic("unreachable")
}

// sign := "" | '+' | '-'
func randSign(r Rand) string {
	switch r.Intn(3) {
	case 0:
		return ""
	case 1:
		return "+"
	case 2:
		return "-"
	}
	panic("unreachable")
}

// digits := digit | digit digits
func randDigits(r Rand) string {
	l := r.Intn(5) + 1
	ss := make([]string, l)
	for i := 0; i < l; i++ {
		ss[i] = randDigit(r)
	}
	return strings.Join(ss, "")
}

// digit := '0' | onenine
func randDigit(r Rand) string {
	return fmt.Sprintf("%d", r.Intn(10))
}

// onenine := '1' . '9'
func randOneNine(r Rand) string {
	return fmt.Sprintf("%d", r.Intn(9)+1)
}

// hex := digit | 'A' . 'F' | 'a' . 'f'
func randHex(r Rand) string {
	s := "01234567890abcdefABCDEF"
	return string([]rune{rune(s[r.Intn(len(s))])})
}

// ws := "" | '0020' ws | '000A' ws | '000D' ws | '0009' ws
func randWs(r Rand) string {
	l := r.Intn(5)
	ss := make([]rune, l)
	for i := 0; i < l; i++ {
		ss[i] = randW(r)
	}
	return string(ss)
}

func randW(r Rand) rune {
	switch r.Intn(4) {
	case 0:
		return '\u0020' // space
	case 1:
		return '\u000A' // \n
	case 2:
		return '\u000D' // \r
	case 3:
		return '\u0009' // \t
	}
	panic("unreachable")
}
