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

// Generate generates many Values, at least one of each kind.
// atleast specifies the minimum number to generate.
// The kinds are object, array, number, string, true, false and null.
func Values(r Rand, atleast int) [][]byte {
	_, vs := values(r, atleast)
	bs := make([][]byte, len(vs))
	for i := range vs {
		bs[i] = []byte(vs[i])
	}
	return bs
}

func values(r Rand, atleast int) (map[byte]int, []string) {
	// kinds is used to track that at least all kinds were generated
	kinds := allKinds()
	ss := []string{}
	for !(eachKindHasNonZero(kinds) && len(ss) >= atleast) {
		s := Value(r, 5)
		ss = append(ss, s)
		kinds[whichKind(s)] += 1
	}
	return kinds, ss
}

func allKinds() map[byte]int {
	return map[byte]int{
		'{': 0,
		'[': 0,
		'0': 0,
		'"': 0,
		't': 0,
		'f': 0,
		'n': 0,
	}
}

func eachKindHasNonZero(kinds map[byte]int) bool {
	for _, v := range kinds {
		if v == 0 {
			return false
		}
	}
	return true
}

func whichKind(s string) byte {
	bs := []byte(s)
	c := bs[0]
	switch c {
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return '0'
	case '{', '[', '"', 't', 'f', 'n':
		return c
	}
	panic("unreachable")
}
