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

	"github.com/katydid/parser-go-json/json/internal/pool"
	"github.com/katydid/parser-go-json/json/rand"
	"github.com/katydid/parser-go/parser/debug"
)

func BenchmarkPoolDefault(b *testing.B) {
	// generate random jsons
	num := 1000
	r := rand.NewRand()
	js := rand.Values(r, num)

	// initialise pool
	jparser := NewParser()

	// exercise buffer pool
	for i := 0; i < num; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			b.Fatalf("seed = %v, err = %v", r.Seed(), err)
		}
		if err := debug.Walk(jparser); err != nil {
			b.Fatalf("seed = %v, err = %v", r.Seed(), err)
		}
	}
	// start benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			b.Fatalf("seed = %v, err = %v", r.Seed(), err)
		}
		if err := debug.Walk(jparser); err != nil {
			b.Fatalf("seed = %v, err = %v", r.Seed(), err)
		}
	}
	b.ReportAllocs()
}

func BenchmarkPoolNone(b *testing.B) {
	// generate random jsons
	num := 1000
	r := rand.NewRand()
	js := rand.Values(r, num)

	// set pool to no pool
	jparser := NewParser()
	jparser.(*jsonParser).pool = pool.None()

	// start benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			b.Fatalf("seed = %v, err = %v", r.Seed(), err)
		}
		if err := debug.Walk(jparser); err != nil {
			b.Fatalf("seed = %v, err = %v", r.Seed(), err)
		}
	}
	b.ReportAllocs()
}
