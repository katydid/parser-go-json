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
	"math/rand"
	"testing"
	"time"

	"github.com/katydid/parser-go-json/json/internal/pool"
)

func BenchmarkPoolDefault(b *testing.B) {
	seed := time.Now().UnixNano()
	// generate random jsons
	num := 1000
	r := rand.New(rand.NewSource(seed))
	js := randJsons(r, num)

	// initialise pool
	jparser := NewJsonParser()

	// exercise buffer pool
	for i := 0; i < num; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			b.Fatalf("seed = %v, err = %v", seed, err)
		}
		walk(jparser)
	}
	// start benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			b.Fatalf("seed = %v, err = %v", seed, err)
		}
		walk(jparser)
	}
	b.ReportAllocs()
}

func BenchmarkPoolNone(b *testing.B) {
	seed := time.Now().UnixNano()
	// generate random jsons
	num := 1000
	r := rand.New(rand.NewSource(seed))
	js := randJsons(r, num)

	// set pool to no pool
	jparser := NewJsonParser()
	jparser.(*jsonParser).pool = pool.None()

	// start benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			b.Fatalf("seed = %v, err = %v", seed, err)
		}
		walk(jparser)
	}
	b.ReportAllocs()
}
