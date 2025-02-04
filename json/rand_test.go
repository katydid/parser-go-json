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
	"testing"

	"github.com/katydid/parser-go-json/json/rand"
	"github.com/katydid/parser-go/parser/debug"
)

func randJSONs(r rand.Rand, num int) [][]byte {
	js := make([][]byte, num)
	for i := 0; i < num; i++ {
		js[i] = randJSON(r)
	}
	return js
}

func randJSON(r rand.Rand) []byte {
	val := rand.Value(r, 5)
	return []byte(val)
}

func TestParseRandom(t *testing.T) {
	r := rand.NewRand()
	num := 10000
	js := randJSONs(r, num)
	jparser := NewParser()

	// warm up buffer pool
	for i := 0; i < num; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			t.Fatalf("seed = %v, err = %v, input = %v", r.Seed(), err, string(js[i%num]))
		}
		if err := debug.Walk(jparser); err != nil {
			t.Fatalf("seed = %v, err = %v", r.Seed(), err)
		}
	}
}
