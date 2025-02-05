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

package token

import (
	"testing"

	"github.com/katydid/parser-go-json/json/internal/pool"
	"github.com/katydid/parser-go-json/json/rand"
)

func TestNoAllocsOnAverage(t *testing.T) {
	r := rand.NewRand()
	num := 100
	js := rand.Values(r, num)

	pool := pool.New()
	tzer := NewTokenizerWithCustomAllocator(nil, pool.Alloc)

	const runsPerTest = 100
	for i := 0; i < num; i++ {
		f := func() {
			tzer.Init(js[i])
			if err := walk(tzer); err != nil {
				t.Fatalf("expected EOF, but got %v", err)
			}
			pool.FreeAll()
		}
		allocs := testing.AllocsPerRun(runsPerTest, f)
		if allocs != 0 {
			t.Errorf("seed = %v, got %v allocs, want 0 allocs", r.Seed(), allocs)
		}
	}
}

func TestNotASingleAllocAfterWarmUp(t *testing.T) {
	r := rand.NewRand()
	num := 100
	js := rand.Values(r, num)

	pool := pool.New()
	tzer := NewTokenizerWithCustomAllocator(nil, pool.Alloc)

	// warm up buffer pool
	for i := 0; i < num; i++ {
		tzer.Init(js[i])
		if err := walk(tzer); err != nil {
			t.Fatalf("expected EOF, but got %v", err)
		}
		pool.FreeAll()
	}
	originalPoolSize := pool.Size()

	const runsPerTest = 1
	for i := 0; i < num; i++ {
		f := func() {
			tzer.Init(js[i])
			if err := walk(tzer); err != nil {
				t.Fatalf("expected EOF, but got %v", err)
			}
			pool.FreeAll()
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
