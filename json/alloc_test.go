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
	"math/rand"
	"testing"
	"time"

	jsonparser "github.com/katydid/parser-go-json/json"
)

func BenchmarkAlloc(b *testing.B) {
	num := 1000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	js := randJsons(r, num)
	jparser := jsonparser.NewJsonParser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			b.Fatal(err)
		}
		walk(jparser)
	}
	b.ReportAllocs()
}
