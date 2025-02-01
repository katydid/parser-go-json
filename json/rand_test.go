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

package json

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestParseRandom(t *testing.T) {
	seed := time.Now().UnixNano()
	seed = 1738425704510828000 // TODO remove this line
	num := 10000
	r := rand.New(rand.NewSource(seed))
	js := randJsons(r, num)
	jparser := NewJsonParser()

	// warm up buffer pool
	for i := 0; i < num; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			t.Fatalf("seed = %v, err = %v, input = %v", seed, err, string(js[i%num]))
		}
		walk(jparser)
	}
}

func randJsons(r *rand.Rand, num int) [][]byte {
	js := make([][]byte, num)
	for i := 0; i < num; i++ {
		js[i] = randJson(r)
	}
	return js
}

func randJson(r *rand.Rand) []byte {
	val := randValue(r, 5)
	return []byte(val)
}

// value := object | array | string | number | "true" | "false" | "null"
func randValue(r *rand.Rand, level int) string {
	maxN := 6
	if level <= 0 {
		// do not generate arrays or objects,
		// since we have generated a deep enough structure and
		// we do not want to endlessly recurse.
		maxN = 4
	}
	switch r.Intn(maxN) {
	case 0:
		return "null"
	case 1:
		return "false"
	case 2:
		return "true"
	case 3:
		return randNumber(r)
	case 4:
		return randString(r)
	case 5:
		return randArray(r, level-1)
	case 6:
		return randObject(r, level-1)
	}
	panic("unreachable")
}

// object := '{' ws '}' | '{' members '}'
// members := member | member ',' members
// member := ws string ws ':' element
func randObject(r *rand.Rand, level int) string {
	l := r.Intn(10)
	if l == 0 {
		return "{" + randWs(r) + "}"
	}
	ss := make([]string, l)
	for i := 0; i < l; i++ {
		ss[i] = randWs(r) + randString(r) + randWs(r) + ":" + randElement(r, level)
	}
	return "{" + strings.Join(ss, ",") + "}"
}

// array := '[' ws ']' | '[' elements ']'
// elements := element | element ',' elements
func randArray(r *rand.Rand, level int) string {
	l := r.Intn(10)
	if l == 0 {
		return "[" + randWs(r) + "]"
	}
	ss := make([]string, l)
	for i := 0; i < l; i++ {
		ss[i] = randElement(r, level)
	}
	return "[" + strings.Join(ss, ",") + "]"
}

// element := ws value ws
func randElement(r *rand.Rand, level int) string {
	return randWs(r) + randValue(r, level) + randWs(r)
}

// string := '"' characters '"'
// characters := "" | character characters
func randString(r *rand.Rand) string {
	ss := make([]string, int(r.Intn(100)))
	for i := range ss {
		ss[i] = randChar(r)
	}
	s := "\"" + strings.Join(ss, "") + "\""
	return s
}

// character := '0020' . '10FFFF' - '"' - '\' | '\' escape
func randChar(r *rand.Rand) string {
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
func randEscape(r *rand.Rand) string {
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

// number := integer fraction exponent
func randNumber(r *rand.Rand) string {
	return randInteger(r) + randFraction(r) + randExponent(r)
}

// integer := digit | onenine digits | '-' digit | '-' onenine digits
func randInteger(r *rand.Rand) string {
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
func randExponent(r *rand.Rand) string {
	switch r.Intn(3) {
	case 0:
		return ""
	case 1:
		return "E" + randSign(r) + randDigits(r)
	case 2:
		return "3" + randSign(r) + randDigits(r)
	}
	panic("unreachable")
}

// fraction := "" | '.' digits
func randFraction(r *rand.Rand) string {
	switch r.Intn(2) {
	case 0:
		return ""
	case 1:
		return "." + randDigits(r)
	}
	panic("unreachable")
}

// sign := "" | '+' | '-'
func randSign(r *rand.Rand) string {
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
func randDigits(r *rand.Rand) string {
	l := r.Intn(5) + 1
	ss := make([]string, l)
	for i := 0; i < l; i++ {
		ss[i] = randDigit(r)
	}
	return strings.Join(ss, "")
}

// digit := '0' | onenine
func randDigit(r *rand.Rand) string {
	return fmt.Sprintf("%d", r.Intn(10))
}

// onenine := '1' . '9'
func randOneNine(r *rand.Rand) string {
	return fmt.Sprintf("%d", r.Intn(9)+1)
}

// hex := digit | 'A' . 'F' | 'a' . 'f'
func randHex(r *rand.Rand) string {
	s := "01234567890abcdefABCDEF"
	return string([]rune{rune(s[r.Intn(len(s))])})
}

// ws := "" | '0020' ws | '000A' ws | '000D' ws | '0009' ws
func randWs(r *rand.Rand) string {
	l := r.Intn(5)
	ss := make([]rune, l)
	for i := 0; i < l; i++ {
		ss[i] = randW(r)
	}
	return string(ss)
}

func randW(r *rand.Rand) rune {
	switch r.Intn(4) {
	case 0:
		return '\u0020'
	case 1:
		return '\u000A'
	case 2:
		return '\u000D'
	case 3:
		return '\u0009'
	}
	panic("unreachable")
}
