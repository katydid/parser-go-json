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
	"testing"

	"github.com/katydid/parser-go-json/json/internal/pool"
	"github.com/katydid/parser-go-json/json/rand"
	"github.com/katydid/parser-go/parser/debug"
)

func TestNoAllocsOnAverage(t *testing.T) {
	r := rand.NewRand()
	num := 100
	js := randJSONs(r, num)
	jparser := NewParser()

	const runsPerTest = 100
	checkNoAllocs := func(f func()) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			if allocs := testing.AllocsPerRun(runsPerTest, f); allocs != 0 {
				t.Errorf("seed = %v, got %v allocs, want 0 allocs", r.Seed(), allocs)
			}
		}
	}
	for i := 0; i < num; i++ {
		t.Run(fmt.Sprintf("%d", i), checkNoAllocs(func() {
			if err := jparser.Init(js[i]); err != nil {
				t.Fatalf("seed = %v, err = %v", r.Seed(), err)
			}
			if err := debug.Walk(jparser); err != nil {
				t.Fatalf("seed = %v, err = %v", r.Seed(), err)
			}
		}))
	}
}

func TestNotASingleAllocAfterWarmUp(t *testing.T) {
	r := rand.NewRand()
	num := 100
	js := randJSONs(r, num)
	pool := pool.New()
	jparser := NewParser()
	jparser.(*jsonParser).pool = pool

	// warm up buffer pool
	for i := 0; i < num; i++ {
		if err := jparser.Init(js[i%num]); err != nil {
			t.Fatalf("seed = %v, err = %v value = %v", r.Seed(), err, string(js[i%num]))
		}
		if err := debug.Walk(jparser); err != nil {
			t.Fatalf("seed = %v, err = %v", r.Seed(), err)
		}
	}
	originalPoolSize := pool.Size()

	const runsPerTest = 1
	for i := 0; i < num; i++ {
		f := func() {
			if err := jparser.Init(js[i]); err != nil {
				t.Fatalf("seed = %v, err = %v", r.Seed(), err)
			}
			if err := debug.Walk(jparser); err != nil {
				t.Fatalf("seed = %v, err = %v", r.Seed(), err)
			}
		}
		allocs := testing.AllocsPerRun(runsPerTest, f)
		if allocs != 0 {
			poolallocs := pool.Size() - originalPoolSize
			// there are sometimes allocations made by the testing framework
			// retry to make sure that the allocation is the parser's fault.
			allocs2 := testing.AllocsPerRun(runsPerTest, f)
			if allocs2 != 0 {
				t.Errorf("seed = %v, got %v allocs, want 0 allocs, pool allocs = %v", r.Seed(), allocs, poolallocs)
			}
		}
	}
}

func BenchmarkAlloc(b *testing.B) {
	num := 1000
	r := rand.NewRand()
	js := randJSONs(r, num)
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

func BenchmarkPoolDefault(b *testing.B) {
	// generate random jsons
	num := 1000
	r := rand.NewRand()
	js := randJSONs(r, num)

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
	js := randJSONs(r, num)

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
