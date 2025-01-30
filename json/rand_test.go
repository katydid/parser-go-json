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

package json_test

import (
	"encoding/json"
	"math/rand"
	"strings"
)

func randJsons(r *rand.Rand, num int) [][]byte {
	js := make([][]byte, num)
	for i := 0; i < num; i++ {
		js[i] = randJson(r)
	}
	return js
}

func randJson(r *rand.Rand) []byte {
	val := randValue(r, 5)
	data, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	return data
}

func randObject(r *rand.Rand, level int) map[string]interface{} {
	l := r.Intn(10)
	ms := make(map[string]interface{})
	for i := 0; i < l; i++ {
		ms[randName(r)] = randValue(r, level)
	}
	return ms
}

func randArray(r *rand.Rand, level int) []interface{} {
	l := r.Intn(10)
	as := make([]interface{}, l)
	for i := 0; i < l; i++ {
		as[i] = randValue(r, level)
	}
	return as
}

func randValue(r *rand.Rand, level int) interface{} {
	maxN := 9
	if level <= 0 {
		// do not generate arrays or objects,
		// since we have generated a deep enough structure and
		// we do not want to endlessly recurse.
		maxN = 7
	}
	switch r.Intn(maxN) {
	case 0:
		return nil
	case 1:
		return bool(r.Intn(2) == 0)
	case 2:
		return int64(r.Int63())
	case 3:
		return uint64(r.Uint64())
	case 4:
		return float64(r.Float64())
	case 5:
		return randString(r)
	case 6:
		return randBytes(r)
	case 7:
		return randArray(r, level-1)
	case 8:
		return randObject(r, level-1)
	}
	panic("unreachable")
}

func randBytes(r *rand.Rand) []byte {
	bs := make([]byte, int(r.Intn(100)))
	for i := range bs {
		bs[i] = byte(r.Intn(255))
	}
	return bs
}

// string := '"' characters '"'
// characters := "" | character characters
func randString(r *rand.Rand) string {
	ss := make([]string, int(r.Intn(100)))
	for i := range ss {
		ss[i] = randChar(r)
	}
	s := strings.Join(ss, "")
	return s
}

// character := '0020' . '10FFFF' - '"' - '\' | '\' escape
func randChar(r *rand.Rand) string {
	switch r.Intn(2) {
	case 0:
		max := int('\U0010FFFF') + 1
		min := int('\u0020')
		ran := int((max - min) - 2)
		random := rune(r.Intn(ran) + min)
		if random >= '"' {
			random += 1
		}
		if random >= '\\' {
			random += 1
		}
		return string([]rune{random})
	case 1:
		return string([]rune{randEscape(r)})
	}
	panic("unreachable")
}

// escape := '"' | '\' | '/' | 'b' | 'f' | 'n' | 'r' | 't' | 'u' hex hex hex hex
// hex := digit | 'A' . 'F' | 'a' . 'f'
func randEscape(r *rand.Rand) rune {
	switch r.Intn(9) {
	case 0:
		return '"'
	case 1:
		return '\\'
	case 2:
		return '/'
	case 3:
		return '\b'
	case 4:
		return '\f'
	case 5:
		return '\n'
	case 6:
		return '\r'
	case 7:
		return '\t'
	case 8:
		// TODO generate better \u characters
		return '\u01aB'
	}
	panic("unreachable")
}

func randName(r *rand.Rand) string {
	ss := make([]byte, int(r.Intn(100)))
	for i := range ss {
		ss[i] = byte(65 + r.Intn(26))
	}
	return string(ss)
}
